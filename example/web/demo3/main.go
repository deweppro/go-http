/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

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
		"demo0": []string{"http://localhost:8090"},
		"demo1": []string{"http://localhost:8091"},
	}}

	cli := webcli.NewClient()
	cli.Debug(true, os.Stdout)
	cli.WithSign(sign)

	protoClient := proto.NewClient(pool)
	protoClient.Client("demo0", cli.Call)
	protoClient.Client("demo1", cli.Call)

	/**
	init server
	*/

	protoServer := proto.NewServer()
	protoServer.Handler("/", 1, (&Demo{}).Index)

	conf := func(port int) *websrv.Config {
		return &websrv.Config{
			HTTP: httpsrv.ConfigItem{Addr: fmt.Sprintf("localhost:%d", port)},
			Headers: websrv.Headers{
				ProxyHeaders:   []string{"X-Forwarded-For", "Accept-Language", "User-Agent"},
				DefaultHeaders: map[string]string{"Content-Type": "text/html"},
			},
		}
	}

	srv := websrv.NewServer(conf(8090), protoServer, logger.Default())
	if err := srv.Up(); err != nil {
		panic(err)
	}
	defer func() {
		if err := srv.Down(); err != nil {
			panic(err)
		}
	}()

	srv2 := websrv.NewServer(conf(8091), protoServer, logger.Default())
	if err := srv2.Up(); err != nil {
		panic(err)
	}
	defer func() {
		if err := srv2.Down(); err != nil {
			panic(err)
		}
	}()

	<-time.After(time.Second * 2)

	/**
	make request
	*/

	req := proto.NewRequest()
	res := proto.NewResponse()

	/**
	case 1
	*/

	if err := req.EncodeJSON(A{Name: "A"}); err != nil {
		fmt.Println("client encode", err)
	}
	if err := protoClient.Call("demo0", req, res); err != nil {
		fmt.Println(err)
	}
	var result A
	if err := res.DecodeJSON(&result); err != nil {
		fmt.Println("client decode", err)
	}
	fmt.Println("UUID is", req.GetUUID() == res.GetUUID())
	fmt.Println("BODY is", result.Name == "A+B")

	/**
	case 2
	*/

	if err := req.EncodeJSON(A{Name: "C"}); err != nil {
		fmt.Println("client encode", err)
	}
	if err := protoClient.Call("demo1", req, res); err != nil {
		fmt.Println(err)
	}
	var result1 A
	if err := res.DecodeJSON(&result1); err != nil {
		fmt.Println("client decode", err)
	}
	fmt.Println("UUID is", req.GetUUID() == res.GetUUID())
	fmt.Println("BODY is", result1.Name == "C+B")

	/**
	case 3
	*/

	if err := protoClient.Call("demo2", req, res); err != nil {
		fmt.Println(err)
	}

	/**
	case 4
	*/

	req.SetVersion(10)
	if err := protoClient.Call("demo1", req, res); err != nil {
		fmt.Println(err)
	}
}
