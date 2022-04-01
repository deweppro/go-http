package routes

import (
	"net/http"
	"strings"

	"github.com/deweppro/go-http/internal"
)

//CtrlFunc interface of controller
type CtrlFunc func(http.ResponseWriter, *http.Request)

//handler model
type handler struct {
	list        map[string]*handler
	methods     map[string]CtrlFunc
	matcher     *matcher
	middlewares []MiddlFunc
}

//newHandler getting new handler
func newHandler() *handler {
	return &handler{
		list:        make(map[string]*handler),
		methods:     make(map[string]CtrlFunc),
		matcher:     newMatcher(),
		middlewares: make([]MiddlFunc, 0),
	}
}

func (v *handler) append(path string) *handler {
	if uh, ok := v.list[path]; ok {
		return uh
	}
	uh := newHandler()
	v.list[path] = uh
	return uh
}

func (v *handler) next(path string, vars internal.VarsData) *handler {
	if uh, ok := v.list[path]; ok {
		return uh
	}
	uri, ok := v.matcher.Match(path, vars)
	if !ok {
		return nil
	}
	if uh, ok := v.list[uri]; ok {
		return uh
	}
	return nil
}

//Route add new route
func (v *handler) Route(path string, ctrl CtrlFunc, methods []string) {
	uris := split(path)
	var uh = v
	for _, uri := range uris {
		if hasMatcher(uri) {
			if err := uh.matcher.Add(uri); err != nil {
				panic(err)
			}
		}
		uh = uh.append(uri)
	}
	for _, m := range methods {
		uh.methods[strings.ToUpper(m)] = ctrl
	}
}

//Middlewares add middleware to route
func (v *handler) Middlewares(path string, middlewares ...MiddlFunc) {
	uris := split(path)
	var uh = v
	for _, uri := range uris {
		uh = v.append(uri)
	}
	uh.middlewares = append(uh.middlewares, middlewares...)
}

//Match find route in tree
func (v *handler) Match(path string, method string) (int, CtrlFunc, internal.VarsData, []MiddlFunc) {
	var uh = v

	uris := split(path)
	midd := append(make([]MiddlFunc, 0), uh.middlewares...)
	vars := internal.VarsData{}

	for _, uri := range uris {
		uh = uh.next(uri, vars)
		if uh == nil {
			return http.StatusNotFound, nil, nil, nil
		}
		midd = append(midd, uh.middlewares...)
	}

	ctrl, ok := uh.methods[method]
	if !ok {
		return http.StatusMethodNotAllowed, nil, nil, nil
	}

	return http.StatusOK, ctrl, vars, midd
}
