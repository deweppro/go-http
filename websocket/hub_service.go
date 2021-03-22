package websocket

import (
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/deweppro/go-http/web/server"
)

var _ http.Handler = (*Hub)(nil)

//Hub connections store
type Hub struct {
	status  int64
	clients map[string]*Conn
	handler Handler
	lock    sync.RWMutex
	wg      sync.WaitGroup
}

//NewHub init hub
func NewHub(handler Handler) *Hub {
	return &Hub{
		status:  server.StatusOff,
		clients: make(map[string]*Conn),
		handler: handler,
		lock:    sync.RWMutex{},
	}
}

//Add connection to hub
func (h *Hub) Add(c *Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.clients[c.UUID()] = c

	h.wg.Add(1)
	c.OnClose(func(uuid string) {
		h.lock.Lock()
		defer h.lock.Unlock()

		delete(h.clients, uuid)
		h.wg.Done()
	})
}

//SendAll broadcast send message
func (h *Hub) SendAll(v []byte) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	for _, c := range h.clients {
		c.Send(v)
	}
}

//CloseAll cloase all connection
func (h *Hub) CloseAll() {
	h.lock.Lock()
	defer h.lock.Unlock()

	for _, c := range h.clients {
		c.Close()
	}
}

//ServeHTTP http handler
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt64(&h.status) != server.StatusOn {
		w.WriteHeader(http.StatusGone)
		return
	}
	conn := NewConn()
	if err := conn.Upgrade(w, r, h.handler); err == nil {
		h.Add(conn)
	}
}

//Up hub
func (h *Hub) Up() error {
	if !atomic.CompareAndSwapInt64(&h.status, server.StatusOff, server.StatusOn) {
		return server.ErrServAlreadyRunning
	}
	return nil
}

//Down hub
func (h *Hub) Down() error {
	if !atomic.CompareAndSwapInt64(&h.status, server.StatusOn, server.StatusOff) {
		return server.ErrServAlreadyStopped
	}
	h.CloseAll()
	h.wg.Wait()
	return nil
}
