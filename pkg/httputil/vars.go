package httputil

import (
	"net/http"
	"strconv"

	"github.com/deweppro/go-http/internal"
	"github.com/deweppro/go-http/pkg/errs"
)

func VarsString(r *http.Request, key string) (string, error) {
	if v := r.Context().Value(internal.VarsKey(key)); v != nil {
		return v.(string), nil
	}
	return "", errs.ErrFailContextKey
}

func VarsInt64(r *http.Request, key string) (int64, error) {
	v, err := VarsString(r, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(v, 10, 64)
}

func VarsFloat64(r *http.Request, key string) (float64, error) {
	v, err := VarsString(r, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(v, 64)
}
