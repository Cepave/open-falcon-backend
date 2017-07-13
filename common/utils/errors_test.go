package utils

import (
	"fmt"
	. "gopkg.in/check.v1"
)

type TestErrorsSuite struct{}

var _ = Suite(&TestErrorsSuite{})

// Tests the capture of panic to error object
func (suite *TestErrorsSuite) TestPanicToSimpleError(c *C) {
	sampleFunc := func() (err error) {
		defer PanicToSimpleError(&err)()

		panic("Sample Error 1")
	}

	c.Assert(sampleFunc(), NotNil)
}

func ExamplePanicToError() {
	sampleFunc := func() (err error) {
		defer PanicToError(
			&err,
			func(p interface{}) error {
				return fmt.Errorf("Customized: %v", p)
			},
		)()

		panic("Good Error!!")
	}

	err := sampleFunc()
	fmt.Println(err)

	// Output:
	// Customized: Good Error!!
}

func ExamplePanicToSimpleError() {
	sampleFunc := func() (err error) {
		defer PanicToSimpleError(&err)()

		panic("Novel Error!!")
	}

	err := sampleFunc()
	fmt.Println(err)

	// Output:
	// Novel Error!!
}

func ExamplePanicToSimpleErrorWrapper() {
	sampleFunc := func(n int) {
		panic(fmt.Sprintf("Value: %d", n))
	}

	testedFunc := PanicToSimpleErrorWrapper(
		func() { sampleFunc(918) },
	)

	fmt.Println(testedFunc())
}

// Tests the capture of error object
func (suite *TestErrorsSuite) TestBuildPanicToError(c *C) {
	testCases := []*struct {
		needPanic    bool
		errorChecker Checker
	}{
		{true, NotNil},
		{false, IsNil},
	}

	for i, testCase := range testCases {
		comment := Commentf("[%d] Test Case: %v", i+1, testCase)
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
