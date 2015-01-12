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

	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	m.ServeHTTP(res, req)

	equals(t, "Hello World", res.Body.String())
}

func TestBasicParams(t *testing.T) {
	m := mix.New()
	m.Get("/pages/:pageId/events/:id", func(rw http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		fmt.Fprint(rw, params.Get("pageId"), params.Get("id"))
	})

	req, _ := http.NewRequest("GET", "/pages/123/events/456/", nil)
	res := httptest.NewRecorder()

	m.ServeHTTP(res, req)

	equals(t, "123456", res.Body.String())
}
