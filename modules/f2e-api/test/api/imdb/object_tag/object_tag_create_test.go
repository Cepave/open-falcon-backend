package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/imdb"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestObjectTagCreate(t *testing.T) {
	routes := SetUpGin()
	var (
		w        *httptest.ResponseRecorder
		r        *http.Request
		postb    map[string]interface{}
		b        []byte
		respBody string
	)
	// clean test data
	db := config.Con()
	var objectTags []imdb.ObjectTag
	db.IMDB.Where("resource_object_id  = ?", 4).Delete(&objectTags)
	Convey("create object new tag", t, func() {
		// string type
		postb = map[string]interface{}{
			"tag_id":             5,
			"resource_object_id": 4,
			"value_text":         "中華台灣機房",
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/imdb/object_tag", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)

		// int type
		postb = map[string]interface{}{
			"tag_id":             4,
			"resource_object_id": 4,
			"value_int":          10,
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/imdb/object_tag", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)

		// description type
		postb = map[string]interface{}{
			"tag_id":             2,
			"resource_object_id": 4,
			"value_text":         "opsuser 000-000-000-000",
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/imdb/object_tag", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)

		// value model type
		postb = map[string]interface{}{
			"tag_id":             1,
			"resource_object_id": 4,
			"value_int":          2,
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/imdb/object_tag", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 200)

		Convey("create object faild with the same value", func() {
			// string type
			postb = map[string]interface{}{
				"tag_id":             5,
				"resource_object_id": 4,
				"value_text":         "中華台灣機房",
			}
			b, _ = json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/imdb/object_tag", &b)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			log.Debug(respBody)
			So(w.Code, ShouldEqual, 400)

			// int type
			postb = map[string]interface{}{
				"tag_id":             4,
				"resource_object_id": 4,
				"value_int":          10,
			}
			b, _ = json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/imdb/object_tag", &b)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			log.Debug(respBody)
			So(w.Code, ShouldEqual, 400)

			// description type
			postb = map[string]interface{}{
				"tag_id":             2,
				"resource_object_id": 4,
				"value_text":         "opsuser 000-000-000-000",
			}
			b, _ = json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/imdb/object_tag", &b)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			log.Debug(respBody)
			So(w.Code, ShouldEqual, 400)

			// value model type
			postb = map[string]interface{}{
				"tag_id":             1,
				"resource_object_id": 4,
				"value_int":          2,
			}
			b, _ = json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/imdb/object_tag", &b)
			routes.ServeHTTP(w, r)
			respBody = w.Body.String()
			log.Debug(respBody)
			So(w.Code, ShouldEqual, 400)
		})
	})
	Convey("create object faild with not exist value", t, func() {
		// string type
		postb = map[string]interface{}{
			"tag_id":             5,
			"resource_object_id": 4,
			"value_int":          1,
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/imdb/object_tag", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 400)

		// value model type
		postb = map[string]interface{}{
			"tag_id":             1,
			"resource_object_id": 4,
			"value_int":          999,
		}
		b, _ = json.Marshal(postb)
		w, r = NewTestContextWithDefaultSession("POST", "/api/v1/imdb/object_tag", &b)
		routes.ServeHTTP(w, r)
		respBody = w.Body.String()
		log.Debug(respBody)
		So(w.Code, ShouldEqual, 400)
	})
}
