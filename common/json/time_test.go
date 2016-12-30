package json

import (
	. "gopkg.in/check.v1"
	"time"
)

type TestTimeSuite struct{}

var _ = Suite(&TestTimeSuite{})

// Tests the serialization of JSON time
func (suite *TestTimeSuite) TestMarshalJSON(c *C) {
	testCases := []*struct {
		sampleTime time.Time
		expectedResult string
	} {
		{ time.Unix(2342377189, 0), "2342377189" },
		{ time.Time{}, "null" },
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i + 1)

		testedJson, _ := JsonTime(testCase.sampleTime).MarshalJSON()
		testedJsonString := string(testedJson)

		c.Logf("Time: %v. JSON; %s", testCase.sampleTime, testedJsonString)
		c.Assert(testedJsonString, Equals, testCase.expectedResult, comment)
	}
}
