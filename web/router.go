/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package web

import (
	"net/http"
	"strings"
	"sync"

	"github.com/deweppro/go-logger"
)

type (
	CallFunc func(*Context) error
	Handler  struct {
		Method     string
		Path       string
		Call       CallFunc
		Middleware CallFunc
		Formatter  FormatterFunc
	}
	ModuleInjecter interface {
		Handlers() []Handler
	}
	Route struct {
		Call       CallFunc
		Middleware CallFunc
		Formatter  FormatterFunc
	}
	Router struct {
		routes map[string]*Route
		log    logger.Logger
		sync.RWMutex
	}
)

func newRouter(log logger.Logger) *Router {
	return &Router{
		routes: make(map[string]*Route),
		log:    log,
	}
}

func (_rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := &Context{Reader: r, Writer: w, Formatter: TextFormatter}

	_rtr.RLock()
	ch := func() *Route {
		if call, ok := _rtr.routes[strings.Join([]string{r.Method, r.URL.Path}, ":")]; ok {
			return call
		} else if call, ok := _rtr.routes[strings.Join([]string{"*", r.URL.Path}, ":")]; ok {
			return call
		} else if call, ok := _rtr.routes["*:"]; ok {
			return call
		}
		return nil
	}()
	_rtr.RUnlock()

	if ch == nil {
		_ = msg.Empty(http.StatusNotFound)
		return
	}

	if ch.Formatter != nil {
		msg.Formatter = ch.Formatter
	}

	if ch.Middleware != nil {
		if err := ch.Middleware(msg); err != nil {
			_ = msg.Error(http.StatusBadRequest, err)
			return
		}
	}

	if err := ch.Call(msg); err != nil {
		_ = msg.Error(http.StatusBadRequest, err)
	}
}

func (_rtr *Router) getRoute(method, path string) *Route {
	_rtr.RLock()
	defer _rtr.RUnlock()

	if call, ok := _rtr.routes[method+":"+path]; ok {
		return call
	} else if call, ok := _rtr.routes["*:"+path]; ok {
		return call
	} else if call, ok := _rtr.routes["*:"]; ok {
		return call
	}
	return nil
}

func (_rtr *Router) AddRoutes(handlers ...Handler) {
	_rtr.Lock()
	defer _rtr.Unlock()

	for _, handler := range handlers {
		_rtr.log.Infof("add route %s:%s", handler.Method, handler.Path)
		_rtr.routes[handler.Method+":"+handler.Path] = &Route{
			Call:       handler.Call,
			Middleware: handler.Middleware,
			Formatter:  handler.Formatter,
		}
	}
}

func (_rtr *Router) InjectRoutes(mod ModuleInjecter) {
	_rtr.AddRoutes(mod.Handlers()...)
}
