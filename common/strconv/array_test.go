package strconv

import (
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Split a string to INT array",
	func(source string, expected []int64) {
		testedResult := SplitStringToIntArray(source, "#")
		Expect(testedResult).To(Equal(expected))
	},
	Entry("Empty string to empty array of int64", "", []int64{}),
	Entry("Viable to array of int64", "123#445#-987#-229", []int64{123, 445, -987, -229}),
)

var _ = DescribeTable("Split a string to UINT array",
	func(source string, expected []uint64) {
		testedResult := SplitStringToUintArray(source, "#")
		Expect(testedResult).To(Equal(expected))
	},
	Entry("Empty string to empty array of int64", "", []uint64{}),
	Entry("Empty string to empty array of uint64", "87#14#44", []uint64{87, 14, 44}),
)
