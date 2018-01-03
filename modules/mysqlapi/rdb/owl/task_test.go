package owl

import (
	t "github.com/Cepave/open-falcon-backend/common/testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("owl_schedule_log", itSkip.PrependBeforeEach(func() {
	Context("Remove deprecated log entries", func() {
		BeforeEach(func() {
			inTx(
				`
				INSERT INTO owl_schedule(
					sch_id, sch_name, sch_lock, sch_modify_time
				)
				VALUES
					(1, 'test-schedule-1', 0, '2014-05-01T00:00:00'),
					(2, 'test-schedule-2', 0, '2014-05-01T00:00:00'),
					(3, 'test-schedule-3', 0, '2014-05-01T00:00:00')
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
				  '2014-05-06T20:14:43', '2014-05-07T20:14:43', 172800,
				  0, 'test-schedule-1'
				),
				(
				  x'7011e902d4a848c184e242e8d71aa961', 1,
				  '2014-04-06T20:14:43', '2014-04-07T20:14:43', 172800,
				  0, 'test-schedule-1'
				),
				(
				  x'349f18f4f89b42568e1e5270987c057d', 2,
				  '2014-05-06T20:14:43', '2014-05-07T20:14:43', 172800,
				  0, 'test-schedule-2'
				),
				(
				  x'bdfa89f3df204071b2f48b359440e0fc', 2,
				  '2014-04-16T20:14:43', NULL, 172800,
				  2, 'test-schedule-2'
				),
				(
				  x'109f1cf4f89b42568e1e5270987c057d', 3,
				  '2014-05-06T20:14:43', '2014-05-07T20:14:43', 172800,
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

		It("Check the number of removed rows", func() {
			sampleTime := t.ParseTimeByGinkgo("2014-04-10T00:00:00+08:00")
			By("Remove 1 entry")
			Expect(RemoveOldScheduleLogs(sampleTime)).To(BeEquivalentTo(1))

			By("Nothing to be removed")
			Expect(RemoveOldQueryObject(sampleTime)).To(BeEquivalentTo(0))
		})
	})
}))
