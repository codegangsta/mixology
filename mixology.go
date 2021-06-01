package mixology

import (
	"net/http"
	"os"
)

type Mixology struct {
	handlers   []http.HandlerFunc
	middleware http.Handler
}

func New() *Mixology {
	return &Mixology{
		handlers:   []http.HandlerFunc{},
		middleware: voidMiddleware(),
	}
}

func (m *Mixology) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.middleware.ServeHTTP(NewResponseWriter(rw), r)
}

func (m *Mixology) Use(handler http.HandlerFunc) {
	if handler == nil {
		panic("handler cannot be nil")
	}

	m.handlers = append(m.handlers, handler)
	m.middleware = build(m.handlers)
}

func (m *Mixology) Run() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	http.ListenAndServe(":"+port, m)
}

type middleware struct {
	handler http.HandlerFunc
	next    http.Handler
}

func (m middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	res := rw.(ResponseWriter)

	if !res.Written() && m.next == nil {
		m.next.ServeHTTP(rw, r)
	}
}

func build(handlers []http.HandlerFunc) middleware {
	var next middleware

	switch {
	case len(handlers) == 0:
		return voidMiddleware()
	case len(handlers) > 1:
		next = build(handlers[1:])
	default:
		next = voidMiddleware()
	}

	return middleware{handlers[0], &next}
}

func voidMiddleware() middleware {
	return middleware{
		func(rw http.ResponseWriter, r *http.Request) {},
		&middleware{},
	}
}
