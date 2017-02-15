package mvc

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"net/http/httptest"
	"strings"
	"encoding/json"

	"gopkg.in/gin-gonic/gin.v1"
	ogin "github.com/Cepave/open-falcon-backend/common/gin"
)

func ExampleMvcBuilder_BuildHandler_httpGet() {
	mvcBuilder := NewMvcBuilder(NewDefaultMvcConfig())

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.GET(
		"/get-1",
		mvcBuilder.BuildHandler(
			func(
				data *struct {
					V1 int8 `mvc:"query[v1]"`
					V2 int32 `mvc:"query[v2]"`
				},
			) string {
				return fmt.Sprintf("V1: %d. V2: %d", data.V1, data.V2)
			},
		),
	)

	req := httptest.NewRequest(http.MethodGet, "/get-1?v1=20&v2=40", nil)
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)

	fmt.Println(resp.Body.String())

	// Output:
	// V1: 20. V2: 40
}
func ExampleMvcBuilder_BuildHandler_httpPost() {
	mvcBuilder := NewMvcBuilder(NewDefaultMvcConfig())

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.POST(
		"/post-1",
		mvcBuilder.BuildHandler(
			func(
				data *struct {
					V1 int8 `mvc:"form[v1]"`
					V2 []int32 `mvc:"form[v2]"`
				},
			) string {
				return fmt.Sprintf("v1: %d. v2: %d,%d", data.V1, data.V2[0], data.V2[1])
			},
		),
	)

	/**
	 * Form data
	 */
	form := url.Values {
		"v1": []string { "17" },
		"v2": []string { "230", "232" },
	}
	// :~)

	req := httptest.NewRequest(http.MethodPost, "/post-1", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)

	fmt.Println(resp.Body.String())

	// Output:
	// v1: 17. v2: 230,232
}

type sampleCar struct {
	Name string `json:"name"`
	Age int `json:"age"`
}
func (car *sampleCar) Bind(c *gin.Context) {
	ogin.BindJson(c, car)
}
func ExampleMvcBuilder_BuildHandler_json() {
	/*
	type sampleCar struct {
		Name string `json:"name"`
		Age int `json:"age"`
	}
	func (car *sampleCar) Bind(c *gin.Context) {
		ogin.BindJson(c, car)
	}
	*/

	mvcBuilder := NewMvcBuilder(NewDefaultMvcConfig())

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.POST(
		"/json-1",
		mvcBuilder.BuildHandler(
			func(car *sampleCar) OutputBody {
				return JsonOutputBody(car)
			},
		),
	)

	rawJson, _ := json.Marshal(&sampleCar{ "GTA-99", 3 })

	req := httptest.NewRequest(http.MethodPost, "/json-1", bytes.NewReader(rawJson))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)

	fmt.Println(resp.Body.String())

	// Output:
	// {"name":"GTA-99","age":3}
}
