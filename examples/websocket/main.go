package main

import (
	"time"

	"github.com/deweppro/go-http/web/server"
	"github.com/deweppro/go-http/websocket"
	"github.com/deweppro/go-logger"
)

func main() {
	hub := websocket.NewHub(handler)
	serv := server.NewCustom(server.ConfigItem{Addr: ":8080"}, hub, logger.Default())

	if err := hub.Up(); err != nil {
		panic(err)
	}
	if err := serv.Up(); err != nil {
		panic(err)
	}

	<-time.After(60 * time.Minute)

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
