package utils

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type car struct{}
type func1 func()

type g1Err struct{}

func (e *g1Err) Error() string {
	return "OK"
}

var _ = Describe("Is viable", func() {
	ch1, ch2 := make(chan bool, 1), make(chan bool, 1)
	ch1 <- true

	var nilErr1 error = (*g1Err)(nil)
	var nilErr2 *g1Err = (*g1Err)(nil)

	DescribeTable("result as expected one",
		func(sampleValue interface{}, expected bool) {
			testedValue := reflect.ValueOf(sampleValue)
			testedResult := ValueExt(testedValue).IsViable()

			Expect(testedResult).To(Equal(expected))
		},
		Entry("30 is viable", 30, true),
		Entry("0 is viable", 0, true),
		Entry("Initialized pointer to *struct is viable", &car{}, true),
		Entry("Nil pointer to *struct is not viable", (*car)(nil), false),
		Entry("Non-empty array is viable", []int{20}, true),
		Entry("Empty array is not viable", []int{}, false),
		Entry("Nil array is not viable", []string(nil), false),
		Entry("Non-empty array(element's type is pointer) is viable", []*car{{}}, true),
		Entry("Nil array(element's type is pointer) is not viable", []*car{}, false),
		Entry("Non-empty map is viable", map[int]bool{20: true}, true),
		Entry("Empty map is not viable", map[int]bool{}, false),
		Entry("Function is viable", func1(func() {}), true),
		Entry("Nil function is not viable", func1(nil), false),
		Entry("Non-empty channel is viable", ch1, true),
		Entry("Empty channel is not viable", ch2, false),
		Entry("Nil error(pure) is not viable", (error)(nil), false),
		Entry("Nil error is not viable", nilErr1, false),
		Entry("Nil error(alias) is not viable", nilErr2, false),
	)
})

var _ = Describe("Generic type converions", func() {
	Context("Converts value to pointer one", func() {
		type weight int
		var v1 int = 20

		It("Conversion for viable value", func() {
			convertedValue := ConvertToTargetType(&v1, new(weight)).(*weight)
			Expect(*convertedValue).To(Equal(weight(20)))
		})

		It("Conversion for nil value", func() {
			Expect(ConvertToTargetType((*string)(nil), new(weight))).To(Equal((*weight)(nil)))
		})
	})

	Context("Conversions for real numbers", func() {
		DescribeTable("result as expected one",
			func(sourceValue interface{}, targetValue interface{}) {
				testedResult := ConvertToTargetType(sourceValue, targetValue)
				Expect(testedResult).To(Equal(targetValue))
			},
			Entry("int8 to int16", int8(10), int16(10)),
			Entry("int8 to int32", int8(11), int32(11)),
			Entry("int8 to int64", int8(12), int64(12)),
			Entry("uint8 to uint16", uint8(13), uint16(13)),
			Entry("uint8 to uint32", uint8(14), uint32(14)),
			Entry("uint8 to uint64", uint8(15), uint64(15)),
			Entry("int8 to uint16", int8(16), uint16(16)),
			Entry("uint64 to int8", uint64(33), int8(33)),
			Entry("float64 to float32", float64(31.77), float32(31.77)),
		)
	})

	Context("Conversions for alias of struct type", func() {
		type car struct {
			age  int
			name string
		}

		type carV1 car

		c1 := car{20, "LBUE-98"}

		It("result struct should be as same as expected one", func() {
			testedCar := ConvertToTargetType(c1, carV1{}).(carV1)
			Expect(testedCar.age).To(Equal(c1.age))
			Expect(testedCar.name).To(Equal(c1.name))
		})
	})
})
