package pool

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/deweppro/go-http/pkg/errs"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var (
	_ ConfigGetter = (*Pool)(nil)
	_ PoolGetter   = (*List)(nil)
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type (
	ConfigGetter interface {
		Get(name string) List
	}

	Pool struct {
		Items map[string][]string `yaml:"services" json:"services"`
	}
)

func (c *Pool) Get(name string) List {
	v, ok := c.Items[name]
	if !ok {
		return List{}
	}
	return v
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type (
	PoolGetter interface {
		Pool() (*url.URL, error)
	}

	List []string
)

func (v List) Pool() (*url.URL, error) {
	l := len(v)
	if l == 0 {
		return nil, errs.ErrEmptyPool
	}
	u, err := url.Parse(v[rand.Intn(len(v))])
	if err != nil || len(u.Scheme) == 0 || len(u.Host) == 0 {
		return nil, errs.ErrInvalidPoolAddress
	}

	return u, nil
}
