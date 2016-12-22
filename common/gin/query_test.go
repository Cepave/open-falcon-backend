package gin

import (
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
	. "gopkg.in/check.v1"
)

type TestGinQuerySuite struct{}

var _ = Suite(&TestGinQuerySuite{})

// Tests the getting of int64 with default value
func (suite *TestGinQuerySuite) TestGetInt64Default(c *C) {
	testCases := []struct {
		queryValue string
		expectedValue int64
		expectedViable bool
		hasError bool
	} {
		{ "v1=2120", 2120, true, false },
		{ "v1=-33867", -33867, true, false },
		{ "", 50, false, false },
		{ "v1=abc", -1, false, true },
	}

	context := &gin.Context{}
	testedWrapper := NewQueryWrapper(context)
	for _, testCase := range testCases {
		request, _ := http.NewRequest("GET", "http://127.0.0.1/?" + testCase.queryValue, nil)
		context.Request = request

		testedParamValue := testedWrapper.GetInt64Default("v1", 50)
		c.Assert(testedParamValue.Viable, Equals, testCase.expectedViable)

		if testCase.hasError {
			c.Assert(testedParamValue.Error, NotNil)
		} else {
			c.Assert(testedParamValue.Error, IsNil)
			c.Assert(testedParamValue.Value, Equals, testCase.expectedValue)
		}
	}
}

// Tests the getting of uint64 with default value
func (suite *TestGinQuerySuite) TestGetUint64Default(c *C) {
	testCases := []struct {
		queryValue string
		expectedValue uint64
		expectedViable bool
		hasError bool
	} {
		{ "v1=33", 33, true, false },
		{ "", 50, false, false },
		{ "v1=-9871", 0, false, true },
	}

	context := &gin.Context{}
	testedWrapper := NewQueryWrapper(context)
	for _, testCase := range testCases {
		request, _ := http.NewRequest("GET", "http://127.0.0.1/?" + testCase.queryValue, nil)
		context.Request = request

		testedParamValue := testedWrapper.GetUint64Default("v1", 50)
		c.Assert(testedParamValue.Viable, Equals, testCase.expectedViable)

		if testCase.hasError {
			c.Assert(testedParamValue.Error, NotNil)
		} else {
			c.Assert(testedParamValue.Error, IsNil)
			c.Assert(testedParamValue.Value, Equals, testCase.expectedValue)
		}
	}
}

// Tests the getting of uint64 with default value
func (suite *TestGinQuerySuite) TestGetBoolDefault(c *C) {
	testCases := []struct {
		queryValue string
		expectedValue bool
		expectedViable bool
		hasError bool
	} {
		{ "v1=true", true, true, false }, // True value
		{ "v1=0", false, true, false }, // False value
		{ "", true, false, false }, // No viable
		{ "v1=abc", true, false, true }, // Cannot be parsed
	}

	context := &gin.Context{}
	testedWrapper := NewQueryWrapper(context)
	for _, testCase := range testCases {
		request, _ := http.NewRequest("GET", "http://127.0.0.1/?" + testCase.queryValue, nil)
		context.Request = request

		testedParamValue := testedWrapper.GetBoolDefault("v1", true)
		c.Assert(testedParamValue.Viable, Equals, testCase.expectedViable)

		if testCase.hasError {
			c.Assert(testedParamValue.Error, NotNil)
		} else {
			c.Assert(testedParamValue.Error, IsNil)
			c.Assert(testedParamValue.Value, Equals, testCase.expectedValue)
		}
	}
}

// Tests the getting of uint64 with default value
func (suite *TestGinQuerySuite) TestGetFloat64Default(c *C) {
	testCases := []struct {
		queryValue string
		expectedValue float64
		expectedViable bool
		hasError bool
	} {
		{ "v1=34.56", 34.56, true, false }, // True value
		{ "v1=-77.819", -77.819, true, false }, // False value
		{ "", 123.45, false, false }, // No viable
		{ "v1=abc", 0, false, true }, // Cannot be parsed
	}

	context := &gin.Context{}
	testedWrapper := NewQueryWrapper(context)
	for _, testCase := range testCases {
		request, _ := http.NewRequest("GET", "http://127.0.0.1/?" + testCase.queryValue, nil)
		context.Request = request

		testedParamValue := testedWrapper.GetFloat64Default("v1", 123.45)
		c.Assert(testedParamValue.Viable, Equals, testCase.expectedViable)

		if testCase.hasError {
			c.Assert(testedParamValue.Error, NotNil)
		} else {
			c.Assert(testedParamValue.Error, IsNil)
			c.Assert(testedParamValue.Value, Equals, testCase.expectedValue)
		}
	}
}
