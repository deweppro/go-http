/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package httpsrv

import (
	"net/http"
	"sync"
	"time"

	"github.com/deweppro/go-http/v2/proto"
	"github.com/deweppro/go-logger"
	"github.com/pkg/errors"
)

const (
	defaultTimeOut = 10 * time.Second
)

type (
	//Server ...
	Server struct {
		conf ConfigItem
		srv  *http.Server
		hand http.Handler
		log  logger.Logger
		wg   sync.WaitGroup
	}
)

//NewServer ...
func NewServer(c *Config, h http.Handler, l logger.Logger) *Server {
	return NewCustomServer(c.HTTP, h, l)
}

//NewCustomServer ...
func NewCustomServer(c ConfigItem, h http.Handler, l logger.Logger) *Server {
	srv := &Server{
		conf: c,
		hand: h,
		log:  l,
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

//Up ...
func (s *Server) Up() error {
	if s.srv != nil {
		return errors.Wrapf(proto.ErrServAlreadyRunning, "on %s", s.conf.Addr)
	}
	if len(s.conf.Addr) == 0 {
		addr, err := proto.RandomPort("localhost")
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
		Handler:      s.hand,
	}
	go s.run()
	return nil
}

//run ...
func (s *Server) run() {
	s.log.Infof("http server started on %s", s.conf.Addr)
	s.wg.Add(1)
	defer func() {
		s.srv = nil
		s.wg.Done()
	}()
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.log.Errorf("http server stopped on %s with error: %s", s.conf.Addr, err.Error())
		return
	}
	s.log.Infof("http server stopped on %s", s.conf.Addr)
}

//Down ...
func (s *Server) Down() error {
	err := s.srv.Close()
	s.wg.Wait()
	return err
}
