# go-http

[![Coverage Status](https://coveralls.io/repos/github/deweppro/go-http/badge.svg?branch=main)](https://coveralls.io/github/deweppro/go-http?branch=main)
[![Release](https://img.shields.io/github/release/deweppro/go-http.svg?style=flat-square)](https://github.com/deweppro/go-http/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/deweppro/go-http)](https://goreportcard.com/report/github.com/deweppro/go-http)
[![Build Status](https://travis-ci.com/deweppro/go-http.svg?branch=main)](https://travis-ci.com/deweppro/go-http)

# How to use it

## Web server

```go
package main

import (
	"net/http"
	"time"

	"github.com/deweppro/go-http/web"
	"github.com/deweppro/go-logger"
)

func main() {

	srv := web.NewCustomServer(web.ConfigItem{Addr: "localhost:8080"}, logger.Default())
	if err := srv.Up(); err != nil {
		panic(err)
	}

	srv.Router().AddRoutes(
		web.Handler{Method: http.MethodGet, Path: "/", Formatter: web.JSONFormatter, Call: func(ctx *web.Context) error {
			return ctx.Encode(func() (int, web.Headers, interface{}) {
				return 200, web.Headers{"x-trace-id": "999-999-999"}, web.ResponseModel{Data: 911}
			})
		}},
	)

	<-time.After(time.Minute)

	if err := srv.Down(); err != nil {
		panic(err)
	}
}

```

## Debug server (pprof)

```go
package main

import (
	"time"

	"github.com/deweppro/go-http/web"
	"github.com/deweppro/go-logger"
)

func main() {

	debug := web.NewCustomDebug(web.ConfigItem{Addr: "localhost:8090"}, logger.Default())
	if err := debug.Up(); err != nil {
		panic(err)
	}

	<-time.After(time.Minute)

	if err := debug.Down(); err != nil {
		panic(err)
	}
}
```

## TCP server (epoll)

```go
package main

import (
	"io"
	"time"

	"github.com/deweppro/go-http/epoll"
	"github.com/deweppro/go-logger"
)

func main() {

	srv := epoll.NewCustomServer(epoll.ConfigItem{Addr: "localhost:8080"}, logger.Default())
	srv.Handler(func(bytes []byte, writer io.Writer) error {
		_, err := writer.Write(bytes)
		return err
	})
	if err := srv.Up(); err != nil {
		panic(err)
	}

	<-time.After(time.Minute)

	if err := srv.Down(); err != nil {
		panic(err)
	}
}

```

## Web server + websocket

```go
package main

import (
	"context"
	"net/http"
	"time"

	"github.com/deweppro/go-http/web"
	"github.com/deweppro/go-http/ws"
	"github.com/deweppro/go-logger"
)

func main() {

	wsock := ws.New(10, 128, logger.Default())

	srv := web.NewCustomServer(web.ConfigItem{Addr: "localhost:8080"}, logger.Default())
	if err := srv.Up(); err != nil {
		panic(err)
	}

	srv.Router().AddRoutes(
		web.Handler{Method: http.MethodGet, Path: "/", Formatter: web.EmptyFormatter, Call: func(ctx *web.Context) error {
			return wsock.Handler(ctx.Writer, ctx.Reader,
				func(out chan<- *ws.Message, in <-chan *ws.Message, ctx context.Context, cncl context.CancelFunc) {
					i := 0
					for {
						select {
						case <-ctx.Done():
							return
						case msg := <-in:
							out <- msg
							if i == 3 {
								cncl()
							}
							i++
						}
					}
				},
			)
		}})

	<-time.After(time.Minute)

	if err := srv.Down(); err != nil {
		panic(err)
	}
}

```