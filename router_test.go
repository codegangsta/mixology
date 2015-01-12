package mix_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codegangsta/mix"
)

func TestBasicRouting(t *testing.T) {
	m := mix.New()
	m.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(404)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	m.ServeHTTP(res, req)

	equals(t, 404, res.Code)
}
