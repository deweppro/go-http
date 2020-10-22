/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

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
