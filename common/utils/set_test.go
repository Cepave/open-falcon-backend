package utils

import (
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unique array of <type>", func() {
	Context("UniqueElements()", func() {
		DescribeTable("result as expected one",
			func(source interface{}, expectedResult interface{}) {
				testedResult := UniqueElements(source)
				Expect(testedResult).To(Equal(expectedResult))
			},
			Entry("array of string",
				[]string{"Z1", "Z2", "Z1", "Z2"},
				[]string{"Z1", "Z2"},
			),
			Entry("array of int",
				[]int{33, 33, 67, 67, 56, 33},
				[]int{33, 67, 56},
			),
		)
	})

	Context("UniqueArrayOfStrings()", func() {
		DescribeTable("result as expected one",
			func(sampleStrings []string, expectedResult []string) {
				testedResult := UniqueArrayOfStrings(sampleStrings)

				sort.Strings(testedResult)
				sort.Strings(expectedResult)

				Expect(testedResult).To(Equal(expectedResult))
			},
			Entry("Duplicated elements", []string{"A", "B", "A", "B"}, []string{"A", "B"}),
			Entry("Source is unique", []string{"C1", "C2", "C3", "C4"}, []string{"C1", "C2", "C3", "C4"}),
			Entry("Contains empty string", []string{"G1", "", "G1", "G2", "", "G2"}, []string{"G1", "", "G2"}),
			Entry("Empty array", []string{}, []string{}),
			Entry("Nil array", nil, nil),
		)
	})
})
