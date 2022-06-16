package routes

import (
	"strings"
)

const separate = "/"

func split(uri string) []string {
	vv := strings.Split(strings.ToLower(uri), separate)
	for i := 0; i < len(vv); i++ {
		if len(vv[i]) == 0 {
			copy(vv[i:], vv[i+1:])
			vv = vv[:len(vv)-1]
			i--
		}
	}
	return vv
}
