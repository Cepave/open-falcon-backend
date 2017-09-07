package json

import (
	"bytes"
	"fmt"

	sjson "github.com/bitly/go-simplejson"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("JsonExt object", func() {
	Context("GetxxxExt() functions", func() {
		sampleJson := sjson.New()
		sampleJson.Set("v1", 20)
		testedExt := ToJsonExt(sampleJson)

		It("Value(existing) should not be nil", func() {
			Expect(testedExt.GetExt("v1").IsNil()).NotTo(BeTrue())
		})

		It("Value(not-existing) should be nil", func() {
			Expect(testedExt.GetExt("v2").IsNil()).To(BeTrue())
		})

		It("Value should be nil(child of existing value)", func() {
			Expect(testedExt.GetPathExt("v1", "ck").IsNil()).To(BeTrue())
		})
	})
})

var _ = Describe("UnmarshalToJson(source)", func() {
	var (
		sampleName      = "j1"
		sampleValue     = "hello445"
		finalJsonString = fmt.Sprintf(`{ "%s": "%s" }`, sampleName, sampleValue)
	)

	Context("Viable JSON Content", func() {
		DescribeTable("The value of property \"js\" is \"hello445\"",
			func(sourceFunc func() interface{}) {
				testedJson := UnmarshalToJson(sourceFunc())
				Expect(testedJson.Get(sampleName).MustString()).To(Equal(sampleValue))
			},
			Entry("*go-simplejson.Json", func() interface{} {
				simpleJson := sjson.New()
				simpleJson.Set(sampleName, sampleValue)
				return simpleJson
			}),
			Entry("SimpleJsonMarshaler", func() interface{} {
				return &simpleJsonMarshaler{
					sampleName, sampleValue,
				}
			}),
			Entry("encoding/json.Marshaler", func() interface{} {
				return &jsonMarshaler{finalJsonString}
			}),
			Entry("string", func() interface{} {
				return finalJsonString
			}),
			Entry("[]byte", func() interface{} {
				return []byte(finalJsonString)
			}),
			Entry("io.Reader", func() interface{} {
				return bytes.NewBuffer([]byte(finalJsonString))
			}),
			Entry("Otherwise type", func() interface{} {
				return &sampleS{sampleValue}
			}),
		)

		DescribeTable("The value of property \"js\" is empty(\"\")",
			func(sourceFunc func() interface{}) {
				testedJson := UnmarshalToJson(sourceFunc())
				Expect(testedJson.Get(sampleName).MustString()).To(BeEmpty())
			},
			Entry("empty string", func() interface{} {
				return ""
			}),
			Entry("empty []byte", func() interface{} {
				return []byte{}
			}),
			Entry("nil []byte", func() interface{} {
				return ([]byte)(nil)
			}),
		)
	})
})

type sampleS struct {
	J1 string `json:"j1"`
}

type jsonMarshaler struct {
	finalString string
}

func (s *jsonMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(s.finalString), nil
}

type simpleJsonMarshaler struct {
	sampleName  string
	sampleValue string
}

func (s *simpleJsonMarshaler) MarshalSimpleJSON() (*sjson.Json, error) {
	newJson := sjson.New()
	newJson.Set(s.sampleName, s.sampleValue)
	return newJson, nil
}
