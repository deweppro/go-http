/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	versionKey         = `Accept`
	versionValueRegexp = `application\/vnd.v(\d+)\+json`
	contentTypeKey     = `Content-Type`
	contentTypeJson    = `application/json; charset=utf-8`
	signKey            = `Signature`
	signValueRegexp    = `keyId=\"(.*)\",algorithm=\"(.*)\",signature=\"(.*)\"`
)

var (
	vercomp  = regexp.MustCompile(versionValueRegexp)
	signcomp = regexp.MustCompile(signValueRegexp)

	ErrInvalidSignature = errors.New(`invalid signature format`)
)

type (
	//Headers model
	Headers map[string]string
	//Context model
	Context struct {
		Writer http.ResponseWriter
		Reader *http.Request
	}
	//Sign model
	Sign struct {
		ID        string
		Algorithm string
		Signature string
	}
	//Decoder interface
	Decoder func(data []byte, v interface{}) error
	//Encoder interface
	Encoder func(v interface{}) ([]byte, error)
)

//Decode decode request body
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

//Write make raw response
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

//JSON make json response
func (c *Context) JSON(code int, model json.Marshaler, heads Headers) error {
	body, err := model.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}
	if heads != nil {
		for k, v := range heads {
			c.Writer.Header().Set(k, v)
		}
	}
	c.Writer.Header().Set(contentTypeKey, contentTypeJson)
	c.Writer.WriteHeader(code)
	_, err = c.Writer.Write(body)
	return errors.Wrap(err, "context write")
}

//Redirect redirect to any url
func (c *Context) Redirect(url string) error {
	c.Writer.Header().Set("Location", url)
	c.Writer.WriteHeader(http.StatusMovedPermanently)
	_, err := c.Writer.Write([]byte{})
	return errors.Wrap(err, "context redirect")
}

//SetCookie setting cookie for response
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

//GetCookies getting cookie of request
func (c *Context) GetCookies() map[string]*http.Cookie {
	result := make(map[string]*http.Cookie)
	for _, v := range c.Reader.Cookies() {
		result[v.Name] = v
	}
	return result
}

//Version getting version of request
func (c *Context) Version() uint64 {
	d := c.Reader.Header.Get(versionKey)
	result := vercomp.FindSubmatch([]byte(d))
	if len(result) != 2 {
		return DefaultVersion
	}
	v, err := strconv.ParseUint(string(result[1]), 10, 32)
	if err != nil || v < 1 {
		return DefaultVersion
	}
	return v
}

//Sign getting signature of request
func (c *Context) Sign() (s Sign, err error) {
	d := c.Reader.Header.Get(signKey)
	r := signcomp.FindSubmatch([]byte(d))
	if len(r) != 4 {
		err = ErrInvalidSignature
		return
	}
	s.ID, s.Algorithm, s.Signature = string(r[1]), string(r[2]), string(r[3])
	return
}
