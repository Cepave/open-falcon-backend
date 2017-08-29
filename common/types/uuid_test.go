package types

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Uuid", func() {
	Context("From string", func() {
		It("Match expected UUID", func() {
			sampleString := "fef97a42-f1d1-4dff-bfe3-83506a901db3"

			var testedUuid Uuid
			testedUuid.MustFromString(sampleString)

			Expect(testedUuid).To(Equal(Uuid{
				0xfe, 0xf9, 0x7a, 0x42, 0xf1, 0xd1, 0x4d, 0xff,
				0xbf, 0xe3, 0x83, 0x50, 0x6a, 0x90, 0x1d, 0xb3,
			}))
		})

		It("Panic while parsing string has error", func() {
			var testedUuid Uuid

			Expect(func() {
				testedUuid.MustFromString("zzz!!")
			}).To(Panic())
		})
	})

	Context("From bytes", func() {
		It("Match expected UUID", func() {
			sampleBytes := []byte{
				0x1e, 0x8b, 0x99, 0x6f, 0x77, 0x24, 0x4f, 0x04,
				0x98, 0x09, 0x91, 0x33, 0xae, 0x0b, 0x9f, 0x7d,
			}

			var testedUuid Uuid
			testedUuid.MustFromBytes(sampleBytes)

			Expect(testedUuid).To(Equal(Uuid{
				0x1e, 0x8b, 0x99, 0x6f, 0x77, 0x24, 0x4f, 0x04,
				0x98, 0x09, 0x91, 0x33, 0xae, 0x0b, 0x9f, 0x7d,
			}))
		})

		It("Panic while reading from []byte has error", func() {
			var testedUuid Uuid

			Expect(func() {
				testedUuid.MustFromBytes([]byte{0x12, 0x44, 0x91})
			}).To(Panic())
		})
	})
})
