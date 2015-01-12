package mix

import (
	"net/http"
	"strings"
)

type Router struct {
	routes []*Route
}

func New() *Router {
	return &Router{}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		ok, params := route.Match2(req.Method, req.URL.Path)
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

func (r *Router) addRoute(method, pattern string, handler http.HandlerFunc) *Route {
	route := &Route{
		method:  method,
		pattern: pattern,
		handler: handler,
	}

	// Grab the /:name/ params a sub with a regex
	// regex := regexp.MustCompile(`:[^/#?()\.\\]+`)
	// pattern = regex.ReplaceAllStringFunc(pattern, func(m string) string {
	// 	return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	// })
	// r2 := regexp.MustCompile(`\*\*`)
	// var index int
	// pattern = r2.ReplaceAllStringFunc(pattern, func(m string) string {
	// 	index++
	// 	return fmt.Sprintf(`(?P<_%d>[^#?]*)`, index)
	// })
	// pattern += `\/?`
	// route.regex = regexp.MustCompile(pattern)

	route.tokens = strings.Split(pattern, "/")[1:]

	r.routes = append(r.routes, route)
	return route
}
