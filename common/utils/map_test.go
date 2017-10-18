package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Abstract Map", func() {
	Context("ToType() of different types for key and value", func() {
		type s2 string

		It("result map should be same as expected", func() {
			sampleAMap := MakeAbstractMap(map[int16]s2{
				1: "Nice",
				2: "Good",
			})

			testedMap := sampleAMap.ToTypeOfTarget(int32(0), "").(map[int32]string)
			Expect(testedMap).To(Equal(
				map[int32]string{
					1: "Nice",
					2: "Good",
				},
			))
		})
	})
})
