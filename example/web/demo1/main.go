/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"bytes"
	"strconv"
	"time"

	"github.com/deweppro/go-http/v2/proto"
	"github.com/deweppro/go-http/v2/servers/httpsrv"
	"github.com/deweppro/go-http/v2/servers/websrv"
	"github.com/deweppro/go-logger"
)

//Simple ...
type Simple struct{}

//Index ...
func (s *Simple) Index(in *proto.Request, out *proto.Response) {
	out.SetStatusCode(proto.StatusCodeOK)
	buf := bytes.Buffer{}
	buf.WriteString("<html><body><pre>")
	buf.WriteString("UUID: " + in.GetUUID() + "\n")
	buf.WriteString("URL: " + in.URL.String() + "\n")
	buf.WriteString("Version: " + strconv.FormatUint(uint64(in.GetVersion()), 10) + "\n")
	buf.WriteString("Meta: " + "\n")
	for k := range in.Meta {
		buf.WriteString(" - " + k + ": " + in.Meta.Get(k) + "\n")
	}
	buf.WriteString("</pre></body></html>")
	out.Body = buf.Bytes()
}

func main() {
	prt := proto.NewServer()
	prt.Handler("/", 1, (&Simple{}).Index)

	conf := &websrv.Config{
		HTTP: httpsrv.ConfigItem{Addr: "localhost:8090"},
		Headers: websrv.Headers{
			ProxyHeaders:   []string{"X-Forwarded-For", "Accept-Language", "User-Agent"},
			DefaultHeaders: map[string]string{"Content-Type": "text/html"},
		},
	}

	srv := websrv.NewServer(conf, prt, logger.Default())
	if err := srv.Up(); err != nil {
		panic(err)
	}

	<-time.After(time.Minute)

	if err := srv.Down(); err != nil {
		panic(err)
	}
}
