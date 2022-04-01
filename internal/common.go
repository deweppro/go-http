package internal

import (
	"io"
	"io/ioutil"
	"net"
	"reflect"
	"strings"

	"github.com/deweppro/go-errors"
	"github.com/deweppro/go-http/pkg/errs"
)

func RandomPort(host string) (string, error) {
	host = strings.Join([]string{host, "0"}, ":")
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return host, errors.Wrap(err, errs.ErrResolveTCPAddress)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return host, errors.Wrap(err, errs.ErrResolveTCPAddress)
	}
	v := l.Addr().String()
	if err = l.Close(); err != nil {
		return host, errors.Wrap(err, errs.ErrResolveTCPAddress)
	}
	return v, nil
}

func FD(c net.Conn) int {
	fd := reflect.Indirect(reflect.ValueOf(c)).FieldByName("fd")
	pfd := reflect.Indirect(fd).FieldByName("pfd")
	return int(pfd.FieldByName("Sysfd").Int())
}

func ValidateAddress(addr string) string {
	hp := strings.Split(addr, ":")
	if len(hp[0]) == 0 {
		hp[0] = "0.0.0.0"
		return strings.Join(hp, ":")
	}
	if len(hp) < 2 || (len(hp) == 2 && len(hp[1]) == 0) {
		v, err := RandomPort(hp[0])
		if err != nil {
			v = hp[0] + ":80"
		}
		return v
	}
	return addr
}

func ReadAll(r io.ReadCloser) ([]byte, error) {
	defer r.Close() //nolint: errcheck
	return ioutil.ReadAll(r)
}
