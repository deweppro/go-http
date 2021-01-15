/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package proto

import (
	"math/rand"
	"time"
)

type (
	Pool struct {
		Items map[string]List `yaml:"services" json:"services"`
	}
	List   []string
	Pooler interface {
		Pool() (string, error)
	}
	Configer interface {
		Get(name string) List
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func (c Pool) Get(name string) List {
	v, ok := c.Items[name]
	if !ok {
		return List{}
	}
	return v
}

func (v List) Pool() (string, error) {
	l := len(v)
	if l == 0 {
		return "", ErrEmptyPool
	}
	return v[rand.Intn(len(v))], nil
}
