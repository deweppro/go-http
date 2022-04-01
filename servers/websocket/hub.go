package websocket

import (
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/deweppro/go-http/pkg/errs"
	"github.com/deweppro/go-http/servers"
	"github.com/deweppro/go-logger"
)

var _ http.Handler = (*Hub)(nil)

type Hub struct {
	status  int64
	clients map[string]*connection
	call    Handler
	log     logger.Logger
	mux     sync.RWMutex
	wg      sync.WaitGroup
}

func New(handler Handler, log logger.Logger) *Hub {
	return &Hub{
		status:  servers.StatusOff,
		clients: make(map[string]*connection),
		call:    handler,
		log:     log,
	}
}

//Add connection to hub
func (v *Hub) Add(c *connection) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.clients[c.UUID()] = c

	v.wg.Add(1)
	c.OnClose(func(uuid string) {
		v.mux.Lock()
		defer v.mux.Unlock()

		delete(v.clients, uuid)
		v.wg.Done()
	})
}

//SendAll broadcast send message
func (v *Hub) SendAll(b []byte) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	for _, c := range v.clients {
		c.Send(b)
	}
}

//CloseAll close all connection
func (v *Hub) CloseAll() {
	v.mux.Lock()
	defer v.mux.Unlock()

	for _, c := range v.clients {
		c.Close()
	}
}

//ServeHTTP http handler
func (v *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt64(&v.status) != servers.StatusOn {
		w.WriteHeader(http.StatusGone)
		return
	}
	conn := newConnection(v.log)
	if err := conn.Upgrade(w, r, v.call); err == nil {
		v.Add(conn)
	} else {
		v.log.Errorf("upgrade for `%s`: %s", conn.UUID(), err.Error())
	}
}

//Up hub
func (v *Hub) Up() error {
	if !atomic.CompareAndSwapInt64(&v.status, servers.StatusOff, servers.StatusOn) {
		return errs.ErrServAlreadyRunning
	}
	return nil
}

//Down hub
func (v *Hub) Down() error {
	if !atomic.CompareAndSwapInt64(&v.status, servers.StatusOn, servers.StatusOff) {
		return errs.ErrServAlreadyStopped
	}
	v.CloseAll()
	v.wg.Wait()
	return nil
}
