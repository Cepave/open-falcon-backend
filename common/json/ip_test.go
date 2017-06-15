package json

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Test type IP", func() {
	DescribeTable("From unmarshalling JSON fromat to making sql driver Value", func(input []byte, expected int) {
		var ip IP
		err := ip.UnmarshalJSON(input)
		Expect(err).To(BeNil())
		actual, _ := ip.Value()
		Expect(actual).To(HaveLen(expected))
	},
		Entry("null value", []byte("null"), 0),
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

	DescribeTable("Marshal JSON format data", func(input string, expected string, expectedErr types.GomegaMatcher) {
		ip := NewIP(input)
		v, err := ip.MarshalJSON()
		Expect(string(v)).To(Equal(expected))
		Expect(err).To(expectedErr)
	},
		Entry("IPv4", "0.0.0.0", `"0.0.0.0"`, BeNil()),
		Entry("IPv4", "10.20.30.40", `"10.20.30.40"`, BeNil()),
		Entry("IPv6", "2001:cdba:0000:0000:0000:0000:3257:9652", `"2001:cdba::3257:9652"`, BeNil()),
		Entry("IPv6", "5cd8:91c3:ed6f:1dc4:8661:a4d:a9ae:d05c", `"5cd8:91c3:ed6f:1dc4:8661:a4d:a9ae:d05c"`, BeNil()),
		Entry("IPv6", "fd32:214f:fbca:97f4::", `"fd32:214f:fbca:97f4::"`, BeNil()),
		Entry("illegal string", "illegal", "null", BeNil()),
		Entry("incomplete IPv4", "10.20.30", "null", BeNil()),
		Entry("incomplete IPv4.", "10.20.30.", "null", BeNil()),
		Entry("incomplete IPv6", "2001:cdba:0000:0000:0000:0000:3257", "null", BeNil()),
		Entry("incomplete IPv6:", "2001:cdba:0000:0000:0000:0000:3257:", "null", BeNil()),
	)
})
