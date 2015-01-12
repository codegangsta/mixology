package mix

import (
	"net/http"

	"github.com/nbio/httpcontext"
)

type Params map[string]string

type paramsHelperKey int

const paramsKey paramsHelperKey = 0

func GetParams(r *http.Request) Params {
	if rv := httpcontext.Get(r, paramsKey); rv != nil {
		return rv.(Params)
	}
	return nil
}

func setParams(r *http.Request, val Params) {
	httpcontext.Set(r, paramsKey, val)
}
