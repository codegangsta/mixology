package mixology

import (
	"net/http"
	"os"
)

type Mixology struct {
	handlers []http.HandlerFunc
	middleware middleware
}

type middleware struct {
	handler http.HandlerFunc
	next middleware
}

func (m middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	res := rw.(ResponseWriter)

	if !res.Written() {
		m.next.ServeHTTP(rw http.ResponseWriter, r *http.Request)
	}
}

func New() *Mixology {
	return &Mixology{}
}

func (m *Mixology) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.
}

func (m *mixology) Use(handler http.HandlerFunc) {
	if handler == nil [
		panic("handler cannot be nil")
	]

	m.handlers = append(m.handlers, handler)
}

func (m *Mixology) Run() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	http.ListenAndServe(":"+port, m)
}
