package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/deweppro/go-http/utils"
	"github.com/deweppro/go-http/web/common"
	"github.com/pkg/errors"
)

type (
	//Client ...
	Client struct {
		cli *http.Client

		headers http.Header
		signer  utils.Signer

		debug  bool
		writer io.Writer
	}
)

//NewClient ...
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

//NewCustomClient ...
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
func (v *Client) WithAuth(s utils.Signer) {
	v.signer = s
}

//Call make request to server
func (v *Client) Call(pool utils.Pooler, method, uri string, body []byte) (int, []byte, error) {
	url, err := pool.Pool()
	if err != nil {
		return 0, nil, errors.Wrap(err, "get address from pool")
	}
	url.Path = uri
	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(body))
	if err != nil {
		return 0, nil, errors.Wrap(err, "create request")
	}
	req.Header.Set("Connection", "keep-alive")
	for k := range v.headers {
		req.Header.Set(k, v.headers.Get(k))
	}
	if v.signer != nil {
		common.SetSignature(req.Header, v.signer, body)
	}
	resp, err := v.cli.Do(req)
	if err != nil {
		return 0, nil, errors.Wrap(err, "make request")
	}
	defer resp.Body.Close() //nolint: errcheck
	b, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, b, err
}
