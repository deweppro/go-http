/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package proto

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type (
	Response struct {
		Common
	}
)

func NewResponse() *Response {
	r := &Response{
		Common: Common{Meta: make(http.Header)},
	}
	return r
}

func (r *Response) WriteToHTTP(w http.ResponseWriter) error {
	for key := range r.Meta {
		w.Header().Set(key, r.Meta.Get(key))
	}
	w.WriteHeader(Code2HTTPCode(r.GetStatusCode()))
	_, err := w.Write(r.Body)
	return err
}

func (r *Response) UpdateFromHTTP(v *http.Response, headers ...string) (err error) {
	r.Meta = v.Header
	for _, item := range append(headers, defaultResponseHeaders...) {
		r.Meta.Set(item, v.Header.Get(item))
	}
	r.Body, err = Reader(v.Body)
	if r.GetStatusCode() != StatusCodeOK {
		switch true {
		case err != nil:
			err = errors.Wrap(err, "body read")
		case len(r.Body) > 0:
			err = fmt.Errorf("%s", r.Body)
		default:
			err = fmt.Errorf("%s", v.Status)
		}
	}
	return err
}

func (r *Response) GetStatusCode() uint {
	code := r.Meta.Get(StatusCodeKey)
	v, err := strconv.ParseUint(code, 10, 32)
	if err != nil {
		return StatusCodeFail
	}
	return uint(v)
}

func (r *Response) SetStatusCode(v uint) {
	r.Meta.Set(StatusCodeKey, strconv.FormatUint(uint64(v), 10))
}