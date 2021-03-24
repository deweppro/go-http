package epoll

import (
	"github.com/pkg/errors"
)

//go:generate easyjson

var (
	ErrEpollEmptyEvents = errors.New("epoll empty event")
)

//easyjson:json
type (
	//Config model
	Config struct {
		Epoll ConfigItem `yaml:"epoll" json:"epoll"`
	}

	//ConfigItem model
	ConfigItem struct {
		Addr string `yaml:"addr" json:"addr"`
	}
)
