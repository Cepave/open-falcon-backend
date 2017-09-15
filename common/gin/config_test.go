package gin

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	sjson "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/h2non/gentleman.v2"

	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ch "gopkg.in/check.v1"
)

var _ = Describe("Redirect http.Server(listening) to gin engine", func() {
	Context("Serve by \"net/http\", redirect to Gin engine", func() {
		sampleBody := `Hello World!!`

		BeforeEach(func() {
			/**
			 * The gin engine used to process real request of HTTP.
			 */
			sampleGin := NewDefaultJsonEngine(&GinConfig{Mode: gin.ReleaseMode})
			sampleGin.GET("/", func(context *gin.Context) {
				context.String(http.StatusOK, sampleBody)
			})
			// :~)

			sampleHandler := func(resp http.ResponseWriter, req *http.Request) {
				// Delegates to Gin engine
				sampleGin.ServeHTTP(resp, req)

				/**
				 * Ordinal code of http handler
				 */
				//resp.Header().Add("Content-Type", "text/plain")
				//resp.Write([]byte(sampleBody))
				// :~)
			}

			go http.ListenAndServe(
				":20301", http.HandlerFunc(sampleHandler),
			)
		})

		It("Should be 200 status and the body must be \"Hello World!!\"", func() {
			resp, err := gentleman.New().URL("http://127.0.0.1:20301").
				Get().
				Send()

			Expect(err).To(Succeed())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			Expect(resp.String()).To(Equal(sampleBody))
		})
	})
})

type TestConfigSuite struct{}

var _ = ch.Suite(&TestConfigSuite{})

// Tests the JSON engine for CORS and exception handler, etc.
func (suite *TestConfigSuite) TestNewDefaultJsonEngine(c *ch.C) {
	testCases := []*struct {
		req        *http.Request
		assertFunc func(*ch.C, *httptest.ResponseRecorder, ch.CommentInterface)
	}{
		{ // Tests the CORS
			httptest.NewRequest(http.MethodOptions, "/simple-1", nil),
			func(c *ch.C, resp *httptest.ResponseRecorder, comment ch.CommentInterface) {
				c.Assert(resp.Header().Get("Access-Control-Allow-Origin"), ch.Equals, "*", comment)
			},
		},
		{ // Tests the error of json binding
			httptest.NewRequest(http.MethodPost, "/json-error-1", nil),
			func(c *ch.C, resp *httptest.ResponseRecorder, comment ch.CommentInterface) {
				c.Assert(resp.Code, ch.Equals, http.StatusBadRequest, comment)

				jsonResult, err := sjson.NewFromReader(resp.Body)
				c.Assert(err, ch.IsNil)

				c.Assert(jsonResult.Get("error_code").MustInt(), ch.Equals, -101, comment)
				c.Assert(jsonResult.Get("error_message").MustString(), ch.Equals, "EOF", comment)
			},
		},
		{ // Tests the error of validation
			httptest.NewRequest(http.MethodGet, "/validation-error-1", nil),
			func(c *ch.C, resp *httptest.ResponseRecorder, comment ch.CommentInterface) {
				c.Assert(resp.Code, ch.Equals, http.StatusBadRequest, comment)

				jsonResult, err := sjson.NewFromReader(resp.Body)
				c.Assert(err, ch.IsNil)

				c.Assert(jsonResult.Get("error_code").MustInt(), ch.Equals, -1, comment)
				c.Assert(jsonResult.Get("error_message").MustString(), ch.Matches, ".*Error:Field validation.*", comment)
			},
		},
		{ // Tests the error of panic
			httptest.NewRequest(http.MethodGet, "/panic-1", nil),
			func(c *ch.C, resp *httptest.ResponseRecorder, comment ch.CommentInterface) {
				c.Assert(resp.Code, ch.Equals, http.StatusInternalServerError, comment)

				jsonResult, err := sjson.NewFromReader(resp.Body)
				c.Assert(err, ch.IsNil)

				c.Assert(jsonResult.Get("error_code").MustInt(), ch.Equals, -1, comment)
				c.Assert(jsonResult.Get("error_message").MustString(), ch.Equals, "HERE WE PANIC!!", comment)
			},
		},
		{ // Tests the error of not-found
			httptest.NewRequest(http.MethodGet, "/not-found", nil),
			func(c *ch.C, resp *httptest.ResponseRecorder, comment ch.CommentInterface) {
				c.Assert(resp.Code, ch.Equals, http.StatusNotFound, comment)

				jsonResult, err := sjson.NewFromReader(resp.Body)
				c.Assert(err, ch.IsNil)

				c.Assert(jsonResult.Get("error_code").MustInt(), ch.Equals, -1, comment)
				c.Assert(jsonResult.Get("uri").MustString(), ch.Equals, "/not-found", comment)
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
