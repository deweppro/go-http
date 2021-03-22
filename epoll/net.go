package epoll

import (
	"net"
	"sync"
)

type NetItem struct {
	Conn  net.Conn
	await bool
	Fd    int
	sync.RWMutex
}

//Await ...
func (n *NetItem) Await(b bool) {
	n.Lock()
	n.await = b
	n.Unlock()
}

//IsAwait ...
func (n *NetItem) IsAwait() bool {
	n.Lock()
	defer n.Unlock()
	return n.await
}
