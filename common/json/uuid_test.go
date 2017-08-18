package json

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("For UUID", func() {
	Context("Marshalling UUID to JSON", func() {
		sampleUuid := uuid.NewV4()

		DescribeTable("String value must match canonical textual representation(8-4-4-4-12)",
			func(testedUuid uuid.UUID, expectedString string) {
				jsonValue, err := Uuid(testedUuid).MarshalJSON()

				Expect(err).To(Succeed())
				Expect(string(jsonValue)).To(Equal(expectedString))
			},
			Entry("Viable UUID", sampleUuid, fmt.Sprintf("\"%s\"", sampleUuid.String())),
			Entry("Nil UUID", uuid.Nil, "null"),
		)
	})

	Context("Unmarshalling UUID from JSON", func() {
		sampleUuid := uuid.NewV4()

		DescribeTable("Result UUID must match expected one",
			func(sampleJson string, expectedUuid uuid.UUID) {
				var testedUuid Uuid

				err := testedUuid.UnmarshalJSON([]byte(sampleJson))
				Expect(err).To(Succeed())
				Expect(uuid.UUID(testedUuid)).To(Equal(expectedUuid))
			},
			Entry("Viable UUID string", fmt.Sprintf("\"%s\"", sampleUuid.String()), sampleUuid),
			Entry("Empty string", "\"\"", uuid.Nil),
			Entry("Null JSON", "null", uuid.Nil),
		)
	})
})
