/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package proto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_PoolConfig(t *testing.T) {
	data := `{"services":{"server0":["http://127.0.0.1"], "server1":["127.0.0.1"], "server2":[]}}`

	c := Pool{}
	require.NoError(t, json.Unmarshal([]byte(data), &c))

	s, v, err := c.Get("server0").Pool()
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1", v)
	require.Equal(t, "http", s)

	s, v, err = c.Get("server1").Pool()
	require.Error(t, err)
	require.EqualError(t, ErrInvalidPoolAddress, err.Error())

	_, _, err = c.Get("server2").Pool()
	require.Error(t, err)
	require.EqualError(t, ErrEmptyPool, err.Error())

	_, _, err = c.Get("server3").Pool()
	require.Error(t, err)
	require.EqualError(t, ErrEmptyPool, err.Error())
}
