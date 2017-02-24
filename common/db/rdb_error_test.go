package db

import (
	"errors"
	//ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestErrorSuite struct{}

var _ = Suite(&TestErrorSuite{})

// Tests the error content(stack processing) of panic
func (suite *TestErrorSuite) TestPanicIfError(c *C) {
	c.Assert(
		func() {
			PanicIfError(errors.New("Sample Error"))
		},
		PanicMatches,
		"(?s:.*TestPanicIfError.*rdb_error_test.go.*Sample Error.*)",
	)
}
