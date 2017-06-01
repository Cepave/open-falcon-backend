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
	dg "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/controller/dashboard_graph"
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

func convertArrInterfaceToArrString(x []gjson.Result) (output []string) {
	output = make([]string, len(x))
	for indx, o := range x {
		output[indx] = o.String()
	}
	return
}

func TestDashboardGraphController(t *testing.T) {
	testGraphId := 4626
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
	Convey("test graph clone api", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("test with no graph id provided", func() {
			postb := map[string]interface{}{}
			b, _ := json.Marshal(postb)
			w, r = newTestContext("POST", "/api/v1/dashboard/graph_clone", &b)
			routes.ServeHTTP(w, r)
			log.Debug(w.Body.String())
			So(w.Code, ShouldEqual, 400)
		})
		Convey("test with not found graph id provided", func() {
			postb := map[string]interface{}{"id": 99998}
			b, _ := json.Marshal(postb)
			w, r = newTestContext("POST", "/api/v1/dashboard/graph_clone", &b)
			routes.ServeHTTP(w, r)
			log.Debug(w.Body.String())
			So(w.Code, ShouldEqual, 400)
		})
		Convey("test clone graph scuessful", func() {
			postb := map[string]interface{}{"id": testGraphId}
			b, _ := json.Marshal(postb)
			w, r = newTestContext("POST", "/api/v1/dashboard/graph_clone", &b)
			routes.ServeHTTP(w, r)
			log.Debug(w.Body.String())
			So(w.Code, ShouldEqual, 200)
			//for reset delete cloned screen & graph
			if w.Code == 200 {
				Convey("delete cloned graph", func() {
					res := struct {
						ID int `json:"id"`
					}{}
					json.Unmarshal(w.Body.Bytes(), &res)
					w, r = newTestContext("DELETE", fmt.Sprintf("/api/v1/dashboard/graph/%d", res.ID), nil)
					routes.ServeHTTP(w, r)
					log.Debug(w.Body.String())
					So(w.Code, ShouldEqual, 200)
				})
			}
		})
	})
	Convey("test update graph api", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("get graph by id", func() {
			w, r = newTestContext("GET", fmt.Sprintf("/api/v1/dashboard/graph/%d", testGraphId), nil)
			routes.ServeHTTP(w, r)
			log.Debug(w.Body.String())
			So(w.Code, ShouldEqual, 200)
			graphTmp := dg.APIDashboardGraphGetOuput{}
			json.Unmarshal(w.Body.Bytes(), &graphTmp)
			Convey("test update graph by id", func() {
				postb := dg.APIGraphUpdateReqData{
					Title:      "testtt",
					Endpoints:  []string{"a1", "a2"},
					Counters:   []string{"c1", "c2"},
					ID:         int64(testGraphId),
					TimeSpan:   int64(9999),
					GraphType:  "s",
					Method:     "sum",
					Position:   int64(22),
					FalconTags: "a=1,b=2",
				}
				b, _ := json.Marshal(postb)
				w, r = newTestContext("PUT", "/api/v1/dashboard/graph", &b)
				routes.ServeHTTP(w, r)
				log.Debug(w.Body.String())
				So(w.Code, ShouldEqual, 200)
				value := gjson.Get(w.Body.String(), "title")
				So(value.String(), ShouldEqual, postb.Title)
				value = gjson.Get(w.Body.String(), "endpoints")
				So(convertArrInterfaceToArrString(value.Array()), ShouldResemble, postb.Endpoints)
				value = gjson.Get(w.Body.String(), "counters")
				So(convertArrInterfaceToArrString(value.Array()), ShouldResemble, postb.Counters)
				value = gjson.Get(w.Body.String(), "timespan")
				So(value.Int(), ShouldResemble, postb.TimeSpan)
				value = gjson.Get(w.Body.String(), "graph_type")
				So(value.String(), ShouldResemble, postb.GraphType)
				value = gjson.Get(w.Body.String(), "method")
				So(value.String(), ShouldResemble, postb.Method)
				value = gjson.Get(w.Body.String(), "position")
				So(value.Int(), ShouldResemble, postb.Position)
				value = gjson.Get(w.Body.String(), "falcon_tags")
				So(value.String(), ShouldResemble, postb.FalconTags)
			})
			Reset(func() {
				postb := dg.APIGraphUpdateReqData{
					ID:         graphTmp.GraphID,
					Title:      graphTmp.Title,
					Endpoints:  graphTmp.Endpoints,
					Counters:   graphTmp.Counters,
					TimeSpan:   graphTmp.TimeSpan,
					GraphType:  graphTmp.GraphType,
					Method:     graphTmp.Method,
					Position:   graphTmp.Position,
					FalconTags: graphTmp.FalconTags,
				}
				b, _ := json.Marshal(postb)
				w, r = newTestContext("PUT", "/api/v1/dashboard/graph", &b)
				routes.ServeHTTP(w, r)
				log.Debug(w.Body.String())
			})
		})
	})
}
