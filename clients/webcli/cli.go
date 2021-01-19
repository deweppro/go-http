/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package webcli

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	proto "github.com/deweppro/go-http/v2"
	"github.com/pkg/errors"
)

type (
	Client struct {
		cli     *http.Client
		headers http.Header
		debug   bool
		writer  io.Writer
	}
)

func NewClient() *Client {
	cli := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
		Timeout: 5 * time.Second,
	}
	return NewCustomClient(cli)
}

func NewCustomClient(cli *http.Client) *Client {
	return &Client{
		cli:     cli,
		headers: make(http.Header),
	}
}

//Debug enable logging of responses
func (v *Client) Debug(is bool, w io.Writer) {
	v.debug, v.writer = is, w
}

func (v *Client) WithHeaders(heads http.Header) {
	v.headers = heads
}

//Call make request to server
func (v *Client) Call(pool proto.Pooler, in *proto.Request, out *proto.Response) error {
	addr, err := pool.Pool()
	if err != nil {
		return errors.Wrap(err, "get address from pool")
	}
	in.URL.Host = addr
	req, err := http.NewRequest(http.MethodPost, in.URL.String(), bytes.NewReader(in.Body))
	if err != nil {
		return errors.Wrap(err, "create request")
	}

	in.SetVersion(in.GetVersion())
	in.UpdateUUID()

	req.Header.Set("Connection", "keep-alive")
	for k := range v.headers {
		req.Header.Set(k, v.headers.Get(k))
	}
	for k := range in.Meta {
		req.Header.Set(k, in.Meta.Get(k))
	}

	resp, err := v.cli.Do(req)
	if err != nil {
		return errors.Wrap(err, "make request")
	}

	err = out.UpdateFromHTTP(resp)
	if v.debug {
		fmt.Fprintf(
			v.writer,
			"code<%d> path<%s> ver<%d> uuid<%s> sign<%s> in<%s> out<%s> err<%+v>\n",
			out.GetStatusCode(), in.URL.String(), in.GetVersion(), in.GetUUID(),
			in.Meta.Get(proto.SignKey), in.Body, out.Body, err,
		)
	}
	return err
}
