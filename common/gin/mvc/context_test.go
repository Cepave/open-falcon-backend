package mvc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/Cepave/open-falcon-backend/common/model"
	ot "github.com/Cepave/open-falcon-backend/common/types"
	sjson "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestContextSuite struct{}

var _ = Suite(&TestContextSuite{})

// Tests the building of handler
func (suite *TestContextSuite) TestBuildHandler(c *C) {
	testCases := []*struct {
		url         string
		handlerFunc interface{}
		req         *http.Request
		assertFunc  func(*C, *httptest.ResponseRecorder, CommentInterface)
	}{
		{ // text/plain output
			"/plain-1",
			func() string {
				return "Simple-OK"
			},
			httptest.NewRequest(http.MethodGet, "/plain-1", nil),
			func(c *C, r *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(r.Code, Equals, http.StatusOK, comment)
				c.Assert(r.Body.String(), Equals, "Simple-OK", comment)
			},
		},
		{ // Json(by json.Marshaler)
			"/json-1",
			func() *sjson.Json {
				jsonBody := sjson.New()
				jsonBody.Set("name", "Blue-Jary")
				jsonBody.Set("weight", 980)

				return jsonBody
			},
			httptest.NewRequest(http.MethodGet, "/json-1", nil),
			func(c *C, r *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(r.Code, Equals, http.StatusOK, comment)

				testedResult := &car{}

				c.Assert(json.Unmarshal(r.Body.Bytes(), testedResult), IsNil, comment)
				c.Assert(testedResult.Name, Equals, "Blue-Jary", comment)
				c.Assert(testedResult.Weight, Equals, 980, comment)
			},
		},
		{ // Json(by JsonOutputBody)
			"/json-2",
			func() OutputBody {
				return JsonOutputBody(&car{Name: "Start-V8", Weight: 1280})
			},
			httptest.NewRequest(http.MethodGet, "/json-2", nil),
			func(c *C, r *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(r.Code, Equals, http.StatusOK, comment)

				testedResult := &car{}

				c.Assert(json.Unmarshal(r.Body.Bytes(), testedResult), IsNil, comment)
				c.Assert(testedResult.Name, Equals, "Start-V8", comment)
				c.Assert(testedResult.Weight, Equals, 1280, comment)
			},
		},
		{ // No return value
			"/no-return",
			func(context *gin.Context) {
				context.String(http.StatusOK, "Hello-83")
			},
			httptest.NewRequest(http.MethodGet, "/no-return", nil),
			func(c *C, r *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(r.Code, Equals, http.StatusOK, comment)
				c.Assert(r.Body.String(), Equals, "Hello-83", comment)
			},
		},
		{ // All-supported input parameters on web objects
			"/full-input/:cg_id",
			func(
				context *gin.Context,
				req *http.Request, resp http.ResponseWriter, url *url.URL,
				header http.Header,
				ginParams gin.Params, ginRespWriter gin.ResponseWriter,
				convSrv ot.ConversionService, validate *validator.Validate,
			) OutputBody {
				jsonBody := sjson.New()

				/**
				 * Binding variable of input
				 */
				jsonBody.Set("context", context != nil)
				jsonBody.Set("req", req != nil)
				jsonBody.Set("http_resp_writer", resp != nil)
				jsonBody.Set("gin_resp_writer", ginRespWriter != nil)
				jsonBody.Set("url", fmt.Sprintf("%s", url))
				jsonBody.Set("form_fv_1", context.PostForm("fv-1"))
				jsonBody.Set("cg_id", ginParams.ByName("cg_id"))
				jsonBody.Set("conv_srv", convSrv != nil)
				jsonBody.Set("validate", validate != nil)
				// :~)

				// Testing header
				header.Set("rh1", "PB-02")

				return JsonOutputBody(jsonBody)
			},
			func() *http.Request {
				form := url.Values{
					"fv-1": []string{"671"},
				}

				req := httptest.NewRequest(http.MethodPost, "/full-input/9810", strings.NewReader(form.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			}(),
			func(c *C, r *httptest.ResponseRecorder, comment CommentInterface) {
				jsonResult, err := sjson.NewFromReader(r.Body)
				c.Assert(err, IsNil, comment)

				c.Assert(jsonResult.Get("context").MustBool(), Equals, true, comment)
				c.Assert(jsonResult.Get("req").MustBool(), Equals, true, comment)
				c.Assert(jsonResult.Get("http_resp_writer").MustBool(), Equals, true, comment)
				c.Assert(jsonResult.Get("gin_resp_writer").MustBool(), Equals, true, comment)
				c.Assert(jsonResult.Get("url").MustString(), Equals, "/full-input/9810", comment)
				c.Assert(jsonResult.Get("form_fv_1").MustString(), Equals, "671", comment)
				c.Assert(jsonResult.Get("cg_id").MustString(), Equals, "9810", comment)
				c.Assert(jsonResult.Get("conv_srv").MustBool(), Equals, true, comment)
				c.Assert(jsonResult.Get("validate").MustBool(), Equals, true, comment)

				c.Assert(r.Header().Get("rh1"), Equals, "PB-02", comment)
			},
		},
		{ // *multipart.Form
			"/file-upload",
			func(form *multipart.Form) *sjson.Json {
				jsonResult := sjson.New()

				jsonResult.Set("form_file", form.File["fn1"][0].Filename)
				jsonResult.Set("form_values", len(form.Value))

				return jsonResult
			},
			func() *http.Request {
				buffer := new(bytes.Buffer)
				multipartWriter := multipart.NewWriter(buffer)

				/**
				 * File upload
				 */
				makeFile(
					c, multipartWriter,
					"fn1", "sample.txt", "This is word!!!",
				)
				// :~)

				/**
				 * Field value
				 */
				multipartWriter.WriteField("gv1", "hello-form-1")
				multipartWriter.WriteField("gv2", "hello-form-2")
				// :~)

				c.Assert(multipartWriter.Close(), IsNil)

				req := httptest.NewRequest(http.MethodPost, "/file-upload", buffer)
				req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

				return req
			}(),
			func(c *C, r *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(r.Code, Equals, http.StatusOK, comment)

				jsonResult, err := sjson.NewFromReader(r.Body)
				c.Assert(err, IsNil)
				c.Assert(jsonResult.Get("form_file").MustString(), Equals, "sample.txt", comment)
				c.Assert(jsonResult.Get("form_values").MustInt(), Equals, 2, comment)
			},
		},
		{ // file by struct tag
			"/sfile-upload",
			func(
				f *struct {
					SingleFile    multipart.File   `mvc:"file[csv_f_1]"`
					MultipleFiles []multipart.File `mvc:"file[csv_f_2]"`
				},
			) *sjson.Json {
				jsonResult := sjson.New()

				b := new(bytes.Buffer)
				b.ReadFrom(f.SingleFile)

				jsonResult.Set("single_file_size", b.Len())
				jsonResult.Set("multiple_files", len(f.MultipleFiles))

				return jsonResult
			},
			func() *http.Request {
				buffer := new(bytes.Buffer)
				multipartWriter := multipart.NewWriter(buffer)

				/**
				 * File upload
				 */
				makeFile(
					c, multipartWriter,
					"csv_f_1", "sample.txt", "Single file",
				)
				makeFile(
					c, multipartWriter,
					"csv_f_2", "sample-1.txt", "M-file-1",
				)
				makeFile(
					c, multipartWriter,
					"csv_f_2", "sample-2.txt", "M-file-2",
				)
				// :~)

				c.Assert(multipartWriter.Close(), IsNil)

				req := httptest.NewRequest(http.MethodPost, "/sfile-upload", buffer)
				req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

				return req
			}(),
			func(c *C, r *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(r.Code, Equals, http.StatusOK, comment)

				jsonResult, err := sjson.NewFromReader(r.Body)
				c.Assert(err, IsNil)
				c.Assert(jsonResult.Get("single_file_size").MustInt(), Equals, 11, comment)
				c.Assert(jsonResult.Get("multiple_files").MustInt(), Equals, 2, comment)
			},
		},
		{ // Binding with json.Unmarshaler
			"/un-json",
			func(car3 *car3) string {
				return fmt.Sprintf("%d-%s", car3.Key, car3.V1)
			},
			func() *http.Request {
				jsonObject := sjson.New()
				jsonObject.Set("key", 383)
				jsonObject.Set("v1", "jp-109")
				rawJson, _ := jsonObject.MarshalJSON()

				req := httptest.NewRequest(http.MethodPost, "/un-json", bytes.NewReader(rawJson))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
			func(c *C, r *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(r.Code, Equals, http.StatusOK, comment)
				c.Assert(r.Body.String(), Equals, "383-jp-109", comment)
			},
		},
		{ // Binding for parameters
			"/all_params/:gd_id",
			func(
				v *struct {
					ParamValue  uint16   `mvc:"param[gd_id]"`
					QueryParam  []int32  `mvc:"query[qa]"`
					CookieValue string   `mvc:"cookie[session_id]"`
					FormAge     []uint16 `mvc:"form[age]"`
					FormAgentId []uint32 `mvc:"form[user_id]"`
					Header      string   `mvc:"header[hv]"`
					Method      string   `mvc:"req[Method]"`
					CarV        *car     `mvc:"key[key-cc-01]"`
					Username    string   `mvc:"basicAuth[username]"`
					Password    string   `mvc:"basicAuth[password]"`
				},
			) *sjson.Json {
				jsonResult := sjson.New()
				jsonResult.Set("param_value", v.ParamValue)
				jsonResult.Set("query_value", v.QueryParam)
				jsonResult.Set("cookie_value", v.CookieValue)

				jsonResult.Set("form_age", v.FormAge)
				jsonResult.Set("form_agent_id", v.FormAgentId)
				jsonResult.Set("header", v.Header)
				jsonResult.Set("method", v.Method)
				jsonResult.Set("car", v.CarV)

				jsonResult.Set("username", v.Username)
				jsonResult.Set("password", v.Password)

				return jsonResult
			},
			func() *http.Request {
				form := url.Values{
					"age":     []string{"6", "33"},
					"user_id": []string{"131", "8716"},
				}

				req := httptest.NewRequest(http.MethodPost, "/all_params/8806?qa=98&qa=99", strings.NewReader(form.Encode()))
				req.SetBasicAuth("ug@mail.com", "pac@9081")
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.Header.Set("hv", "32768")
				req.AddCookie(&http.Cookie{Name: "session_id", Value: "pk228c122"})

				return req
			}(),
			func(c *C, r *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(r.Code, Equals, http.StatusOK, comment)

				jsonResult, err := sjson.NewFromReader(r.Body)
				c.Assert(err, IsNil, comment)

				// URL param value
				c.Assert(jsonResult.Get("param_value").MustUint64(), Equals, uint64(8806), comment)

				// Cookie value
				c.Assert(jsonResult.Get("cookie_value").MustString(), Equals, "pk228c122", comment)

				/**
				 * Query value
				 */
				queryValues := jsonResult.Get("query_value").MustArray()
				c.Assert(queryValues[0], Equals, json.Number("98"), comment)
				c.Assert(queryValues[1], Equals, json.Number("99"), comment)
				// :~)

				/**
				 * Form values
				 */
				formAge := jsonResult.Get("form_age").MustArray()
				c.Assert(formAge[0], Equals, json.Number("6"), comment)
				c.Assert(formAge[1], Equals, json.Number("33"), comment)
				formAgentId := jsonResult.Get("form_agent_id").MustArray()
				c.Assert(formAgentId[0], Equals, json.Number("131"), comment)
				c.Assert(formAgentId[1], Equals, json.Number("8716"), comment)
				// :~)

				c.Assert(jsonResult.Get("header").MustString(), Equals, "32768", comment)
				c.Assert(jsonResult.Get("method").MustString(), Equals, "POST", comment)

				c.Assert(jsonResult.GetPath("car", "name").MustString(), Equals, "Cool!", comment)
				c.Assert(jsonResult.GetPath("car", "weight").MustInt(), Equals, 88, comment)

				c.Assert(jsonResult.GetPath("username").MustString(), Equals, "ug@mail.com", comment)
				c.Assert(jsonResult.GetPath("password").MustString(), Equals, "pac@9081", comment)
			},
		},
		{ // Binding for ContextBinder
			"/context-binder",
			func(c *car2) string {
				return fmt.Sprintf("%s-%d", c.name, c.age)
			},
			httptest.NewRequest(http.MethodPost, "/context-binder", nil),
			func(c *C, r *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(r.Code, Equals, http.StatusOK, comment)
				c.Assert(r.Body.String(), Equals, "context-name-55", comment)
			},
		},
		{ // Binding for paging
			"/test-paging",
			func(
				p *struct {
					Paging *model.Paging
				},
			) (string, *model.Paging) {
				p.Paging.PageMore = true
				p.Paging.TotalCount = 190

				return fmt.Sprintf("Size: %d", p.Paging.Size), p.Paging
			},
			httptest.NewRequest(http.MethodGet, "/test-paging", nil),
			func(c *C, r *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(r.Code, Equals, http.StatusOK, comment)

				c.Assert(r.Body.String(), Equals, "Size: 64", comment)
				c.Assert(r.Header().Get("total-count"), Equals, "190", comment)
				c.Assert(r.Header().Get("page-more"), Equals, "true", comment)
			},
		},
		{ // Nested struct
			"/test-nest-struct",
			func(
				p *struct {
					Name string `mvc:"query[name1]"`
					N1   *struct {
						Name1 string `mvc:"query[name2]"`
					}
					N2 struct {
						Name1 string `mvc:"query[name3]"`
					}
				},
			) OutputBody {
				return JsonOutputBody(p)
			},
			httptest.NewRequest(http.MethodGet, "/test-nest-struct?name1=v10c10&name2=v10c20&name3=v10c30", nil),
			func(c *C, r *httptest.ResponseRecorder, comment CommentInterface) {
				c.Assert(r.Code, Equals, http.StatusOK, comment)

				jsonResult, err := sjson.NewFromReader(r.Body)
				c.Assert(err, IsNil, comment)

				c.Logf("Json of nested result: %v", jsonResult)

				c.Assert(jsonResult.Get("Name").MustString(), Equals, "v10c10")
				c.Assert(jsonResult.GetPath("N1", "Name1").MustString(), Equals, "v10c20")
				c.Assert(jsonResult.GetPath("N2", "Name1").MustString(), Equals, "v10c30")
			},
		},
	}

	testedBuilder := NewMvcBuilder(NewDefaultMvcConfig())

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		gin.SetMode(gin.ReleaseMode)
		engine := gin.Default()
		engine.Use(func(context *gin.Context) {
			context.Set("key-cc-01", &car{"Cool!", 88})
		})
		engine.Any(
			testCase.url,
			testedBuilder.BuildHandler(testCase.handlerFunc),
		)

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, testCase.req)

		c.Logf("Status: [%d]. Output: << %s >>", recorder.Code, recorder.Body)

		testCase.assertFunc(c, recorder, comment)
	}
}

// Tests the validation of struct value
func (suite *TestContextSuite) TestValidation(c *C) {
	testedBuilder := NewMvcBuilder(NewDefaultMvcConfig())

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	var err error
	engine.Use(
		func(c *gin.Context) {
			defer func() {
				p := recover()
				if p != nil {
					err = p.(error)
				}
			}()
			c.Next()
		},
	)
	engine.Any(
		"/input-validate",
		testedBuilder.BuildHandler(func(
			data *struct {
				Name string `mvc:"query[name]" validate:"min=4"`
			},
		) string {
			return "OK"
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/input-validate?name=bbb", nil)
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)

	c.Logf("Validation error: %v", err)
	c.Assert(err, NotNil)
}

type car struct {
	Name   string `json:"name"`
	Weight int    `json:"weight"`
}

type car2 struct {
	name string
	age  int
}

func (c *car2) Bind(context *gin.Context) {
	c.name = "context-name"
	c.age = 55
}

type car3 struct {
	Key int32
	V1  string
}

func (c *car3) UnmarshalJSON(rawJson []byte) error {
	jsonObject, err := sjson.NewJson(rawJson)
	if err != nil {
		return err
	}

	c.Key = int32(jsonObject.Get("key").MustInt())
	c.V1 = jsonObject.Get("v1").MustString()

	return nil
}

type box struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

func (b *box) UnmarshalJSON([]byte) error {
	return nil
}

type sampleParamBinding struct{}

func makeFile(
	c *C,
	writer *multipart.Writer,
	fieldName string, filename string, fileContent string,
) {
	fn1Writer, err := writer.CreateFormFile(fieldName, filename)
	c.Assert(err, IsNil)
	bufWriter := bufio.NewWriter(fn1Writer)
	bufWriter.WriteString(fileContent)
	c.Assert(bufWriter.Flush(), IsNil)
}
