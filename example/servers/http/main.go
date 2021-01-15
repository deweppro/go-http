/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package main

import (
	http2 "net/http"
	"time"

	"github.com/deweppro/go-http/v2/servers/http"
	"github.com/deweppro/go-logger"
)

type Simple struct{}

func (s *Simple) ServeHTTP(w http2.ResponseWriter, r *http2.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello world"))
}

func main() {
	simple := &Simple{}
	srv := http.NewCustomServer(http.ConfigItem{Addr: "localhost:8090"}, simple, logger.Default())
	if err := srv.Up(); err != nil {
		panic(err)
	}

	<-time.After(time.Minute)

	if err := srv.Down(); err != nil {
		panic(err)
	}
}
