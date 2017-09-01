package test

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tidwall/gjson"

	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	. "github.com/smartystreets/goconvey/convey"
)

/*  convered routes test
*	tmpr.PUT("", UpdateTemplate)
 */

func TestTplDelete(t *testing.T) {
	routes := SetUpGin()
	Convey("delete a template", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("SetUpGin template ok", func() {
			postb := map[string]interface{}{
				"name":      "mytpl99",
				"parent_id": 0,
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/template", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "template created")
			Tid := gjson.Get(respBody, "template.id")
			So(w.Code, ShouldEqual, 200)
			Convey("delete template ok", func() {
				uri := "/api/v1/template/" + Tid.String()
				w, r = NewTestContextWithDefaultSession("DELETE", uri, nil)
				routes.ServeHTTP(w, r)
				respBody = w.Body.String()
				So(respBody, ShouldContainSubstring, "has been delete")
				So(w.Code, ShouldEqual, 200)
			})
		})
	})
}

// func TestTplDeleteHelper(t *testing.T) {
// 	routes := SetUpGin()
// 	Convey("helper for delete templates", t, func() {
// 		var (
// 			w *httptest.ResponseRecorder
// 			r *http.Request
// 		)
// 		for _, v := range []int{99, 100} {
// 			uri := fmt.Sprintf("/api/v1/template/%d", v)
// 			w, r = NewTestContextWithDefaultSession("DELETE", uri, nil)
// 			routes.ServeHTTP(w, r)
// 			respBody := w.Body.String()
// 			So(respBody, ShouldContainSubstring, "has been delete")
// 			So(w.Code, ShouldEqual, 200)
// 		}
// 	})
// }
