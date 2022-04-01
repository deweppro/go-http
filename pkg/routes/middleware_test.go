package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnit_ThrottlingMiddleware(t *testing.T) {
	rec := httptest.NewRecorder()
	ThrottlingMiddleware(0)(func(w http.ResponseWriter, r *http.Request) {
	})(rec, &http.Request{})
	resp := rec.Result()
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
}
