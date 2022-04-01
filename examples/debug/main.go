package main

import (
	"time"

	"github.com/deweppro/go-http/servers"
	"github.com/deweppro/go-http/servers/debug"
	"github.com/deweppro/go-logger"
)

func main() {
	logger.Default().SetLevel(logger.LevelDebug)

	conf := servers.Config{Addr: ":8080"}
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
