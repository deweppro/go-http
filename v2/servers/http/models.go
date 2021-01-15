/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package http

import (
	"time"
)

//go:generate easyjson

//easyjson:json
type (
	//Config ...
	Config struct {
		HTTP ConfigItem `yaml:"http" json:"http"`
	}
	//ConfigItem ...
	ConfigItem struct {
		Addr         string        `yaml:"addr" json:"addr"`
		ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
		IdleTimeout  time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
	}
)
