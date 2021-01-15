/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

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
