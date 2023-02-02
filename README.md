# go-http

## Deprecated. Use https://github.com/deweppro/go-sdk

## Install

```sh
go get -u github.com/deweppro/go-http
```

## Examples

see more [here](examples)

### Debug

You can add http server with routes to receive data from pprof:

```go
import (
	"github.com/deweppro/go-http/servers"
	"github.com/deweppro/go-http/servers/debug"
	"github.com/deweppro/go-logger"
)
// use the configuration, and specify a host with a port 
// (example: localhost:8090, :8090 for bind 0.0.0.0:8090), 
// or specify only a host (example: localhost) to get a random port.
conf := servers.Config{Addr: ":8080"}
// сreate object and pass the config and logger to it.
serv := debug.New(conf, logger.Default())
serv.Up() // сall to start the server.
serv.Down() //  сall to stop the server.
```

### Routes for http

Create instance (it match to the interface `http.Handler`)

```go
import "github.com/deweppro/go-http/pkg/routes"

route := routes.NewRouter()
```

and add URL path and handler

```go
route.Route(<URL>, <Handler>, <Methods>...)

// example:
route.Route("/", IndexHandler, http.MethodGet) // only GET
route.Route("/page", PageHandler, http.MethodGet, http.MethodPost) // GET + POST
route.Route("/page-{id}", PageHandler, http.MethodGet, http.MethodPost) // GET + POST
route.Route("/page-{id:\d+}", PageHandler, http.MethodGet, http.MethodPost) // GET + POST
```

handler must match the interface: `CtrlFunc: func(http.ResponseWriter,*http.Request)`

```go
func IndexHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hello"))
}
```

Adding Global Middleware: (middleware must match the interface: `func(CtrlFunc)CtrlFunc`, where `CtrlFunc = func(http.ResponseWriter,*http.Request)`)

```go
route.Global(<Middleware>...)

// examples:
route.Global(
	// default: Recovery for route handler pain
    routes.RecoveryMiddleware(logger),
    // default: Limit on the number of requests
    routes.ThrottlingMiddleware(1000),
    // default: Setting up for сross-origin resource sharing (CORS)
    routes.CORSMiddleware(routes.CORSConfig{
        Age:     100,
        Origin:  []string{"localhost"},
        Methods: []string{http.MethodGet, http.MethodPost},
        Headers: []string{"X-Token"},
    }),
    // custom: for example, updating headers
    func(ctrlFunc routes.CtrlFunc) routes.CtrlFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("X-Trace", "0000")
            ctrlFunc(w, r)
        }
    },
)
```

you can also add middleware for any URL level:

```go
route.Middlewares(<URL PREFIX>, <Middleware>...)

// examples
// we have added several URL paths
route.Route("/pages/book/book-{id:\d+}", BookHandler, http.MethodGet)
route.Route("/pages/book/book-{title:[a-z]+}", BookHandler, http.MethodGet)
route.Route("/pages/book/list/all", ListBookHandler, http.MethodGet)
// and we want that when calling any page starting 
// from /pages/book level, the necessary middlewares are executed
route.Middlewares("/pages/book", 
	func(ctrlFunc routes.CtrlFunc) routes.CtrlFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("X-Book", "1")
            ctrlFunc(w, r)
        }
    }, 
)
```

### Web server

You can add web server:

```go
import (
    "github.com/deweppro/go-http/servers"
    "github.com/deweppro/go-http/servers/web"
)

// use the configuration, and specify a host with a port 
// (example: localhost:8080, :8080 for bind 0.0.0.0:8080), 
// or specify only a host (example: localhost) to get a random port.
conf := servers.Config{Addr: ":8080"}
// сreate object and pass the config, handler and logger to it.
// handler must match the interface - http.Handler
serv := server.New(conf, handler, logger)
serv.Up() // сall to start the server.
serv.Down() //  сall to stop the server.
```

### Full example

```go
package main

import (
	"net/http"
	"time"

	"github.com/deweppro/go-http/pkg/httputil/enc"
	"github.com/deweppro/go-http/pkg/routes"
	"github.com/deweppro/go-http/servers"
	"github.com/deweppro/go-http/servers/web"
	"github.com/deweppro/go-logger"
)

func main() {
	conf := servers.Config{Addr: ":8080"}
	route := routes.NewRouter()
	serv := web.New(conf, route, logger.Default())

	route.Route("/", func(w http.ResponseWriter, r *http.Request) {
		enc.Raw(w, []byte("Hello"))
	}, http.MethodGet, http.MethodPost)
	
	route.Route("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("*")
	}, http.MethodGet, http.MethodPost)

	route.Global(
		routes.RecoveryMiddleware(logger.Default()),
		routes.ThrottlingMiddleware(1000),
		routes.CORSMiddleware(routes.CORSConfig{
			Age:     100,
			Origin:  []string{"localhost"},
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

	<-time.After(60 * time.Second)

	if err := serv.Down(); err != nil {
		panic(err)
	}

	logger.Close()
}

```

## License

BSD-3-Clause License. See the LICENSE file for details.
