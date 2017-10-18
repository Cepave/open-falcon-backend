package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type sampleKey string

func (k sampleKey) GetKey() interface{} {
	return k
}

var _ = Describe("Grouping by NewGroupingProcessorOfTargetType()", func() {
	Context("Grouping operations", func() {
		testedProcessor := NewGroupingProcessorOfTargetType(int(0))

		testedProcessor.Put(sampleKey("GD-1"), 20)
		testedProcessor.Put(sampleKey("GD-1"), 30)
		testedProcessor.Put(sampleKey("GD-1"), 40)
		testedProcessor.Put(sampleKey("GD-2"), 70)
		testedProcessor.Put(sampleKey("GD-2"), 80)

		It("Number of keys should be 2", func() {
			Expect(testedProcessor.Keys()).To(HaveLen(2))
		})

		It("Key object should be as same as input", func() {
			Expect(testedProcessor.KeyObject(sampleKey("GD-1"))).To(Equal(sampleKey("GD-1")))
			Expect(testedProcessor.KeyObject(sampleKey("GD-2"))).To(Equal(sampleKey("GD-2")))
		})

		It("Number of children should be as same as input", func() {
			Expect(testedProcessor.Children(sampleKey("GD-1"))).To(HaveLen(3))
			Expect(testedProcessor.Children(sampleKey("GD-2"))).To(HaveLen(2))
		})

		It("Children should be array of <type>", func() {
			intValues := testedProcessor.Children(sampleKey("GD-1")).([]int)
			Expect(intValues).To(Equal([]int{20, 30, 40}))
		})
	})
})
