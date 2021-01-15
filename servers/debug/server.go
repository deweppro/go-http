/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package debug

import (
	http2 "net/http"
	"net/http/pprof"

	"github.com/deweppro/go-http/v2/servers/http"
	"github.com/deweppro/go-logger"
)

type (
	Debug struct {
		srv *http.Server
		log logger.Logger
	}
)

func New(conf *Config, log logger.Logger) *Debug {
	return NewCustom(conf.Debug, log)
}

func NewCustom(conf http.ConfigItem, log logger.Logger) *Debug {
	d := &Debug{log: log}
	d.srv = http.NewCustomServer(conf, d, log)
	return d
}

func (d *Debug) ServeHTTP(w http2.ResponseWriter, r *http2.Request) {
	switch r.URL.Path {
	case "/debug/pprof/", "/debug/pprof/goroutine", "/debug/pprof/allocs",
		"/debug/pprof/block", "/debug/pprof/heap", "/debug/pprof/mutex",
		"/debug/pprof/threadcreate":
		pprof.Index(w, r)
	case "", "/", "/goroutine", "/allocs", "/block", "/heap", "/mutex", "/threadcreate":
		pprof.Index(w, r)
	case "/cmdline", "/debug/pprof/cmdline":
		pprof.Cmdline(w, r)
	case "/profile", "/debug/pprof/profile":
		pprof.Profile(w, r)
	case "/symbol", "/debug/pprof/symbol":
		pprof.Symbol(w, r)
	case "/trace", "/debug/pprof/trace":
		pprof.Trace(w, r)
	default:
		d.log.Errorf("fail request to: %s", r.URL.Path)
		w.WriteHeader(http2.StatusNotFound)
	}
}

func (d *Debug) Up() error {
	return d.srv.Up()
}

func (d *Debug) Down() error {
	return d.srv.Down()
}
