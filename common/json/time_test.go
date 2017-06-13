package json

import (
	"time"

	. "gopkg.in/check.v1"
)

type TestTimeSuite struct{}

var _ = Suite(&TestTimeSuite{})

// Tests the serialization of JSON time
func (suite *TestTimeSuite) TestMarshalJSON(c *C) {
	testCases := []*struct {
		sampleTime     time.Time
		expectedResult string
	}{
		{time.Unix(2342377189, 0), "2342377189"},
		{time.Time{}, "null"},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedJson, _ := JsonTime(testCase.sampleTime).MarshalJSON()
		testedJsonString := string(testedJson)

		c.Logf("Time: %v. JSON; %s", testCase.sampleTime, testedJsonString)
		c.Assert(testedJsonString, Equals, testCase.expectedResult, comment)
	}
}

// Tests the deserialization of JSON time
func (suite *TestTimeSuite) TestUnmarshalJSON(c *C) {
	testCases := []*struct {
		input    []byte
		expected time.Time
	}{
		{[]byte("2342377189"), time.Unix(2342377189, 0)},
		{[]byte("null"), time.Time{}},
		{[]byte("0"), time.Unix(0, 0)},
		{[]byte("-2342377189"), time.Unix(-2342377189, 0)},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		var testedJson JsonTime
		err := testedJson.UnmarshalJSON(testCase.input)
		c.Assert(err, IsNil)

		c.Logf("JSON; %s, Time: %v", testCase.input, testCase.expected)
		c.Assert(time.Time(testedJson), Equals, testCase.expected, comment)
	}
}
