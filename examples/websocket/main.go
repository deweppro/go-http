package main

import (
	"time"

	"github.com/deweppro/go-http/servers"
	"github.com/deweppro/go-http/servers/web"
	"github.com/deweppro/go-http/servers/websocket"
	"github.com/deweppro/go-logger"
)

func main() {
	logger.Default().SetLevel(logger.LevelDebug)

	hub := websocket.New(handler, logger.Default())
	serv := web.New(servers.Config{Addr: ":8080"}, hub, logger.Default())

	if err := hub.Up(); err != nil {
		panic(err)
	}
	if err := serv.Up(); err != nil {
		panic(err)
	}

	<-time.After(60 * time.Second)

	if err := serv.Down(); err != nil {
		panic(err)
	}
	if err := hub.Down(); err != nil {
		panic(err)
	}

	logger.Close()
}

func handler(r []byte, w websocket.Connector) {
	result := []byte("=> ")
	result = append(result, r...)
	w.Send(result)
}
