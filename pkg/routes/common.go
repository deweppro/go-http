package routes

import (
	"strings"
)

const separate = "/"

func split(uri string) []string {
	vv := strings.Split(strings.ToLower(uri), separate)
	result := make([]string, 0, len(vv))
	for _, v := range vv {
		if len(v) > 0 {
			result = append(result, v)
		}
	}
	return result
}
