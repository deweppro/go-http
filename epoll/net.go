/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package epoll

import (
	"net"
	"sync"
)

type netConnItem struct {
	Conn  net.Conn
	await bool
	Fd    int
	sync.RWMutex
}

func (n *netConnItem) Await(b bool) {
	n.Lock()
	n.await = b
	n.Unlock()
}

func (n *netConnItem) IsAwait() bool {
	n.Lock()
	defer n.Unlock()
	return n.await
}
