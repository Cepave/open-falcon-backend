package test

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

/*  convered routes test
*	tmpr.POST("", CreateTemplate)
*	tmpr.POST("/action", CreateActionToTmplate)
*	tmpr.PUT("/action", UpdateActionToTmplate)
*	tmpr.DELETE("/:tpl_id", DeleteTemplate)
 */

func TestActionCreate(t *testing.T) {
	routes := SetUpGin()
	Convey("create a action", t, func() {
		log.Debug("create a action")
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		// create a template for action test
		postb := map[string]interface{}{
			"name":      "testtp3",
			"parent_id": 0,
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/template", &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		So(respBody, ShouldContainSubstring, "template created")
		So(w.Code, ShouldEqual, 200)
		tplId := gjson.Get(respBody, "template.id")
		CleanSession(r)
		Convey("create action ok", func() {
			postb := map[string]interface{}{
				"uic":                  "teamA,teamB",
				"url":                  "localhost:9999/v1/fix",
				"callback":             1,
				"before_callback_sms":  0,
				"after_callback_sms":   0,
				"before_callback_mail": 0,
				"after_callback_mail":  0,
				"tpl_id":               tplId.Int(),
			}
			b, _ := json.Marshal(postb)
			w, r := NewTestContextWithDefaultSession("POST", "/api/v1/template/action", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "action is created and bind to")
			So(w.Code, ShouldEqual, 200)
			aId := gjson.Get(respBody, "action.id")
			// test create twice failed
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/template/action", &b)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			So(respBody, ShouldContainSubstring, "already exist")
			So(w.Code, ShouldEqual, 400)
			CleanSession(r)
			Convey("update action ok", func() {
				postb := map[string]interface{}{
					"id":                   aId.Int(),
					"uic":                  "teamA",
					"url":                  "localhost:9999/v1/fix",
					"callback":             0,
					"before_callback_sms":  1,
					"after_callback_sms":   1,
					"before_callback_mail": 1,
					"after_callback_mail":  1,
				}
				b, _ := json.Marshal(postb)
				w, r := NewTestContextWithDefaultSession("PUT", "/api/v1/template/action", &b)
				routes.ServeHTTP(w, r)
				respBody := w.Body.String()
				So(respBody, ShouldContainSubstring, "action is updated")
				Lbody := gjson.Get(respBody, "action.UIC")
				So(Lbody.String(), ShouldEqual, postb["uic"])
				So(w.Code, ShouldEqual, 200)

				// delete this template
				w, r = NewTestContextWithDefaultSession("DELETE", fmt.Sprintf("/api/v1/template/%d", tplId.Int()), nil)
				routes.ServeHTTP(w, r)
				respBody = w.Body.String()
			})
		})
	})
}
