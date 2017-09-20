package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CartisianProduct", func() {
	Context("normal situation", func() {
		It("Number of records and detail of records(1st and 2nd)", func() {
			testedResult := CartesianProduct(
				[]interface{}{"a1", "a2"},
				[]interface{}{"b1"},
				[]interface{}{"c1", "c2"},
			)

			GinkgoT().Logf("Result: %#v", testedResult)

			Expect(testedResult).To(HaveLen(4))
			Expect(testedResult[0]).To(Equal([]interface{}{"a1", "b1", "c1"}))
			Expect(testedResult[1]).To(Equal([]interface{}{"a1", "b1", "c2"}))
		})
	})

	Context("Empty set", func() {
		Context("Everything is empty", func() {
			It("Zero record", func() {
				Expect(CartesianProduct()).To(HaveLen(0))
				Expect(CartesianProduct([]interface{}{})).To(HaveLen(0))
				Expect(CartesianProduct([]interface{}{}, []interface{}{})).To(HaveLen(0))
			})
		})
		Context("Some set are empty", func() {
			aSet := []interface{}{"a1", "a2"}
			bSet := []interface{}{"b1"}
			cSet := []interface{}{}

			It("Zero record", func() {
				Expect(CartesianProduct(aSet, bSet, cSet)).To(HaveLen(0))
				Expect(CartesianProduct(cSet, bSet)).To(HaveLen(0))
				Expect(CartesianProduct(aSet, cSet, bSet)).To(HaveLen(0))
			})
		})
	})
})
