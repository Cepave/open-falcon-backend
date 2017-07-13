package types

import (
	"reflect"

	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"

	. "github.com/Cepave/open-falcon-backend/common/reflect/types"
	. "gopkg.in/check.v1"
)

type TestStringConvertSuite struct{}

var _ = Suite(&TestStringConvertSuite{})

// Tests the conversion from string to float types
func (suite *TestStringConvertSuite) TestConvertStringToFloat(c *C) {
	testCases := []*struct {
		value          string
		targetType     reflect.Type
		expectedResult interface{}
	}{
		{"93.23", TypeOfFloat32, float32(93.23)},
		{"101.07", TypeOfFloat64, float64(101.07)},
		{"", TypeOfFloat64, float64(0)},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedValue := convertStringToFloat(testCase.value, testCase.targetType)
		c.Assert(testedValue, Equals, testCase.expectedResult, comment)
	}
}

// Tests the conversion from string to int types
func (suite *TestStringConvertSuite) TestConvertStringToInt(c *C) {
	testCases := []*struct {
		value          string
		targetType     reflect.Type
		expectedResult interface{}
	}{
		{"33", TypeOfInt, 33},
		{"-34", TypeOfInt8, int8(-34)},
		{"35", TypeOfInt16, int16(35)},
		{"-36", TypeOfInt32, int32(-36)},
		{"37", TypeOfInt64, int64(37)},
		{"", TypeOfInt64, int64(0)},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedValue := convertStringToInt(testCase.value, testCase.targetType)
		c.Assert(testedValue, Equals, testCase.expectedResult, comment)
	}
}

// Tests the conversion from string to uint types
func (suite *TestStringConvertSuite) TestConvertStringToUint(c *C) {
	testCases := []*struct {
		value          string
		targetType     reflect.Type
		expectedResult interface{}
	}{
		{"33", TypeOfUint, uint(33)},
		{"34", TypeOfUint8, uint8(34)},
		{"35", TypeOfUint16, uint16(35)},
		{"36", TypeOfUint32, uint32(36)},
		{"37", TypeOfUint64, uint64(37)},
		{"38", TypeOfByte, byte(38)},
		{"", TypeOfUint64, uint64(0)},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedValue := convertStringToUint(testCase.value, testCase.targetType)
		c.Assert(testedValue, Equals, testCase.expectedResult, comment)
	}
}

// Tests the conversion from string to uint types
func (suite *TestStringConvertSuite) TestConvertStringTobool(c *C) {
	testCases := []*struct {
		value          string
		targetType     reflect.Type
		expectedResult interface{}
	}{
		{"t", TypeOfBool, true},
		{"true", TypeOfBool, true},
		{"y", TypeOfBool, true},
		{"yes", TypeOfBool, true},
		{"1", TypeOfBool, true},
		{"100", TypeOfBool, true},
		{"no", TypeOfBool, false},
		{"false", TypeOfBool, false},
		{"0", TypeOfBool, false},
		{"", TypeOfBool, false},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedValue := convertStringToBool(testCase.value, testCase.targetType)
		c.Assert(testedValue, Equals, testCase.expectedResult, comment)
	}
}
