/*
 * Copyright (c) 2021 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package go_http

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_PoolConfig(t *testing.T) {
	data := `{"services":{"server1":["127.0.0.1"], "server2":[]}}`

	c := Pool{}
	require.NoError(t, json.Unmarshal([]byte(data), &c))

	v, err := c.Get("server1").Pool()
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1", v)

	_, err = c.Get("server2").Pool()
	require.Error(t, err)

	_, err = c.Get("server3").Pool()
	require.Error(t, err)
}
