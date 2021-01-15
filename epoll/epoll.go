/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package epoll

import (
	"net"
	"sync"
	"syscall"

	"github.com/deweppro/go-logger"

	serv "github.com/deweppro/go-http"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type (
	netConnMap     map[int]*netConnItem
	NetConns       []*netConnItem
	unixEventSlice []unix.EpollEvent

	//Epoll ...
	Epoll struct {
		epfd   int
		conn   netConnMap
		events unixEventSlice
		elist  NetConns
		log    logger.Logger
		sync.RWMutex
	}
)

const (
	cEpollEvents     = unix.POLLIN | unix.POLLRDHUP | unix.POLLERR | unix.POLLHUP | unix.POLLNVAL
	cEventCount      = 100
	cEventIntervalms = 500
)

//NewEpoll ...
func NewEpoll(log logger.Logger) (*Epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &Epoll{
		epfd:   fd,
		conn:   make(netConnMap),
		events: make(unixEventSlice, cEventCount),
		elist:  make(NetConns, cEventCount),
		log:    log,
	}, nil
}

//AddOrClose ...
func (e *Epoll) AddOrClose(c net.Conn) error {
	fd := serv.GetFD(c)
	err := unix.EpollCtl(e.epfd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: cEpollEvents, Fd: int32(fd)})
	if err != nil {
		if er := c.Close(); er != nil {
			err = errors.Wrap(err, er.Error())
		}
		return err
	}
	e.Lock()
	e.conn[fd] = &netConnItem{Conn: c, Fd: fd}
	e.Unlock()
	e.log.Infof("epoll open connect: %s", c.RemoteAddr().String())
	return nil
}

func (e *Epoll) removeFD(fd int) error {
	return unix.EpollCtl(e.epfd, syscall.EPOLL_CTL_DEL, fd, nil)
}

//Close ...
func (e *Epoll) Close(c *netConnItem) error {
	e.Lock()
	defer e.Unlock()
	return e.closeConn(c)
}

func (e *Epoll) closeConn(c *netConnItem) error {
	if err := e.removeFD(c.Fd); err != nil {
		return err
	}
	delete(e.conn, c.Fd)
	e.log.Infof("epoll close connect: %s", c.Conn.RemoteAddr().String())
	return c.Conn.Close()
}

//CloseAll ...
func (e *Epoll) CloseAll() (er error) {
	e.Lock()
	defer e.Unlock()

	for _, conn := range e.conn {
		if err := e.closeConn(conn); err != nil {
			er = errors.Wrap(er, err.Error())
		}
	}
	e.conn = make(netConnMap)
	return
}

func (e *Epoll) getConn(fd int) (*netConnItem, bool) {
	e.RLock()
	conn, ok := e.conn[fd]
	e.RUnlock()
	return conn, ok
}

//Wait ...
func (e *Epoll) Wait() (NetConns, error) {
	n, err := unix.EpollWait(e.epfd, e.events, cEventIntervalms)
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, serv.ErrEpollEmptyEvents
	}

	e.elist = e.elist[:0]
	for i := 0; i < n; i++ {
		fd := int(e.events[i].Fd)
		conn, ok := e.getConn(fd)
		if !ok {
			_ = e.removeFD(fd)
			continue
		}
		if conn.IsAwait() {
			continue
		}
		conn.Await(true)

		switch e.events[i].Events {
		case unix.POLLIN:
			e.elist = append(e.elist, conn)
		default:
			if err := e.Close(conn); err != nil {
				e.log.Errorf("epoll close connect: %s", err.Error())
			}
		}
	}

	return e.elist, nil
}
