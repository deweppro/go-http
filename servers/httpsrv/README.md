# Using blank of http server

```go
package main

import (
	"net/http"
	"time"

	"github.com/deweppro/go-http/v2/servers/httpsrv"
	"github.com/deweppro/go-logger"
)

type Simple struct{}

func (s *Simple) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello world"))
}

func main() {
	simple := &Simple{}
	srv := httpsrv.NewCustomServer(httpsrv.ConfigItem{Addr: "localhost:8090"}, simple, logger.Default())
	if err := srv.Up(); err != nil {
		panic(err)
	}

	<-time.After(time.Minute)

	if err := srv.Down(); err != nil {
		panic(err)
	}
}

```