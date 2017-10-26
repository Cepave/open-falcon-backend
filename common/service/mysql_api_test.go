package service

import (
	"net/http"

	ojson "github.com/Cepave/open-falcon-backend/common/json"
	apiModel "github.com/Cepave/open-falcon-backend/common/model/mysqlapi"
	mock "github.com/Cepave/open-falcon-backend/common/testing/http/gock"

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
			Rdb: &apiModel.AllRdbHealth{
				Dsn:             "root:!hide password!@tcp(mysql.mock.com:3306)/falcon_portal_it",
				OpenConnections: 12,
				PingResult:      0,
				PingMessage:     "",
				Portal: &apiModel.Rdb{
					Dsn:             "root:!hide password!@tcp(mysql.mock.com:3306)/falcon_portal_it",
					OpenConnections: 12,
					PingResult:      0,
					PingMessage:     "",
				},
				Graph: &apiModel.Rdb{
					Dsn:             "root:!hide password!@tcp(mysql.mock.com:3306)/graph_it",
					OpenConnections: 12,
					PingResult:      0,
					PingMessage:     "",
				},
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

	Context("[/health] of MySqlAPI", func() {
		It("viable content for 200 status", func() {
			gockConfig.New().
				Get(PATH).
				Reply(http.StatusOK).
				JSON(mockResp)
			view := testedSrv.GetHealth()

			GinkgoT().Logf("JSON Response for 200: %s", ojson.MarshalJSON(view))

			Expect(view.Address).To(Equal(gockConfig.GetUrl()))
			Expect(view.Message).To(BeEmpty())
		})

		It("nil value for 404 status", func() {
			gockConfig.New().
				Get(PATH).
				Reply(http.StatusNotFound)
			view := testedSrv.GetHealth()
			GinkgoT().Logf("JSON Response for 404: %s", ojson.MarshalJSON(view))

			Expect(view.Address).To(Equal(gockConfig.GetUrl()))
			Expect(view.Message).To(Not(BeEmpty()))
			Expect(view.Response).To(BeNil())
		})
	})
})
