package ginkgo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type Zoo string

func (z Zoo) MarshalJSON() ([]byte, error) {
	return []byte("[8, 11, 2]"), nil
}

var _ = Describe("Tests MatchJson", func() {
	DescribeTable("Tests matching result is true",
		func(actual interface{}, expected interface{}) {
			Expect(actual).To(MatchJson(expected))
		},
		Entry("JSON string == JSON string", "[8, 11, 2]", "[8, 11, 2]"),
		Entry("JSON string == json.Marshaler", "[8, 11, 2]", Zoo("")),
		Entry("json.Marshaler == JSON string", Zoo(""), "[8, 11, 2]"),
	)

	DescribeTable("Tests matching result is false",
		func(actual interface{}, expected interface{}) {
			Expect(actual).ToNot(MatchJson(expected))
		},
		Entry("JSON string != JSON string", "[8, 11, 2]", "[9, 11, 2]"),
		Entry("JSON string != json.Marshaler", "[9, 11, 2]", Zoo("")),
		Entry("json.Marshaler != JSON string", Zoo(""), "[9, 11, 2]"),
	)
})
