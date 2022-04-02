package main

import (
	"io"
	"time"

	"github.com/deweppro/go-http/servers/epoll"

	"github.com/deweppro/go-http/servers"
	"github.com/deweppro/go-logger"
)

func main() {
	logger.Default().SetLevel(logger.LevelDebug)

	conf := servers.Config{Addr: ":8080"}
	serv := epoll.New(conf, handler, nil, logger.Default())

	if err := serv.Up(); err != nil {
		panic(err)
	}

	<-time.After(60 * time.Second)

	if err := serv.Down(); err != nil {
		panic(err)
	}

	logger.Close()
}

func handler(r []byte, w io.Writer) error {
	var result []byte
	result = append(result, []byte("> ")...)
	result = append(result, r...)
	result = append(result, []byte("\n")...)
	_, err := w.Write(result)
	return err
}
