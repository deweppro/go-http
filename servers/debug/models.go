/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package debug

import (
	"github.com/deweppro/go-http/v2/servers/http"
)

//go:generate easyjson

//easyjson:json
type (
	//DebugConfig ...
	Config struct {
		Debug http.ConfigItem `yaml:"debug" json:"debug"`
	}
)
