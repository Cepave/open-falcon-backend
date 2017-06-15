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
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
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
	return
}

func TestApiToekn(t *testing.T) {
	viper.Set("services", map[string]interface{}{
		"test":  "test123",
		"test2": "test456",
	})
	viper.Set("enable_services", true)
	viper.AddConfigPath(".")
	viper.SetConfigFile("./cfg_test.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err.Error())
	}
	gin.SetMode(gin.TestMode)
	log.SetLevel(log.DebugLevel)
	config.InitDB(viper.GetBool("db.db_debug"))
	config.Init()
	routes := controller.StartGin(":8088", gin.Default(), true)
	apiClinet := config.ApiClient
	Convey("read token keys map", t, func() {
		So(len(apiClinet.Keys()), ShouldEqual, 2)
	})
	Convey("test token auth", t, func() {
		So(apiClinet.AuthToken("test", "test123"), ShouldEqual, true)
		So(apiClinet.AuthToken("notfound", "test123"), ShouldEqual, false)
		So(apiClinet.AuthToken("test", "notfoundtoken"), ShouldEqual, false)
		So(apiClinet.AuthToken("test2", "test456"), ShouldEqual, true)
	})
	Convey("Create User before test", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		// postb := map[string]interface{}{
		// 	"name":     "root",
		// 	"password": "testroot",
		// 	"email":    "xxx@xxx.com",
		// 	"cnname":   "myname",
		// }
		// b, _ := json.Marshal(postb)
		// w, r := newTestContext("POST", "/api/v1/user/create", &b)
		// routes.ServeHTTP(w, r)
		// fmt.Printf("body: %s", w.Body)
		// So(w.Code, ShouldEqual, 200)
		Convey("test session auth with token", func() {
			w, r = newTestContext("GET", "/api/v1/user/auth_session", nil)
			r.Header.Set("Apitoken", "{\"name\":\"root\",\"sig\":\"590c949340f811e7b6de001500c6ca5a\"}")
			routes.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, 200)
		})
		Convey("test session auth with servies token", func() {
			w, r = newTestContext("GET", "/api/v1/user/auth_session", nil)
			r.Header.Set("Apitoken", "{\"name\":\"test\",\"sig\":\"test123\"}")
			routes.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, 200)
			So(w.Body.String(), ShouldContainSubstring, "servies token")
		})
		Convey("test session auth with services not vaild", func() {
			w, r = newTestContext("GET", "/api/v1/user/auth_session", nil)
			r.Header.Set("Apitoken", "{\"name\":\"test\",\"sig\":\"notfoundtoken\"}")
			routes.ServeHTTP(w, r)
			So(w.Code, ShouldNotEqual, 200)
		})
	})
	Convey("test not allow action on apitoken", t, func() {
		postb := map[string]interface{}{"pid": 0, "name": "screen_test"}
		b, _ := json.Marshal(postb)
		w, r := newTestContext("POST", "/api/v1/dashboard/screen", &b)
		r.Header.Set("Apitoken", "{\"name\":\"test\",\"sig\":\"test123\"}")
		routes.ServeHTTP(w, r)
		So(w.Code, ShouldNotEqual, 200)
		So(w.Body.String(), ShouldContainSubstring, "services token no support")
	})
}
