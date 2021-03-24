package main

import (
	"time"

	"github.com/deweppro/go-http/web/debug"
	"github.com/deweppro/go-http/web/server"
	"github.com/deweppro/go-logger"
)

func main() {
	conf := &debug.Config{Debug: server.ConfigItem{Addr: ":8080"}}
	serv := debug.New(conf, logger.Default())
	if err := serv.Up(); err != nil {
		panic(err)
	}
	<-time.After(60 * time.Second)
	if err := serv.Down(); err != nil {
		panic(err)
	}
	logger.Close()
}
