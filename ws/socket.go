/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package ws

import (
	"context"
	"net/http"

	pool "github.com/deweppro/go-chan-pool"
	"github.com/deweppro/go-logger"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type (
	//Message ...
	Message struct {
		Type int
		Text []byte
	}
	//Handler ...
	Handler func(out chan<- *Message, in <-chan *Message, ctx context.Context, cncl context.CancelFunc)
	//WebSocket ...
	WebSocket struct {
		buf   int
		count int
		pool  pool.ChanPool
		serv  websocket.Upgrader
		log   logger.Logger
	}
)

//New ...
func NewServer(countMsg, bufferSize int, log logger.Logger) *WebSocket {
	return &WebSocket{
		log:   log,
		count: countMsg,
		buf:   bufferSize,
		pool: pool.ChanPool{
			Size: countMsg,
			New: func() interface{} {
				return &Message{}
			},
		},
		serv: websocket.Upgrader{
			ReadBufferSize:  bufferSize,
			WriteBufferSize: bufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *WebSocket) Handler(w http.ResponseWriter, r *http.Request, handler Handler) error {
	conn, err := s.serv.Upgrade(w, r, nil)
	if err != nil {
		return errors.Wrap(err, "failed to set websocket upgrade")
	}

	ctx, cncl := context.WithCancel(context.Background())
	defer cncl()

	in := make(chan *Message, s.count)
	out := make(chan *Message, s.count)
	go handler(out, in, ctx, cncl)

	go func() {
		for {
			select {
			case <-ctx.Done():
				err := conn.WriteMessage(websocket.CloseMessage, []byte(`Bye bye!`))
				if err != nil && errors.Is(err, websocket.ErrCloseSent) {
					s.log.Errorf("websocket write message: %s", err.Error())
				}
				return
			case d := <-out:
				if err := conn.WriteMessage(d.Type, d.Text); err != nil {
					s.log.Errorf("websocket write message: %s", err.Error())
				}
				d.Type = 0
				d.Text = d.Text[:0]
				s.pool.Put(d)
			}
		}
	}()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			return errors.Wrap(err, "websocket read message")
		}
		rm := s.pool.Get().(*Message)
		rm.Type = msgType
		rm.Text = msg
		in <- rm
	}
}
