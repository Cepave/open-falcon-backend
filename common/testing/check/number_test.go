package check

import (
	. "gopkg.in/check.v1"
)

type TestNumberSuite struct{}

var _ = Suite(&TestNumberSuite{})

// Tests the larger than checker
func (suite *TestNumberSuite) TestLargerThan(c *C) {
	testCases := []*struct {
		sampleLeft, sampleRight int
		expectedResult          bool
	}{
		{4, 3, true},
		{3, 3, false},
		{2, 3, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult, _ := LargerThan.Check(
			[]interface{}{testCase.sampleLeft, testCase.sampleRight},
			[]string{},
		)

		c.Assert(testedResult, Equals, testCase.expectedResult, comment)
	}
}

// Tests the larger than or equal checker
func (suite *TestNumberSuite) TestLargerThanOrEqualTo(c *C) {
	testCases := []*struct {
		sampleLeft, sampleRight int
		expectedResult          bool
	}{
		{4, 3, true},
		{3, 3, true},
		{2, 3, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult, _ := LargerThanOrEqualTo.Check(
			[]interface{}{testCase.sampleLeft, testCase.sampleRight},
			[]string{},
		)

		c.Assert(testedResult, Equals, testCase.expectedResult, comment)
	}
}

// Tests the smaller than checker
func (suite *TestNumberSuite) TestSmallerThan(c *C) {
	testCases := []*struct {
		sampleLeft, sampleRight int
		expectedResult          bool
	}{
		{3, 4, true},
		{3, 3, false},
		{4, 3, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult, _ := SmallerThan.Check(
			[]interface{}{testCase.sampleLeft, testCase.sampleRight},
			[]string{},
		)

		c.Assert(testedResult, Equals, testCase.expectedResult, comment)
	}
}

// Tests the smaller than or equal checker
func (suite *TestNumberSuite) TestSmallerThanOrEqualTo(c *C) {
	testCases := []*struct {
		sampleLeft, sampleRight int
		expectedResult          bool
	}{
		{2, 3, true},
		{3, 3, true},
		{4, 3, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult, _ := SmallerThanOrEqualTo.Check(
			[]interface{}{testCase.sampleLeft, testCase.sampleRight},
			[]string{},
		)

		c.Assert(testedResult, Equals, testCase.expectedResult, comment)
	}
}

// Tests the compare for different types
func (suite *TestNumberSuite) TestPerformCompare(c *C) {
	testCases := []*struct {
		leftValue, rightValue interface{}
		expectedResult        bool
	}{
		{int(3), int8(44), false},
		{uint(20), uint8(7), true},
		{float32(20.509), float64(7.87), true},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		expectedResult := performCompare(
			testCase.leftValue, testCase.rightValue,
			largerThan,
		)
		c.Assert(expectedResult, Equals, testCase.expectedResult, comment)
	}
}
