/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

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
