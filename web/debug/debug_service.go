package debug

import (
	"net/http"
	"net/http/pprof"

	"github.com/deweppro/go-http/web/routes"
	"github.com/deweppro/go-http/web/server"
	"github.com/deweppro/go-logger"
)

//Debug service model
type Debug struct {
	server *server.Server
	route  *routes.Router
}

//New init service
func New(conf *Config, log logger.Logger) *Debug {
	r := routes.NewRouter()
	return &Debug{
		server: server.NewCustom(conf.Debug, r, log),
		route:  r,
	}
}

//Up start service
func (o *Debug) Up() error {
	o.route.Route("/debug/pprof", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/goroutine", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/allocs", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/block", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/heap", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/mutex", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/threadcreate", pprof.Index, http.MethodGet)
	o.route.Route("/debug/pprof/cmdline", pprof.Cmdline, http.MethodGet)
	o.route.Route("/debug/pprof/profile", pprof.Profile, http.MethodGet)
	o.route.Route("/debug/pprof/symbol", pprof.Symbol, http.MethodGet)
	o.route.Route("/debug/pprof/trace", pprof.Trace, http.MethodGet)
	return o.server.Up()
}

//Down stop service
func (o *Debug) Down() error {
	return o.server.Down()
}
