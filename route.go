package mix

import "net/http"

type Route struct {
	method  string
	pattern string
	name    string
	handler http.HandlerFunc
	tokens  []string
}

func (r *Route) Match(method string, tokens []string) (bool, map[string]string) {
	if !r.MatchMethod(method) {
		return false, nil
	}

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

func (r *Route) Name(name string) {
	r.name = name
}
