# Using the server and client with the proto

```go
package main

import (
	"fmt"
	"os"
	"time"

	proto "github.com/deweppro/go-http/v2"
	"github.com/deweppro/go-http/v2/clients/webcli"
	"github.com/deweppro/go-http/v2/servers/httpsrv"
	"github.com/deweppro/go-http/v2/servers/websrv"
	"github.com/deweppro/go-logger"
)

type A struct {
	Name string `json:"name"`
}

type B struct {
	Name string `json:"name"`
}

type Demo struct{}

func (s *Demo) Index(in *proto.Request, out *proto.Response) {
	var model B
	if err := in.DecodeJSON(&model); err != nil {
		fmt.Println("Index Decode", err)
	}

	out.SetStatusCode(proto.StatusCodeOK)
	if err := out.EncodeJSON(B{Name: model.Name + "+B"}); err != nil {
		fmt.Println("Index Encode", err)
	}
}

func main() {

	/**
	init client
	*/

	sign := proto.NewSigner("1", "1234567890")
	pool := proto.Pool{Items: map[string]proto.List{
		"demo": []string{"http://localhost:8090"},
	}}

	cli := webcli.NewClient()
	cli.Debug(true, os.Stdout)

	protoClient := proto.NewClient(pool)
	protoClient.Client("demo", cli.Call)

	/**
	init server
	*/

	protoServer := proto.NewServer()
	protoServer.Handler("/", 1, (&Demo{}).Index)

	conf := &websrv.Config{
		HTTP: httpsrv.ConfigItem{Addr: `localhost:8080`},
		Headers: websrv.Headers{
			ProxyHeaders:   []string{"X-Forwarded-For", "Accept-Language", "User-Agent"},
			DefaultHeaders: map[string]string{"Content-Type": "text/html"},
		},
	}

	srv := websrv.NewServer(conf, protoServer, logger.Default())
	if err := srv.Up(); err != nil {
		panic(err)
	}
	defer func() {
		if err := srv.Down(); err != nil {
			panic(err)
		}
	}()

	<-time.After(time.Second * 2)

	/**
	make request
	*/

	req := proto.NewRequest()
	res := proto.NewResponse()

	if err := req.EncodeJSON(A{Name: "A"}); err != nil {
		fmt.Println("client encode", err)
	}
	req.CreateSign(sign)
	if err := protoClient.Call("demo", req, res); err != nil {
		fmt.Println(err)
	}

	var result A
	if err := res.DecodeJSON(&result); err != nil {
		fmt.Println("client decode", err)
	}

	fmt.Println("UUID is", req.GetUUID() == res.GetUUID())
	fmt.Println("BODY is", result.Name == "A+B")
}
```