package routes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/deweppro/go-http/internal"
)

var match = regexp.MustCompile(`\{([a-z0-9]+)\:?([^{}]*)\}`)

type matcher struct {
	incr    int
	keys    map[string]string
	links   map[string]string
	pattern string
	rex     *regexp.Regexp
}

func newMatcher() *matcher {
	return &matcher{
		incr:    1,
		pattern: "",
		keys:    make(map[string]string),
		links:   make(map[string]string),
	}
}

func (v *matcher) Add(vv string) error {
	result := vv

	patterns := match.FindAllString(vv, -1)
	for _, pattern := range patterns {
		res := match.FindAllStringSubmatch(pattern, 1)[0]

		key := fmt.Sprintf("k%d", v.incr)
		rex := ".+"
		if len(res) == 3 && len(res[2]) > 0 {
			rex = res[2]
		}
		result = strings.Replace(result, res[0], fmt.Sprintf("(?P<%s>%s)", key, rex), 1)

		v.links[key] = vv
		v.keys[key] = res[1]
		v.incr++
	}

	result = "^" + result + "$"

	if _, err := regexp.Compile(result); err != nil {
		return fmt.Errorf("regex compilation error for `%s`: %w", vv, err)
	}

	if len(v.pattern) != 0 {
		v.pattern += "|"
	}
	v.pattern += result
	v.rex = regexp.MustCompile(v.pattern)
	return nil
}

func (v *matcher) Match(vv string, vars internal.VarsData) (string, bool) {
	if v.rex == nil {
		return "", false
	}

	matches := v.rex.FindStringSubmatch(vv)
	if len(matches) == 0 {
		return "", false
	}

	link := ""
	for indx, name := range v.rex.SubexpNames() {
		val := matches[indx]
		if len(val) == 0 {
			continue
		}
		if l, ok := v.links[name]; ok {
			link = l
		}
		if key, ok := v.keys[name]; ok {
			vars[key] = val
		}
	}

	return link, true
}

func hasMatcher(vv string) bool {
	return match.MatchString(vv)
}
