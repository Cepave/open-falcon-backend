package rpc

import (
	"github.com/Cepave/open-falcon-backend/common/model"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"

	. "gopkg.in/check.v1"
)

type TestRpcNqmAgentSuite struct{}

var _ = Suite(&TestRpcNqmAgentSuite{})

/**
 * Tests the data validation for ping task
 */
func (suite *TestRpcNqmAgentSuite) TestValidatePingTask(c *C) {
	var testCases = []*struct {
		connectionId string
		hostname     string
		ipAddress    string
		checker      Checker
	}{
		{"120.49.58.19", "localhost.localdomain", "120.49.58.19", IsNil},
		{"", "localhost.localdomain", "120.49.58.19", NotNil},
		{"120.49.58.19", "", "120.49.58.19", NotNil},
		{"120.49.58.19", "host1.com.cn", "900.49.58.19", NotNil}, // IP address cannot be parsed
		{"120.49.58.19", "localhost.localdomain", "", NotNil},
	}

	for i, testCase := range testCases {
		ocheck.LogTestCase(c, testCase)
		comment := ocheck.TestCaseComment(i)

		err := validatePingTask(
			&model.NqmTaskRequest{
				ConnectionId: testCase.connectionId,
				Hostname:     testCase.hostname,
				IpAddress:    testCase.ipAddress,
			},
		)

		c.Assert(err, testCase.checker, comment)
	}
}
