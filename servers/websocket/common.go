package websocket

import "net/http"

type (
	//Connector interface
	Connector interface {
		UUID() string
		Headers() http.Header
		Send(v []byte)
		OnClose(cb func(string))
		Close()
	}
	//Handler request processor
	Handler func([]byte, Connector)
)
