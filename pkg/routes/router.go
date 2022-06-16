package routes

import (
	"context"
	"net/http"
	"sync"

	"github.com/deweppro/go-http/internal"
)

var _ http.Handler = (*Router)(nil)

//Router model
type Router struct {
	handler *handler
	lock    sync.RWMutex
}

//NewRouter init new router
func NewRouter() *Router {
	return &Router{
		handler: newHandler(),
	}
}

//Route add new route
func (v *Router) Route(path string, ctrl CtrlFunc, methods ...string) {
	v.lock.Lock()
	v.handler.Route(path, ctrl, methods)
	v.lock.Unlock()
}

//Global add global middlewares
func (v *Router) Global(middlewares ...MiddlFunc) {
	v.lock.Lock()
	v.handler.Middlewares("", middlewares...)
	v.lock.Unlock()
}

//Middlewares add middlewares to route
func (v *Router) Middlewares(path string, middlewares ...MiddlFunc) {
	v.lock.Lock()
	v.handler.Middlewares(path, middlewares...)
	v.lock.Unlock()
}

//NoFoundHandler handler call if route not found
func (v *Router) NoFoundHandler(call CtrlFunc) {
	v.lock.Lock()
	v.handler.NoFoundHandler(call)
	v.lock.Unlock()
}

//ServeHTTP http interface
func (v *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.lock.RLock()
	defer v.lock.RUnlock()

	code, next, vars, midd := v.handler.Match(r.URL.Path, r.Method)
	if code != http.StatusOK {
		next = codeHandler(code)
	}

	ctx := r.Context()
	for key, val := range vars {
		ctx = context.WithValue(ctx, internal.VarsKey(key), val)
	}

	for i := len(midd) - 1; i >= 0; i-- {
		next = midd[i](next)
	}
	next(w, r.WithContext(ctx))
}

func codeHandler(code int) CtrlFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	}
}
