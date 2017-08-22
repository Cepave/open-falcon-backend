package http

import (
	"net/http"
	"net/http/httptest"

	sjson "github.com/bitly/go-simplejson"
	"github.com/dghubble/sling"
	"gopkg.in/h2non/gock.v1"

	ogt "./gentleman"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResponseResult and sling", func() {
	var sampleServer *httptest.Server
	var slingClient *sling.Sling

	BeforeEach(func() {
		sampleServer = httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("[3, 55, 17]"))
			}),
		)
		slingClient = sling.New().Base(sampleServer.URL).Get("")
	})
	AfterEach(func() {
		sampleServer.Close()
	})

	Context("200 status", func() {
		It("Match body(string)", func() {
			testedResult := NewResponseResultBySling(slingClient)
			Expect(testedResult.Response.StatusCode).To(Equal(200))

			Expect(testedResult.GetBodyAsString()).To(Equal("[3, 55, 17]"))
			By("Get body again")
			Expect(testedResult.GetBodyAsString()).To(Equal("[3, 55, 17]"))
		})

		It("Match Json(string)", func() {
			testedResult := NewResponseResultBySling(slingClient)
			Expect(testedResult.Response.StatusCode).To(Equal(200))

			expectedJson, _ := sjson.NewJson([]byte(`[3, 55, 17]`))
			Expect(testedResult.GetBodyAsJson()).To(Equal(expectedJson))
			By("Get body again")
			Expect(testedResult.GetBodyAsJson()).To(Equal(expectedJson))
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
