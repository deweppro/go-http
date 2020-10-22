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

func (_epoll *Epoll) AddOrClose(c net.Conn) error {
	fd := serv.GetFD(c)
	err := unix.EpollCtl(_epoll.epfd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: cEpollEvents, Fd: int32(fd)})
	if err != nil {
		if er := c.Close(); er != nil {
			err = errors.Wrap(err, er.Error())
		}
		return err
	}
	_epoll.Lock()
	_epoll.conn[fd] = &netConnItem{Conn: c, Fd: fd}
	_epoll.Unlock()
	_epoll.log.Infof("epoll open connect: %s", c.RemoteAddr().String())
	return nil
}

func (_epoll *Epoll) removeFD(fd int) error {
	return unix.EpollCtl(_epoll.epfd, syscall.EPOLL_CTL_DEL, fd, nil)
}

func (_epoll *Epoll) Close(c *netConnItem) error {
	_epoll.Lock()
	defer _epoll.Unlock()
	return _epoll.closeConn(c)
}

func (_epoll *Epoll) closeConn(c *netConnItem) error {
	if err := _epoll.removeFD(c.Fd); err != nil {
		return err
	}
	delete(_epoll.conn, c.Fd)
	_epoll.log.Infof("epoll close connect: %s", c.Conn.RemoteAddr().String())
	return c.Conn.Close()
}

func (_epoll *Epoll) CloseAll() (er error) {
	_epoll.Lock()
	defer _epoll.Unlock()

	for _, conn := range _epoll.conn {
		if err := _epoll.closeConn(conn); err != nil {
			er = errors.Wrap(er, err.Error())
		}
	}
	_epoll.conn = make(netConnMap)
	return
}

func (_epoll *Epoll) getConn(fd int) (*netConnItem, bool) {
	_epoll.RLock()
	conn, ok := _epoll.conn[fd]
	_epoll.RUnlock()
	return conn, ok
}

func (_epoll *Epoll) Wait() (NetConns, error) {
	n, err := unix.EpollWait(_epoll.epfd, _epoll.events, cEventIntervalms)
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, serv.ErrEpollEmptyEvents
	}

	_epoll.elist = _epoll.elist[:0]
	for i := 0; i < n; i++ {
		fd := int(_epoll.events[i].Fd)
		conn, ok := _epoll.getConn(fd)
		if !ok {
			_ = _epoll.removeFD(fd)
			continue
		}
		if conn.IsAwait() {
			continue
		}
		conn.Await(true)

		switch _epoll.events[i].Events {
		case unix.POLLIN:
			_epoll.elist = append(_epoll.elist, conn)
		default:
			if err := _epoll.Close(conn); err != nil {
				_epoll.log.Errorf("epoll close connect: %s", err.Error())
			}
		}
	}

	return _epoll.elist, nil
}
