package routes

import (
	"net/http"
	"strings"
)

//Handler model
type Handler struct {
	list        map[string]*Handler
	methods     map[string]struct{}
	middlewares []MiddlFunc
	controller  CtrlFunc
}

//NewHandler getting new handler
func NewHandler() *Handler {
	return &Handler{
		list:        make(map[string]*Handler),
		methods:     make(map[string]struct{}),
		middlewares: make([]MiddlFunc, 0),
	}
}

//Route add new route
func (h *Handler) Route(uris []string, pos int, ctrl CtrlFunc, methods []string) {
	if pos >= len(uris) {
		return
	}
	uri := uris[pos]
	uh, ok := h.list[uri]
	if !ok {
		uh = NewHandler()
		h.list[uri] = uh
	}
	if pos == len(uris)-1 {
		uh.controller = ctrl
		for _, m := range methods {
			uh.methods[strings.ToUpper(m)] = struct{}{}
		}
		return
	}
	uh.Route(uris, pos+1, ctrl, methods)
}

//Middlewares add middleware to route
func (h *Handler) Middlewares(uris []string, pos int, middlewares []MiddlFunc) {
	if pos >= len(uris) {
		return
	}
	uri := uris[pos]
	uh, ok := h.list[uri]
	if !ok {
		uh = NewHandler()
		h.list[uri] = uh
	}
	if pos == len(uris)-1 {
		for _, m := range middlewares {
			uh.middlewares = append(uh.middlewares, m)
		}
		return
	}
	uh.Middlewares(uris, pos+1, middlewares)
}

//Match find route in tree
func (h *Handler) Match(uris []string, pos int, method string) (code int, ctrl CtrlFunc, midd []MiddlFunc) {
	if pos >= len(uris) {
		return http.StatusBadRequest, nil, nil
	}
	if uh, ok := h.list[uris[pos]]; ok {
		midd = append(midd, uh.middlewares...)

		if pos == len(uris)-1 {
			if _, ok = uh.methods[method]; !ok {
				return http.StatusMethodNotAllowed, nil, nil
			}
			code, ctrl = http.StatusOK, uh.controller
			return
		}

		co, ct, mi := uh.Match(uris, pos+1, method)
		code, ctrl = co, ct
		midd = append(midd, mi...)
		return
	}
	return http.StatusNotFound, nil, nil
}
