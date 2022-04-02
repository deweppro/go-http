package errs

import "github.com/deweppro/go-errors"

var (
	ErrResolveTCPAddress  = errors.New("resolve tcp address")
	ErrInvalidSignature   = errors.New("invalid signature header")
	ErrEmptyPool          = errors.New("empty pool")
	ErrInvalidPoolAddress = errors.New("invalid address")
	ErrServAlreadyRunning = errors.New("server already running")
	ErrServAlreadyStopped = errors.New("server already stopped")
	ErrInvalidPoolType    = errors.New("invalid data type from pool")
	ErrEpollEmptyEvents   = errors.New("epoll empty event")
	ErrFailContextKey     = errors.New("context key is not found")
)
