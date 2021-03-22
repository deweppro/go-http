package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deweppro/go-http/web/routes"
	"github.com/stretchr/testify/require"
)

func TestUnit_ThrottlingMiddleware(t *testing.T) {
	rec := httptest.NewRecorder()
	routes.ThrottlingMiddleware(0)(func(w http.ResponseWriter, r *http.Request) {
	})(rec, &http.Request{})
	resp := rec.Result()
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
}
