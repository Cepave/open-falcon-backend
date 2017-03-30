package mvc

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/common/utils"
	ot "github.com/Cepave/open-falcon-backend/common/types"
	otest "github.com/Cepave/open-falcon-backend/common/testing"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestStructTagsSuite struct{}

var _ = Suite(&TestStructTagsSuite{})

// Tests the building of loader
func (suite *TestStructTagsSuite) TestBuildParamLoader(c *C) {
	type sampleCar struct {
		F1 string `mvc:"query[easy-1]"`
		F11 string `mvc:"query[easy-2] default[something-1]"` // Default value
		F2 int `mvc:"cookie[kc1]"`
		F3 float32 `mvc:"param[kc1]"`
		F4 int8 `mvc:"form[kc1]"`
		F5 string `mvc:"header[sure]"`
		F6 string `mvc:"req[ContentLength]"`
		F7 int32 `mvc:"key[col-1]"`
		F77 int32 `mvc:"key[col-2] default[-13]"` // Default value on *gin.Context key
		FA1 []int32 `mvc:"query[easy-1]"` // Normal array
		FA11 []uint8 `mvc:"query[noV] default[9,8,13]"` // default value of array
		FA2 []string `mvc:"key[ks1]"` // default value on *gin.Context key
		FA22 []uint16 `mvc:"key[ks2] default[12,76,33]"` // default value on *gin.Context key
		FE1 int `mvc:"query[c1]"` // empty to default value of golang
		None string
	}

	testCases := []*struct {
		fieldName string
		setupFunc func(*gin.Context)
		expectedValue interface{}
	} {
		{
			"None",
			func(ginC *gin.Context) {},
			"",
		},
		{
			"F1",
			func(ginC *gin.Context) {
				ginC.Request = &http.Request{
					URL: otest.ParseRequestUri(c, "/query?easy-1=38"),
				}
			},
			"38",
		},
		{
			"F11",
			func(ginC *gin.Context) {
				ginC.Request = &http.Request{
					URL: otest.ParseRequestUri(c, "/query"),
				}
			},
			"something-1",
		},
		{
			"F2",
			func(ginC *gin.Context) {
				ginC.Request = &http.Request{ Header: make(http.Header) }
				ginC.Request.Header.Add("Cookie", "kc1=129")
			},
			129,
		},
		{
			"F3",
			func(ginC *gin.Context) {
				ginC.Params = make([]gin.Param, 1)
				ginC.Params[0] = gin.Param { Key: "kc1", Value: "54.61" }
			},
			float32(54.61),
		},
		{
			"F4",
			func(ginC *gin.Context) {
				ginC.Request = &http.Request{}
				ginC.Request.PostForm, _ = url.ParseQuery("kc1=17")
			},
			int8(17),
		},
		{
			"F5",
			func(ginC *gin.Context) {
				ginC.Request = &http.Request{ Header: make(http.Header) }
				ginC.Request.Header.Set("sure", "gc-091-KC1")
			},
			"gc-091-KC1",
		},
		{
			"F6",
			func(ginC *gin.Context) {
				ginC.Request = &http.Request{}
				ginC.Request.ContentLength = 98761
			},
			"98761",
		},
		{
			"F7",
			func(ginC *gin.Context) {
				ginC.Set("col-1", "981")
			},
			int32(981),
		},
		{
			"F77",
			func(ginC *gin.Context) {},
			int32(-13),
		},
		{
			"FA1",
			func(ginC *gin.Context) {
				ginC.Request = &http.Request{
					URL: otest.ParseRequestUri(c, "/query?easy-1=38&easy-1=91"),
				}
			},
			[]int32 { 38, 91 },
		},
		{
			"FA11",
			func(ginC *gin.Context) {
				ginC.Request = &http.Request{
					URL: otest.ParseRequestUri(c, "/query"),
				}
			},
			[]uint8 { 9, 8, 13 },
		},
		{
			"FA2",
			func(ginC *gin.Context) {
				ginC.Set("ks1", []int{ 81, 71, 62 })
			},
			[]string { "81", "71", "62" },
		},
		{
			"FA22",
			func(ginC *gin.Context) {},
			[]uint16 { 12, 76, 33 },
		},
		{
			"FE1",
			func(ginC *gin.Context) {
				ginC.Request = &http.Request{
					URL: otest.ParseRequestUri(c, "/query"),
				}
			},
			int(0),
		},
	}

	sampleType := reflect.TypeOf(sampleCar{})
	convSrv := ot.NewDefaultConversionService()
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		fieldType, _ := sampleType.FieldByName(testCase.fieldName)
		paramLoader := buildParamLoader(fieldType, convSrv)

		context := &gin.Context {}
		testCase.setupFunc(context)

		if paramLoader == nil {
			continue
		}
		testedValue := paramLoader(context)

		c.Assert(testedValue, DeepEquals, testCase.expectedValue, comment)
	}
}

// Tests the loading of page
func (suite *TestStructTagsSuite) TestPaging(c *C) {
	type pagingSample struct {
		P1 *model.Paging `mvc:"pageSize[17] pageOrderBy[ak_1:bd_2]"`
	}

	testCases := []*struct {
		requestSetup func(*http.Request)
		expectedSize int32
		expectedOrderBy []*model.OrderByEntity
	} {
		{
			func(req *http.Request) {},
			17,
			[]*model.OrderByEntity {
				{ Expr: "ak_1", Direction: utils.DefaultDirection },
				{ Expr: "bd_2", Direction: utils.DefaultDirection },
			},
		},
		{
			func(req *http.Request) {
				req.Header = make(http.Header)
				req.Header.Set("page-size", "39")
				req.Header.Set("order-by", "cp_1#asc:cp_2#desc")
			},
			39,
			[]*model.OrderByEntity {
				{ Expr: "cp_1", Direction: utils.Ascending },
				{ Expr: "cp_2", Direction: utils.Descending },
			},
		},
	}

	fieldType, _ := reflect.TypeOf(pagingSample{}).FieldByName("P1")
	convSrv := ot.NewDefaultConversionService()

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		context := &gin.Context {
			Request: httptest.NewRequest(http.MethodPost, "/", nil),
		}
		testCase.requestSetup(context.Request)

		paramLoader := buildParamLoader(fieldType, convSrv)
		testedPaging := paramLoader(context).(*model.Paging)

		c.Logf("Result paging: %#v", testedPaging)
		c.Assert(testedPaging.Size, Equals, testCase.expectedSize, comment)
		c.Assert(testedPaging.OrderBy, DeepEquals, testCase.expectedOrderBy, comment)
	}
}
