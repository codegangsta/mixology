package mix_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codegangsta/mix"
)

func TestBasicRouting(t *testing.T) {
	m := mix.New()
	m.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello World"))
	})

	equals(t, "Hello World", req(m, "GET", "/").Body.String())
}

func TestSlashPathsRouting(t *testing.T) {
	m := mix.New()
	m.Get("/hello/world", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello World"))
	})

	equals(t, "Hello World", req(m, "GET", "/hello/world/").Body.String())
}

func TestBasicParams(t *testing.T) {
	m := mix.New()
	m.Get("/pages/:pageId/events/:id", func(rw http.ResponseWriter, r *http.Request) {
		params := mix.Params(r)
		fmt.Fprint(rw, params["pageId"], params["id"])
	})

	res := req(m, "GET", "/pages/123/events/456")
	equals(t, "123456", res.Body.String())
}

func TestMethods(t *testing.T) {
	m := mix.New()

	fn := func(name string) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprint(rw, "method: ", name)
		}
	}

	m.Get("/", fn("GET"))
	m.Post("/", fn("POST"))
	m.Put("/", fn("PUT"))
	m.Patch("/", fn("PATCH"))
	m.Option("/", fn("OPTION"))
	m.Delete("/", fn("DELETE"))

	for _, method := range []string{"GET", "POST", "PUT", "PATCH", "OPTION", "DELETE"} {
		equals(t, "method: "+method, req(m, method, "/").Body.String())
	}
}

func TestHead(t *testing.T) {
	m := mix.New()
	m.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, "GET")
	})

	equals(t, "GET", req(m, "HEAD", "/").Body.String())
}

func req(handler http.Handler, method, path string) *httptest.ResponseRecorder {
	r, _ := http.NewRequest(method, path, nil)
	rw := httptest.NewRecorder()

	handler.ServeHTTP(rw, r)
	return rw
}
