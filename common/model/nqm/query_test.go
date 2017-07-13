package nqm

import (
	. "gopkg.in/check.v1"
)

type TestQuerySuite struct{}

var _ = Suite(&TestQuerySuite{})

// Tests the checking if the filter has descriptive information of agent
func (suite *TestQuerySuite) TestHasAgentDescriptive(c *C) {
	testCases := []*struct {
		sampleFilter *AgentFilter
		expected     bool
	}{
		{
			&AgentFilter{}, false,
		},
		{
			&AgentFilter{
				Name: []string{"A"},
			},
			true,
		},
		{
			&AgentFilter{
				Hostname: []string{"A"},
			},
			true,
		},
		{
			&AgentFilter{
				IpAddress: []string{"A"},
			},
			true,
		},
		{
			&AgentFilter{
				ConnectionId: []string{"A"},
			},
			true,
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := testCase.sampleFilter.HasAgentDescriptive()
		c.Assert(testedResult, Equals, testCase.expected, comment)
	}
}

// Tests the checking if the filter has descriptive information of target
func (suite *TestQuerySuite) TestHasTargetDescriptive(c *C) {
	testCases := []*struct {
		sampleFilter *TargetFilter
		expected     bool
	}{
		{
			&TargetFilter{}, false,
		},
		{
			&TargetFilter{
				Name: []string{"A"},
			},
			true,
		},
		{
			&TargetFilter{
				Host: []string{"A"},
			},
			true,
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedResult := testCase.sampleFilter.HasTargetDescriptive()
		c.Assert(testedResult, Equals, testCase.expected, comment)
	}
}
