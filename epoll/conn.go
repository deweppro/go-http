package epoll

import (
	"bytes"
	"io"

	cpool "github.com/deweppro/go-chan-pool"
	"github.com/pkg/errors"
)

var (
	ErrInvalidPoolType = errors.New("invalid data type from pool")

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
	b, ok := pool.Get().([]byte)
	if !ok {
		return ErrInvalidPoolType
	}
	defer pool.Put(b[:0])

	for {
		if len(b) == cap(b) {
			b = append(b, 0)[:len(b)]
		}
		n, err = conn.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
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
