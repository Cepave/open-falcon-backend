package http

import (
	"net/http"
	"net/http/httptest"

	"github.com/Cepave/open-falcon-backend/modules/hbs/g"

	oModel "github.com/Cepave/open-falcon-backend/common/model"
	tg "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"
	tHttp "github.com/Cepave/open-falcon-backend/common/testing/http"
	mock "github.com/Cepave/open-falcon-backend/common/testing/http/gock"
	apiModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	itClient = tHttp.GentlemanClientConf{itConfig}

	gockConfig       = mock.GockConfigBuilder.NewConfigByRandom()
	fakeServerConfig = &tHttp.FakeServerConfig{"127.0.0.1", 6040}
)

var _ = Describe("[Intg] Test GET on /api/v1/health", itEnabled(func() {
	const PATH = "/health"
	var (
		mockMysqlApiServer *httptest.Server
		mockResp           = &apiModel.HealthView{
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

	BeforeEach(func() {
		mockMysqlApiServer = gockConfig.HttpTest.NewServer(fakeServerConfig)
		mockMysqlApiServer.Start()
	})
	AfterEach(func() {
		mockMysqlApiServer.Close()
		mockMysqlApiServer = nil
		gockConfig.Off()
	})

	It("should return same Mysql-API info as gived", func() {
		gockConfig.New().
			Get(PATH).
			Reply(http.StatusOK).
			JSON(mockResp)

		resp, err := itClient.NewClient().Path("/api/v1/health").Get().Send()
		Expect(err).To(Succeed())
		defer resp.Close()

		Expect(resp).To(tg.MatchHttpStatus(http.StatusOK))
		GinkgoT().Logf("JSON Result: %s", resp.String())

		h := &g.HealthView{}
		Expect(resp.JSON(h)).To(Succeed())
		expVal := &g.HealthView{
			MysqlApi: &oModel.MysqlApi{
				Address:  fakeServerConfig.GetUrlString(),
				Message:  "",
				Response: mockResp,
			},
		}

		Expect(h.HealthCheck).To(Equal(expVal.MysqlApi.Response.Rdb.PingResult))
		Expect(h.MysqlApi).To(Equal(expVal.MysqlApi))
	})
}))
