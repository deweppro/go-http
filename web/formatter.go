/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	_ FormatterFunc = TextFormatter
	_ FormatterFunc = RpcFormatter
	_ FormatterFunc = JSONFormatter
	_ FormatterFunc = EmptyFormatter
)

type FormatterFunc func(m *Context, code int, headers Headers, body interface{}) error

func EmptyFormatter(*Context, int, Headers, interface{}) error {
	return nil
}

func TextFormatter(m *Context, code int, headers Headers, body interface{}) error {
	var result []byte
	switch body.(type) {
	case error:
		result = []byte(body.(error).Error())
	case []byte:
		result = body.([]byte)
	case nil:
		result = []byte{}
	default:
		result = []byte(fmt.Sprintf("%+v", body))
	}

	if headers != nil {
		for k, v := range headers {
			m.Writer.Header().Set(k, v)
		}
	}
	m.Writer.WriteHeader(code)
	_, err := m.Writer.Write(result)
	return err
}

func RpcFormatter(m *Context, code int, headers Headers, body interface{}) error {
	var (
		result []byte
		err    error
	)
	switch body.(type) {
	case RpcResponseErrorModel:
		result, err = json.Marshal(body.(RpcResponseErrorModel))
	case RpcResponseModel:
		result, err = json.Marshal(body.(RpcResponseModel))
	case error:
		result, err = json.Marshal(RpcResponseErrorModel{
			ID: "",
			Error: RpcErrorBodyModel{
				Code:    code,
				Message: body.(error).Error(),
			},
		})
	case []byte:
		result, err = json.Marshal(RpcResponseModel{
			ID:     "",
			Result: string(body.([]byte)),
		})
	case json.Marshaler:
		result, err = json.Marshal(RpcResponseModel{
			ID:     "",
			Result: body.(json.Marshaler),
		})
	default:
		result, err = json.Marshal(RpcResponseModel{
			ID:     "",
			Result: body,
		})
	}

	if err != nil {
		return err
	}

	m.Writer.Header().Set("Content-Type", "application/json")
	if headers != nil {
		for k, v := range headers {
			m.Writer.Header().Set(k, v)
		}
	}
	m.Writer.WriteHeader(http.StatusOK)
	_, err = m.Writer.Write(result)
	return err
}

func JSONFormatter(m *Context, code int, headers Headers, body interface{}) error {
	var (
		result []byte
		err    error
	)
	switch body.(type) {
	case ResponseErrorModel:
		result, err = json.Marshal(body.(ResponseErrorModel))
	case ResponseModel:
		result, err = json.Marshal(body.(ResponseModel))
	case error:
		result, err = json.Marshal(ResponseErrorModel{
			Error: body.(error).Error(),
		})
	case []byte:
		result, err = json.Marshal(ResponseModel{
			Data: string(body.([]byte)),
		})
	case json.Marshaler:
		result, err = json.Marshal(ResponseModel{
			Data: body.(json.Marshaler),
		})
	default:
		result, err = json.Marshal(ResponseModel{
			Data: body,
		})
	}

	if err != nil {
		return err
	}

	m.Writer.Header().Set("Content-Type", "application/json")
	if headers != nil {
		for k, v := range headers {
			m.Writer.Header().Set(k, v)
		}
	}
	m.Writer.WriteHeader(http.StatusOK)
	_, err = m.Writer.Write(result)
	return err
}
