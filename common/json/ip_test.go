package json

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test type IP", func() {
	DescribeTable("From unmarshalling JSON fromat to making sql driver Value", func(input []byte, expected int) {
		var ip IP
		err := ip.UnmarshalJSON(input)
		Expect(err).To(BeNil())
		actual, _ := ip.Value()
		Expect(actual).To(HaveLen(expected))
	},
		Entry("IPv4", []byte(`"0.0.0.0"`), 4),
		Entry("IPv4", []byte(`"10.20.30.40"`), 4),
		Entry("IPv6", []byte(`"2001:cdba:0000:0000:0000:0000:3257:9652"`), 16),
	)

	DescribeTable("Returns error for illegal inputs", func(input []byte) {
		var ip IP
		err := ip.UnmarshalJSON(input)
		Expect(err).NotTo(BeNil())
	},
		Entry("illegal string", []byte(`"illegal"`)),
		Entry("incomplete IPv4", []byte(`"10.20.30"`)),
		Entry("incomplete IPv4.", []byte(`"10.20.30."`)),
		Entry("incomplete IPv6", []byte(`"2001:cdba:0000:0000:0000:0000:3257"`)),
		Entry("incomplete IPv6:", []byte(`"2001:cdba:0000:0000:0000:0000:3257:"`)),
	)
})
