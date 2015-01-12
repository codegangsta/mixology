package mix

import (
	"net/http"

	"github.com/nbio/httpcontext"
)

type paramsHelperKey int

const paramsKey paramsHelperKey = 0

func Params(r *http.Request) map[string]string {
	if rv := httpcontext.Get(r, paramsKey); rv != nil {
		return rv.(map[string]string)
	}
	return nil
}

func setParams(r *http.Request, val map[string]string) {
	httpcontext.Set(r, paramsKey, val)
}
