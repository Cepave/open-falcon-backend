package utils

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Abstract Array", func() {
	Context("Filter by customized function", func() {
		DescribeTable("result as expected one",
			func(testedFunc FilterFunc, sampleData interface{}, expectedData interface{}) {
				testedResult := MakeAbstractArray(sampleData).
					FilterWith(testedFunc)

				Expect(testedResult.GetArray()).To(Equal(expectedData))
			},
			Entry("middle element is filtered",
				TypedFuncToFilter(func(v string) bool { return v == "ok" }),
				[]string{"ok", "skip", "ok"}, []string{"ok", "ok"},
			),
			Entry("heading elements are filtered",
				TypedFuncToFilter(func(v int32) bool { return v > 20 }),
				[]int16{7, 18, 22, 98}, []int16{22, 98},
			),
		)
	})

	Context("Map by customized function", func() {
		DescribeTable("result as expected one",
			func(testedFunc MapperFunc, targetType reflect.Type, sampleData interface{}, expectedData interface{}) {
				testedResult := MakeAbstractArray(sampleData).
					MapTo(testedFunc, targetType)
				Expect(testedResult.GetArray()).To(Equal(expectedData))
			},
			Entry("value + 3",
				TypedFuncToMapper(func(v int8) int8 { return v + 3 }),
				reflect.TypeOf(int8(0)),
				[]int8{1, 3, 5}, []int8{4, 6, 8},
			),
			Entry("int64 to int8(+2)", // Tests the type conversion
				TypedFuncToMapper(func(v int64) int8 { return int8(v + 2) }),
				reflect.TypeOf(int8(0)),
				[]int64{11, 12}, []int8{13, 14},
			),
			Entry("string appeended with \"ok\"",
				TypedFuncToMapper(func(v string) string { return v + "ok" }),
				reflect.TypeOf(""),
				[]string{"g1:", "g2:"}, []string{"g1:ok", "g2:ok"},
			),
		)
	})

	Context("Unique filter", func() {
		DescribeTable("result as unique one",
			func(targetType reflect.Type, sampleData interface{}, expectedData interface{}) {
				testedResult := MakeAbstractArray(sampleData).
					FilterWith(NewUniqueFilter(targetType))
				Expect(testedResult.GetArray()).To(Equal(expectedData))
			},
			Entry("Remove one duplicated value",
				TypeOfString,
				[]string{"A1", "B1", "A2", "B1"},
				[]string{"A1", "B1", "A2"},
			),
			Entry("Remove two duplicated values",
				TypeOfInt16,
				[]int16{10, 20, 10, 20},
				[]int16{10, 20},
			),
		)
	})

	Context("Filter by NewDomainFilter()", func() {
		DescribeTable("filtered result as expected one",
			func(domain interface{}, sampleData interface{}, expectedData interface{}) {
				testedResult := MakeAbstractArray(sampleData).
					FilterWith(NewDomainFilter(domain))
				Expect(testedResult.GetArray()).To(Equal(expectedData))
			},
			Entry("last elements are removed",
				map[int]bool{1: true, 2: false},
				[]int{1, 2, 3, 4},
				[]int{1, 2},
			),
			Entry("middle element is removed",
				map[string]bool{"G1": true, "G2": false},
				[]string{"G1", "G3", "G2"},
				[]string{"G1", "G2"},
			),
		)
	})

	Context("GetArrayAsTargetType()", func() {
		DescribeTable("result as expected one",
			func(sourceArray, targetValue, expectedResult interface{}) {
				testedResult := MakeAbstractArray(sourceArray).
					GetArrayAsTargetType(targetValue)

				Expect(testedResult).To(Equal(expectedResult))
			},
			Entry("[]int16 to []uint32", []int16{11, 16}, uint32(0), []uint32{11, 16}),
			Entry("[]int64 to []int8", []int64{-1, -11}, int8(0), []int8{-1, -11}),
			Entry("[]int32 to []int32", []int32{-13, 109}, int32(0), []int32{-13, 109}),
			Entry("(empty array) []int32 to []int8", []int32{}, int8(0), []int8{}),
			Entry("(nil array) []int32 to []int16", []int32(nil), int16(0), []int16{}),
		)
	})

	Context("BatchProcess()", func() {
		sampleData := []int{10, 11, 22, 33, 71, 32}

		DescribeTable("The result slice should be as same as original",
			func(batchSize int) {
				receivedData := make([]int, 0, len(sampleData))

				MakeAbstractArray(sampleData).
					BatchProcess(
						batchSize,
						func(batch interface{}) {
							receivedData = append(receivedData, batch.([]int)...)
						},
						func(rest interface{}) {
							receivedData = append(receivedData, rest.([]int)...)
						},
					)

				Expect(receivedData).To(Equal(sampleData))
			},
			Entry("Perfect batch(no rest data)", 3),
			Entry("Perfect batch(1 batch)", 6),
			Entry("batch with rest", 4),
			Entry("batch is larger than original size(everything is put into rest)", 7),
			Entry("batch size is 1", 1),
		)
	})
})

var _ = Describe("FlattenToSlice()", func() {
	DescribeTable("Result should be matched expected one",
		func(sampleSlice interface{}, convertFunc func(v interface{}) []interface{}, expectedResult interface{}) {
			testedResult := FlattenToSlice(sampleSlice, convertFunc)

			Expect(testedResult).To(Equal(expectedResult))
		},
		Entry("Normal flatten",
			[]int{20, 30, 40},
			func(v interface{}) []interface{} {
				return []interface{}{
					v, v.(int) + 2,
				}
			},
			[]interface{}{
				20, 22, 30, 32, 40, 42,
			},
		),
		Entry("Empty slice",
			[]int{},
			func(v interface{}) []interface{} {
				return []interface{}{v, v}
			},
			[]interface{}{},
		),
	)
})
