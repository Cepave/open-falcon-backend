package utils

import (
	. "gopkg.in/check.v1"
	"sort"
)

type TestSetSuite struct{}

var _ = Suite(&TestSetSuite{})

// Tests the unique prcessing for a array of strings
func (suite *TestSetSuite) TestUniqueArrayOfString(c *C) {
	testCases := []struct {
		sampleStrings []string
		expectedResult []string
	} {
		{ []string{ "A", "B", "A", "B" }, []string{ "A", "B" } },
		{ []string{ "C1", "C2", "C3", "C4" }, []string{ "C1", "C2", "C3", "C4" } },
		{ []string{}, []string{} },
		{ nil, []string{} },
	}

	for i, testCase := range testCases {
		testedResult := UniqueArrayOfStrings(testCase.sampleStrings)

		sort.Strings(testedResult)
		sort.Strings(testCase.expectedResult)

		c.Assert(testedResult, DeepEquals, testCase.expectedResult, Commentf("Test Case: %d", i))
	}
}
