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

type Server struct {
	conf     ConfigItem
	handler  ConnectionHandler
	listener net.Listener
	epoll    *Epoll
	ctx      context.Context
	cncl     context.CancelFunc
	log      logger.Logger
}

func NewServer(conf EpollConfig, log logger.Logger) *Server {
	return NewCustomServer(conf.Epoll, log)
}

func NewCustomServer(conf ConfigItem, log logger.Logger) *Server {
	return &Server{
		conf: conf,
		log:  log,
	}
}

func (_srv *Server) Handler(h ConnectionHandler) {
	_srv.handler = h
}

func (_srv *Server) Up() (err error) {
	if _srv.listener != nil {
		err = serv.ErrServAlreadyRunning
		return
	}
	_srv.ctx, _srv.cncl = context.WithCancel(context.Background())
	if len(_srv.conf.Addr) == 0 {
		if _srv.conf.Addr, err = serv.RandomPort("localhost"); err != nil {
			return
		}
	}
	_srv.listener, err = net.Listen("tcp", _srv.conf.Addr)
	if err != nil {
		return
	}
	if _srv.epoll, err = NewEpoll(_srv.log); err != nil {
		return
	}
	_srv.log.Infof("tcp server started on %s", _srv.conf.Addr)
	go _srv.connAccept()
	go _srv.epollAccept()
	return nil
}

func (_srv *Server) connAccept() {
	for {
		conn, err := _srv.listener.Accept()
		if err != nil {
			select {
			case <-_srv.ctx.Done():
				return
			default:
				_srv.log.Errorf("epoll conn accept: %s", err.Error())

				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					time.Sleep(1 * time.Second)
					continue
				}
				return
			}
		}

		if err := _srv.epoll.AddOrClose(conn); err != nil {
			_srv.log.Errorf("epoll add conn: %s", err.Error())
		}
	}
}

func (_srv *Server) epollAccept() {
	for {
		select {
		case <-_srv.ctx.Done():
			return
		default:
			clist, err := _srv.epoll.Wait()
			switch err {
			case nil:
			case serv.ErrEpollEmptyEvents:
				continue
			case unix.EINTR:
				continue
			default:
				_srv.log.Errorf("epoll accept conn: %s", err.Error())
				continue
			}

			for _, c := range clist {
				go func(conn *netConnItem) {
					defer conn.Await(false)
					er := connection(conn.Conn, _srv.handler)
					switch er {
					case nil:
					case io.EOF:
						_ = _srv.epoll.Close(conn)
					default:
						_srv.log.Errorf("epoll bad conn from %s: %s", conn.Conn.RemoteAddr().String(), err.Error())
						_ = _srv.epoll.Close(conn)
					}
				}(c)
			}
		}
	}
}

func (_srv *Server) Down() (err error) {
	if _srv.listener == nil {
		return
	}
	defer func() {
		_srv.listener = nil
	}()
	_srv.cncl()
	if e := _srv.epoll.CloseAll(); e != nil {
		err = errors.Wrap(err, e.Error())
	}
	if e := _srv.listener.Close(); e != nil {
		err = errors.Wrap(err, e.Error())
	}
	_srv.log.Infof("tcp server stopped on %s", _srv.conf.Addr)
	return
}
