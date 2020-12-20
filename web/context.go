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
	//Headers ...
	Headers map[string]string
	//Context ...
	Context struct {
		Writer http.ResponseWriter
		Reader *http.Request
	}
	//Decoder ...
	Decoder func(data []byte, v interface{}) error
	//Encoder ...
	Encoder func(v interface{}) ([]byte, error)
)

//Decode ...
func (c *Context) Decode(model interface{}, call Decoder) error {
	data, err := ioutil.ReadAll(c.Reader.Body)
	if err != nil {
		return errors.Wrap(err, "context decode body read")
	}
	if err = c.Reader.Body.Close(); err != nil {
		return errors.Wrap(err, "context decode body close")
	}
	return errors.Wrap(call(data, model), "context decode call")
}

//Write ...
func (c *Context) Write(code int, body []byte, heads Headers) error {
	if heads != nil {
		for k, v := range heads {
			c.Writer.Header().Set(k, v)
		}
	}
	c.Writer.WriteHeader(code)
	_, err := c.Writer.Write(body)
	return errors.Wrap(err, "context write")
}

//Redirect ...
func (c *Context) Redirect(url string) error {
	c.Writer.Header().Set("Location", url)
	c.Writer.WriteHeader(http.StatusMovedPermanently)
	_, err := c.Writer.Write([]byte{})
	return errors.Wrap(err, "context redirect")
}

//SetCookie ...
func (c *Context) SetCookie(key, value string, ttl time.Duration) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		Domain:   c.Reader.Host,
		Expires:  time.Now().Add(ttl),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

//GetCookies ...
func (c *Context) GetCookies() map[string]*http.Cookie {
	result := make(map[string]*http.Cookie)
	for _, c := range c.Reader.Cookies() {
		result[c.Name] = c
	}
	return result
}
