package graph

import (
	"net/http"

	mock "github.com/Cepave/open-falcon-backend/common/testing/http/gock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var gockConfig = mock.GockConfigBuilder.NewConfigByRandom()

var _ = Describe("RESTful client for graph", func() {
	mysqlApiConfig := gockConfig.NewRestfulClientConfig()

	testedSrv := NewGraphService(
		&GraphServiceConfig{mysqlApiConfig},
	)

	AfterEach(func() {
		gockConfig.Off()
	})

	Context("[POST] /api/v1/graph/endpoint-index/vacuum", func() {
		sampleTime := 20879123

		BeforeEach(func() {
			gockConfig.New().
				MatchParam("for_days", "4").
				Post("/api/v1/graph/endpoint-index/vacuum").
				Reply(http.StatusOK).
				JSON(map[string]interface{}{
					"before_time": sampleTime, // Unix time before which the access time of older query objects applied
					"affected_rows": map[string]interface{}{
						"endpoints": 10,
						"tags":      20,
						"counters":  30,
					},
				})
		})

		It("Result must match expected one", func() {
			result := testedSrv.VacuumIndex(
				&VacuumIndexConfig{BeforeDays: 4},
			)

			Expect(result.GetBeforeTime().Unix()).To(BeEquivalentTo(sampleTime))
			Expect(result.AffectedRows).To(
				PointTo(MatchAllFields(Fields{
					"Endpoints": BeEquivalentTo(10),
					"Tags":      BeEquivalentTo(20),
					"Counters":  BeEquivalentTo(30),
				})),
			)
		})
	})
})
