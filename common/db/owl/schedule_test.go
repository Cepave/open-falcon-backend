package owl

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Tests AcquireLock(...)", itSkip.PrependBeforeEach(func() {

	var (
		scheduleName = "test-schedule-3"
	)

	AfterEach(inTx(
		`
		DELETE sl
		FROM owl_schedule sch
		LEFT JOIN owl_schedule_log sl
		ON sch.sch_id = sl.sl_sch_id
		WHERE sch_name LIKE 'test-schedule-%'
		`,
		`DELETE FROM owl_schedule WHERE sch_name LIKE 'test-schedule-%'`,
	))

	Context("when schedule non-exists", func() {
		It("should return uuid correctly", func() {
			s := &Schedule{
				Name:    scheduleName,
				Timeout: 1,
			}
			err := AcquireLock(s)

			GinkgoT().Logf("UUID=%v", s.Uuid)
			Expect(err).NotTo(HaveOccurred())
		})
	})

}))
