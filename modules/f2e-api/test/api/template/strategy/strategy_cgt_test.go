package test

import (
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tidwall/gjson"

	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	. "github.com/smartystreets/goconvey/convey"
)

/*  convered routes test
*	strr.POST("", CreateStrategy)
* strr.GET("", GetStrategys)
 */

func TestStrategyCreate(t *testing.T) {
	routes := SetUpGin()
	var (
		w *httptest.ResponseRecorder
		r *http.Request
	)
	Convey("create a strategy", t, func() {
		// create a template for action test
		postb := map[string]interface{}{
			"metric":      "net.if.in.bits",
			"tags":        "iface=all",
			"max_step":    3,
			"priority":    1,
			"func":        "all(#3)",
			"op":          ">",
			"right_value": "10",
			"note":        "a test strategy0",
			"run_begin":   "",
			"run_end":     "",
			"tpl_id":      1,
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/strategy", &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		So(respBody, ShouldContainSubstring, "stragtegy created")
		postb = map[string]interface{}{
			"metric":      "cpu.busy",
			"tags":        "",
			"max_step":    3,
			"priority":    2,
			"func":        "all(#3)",
			"op":          ">",
			"right_value": "10",
			"note":        "a test strategy1",
			"run_begin":   "00:00",
			"run_end":     "24:00",
			"tpl_id":      1,
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/strategy", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		So(respBody, ShouldContainSubstring, "stragtegy created")
		postb = map[string]interface{}{
			"metric":      "mem.memfree.percent",
			"tags":        "",
			"max_step":    3,
			"priority":    2,
			"func":        "all(#5)",
			"op":          "<",
			"right_value": "20",
			"note":        "a test strategy2",
			"run_begin":   "",
			"run_end":     "",
			"tpl_id":      2,
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/strategy", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		So(respBody, ShouldContainSubstring, "stragtegy created")
		postb = map[string]interface{}{
			"metric":      "mem.memfree.percent",
			"tags":        "",
			"max_step":    3,
			"priority":    2,
			"func":        "all(#5)",
			"op":          "<",
			"right_value": "20",
			"note":        "a test strategy3",
			"run_begin":   "00:00",
			"run_end":     "05:00",
			"tpl_id":      3,
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/strategy", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		So(respBody, ShouldContainSubstring, "stragtegy created")
	})
}

func TestStrategyGet(t *testing.T) {
	routes := SetUpGin()
	var (
		w *httptest.ResponseRecorder
		r *http.Request
	)
	Convey("get strategy list", t, func() {
		tid := 1
		w, r = NewTestContextWithDefaultSession("GET", fmt.Sprintf("/api/v1/strategy?tid=%d", tid), nil)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		ids := gjson.Get(respBody, "#.id")
		So(len(ids.Array()), ShouldBeGreaterThan, 1)
	})
	Convey("get a strategy", t, func() {
		sid := 1
		w, r = NewTestContextWithDefaultSession("GET", fmt.Sprintf("/api/v1/strategy/%d", sid), nil)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		id := gjson.Get(respBody, "id")
		So(id.Int(), ShouldEqual, 1)
	})
}

func TestUpdateStrategy(t *testing.T) {
	routes := SetUpGin()
	var (
		w *httptest.ResponseRecorder
		r *http.Request
	)
	Convey("update a strategy", t, func() {
		postb := map[string]interface{}{
			"id":          3,
			"metric":      "net.if.in.bits",
			"tags":        "iface=eth0",
			"max_step":    3,
			"priority":    2,
			"func":        "all(#3)",
			"op":          "<",
			"right_value": "1000000",
			"note":        "a test strategy3",
			"run_begin":   "01:00",
			"run_end":     "02:00",
			"tpl_id":      2,
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("PUT", "/api/v1/strategy", &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		metirc := gjson.Get(respBody, "strategy.metric")
		So(metirc.String(), ShouldEqual, "net.if.in.bits")
	})
}
