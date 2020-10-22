/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package web

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type (
	Headers map[string]string
	Context struct {
		Writer    http.ResponseWriter
		Reader    *http.Request
		Formatter FormatterFunc
	}
)

func (_ctx *Context) Empty(code int) error {
	return errors.Wrap(_ctx.Formatter(_ctx, code, nil, nil), "context empty")
}

func (_ctx *Context) Error(code int, err error) error {
	return errors.Wrap(_ctx.Formatter(_ctx, code, nil, err), "context error")
}

func (_ctx *Context) Decode(call func(http.Header, []byte) error) error {
	data, err := ioutil.ReadAll(_ctx.Reader.Body)
	if err != nil {
		return errors.Wrap(err, "context decode body read")
	}
	if err = _ctx.Reader.Body.Close(); err != nil {
		return errors.Wrap(err, "context decode body close")
	}
	return errors.Wrap(call(_ctx.Reader.Header, data), "context decode call")
}

func (_ctx *Context) Encode(call func() (int, Headers, interface{})) error {
	code, headers, body := call()
	return errors.Wrap(_ctx.Formatter(_ctx, code, headers, body), "context encode formatter")
}

func (_ctx *Context) Redirect(url string) error {
	_ctx.Writer.Header().Set("Location", url)
	_ctx.Writer.WriteHeader(http.StatusMovedPermanently)
	_, err := _ctx.Writer.Write([]byte{})
	return errors.Wrap(err, "context redirect")
}

func (_ctx *Context) SetCookie(key, value string, ttl time.Duration) {
	http.SetCookie(_ctx.Writer, &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		Domain:   _ctx.Reader.Host,
		Expires:  time.Now().Add(ttl),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (_ctx *Context) GetCookies() map[string]*http.Cookie {
	result := make(map[string]*http.Cookie)
	for _, c := range _ctx.Reader.Cookies() {
		result[c.Name] = c
	}
	return result
}
