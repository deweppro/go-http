package routes

import (
	"net/http"
	"sync"
)

//Router model
type Router struct {
	middlewares []MiddlFunc
	handlers    *Handler
	lock        sync.RWMutex
}

//NewRouter init new router
func NewRouter() *Router {
	return &Router{
		middlewares: make([]MiddlFunc, 0),
		handlers:    NewHandler(),
	}
}

//Route add new route
func (o *Router) Route(uri string, cb CtrlFunc, methods ...string) {
	o.lock.Lock()
	o.handlers.Route(SplitURI(uri), 0, cb, methods)
	o.lock.Unlock()
}

//Global add global middlewares
func (o *Router) Global(middlewares ...MiddlFunc) {
	o.lock.Lock()
	o.middlewares = append(o.middlewares, middlewares...)
	o.lock.Unlock()
}

//Middlewares add middlewares to route
func (o *Router) Middlewares(prefix string, middlewares ...MiddlFunc) {
	o.lock.Lock()
	o.handlers.Middlewares(SplitURI(prefix), 0, middlewares)
	o.lock.Unlock()
}

//ServeHTTP http interface
func (o *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	next := o.match()
	for i := len(o.middlewares) - 1; i >= 0; i-- {
		next = o.middlewares[i](next)
	}
	next(w, r)
}

func (o *Router) match() CtrlFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o.lock.RLock()
		defer o.lock.RUnlock()
		code, ctrl, mids := o.handlers.Match(SplitURI(r.RequestURI), 0, r.Method)
		if code != http.StatusOK {
			w.WriteHeader(code)
			return
		}

		next := ctrl
		for i := len(mids) - 1; i >= 0; i-- {
			next = mids[i](next)
		}
		next(w, r)
	}
}
