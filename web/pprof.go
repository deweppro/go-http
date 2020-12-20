/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package web

import (
	"net/http"
	"net/http/pprof"

	"github.com/deweppro/go-logger"
)

type (
	//Debug ...
	Debug struct {
		srv *Server
	}
)

var _ RouteInjector = (*Debug)(nil)

//NewDebug ...
func NewDebug(conf *DebugConfig, log logger.Logger) *Debug {
	return NewCustomDebug(conf.Debug, log)
}

//NewCustomDebug ...
func NewCustomDebug(conf ConfigItem, log logger.Logger) *Debug {
	debug := &Debug{srv: NewCustomServer(conf, log)}
	debug.srv.Router().InjectRoutes(debug)
	return debug
}

//Handlers ...
func (d *Debug) Handlers() []Handler {
	return []Handler{
		{Method: []string{http.MethodGet}, Path: "/debug/pprof", Call: d.call(pprof.Index)},
		{Method: []string{http.MethodGet}, Path: "/debug/pprof/goroutine", Call: d.call(pprof.Index)},
		{Method: []string{http.MethodGet}, Path: "/debug/pprof/allocs", Call: d.call(pprof.Index)},
		{Method: []string{http.MethodGet}, Path: "/debug/pprof/block", Call: d.call(pprof.Index)},
		{Method: []string{http.MethodGet}, Path: "/debug/pprof/heap", Call: d.call(pprof.Index)},
		{Method: []string{http.MethodGet}, Path: "/debug/pprof/heap", Call: d.call(pprof.Index)},
		{Method: []string{http.MethodGet}, Path: "/debug/pprof/mutex", Call: d.call(pprof.Index)},
		{Method: []string{http.MethodGet}, Path: "/debug/pprof/threadcreate", Call: d.call(pprof.Index)},
		{Method: []string{http.MethodGet}, Path: "/debug/pprof/cmdline", Call: d.call(pprof.Cmdline)},
		{Method: []string{http.MethodGet}, Path: "/debug/pprof/profile", Call: d.call(pprof.Profile)},
		{Method: []string{http.MethodGet}, Path: "/debug/pprof/symbol", Call: d.call(pprof.Symbol)},
		{Method: []string{http.MethodGet}, Path: "/debug/pprof/trace", Call: d.call(pprof.Trace)},
	}
}

func (d *Debug) call(cb func(http.ResponseWriter, *http.Request)) Caller {
	return func(message *Context) error {
		cb(message.Writer, message.Reader)
		return nil
	}
}

//Up ...
func (d *Debug) Up() error {
	return d.srv.Up()
}

//Down ...
func (d *Debug) Down() error {
	return d.srv.Down()
}
