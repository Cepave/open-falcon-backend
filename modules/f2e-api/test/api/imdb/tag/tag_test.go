package test

import (
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/imdb"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tidwall/gjson"
)

func TestTagCreate(t *testing.T) {
	routes := SetUpGin()
	// clean test data
	db := config.Con()
	var tags []imdb.Tag
	// tag := imdb.Tag{}
	aname := []string{"tag1"}
	db.IMDB.Where("name in (?)", aname).Delete(&tags)
	var (
		w *httptest.ResponseRecorder
		r *http.Request
	)
	Convey("create a new tag", t, func() {
		postb := map[string]interface{}{
			"name":        "tag1",
			"tag_type_id": 1,
			"description": "hello test",
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/imdb/tag", &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		So(w.Code, ShouldEqual, 200)
		// insert the same data again, for test uniq check
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/imdb/tag", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 400)
		// test tag type not existing
		postb = map[string]interface{}{
			"name":        "tag2",
			"tag_type_id": 99,
			"description": "hello test",
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/imdb/tag", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 400)
	})
}

func TestTagUpdate(t *testing.T) {
	routes := SetUpGin()
	var (
		w *httptest.ResponseRecorder
		r *http.Request
	)
	Convey("create a new tag", t, func() {
		postb := map[string]interface{}{
			"description": "hello test2",
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("PUT", fmt.Sprintf("/api/v1/imdb/tag/%d", 2), &b)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)
		gjson := gjson.Get(respBody, "data.description")
		So(gjson.String(), ShouldEqual, "hello test2")
	})
}

func TestTagDelete(t *testing.T) {
	routes := SetUpGin()
	var (
		w *httptest.ResponseRecorder
		r *http.Request
	)
	var tag imdb.Tag
	db := config.Con()
	db.IMDB.Where("name = ?", "tag1").Find(&tag)
	// if record not exist, create one for testing
	if tag.ID == 0 {
		tag = imdb.Tag{
			Name:      "tag1",
			TagTypeId: 1,
		}
		db.IMDB.Save(&tag)
	}
	Convey("create a new tag", t, func() {
		w, r = NewTestContextWithDefaultSession("DELETE", fmt.Sprintf("/api/v1/imdb/tag/%d", tag.ID), nil)
		routes.ServeHTTP(w, r)
		respBody := w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)

		// test - can not delete a default tag
		w, r = NewTestContextWithDefaultSession("DELETE", fmt.Sprintf("/api/v1/imdb/tag/%d", 5), nil)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 400)
	})
}
