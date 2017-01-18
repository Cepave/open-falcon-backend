package strconv

import (
	. "gopkg.in/check.v1"
)

type TestArraySuite struct{}

var _ = Suite(&TestArraySuite{})

// Tests the spliting of a string to int array
func (suite *TestArraySuite) TestSplitStringToIntArray(c *C) {
	testCases := []*struct {
		values string
		expectedResult []int64
	} {
		{ "", []int64 {} },
		{ "123#445#-987#-229", []int64 { 123, 445, -987, -229 } },
	}

	for _, testCase := range testCases {
		testedResult := SplitStringToIntArray(testCase.values, "#")
		c.Assert(testedResult, DeepEquals, testCase.expectedResult)
	}
}
// Tests the spliting of a string to uint array
func (suite *TestArraySuite) TestSplitStringToUintArray(c *C) {
	testCases := []*struct {
		values string
		expectedResult []uint64
	} {
		{ "", []uint64 {} },
		{ "1,3,5,7,9", []uint64 { 1, 3, 5, 7, 9 } },
	}

	for _, testCase := range testCases {
		testedResult := SplitStringToUintArray(testCase.values, ",")
		c.Assert(testedResult, DeepEquals, testCase.expectedResult)
	}
}
