package http_test

import (
	"net/http"
	"testing"

	mock "github.com/Cepave/open-falcon-backend/common/testing/http/gock"

	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	"github.com/Cepave/open-falcon-backend/common/http/client"
	tg "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var gockConfig = mock.GockConfigBuilder.NewConfigByRandom()

var _ = Describe("Build a new client", func() {
	var (
		testedSrv    *oHttp.ApiService
		clientConfig *client.HttpClientConfig
	)

	JustBeforeEach(func() {
		testedSrv = oHttp.NewApiService(
			&oHttp.RestfulClientConfig{
				HttpClientConfig: clientConfig,
			},
		)
	})

	AfterEach(func() {
		gockConfig.Off()
	})

	Context("Without resource(sub-path)", func() {
		BeforeEach(func() {
			clientConfig = client.NewDefaultConfig()
			clientConfig.Url = gockConfig.GetUrl()

			gockConfig.New().Get("/ball").
				Reply(http.StatusOK).
				JSON(map[string]interface{}{
					"v1": 918,
				})
		})

		It("Response is success", func() {
			resp, err := testedSrv.NewClient().
				Use(gockConfig.GentlemanT.Plugin()).
				Get().
				AddPath("/ball").
				Send()

			Expect(err).To(Succeed())
			Expect(resp).To(tg.MatchHttpStatus(http.StatusOK))
		})
	})

	Context("With resource(sub-path)", func() {
		BeforeEach(func() {
			clientConfig = client.NewDefaultConfig()
			clientConfig.Url = gockConfig.GetUrl() + "/v13"

			gockConfig.New().Get("/v13/ball").
				Reply(http.StatusOK).
				JSON(map[string]interface{}{
					"v1": 918,
				})
		})

		It("Response is success", func() {
			resp, err := testedSrv.NewClient().
				Use(gockConfig.GentlemanT.Plugin()).
				Get().
				AddPath("/ball").
				Send()

			Expect(err).To(Succeed())
			Expect(resp).To(tg.MatchHttpStatus(http.StatusOK))
		})
	})
})

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}
