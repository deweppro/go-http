/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package epoll

import (
	"context"
	"io"
	"net"
	"time"

	serv "github.com/deweppro/go-http"
	"github.com/deweppro/go-logger"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

//Server ...
type Server struct {
	conf     ConfigItem
	handler  ConnectionHandler
	listener net.Listener
	epoll    *Epoll
	ctx      context.Context
	cncl     context.CancelFunc
	log      logger.Logger
}

//NewServer ...
func NewServer(conf EpollConfig, log logger.Logger) *Server {
	return NewCustomServer(conf.Epoll, log)
}

//NewCustomServer ...
func NewCustomServer(conf ConfigItem, log logger.Logger) *Server {
	return &Server{
		conf: conf,
		log:  log,
	}
}

//Handler ...
func (s *Server) Handler(h ConnectionHandler) {
	s.handler = h
}

//Up ...
func (s *Server) Up() (err error) {
	if s.listener != nil {
		err = serv.ErrServAlreadyRunning
		return
	}
	s.ctx, s.cncl = context.WithCancel(context.Background())
	if len(s.conf.Addr) == 0 {
		if s.conf.Addr, err = serv.RandomPort("localhost"); err != nil {
			return
		}
	}
	s.listener, err = net.Listen("tcp", s.conf.Addr)
	if err != nil {
		return
	}
	if s.epoll, err = NewEpoll(s.log); err != nil {
		return
	}
	s.log.Infof("tcp server started on %s", s.conf.Addr)
	go s.connAccept()
	go s.epollAccept()
	return nil
}

func (s *Server) connAccept() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				return
			default:
				s.log.Errorf("epoll conn accept: %s", err.Error())

				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					time.Sleep(1 * time.Second)
					continue
				}
				return
			}
		}

		if err := s.epoll.AddOrClose(conn); err != nil {
			s.log.Errorf("epoll add conn: %s", err.Error())
		}
	}
}

func (s *Server) epollAccept() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			clist, err := s.epoll.Wait()
			switch err {
			case nil:
			case serv.ErrEpollEmptyEvents:
				continue
			case unix.EINTR:
				continue
			default:
				s.log.Errorf("epoll accept conn: %s", err.Error())
				continue
			}

			for _, c := range clist {
				go func(conn *netConnItem) {
					defer conn.Await(false)
					er := connection(conn.Conn, s.handler)
					switch er {
					case nil:
					case io.EOF:
						_ = s.epoll.Close(conn)
					default:
						s.log.Errorf("epoll bad conn from %s: %s", conn.Conn.RemoteAddr().String(), err.Error())
						_ = s.epoll.Close(conn)
					}
				}(c)
			}
		}
	}
}

//Down ...
func (s *Server) Down() (err error) {
	if s.listener == nil {
		return
	}
	defer func() {
		s.listener = nil
	}()
	s.cncl()
	if e := s.epoll.CloseAll(); e != nil {
		err = errors.Wrap(err, e.Error())
	}
	if e := s.listener.Close(); e != nil {
		err = errors.Wrap(err, e.Error())
	}
	s.log.Infof("tcp server stopped on %s", s.conf.Addr)
	return
}
