package check

import (
	. "gopkg.in/check.v1"
)

func LogTestCase(c *C, testCase interface{}) {
	c.Logf("Test Case: [%v]", testCase)
}

func TestCaseComment(index int) CommentInterface {
	return Commentf("Test Case: %d", index+1)
}
