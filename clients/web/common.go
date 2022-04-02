package web

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/deweppro/go-errors"
	"github.com/deweppro/go-http/internal"
	"github.com/deweppro/go-http/pkg/errs"
	"github.com/deweppro/go-http/pkg/pool"
	"github.com/deweppro/go-http/pkg/signature"
)

//Client ...
type Client struct {
	cli *http.Client

	headers http.Header
	signer  signature.SignGetter

	debug  bool
	writer io.Writer
}

func New() *Client {
	cli := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
		Timeout: 5 * time.Second,
	}
	return NewCustom(cli)
}

func NewCustom(cli *http.Client) *Client {
	return &Client{
		cli:     cli,
		headers: make(http.Header),
	}
}

//Debug enable logging of responses
func (v *Client) Debug(is bool, w io.Writer) {
	v.debug, v.writer = is, w
}

//WithHeaders setting headers
func (v *Client) WithHeaders(heads http.Header) {
	v.headers = heads
}

//WithAuth sitting auth
func (v *Client) WithAuth(s signature.SignGetter) {
	v.signer = s
}

//Call make request to server
func (v *Client) Call(method, uri string, body []byte) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Connection", "keep-alive")
	for k := range v.headers {
		req.Header.Set(k, v.headers.Get(k))
	}
	if v.signer != nil {
		signature.Encode(req.Header, v.signer, body)
	}
	resp, err := v.cli.Do(req)
	if err != nil {
		return 0, nil, err
	}
	b, err := internal.ReadAll(resp.Body)
	return resp.StatusCode, b, err
}

//CallPool make request to pool of server
func (v *Client) CallPool(pool pool.PoolGetter, method, uri string, body []byte) (int, []byte, error) {
	url, err := pool.Pool()
	if err != nil {
		return 0, nil, errors.Wrap(err, errs.ErrEmptyPool)
	}
	url.Path = uri
	return v.Call(method, url.String(), body)
}
