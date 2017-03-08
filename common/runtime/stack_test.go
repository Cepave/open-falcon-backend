package runtime

import (
	//ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestStackSuite struct{}

var _ = Suite(&TestStackSuite{})

// Tests the information of caller
func (suite *TestStackSuite) TestGetCallerInfo(c *C) {
	testedInfo := info()

	c.Logf("Caller info: %#v", testedInfo)
	c.Assert(testedInfo.File, Matches, ".*stack_test.go")
}

func info() *CallerInfo {
	return GetCallerInfo()
}
