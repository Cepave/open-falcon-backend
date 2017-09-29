package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"testing"

	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestObjectTagDelete(t *testing.T) {
	routes := SetUpGin()
	var (
		w        *httptest.ResponseRecorder
		r        *http.Request
		respBody string
		postb    map[string]interface{}
		b        []byte
		surl     string
	)
	Convey("delete a object tag", t, func() {
		postb = map[string]interface{}{
			"test": false,
		}
		b, _ = json.Marshal(postb)
		surl = fmt.Sprintf("/api/v1/imdb/object_tag/%v", 17)
		w, r = NewTestContextWithDefaultSession("DELETE", surl, &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
	})
}
