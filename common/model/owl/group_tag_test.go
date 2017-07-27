package owl

import (
	. "gopkg.in/check.v1"
)

type TestGroupTagSuite struct{}

var _ = Suite(&TestGroupTagSuite{})

// Tests the splitting for group tags
func (suite *TestGroupTagSuite) TestSplitToArryOfGroupTags(c *C) {
	testCases := []*struct {
		sampleIds      string
		sampleNames    string
		expectedResult []*GroupTag
	}{
		{"12,34,81", "GT-1#GT-2#GT-3",
			[]*GroupTag{
				{
					12, "GT-1",
				},
				{
					34, "GT-2",
				},
				{
					81, "GT-3",
				},
			},
		},
		{"", "", []*GroupTag{}},
	}

	for i, testCase := range testCases {
		testedResult := SplitToArrayOfGroupTags(
			testCase.sampleIds, ",",
			testCase.sampleNames, "#",
		)

		c.Assert(testedResult, DeepEquals, testCase.expectedResult, Commentf("Test Case: %d", i+1))
	}
}
