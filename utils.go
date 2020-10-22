/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package serv

import (
	"net"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrServAlreadyRunning = errors.New("server is already running")
	ErrEpollEmptyEvents   = errors.New("epoll events is empty")
)

func RandomPort(host string) (string, error) {
	host = strings.Join([]string{host, "0"}, ":")
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return host, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return host, err
	}
	return l.Addr().String(), l.Close()
}

func GetFD(c net.Conn) int {
	fd := reflect.Indirect(reflect.ValueOf(c)).FieldByName("fd")
	pfd := reflect.Indirect(fd).FieldByName("pfd")
	return int(pfd.FieldByName("Sysfd").Int())
}
