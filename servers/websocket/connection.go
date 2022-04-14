package websocket

import (
	"context"
	"net/http"
	"time"

	"github.com/deweppro/go-logger"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

var _ Connector = (*connection)(nil)

const (
	pongWait   = 60 * time.Second
	pingPeriod = pongWait / 3
)

var upgrader = ws.Upgrader{
	EnableCompression: true,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type connection struct {
	uuid        string
	onCloseFunc []func(string)

	conn    *ws.Conn
	headers http.Header
	handler Handler

	sendChan  chan []byte
	ctx       context.Context
	closeFunc context.CancelFunc
	log       logger.Logger
}

func newConnection(log logger.Logger) *connection {
	ctx, cncl := context.WithCancel(context.TODO())
	return &connection{
		uuid:        uuid.NewString(),
		onCloseFunc: make([]func(string), 0),
		sendChan:    make(chan []byte, 128),
		closeFunc:   cncl,
		ctx:         ctx,
		log:         log,
	}
}

//UUID get conn unique ID
func (v *connection) UUID() string {
	return v.uuid
}

//Headers request headers
func (v *connection) Headers() http.Header {
	return v.headers
}

//Close connection
func (v *connection) Close() {
	v.closeFunc()
	for _, fn := range v.onCloseFunc {
		fn(v.UUID())
	}
}

//OnClose event on close connection
func (v *connection) OnClose(cb func(string)) {
	v.onCloseFunc = append(v.onCloseFunc, cb)
}

//Send message
func (v *connection) Send(b []byte) {
	select {
	case v.sendChan <- b:
	default:
	}
}

func (v *connection) pumpWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-v.ctx.Done():
			msg := ws.FormatCloseMessage(ws.CloseNormalClosure, "Bye")
			if err := v.conn.WriteMessage(ws.CloseMessage, msg); err != nil {
				v.log.WithFields(logger.Fields{
					"err": err.Error(), "uuid": v.UUID(),
				}).Errorf("close conn")
			}
			return
		case msg := <-v.sendChan:
			if err := v.conn.WriteMessage(ws.TextMessage, msg); err != nil {
				v.log.WithFields(logger.Fields{
					"err": err.Error(), "uuid": v.UUID(),
				}).Errorf("send message")
				return
			}
		case <-ticker.C:
			if err := v.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				v.log.WithFields(logger.Fields{
					"err": err.Error(), "uuid": v.UUID(),
				}).Errorf("send ping")
				return
			}
		}
	}
}

func (v *connection) pumpRead() {
	defer v.Close()
	for {
		_, message, err := v.conn.ReadMessage()
		if err != nil {
			v.log.WithFields(logger.Fields{
				"err": err.Error(), "uuid": v.UUID(),
			}).Errorf("read message")
			return
		}
		v.handler(message, v)
	}
}

//Upgrade tcp connection
func (v *connection) Upgrade(w http.ResponseWriter, r *http.Request, handler Handler) (err error) {
	v.conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		v.Close()
		return
	}

	v.handler = handler
	v.headers = r.Header
	if err = v.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		v.Close()
		return
	}
	v.conn.SetPongHandler(func(string) error {
		return v.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	go v.pumpWrite()
	go v.pumpRead()
	return
}
