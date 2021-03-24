package main

import (
	"net/http"
	"time"

	"github.com/deweppro/go-http/web/routes"
	"github.com/deweppro/go-http/web/server"
	"github.com/deweppro/go-logger"
)

func main() {
	conf := &server.Config{HTTP: server.ConfigItem{Addr: ":8080"}}
	route := routes.NewRouter()
	serv := server.New(conf, route, logger.Default())

	route.Route("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello"))
	}, http.MethodGet, http.MethodPost)

	route.Route("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("*")
	}, http.MethodGet, http.MethodPost)

	route.Global(
		routes.RecoveryMiddleware(logger.Default()),
		routes.ThrottlingMiddleware(1000),
		routes.CORSMiddleware(routes.CORSConfig{
			Age:     100,
			Origin:  "localhost",
			Methods: []string{http.MethodGet, http.MethodPost},
			Headers: []string{"X-Token"},
		}),
		func(ctrlFunc routes.CtrlFunc) routes.CtrlFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Trace", "0000")
				ctrlFunc(w, r)
			}
		},
	)

	if err := serv.Up(); err != nil {
		panic(err)
	}

	<-time.After(3 * time.Second)

	if err := serv.Down(); err != nil {
		panic(err)
	}

	logger.Close()
}
