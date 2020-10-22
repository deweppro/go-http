/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package epoll

import (
	"io"

	cpool "github.com/deweppro/go-chan-pool"
)

var (
	endByteUnix = byte('\n')
	endByteWin  = byte('\r')

	pool = cpool.ChanPool{
		Size: cEventCount,
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}
)

type (
	ConnectionHandler func([]byte, io.Writer) error
)

func connection(conn io.ReadWriter, handler ConnectionHandler) (err error) {
	var (
		n   int
		end int
	)

	data := pool.Get().([]byte)
	defer pool.Put(data)

	for l := 0; l < 100; l++ {
		n, err = conn.Read(data[end:cap(data)])
		switch err {
		case nil:
		case io.EOF:
			err = nil
			break
		default:
			return
		}
		if n <= 0 {
			continue
		}
		end += n
		if len(data[0:end]) == cap(data) {
			data = append(data[0:end], make([]byte, cap(data)*2)...)
		} else {
			break
		}
	}

	if end > 0 && data[end-1] == endByteUnix {
		end--
	}
	if end > 0 && data[end-1] == endByteWin {
		end--
	}
	if end > 0 {
		err = handler(data[0:end], conn)
	}

	return
}
