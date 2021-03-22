package utils

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrEmptyPool          = errors.New("empty pool")
	ErrInvalidPoolAddress = errors.New("invalid address")
)

type (
	//Pool ...
	Pool struct {
		Items map[string]List `yaml:"services" json:"services"`
	}
	//List ...
	List []string
	//Pooler ...
	Pooler interface {
		Pool() (*url.URL, error)
	}
	//Configer ...
	Configer interface {
		Get(name string) List
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

//Get ...
func (c Pool) Get(name string) List {
	v, ok := c.Items[name]
	if !ok {
		return List{}
	}
	return v
}

//Pool ...
func (v List) Pool() (*url.URL, error) {
	l := len(v)
	if l == 0 {
		return nil, ErrEmptyPool
	}
	u, err := url.Parse(v[rand.Intn(len(v))])
	if err != nil || len(u.Scheme) == 0 || len(u.Host) == 0 {
		return nil, ErrInvalidPoolAddress
	}

	return u, nil
}
