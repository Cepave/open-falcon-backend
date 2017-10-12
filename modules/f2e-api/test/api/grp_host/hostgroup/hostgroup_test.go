package test

import (
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/tidwall/gjson"

	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

/*  convered routes test
// hostr.GET("/hostgroup", GetHostGroups)
// hostr.POST("/hostgroup", CrateHostGroup)
// hostr.POST("/hostgroup/host", BindHostToHostGroup)
// hostr.PUT("/hostgroup/host", UnBindAHostToHostGroup)
// hostr.GET("/hostgroup/:host_group", GetHostGroup)
// hostr.DELETE("/hostgroup/:host_group", DeleteHostGroup)
*/

func TestHostGroupCreate(t *testing.T) {
	routes := SetUpGin()
	Convey("create hostgroup", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		postb := map[string]interface{}{
			"name": "hostgroup1",
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/hostgroup", &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
		postb = map[string]interface{}{
			"name": "ahostgroup2",
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/hostgroup", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
		postb = map[string]interface{}{
			"name": "testhostgroup3",
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/hostgroup", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		So(w.Code, ShouldEqual, 200)
		postb = map[string]interface{}{
			"name": "hosttestgroup4",
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/hostgroup", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
	})
}

func TestGetHostGroup(t *testing.T) {
	routes := SetUpGin()
	Convey("get list of hostgroup", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		uriBase := "/api/v1/hostgroup"
		w, r = NewTestContextWithDefaultSession("GET", uriBase, nil)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		log.Debug(respBody)
		rCheck := gjson.Get(respBody, "#.id")
		So(len(rCheck.Array()), ShouldEqual, 4)
		So(w.Code, ShouldEqual, 200)

		/* test with pagging */
		mixUri := fmt.Sprintf("%s?limit=%d&page=%d", uriBase, 1, 1)
		w, r = NewTestContextWithDefaultSession("GET", mixUri, nil)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		rCheck = gjson.Get(respBody, "#.id")
		So(len(rCheck.Array()), ShouldEqual, 1)

		/* test with regex query */
		// params decode
		u, _ := url.Parse(".*test.*")
		mixUri = fmt.Sprintf("%s?q=%s", uriBase, u.String())
		w, r = NewTestContextWithDefaultSession("GET", mixUri, nil)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		rCheck = gjson.Get(respBody, "#.id")
		So(len(rCheck.Array()), ShouldEqual, 2)
	})
}

func TestBindHostGroupCreate1(t *testing.T) {
	routes := SetUpGin()
	Convey("bind host to hostgroup 1", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		postb := map[string]interface{}{
			"hostgroup_id": 1,
			"hosts":        []string{"a", "b", "c", "d"},
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/hostgroup/host", &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
	})
}

func TestBindHostGroupCreate2(t *testing.T) {
	routes := SetUpGin()
	Convey("bind host to hostgroup 2", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		postb := map[string]interface{}{
			"hostgroup_id": 1,
			"hosts":        []string{"a", "b", "c", "e"},
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/hostgroup/host", &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
	})
}

func TestUnBindHostGroupCreate1(t *testing.T) {
	routes := SetUpGin()
	Convey("unbind host to hostgroup", t, func() {
		// perpare hostgroup
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		postb := map[string]interface{}{
			"hostgroup_id": 2,
			"hosts":        []string{"a", "b", "c", "d"},
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/hostgroup/host", &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
		rCheck := gjson.Get(respBody, "hosts.#")
		var ids []int64
		for _, rt := range rCheck.Array() {
			ids = append(ids, rt.Int())
		}
		ubindids := ids[0]
		postb = map[string]interface{}{
			"hostgroup_id": 1,
			"host_id":      ubindids,
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("PUT", "/api/v1/hostgroup/host", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
	})
}

func TestDeleteHostGroup(t *testing.T) {
	routes := SetUpGin()
	Convey("delete a hostgroup", t, func() {
		// perpare hostgroup
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		// perpare hostgroup for delete
		postb := map[string]interface{}{
			"name": "hostgroupdelme",
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/hostgroup", &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		log.Debug(respBody)
		hid := gjson.Get(respBody, "id")
		So(w.Code, ShouldEqual, 200)
		Convey("delete a exist hostgroup", func() {
			urlBase := fmt.Sprintf("/api/v1/hostgroup/%d", hid.Int())
			w, r = NewTestContextWithDefaultSession("DELETE", urlBase, nil)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			log.Debug(respBody)
			So(w.Code, ShouldEqual, 200)
		})
	})
}
