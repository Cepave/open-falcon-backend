package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Shorten String(by size)", func() {
	sampleText := "Hello World! 明月幾時有"
	sampleSize := len(sampleText)

	Context("Arguments are viable", func() {
		DescribeTable("Result as expected",
			func(maxSize int, expected string) {
				testedResult := ShortenStringToSize(sampleText, "...", maxSize)

				Expect(testedResult).To(Equal(expected))
			},
			Entry("Source is equvalent maximum", sampleSize, sampleText),
			Entry("Source is shorter than maximum", sampleSize+1, sampleText),
			Entry("Source is longer than maximum", 6, "Hel ... 幾時有"),
			Entry("Source is longer than maximum(odd size)", 7, "Hell ... 幾時有"),
		)
	})

	Context("Source string is empty", func() {
		It("Result is empty", func() {
			testedResult := ShortenStringToSize("", "|", 20)
			Expect(testedResult).To(Equal(""))
		})
	})

	Context("Maximuim size is 1", func() {
		It("Result is viable", func() {
			testedResult := ShortenStringToSize("東風不與周郎便，銅雀春深鎖二喬", "|", 1)
			Expect(testedResult).To(Equal("東 | 喬"))
		})
	})

	Context("Maximuim size is 0", func() {
		It("Should be panic", func() {
			Expect(func() { ShortenStringToSize("折戟沉沙鐵未銷，自將磨洗認前朝", "|", 0) }).To(Panic())
		})
	})

	Context("More is empty", func() {
		It("Heading and tailing is connected by a space", func() {
			testedResult := ShortenStringToSize("落魄江湖載酒行", "", 4)
			Expect(testedResult).To(Equal("落魄 酒行"))
		})
	})
})
