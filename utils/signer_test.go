package utils_test

import (
	"testing"

	"github.com/deweppro/go-http/utils"
	"github.com/stretchr/testify/require"
)

func TestUnit_Signature(t *testing.T) {
	sign := utils.NewSHA256("123", "456")

	body := []byte("hello")
	hash := "b7089b0463bf766946fc467102671dbe91659f17a7a19145cd68138c36b00555"

	require.Equal(t, "123", sign.ID())
	require.Equal(t, hash, sign.CreateString(body))
	require.True(t, sign.Validate(body, hash))
}

func TestUnit_SignatureStorage(t *testing.T) {
	store := utils.NewSignatureStorage()

	store.Add(utils.NewSHA256("1", "0"))
	store.Add(utils.NewSHA256("2", "0"))
	store.Add(utils.NewSHA256("3", "0"))
	store.Add(utils.NewSHA256("5", "0"))
	require.Equal(t, 4, store.Count())

	store.Add(utils.NewMD5("5", "0"))
	require.Equal(t, 4, store.Count())

	require.Nil(t, store.Get("4"))
	s := store.Get("5")
	require.NotNil(t, s)
	require.Equal(t, "5", s.ID())
	require.Equal(t, "hmac-md5", s.Algorithm())
}
