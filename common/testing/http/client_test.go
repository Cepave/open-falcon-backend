package http_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/dghubble/sling"

	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	mock "github.com/Cepave/open-falcon-backend/common/testing/http/gock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Test Suite")
}

var gockConfig = mock.GockConfigBuilder.NewConfigByRandom()

var _ = Describe("ResponseResult and sling", func() {
	var slingClient *sling.Sling

	BeforeEach(func() {
		gockConfig.New().
			Reply(http.StatusOK).
			JSON([]int{3, 55, 17})

		slingClient = sling.New().Base(gockConfig.GetUrl()).Get("res-2")
	})
	AfterEach(func() {
		gockConfig.Off()
	})

	Context("200 status", func() {
		It("Match body(string)", func() {
			testedResult := tHttp.NewResponseResultBySling(slingClient)
			Expect(testedResult.Response.StatusCode).To(Equal(200))

			GinkgoT().Logf("String body: %v", testedResult.GetBodyAsString())

			Expect(strings.TrimSpace(testedResult.GetBodyAsString())).To(Equal("[3,55,17]"))
			By("Get body again")
			Expect(strings.TrimSpace(testedResult.GetBodyAsString())).To(Equal("[3,55,17]"))
		})
		It("Match body(JSON)", func() {
			testedResult := tHttp.NewResponseResultBySling(slingClient)
			Expect(testedResult.Response.StatusCode).To(Equal(200))

			Expect(testedResult.GetBodyAsJson().MustArray()).To(HaveLen(3))
			By("Get body again")
			testedJson := testedResult.GetBodyAsJson()
			Expect(testedJson.GetIndex(0).MustInt()).To(BeEquivalentTo(3))
			Expect(testedJson.GetIndex(2).MustInt()).To(BeEquivalentTo(17))
		})
	})
})

var _ = Describe("Library of Gentleman HTTP client(h2non/gentleman)", func() {
	Context("200 Response(JSON)", func() {
		BeforeEach(func() {
			gockConfig.New().
				Get("/key").
				Reply(http.StatusOK).
				JSON(map[string]interface{}{
					"name": "Bob",
					"age":  30,
				})
		})

		AfterEach(func() {
			gockConfig.Off()
		})

		It("Match status and body", func() {
			clientConf := &tHttp.GentlemanClientConf{
				&tHttp.HttpClientConfig{
					Ssl:      false,
					Host:     gockConfig.Host,
					Port:     gockConfig.Port,
					Resource: "res-1",
				},
			}

			resp, err := clientConf.NewClient().
				Use(gockConfig.GentlemanT.Plugin()).
				Get().
				Path("/key").
				Send()

			Expect(err).To(Succeed())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var jsonResult = make(map[string]interface{})
			err = resp.JSON(&jsonResult)
			Expect(err).To(Succeed(), "Parse JSON error")

			Expect(jsonResult["name"]).To(Equal("Bob"))
			Expect(jsonResult["age"]).To(BeEquivalentTo(30))
		})
	})
})
