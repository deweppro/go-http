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
	"strconv"
	"time"

	"github.com/google/uuid"
)

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
)

var (
	vercomp                = regexp.MustCompile(versionValueRegexp)
	signcomp               = regexp.MustCompile(signValueRegexp)
	defaultRequestHeaders  = []string{VersionKey, UUIDKey, SignKey}
	defaultResponseHeaders = []string{UUIDKey, StatusCodeKey, SignKey, ContentTypeKey}
)

type (
	Common struct {
		Meta http.Header
		Body []byte
	}

	Sign struct {
		ID        string
		Algorithm string
		Signature string
	}
)

func (o *Common) SetMeta(m map[string]string) {
	if m != nil {
		for k, v := range m {
			o.Meta.Set(k, v)
		}
	}
}

func (o *Common) CreateSign(s *Signer) {
	o.Meta.Set(SignKey, fmt.Sprintf(signValueTmpl, s.ID(), s.CreateString(o.Body)))
}

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

func (o *Common) CreateUUID() string {
	rnd, err := uuid.NewRandom()
	if err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 16)
	}
	return rnd.String()
}

func (o *Common) GetUUID() (v string) {
	v = o.Meta.Get(UUIDKey)
	if len(v) > 0 {
		return
	}
	v = o.CreateUUID()
	o.SetUUID(v)
	return
}

func (o *Common) SetUUID(v string) {
	o.Meta.Set(UUIDKey, v)
}

func (o *Common) DecodeJSON(v interface{}) error {
	return json.Unmarshal(o.Body, v)
}

func (o *Common) EncodeJSON(v interface{}) (err error) {
	o.Meta.Set(ContentTypeKey, ContentTypeJSONValue)
	o.Body, err = json.Marshal(v)
	return
}

func (o *Common) DecodeGob(v interface{}) error {
	buf := bytes.NewBuffer(o.Body)
	return gob.NewDecoder(buf).Decode(v)
}

func (o *Common) EncodeGob(v interface{}) (err error) {
	o.Meta.Set(ContentTypeKey, ContentTypeBinaryValue)
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(v)
	o.Body = buf.Bytes()
	return
}

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
	default:
		return http.StatusOK
	}
}

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
	default:
		return StatusCodeOK
	}
}
