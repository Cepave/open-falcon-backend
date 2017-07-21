package model

import (
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestPagingSuite struct{}

var _ = Suite(&TestPagingSuite{})

// Tests the setting of total count for paging
func (suite *TestPagingSuite) TestSetTotalCount(c *C) {
	testCases := []*struct {
		totalCount int32
		expected   bool
	}{
		{21, true},
		{20, false},
		{19, false},
		{0, false},
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedPaging := NewUndefinedPaging()
		testedPaging.Position = 2
		testedPaging.Size = 10

		testedPaging.SetTotalCount(testCase.totalCount)

		c.Assert(testedPaging.PageMore, Equals, testCase.expected, comment)
	}
}

// Tests the extracting of a page of data
func (suite *TestPagingSuite) TestExtractPage(c *C) {
	testCases := []*struct {
		size         int32
		position     int32
		expectedSize int
	}{
		{10, 1, 10},
		{10, 3, 10},
		{10, 4, 0},
		{10, 100, 0},
		{11, 1, 11},
		{11, 2, 11},
		{11, 3, 8},
		{11, 100, 0},
	}

	sampleArray := make([]int, 30)
	for i := 0; i < len(sampleArray); i++ {
		sampleArray[i] = i + 1
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		paging := &Paging{
			Position: testCase.position,
			Size:     testCase.size,
		}

		testedSlice := ExtractPage(sampleArray, paging).([]int)
		c.Logf("Got slice: %v", testedSlice)
		c.Assert(testedSlice, HasLen, testCase.expectedSize, comment)
	}
}
