package owl

import (
	"fmt"
	"net/http"

	json "github.com/Cepave/open-falcon-backend/common/json"
	ogko "github.com/Cepave/open-falcon-backend/common/testing/ginkgo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("[POST] on /owl/task/log/clear", itSkip.PrependBeforeEach(func() {
	Context("Remove older log entries", func() {
		BeforeEach(func() {
			inTx(
				`
					INSERT INTO owl_schedule(
					  sch_id, sch_name, sch_lock, sch_modify_time
					)
					VALUES
					(
					  1, 'test-schedule-1', 0, DATE_SUB(NOW(), INTERVAL 100 DAY)
					),
					(
					  2, 'test-schedule-2', 0, DATE_SUB(NOW(), INTERVAL 100 DAY)
					),
					(
					  3, 'test-schedule-3', 0, DATE_SUB(NOW(), INTERVAL 100 DAY)
					)
				`,
				`
				INSERT INTO owl_schedule_log(
					sl_uuid, sl_sch_id,
					sl_start_time, sl_end_time, sl_timeout,
					sl_status, sl_message
				)
				VALUES
				(
					x'209f18f4f89b42568e1e5270987c057d', 1,
					DATE_SUB(NOW(), INTERVAL 65 DAY), DATE_SUB(NOW(), INTERVAL 64 DAY), 172800,
					0, 'test-schedule-1'
				),
				(
					x'7011e902d4a848c184e242e8d71aa961', 1,
					DATE_SUB(NOW(), INTERVAL 95 DAY), DATE_SUB(NOW(), INTERVAL 94 DAY), 172800,
					0, 'test-schedule-1'
				),
				(
					x'349f18f4f89b42568e1e5270987c057d', 2,
					DATE_SUB(NOW(), INTERVAL 65 DAY), DATE_SUB(NOW(), INTERVAL 64 DAY), 172800,
					0, 'test-schedule-2'
				),
				(
					x'bdfa89f3df204071b2f48b359440e0fc', 2,
					DATE_SUB(NOW(), INTERVAL 85 DAY), NULL, 172800,
					2, 'test-schedule-2'
				),
				(
					x'109f1cf4f89b42568e1e5270987c057d', 3,
					DATE_SUB(NOW(), INTERVAL 65 DAY), DATE_SUB(NOW(), INTERVAL 64 DAY), 172800,
					0, 'test-schedule-3'
				)
				`,
			)
		})

		AfterEach(func() {
			inTx(
				`
				DELETE FROM owl_schedule_log
				WHERE sl_message LIKE 'test-schedule%'
				`,
				`
				DELETE FROM owl_schedule
				WHERE sch_name LIKE 'test-schedule%'
				`,
			)
		})

		DescribeTable("Check number of affected rows",
			func(forDays int, expectedAffectedRows int) {
				resp, err := httpClient.NewClient().
					Post().AddPath("/api/v1/owl/task/log/clear").
					AddQuery("for_days", fmt.Sprintf("%d", forDays)).
					Send()

				Expect(err).To(Succeed())

				jsonBody := gtResp(resp).MustGetJson()
				GinkgoT().Logf("[/owl/task/log/clear] JSON Result:\n%s", json.MarshalPrettyJSON(jsonBody))

				Expect(resp).To(ogko.MatchHttpStatus(http.StatusOK))
				Expect(jsonBody.Get("affected_rows").MustInt()).To(BeEquivalentTo(expectedAffectedRows))
			},
			Entry("Delete one row", 90, 1),
			Entry("Delete one row", 80, 2),
			Entry("Delete no row", 120, 0),
		)
	})
}))
