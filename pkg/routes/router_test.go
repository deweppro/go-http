package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deweppro/go-http/pkg/routes"
	"github.com/stretchr/testify/require"
)

func TestUnit_Route1(t *testing.T) {
	result := new(string)
	r := routes.NewRouter()
	r.Global(func(c routes.CtrlFunc) routes.CtrlFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			*result += "1"
			c(w, r)
		}
	})
	r.Global(func(c routes.CtrlFunc) routes.CtrlFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			*result += "2"
			c(w, r)
		}
	})
	r.Global(func(c routes.CtrlFunc) routes.CtrlFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			*result += "3"
			c(w, r)
		}
	})
	r.Route("/", func(w http.ResponseWriter, r *http.Request) {
		*result += "Ctrl"
	}, http.MethodGet)
	r.Middlewares("/test", func(c routes.CtrlFunc) routes.CtrlFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			*result += "4"
			c(w, r)
		}
	})
	r.Middlewares("/", func(c routes.CtrlFunc) routes.CtrlFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			*result += "5"
			c(w, r)
		}
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	require.Equal(t, "1235Ctrl", *result)
}

func TestUnit_Route2(t *testing.T) {
	r := routes.NewRouter()
	r.Route("/{id}", func(w http.ResponseWriter, r *http.Request) {}, http.MethodGet)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/aaa/bbb/ccc/eee/ggg/fff/kkk", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, 404, w.Result().StatusCode)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/aaa/", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Result().StatusCode)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/aaa", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Result().StatusCode)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/aaa?a=1", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Result().StatusCode)
}

func BenchmarkRouter(b *testing.B) {
	b.ReportAllocs()

	serv := routes.NewRouter()
	serv.Route("/{id0}/{id1}/{id2}/{id3}/{id4}/{id5}/{id6}", func(w http.ResponseWriter, r *http.Request) {}, http.MethodGet)
	r := httptest.NewRequest("GET", "/aaa/bbb/ccc/eee/ggg/fff/kkk", nil)

	b.ResetTimer()
	b.ReportAllocs()

	b.Run("", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			w := httptest.NewRecorder()
			serv.ServeHTTP(w, r)
			if w.Result().StatusCode != http.StatusOK {
				b.Fatalf("invalid code: %d", w.Result().StatusCode)
			}
			w.Flush()
		}
	})
}
