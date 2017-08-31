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
*	tmpr.PUT("", UpdateTemplate)
 */

func TestTplUpdate(t *testing.T) {
	routes := SetUpGin()
	Convey("update a template", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("update template ok", func() {
			postb := map[string]interface{}{
				"name":      "mytpl0",
				"parent_id": 0,
				"tpl_id":    2,
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("PUT", "/api/v1/template", &b)
			r = SetDefaultAdminSession(r)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "template updated")
			So(w.Code, ShouldEqual, 200)
		})
		Convey("update template faild", func() {
			postb := map[string]interface{}{
				"name":      "mytpl0",
				"parent_id": 0,
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("PUT", "/api/v1/template", &b)
			routes.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, 400)
		})
		Convey("update template faild with no permission", func() {
			postb := map[string]interface{}{
				"name":      "mytpl0",
				"parent_id": 0,
				"tpl_id":    2,
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("PUT", "/api/v1/template", &b)
			routes.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, 400)
		})
	})
}
