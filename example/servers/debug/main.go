/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"time"

	"github.com/deweppro/go-http/v2/servers/debug"
	"github.com/deweppro/go-http/v2/servers/http"
	"github.com/deweppro/go-logger"
)

func main() {
	dbg := debug.NewCustom(http.ConfigItem{Addr: "localhost:8090"}, logger.Default())
	if err := dbg.Up(); err != nil {
		panic(err)
	}

	<-time.After(time.Minute)

	if err := dbg.Down(); err != nil {
		panic(err)
	}
}
