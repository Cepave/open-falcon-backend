package model

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestHeartbeatSuite struct{}

var _ = Suite(&TestHeartbeatSuite{})

func (suite *TestHeartbeatSuite) TestIPStringOfNqmAgentHeartbeatRequest(c *C) {
	var testedCases = []*struct {
		input    NqmAgentHeartbeatRequest
		expected int
	}{
		{input: NqmAgentHeartbeatRequest{IpAddress: "0.0.0.0"}, expected: 4},
		{input: NqmAgentHeartbeatRequest{IpAddress: "10.20.30.40"}, expected: 4},
		{input: NqmAgentHeartbeatRequest{IpAddress: "2001:cdba:0000:0000:0000:0000:3257:9652"}, expected: 16},
	}

	for i, v := range testedCases {
		expectedValue, err := v.input.IpAddress.Value()
		c.Assert(err, IsNil)
		c.Assert(len(expectedValue.([]byte)), Equals, v.expected, Commentf("Test Case: %d", i+1))
	}
}
