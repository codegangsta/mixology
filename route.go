package mix

import (
	"net/http"
	"regexp"
	"strings"
)

type Route struct {
	method  string
	regex   *regexp.Regexp
	pattern string
	handler http.HandlerFunc
	tokens  []string
}

func (r *Route) Match(method, path string) (bool, Params) {
	if !r.MatchMethod(method) {
		return false, nil
	}

	matches := r.regex.FindStringSubmatch(path)
	if len(matches) > 0 && matches[0] == path {
		params := Params{}
		for i, name := range r.regex.SubexpNames() {
			if len(name) > 0 {
				params[name] = matches[i]
			}
		}
		return true, params
	}
	return false, nil
}

func (r *Route) Match2(method, path string) (bool, Params) {
	if !r.MatchMethod(method) {
		return false, nil
	}

	tokens := strings.Split(path, "/")[1:]
	if len(tokens) != len(r.tokens) {
		return false, nil
	}

	var params map[string]string

	// loop over each token and find a match
	for i, t := range r.tokens {
		// it's a variable
		if len(t) > 0 && t[:1] == ":" {
			if params == nil {
				params = make(map[string]string)
			}
			// do params matching
			params[t[1:]] = tokens[i]
		} else if t != tokens[i] {
			return false, nil
		}
	}

	return true, params
}

func (r *Route) MatchMethod(method string) bool {
	return r.method == "*" || method == r.method || (method == "HEAD" && r.method == "GET")
}

func (r *Route) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(rw, req)
}
