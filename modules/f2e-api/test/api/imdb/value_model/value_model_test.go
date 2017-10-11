package test

import (
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tidwall/gjson"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/imdb"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValueModelGet(t *testing.T) {
	routes := SetUpGin()
	var (
		w        *httptest.ResponseRecorder
		r        *http.Request
		respBody string
		rCheck   gjson.Result
	)
	Convey("test get value_models", t, func() {
		w, r = NewTestContextWithDefaultSession("GET", fmt.Sprintf("/api/v1/imdb/value_model/%v", 1), nil)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		So(w.Code, ShouldEqual, 200)
		log.Debug(respBody)
		w, r = NewTestContextWithDefaultSession("GET", fmt.Sprintf("/api/v1/imdb/value_model/%v?limit=1", 1), nil)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		So(w.Code, ShouldEqual, 200)
		log.Debug(respBody)
		rCheck = gjson.Get(respBody, "data.#")
		So(len(rCheck.Array()), ShouldEqual, 1)
		w, r = NewTestContextWithDefaultSession("GET", fmt.Sprintf("/api/v1/imdb/value_model/%v?q=c01.i02", 1), nil)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		So(w.Code, ShouldEqual, 200)
		log.Debug(respBody)
		rCheck = gjson.Get(respBody, "data.0.value")
		So(rCheck.String(), ShouldEqual, "c01.i02")
	})
}

func TestValueModelCreate(t *testing.T) {
	routes := SetUpGin()
	tagId := 1
	// clean test data
	db := config.Con()
	var vmd []imdb.ValueModel
	// tag := imdb.Tag{}
	aname := []string{"c01.i77"}
	db.IMDB.Where("tag_id = ? and value = ?", tagId, aname).Delete(&vmd)
	var (
		w *httptest.ResponseRecorder
		r *http.Request
	)
	Convey("create a new value_model", t, func() {
		postb := map[string]interface{}{
			"value": "c01.i77",
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", fmt.Sprintf("/api/v1/imdb/value_model/%d", tagId), &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		So(w.Code, ShouldEqual, 200)
		// insert the same data again, for test uniq check
		w, r = NewTestContextWithDefaultSession("POST", fmt.Sprintf("/api/v1/imdb/value_model/%d", tagId), &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 400)
	})
}

func TestValueModelUpdate(t *testing.T) {
	routes := SetUpGin()
	tagId := 1
	var (
		w *httptest.ResponseRecorder
		r *http.Request
	)
	Convey("update a new value_model", t, func() {
		postb := map[string]interface{}{
			"id":    1,
			"value": "c01.i09",
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("PUT", fmt.Sprintf("/api/v1/imdb/value_model/%d", tagId), &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		So(w.Code, ShouldEqual, 200)
		dvalue := gjson.Get(respBody, "data.value")
		So(dvalue.String(), ShouldEqual, "c01.i09")

		// test with update value with exist value & tag_id [test uniq key]
		w, r = NewTestContextWithDefaultSession("PUT", fmt.Sprintf("/api/v1/imdb/value_model/%d", tagId), &b)
		postb = map[string]interface{}{
			"id":    1,
			"value": "c01.i02",
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("PUT", fmt.Sprintf("/api/v1/imdb/value_model/%d", tagId), &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 400)
	})
}

func TestValueModelDelete(t *testing.T) {
	routes := SetUpGin()
	tagId := 1
	var (
		w *httptest.ResponseRecorder
		r *http.Request
	)
	Convey("delete value_models", t, func() {
		postb := map[string]interface{}{
			"value_model_ids": []int{1, 2},
			// this flag for test only
			"test": true,
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("DELETE", fmt.Sprintf("/api/v1/imdb/value_models/%d", tagId), &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
		Dno := gjson.Get(respBody, "data.deleted_number_of_object_tags")
		So(Dno.Int(), ShouldBeGreaterThan, 0)
	})
}
