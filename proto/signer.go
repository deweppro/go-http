/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package proto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"sync"
)

type (
	//Signer ...
	Signer struct {
		id string
		h  hash.Hash
		l  sync.Mutex
	}
)

//NewSigner ...
func NewSigner(id, secret string) *Signer {
	return NewCustomSigner(id, secret, sha256.New)
}

//NewCustomSigner ...
func NewCustomSigner(id, secret string, h func() hash.Hash) *Signer {
	return &Signer{
		id: id,
		h:  hmac.New(h, []byte(secret)),
	}
}

//ID ...
func (s *Signer) ID() string {
	return s.id
}

//Create ...
func (s *Signer) Create(b []byte) []byte {
	s.l.Lock()
	defer func() {
		s.h.Reset()
		s.l.Unlock()
	}()
	s.h.Write(b)
	return s.h.Sum(nil)
}

//CreateString ...
func (s *Signer) CreateString(b []byte) string {
	return hex.EncodeToString(s.Create(b))
}

//Validate ...
func (s *Signer) Validate(b []byte, ex string) bool {
	return hmac.Equal(s.Create(b), []byte(ex))
}
