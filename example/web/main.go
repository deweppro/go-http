/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

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
