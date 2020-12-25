/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package web

import (
	"net/http"
	"sync"

	"github.com/deweppro/go-logger"
)

const DefaultVersion uint64 = 1

type (
	//Caller ...
	Caller func(*Context) error
	//MiddlewareCaller ...
	MiddlewareCaller func(*Context) (int, error)
	//Handler ...
	Handler struct {
		Method     []string
		Path       string
		Call       VerCaller
		Middleware MiddlewareCaller
	}
	VerCaller map[uint64]Caller
	//RouteItem ...
	RouteItem struct {
		Call       VerCaller
		Middleware MiddlewareCaller
	}
	//Route ...
	Route struct {
		routes map[string]*RouteItem
		log    logger.Logger
		sync.RWMutex
	}
	//RouteInjector ...
	RouteInjector interface {
		Handlers() []Handler
	}
	//Router ...
	Router interface {
		AddRoutes(handlers ...Handler)
		InjectRoutes(mod RouteInjector)
	}
)

func newRouter(log logger.Logger) *Route {
	return &Route{
		routes: make(map[string]*RouteItem),
		log:    log,
	}
}

//ServeHTTP ...
func (o *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := &Context{Reader: r, Writer: w}
	route := o.getRoute(r.Method, r.URL.Path)

	if route == nil {
		if err := msg.Write(http.StatusNotFound, nil, nil); err != nil {
			o.log.Errorf("route not found: %s", err.Error())
		}
		return
	}

	call, ok := route.Call[msg.Version()]
	if !ok || call == nil {
		if err := msg.Write(http.StatusNotFound, nil, nil); err != nil {
			o.log.Errorf("route version not found: %s", err.Error())
		}
		return
	}

	if route.Middleware != nil {
		if code, err := route.Middleware(msg); err != nil {
			if er := msg.Write(code, []byte(err.Error()), nil); er != nil {
				o.log.Errorf("middleware: %s", er.Error())
			}
			return
		}
	}

	if err := call(msg); err != nil {
		if er := msg.Write(http.StatusInternalServerError, []byte(err.Error()), nil); er != nil {
			o.log.Errorf("call: %s", er.Error())
		}
	}
}

func (o *Route) getRoute(method, path string) *RouteItem {
	o.RLock()
	defer o.RUnlock()

	if call, ok := o.routes[method+":"+path]; ok {
		return call
	}
	if call, ok := o.routes[method+":*"]; ok {
		return call
	}
	return nil
}

//AddRoutes ...
func (o *Route) AddRoutes(handlers ...Handler) {
	o.Lock()
	defer o.Unlock()

	for _, handler := range handlers {
		for _, method := range handler.Method {
			o.log.Infof("add route %s:%s", handler.Method, handler.Path)
			o.routes[method+":"+handler.Path] = &RouteItem{
				Call:       handler.Call,
				Middleware: handler.Middleware,
			}
		}

	}
}

//InjectRoutes ...
func (o *Route) InjectRoutes(mod RouteInjector) {
	o.AddRoutes(mod.Handlers()...)
}
