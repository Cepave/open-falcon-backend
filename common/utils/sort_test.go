package utils

import (
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Compare values with flag of checking nil", func() {
	sampleValue := 10

	DescribeTable("result as expected one",
		func(left *int, right *int, direction byte, expectedResult int, expectedHasNil bool) {
			testedResult, testedHasNil := CompareNil(left, right, direction)

			Expect(testedResult).To(Equal(expectedResult))
			Expect(testedHasNil).To(Equal(expectedHasNil))
		},
		Entry("Equivalence(not-nil)", &sampleValue, &sampleValue, DefaultDirection, SeqEqual, false),
		Entry("Equivalence(nil)", nil, nil, DefaultDirection, SeqEqual, true),
		Entry("Lower(right operand is nil)", &sampleValue, nil, DefaultDirection, SeqLower, true),
		Entry("Higher(left operand is nil)", nil, &sampleValue, DefaultDirection, SeqHigher, true),
		Entry("Equivalence(Descending, not-nil)", &sampleValue, &sampleValue, Descending, SeqEqual, false),
		Entry("Equivalence(Descending, nil)", nil, nil, Descending, SeqEqual, true),
		Entry("Higher(Descending, right operand is nil)", &sampleValue, nil, Descending, SeqHigher, true),
		Entry("lower(Descending, left operand is nil)", nil, &sampleValue, Descending, SeqLower, true),
	)
})

var _ = Describe("Compare typed values", func() {
	Context("String values", func() {
		DescribeTable("result as expected one",
			func(left string, right string, direction byte, expectedResult int) {
				testedResult := CompareString(left, right, direction)
				Expect(testedResult).To(Equal(expectedResult))
			},
			Entry("Equivalence", "A", "A", DefaultDirection, SeqEqual),
			Entry("Higher", "A", "B", DefaultDirection, SeqHigher),
			Entry("Lower", "B", "A", DefaultDirection, SeqLower),
			Entry("Descending(Equivalence)", "A", "A", Descending, SeqEqual),
			Entry("Descending(Lower)", "A", "B", Descending, SeqLower),
			Entry("Descending(Higher)", "B", "A", Descending, SeqHigher),
			Entry("Lower(For <UNDEFINED>)", "廣東", "<UNDEFINED>", DefaultDirection, SeqLower),
		)
	})

	Context("Interger values", func() {
		Context("Int64 values", func() {
			DescribeTable("result as expected one",
				func(left int, right int, direction byte, expectedResult int) {
					testedResult := CompareInt(int64(left), int64(right), direction)
					Expect(testedResult).To(Equal(expectedResult))
				},
				Entry("Equivalence", 10, 10, DefaultDirection, SeqEqual),
				Entry("Higher", 10, 20, DefaultDirection, SeqHigher),
				Entry("Lower", 20, 10, DefaultDirection, SeqLower),
				Entry("Descending(Equivalence)", 10, 10, Descending, SeqEqual),
				Entry("Descending(Lower)", 10, 20, Descending, SeqLower),
				Entry("Descending(Higher)", 20, 10, Descending, SeqHigher),
			)
		})

		Context("Uint64 values", func() {
			DescribeTable("result as expected one",
				func(left int, right int, direction byte, expectedResult int) {
					testedResult := CompareUint(uint64(left), uint64(right), direction)
					Expect(testedResult).To(Equal(expectedResult))
				},
				Entry("Equivalence", 10, 10, DefaultDirection, SeqEqual),
				Entry("Higher", 10, 20, DefaultDirection, SeqHigher),
				Entry("Lower", 20, 10, DefaultDirection, SeqLower),
				Entry("Descending(Equivalence)", 10, 10, Descending, SeqEqual),
				Entry("Descending(Lower)", 10, 20, Descending, SeqLower),
				Entry("Descending(Higher)", 20, 10, Descending, SeqHigher),
			)
		})
	})

	Context("Float values", func() {
		DescribeTable("result as expected one",
			func(left float64, right float64, direction byte, expectedResult int) {
				testedResult := CompareFloat(left, right, direction)
				Expect(testedResult).To(Equal(expectedResult))
			},
			Entry("Equivalence", 10.33, 10.33, DefaultDirection, SeqEqual),
			Entry("Higher", 10.19, 20.23, DefaultDirection, SeqHigher),
			Entry("Lower", 20.87, 10.23, DefaultDirection, SeqLower),
			Entry("Descending(Equivalence)", 10.05, 10.05, Descending, SeqEqual),
			Entry("Descending(Lower)", 10.11, 20.98, Descending, SeqLower),
			Entry("Descending(Higher)", 20.44, 10.02, Descending, SeqHigher),
		)
	})

	Context("Compare IP address", func() {
		DescribeTable("result as expected one",
			func(leftIp net.IP, rightIp net.IP, direction byte, expected int) {
				testedResult := CompareIpAddress(leftIp, rightIp, direction)
				Expect(testedResult).To(Equal(expected))
			},
			Entry("Higher(1st byte)", net.ParseIP("0.0.0.0"), net.ParseIP("20.0.0.0"), DefaultDirection, SeqHigher),
			Entry("Lower(1st byte)", net.ParseIP("20.0.0.0"), net.ParseIP("0.0.0.0"), DefaultDirection, SeqLower),
			Entry("Higher(1st byte)", net.ParseIP("40.20.30.40"), net.ParseIP("109.20.30.40"), DefaultDirection, SeqHigher),
			Entry("Lower(1st byte)", net.ParseIP("109.20.30.40"), net.ParseIP("40.20.30.40"), DefaultDirection, SeqLower),
			Entry("Equivalence(nil)", nil, nil, DefaultDirection, SeqEqual),
			Entry("Higher(left operand is nil)", nil, net.ParseIP("10.20.30.40"), DefaultDirection, SeqHigher),
			Entry("Higher(right operand is nil)", net.ParseIP("10.20.30.40"), nil, DefaultDirection, SeqLower),
		)
	})
})
