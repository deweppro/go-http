package utils

import (
	"net"
	"reflect"
	"strings"
)

//RandomPort getting random port
func RandomPort(host string) (string, error) {
	host = strings.Join([]string{host, "0"}, ":")
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return host, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return host, err
	}
	return l.Addr().String(), l.Close()
}

//FD getting file descriptor
func FD(c net.Conn) int {
	fd := reflect.Indirect(reflect.ValueOf(c)).FieldByName("fd")
	pfd := reflect.Indirect(fd).FieldByName("pfd")
	return int(pfd.FieldByName("Sysfd").Int())
}
