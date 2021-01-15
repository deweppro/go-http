/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package go_http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"sync"
)

type (
	Signer struct {
		id string
		h  hash.Hash
		l  sync.Mutex
	}
)

func NewSigner(id, secret string) *Signer {
	return NewCustomSigner(id, secret, sha256.New)
}

func NewCustomSigner(id, secret string, h func() hash.Hash) *Signer {
	return &Signer{
		id: id,
		h:  hmac.New(h, []byte(secret)),
	}
}

func (s *Signer) ID() string {
	return s.id
}

func (s *Signer) Create(b []byte) []byte {
	s.l.Lock()
	defer func() {
		s.h.Reset()
		s.l.Unlock()
	}()
	s.h.Write(b)
	return s.h.Sum(nil)
}

func (s *Signer) CreateString(b []byte) string {
	return hex.EncodeToString(s.Create(b))
}

func (s *Signer) Validate(b []byte, ex string) bool {
	return hmac.Equal(s.Create(b), []byte(ex))
}
