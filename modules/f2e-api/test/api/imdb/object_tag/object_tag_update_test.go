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
	"github.com/tidwall/gjson"
)

func TestObjectTagUpdate(t *testing.T) {
	routes := SetUpGin()
	var (
		w        *httptest.ResponseRecorder
		r        *http.Request
		postb    map[string]interface{}
		b        []byte
		respBody string
		surl     string
		rCheck   gjson.Result
	)
	// resourceObjId := 20
	Convey("update object new tag", t, func() {
		// string type
		surl = fmt.Sprintf("/api/v1/imdb/object_tag/%v", 7)
		postb = map[string]interface{}{
			"value_text": "127.0.0.99",
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("PUT", surl, &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
		rCheck = gjson.Get(respBody, "data.value")
		So(rCheck.String(), ShouldEqual, postb["value_text"])

		// int type
		surl = fmt.Sprintf("/api/v1/imdb/object_tag/%v", 13)
		postb = map[string]interface{}{
			"value_int": 30,
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("PUT", surl, &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
		rCheck = gjson.Get(respBody, "data.value")
		So(rCheck.Int(), ShouldEqual, postb["value_int"])

		// description type
		surl = fmt.Sprintf("/api/v1/imdb/object_tag/%v", 12)
		postb = map[string]interface{}{
			"value_text": "ooo 000-000-000",
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("PUT", surl, &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
		rCheck = gjson.Get(respBody, "data.value")
		So(rCheck.String(), ShouldEqual, postb["value_text"])

		// value model type
		surl = fmt.Sprintf("/api/v1/imdb/object_tag/%v", 14)
		// c01.i02
		postb = map[string]interface{}{
			"value_int": 2,
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("PUT", surl, &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
		rCheck = gjson.Get(respBody, "data.value")
		So(rCheck.String(), ShouldEqual, "c01.i02")
	})
}
