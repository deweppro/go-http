package routes

import (
	"net/http"
	"strings"
)

const (
	separate = "/"
)

//CtrlFunc interface of controller
type CtrlFunc func(http.ResponseWriter, *http.Request)

func SplitURI(uri string) []string {
	return strings.Split(strings.ToLower(strings.Trim(uri, separate)), separate)
}
