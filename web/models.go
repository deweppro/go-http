/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package web

//go:generate easyjson

import (
	"time"
)

//easyjson:json
type (
	/*
		http server config
	*/
	DebugConfig struct {
		Debug ConfigItem `yaml:"debug" json:"debug"`
	}
	HTTPConfig struct {
		HTTP ConfigItem `yaml:"http" json:"http"`
	}
	ConfigItem struct {
		Addr         string        `yaml:"addr" json:"addr"`
		ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
		IdleTimeout  time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
	}

	/*
		json rpc model
	*/
	RpcRequestModel struct {
		ID     string      `json:"id"`
		Method string      `json:"method"`
		Params interface{} `json:"params,omitempty"`
	}
	RpcResponseModel struct {
		ID     string      `json:"id"`
		Result interface{} `json:"result"`
	}
	RpcResponseErrorModel struct {
		ID    string            `json:"id"`
		Error RpcErrorBodyModel `json:"error"`
	}
	RpcErrorBodyModel struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	/*
		json model
	*/
	ResponseModel struct {
		Data interface{} `json:"data"`
	}
	ResponseErrorModel struct {
		Error string `json:"error"`
	}
)
