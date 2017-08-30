package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tidwall/gjson"
)

/*  convered routes test
*	tmpr.POST("", CreateTemplate)
 */

func TestTplGet(t *testing.T) {
	routes := SetUpGin()
	Convey("test get info of template", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("get template list", func() {
			w, r = NewTestContextWithDefaultSession("GET", "/api/v1/template", nil)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			LCheck := gjson.Get(respBody, "templates.#.template")
			So(len(LCheck.Array()), ShouldBeGreaterThan, 2)
			So(w.Code, ShouldEqual, 200)
			CleanSession(r)
			Convey("test template get id", func() {
				Tid := gjson.Get(respBody, "templates.1.template.id")
				uri := "/api/v1/template/" + Tid.String()
				w, r = NewTestContextWithDefaultSession("GET", uri, nil)
				routes.ServeHTTP(w, r)
				respBody := w.Body.String()
				RCheck := gjson.Get(respBody, "template.id")
				So(RCheck.String(), ShouldEqual, Tid.String())
				So(w.Code, ShouldEqual, 200)
			})
		})
	})
}
