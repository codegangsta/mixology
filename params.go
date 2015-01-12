package mix

import (
	"net/http"

	"github.com/gorilla/context"
)

type Params map[string]string

type paramsHelperKey int

const paramsKey paramsHelperKey = 0

func GetParams(r *http.Request) Params {
	if rv := context.Get(r, paramsKey); rv != nil {
		return rv.(Params)
	}
	return nil
}

func setParams(r *http.Request, val Params) {
	context.Set(r, paramsKey, val)
}
