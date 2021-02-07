/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package debugsrv

import (
	"github.com/deweppro/go-http/v2/servers/httpsrv"
)

//go:generate easyjson

//easyjson:json
type (
	//Config model
	Config struct {
		Debug httpsrv.ConfigItem `yaml:"debug" json:"debug"`
	}
)
