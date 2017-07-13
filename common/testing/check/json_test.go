package check

import (
	. "gopkg.in/check.v1"
)

type TestJsonSuite struct{}

var _ = Suite(&TestJsonSuite{})

// Tests the checking of JSON
func (suite *TestJsonSuite) Test(c *C) {
	testCases := []*struct {
		obtainedJson   interface{}
		expectedJson   interface{}
		expectedResult bool
	}{
		{`[1, 3, 5]`, `[1, 3, 5]`, true},
		{`"Easy"`, `"Easy"`, true},
		{`"Easy"`, `"Easy2"`, false},
		{`38`, `38`, true},
		{`38`, `39`, false},
		{`{ "a": 10, "b": 20 }`, `{ "b": 20, "a": 10 }`, true},
		{`[1, 3, 5]`, `[3, 1, 5]`, false},
		{`{ "a": 10, "b": 20 }`, `{ "b": 21, "a": 11 }`, false},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult, _ := JsonEquals.Check(
			[]interface{}{testCase.obtainedJson, testCase.expectedJson},
			[]string{"obtained", "expected"},
		)

		c.Assert(testedResult, Equals, testCase.expectedResult, comment)
	}
}
