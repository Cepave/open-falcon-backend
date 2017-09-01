package test

import (
	"encoding/json"
	"fmt"
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

func TestTplClone(t *testing.T) {
	routes := SetUpGin()
	Convey("clone a template", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		postb := map[string]interface{}{
			"name":      "mytplclone1",
			"parent_id": 0,
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/template", &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		So(respBody, ShouldContainSubstring, "template created")
		So(w.Code, ShouldEqual, 200)
		tplId := gjson.Get(respBody, "template.id")
		postb = map[string]interface{}{
			"uic":                  "teamA,teamB",
			"url":                  "localhost:9999/v1/fix",
			"callback":             1,
			"before_callback_sms":  0,
			"after_callback_sms":   0,
			"before_callback_mail": 0,
			"after_callback_mail":  0,
			"tpl_id":               tplId.Int(),
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/template/action", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		So(respBody, ShouldContainSubstring, "action is created and bind to")
		CleanSession(r)
		Convey("clone template ok", func() {
			postb := map[string]interface{}{
				"name": "mytplclone1_copy",
				"id":   tplId.Int(),
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/template/clone_tpl", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "cloned")
			//{"message":"template is cloned","tpl_id":95}
			tplClonedId := gjson.Get(respBody, "tpl_id")
			So(w.Code, ShouldEqual, 200)
			w, r = NewTestContextWithDefaultSession("GET", fmt.Sprintf("/api/v1/template/%d", tplClonedId.Int()), nil)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			actId := gjson.Get(respBody, "template.action_id")
			So(actId.Int(), ShouldNotEqual, 0)
			// delete templates
			for _, v := range []int64{tplId.Int(), tplClonedId.Int()} {
				uri := fmt.Sprintf("/api/v1/template/%d", v)
				w, r = NewTestContextWithDefaultSession("DELETE", uri, nil)
				routes.ServeHTTP(w, r)
				respBody := w.Body.String()
				So(respBody, ShouldContainSubstring, "has been delete")
				So(w.Code, ShouldEqual, 200)
			}
		})
	})
}
