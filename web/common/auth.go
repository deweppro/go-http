package common

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/deweppro/go-http/utils"
	"github.com/pkg/errors"
)

const (
	SignHeader      = `Signature`
	signValueRegexp = `keyId=\"(.*)\",algorithm=\"(.*)\",signature=\"(.*)\"`
	signValueTmpl   = `keyId="%s",algorithm="%s",signature="%s"`
)

var (
	signcomp = regexp.MustCompile(signValueRegexp)
	//ErrInvalidSignature if invalid header format
	ErrInvalidSignature = errors.New("invalid signature header")
)

//Signature model
type Signature struct {
	ID   string
	Alg  string
	Hash string
}

//GetSignature getting signature from header
func GetSignature(h http.Header) (s Signature, err error) {
	d := h.Get(SignHeader)
	r := signcomp.FindSubmatch([]byte(d))
	if len(r) != 4 {
		err = ErrInvalidSignature
		return
	}
	s.ID, s.Alg, s.Hash = string(r[1]), string(r[2]), string(r[3])
	return
}

//SetSignature make and setting signature to header
func SetSignature(h http.Header, s utils.Signer, body []byte) {
	h.Set(SignHeader, fmt.Sprintf(signValueTmpl, s.ID(), s.Algorithm(), s.CreateString(body)))
}
