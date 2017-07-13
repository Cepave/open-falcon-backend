package utils

import (
	. "gopkg.in/check.v1"
)

type TestCompareSuite struct{}

var _ = Suite(&TestCompareSuite{})

// Tests the compare of two arrays
func (suite *TestCompareSuite) TestAreArrayOfStringsSame(c *C) {
	testCases := []*struct {
		leftArray      []string
		rightArray     []string
		expectedResult bool
	}{
		{[]string{"A", "B"}, []string{"A", "B"}, true},
		{[]string{"A", "B"}, []string{"B", "A"}, true},
		{[]string{}, []string{}, true},
		{nil, nil, true},
		{[]string{}, nil, true},
		{[]string{"A", "B"}, []string{"C", "B"}, false},
		{[]string{"A", "B"}, []string{}, false},
		{[]string{"A", "B"}, nil, false},
	}

	for i, testCase := range testCases {
		testedResult := AreArrayOfStringsSame(testCase.leftArray, testCase.rightArray)
		c.Assert(testedResult, Equals, testCase.expectedResult, Commentf("Test Case: %d", i+1))
	}
}
