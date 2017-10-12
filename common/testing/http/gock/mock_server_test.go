package gock

import (
	"net/http"
	"net/http/httptest"

	cl "github.com/Cepave/open-falcon-backend/common/http/client"
	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Http test server", func() {
	gockConfig := GockConfigBuilder.NewConfigByRandom()
	fakeServerConfig := &tHttp.FakeServerConfig{"127.0.0.1", 22081}

	var server *httptest.Server

	BeforeEach(func() {
		server = gockConfig.HttpTest.NewServer(fakeServerConfig)
		server.Start()
	})
	AfterEach(func() {
		server.Close()
		gockConfig.Off()
	})

	Context("Normal matching of request and response", func() {
		reqJson := map[string]interface{}{
			"new-1": "hello",
			"new-2": 33,
		}

		BeforeEach(func() {
			gockConfig.New().Post("/resource-1").
				MatchHeader("owl_key", "ks3").
				MatchType("json").
				MatchParam("q1", "v1").
				JSON(reqJson).
				Reply(http.StatusOK).
				SetHeader("owl_resp1", "9081").
				JSON(
					map[string]interface{}{
						"rep-1": 33,
					},
				)
		})

		It("Matching request and expected response", func() {
			client := (&tHttp.GentlemanClientConf{
				&tHttp.HttpClientConfig{
					Host: "127.0.0.1", Port: 22081,
				},
			}).NewClient()

			resp, err := client.Path("/resource-1").Post().
				AddHeader("owl_key", "ks3").
				AddQuery("q1", "v1").
				JSON(reqJson).
				Send()

			Expect(err).To(Succeed())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("owl_resp1")).To(Equal("9081"))

			jsonBody := cl.ToGentlemanResp(resp).MustGetJson()
			Expect(jsonBody.Get("rep-1").MustInt()).To(BeEquivalentTo(33))
		})
	})
})
