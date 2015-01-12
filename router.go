package mix

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

type Router struct {
	routes []*Route
}

func New() *Router {
	return &Router{}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		ok, params := route.Match(req.Method, req.URL.Path)
		if ok {
			// Add Params
			if len(params) > 0 {
				req.URL.RawQuery = url.Values(params).Encode() + "&" + req.URL.RawQuery
			}
			route.ServeHTTP(rw, req)
			return
		}
	}
}

func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.addRoute("GET", path, handler)
}

func (r *Router) addRoute(method, pattern string, handler http.HandlerFunc) *Route {
	route := &Route{
		method:  method,
		pattern: pattern,
		handler: handler,
	}

	// Grab the /:name/ params a sub with a regex
	regex := regexp.MustCompile(`:[^/#?()\.\\]+`)
	pattern = regex.ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})
	// r2 := regexp.MustCompile(`\*\*`)
	// var index int
	// pattern = r2.ReplaceAllStringFunc(pattern, func(m string) string {
	// 	index++
	// 	return fmt.Sprintf(`(?P<_%d>[^#?]*)`, index)
	// })
	pattern += `\/?`
	route.regex = regexp.MustCompile(pattern)

	r.routes = append(r.routes, route)
	return route
}

type Route struct {
	method  string
	regex   *regexp.Regexp
	pattern string
	handler http.HandlerFunc
}

func (r *Route) Match(method, path string) (bool, url.Values) {
	matches := r.regex.FindStringSubmatch(path)
	if len(matches) > 0 && matches[0] == path {
		params := url.Values{}
		for i, name := range r.regex.SubexpNames() {
			if len(name) > 0 {
				params[name] = append(params[name], matches[i])
			}
		}
		return true, params
	}
	return false, nil
}

func (r *Route) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(rw, req)
}
