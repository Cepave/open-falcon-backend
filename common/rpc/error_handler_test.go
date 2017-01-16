package rpc

import (
	"errors"
	ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestErrorSuite struct{}

var _ = Suite(&TestErrorSuite{})

// Tests the capture of error object
func (suite *TestErrorSuite) TestHandleError(c *C) {
	testCases := []*struct {
		sampleError interface{}
		expectedError string
	} {
		{ "P1", "P1" },
		{ errors.New("E1"), "E1" },
	}

	for i, testCase := range testCases {
		comment := ocheck.TestCaseComment(i)
		ocheck.LogTestCase(c, testCase)

		testedError := samplePanic(testCase.sampleError)

		c.Assert(testedError.Error(), Matches, ".*" + testCase.expectedError + ".*", comment)
	}
}

func samplePanic(samplePanic interface{}) (err error) {
	defer HandleError(&err)()
	panic(samplePanic)
}
