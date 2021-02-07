/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package proto

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

//nolint: golint
const (
	VersionKey                = `Accept`
	versionValueRegexp        = `application\/vnd.v(\d+)\+json`
	versionValueTmpl          = `application/vnd.v%d+json`
	versionDefault     uint64 = 1

	UUIDKey       = `X-Request-ID`
	StatusCodeKey = `X-Status`

	SignKey         = `Signature`
	signValueRegexp = `keyId=\"(.*)\",algorithm=\"(.*)\",signature=\"(.*)\"`
	signValueTmpl   = `keyId="%s",algorithm="hmac-sha256",signature="%s"`

	ContentTypeKey         = `Content-Type`
	ContentTypeJSONValue   = `application/json; charset=utf-8`
	ContentTypeBinaryValue = `application/octet-stream`

	StatusCodeFail         uint = 0
	StatusCodeOK           uint = 1
	StatusCodeNotFound     uint = 2
	StatusCodeUnauthorized uint = 3
	StatusCodeForbidden    uint = 4
	StatusCodeServerError  uint = 5
	StatusCodeRedirect     uint = 6
)

var (
	vercomp                = regexp.MustCompile(versionValueRegexp)
	signcomp               = regexp.MustCompile(signValueRegexp)
	defaultRequestHeaders  = []string{VersionKey, UUIDKey, SignKey}
	defaultResponseHeaders = []string{UUIDKey, StatusCodeKey, SignKey, ContentTypeKey}
)

type (
	//Common ...
	Common struct {
		cookies map[string]*http.Cookie
		Meta    http.Header
		Body    []byte
	}
	//Sign ...
	Sign struct {
		ID        string
		Algorithm string
		Signature string
	}
)

//SetMeta ...
func (o *Common) SetMeta(m map[string]string) {
	if m != nil {
		for k, v := range m {
			o.Meta.Set(k, v)
		}
	}
}

//CreateSign ...
func (o *Common) CreateSign(s *Signer) {
	o.Meta.Set(SignKey, fmt.Sprintf(signValueTmpl, s.ID(), s.CreateString(o.Body)))
}

//ValidateSign ...
func (o *Common) ValidateSign(s *Signer) bool {
	sign, err := o.GetSignature()
	if err != nil {
		return false
	}
	if sign.ID != s.ID() {
		return false
	}
	return s.Validate(o.Body, sign.Signature)
}

//GetSignature ...
func (o *Common) GetSignature() (s Sign, err error) {
	d := o.Meta.Get(SignKey)
	r := signcomp.FindSubmatch([]byte(d))
	if len(r) != 4 {
		err = ErrInvalidSignature
		return
	}
	s.ID, s.Algorithm, s.Signature = string(r[1]), string(r[2]), string(r[3])
	return
}

//GetUUID ...
func (o *Common) GetUUID() string {
	return o.Meta.Get(UUIDKey)
}

//UpdateUUID ...
func (o *Common) UpdateUUID() {
	o.SetUUID(CreateUUID())
}

//SetUUID ...
func (o *Common) SetUUID(v string) {
	o.Meta.Set(UUIDKey, v)
}

//DecodeJSON ...
func (o *Common) DecodeJSON(v interface{}) error {
	return json.Unmarshal(o.Body, v)
}

//EncodeJSON ...
func (o *Common) EncodeJSON(v interface{}) (err error) {
	o.Meta.Set(ContentTypeKey, ContentTypeJSONValue)
	o.Body, err = json.Marshal(v)
	return
}

//DecodeGob ...
func (o *Common) DecodeGob(v interface{}) error {
	buf := bytes.NewBuffer(o.Body)
	return gob.NewDecoder(buf).Decode(v)
}

//EncodeGob ...
func (o *Common) EncodeGob(v interface{}) (err error) {
	o.Meta.Set(ContentTypeKey, ContentTypeBinaryValue)
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(v)
	o.Body = buf.Bytes()
	return
}

//GetCookie ...
func (o *Common) GetCookie(name string) (*http.Cookie, error) {
	if v, ok := o.cookies[name]; ok {
		return v, nil
	}
	return nil, ErrCookieNotFound
}

//SetCookie ...
func (o *Common) SetCookie(v ...*http.Cookie) {
	for _, item := range v {
		o.cookies[item.Name] = item
	}
}

//Code2HTTPCode ...
func Code2HTTPCode(v uint) int {
	switch v {
	case StatusCodeFail:
		return http.StatusBadRequest
	case StatusCodeOK:
		return http.StatusOK
	case StatusCodeNotFound:
		return http.StatusNotFound
	case StatusCodeUnauthorized:
		return http.StatusUnauthorized
	case StatusCodeForbidden:
		return http.StatusForbidden
	case StatusCodeServerError:
		return http.StatusInternalServerError
	case StatusCodeRedirect:
		return http.StatusMovedPermanently
	default:
		return http.StatusOK
	}
}

//HTTPCode2Code ...
func HTTPCode2Code(v int) uint {
	switch v {
	case http.StatusBadRequest:
		return StatusCodeFail
	case http.StatusOK:
		return StatusCodeOK
	case http.StatusNotFound:
		return StatusCodeNotFound
	case http.StatusUnauthorized:
		return StatusCodeUnauthorized
	case http.StatusForbidden:
		return StatusCodeForbidden
	case http.StatusInternalServerError:
		return StatusCodeServerError
	case http.StatusMovedPermanently:
		return StatusCodeRedirect
	default:
		return StatusCodeOK
	}
}