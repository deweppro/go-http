package epoll

import (
	"net"
	"sync"
)

type netItem struct {
	Conn  net.Conn
	await bool
	Fd    int
	mux   sync.RWMutex
}

func (v *netItem) Await(b bool) {
	v.mux.Lock()
	v.await = b
	v.mux.Unlock()
}

func (v *netItem) IsAwait() bool {
	v.mux.RLock()
	is := v.await
	v.mux.RUnlock()
	return is
}
