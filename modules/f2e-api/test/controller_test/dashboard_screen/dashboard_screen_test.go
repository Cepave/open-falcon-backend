package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"bytes"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	yaagGin "github.com/masato25/yaag/gin"
	"github.com/masato25/yaag/yaag"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

//https://github.com/gin-gonic/gin/issues/580#issuecomment-242168245
func newTestContext(method, path string, body *[]byte) (w *httptest.ResponseRecorder, r *http.Request) {
	w = httptest.NewRecorder()
	if body == nil {
		r, _ = http.NewRequest(method, path, bytes.NewBuffer(nil))
	} else {
		r, _ = http.NewRequest(method, path, bytes.NewBuffer(*body))
	}
	//set json post as default
	r.Header.Set("Content-Type", "application/json")
	r.PostForm = url.Values{}
	setSession(r)
	return
}

func setSession(r *http.Request) {
	r.Header.Set("Apitoken", "{\"name\":\"root\",\"sig\":\"590c949340f811e7b6de001500c6ca5a\"}")
}
func TestDashboardScreenController(t *testing.T) {
	viper.Set("services", map[string]interface{}{
		"test":  "test123",
		"test2": "test456",
	})
	viper.Set("enable_services", true)
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../")
	viper.SetConfigName("cfg_test")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err.Error())
	}
	gin.SetMode(gin.TestMode)
	log.SetLevel(log.DebugLevel)
	config.InitDB(viper.GetBool("db.db_debug"))
	config.Init()
	routes := gin.Default()
	if viper.GetBool("gen_doc") {
		yaag.Init(&yaag.Config{
			On:       true,
			DocTitle: "Gin",
			DocPath:  viper.GetString("gen_doc_path"),
			BaseUrls: map[string]string{"Production": "/api/v1", "Staging": "/api/v1"},
		})
		routes.Use(yaagGin.Document())
	}
	routes = controller.StartGin(":8088", routes, true)
	Convey("test screen clone api", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("test with no screen id provided", func() {
			postb := map[string]interface{}{}
			b, _ := json.Marshal(postb)
			w, r = newTestContext("POST", "/api/v1/dashboard/screen_clone", &b)
			routes.ServeHTTP(w, r)
			log.Debug(w.Body.String())
			So(w.Code, ShouldEqual, 400)
		})
		Convey("test with not found screen id provided", func() {
			postb := map[string]interface{}{"id": 99998}
			b, _ := json.Marshal(postb)
			w, r = newTestContext("POST", "/api/v1/dashboard/screen_clone", &b)
			routes.ServeHTTP(w, r)
			log.Debug(w.Body.String())
			So(w.Code, ShouldEqual, 400)
		})
		Convey("test clone screen scuessful", func() {
			postb := map[string]interface{}{"id": 965}
			b, _ := json.Marshal(postb)
			w, r = newTestContext("POST", "/api/v1/dashboard/screen_clone", &b)
			routes.ServeHTTP(w, r)
			log.Debug(w.Body.String())
			So(w.Code, ShouldEqual, 200)
			//for reset delete cloned screen & graph
			Reset(func() {
				value := gjson.Get(w.Body.String(), "id")
				w, r = newTestContext("DELETE", fmt.Sprintf("/api/v1/dashboard/screen/%d", value.Int()), nil)
				routes.ServeHTTP(w, r)
				log.Debug(w.Body.String())
			})
		})
	})
}
