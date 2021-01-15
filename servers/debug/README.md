# init debug server with custom setting

```go
package main

import (
	"time"

	"github.com/deweppro/go-http/v2/servers/debug"
	"github.com/deweppro/go-http/v2/servers/http"
	"github.com/deweppro/go-logger"
)

func main() {
	dbg := debug.NewCustom(http.ConfigItem{Addr: "localhost:8090"}, logger.Default())
	if err := dbg.Up(); err != nil {
		panic(err)
	}

	<-time.After(time.Minute)

	if err := dbg.Down(); err != nil {
		panic(err)
	}
}

```