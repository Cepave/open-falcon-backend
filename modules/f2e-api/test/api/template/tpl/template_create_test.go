package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	. "github.com/smartystreets/goconvey/convey"
)

/*  convered routes test
*	tmpr.POST("", CreateTemplate)
 */

func TestTplCreate(t *testing.T) {
	routes := SetUpGin()
	Convey("create a template", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("create template ok", func() {
			postb := map[string]interface{}{
				"name":      "mytpl4",
				"parent_id": 0,
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/template", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "template created")
			So(w.Code, ShouldEqual, 200)
			postb = map[string]interface{}{
				"name":      "mytpl5",
				"parent_id": 0,
			}
			b, _ = json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/template", &b)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			So(respBody, ShouldContainSubstring, "template created")
			So(w.Code, ShouldEqual, 200)
			postb = map[string]interface{}{
				"name":      "mytpl6",
				"parent_id": 0,
			}
			b, _ = json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/template", &b)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			So(respBody, ShouldContainSubstring, "template created")
			So(w.Code, ShouldEqual, 200)
		})
	})
}
