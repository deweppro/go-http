package routes

import (
	"net/http"
	"testing"

	"github.com/deweppro/go-http/internal"

	"github.com/stretchr/testify/require"
)

func TestUnit_NewHandler(t *testing.T) {
	h := newHandler()
	h.Route("/aaa/{id}", func(_ http.ResponseWriter, _ *http.Request) {}, []string{http.MethodPost})
	h.Route("", func(_ http.ResponseWriter, _ *http.Request) {}, []string{http.MethodPost})

	code, ctrl, vars, midd := h.Match("/aaa/bbb", http.MethodPost)
	require.Equal(t, 200, code)
	require.NotNil(t, ctrl)
	require.Equal(t, 0, len(midd))
	require.Equal(t, internal.VarsData{"id": "bbb"}, vars)

	h.Middlewares("/aaa", ThrottlingMiddleware(0))
	h.Middlewares("", ThrottlingMiddleware(0))

	code, ctrl, vars, midd = h.Match("/aaa/ccc", http.MethodGet)
	require.Equal(t, http.StatusMethodNotAllowed, code)
	require.Nil(t, ctrl)
	require.Equal(t, 1, len(midd))
	require.Equal(t, internal.VarsData(nil), vars)

	code, ctrl, vars, midd = h.Match("/aaa/bbb", http.MethodPost)
	require.Equal(t, http.StatusOK, code)
	require.NotNil(t, ctrl)
	require.Equal(t, 2, len(midd))
	require.Equal(t, internal.VarsData{"id": "bbb"}, vars)

	code, ctrl, vars, midd = h.Match("", http.MethodPost)
	require.Equal(t, http.StatusOK, code)
	require.NotNil(t, ctrl)
	require.Equal(t, 1, len(midd))
	require.Equal(t, internal.VarsData{}, vars)

	h.Middlewares("/www/www/www", ThrottlingMiddleware(0))

	code, ctrl, vars, midd = h.Match("/www/www/www", http.MethodPost)
	require.Equal(t, http.StatusNotFound, code)
	require.Nil(t, ctrl)
	require.Equal(t, 1, len(midd))
	require.Equal(t, internal.VarsData(nil), vars)

	code, ctrl, vars, midd = h.Match("/test", http.MethodGet)
	require.Equal(t, http.StatusNotFound, code)
	require.Nil(t, ctrl)
	require.Equal(t, 1, len(midd))
	require.Equal(t, internal.VarsData(nil), vars)

	h.NoFoundHandler(func(_ http.ResponseWriter, _ *http.Request) {})

	code, ctrl, vars, midd = h.Match("/test", http.MethodGet)
	require.Equal(t, http.StatusOK, code)
	require.NotNil(t, ctrl)
	require.Equal(t, 1, len(midd))
	require.Equal(t, internal.VarsData(nil), vars)
}

func TestUnit_NewHandler2(t *testing.T) {
	h := newHandler()
	h.Route("/api/v1/data", func(_ http.ResponseWriter, _ *http.Request) {}, []string{http.MethodGet})

	h.Middlewares("/api/v1", ThrottlingMiddleware(0))

	code, ctrl, vars, midd := h.Match("/api/v1/data", http.MethodGet)
	require.Equal(t, http.StatusOK, code)
	require.NotNil(t, ctrl)
	require.Equal(t, 1, len(midd))
	require.Equal(t, internal.VarsData{}, vars)

}
