package graph

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("buildBeforeTime() by multiple and exclusive parameters", func() {
	sampleTime, _ := time.Parse(time.RFC3339, "2015-10-15T10:20:20+08:00")

	DescribeTable("Result time should be as expected",
		func(parameter *relativeTimeParams, expectedTimeString string) {
			expectedTime, err := time.Parse(time.RFC3339, expectedTimeString)
			Expect(err).To(Succeed())

			Expect(buildBeforeTime(sampleTime, parameter)).To(Equal(expectedTime))
		},
		Entry("Default value of \"ForDays\"", &relativeTimeParams{-1, 0}, "2015-10-01T10:20:20+08:00"),
		Entry("Assign value of \"ForDays\"", &relativeTimeParams{2, 0}, "2015-10-13T10:20:20+08:00"),
		Entry("Assign value of both of \"ForDays\" and \"ForMinutes\"", &relativeTimeParams{4, 17}, "2015-10-11T10:20:20+08:00"),
		Entry("Assign value of \"ForMinutes\"", &relativeTimeParams{-1, 11}, "2015-10-15T10:09:20+08:00"),
	)
})
