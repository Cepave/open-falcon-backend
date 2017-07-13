package mvc

import (
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestParamSuite struct{}

var _ = Suite(&TestParamSuite{})

// Tests the getting value in *gin.Context(by key)
func (suite *TestParamSuite) TestContextKeyGetter(c *C) {
	testCases := []*struct {
		sampleValue   interface{}
		defaultValue  interface{}
		expectedValue interface{}
	}{
		{20, 30, 20},
		{-1, 30, 30}, // No setting of value on key
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		sampleContext := &gin.Context{}

		if testCase.sampleValue.(int) != -1 {
			sampleContext.Set("data-1", testCase.sampleValue)
		}

		testedValue := keyGetter.getValue(
			sampleContext, "data-1", testCase.defaultValue,
		)

		c.Assert(testedValue, Equals, testCase.expectedValue, comment)
	}
}

// Test the various getters:
//
// 1. Cookie Getter
// 2. URI Param Getter
// 3. Header Getter
type testCaseOfGetterOnArrayValue struct {
	paramName     string
	paramValue    []string
	defaultValue  []string
	expectedValue []string
}

func (suite *TestParamSuite) TestGettersForArrayValue(c *C) {
	testCases := []*struct {
		name         string
		contextSetup func(*gin.Context, *testCaseOfGetterOnArrayValue)
	}{
		{
			"query",
			func(context *gin.Context, testCase *testCaseOfGetterOnArrayValue) {
				context.Request = httptest.NewRequest(
					"GET", "/query?"+buildUriValuesAsString(testCase.paramName, testCase.paramValue), nil,
				)
			},
		},
		{
			"form",
			func(context *gin.Context, testCase *testCaseOfGetterOnArrayValue) {
				context.Request = httptest.NewRequest(
					"POST", "/form-post", buildUriValuesAsBody(testCase.paramName, testCase.paramValue),
				)
				context.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			},
		},
		{
			"cookie",
			func(context *gin.Context, testCase *testCaseOfGetterOnArrayValue) {
				context.Request = httptest.NewRequest("GET", "/cookie-getter", nil)

				if len(testCase.paramValue) > 0 {
					context.Request.Header.Set("Cookie", fmt.Sprintf("%s=%s", testCase.paramName, testCase.paramValue[0]))
					testCase.expectedValue = []string{testCase.expectedValue[0]}
				}
			},
		},
		{
			"param",
			func(context *gin.Context, testCase *testCaseOfGetterOnArrayValue) {
				context.Params = make([]gin.Param, 1)

				if len(testCase.paramValue) > 0 {
					context.Params[0] = gin.Param{Key: testCase.paramName, Value: testCase.paramValue[0]}
					testCase.expectedValue = []string{testCase.expectedValue[0]}
				}
			},
		},
		{
			"header",
			func(context *gin.Context, testCase *testCaseOfGetterOnArrayValue) {
				context.Request = httptest.NewRequest("GET", "/header-getter", nil)
				for _, v := range testCase.paramValue {
					context.Request.Header.Add(testCase.paramName, v)
				}
			},
		},
	}

	for _, testCase := range testCases {
		ocheck.LogTestCase(c, testCase)
		gettersArrayValueTestingStub(c, paramGetters[testCase.name], testCase.contextSetup)
	}
}

func buildUriValuesAsBody(name string, values []string) io.Reader {
	return strings.NewReader(buildUriValuesAsString(name, values))
}
func buildUriValuesAsString(name string, values []string) string {
	return buildUriValues(name, values).Encode()
}
func buildUriValues(name string, values []string) url.Values {
	result := make(url.Values)

	for _, v := range values {
		result.Add(name, v)
	}

	return result
}

// Test the getting value of URI parameters
func gettersArrayValueTestingStub(c *C, testedGetter paramGetter, contextSetup func(*gin.Context, *testCaseOfGetterOnArrayValue)) {
	testCases := []*testCaseOfGetterOnArrayValue{
		{"v2", []string{}, []string{"74", "78"}, []string{"74", "78"}},
		{"v2", []string{}, []string{}, []string{}},
		{"v2", []string{"23", "66"}, []string{"77"}, []string{"23", "66"}},
	}

	sampleContext := &gin.Context{}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		c.Logf("\tParam Data: %v", testCase)

		contextSetup(sampleContext, testCase)

		testedValue := testedGetter.getParamAsArray(sampleContext, testCase.paramName, testCase.defaultValue)
		c.Assert(testedValue, DeepEquals, testCase.expectedValue, comment)
	}
}

// Test the various getters:
//
// 1. Cookie Getter
// 2. URI Param Getter
// 3. Header Getter
type testCaseOfGetterOnSingleValue struct {
	paramName     string
	paramValue    string
	defaultValue  string
	expectedValue string
}

func (suite *TestParamSuite) TestGettersForSingleValue(c *C) {
	testCases := []*struct {
		name         string
		contextSetup func(*gin.Context, *testCaseOfGetterOnSingleValue)
	}{
		{
			"query",
			func(context *gin.Context, testCase *testCaseOfGetterOnSingleValue) {
				context.Request = httptest.NewRequest(
					"GET", "/query?"+buildUriValuesAsString(testCase.paramName, []string{testCase.paramValue}), nil,
				)
			},
		},
		{
			"form",
			func(context *gin.Context, testCase *testCaseOfGetterOnSingleValue) {
				context.Request = httptest.NewRequest("POST", "/post-form", buildUriValuesAsBody(testCase.paramName, []string{testCase.paramValue}))
				context.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			},
		},
		{
			"cookie",
			func(context *gin.Context, testCase *testCaseOfGetterOnSingleValue) {
				context.Request = httptest.NewRequest("GET", "/cookie-getter", nil)
				context.Request.Header.Add("Cookie", fmt.Sprintf("%s=%s", testCase.paramName, testCase.paramValue))
			},
		},
		{
			"param",
			func(context *gin.Context, testCase *testCaseOfGetterOnSingleValue) {
				context.Params = make([]gin.Param, 1)
				context.Params[0] = gin.Param{Key: testCase.paramName, Value: testCase.paramValue}
			},
		},
		{
			"header",
			func(context *gin.Context, testCase *testCaseOfGetterOnSingleValue) {
				context.Request = httptest.NewRequest("GET", "/header-getter", nil)
				context.Request.Header.Set(testCase.paramName, testCase.paramValue)
			},
		},
	}

	for _, testCase := range testCases {
		ocheck.LogTestCase(c, testCase)
		gettersSingleValueTestingStub(c, paramGetters[testCase.name], testCase.contextSetup)
	}
}

// Test the getting value of URI parameters
func gettersSingleValueTestingStub(c *C, testedGetter paramGetter, contextSetup func(*gin.Context, *testCaseOfGetterOnSingleValue)) {
	testCases := []*testCaseOfGetterOnSingleValue{
		{"v1", "", "19", "19"},
		{"v1", "", "", ""},
		{"v1", "23", "77", "23"},
	}

	sampleContext := &gin.Context{}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		c.Logf("\tParam Data: %v", testCase)

		contextSetup(sampleContext, testCase)

		testedValue := testedGetter.getParam(sampleContext, testCase.paramName, testCase.defaultValue)
		c.Assert(testedValue, Equals, testCase.expectedValue, comment)
	}
}

// Tests the bool getter for checking of viable value
func (suite *TestParamSuite) TestBoolParamGetter(c *C) {
	queryContextSetup := func(context *gin.Context) {
		context.Request = httptest.NewRequest("GET", "/query?g1=33&g2=", nil)

		context.Request.Header.Add("Cookie", "ck1=28")
		context.Request.Header.Add("Cookie", "ck2=")
		context.Request.Header.Add("hd1", "ss90")
		context.Request.Header.Add("hd2", "  ")

		context.Set("key-1", "v1")
		context.Set("key-2", nil)
		context.Set("key-es", "  ")
	}
	formContextSetup := func(context *gin.Context) {
		values := make(url.Values)
		values.Add("fc1", "10")
		values.Add("fc2", "  ")

		context.Request = httptest.NewRequest("POST", "/form-post", strings.NewReader(values.Encode()))
		context.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	testCases := []*struct {
		testedFunc     string
		sampleParam    string
		contextSetup   func(*gin.Context)
		expectedResult bool
	}{
		{"query", "g1", queryContextSetup, true},
		{"query", "g2", queryContextSetup, false}, // Empty value
		{"query", "g3", queryContextSetup, false}, // No such param
		{"form", "fc1", formContextSetup, true},
		{"form", "fc2", formContextSetup, false}, // Empty value
		{"form", "fc3", formContextSetup, false}, // No such param
		{"cookie", "ck1", queryContextSetup, true},
		{"cookie", "ck2", queryContextSetup, false}, // Empty value
		{"cookie", "ck3", queryContextSetup, false}, // No such param
		{"header", "hd1", queryContextSetup, true},
		{"header", "hd2", queryContextSetup, false}, // Empty value
		{"header", "hd3", queryContextSetup, false}, // No such param
		{"key", "key-1", queryContextSetup, true},
		{"key", "key-es", queryContextSetup, false}, // Empty string
		{"key", "key-2", queryContextSetup, false},  // nil value
		{"key", "key-3", queryContextSetup, false},  // No such param
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		sampleContext := &gin.Context{}

		testCase.contextSetup(sampleContext)

		testedResult := paramCheckers[testCase.testedFunc](sampleContext, testCase.sampleParam)

		c.Assert(testedResult, Equals, testCase.expectedResult, comment)
	}
}
