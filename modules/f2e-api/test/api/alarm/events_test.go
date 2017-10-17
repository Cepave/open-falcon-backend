package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"

	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tidwall/gjson"
)

func TestEvents(t *testing.T) {
	routes := SetUpGin()
	var (
		w        *httptest.ResponseRecorder
		r        *http.Request
		postb    map[string]interface{}
		b        []byte
		respBody string
		rCheck   gjson.Result
	)
	Convey("test get events with id", t, func() {
		// string type
		postb = map[string]interface{}{
			"startTime": 1466600460,
			"endTime":   1467291660,
			"event_id":  "s_50_6438ac68b30e2712fb8f00d894c46e21",
			"page":      1,
			"limit":     10,
		}

		Convey("test limit params", func() {
			postb2 := postb
			b, _ = json.Marshal(postb2)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/alarm/events", &b)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			log.Debug(respBody)
			rCheck = gjson.Parse(respBody)
			So(len(rCheck.Array()), ShouldEqual, 2)
			postb2["limit"] = 1
			b, _ = json.Marshal(postb2)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/alarm/events", &b)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			log.Debug(respBody)
			rCheck = gjson.Parse(respBody)
			So(len(rCheck.Array()), ShouldEqual, 1)
		})
		Convey("test startTime & endTime", func() {
			postb["endTime"] = 1467637260
			b, _ = json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/alarm/events", &b)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			log.Debug(respBody)
			rCheck = gjson.Parse(respBody)
			So(len(rCheck.Array()), ShouldEqual, 10)
		})
		Convey("test status filter", func() {
			postb["endTime"] = 1467637260
			postb["status"] = 1
			b, _ = json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/alarm/events", &b)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			log.Debug(respBody)
			rCheck = gjson.Parse(respBody)
			So(len(rCheck.Array()), ShouldEqual, 5)

			postb["status"] = 0
			b, _ = json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/alarm/events", &b)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			log.Debug(respBody)
			rCheck = gjson.Parse(respBody)
			So(len(rCheck.Array()), ShouldEqual, 5)

			postb["status"] = -1
			b, _ = json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/alarm/events", &b)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			log.Debug(respBody)
			rCheck = gjson.Parse(respBody)
			So(len(rCheck.Array()), ShouldEqual, 10)
		})
	})
}
