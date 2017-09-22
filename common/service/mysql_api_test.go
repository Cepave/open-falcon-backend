package service

import (
	"net/http"

	mock "github.com/Cepave/open-falcon-backend/common/testing/http/gock"
	apiModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var gockConfig = mock.GockConfigBuilder.NewConfigByRandom()

var _ = Describe("[Unit] GetHealth()", func() {
	const PATH = "/health"
	var (
		mysqlApiConfig = gockConfig.NewRestfulClientConfig()

		testedSrv = NewMysqlApiService(
			&MysqlApiServiceConfig{mysqlApiConfig},
		)

		mockResp = apiModel.HealthView{
			Rdb: &apiModel.Rdb{
				Dsn:             "root:!hide password!@tcp(mysql.mock.com:3306)/falcon_portal",
				OpenConnections: 12,
				PingResult:      0,
				PingMessage:     "",
			},
			Http: &apiModel.Http{
				Listening: "0.0.0.0:12345",
			},
			Nqm: &apiModel.Nqm{
				&apiModel.Heartbeat{
					Count: 98765,
				},
			},
		}
	)

	AfterEach(func() {
		gockConfig.Off()
	})

	Context("when the connection is stable", func() {
		It("should return correct response if 200", func() {
			gockConfig.New().
				Get(PATH).
				Reply(http.StatusOK).
				JSON(mockResp)
			view := testedSrv.GetHealth()
			GinkgoT().Logf("Response of GetHealth: %#v", view)

			Expect(view.Address).To(Equal(gockConfig.GetUrl()))
			Expect(view.Message).To(BeEmpty())
		})

		It("should return nil response if 404", func() {
			gockConfig.New().
				Get(PATH).
				Reply(http.StatusNotFound)
			view := testedSrv.GetHealth()
			GinkgoT().Logf("Response of GetHealth: %#v", view)

			Expect(view.Address).To(Equal(gockConfig.GetUrl()))
			Expect(view.Message).To(Not(BeEmpty()))
			Expect(view.Response).To(BeNil())
		})
	})
})
