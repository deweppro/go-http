# web server with proto

```go
package main

import (
	"bytes"
	"strconv"
	"time"

	"github.com/deweppro/go-http/v2"
	"github.com/deweppro/go-http/v2/servers/http"
	"github.com/deweppro/go-http/v2/servers/web"
	"github.com/deweppro/go-logger"
)

type Simple struct{}

func (s *Simple) Index(in proto.Request, out *proto.Response) {
	out.Code = proto.StatusCodeOK
	buf := bytes.Buffer{}
	buf.WriteString("<html><body><pre>")
	buf.WriteString("UUID: " + in.UUID + "\n")
	buf.WriteString("Path: " + in.Path + "\n")
	buf.WriteString("Version: " + strconv.FormatUint(uint64(in.Version), 10) + "\n")
	buf.WriteString("Meta: " + "\n")
	for k := range in.Meta {
		buf.WriteString(" - " + k + ": " + in.Meta.Get(k) + "\n")
	}
	buf.WriteString("</pre></body></html>")
	out.Body = buf.Bytes()
}

func main() {
	prt := proto.NewProto()
	prt.Handler("/", 1, (&Simple{}).Index)

	conf := &web.Config{
		HTTP: http.ConfigItem{Addr: "localhost:8090"},
		Headers: web.Headers{
			ProxyHeaders:   []string{"X-Forwarded-For", "Accept-Language", "User-Agent"},
			DefaultHeaders: map[string]string{"Content-Type": "text/html"},
		},
	}

	srv := web.NewServer(conf, prt, logger.Default())
	if err := srv.Up(); err != nil {
		panic(err)
	}

	<-time.After(time.Minute)

	if err := srv.Down(); err != nil {
		panic(err)
	}
}
```