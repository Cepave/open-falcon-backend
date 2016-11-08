package nqm

import (
	"github.com/Cepave/open-falcon-backend/common/model"
	. "gopkg.in/check.v1"
)

type TestNqmAgentSuite struct{}

var _ = Suite(&TestNqmAgentSuite{})

// Tests the length of bytes for IP address
type ipAddressTestCase struct {
	sampleIpAddress string
	expectedLength  int
}

func (suite *TestNqmAgentSuite) TestIpAddress(c *C) {
	var testedCases = []ipAddressTestCase{
		{"10.20.30.40", 4},
		{"2001:cdba:0000:0000:0000:0000:3257:9652", 16},
	}

	for _, v := range testedCases {
		testedAgent := NewNqmAgent(&model.NqmTaskRequest{
			IpAddress: v.sampleIpAddress,
		})

		c.Assert(len(testedAgent.IpAddress), Equals, v.expectedLength)
	}
}
