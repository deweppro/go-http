/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package epoll

//go:generate easyjson

//easyjson:json
type (
	//EpollConfig ...
	EpollConfig struct {
		Epoll ConfigItem `yaml:"epoll" json:"epoll"`
	}
	//ConfigItem ...
	ConfigItem struct {
		Addr string `yaml:"addr" json:"addr"`
	}
)
