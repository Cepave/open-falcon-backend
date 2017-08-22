package http

import (
	"net/http"
	"strings"

	ogt "./gentleman"
	"github.com/dghubble/sling"
	"gopkg.in/h2non/gock.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResponseResult and sling", func() {
	var slingClient *sling.Sling

	BeforeEach(func() {
		gock.New("http://mock.testing.gt:13770/res-2").
			Reply(http.StatusOK).
			JSON([]int{3, 55, 17})

		slingClient = sling.New().Base("http://mock.testing.gt:13770").Get("res-2")
	})
	AfterEach(func() {
		ogt.DisableMock()
	})

	Context("200 status", func() {
		It("Match body(string)", func() {
			testedResult := NewResponseResultBySling(slingClient)
			Expect(testedResult.Response.StatusCode).To(Equal(200))

			GinkgoT().Logf("String body: %v", testedResult.GetBodyAsString())

			Expect(strings.TrimSpace(testedResult.GetBodyAsString())).To(Equal("[3,55,17]"))
			By("Get body again")
			Expect(strings.TrimSpace(testedResult.GetBodyAsString())).To(Equal("[3,55,17]"))
		})
		It("Match body(JSON)", func() {
			testedResult := NewResponseResultBySling(slingClient)
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
			gock.New("http://mock.testing.gt:10770/res-1").
				Get("/key").
				Reply(http.StatusOK).
				JSON(map[string]interface{}{
					"name": "Bob",
					"age":  30,
				})
		})

		AfterEach(func() {
			ogt.DisableMock()
		})

		It("Match status and body", func() {
			clientConf := &GentlemanClientConf{
				&HttpClientConfig{
					Ssl:      false,
					Host:     "mock.testing.gt",
					Port:     10770,
					Resource: "res-1",
				},
			}

			resp, err := clientConf.NewClient().
				Use(ogt.MockPlugin).
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
