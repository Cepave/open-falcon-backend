package test

import (
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	jconf "github.com/Cepave/open-falcon-backend/modules/f2e-api/lambda_extends/conf"
	. "github.com/Cepave/open-falcon-backend/modules/f2e-api/test_utils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func TestFuncAvgCompare(t *testing.T) {
	routes := SetUpGin()
	if viper.GetBool("lambda_extends.enable") {
		jconf.ReadConf()
	}
	Convey("query func", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		// create a template for action test
		postb := map[string]interface{}{
			"from":  1504240064,
			"until": 1504250000,
			"endpoints": []string{
				"hostA", "hostB", "hostC", "hostD", "hostE",
			},
			"metrics": []string{
				"cpu.idle",
			},
		}
		Convey("query with func - avgCompare with '>'", func() {
			postb2 := postb
			postb2["func"] = map[string]interface{}{
				"function": "avgCompare",
				"cond":     ">",
			}
			b, _ := json.Marshal(postb2)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/lambdaq/q", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			jp := gjson.Parse(respBody)
			enp := jp.Get("#.endpoint")
			So(enp.Array()[0].String(), ShouldEqual, "hostA")
			So(len(enp.Array()), ShouldEqual, 3)
			So(w.Code, ShouldEqual, 200)
		})
		Convey("query with func - avgCompare with '<'", func() {
			postb2 := postb
			postb2["func"] = map[string]interface{}{
				"function": "avgCompare",
				"cond":     "<",
			}
			b, _ := json.Marshal(postb2)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/lambdaq/q", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			jp := gjson.Parse(respBody)
			enp := jp.Get("#.endpoint")
			So(enp.Array()[0].String(), ShouldEqual, "hostC")
			So(len(enp.Array()), ShouldEqual, 2)
			So(w.Code, ShouldEqual, 200)
		})
	})
}
