package epoll

import (
	"bytes"
	"io"

	cpool "github.com/deweppro/go-chan-pool"
)

var (
	pool = cpool.ChanPool{
		Size: eventCount,
		New: func() interface{} {
			return make([]byte, 0, 1024)
		},
	}
)

type (
	//ConnectionHandler ...
	ConnectionHandler func([]byte, io.Writer) error
)

func connection(conn io.ReadWriter, handler ConnectionHandler, eof []byte) error {
	var (
		n   int
		err error
		l   = len(eof)
	)
	b := pool.Get().([]byte)
	defer pool.Put(b[:0])

	for {
		if len(b) == cap(b) {
			b = append(b, 0)[:len(b)]
		}
		n, err = conn.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return err
		}
		if len(b) < l {
			return io.EOF
		}
		if bytes.Equal(eof, b[len(b)-l:]) {
			b = b[:len(b)-l]
			break
		}
	}
	return handler(b, conn)
}
