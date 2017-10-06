package ginkgo

import (
	"net/http"
	"net/http/httptest"

	ohttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	"github.com/Cepave/open-falcon-backend/common/testing/http/gock"
	gt "gopkg.in/h2non/gentleman.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("MatchHttpStatus(int)", func() {
	respRecorder := httptest.NewRecorder()
	respRecorder.WriteHeader(http.StatusOK)

	DescribeTable("Matching result is true",
		func(actual interface{}) {
			Expect(actual).To(MatchHttpStatus(200))
		},
		Entry("By *http.Response", respRecorder.Result()),
		Entry("By *testing/http.ResponseResult", ohttp.NewResponseResultByResponse(
			respRecorder.Result(),
		)),
		Entry("By *h2non/gentleman.ResponseResult", &gt.Response{
			RawResponse: respRecorder.Result(),
		}),
	)

	DescribeTable("Matching result is false",
		func(actual interface{}) {
			Expect(actual).ToNot(MatchHttpStatus(400))
		},
		Entry("By *http.Response", respRecorder.Result()),
		Entry("By *testing/http.ResponseResult", ohttp.NewResponseResultByResponse(
			respRecorder.Result(),
		)),
		Entry("By *h2non/gentleman.ResponseResult", &gt.Response{
			RawResponse: respRecorder.Result(),
		}),
	)
})

var _ = Describe("MatchHttpBodyAsJson(interface{})", func() {
	sampleJson := `{ "name": "joe", "age": 33 }`

	Context("For \"*http.Response\"", func() {
		newResp := func() *http.Response {
			respRecorder := httptest.NewRecorder()
			respRecorder.Header().Set("Content-Type", "application/json")
			respRecorder.WriteHeader(http.StatusOK)
			respRecorder.WriteString(sampleJson)

			return respRecorder.Result()
		}
		DescribeTable("Matching result is true",
			func(actual interface{}) {
				Expect(actual).To(MatchHttpBodyAsJson(sampleJson))
			},
			Entry("By *http.Response", newResp()),
			Entry("By *testing/http.ResponseResult", ohttp.NewResponseResultByResponse(newResp())),
		)

		DescribeTable("Matching result is false",
			func(actual interface{}) {
				Expect(actual).ToNot(MatchHttpBodyAsJson(`{ "name": "joe", "age": 34 }`))
			},
			Entry("By *http.Response", newResp()),
			Entry("By *testing/http.ResponseResult", ohttp.NewResponseResultByResponse(newResp())),
		)
	})

	Context("For \"*h2non/gentleman.ResponseResult\"", func() {
		var gockConfig *gock.GockConfig
		var jsonResp map[string]interface{}

		AfterEach(func() {
			gockConfig.Off()
		})

		JustBeforeEach(func() {
			GinkgoT().Logf("Set up Mock")
			gockConfig = gock.GockConfigBuilder.NewConfigByRandom()
			gockConfig.New().Get("/gg-sample-1").
				Reply(http.StatusOK).
				JSON(jsonResp)
		})

		Context("Matching result is true", func() {
			BeforeEach(func() {
				GinkgoT().Logf("Set up JSON")
				jsonResp = map[string]interface{}{
					"name": "joe",
					"age":  33,
				}
			})
			It("The JSON should match expected one", func() {
				resp, err := gockConfig.GentlemanT.NewClient().
					Get().
					Path("/gg-sample-1").
					Send()
				Expect(err).To(Succeed())

				Expect(resp).To(MatchHttpBodyAsJson(sampleJson))
			})
		})

		Context("Matching result is false", func() {
			BeforeEach(func() {
				jsonResp = map[string]interface{}{
					"name": "bob",
					"age":  33,
				}
			})

			It("The JSON should not match expected one", func() {
				resp, err := gockConfig.GentlemanT.NewClient().
					Get().
					Path("/gg-sample-1").
					Send()
				Expect(err).To(Succeed())

				Expect(resp).NotTo(MatchHttpBodyAsJson(sampleJson))
			})
		})
	})
})
