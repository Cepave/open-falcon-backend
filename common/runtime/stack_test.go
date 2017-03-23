package runtime

import (
	"fmt"
	//ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestStackSuite struct{}

var _ = Suite(&TestStackSuite{})

func ExampleGetCallerInfo() {
	f := func() *CallerInfo {
		return GetCallerInfo()
	}

	// This is line 18 of stack_test.go
	callerInfo := f()
	fmt.Printf("%s", callerInfo)

	// Output:
	// github.com/Cepave/open-falcon-backend/common/runtime/stack_test.go:18
}

func ExampleGetCallerInfoWithDepth() {
	f2 := func() *CallerInfo {
		return GetCallerInfoWithDepth(1)
	}
	f1 := func() *CallerInfo {
		return f2()
	}

	// This is line 35 of stack_test.go
	callerInfo := f1()
	fmt.Printf("%s", callerInfo)

	// Output:
	// github.com/Cepave/open-falcon-backend/common/runtime/stack_test.go:35
}

// Tests the information of caller
func (suite *TestStackSuite) TestGetCallerInfo(c *C) {
	testedInfo := info()

	c.Logf("Caller info: %#v", testedInfo)
	c.Assert(testedInfo.file, Matches, ".*stack_test.go")
}

func info() *CallerInfo {
	return GetCallerInfo()
}
