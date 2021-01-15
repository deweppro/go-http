/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package proto

import (
	"fmt"
)

type (
	Server struct {
		routes map[string]SrvCaller
	}
	SrvCaller func(in *Request, out *Response)
)

func NewServer() *Server {
	return &Server{routes: make(map[string]SrvCaller)}
}

func (o *Server) Call(in *Request, out *Response) {
	if c, ok := o.routes[o.route(in.Path, in.GetVersion())]; ok {
		c(in, out)
		return
	}

	out.SetStatusCode(StatusCodeNotFound)
}

func (o *Server) Handler(path string, ver uint64, call SrvCaller) {
	o.routes[o.route(path, ver)] = call
}

func (o *Server) route(path string, ver uint64) string {
	return fmt.Sprintf("%s::%d", path, ver)
}
