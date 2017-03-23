package nqm

import (
	"github.com/leebenson/conform"

	. "gopkg.in/check.v1"
)

type TestPingtaskSuite struct{}

var _ = Suite(&TestPingtaskSuite{})

// Tests validation of pingtasks
func (suite *TestPingtaskSuite) TestPingtaskModify(c *C) {
	testCase := &PingtaskModify{
		Name:    " 台灣 ",
		Comment: " 測試用 ",
	}

	conform.Strings(testCase)

	c.Assert(testCase.Name, Equals, "台灣")
	c.Assert(testCase.Comment, Equals, "測試用")
}
