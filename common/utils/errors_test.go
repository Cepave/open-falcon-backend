package utils

import (
	. "gopkg.in/check.v1"
)

type TestErrorSuite struct{}

var _ = Suite(&TestErrorSuite{})

// Tests the capture of error object
func (suite *TestErrorSuite) TestBuildPanicToError(c *C) {
	testCases := []*struct {
		needPanic bool
		errorChecker Checker
	} {
		{ true, NotNil },
		{ false, IsNil },
	}

	for i, testCase := range testCases {
		comment := Commentf("[%d] Test Case: %v", i + 1, testCase)
		c.Logf("%s", comment.CheckCommentString())

		var err error

		needPanic := testCase.needPanic
		testedFunc := BuildPanicToError(
			func() {
				samplePanic(needPanic)
			},
			&err,
		)
		testedFunc()

		c.Assert(err, testCase.errorChecker, comment)
	}
}

func samplePanic(needPanic bool) {
	if needPanic {
		panic("We are panic!!")
	}
}
