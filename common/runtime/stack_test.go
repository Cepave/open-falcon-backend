package runtime

import (
	"fmt"
	"regexp"

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
	replaceRegExp, _ := regexp.Compile(`github\.com/.+$`)

	fmt.Printf("%s", replaceRegExp.FindStringSubmatch(callerInfo.String())[0])

	// Output:
	// github.com/Cepave/open-falcon-backend/common/runtime/stack_test.go:20
}

func ExampleGetCallerInfoWithDepth() {
	f2 := func() *CallerInfo {
		return GetCallerInfoWithDepth(1)
	}
	f1 := func() *CallerInfo {
		return f2()
	}

	callerInfo := f1()
	fmt.Printf("%s", callerInfo)
}

// Tests the information of caller
func (suite *TestStackSuite) TestGetCallerInfo(c *C) {
	testedInfo := info()

	c.Logf("Caller info: %#v", testedInfo)
	c.Assert(testedInfo.rawFile, Matches, ".*stack_test.go")
}

// Tests the getting of caller informations
func (suite *TestStackSuite) TestGetCallerInfoStack(c *C) {
	stack := GetCallerInfoStack(0, 10)

	for _, callerInfo := range stack {
		c.Logf("%#v", callerInfo.GetFile())
	}
}

func info() *CallerInfo {
	return GetCallerInfo()
}
