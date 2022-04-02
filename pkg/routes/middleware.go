package routes

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/deweppro/go-logger"
)

//MiddlFunc interface of middleware
type MiddlFunc func(CtrlFunc) CtrlFunc

//CORSConfig model
type CORSConfig struct {
	Age     int      `yaml:"age"`
	Origin  []string `yaml:"origin"`
	Methods []string `yaml:"methods"`
	Headers []string `yaml:"headers"`
}

//CORSMiddleware setting Cross-Origin Resource Sharing (CORS)
func CORSMiddleware(conf CORSConfig) func(c CtrlFunc) CtrlFunc {
	h := make(http.Header)
	if conf.Age > 0 {
		h.Add("Access-Control-Max-Age", strconv.FormatInt(int64(conf.Age), 10))
	}
	if len(conf.Methods) > 0 {
		h.Add("Access-Control-Allow-Methods", strings.Join(conf.Methods, ", "))
	}
	if len(conf.Headers) > 0 {
		h.Add("Access-Control-Allow-Headers", strings.Join(conf.Headers, ", "))
	}
	h.Add("Vary", "Origin")
	return func(c CtrlFunc) CtrlFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			uri, err := url.Parse(r.Referer())
			if err != nil {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			for key := range h {
				w.Header().Set(key, h.Get(key))
			}
			if len(conf.Origin) > 0 {
				for _, v := range conf.Origin {
					if v == uri.Host {
						w.Header().Set("Access-Control-Allow-Origin", uri.Scheme+"://"+uri.Host)
					}
				}
			} else {
				w.Header().Set("Access-Control-Allow-Origin", uri.Scheme+"://"+uri.Host)
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			c(w, r)
		}
	}
}

//ThrottlingMiddleware limits active requests
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
