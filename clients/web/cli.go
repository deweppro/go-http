/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package web

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	proto "github.com/deweppro/go-http/v2"
	"github.com/pkg/errors"
)

type (
	Client struct {
		cli     *http.Client
		headers http.Header
		sign    *proto.Signer
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

//WithSign sign request
func (v *Client) WithSign(s *proto.Signer) {
	v.sign = s
}

//Call make request to server
func (v *Client) Call(pool proto.Pooler, in *proto.Request, out *proto.Response) error {
	add, err := pool.Pool()
	if err != nil {
		return errors.Wrap(err, "get address from pool")
	}
	req, err := http.NewRequest(http.MethodPost, add+"/"+strings.TrimLeft(in.Path, "/"), bytes.NewReader(in.Body))
	if err != nil {
		return errors.Wrap(err, "create request")
	}

	in.SetVersion(in.GetVersion())
	if v.sign != nil {
		in.CreateSign(v.sign)
	}
	in.SetUUID(in.CreateUUID())

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
			out.GetStatusCode(), in.Path, in.GetVersion(), in.GetUUID(),
			in.Meta.Get(proto.SignKey), in.Body, out.Body, err,
		)
	}
	return err
}
