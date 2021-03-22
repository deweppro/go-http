package routes_test

import (
	"net/http"
	"testing"

	"github.com/deweppro/go-http/web/routes"
	"github.com/stretchr/testify/require"
)

func TestUnit_Route(t *testing.T) {
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
	r.ServeHTTP(nil, &http.Request{Method: "GET", RequestURI: "/"})
	require.Equal(t, "1235Ctrl", *result)
}

func BenchmarkRouter(b *testing.B) {
	b.ReportAllocs()

	r := routes.NewRouter()
	r.Route("/aaa/bbb/ccc/eee/ggg/fff/kkk", func(w http.ResponseWriter, r *http.Request) {}, http.MethodGet)
	req := &http.Request{Method: "GET", RequestURI: "/aaa/bbb/ccc/eee/ggg/fff/kkk"}

	b.ResetTimer()

	b.Run("", func(b *testing.B) {
		r.ServeHTTP(nil, req)
	})
}
