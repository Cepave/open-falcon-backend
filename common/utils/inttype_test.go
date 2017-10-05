package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("64 bits integer to integer of lesser bits", func() {
	Context("Int64 to interget fo lesser bits", func() {
		It("IntTo8()", func() {
			Expect(IntTo8([]int64{10, -44})).To(Equal([]int8{10, -44}))
			Expect(IntTo8([]int64(nil))).To(Equal([]int8(nil)))
			Expect(IntTo8([]int64{})).To(Equal([]int8{}))
		})
		It("IntTo16()", func() {
			Expect(IntTo16([]int64{234, -918})).To(Equal([]int16{234, -918}))
			Expect(IntTo16([]int64(nil))).To(Equal([]int16(nil)))
			Expect(IntTo16([]int64{})).To(Equal([]int16{}))
		})
		It("IntTo32()", func() {
			Expect(IntTo32([]int64{12234, -76918})).To(Equal([]int32{12234, -76918}))
			Expect(IntTo32([]int64(nil))).To(Equal([]int32(nil)))
			Expect(IntTo32([]int64{})).To(Equal([]int32{}))
		})
	})
	Context("Uint64 to interget fo lesser bits", func() {
		It("UintTo8()", func() {
			Expect(UintTo8([]uint64{10, 44})).To(Equal([]uint8{10, 44}))
			Expect(UintTo8([]uint64(nil))).To(Equal([]uint8(nil)))
			Expect(UintTo8([]uint64{})).To(Equal([]uint8{}))
		})
		It("UintTo16()", func() {
			Expect(UintTo16([]uint64{234, 918})).To(Equal([]uint16{234, 918}))
			Expect(UintTo16([]uint64(nil))).To(Equal([]uint16(nil)))
			Expect(UintTo16([]uint64{})).To(Equal([]uint16{}))
		})
		It("UintTo32()", func() {
			Expect(UintTo32([]uint64{12234, 76918})).To(Equal([]uint32{12234, 76918}))
			Expect(UintTo32([]uint64(nil))).To(Equal([]uint32(nil)))
			Expect(UintTo32([]uint64{})).To(Equal([]uint32{}))
		})
	})
})

var _ = Describe("Sort and unique array of number", func() {
	Context("Array of Int64", func() {
		DescribeTable("result array as expected",
			func(source []int64, expectedResult []int64) {
				testedResult := SortAndUniqueInt64(source)
				Expect(testedResult).To(Equal(expectedResult))
			},
			Entry("normal array",
				[]int64{30, 30, -10, 10, -7, 22},
				[]int64{-10, -7, 10, 22, 30},
			),
			Entry("nil array", nil, nil),
			Entry("empty array", []int64{}, []int64{}),
		)
	})

	Context("Array of Uint64", func() {
		DescribeTable("result array as expected",
			func(source []uint64, expectedResult []uint64) {
				testedResult := SortAndUniqueUint64(source)
				Expect(testedResult).To(Equal(expectedResult))
			},
			Entry("normal array",
				[]uint64{30, 30, 10, 10, 7, 22},
				[]uint64{7, 10, 22, 30},
			),
			Entry("nil array", nil, nil),
			Entry("empty array", []uint64{}, []uint64{}),
		)
	})
})
