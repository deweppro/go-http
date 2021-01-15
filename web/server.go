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

type (
	//Server ...
	Server struct {
		conf  ConfigItem
		srv   *http.Server
		route *Route
		log   logger.Logger
		wg    sync.WaitGroup
	}
)

//NewServer ...
func NewServer(conf *HTTPConfig, log logger.Logger) *Server {
	return NewCustomServer(conf.HTTP, log)
}

//NewCustomServer ...
func NewCustomServer(conf ConfigItem, log logger.Logger) *Server {
	srv := &Server{
		conf:  conf,
		route: newRouter(log),
		log:   log,
	}
	srv.checkConf()
	return srv
}

func (s *Server) checkConf() {
	if s.conf.ReadTimeout == 0 {
		s.conf.ReadTimeout = defaultTimeOut
	}
	if s.conf.WriteTimeout == 0 {
		s.conf.WriteTimeout = defaultTimeOut
	}
}

//Router ...
func (s *Server) Router() Router {
	return s.route
}

//Up ...
func (s *Server) Up() error {
	if s.srv != nil {
		return errors.Wrapf(serv.ErrServAlreadyRunning, "on %s", s.conf.Addr)
	}
	if len(s.conf.Addr) == 0 {
		addr, err := serv.RandomPort("localhost")
		if err != nil {
			return errors.Wrap(err, "get random port")
		}
		s.conf.Addr = addr
	}
	s.srv = &http.Server{
		Addr:         s.conf.Addr,
		ReadTimeout:  s.conf.ReadTimeout,
		WriteTimeout: s.conf.WriteTimeout,
		IdleTimeout:  s.conf.IdleTimeout,
		Handler:      s.route,
	}
	s.log.Infof("http server started on %s", s.conf.Addr)
	go func() {
		s.wg.Add(1)
		defer func() {
			s.srv = nil
			s.wg.Done()
		}()
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			s.log.Errorf("http server stopped on %s with error: %+v", s.conf.Addr, err)
			return
		}
		s.log.Infof("http server stopped on %s", s.conf.Addr)
	}()
	return nil
}

//Down ...
func (s *Server) Down() error {
	err := s.srv.Close()
	s.wg.Wait()
	return err
}
