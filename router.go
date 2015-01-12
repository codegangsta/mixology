package mix

import (
	"net/http"
	"strings"
)

type Router struct {
	routes []*Route
	groups []string
}

func New() *Router {
	return &Router{}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := r.tokenize(req.URL.Path)
	for _, route := range r.routes {
		ok, params := route.Match(req.Method, path)
		if ok {
			setParams(req, params)
			route.ServeHTTP(rw, req)
			return
		}
	}
}

func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.addRoute("GET", path, handler)
}

func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.addRoute("POST", path, handler)
}

func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.addRoute("PUT", path, handler)
}

func (r *Router) Patch(path string, handler http.HandlerFunc) {
	r.addRoute("PATCH", path, handler)
}

func (r *Router) Option(path string, handler http.HandlerFunc) {
	r.addRoute("OPTION", path, handler)
}

func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.addRoute("DELETE", path, handler)
}

func (r *Router) Group(pattern string, fn func(r *Router)) {
	r.groups = append(r.groups, r.sanitize(pattern))
	fn(r)
	r.groups = r.groups[:len(r.groups)-1]
}

func (r *Router) addRoute(method, pattern string, handler http.HandlerFunc) *Route {
	if len(r.groups) > 0 {
		pattern = strings.Join(r.groups, "/") + pattern
	}

	route := &Route{
		method:  method,
		pattern: pattern,
		handler: handler,
	}

	route.tokens = r.tokenize(pattern)

	r.routes = append(r.routes, route)
	return route
}

func (r *Router) tokenize(path string) []string {
	return strings.Split(r.sanitize(path), "/")[1:]
}

// Manually trimming strings for performance reasons
func (r *Router) sanitize(path string) string {
	last := len(path) - 1
	if path[last] == '/' {
		path = path[:last]
		last--
	}
	if last >= 0 && path[0] == '/' {
		path = path[1:]
	}
	return path
}
