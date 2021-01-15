# Using debug server with custom setting

```go
package main

import (
	"time"

	"github.com/deweppro/go-http/v2/servers/debugsrv"
	"github.com/deweppro/go-http/v2/servers/httpsrv"
	"github.com/deweppro/go-logger"
)

func main() {
	dbg := debugsrv.NewCustom(httpsrv.ConfigItem{Addr: "localhost:8090"}, logger.Default())

	if err := dbg.Up(); err != nil {
		panic(err)
	}

	<-time.After(time.Minute)

	if err := dbg.Down(); err != nil {
		panic(err)
	}
}
```