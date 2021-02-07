/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"net/http"
	"time"

	"github.com/deweppro/go-http/v2/servers/httpsrv"
	"github.com/deweppro/go-logger"
)

//Simple ...
type Simple struct{}

//ServeHTTP ...
func (s *Simple) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello world"))
}

func main() {
	simple := &Simple{}
	srv := httpsrv.NewCustomServer(httpsrv.ConfigItem{Addr: "localhost:8090"}, simple, logger.Default())
	if err := srv.Up(); err != nil {
		panic(err)
	}

	<-time.After(time.Minute)

	if err := srv.Down(); err != nil {
		panic(err)
	}
}
