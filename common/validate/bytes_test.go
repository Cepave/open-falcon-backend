package validate

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("TagNonZeroSlice(non_zero_slice)", func() {
	validator := newValidator()

	Context("Slice or array of various integer type", func() {
		DescribeTable("Check the validation result",
			func(sampleVariable interface{}, passed bool) {
				err := validator.Var(sampleVariable, TagNonZeroSlice)

				if passed {
					Expect(err).To(Succeed())
				} else {
					Expect(err).To(HaveOccurred())
				}
			},
			Entry("int8 slice(zero)", []int8{0, 0}, false),
			Entry("uint32 array[3](zero)", [3]uint32{0, 0, 0}, false),
			Entry("int8 slice", []int8{1, 2}, true),
			Entry("int16 empty slice", []int16{}, true),
			Entry("int16 nil slice", []int16(nil), true),
			Entry("int64 array", [2]int64{29, 33}, true),
			Entry("uint32 slice", []uint32{9, 0}, true),
		)
	})
})
