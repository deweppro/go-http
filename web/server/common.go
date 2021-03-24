package server

import "github.com/pkg/errors"

var (
	//ErrServAlreadyRunning error than run method UP if server already started
	ErrServAlreadyRunning = errors.New("server already running")
	//ErrServAlreadyStopped error than run method Down if server already stopped
	ErrServAlreadyStopped = errors.New("server already stopped")
)
