/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package proto

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type (
	Request struct {
		Common
		URL *url.URL
	}
)

func NewRequest() *Request {
	r := &Request{
		Common: Common{Meta: make(http.Header)},
	}
	return r
}

func (r *Request) UpdateFromHTTP(v *http.Request, headers ...string) (err error) {
	r.URL = v.URL
	r.Meta = make(http.Header)
	for _, item := range append(headers, defaultRequestHeaders...) {
		r.Meta.Set(item, v.Header.Get(item))
	}
	r.Body, err = Reader(v.Body)
	return
}

func (r *Request) GetVersion() uint64 {
	d := r.Meta.Get(VersionKey)
	result := vercomp.FindSubmatch([]byte(d))
	if len(result) != 2 {
		return versionDefault
	}
	v, err := strconv.ParseUint(string(result[1]), 10, 32)
	if err != nil || v < 1 {
		return versionDefault
	}
	return v
}

func (r *Request) SetVersion(v uint64) {
	r.Meta.Set(VersionKey, fmt.Sprintf(versionValueTmpl, v))
}
