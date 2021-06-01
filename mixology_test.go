package mixology

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicMiddleware(t *testing.T) {
	result := []string{}
	response := httptest.NewRecorder()

	// create a new mixology
	m := New()

	m.Use(func(rw http.ResponseWriter, r *http.Request) {
		result = append(result, "One")
	})

	m.Use(func(rw http.ResponseWriter, r *http.Request) {
		result = append(result, "Two")
	})

	m.Use(func(rw http.ResponseWriter, r *http.Request) {
		result = append(result, "Three")
	})

	m.ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, []string{"One", "Two", "Three"})
}
