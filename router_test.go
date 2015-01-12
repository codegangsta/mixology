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

func TestBasicParams(t *testing.T) {
	m := mix.New()
	m.Get("/pages/:pageId/events/:id", func(rw http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		fmt.Fprint(rw, params.Get("pageId"), params.Get("id"))
	})

	res := req(m, "GET", "/pages/123/events/456/")
	equals(t, "123456", res.Body.String())
}

func req(handler http.Handler, method, path string) *httptest.ResponseRecorder {
	r, _ := http.NewRequest(method, path, nil)
	rw := httptest.NewRecorder()

	handler.ServeHTTP(rw, r)
	return rw
}
