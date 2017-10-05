package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("AreArrayOfStringsSame(interface{}, interface{})", func() {
	DescribeTable("result as expected one",
		func(leftArray []string, rightArray []string, expectedResult bool) {
			testedResult := AreArrayOfStringsSame(leftArray, rightArray)
			Expect(testedResult).To(Equal(expectedResult))
		},
		Entry("Same array", []string{"A", "B"}, []string{"A", "B"}, true),
		Entry("Differ sequence", []string{"A", "B"}, []string{"B", "A"}, true),
		Entry("Empty arrays", []string{}, []string{}, true),
		Entry("nil arrays", nil, nil, true),
		Entry("Empty array and nil array", []string{}, nil, true),
		Entry("Different arrays", []string{"A", "B"}, []string{"C", "B"}, false),
		Entry("Viable array and empty array", []string{"A", "B"}, []string{}, false),
		Entry("Viable array and nil array", []string{"A", "B"}, nil, false),
	)
})
