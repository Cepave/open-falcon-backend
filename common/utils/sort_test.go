package utils

import (
	. "gopkg.in/check.v1"
	"net"
)

type TestSortSuite struct{}

var _ = Suite(&TestSortSuite{})

// Tests the comparison of nil values
func (suite *TestSortSuite) TestCompareNil(c *C) {
	sampleValue := 10

	testCases := []*struct {
		left           *int
		right          *int
		direction      byte
		expectedResult int
		expectedHasNil bool
	}{
		{&sampleValue, &sampleValue, DefaultDirection, SeqEqual, false},
		{nil, nil, DefaultDirection, SeqEqual, true},
		{&sampleValue, nil, DefaultDirection, SeqLower, true},
		{nil, &sampleValue, DefaultDirection, SeqHigher, true},
		{&sampleValue, &sampleValue, Descending, SeqEqual, false},
		{nil, nil, Descending, SeqEqual, true},
		{&sampleValue, nil, Descending, SeqHigher, true},
		{nil, &sampleValue, Descending, SeqLower, true},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)
		c.Logf("Test Case Data: %v", testCase)

		testedResult, testedHasNil := CompareNil(testCase.left, testCase.right, testCase.direction)

		c.Assert(testedResult, Equals, testCase.expectedResult, comment)
		c.Assert(testedHasNil, Equals, testCase.expectedHasNil, comment)
	}
}

// Tests the comparison of string values
func (suite *TestSortSuite) TestCompareString(c *C) {
	testCases := []*struct {
		left           string
		right          string
		direction      byte
		expectedResult int
	}{
		{"A", "A", DefaultDirection, SeqEqual},
		{"A", "B", DefaultDirection, SeqHigher},
		{"B", "A", DefaultDirection, SeqLower},
		{"A", "A", Descending, SeqEqual},
		{"A", "B", Descending, SeqLower},
		{"B", "A", Descending, SeqHigher},
		{"廣東", "<UNDEFINED>", DefaultDirection, SeqLower},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)
		c.Logf("Test Case Data: %v", testCase)

		testedResult := CompareString(testCase.left, testCase.right, testCase.direction)
		c.Assert(testedResult, Equals, testCase.expectedResult, comment)
	}
}

// Tests the comparison of int values
func (suite *TestSortSuite) TestCompareInt(c *C) {
	testCases := []*struct {
		left           int8
		right          int8
		direction      byte
		expectedResult int
	}{
		{10, 10, DefaultDirection, SeqEqual},
		{10, 20, DefaultDirection, SeqHigher},
		{20, 10, DefaultDirection, SeqLower},
		{10, 10, Descending, SeqEqual},
		{10, 20, Descending, SeqLower},
		{20, 10, Descending, SeqHigher},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)
		c.Logf("Test Case Data: %v", testCase)

		testedResult := CompareInt(int64(testCase.left), int64(testCase.right), testCase.direction)
		c.Assert(testedResult, Equals, testCase.expectedResult, comment)
	}
}

// Tests the comparison of uint values
func (suite *TestSortSuite) TestCompareUint(c *C) {
	testCases := []*struct {
		left           uint8
		right          uint8
		direction      byte
		expectedResult int
	}{
		{10, 10, DefaultDirection, SeqEqual},
		{10, 20, DefaultDirection, SeqHigher},
		{20, 10, DefaultDirection, SeqLower},
		{10, 10, Descending, SeqEqual},
		{10, 20, Descending, SeqLower},
		{20, 10, Descending, SeqHigher},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)
		c.Logf("Test Case Data: %v", testCase)

		testedResult := CompareUint(uint64(testCase.left), uint64(testCase.right), testCase.direction)
		c.Assert(testedResult, Equals, testCase.expectedResult, comment)
	}
}

// Tests the comparison of float values
func (suite *TestSortSuite) TestCompareFloat(c *C) {
	testCases := []*struct {
		left           float64
		right          float64
		direction      byte
		expectedResult int
	}{
		{10.33, 10.33, DefaultDirection, SeqEqual},
		{10.19, 10.76, DefaultDirection, SeqHigher},
		{10.76, 10.19, DefaultDirection, SeqLower},
		{10.33, 10.33, Descending, SeqEqual},
		{10.19, 10.76, Descending, SeqLower},
		{10.76, 10.19, Descending, SeqHigher},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)
		c.Logf("Test Case Data: %v", testCase)

		testedResult := CompareFloat(testCase.left, testCase.right, testCase.direction)
		c.Assert(testedResult, Equals, testCase.expectedResult, comment)
	}
}

// Tests the comparison for IP Address
func (suite *TestSortSuite) TestCompareIpAddress(c *C) {
	testCases := []*struct {
		leftIp    net.IP
		rightIp   net.IP
		direction byte
		expected  int
	}{
		{net.ParseIP("0.0.0.0"), net.ParseIP("20.0.0.0"), DefaultDirection, SeqHigher},
		{net.ParseIP("20.0.0.0"), net.ParseIP("0.0.0.0"), DefaultDirection, SeqLower},
		{net.ParseIP("40.20.30.40"), net.ParseIP("109.20.30.40"), DefaultDirection, SeqHigher},
		{net.ParseIP("109.20.30.40"), net.ParseIP("40.20.30.40"), DefaultDirection, SeqLower},
		{nil, nil, DefaultDirection, SeqEqual},
		{nil, net.ParseIP("10.20.30.40"), DefaultDirection, SeqHigher},
		{net.ParseIP("10.20.30.40"), nil, DefaultDirection, SeqLower},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)
		c.Logf("Test Case Data: %v", testCase)

		testedResult := CompareIpAddress(
			testCase.leftIp, testCase.rightIp,
			testCase.direction,
		)
		c.Assert(testedResult, Equals, testCase.expected, comment)
	}
}
