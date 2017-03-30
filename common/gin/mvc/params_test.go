package mvc

import (
	"fmt"
	"net/http"
	"net/url"

	"gopkg.in/gin-gonic/gin.v1"

	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	otest "github.com/Cepave/open-falcon-backend/common/testing"
	. "gopkg.in/check.v1"
)

type TestParamSuite struct{}

var _ = Suite(&TestParamSuite{})

// Tests the getting value in *gin.Context(by key)
func (suite *TestParamSuite) TestContextKeyGetter(c *C) {
	testCases := []*struct {
		sampleValue interface{}
		defaultValue interface{}
		expectedValue interface{}
	} {
		{ 20, 30, 20 },
		{ -1, 30, 30 }, // No setting of value on key
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		sampleContext := &gin.Context {}

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
	paramName string
	paramValue []string
	defaultValue []string
	expectedValue []string
}
func (suite *TestParamSuite) TestGettersForArrayValue(c *C) {
	testCases := []*struct {
		name string
		contextSetup func(*gin.Context, *testCaseOfGetterOnArrayValue)
	} {
		{
			"query",
			func(context *gin.Context, testCase *testCaseOfGetterOnArrayValue) {
				context.Request = &http.Request{}
				context.Request.URL = otest.ParseRequestUri(
					c, "/query?" + buildUriValues(testCase.paramName, testCase.paramValue),
				)
			},
		},
		{
			"form",
			func(context *gin.Context, testCase *testCaseOfGetterOnArrayValue) {
				context.Request = &http.Request{}
				context.Request.PostForm, _ = url.ParseQuery(
					buildUriValues(testCase.paramName, testCase.paramValue),
				)
			},
		},
		{
			"cookie",
			func(context *gin.Context, testCase *testCaseOfGetterOnArrayValue) {
				context.Request = &http.Request{ Header: make(http.Header), }

				if len(testCase.paramValue) > 0 {
					context.Request.Header.Set("Cookie", fmt.Sprintf("%s=%s", testCase.paramName, testCase.paramValue[0]))
					testCase.expectedValue = []string{ testCase.expectedValue[0] }
				}
			},
		},
		{
			"param",
			func(context *gin.Context, testCase *testCaseOfGetterOnArrayValue) {
				context.Params = make([]gin.Param, 1)

				if len(testCase.paramValue) > 0 {
					context.Params[0] = gin.Param { Key: testCase.paramName, Value: testCase.paramValue[0] }
					testCase.expectedValue = []string{ testCase.expectedValue[0] }
				}
			},
		},
		{
			"header",
			func(context *gin.Context, testCase *testCaseOfGetterOnArrayValue) {
				context.Request = &http.Request { Header: make(http.Header), }
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

func buildUriValues(name string, values []string) string {
	result := make(url.Values)

	for _, v := range values {
		result.Add(name, v)
	}

	return result.Encode()
}

// Test the getting value of URI parameters
func gettersArrayValueTestingStub(c *C, testedGetter paramGetter, contextSetup func(*gin.Context, *testCaseOfGetterOnArrayValue)) {
	testCases := []*testCaseOfGetterOnArrayValue {
		{ "v2", []string{}, []string{ "74", "78" }, []string{ "74", "78" } },
		{ "v2", []string{}, []string{}, []string{} },
		{ "v2", []string{ "23", "66" }, []string{ "77" }, []string{ "23", "66" } },
	}

	sampleContext := &gin.Context {}

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
	paramName string
	paramValue string
	defaultValue string
	expectedValue string
}
func (suite *TestParamSuite) TestGettersForSingleValue(c *C) {
	testCases := []*struct {
		name string
		contextSetup func(*gin.Context, *testCaseOfGetterOnSingleValue)
	} {
		{
			"query",
			func(context *gin.Context, testCase *testCaseOfGetterOnSingleValue) {
				context.Request = &http.Request{}
				context.Request.URL = otest.ParseRequestUri(
					c, "/query?" + fmt.Sprintf("%s=%s", testCase.paramName, testCase.paramValue),
				)
			},
		},
		{
			"form",
			func(context *gin.Context, testCase *testCaseOfGetterOnSingleValue) {
				context.Request = &http.Request{}
				context.Request.PostForm, _ = url.ParseQuery(
					fmt.Sprintf("%s=%s", testCase.paramName, testCase.paramValue),
				)
			},
		},
		{
			"cookie",
			func(context *gin.Context, testCase *testCaseOfGetterOnSingleValue) {
				context.Request = &http.Request{ Header: make(http.Header), }
				context.Request.Header.Add("Cookie", fmt.Sprintf("%s=%s", testCase.paramName, testCase.paramValue))
			},
		},
		{
			"param",
			func(context *gin.Context, testCase *testCaseOfGetterOnSingleValue) {
				context.Params = make([]gin.Param, 1)
				context.Params[0] = gin.Param { Key: testCase.paramName, Value: testCase.paramValue }
			},
		},
		{
			"header",
			func(context *gin.Context, testCase *testCaseOfGetterOnSingleValue) {
				context.Request = &http.Request { Header: make(http.Header), }
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
	testCases := []*testCaseOfGetterOnSingleValue {
		{ "v1", "", "19", "19" },
		{ "v1", "", "", "" },
		{ "v1", "23", "77", "23" },
	}

	sampleContext := &gin.Context {}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		c.Logf("\tParam Data: %v", testCase)

		contextSetup(sampleContext, testCase)

		testedValue := testedGetter.getParam(sampleContext, testCase.paramName, testCase.defaultValue)
		c.Assert(testedValue, Equals, testCase.expectedValue, comment)
	}
}
