package utils

import (
	. "gopkg.in/check.v1"
	"sort"
)

type TestSetSuite struct{}

var _ = Suite(&TestSetSuite{})

// Tests the unique processing for a array of types which are valid as a key of map
func (suite *TestSetSuite) TestUniqueElements(c *C) {
	testCases := []*struct {
		source         interface{}
		expectedResult interface{}
	}{
		{
			[]string{"Z1", "Z2", "Z1", "Z2"},
			[]string{"Z1", "Z2"},
		},
		{
			[]int{33, 33, 67, 67, 56, 33},
			[]int{33, 67, 56},
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := UniqueElements(testCase.source)

		c.Assert(testedResult, DeepEquals, testCase.expectedResult, comment)
	}
}

// Tests the unique prcessing for a array of strings
func (suite *TestSetSuite) TestUniqueArrayOfString(c *C) {
	testCases := []*struct {
		sampleStrings  []string
		expectedResult []string
	}{
		{[]string{"A", "B", "A", "B"}, []string{"A", "B"}},
		{[]string{"C1", "C2", "C3", "C4"}, []string{"C1", "C2", "C3", "C4"}},
		{[]string{"C1", "C1", "C3", "C3", "C1", "C2", "C3", "C2"}, []string{"C1", "C3", "C2"}},
		{[]string{"G1", "", "G1", "G2", "", "G2"}, []string{"G1", "", "G2"}},
		{[]string{}, []string{}},
		{nil, nil},
	}

	for i, testCase := range testCases {
		testedResult := UniqueArrayOfStrings(testCase.sampleStrings)

		sort.Strings(testedResult)
		sort.Strings(testCase.expectedResult)

		c.Assert(testedResult, DeepEquals, testCase.expectedResult, Commentf("Test Case: %d", i))
	}
}
