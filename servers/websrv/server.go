/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package websrv

import (
	"net/http"

	"github.com/deweppro/go-http/v2/servers/httpsrv"

	proto "github.com/deweppro/go-http/v2"
	"github.com/deweppro/go-logger"
)

type (
	Server struct {
		conf *Config
		srv  *httpsrv.Server
		hub  *proto.Server
		log  logger.Logger
	}
)

func NewServer(c *Config, h *proto.Server, l logger.Logger) *Server {
	s := &Server{
		conf: c,
		hub:  h,
		log:  l,
	}
	s.srv = httpsrv.NewCustomServer(c.HTTP, s, l)
	return s
}

func (s *Server) Up() error {
	return s.srv.Up()
}

func (s *Server) Down() error {
	return s.srv.Down()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := proto.NewRequest()
	res := proto.NewResponse()
	res.SetMeta(s.conf.Headers.DefaultHeaders)

	err := req.UpdateFromHTTP(r, s.conf.Headers.ProxyHeaders...)
	res.SetUUID(req.GetUUID())

	if err == nil {
		s.hub.Call(req, res)
	} else {
		res.SetStatusCode(proto.StatusCodeFail)
		res.Body = []byte(err.Error())
		s.log.Errorf("request read error: %s", err.Error())
	}

	if err = res.WriteToHTTP(w); err != nil {
		s.log.Errorf("response write error: %s", err.Error())
	}
}
