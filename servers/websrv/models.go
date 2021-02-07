/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package websrv

import (
	"github.com/deweppro/go-http/v2/servers/httpsrv"
)

//go:generate easyjson

//easyjson:json
type (
	//Config model
	Config struct {
		HTTP    httpsrv.ConfigItem `yaml:"http" json:"http"`
		Headers Headers            `yaml:"headers" json:"headers"`
	}
	//Headers model
	Headers struct {
		ProxyHeaders   []string          `yaml:"proxy_headers" json:"proxy_headers"`
		DefaultHeaders map[string]string `yaml:"default_headers" json:"default_headers"`
	}
)
