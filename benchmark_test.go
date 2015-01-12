package mix_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codegangsta/mix"
)

func Benchmark1Route(b *testing.B) {
	m := mix.New()
	m.Get("/pages", func(http.ResponseWriter, *http.Request) {
	})
	req, _ := http.NewRequest("GET", "/pages", nil)
	res := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		m.ServeHTTP(res, req)
	}
}

func Benchmark1RouteWithParams(b *testing.B) {
	m := mix.New()
	m.Get("/pages/:pageId/events/:id", func(http.ResponseWriter, *http.Request) {
	})

	req, _ := http.NewRequest("GET", "/pages/123/events/456", nil)
	res := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		m.ServeHTTP(res, req)
	}
}
