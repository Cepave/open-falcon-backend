package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"bytes"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	log "github.com/sirupsen/logrus"
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

func convertArrInterfaceToArrString(x []gjson.Result) (output []string) {
	output = make([]string, len(x))
	for indx, o := range x {
		output[indx] = o.String()
	}
	return
}

func TestDashboardGraphController(t *testing.T) {
	// testGraphId := 4626
	viper.Set("services", map[string]interface{}{
		"test":  "test123",
		"test2": "test456",
	})
	viper.Set("enable_services", true)
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../../")
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
	Convey("create graph with existing id", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("invaild values", func() {
			postb := map[string]interface{}{
				"screen_id":     1247,
				"title":         "testiv",
				"endpoints":     []string{},
				"counters":      []string{},
				"graph_type":    "notv",
				"method":        "notv",
				"time_range":    "notv",
				"y_scale":       "notv",
				"sample_method": "notv",
				"sort_by":       "notv",
			}
			b, _ := json.Marshal(postb)
			w, r = newTestContext("POST", "/api/v1/dashboard/graph", &b)
			routes.ServeHTTP(w, r)
			So(w.Body.String(), ShouldEqual, "")
			So(w.Code, ShouldEqual, 400)
		})
	})
}
