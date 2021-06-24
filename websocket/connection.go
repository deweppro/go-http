package websocket

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

var _ Connector = (*Conn)(nil)

const (
	pongWait   = 60 * time.Second
	pingPeriod = pongWait / 3
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type (
	//Conn client connect
	Conn struct {
		uuid        string
		onCloseFunc []func(string)

		conn    *ws.Conn
		headers http.Header
		handler Handler

		sendChan  chan []byte
		ctx       context.Context
		closeFunc context.CancelFunc
	}
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

//NewConn init new connection
func NewConn() *Conn {
	ctx, cncl := context.WithCancel(context.TODO())
	return &Conn{
		uuid:        uuid.NewString(),
		onCloseFunc: make([]func(string), 0),
		sendChan:    make(chan []byte, 128),
		closeFunc:   cncl,
		ctx:         ctx,
	}
}

//UUID get conn unique ID
func (c *Conn) UUID() string {
	return c.uuid
}

//Headers request headers
func (c *Conn) Headers() http.Header {
	return c.headers
}

//Close connection
func (c *Conn) Close() {
	c.closeFunc()
	for _, fn := range c.onCloseFunc {
		fn(c.UUID())
	}
}

//OnClose event on close connection
func (c *Conn) OnClose(cb func(string)) {
	c.onCloseFunc = append(c.onCloseFunc, cb)
}

//Send message
func (c *Conn) Send(v []byte) {
	select {
	case c.sendChan <- v:
	default:
	}
}

func (c *Conn) pumpWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-c.ctx.Done():
			c.conn.WriteMessage(ws.CloseMessage, ws.FormatCloseMessage(ws.CloseNormalClosure, "Bye")) //nolint: errcheck
			return
		case msg := <-c.sendChan:
			if err := c.conn.WriteMessage(ws.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Conn) pumpRead() {
	defer c.Close()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		c.handler(message, c)
	}
}

//Upgrade tcp connection
func (c *Conn) Upgrade(w http.ResponseWriter, r *http.Request, handler Handler) (err error) {
	c.conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.Close()
		return
	}

	c.handler = handler
	c.headers = r.Header
	c.conn.SetReadDeadline(time.Now().Add(pongWait))                                                           //nolint: errcheck
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil }) //nolint: errcheck

	go c.pumpWrite()
	go c.pumpRead()
	return
}
