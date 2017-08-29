package types

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VarBytes", func() {
	Context("Convert to base64 string", func() {
		It("Match string of base64", func() {
			testedBytes := VarBytes{
				0x31, 0xa4, 0xee, 0xf8, 0xc0, 0xd6, 0xa5, 0xd0,
				0x02, 0xcc, 0x0c, 0x41, 0x92, 0x1e, 0xc8, 0xa9,
			}

			Expect(testedBytes.ToBase64()).To(Equal("MaTu+MDWpdACzAxBkh7IqQ=="))
		})
	})

	Context("Converted from base64 string", func() {
		It("Match bytes of base64", func() {
			var sampleBytes VarBytes

			err := sampleBytes.FromBase64("ASeh2RfCBc2HTTeB5Zr77A==")

			Expect(err).To(Succeed())
			Expect(sampleBytes).To(Equal(
				VarBytes{
					0x01, 0x27, 0xa1, 0xd9, 0x17, 0xc2, 0x05, 0xcd,
					0x87, 0x4d, 0x37, 0x81, 0xe5, 0x9a, 0xfb, 0xec,
				},
			))
		})

		It("Error has occur", func() {
			var sampleBytes VarBytes

			err := sampleBytes.FromBase64("!!")
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Bytes16", func() {
	Context("Convert to base64 string", func() {
		It("Match string of base64", func() {
			testedBytes := Bytes16{
				0x31, 0xa4, 0xee, 0xf8, 0xc0, 0xd6, 0xa5, 0xd0,
				0x02, 0xcc, 0x0c, 0x41, 0x92, 0x1e, 0xc8, 0xa9,
			}

			Expect(testedBytes.ToBase64()).To(Equal("MaTu+MDWpdACzAxBkh7IqQ=="))
		})
	})

	Context("Converted from base64 string", func() {
		It("Match bytes of base64", func() {
			var sampleBytes Bytes16

			err := sampleBytes.FromBase64("ASeh2RfCBc2HTTeB5Zr77A==")

			Expect(err).To(Succeed())
			Expect(sampleBytes).To(Equal(
				Bytes16{
					0x01, 0x27, 0xa1, 0xd9, 0x17, 0xc2, 0x05, 0xcd,
					0x87, 0x4d, 0x37, 0x81, 0xe5, 0x9a, 0xfb, 0xec,
				},
			))
		})

		It("Error has occur", func() {
			var sampleBytes Bytes16

			err := sampleBytes.FromBase64("!!")
			Expect(err).To(HaveOccurred())
		})
	})
})
