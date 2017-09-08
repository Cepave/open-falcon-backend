package json

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("For Bytes16", func() {
	Context("JSON marshalling", func() {
		DescribeTable("Result BASE64 must match expected one",
			func(srcBytes Bytes16, expectedBase64 string) {
				testedResult, err := srcBytes.MarshalJSON()

				Expect(err).To(Succeed())
				Expect(string(testedResult)).To(Equal(expectedBase64))
			},
			Entry("16 bytes data", [16]byte{
				0x90, 0x39, 0xd0, 0x44, 0xbb, 0xc9, 0x21, 0x46, 0xf6, 0x95, 0x47, 0x66, 0x2f, 0x24, 0x6b, 0x23,
			}, "\"kDnQRLvJIUb2lUdmLyRrIw==\""),
			Entry("empty byte", [16]byte{}, "\"AAAAAAAAAAAAAAAAAAAAAA==\""),
		)
	})

	Context("JSON Unmarshalling", func() {
		DescribeTable("Result byte array[16] must match expected one",
			func(srcBase64 string, expected Bytes16) {
				var testedBytes Bytes16
				err := testedBytes.UnmarshalJSON([]byte(srcBase64))

				Expect(err).To(Succeed())
				Expect(testedBytes).To(Equal(expected))
			},
			Entry("16 bytes data",
				"\"kDnQRLvJIUb2lUdmLyRrIw==\"",
				[16]byte{
					0x90, 0x39, 0xd0, 0x44, 0xbb, 0xc9, 0x21, 0x46, 0xf6, 0x95, 0x47, 0x66, 0x2f, 0x24, 0x6b, 0x23,
				},
			),
			Entry("empty string", "\"\"", [16]byte{}),
			Entry("null value", "null", [16]byte{}),
		)
	})
})

var _ = Describe("For VarBytes", func() {
	Context("JSON marshalling", func() {
		DescribeTable("Result BASE64 must match expected one",
			func(srcBytes VarBytes, expectedBase64 string) {
				testedResult, err := srcBytes.MarshalJSON()

				Expect(err).To(Succeed())
				Expect(string(testedResult)).To(Equal(expectedBase64))
			},
			Entry("Viable bytes data", []byte{
				0x90, 0x39, 0xd0, 0x44, 0xbb, 0xc9, 0x21, 0x46, 0xf6, 0x95, 0x47, 0x66, 0x2f, 0x24, 0x6b, 0x23,
			}, "\"kDnQRLvJIUb2lUdmLyRrIw==\""),
			Entry("empty var bytes", []byte{}, "null"),
			Entry("nil var bytes", []byte(nil), "null"),
		)
	})

	Context("JSON Unmarshalling", func() {
		DescribeTable("Result byte slice must match expected one",
			func(srcBase64 string, expected VarBytes) {
				var testedBytes VarBytes
				err := testedBytes.UnmarshalJSON([]byte(srcBase64))

				Expect(err).To(Succeed())
				Expect(testedBytes).To(Equal(expected))
			},
			Entry("Viable bytes data",
				"\"kDnQRLvJIUb2lUdmLyRrIw==\"",
				[]byte{
					0x90, 0x39, 0xd0, 0x44, 0xbb, 0xc9, 0x21, 0x46, 0xf6, 0x95, 0x47, 0x66, 0x2f, 0x24, 0x6b, 0x23,
				},
			),
			Entry("empty string", "\"\"", []byte(nil)),
			Entry("null value", "null", []byte(nil)),
		)
	})
})
