package mix_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codegangsta/mix"
	"github.com/codegangsta/negroni"
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

func TestGroup(t *testing.T) {
	m := mix.New()
	m.Group("/admin", func(r *mix.Router) {
		r.Get("/pages", func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprint(rw, "pages route")
		})
		r.Get("/events/:id", func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprint(rw, mix.Params(r)["id"])
		})
	})
	res := req(m, "GET", "/admin/pages")
	equals(t, "pages route", res.Body.String())
	res = req(m, "GET", "/admin/events/123")
	equals(t, "123", res.Body.String())
}

func TestNestedGroup(t *testing.T) {
	m := mix.New()
	m.Group("/one", func(r *mix.Router) {
		m.Group("/two", func(r *mix.Router) {
			r.Get("/pages", func(rw http.ResponseWriter, r *http.Request) {
				fmt.Fprint(rw, "pages route")
			})
			r.Get("/events", func(rw http.ResponseWriter, r *http.Request) {
				fmt.Fprint(rw, "events route")
			})
		})
	})
	res := req(m, "GET", "/one/two/pages")
	equals(t, "pages route", res.Body.String())
	res = req(m, "GET", "/one/two/events")
	equals(t, "events route", res.Body.String())
}

func TestGroupMiddleware(t *testing.T) {
	m := mix.New()

	auth := func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fmt.Fprint(rw, "middleware ")
		next(rw, r)
	}

	m.Group("/admin", func(r *mix.Router) {
		r.Get("/foo", func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprint(rw, "foo route")
		})
	}, negroni.HandlerFunc(auth))

	res := req(m, "GET", "/admin/foo")
	equals(t, "middleware foo route", res.Body.String())
}

func TestDefaultNotFound(t *testing.T) {
	m := mix.New()
	res := req(m, "GET", "/not-here")
	equals(t, 404, res.Code)
	equals(t, "404 page not found\n", res.Body.String())
}

func TestCustomNotFound(t *testing.T) {
	m := mix.New()
	m.NotFound = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Custom Not Found!", 405)
	}
	res := req(m, "GET", "/not-here")
	equals(t, 405, res.Code)
	equals(t, "Custom Not Found!\n", res.Body.String())
}

func req(handler http.Handler, method, path string) *httptest.ResponseRecorder {
	r, _ := http.NewRequest(method, path, nil)
	rw := httptest.NewRecorder()

	handler.ServeHTTP(rw, r)
	return rw
}
