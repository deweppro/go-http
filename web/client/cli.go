/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/deweppro/go-http/web"
	"github.com/pkg/errors"
)

const (
	version = `application/vnd.v%d+json`
)

type (
	//Client client connecting to the server
	Client struct {
		cli     *http.Client
		headers web.Headers
		debug   bool
		writer  io.Writer
	}
)

//NewClient init client
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

//NewCustomClient init client
func NewCustomClient(cli *http.Client) *Client {
	return &Client{
		cli: cli,
	}
}

//Debug enable logging of responses
func (v *Client) Debug(is bool, w io.Writer) {
	v.debug, v.writer = is, w
}

//Debug enable logging of responses
func (v *Client) WithHeaders(heads web.Headers) {
	v.headers = heads
}

//Call make request to server
func (v *Client) Call(method, url string, ver uint64, in json.Marshaler, out json.Unmarshaler) (int, error) {
	var (
		body []byte
		err  error
	)

	if in != nil {
		body, err = in.MarshalJSON()
		if err != nil {
			return 0, errors.Wrap(err, "marshal request")
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return 0, errors.Wrap(err, "create request")
	}

	req.Header.Set("Connection", "keep-alive")
	if v.headers != nil {
		for k, v := range v.headers {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Accept", fmt.Sprintf(version, ver))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := v.cli.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "make request")
	}

	code := resp.StatusCode
	switch code {
	case 200:
		body, err = v.readBody(resp.Body, out)
	default:
		body, err = v.readBody(resp.Body, nil)
		if err == nil {
			err = errors.New(string(body))
		}
	}

	v.writeDebug(code, method, url, ver, body, err)

	switch err {
	case nil:
		return code, nil
	case io.EOF:
		return code, errors.New(resp.Status)
	default:
		return code, err
	}
}

func (v *Client) readBody(rc io.ReadCloser, resp json.Unmarshaler) ([]byte, error) {
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, errors.Wrap(err, "read response")
	}
	if resp != nil {
		err = resp.UnmarshalJSON(b)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal response")
		}
	}
	return b, nil
}

func (v *Client) writeDebug(code int, method, path string, ver uint64, body []byte, err error) {
	if v.debug {
		fmt.Fprintf(v.writer, "[%d] %s:%s ver:%d err: %+v raw:%s \n", code, method, path, ver, err, body)
	}
}
