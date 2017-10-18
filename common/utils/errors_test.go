package utils

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type TestErrorsSuite struct{}

var _ = Describe("Capture panic with &err object", func() {
	sampleFunc := func() (err error) {
		defer PanicToSimpleError(&err)()
		panic("Sample Error 1")
	}

	It("Error object should not be nil", func() {
		err := sampleFunc()
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("Put content of &err while panic", func() {
	DescribeTable("result error as expected",
		func(needPanic bool) {
			var err error

			testedFunc := BuildPanicToError(
				func() {
					if needPanic {
						panic("We are panic!!")
					}
				},
				&err,
			)
			testedFunc()

			if needPanic {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).To(Succeed())
			}
		},
		Entry("Got panic", true),
		Entry("Nothing panic", false),
	)
})

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
