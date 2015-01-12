package mix

import (
	"fmt"
	"net/http"
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
		ok, _ := route.Match(req.Method, req.URL.Path)
		if ok {
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

func (r *Route) Match(method, path string) (bool, map[string]string) {
	matches := r.regex.FindStringSubmatch(path)
	if len(matches) > 0 && matches[0] == path {
		params := make(map[string]string)
		for i, name := range r.regex.SubexpNames() {
			if len(name) > 0 {
				params[name] = matches[i]
			}
		}
		return true, params
	}
	return false, nil
}

func (r *Route) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(rw, req)
}
