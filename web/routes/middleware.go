package routes

import (
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/deweppro/go-logger"
)

//MiddlFunc interface of middleware
type MiddlFunc func(CtrlFunc) CtrlFunc

//CORSConfig model
type CORSConfig struct {
	Age     int
	Origin  string
	Methods []string
	Headers []string
}

//CORSMiddleware setting Cross-Origin Resource Sharing (CORS)
func CORSMiddleware(conf CORSConfig) func(c CtrlFunc) CtrlFunc {
	h := make(http.Header)
	if conf.Age > 0 {
		h.Add("Access-Control-Max-Age", strconv.FormatInt(int64(conf.Age), 10))
	}
	if len(conf.Origin) > 0 {
		h.Add("Access-Control-Allow-Origin", conf.Origin)
	}
	if len(conf.Methods) > 0 {
		h.Add("Access-Control-Allow-Methods", strings.Join(conf.Methods, ", "))
	}
	if len(conf.Headers) > 0 {
		h.Add("Access-Control-Allow-Headers", strings.Join(conf.Headers, ", "))
	}
	return func(c CtrlFunc) CtrlFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				for key := range h {
					w.Header().Set(key, h.Get(key))
				}
				return
			}
			c(w, r)
		}
	}
}

//Throttling limits active requests
func ThrottlingMiddleware(max int64) func(c CtrlFunc) CtrlFunc {
	var i int64
	return func(c CtrlFunc) CtrlFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if atomic.LoadInt64(&i) >= max {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			atomic.AddInt64(&i, 1)
			c(w, r)
			atomic.AddInt64(&i, -1)
		}
	}
}

//RecoveryMiddleware recovery go panic and write to log
func RecoveryMiddleware(log logger.Logger) func(c CtrlFunc) CtrlFunc {
	return func(c CtrlFunc) CtrlFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("Recovered: %+v", err)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			c(w, r)
		}
	}
}
