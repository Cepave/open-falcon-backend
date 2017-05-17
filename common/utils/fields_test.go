package utils

import (
	. "gopkg.in/check.v1"
)

var _ = Suite(&TestSetSuite{})

func (suite *TestSetSuite) TestDictedFieldstring(c *C) {
	testCases := []struct {
		sampleStrings  string
		expectedResult map[string]interface{}
	}{
		{"a=Mike has a hard penis!, digit=5", map[string]interface{}{"a": "Mike has a hard penis!", "digit": "5"}},
		{"", map[string]interface{}{}},
	}

	for i, testCase := range testCases {
		testedResult := DictedFieldstring(testCase.sampleStrings)
		c.Assert(testedResult, DeepEquals, testCase.expectedResult, Commentf("Test Case: %d", i))
	}
}
