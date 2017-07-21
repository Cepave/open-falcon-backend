package gin

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	sjson "github.com/bitly/go-simplejson"
	"gopkg.in/go-playground/validator.v9"

	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestConfigSuite struct{}

var _ = Suite(&TestConfigSuite{})

// Tests the JSON engine for CORS and exception handler, etc.
func (suite *TestConfigSuite) TestNewDefaultJsonEngine(c *C) {
	testCases := []*struct {
		req        *http.Request
		assertFunc func(*C, *httptest.ResponseRecorder, CommentInterface)
	}{
		{ // Tests the CORS
			httptest.NewRequest(http.MethodOptions, "/simple-1", nil),
			func(c *C, resp *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(resp.Header().Get("Access-Control-Allow-Origin"), Equals, "*", comment)
			},
		},
		{ // Tests the error of json binding
			httptest.NewRequest(http.MethodPost, "/json-error-1", nil),
			func(c *C, resp *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(resp.Code, Equals, http.StatusBadRequest, comment)

				jsonResult, err := sjson.NewFromReader(resp.Body)
				c.Assert(err, IsNil)

				c.Assert(jsonResult.Get("error_code").MustInt(), Equals, -101, comment)
				c.Assert(jsonResult.Get("error_message").MustString(), Equals, "EOF", comment)
			},
		},
		{ // Tests the error of validation
			httptest.NewRequest(http.MethodGet, "/validation-error-1", nil),
			func(c *C, resp *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(resp.Code, Equals, http.StatusBadRequest, comment)

				jsonResult, err := sjson.NewFromReader(resp.Body)
				c.Assert(err, IsNil)

				c.Assert(jsonResult.Get("error_code").MustInt(), Equals, -1, comment)
				c.Assert(jsonResult.Get("error_message").MustString(), Matches, ".*Error:Field validation.*", comment)
			},
		},
		{ // Tests the error of panic
			httptest.NewRequest(http.MethodGet, "/panic-1", nil),
			func(c *C, resp *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(resp.Code, Equals, http.StatusInternalServerError, comment)

				jsonResult, err := sjson.NewFromReader(resp.Body)
				c.Assert(err, IsNil)

				c.Assert(jsonResult.Get("error_code").MustInt(), Equals, -1, comment)
				c.Assert(jsonResult.Get("error_message").MustString(), Equals, "HERE WE PANIC!!", comment)
			},
		},
		{ // Tests the error of not-found
			httptest.NewRequest(http.MethodGet, "/not-found", nil),
			func(c *C, resp *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(resp.Code, Equals, http.StatusNotFound, comment)

				jsonResult, err := sjson.NewFromReader(resp.Body)
				c.Assert(err, IsNil)

				c.Assert(jsonResult.Get("error_code").MustInt(), Equals, -1, comment)
				c.Assert(jsonResult.Get("uri").MustString(), Equals, "/not-found", comment)
			},
		},
	}

	engine := NewDefaultJsonEngine(&GinConfig{Mode: gin.ReleaseMode})
	engine.GET("/sample-1", func(context *gin.Context) {
		context.String(http.StatusOK, "OK")
	})
	engine.POST("/json-error-1", func(context *gin.Context) {
		type car struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		BindJson(context, &car{})
	})
	engine.GET("/validation-error-1", func(context *gin.Context) {
		type car struct {
			Name string `validate:"min=10"`
		}

		ConformAndValidateStruct(&car{"cc"}, validator.New())
	})
	engine.GET("/panic-1", func(context *gin.Context) {
		panic("HERE WE PANIC!!")
	})

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testCase.req.Header.Set("Origin", "http://non-local/")

		resp := httptest.NewRecorder()
		engine.ServeHTTP(resp, testCase.req)

		testCase.assertFunc(c, resp, comment)
	}
}

func ExampleNewDefaultJsonEngine() {
	engine := NewDefaultJsonEngine(&GinConfig{Mode: gin.ReleaseMode})

	engine.GET(
		"/car/:car_id",
		func(c *gin.Context) {
			c.String(http.StatusOK, "Car Id: "+c.Param("car_id"))
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/car/91", nil)
	resp := httptest.NewRecorder()

	engine.ServeHTTP(resp, req)

	fmt.Println(resp.Body)

	// Output:
	// Car Id: 91
}
