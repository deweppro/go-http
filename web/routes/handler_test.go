package routes_test

import (
	"net/http"
	"testing"

	"github.com/deweppro/go-http/web/routes"
	"github.com/stretchr/testify/require"
)

func TestUnit_NewHandler(t *testing.T) {
	h := routes.NewHandler()
	h.Route(routes.SplitURI("/aaa/bbb"), 0, func(_ http.ResponseWriter, _ *http.Request) {}, []string{"post"})
	h.Route(nil, 0, func(_ http.ResponseWriter, _ *http.Request) {}, []string{"post"})

	code, c, m := h.Match(routes.SplitURI("/aaa/bbb"), 0, http.MethodPost)
	require.Equal(t, 200, code)
	require.NotNil(t, c)
	require.Equal(t, 0, len(m))

	h.Middlewares(routes.SplitURI("/aaa"), 0, []routes.MiddlFunc{routes.ThrottlingMiddleware(0)})

	code, c, m = h.Match(routes.SplitURI("/aaa/bbb"), 0, http.MethodGet)
	require.Equal(t, 405, code)
	require.Nil(t, c)
	require.Equal(t, 1, len(m))

	code, c, m = h.Match(routes.SplitURI("/test"), 0, http.MethodGet)
	require.Equal(t, 404, code)
	require.Nil(t, c)
	require.Equal(t, 0, len(m))

	code, c, m = h.Match(nil, 0, http.MethodPost)
	require.Equal(t, 400, code)
	require.Nil(t, c)
	require.Equal(t, 0, len(m))
}
