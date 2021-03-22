package utils_test

import (
	"encoding/json"
	"testing"

	"github.com/deweppro/go-http/utils"
	"github.com/stretchr/testify/require"
)

func TestUnit_PoolConfig(t *testing.T) {
	data := `{"services":{"server0":["http://127.0.0.1"], "server1":["127.0.0.1"], "server2":[]}}`

	c := utils.Pool{}
	require.NoError(t, json.Unmarshal([]byte(data), &c))

	uri, err := c.Get("server0").Pool()
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1", uri.Host)
	require.Equal(t, "http", uri.Scheme)

	uri, err = c.Get("server1").Pool()
	require.Error(t, err)
	require.EqualError(t, utils.ErrInvalidPoolAddress, err.Error())

	_, err = c.Get("server2").Pool()
	require.Error(t, err)
	require.EqualError(t, utils.ErrEmptyPool, err.Error())

	_, err = c.Get("server3").Pool()
	require.Error(t, err)
	require.EqualError(t, utils.ErrEmptyPool, err.Error())
}
