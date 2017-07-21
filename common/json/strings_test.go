package json

import (
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestStringsSuite struct{}

var _ = Suite(&TestStringsSuite{})

// Tests the marshalling of JSON
func (suite *TestStringsSuite) TestMarshalJSONOfJsonString(c *C) {
	testCases := []*struct {
		sample   []JsonString
		expected string
	}{
		{
			[]JsonString{"A", "", "B"},
			`[ "A", null, "B" ]`,
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := MarshalJSON(testCase.sample)
		c.Logf("JSON Result: %s", testedResult)

		c.Assert(testedResult, ocheck.JsonEquals, testCase.expected, comment)
	}
}
