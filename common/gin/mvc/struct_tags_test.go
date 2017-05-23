package mvc

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/common/utils"
	rt "github.com/Cepave/open-falcon-backend/common/reflect/types"
	ot "github.com/Cepave/open-falcon-backend/common/types"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestStructTagsSuite struct{}

var _ = Suite(&TestStructTagsSuite{})

// Tests the building of loader
func (suite *TestStructTagsSuite) TestBuildParamLoader(c *C) {
	getReqSetup := func(context *gin.Context) {
		context.Request = httptest.NewRequest("GET", "/query?qv-1=hello3&qv-2=56&nv=32&nv=76", nil)
		context.Request.Header.Add("Cookie", "ck_1=hello98")
		context.Request.Header.Add("Cookie", "ck_2=77")

		context.Request.Header.Add("hvv1", "hello24")
		context.Request.Header.Add("hvv2", "278")

		context.Set("kg-1", "hello45")
		context.Set("kg-2", "871")
		context.Set("kg-3", nil)
	}
	postReqSetup := func(context *gin.Context) {
		params := make(url.Values)
		params.Add("fv-1", "hello256")
		params.Add("fv-2", "23")
		params.Add("fv-2", "56")

		context.Request = httptest.NewRequest(
			"POST", "/query?qv-1=hello3&qv-2=56&nv=32&nv=76",
			strings.NewReader(params.Encode()),
		)
		context.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	testCases := []*struct {
		sampleType reflect.Type
		sampleTag reflect.StructTag
		setupFunc func(*gin.Context)
		expectedValue interface{}
	} {
		/**
		 * Query parameters
		 */
		{ // Viable value
			rt.TypeOfString, `mvc:"query[qv-1]"`,
			getReqSetup, "hello3",
		},
		{ // Default value
			rt.TypeOfString, `mvc:"query[qv-3] default[sde-13]"`,
			getReqSetup, "sde-13",
		},
		{ // Default value of go lang
			rt.TypeOfInt, `mvc:"query[qvn-3]"`,
			getReqSetup, 0,
		},
		{ // Viable value with type conversion
			rt.TypeOfInt16, `mvc:"query[qv-2]"`,
			getReqSetup, int16(56),
		},
		{ // Checked param(true)
			rt.TypeOfBool, `mvc:"query[?qv-1]"`,
			getReqSetup, true,
		},
		{ // Checked param(false)
			rt.TypeOfBool, `mvc:"query[?qv-3]"`,
			getReqSetup, false,
		},
		{ // Viable slice
			rt.STypeOfString, `mvc:"query[nv]"`,
			getReqSetup, []string{ "32", "76" },
		},
		{ // Vialbe slice with type conversion
			rt.STypeOfUint16, `mvc:"query[nv]"`,
			getReqSetup, []uint16{ 32, 76 },
		},
		{ // Default value
			rt.STypeOfUint32, `mvc:"query[nv9] default[88,76,39]"`,
			getReqSetup, []uint32{ 88, 76, 39 },
		},
		// :~)

		/**
		 * Form parameters
		 */
		{ // Viable value
			rt.TypeOfString, `mvc:"form[fv-1]"`,
			postReqSetup, "hello256",
		},
		{ // Default value
			rt.TypeOfString, `mvc:"form[fv-3] default[ffcc-013]"`,
			postReqSetup, "ffcc-013",
		},
		{ // Default value of go lang
			rt.TypeOfUint8, `mvc:"form[nff]"`,
			postReqSetup, uint8(0),
		},
		{ // Viable value with type conversion
			rt.TypeOfInt16, `mvc:"form[fv-2]"`,
			postReqSetup, int16(23),
		},
		{ // Checked param(true)
			rt.TypeOfBool, `mvc:"form[?fv-1]"`,
			postReqSetup, true,
		},
		{ // Checked param(false)
			rt.TypeOfBool, `mvc:"form[?nfv03]"`,
			postReqSetup, false,
		},
		{ // Viable slice
			rt.STypeOfString, `mvc:"form[fv-2]"`,
			postReqSetup, []string{ "23", "56" },
		},
		{ // Vialbe slice with type conversion
			rt.STypeOfUint16, `mvc:"form[fv-2]"`,
			postReqSetup, []uint16{ 23, 56 },
		},
		{ // Default value
			rt.STypeOfUint32, `mvc:"form[ng55] default[28,176,89]"`,
			postReqSetup, []uint32{ 28, 176, 89 },
		},
		// :~)

		/**
		 * Headers
		 */
		{ // Viable value
			rt.TypeOfString, `mvc:"header[hvv1]"`,
			getReqSetup, "hello24",
		},
		{ // Default value
			rt.TypeOfString, `mvc:"header[hvv3] default[gd-103]"`,
			getReqSetup, "gd-103",
		},
		{ // Default value of go lang
			rt.TypeOfInt16, `mvc:"query[hvv722]"`,
			getReqSetup, int16(0),
		},
		{ // Viable value with type conversion
			rt.TypeOfInt32, `mvc:"header[hvv2]"`,
			getReqSetup, int32(278),
		},
		{ // Checked param(true)
			rt.TypeOfBool, `mvc:"header[?hvv1]"`,
			getReqSetup, true,
		},
		{ // Checked param(false)
			rt.TypeOfBool, `mvc:"header[?hvv3]"`,
			getReqSetup, false,
		},
		// :~)

		/**
		 * Cookies
		 */
		{ // Viable value
			rt.TypeOfString, `mvc:"cookie[ck_1]"`,
			getReqSetup, "hello98",
		},
		{ // Default value
			rt.TypeOfString, `mvc:"cookie[ck_3] default[gd-33]"`,
			getReqSetup, "gd-33",
		},
		{ // Default value of go lang
			rt.TypeOfString, `mvc:"query[ck_12]"`,
			getReqSetup, "",
		},
		{ // Viable value with type conversion
			rt.TypeOfUint64, `mvc:"cookie[ck_2]"`,
			getReqSetup, uint64(77),
		},
		{ // Checked param(true)
			rt.TypeOfBool, `mvc:"cookie[?ck_1]"`,
			getReqSetup, true,
		},
		{ // Checked param(false)
			rt.TypeOfBool, `mvc:"cookie[?ck_3]"`,
			getReqSetup, false,
		},

		/**
		 * key/value in gin.Context
		 */
		{ // Viable value
			rt.TypeOfString, `mvc:"key[kg-1]"`,
			getReqSetup, "hello45",
		},
		{ // Default value
			rt.TypeOfString, `mvc:"key[kg-5] default[sampleKey1]"`,
			getReqSetup, "sampleKey1",
		},
		{ // Default value of go lang
			rt.TypeOfString, `mvc:"key[kg-23]"`,
			getReqSetup, "",
		},
		{ // Viable value with type conversion
			rt.TypeOfInt32, `mvc:"key[kg-2]"`,
			getReqSetup, int32(871),
		},
		{ // Checked param(true)
			rt.TypeOfBool, `mvc:"key[?kg-1]"`,
			getReqSetup, true,
		},
		{ // Checked param(false)
			rt.TypeOfBool, `mvc:"key[?kg-3]"`,
			getReqSetup, false,
		},
		{ // Checked param(false)
			rt.TypeOfBool, `mvc:"key[?kg-5]"`,
			getReqSetup, false,
		},
		// :~)
		{ // Not a MVC field
			rt.TypeOfBool, ``,
			getReqSetup, nil,
		},
	}

	convSrv := ot.NewDefaultConversionService()
	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		fieldType := reflect.StructField {
			Name: "FieldValue1",
			Type: testCase.sampleType,
			Tag: testCase.sampleTag,
		}

		paramLoader := buildParamLoader(fieldType, convSrv)
		if paramLoader == nil {
			continue
		}

		context := &gin.Context {}
		testCase.setupFunc(context)

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
