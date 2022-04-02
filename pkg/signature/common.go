package signature

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/deweppro/go-http/pkg/errs"
)

const (
	SignHeader      = `Signature`
	signValueRegexp = `keyId=\"(.*)\",algorithm=\"(.*)\",signature=\"(.*)\"`
	signValueTmpl   = `keyId="%s",algorithm="%s",signature="%s"`
)

var rex = regexp.MustCompile(signValueRegexp)

type Data struct {
	ID   string
	Alg  string
	Hash string
}

//Decode getting signature from header
func Decode(h http.Header) (s Data, err error) {
	d := h.Get(SignHeader)
	r := rex.FindSubmatch([]byte(d))
	if len(r) != 4 {
		err = errs.ErrInvalidSignature
		return
	}
	s.ID, s.Alg, s.Hash = string(r[1]), string(r[2]), string(r[3])
	return
}

//Encode make and setting signature to header
func Encode(h http.Header, s SignGetter, body []byte) {
	h.Set(SignHeader, fmt.Sprintf(signValueTmpl, s.ID(), s.Algorithm(), s.CreateString(body)))
}
