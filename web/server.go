/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package web

import (
	"net/http"
	"sync"
	"time"

	serv "github.com/deweppro/go-http"

	"github.com/deweppro/go-logger"
	"github.com/pkg/errors"
)

const (
	defaultTimeOut = 10 * time.Second
)

type Server struct {
	conf  ConfigItem
	srv   *http.Server
	route *Router
	log   logger.Logger
	wg    sync.WaitGroup
}

func NewServer(conf *HTTPConfig, log logger.Logger) *Server {
	return NewCustomServer(conf.HTTP, log)
}

func NewCustomServer(conf ConfigItem, log logger.Logger) *Server {
	srv := &Server{
		conf:  conf,
		route: newRouter(log),
		log:   log,
	}
	srv.checkConf()
	return srv
}

func (_srv *Server) checkConf() {
	if _srv.conf.ReadTimeout == 0 {
		_srv.conf.ReadTimeout = defaultTimeOut
	}
	if _srv.conf.WriteTimeout == 0 {
		_srv.conf.WriteTimeout = defaultTimeOut
	}
}

func (_srv *Server) Router() *Router {
	return _srv.route
}

func (_srv *Server) Up() error {
	if _srv.srv != nil {
		return errors.Wrapf(serv.ErrServAlreadyRunning, "on %s", _srv.conf.Addr)
	}
	if len(_srv.conf.Addr) == 0 {
		addr, err := serv.RandomPort("localhost")
		if err != nil {
			return errors.Wrap(err, "get random port")
		}
		_srv.conf.Addr = addr
	}
	_srv.srv = &http.Server{
		Addr:         _srv.conf.Addr,
		ReadTimeout:  _srv.conf.ReadTimeout,
		WriteTimeout: _srv.conf.WriteTimeout,
		IdleTimeout:  _srv.conf.IdleTimeout,
		Handler:      _srv.route,
	}
	_srv.log.Infof("http server started on %s", _srv.conf.Addr)
	go func() {
		_srv.wg.Add(1)
		defer func() {
			_srv.srv = nil
			_srv.wg.Done()
		}()
		if err := _srv.srv.ListenAndServe(); err != http.ErrServerClosed {
			_srv.log.Errorf("http server stopped on %s with error: %+v", _srv.conf.Addr, err)
			return
		}
		_srv.log.Infof("http server stopped on %s", _srv.conf.Addr)
	}()
	return nil
}

func (_srv *Server) Down() error {
	err := _srv.srv.Close()
	_srv.wg.Wait()
	return err
}
