package owl

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[RESTful client] Clear Task Log Entries", func() {
	mysqlApiConfig := gockConfig.NewRestfulClientConfig()

	testedSrv := NewClearLogService(
		ClearLogServiceConfig{mysqlApiConfig},
	)

	AfterEach(func() {
		gockConfig.Off()
	})

	Context("Clear Task Log Entries", func() {
		// 2017-05-14
		sampleTime := 1494720000

		BeforeEach(func() {
			gockConfig.New().
				Post("/api/v1/owl/task/log/clear").
				MatchParam("for_days", "91").
				Reply(http.StatusOK).
				JSON(map[string]interface{}{
					"before_time":   sampleTime,
					"affected_rows": 104,
				})
		})

		It("Affectd rows should match expected", func() {
			testedResult := testedSrv.ClearLogEntries(91)

			Expect(testedResult.AffectedRows).To(Equal(104))
			Expect(testedResult.GetBeforeTime().Unix()).To(BeEquivalentTo(sampleTime))
		})
	})
})
