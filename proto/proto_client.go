/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package proto

import "github.com/pkg/errors"

type (
	//Client ...
	Client struct {
		conf Configer
		cli  map[string]CliCaller
	}
	//CliCaller ...
	CliCaller func(pool Pooler, in *Request, out *Response) error
)

//NewClient ...
func NewClient(conf Configer) *Client {
	return &Client{
		conf: conf,
		cli:  make(map[string]CliCaller),
	}
}

//Call ...
func (o *Client) Call(name string, in *Request, out *Response) error {
	if c, ok := o.cli[name]; ok {
		return c(o.conf.Get(name), in, out)
	}
	return errors.Wrap(ErrClientNotFound, name)
}

//Client ...
func (o *Client) Client(name string, call CliCaller) {
	o.cli[name] = call
}
