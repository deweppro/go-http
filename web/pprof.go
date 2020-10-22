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
	Debug struct {
		srv *Server
	}
)

func NewDebug(conf *DebugConfig, log logger.Logger) *Debug {
	return NewCustomDebug(conf.Debug, log)
}

func NewCustomDebug(conf ConfigItem, log logger.Logger) *Debug {
	debug := &Debug{srv: NewCustomServer(conf, log)}
	debug.srv.Router().InjectRoutes(debug)
	return debug
}

func (_debug *Debug) Handlers() []Handler {
	return []Handler{
		{Method: http.MethodGet, Path: "/debug/pprof", Call: _debug.call(pprof.Index)},
		{Method: http.MethodGet, Path: "/debug/pprof/goroutine", Call: _debug.call(pprof.Index)},
		{Method: http.MethodGet, Path: "/debug/pprof/allocs", Call: _debug.call(pprof.Index)},
		{Method: http.MethodGet, Path: "/debug/pprof/block", Call: _debug.call(pprof.Index)},
		{Method: http.MethodGet, Path: "/debug/pprof/heap", Call: _debug.call(pprof.Index)},
		{Method: http.MethodGet, Path: "/debug/pprof/heap", Call: _debug.call(pprof.Index)},
		{Method: http.MethodGet, Path: "/debug/pprof/mutex", Call: _debug.call(pprof.Index)},
		{Method: http.MethodGet, Path: "/debug/pprof/threadcreate", Call: _debug.call(pprof.Index)},
		{Method: http.MethodGet, Path: "/debug/pprof/cmdline", Call: _debug.call(pprof.Cmdline)},
		{Method: http.MethodGet, Path: "/debug/pprof/profile", Call: _debug.call(pprof.Profile)},
		{Method: http.MethodGet, Path: "/debug/pprof/symbol", Call: _debug.call(pprof.Symbol)},
		{Method: http.MethodGet, Path: "/debug/pprof/trace", Call: _debug.call(pprof.Trace)},
	}
}

func (_debug *Debug) call(cb func(http.ResponseWriter, *http.Request)) CallFunc {
	return func(message *Context) error {
		cb(message.Writer, message.Reader)
		return nil
	}
}

func (_debug *Debug) Up() error {
	return _debug.srv.Up()
}

func (_debug *Debug) Down() error {
	return _debug.srv.Down()
}
