package owl

import (
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	json "github.com/Cepave/open-falcon-backend/common/json"
	ogk "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Input QueryObject", func() {
	var validateSrv = mvc.NewDefaultMvcConfig().Validator

	var sampleQuery *QueryObject

	BeforeEach(func() {
		sampleQuery = &QueryObject{
			NamedId: "sample-1",
			Content: json.VarBytes{0x1, 0x2, 0x3, 0x4},
			Md5Content: json.Bytes16{
				0x84, 0x19, 0x3f, 0x74, 0xda, 0x78, 0x2e, 0x03,
				0x91, 0x64, 0x77, 0xda, 0xca, 0x7e, 0x6c, 0x3c,
			},
		}
	})

	Context("Check passed data", func() {
		It("No validation error", func() {
			Expect(validateSrv.Struct(sampleQuery)).To(Succeed())
		})
	})

	Context("Check field rule", func() {
		type setErrorField func(queryObject *QueryObject)

		DescribeTable("Match of field name of validation error",
			func(fieldSetter setErrorField, expectedField string) {
				fieldSetter(sampleQuery)
				Expect(validateSrv.Struct(sampleQuery)).To(ogk.MatchFieldErrorOnName(expectedField))
			},
			Entry("QueryObject.NamedId",
				func(queryObject *QueryObject) {
					queryObject.NamedId = ""
				},
				"NamedId",
			),
			Entry("QueryObject.Content",
				func(queryObject *QueryObject) {
					queryObject.Content = nil
				},
				"Content",
			),
			Entry("QueryObject.Md5Content",
				func(queryObject *QueryObject) {
					queryObject.Md5Content = json.Bytes16{}
				},
				"Md5Content",
			),
		)
	})
})
