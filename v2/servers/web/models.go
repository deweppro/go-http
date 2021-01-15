/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package web

import (
	"github.com/deweppro/go-http/v2/servers/http"
)

//go:generate easyjson

//easyjson:json
type (
	//Config ...
	Config struct {
		HTTP    http.ConfigItem `yaml:"http" json:"http"`
		Headers Headers         `yaml:"headers" json:"headers"`
	}
	//Headers ...
	Headers struct {
		ProxyHeaders   []string          `yaml:"proxy_headers" json:"proxy_headers"`
		DefaultHeaders map[string]string `yaml:"default_headers" json:"default_headers"`
	}
)
