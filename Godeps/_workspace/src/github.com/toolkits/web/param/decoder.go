package param

import (
	"encoding/json"
	"net/http"

	"github.com/open-falcon/fe/Godeps/_workspace/src/github.com/toolkits/web/errors"
)

func ParseJson(r *http.Request, v interface{}) {
	if r.ContentLength == 0 {
		panic(errors.BadRequestError("content is blank"))
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	if err != nil {
		panic(errors.BadRequestError(err.Error()))
	}
}
