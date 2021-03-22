package common

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

const (
	VersionHeader      = `Accept`
	versionValueRegexp = `application\/vnd.v(\d+)\+json`
	versionValueTmpl   = `application/vnd.v%d+json`
)

var vercomp = regexp.MustCompile(versionValueRegexp)

//GetVersion getting version from header
func GetVersion(h http.Header) uint64 {
	d := h.Get(VersionHeader)
	r := signcomp.FindSubmatch([]byte(d))
	if len(r) == 2 {
		if v, err := strconv.ParseUint(string(r[1]), 10, 64); err == nil {
			return v
		}
	}
	return 0
}

//SetVersion setting version to header
func SetVersion(h http.Header, v uint64) {
	h.Set(VersionHeader, fmt.Sprintf(versionValueTmpl, v))
}
