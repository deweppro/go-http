package version

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

var rex = regexp.MustCompile(versionValueRegexp)

//Decode getting version from header
func Decode(h http.Header) uint64 {
	d := h.Get(VersionHeader)
	r := rex.FindSubmatch([]byte(d))
	if len(r) == 2 {
		if v, err := strconv.ParseUint(string(r[1]), 10, 64); err == nil {
			return v
		}
	}
	return 0
}

//Encode setting version to header
func Encode(h http.Header, v uint64) {
	h.Set(VersionHeader, fmt.Sprintf(versionValueTmpl, v))
}
