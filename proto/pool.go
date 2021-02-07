/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package proto

import (
	"math/rand"
	"net/url"
	"time"
)

type (
	//Pool ...
	Pool struct {
		Items map[string]List `yaml:"services" json:"services"`
	}
	//List ...
	List []string
	//Pooler ...
	Pooler interface {
		Pool() (string, string, error)
	}
	//Configer ...
	Configer interface {
		Get(name string) List
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

//Get ...
func (c Pool) Get(name string) List {
	v, ok := c.Items[name]
	if !ok {
		return List{}
	}
	return v
}

//Pool ...
func (v List) Pool() (string, string, error) {
	l := len(v)
	if l == 0 {
		return "", "", ErrEmptyPool
	}
	u, err := url.Parse(v[rand.Intn(len(v))])
	if err != nil || len(u.Scheme) == 0 || len(u.Host) == 0 {
		return "", "", ErrInvalidPoolAddress
	}

	return u.Scheme, u.Host, nil
}
