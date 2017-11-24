package owl

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
)

var _ = Describe("Tests AcquireLock(...)", itSkip.PrependBeforeEach(func() {

	var (
		now          time.Time
		scheduleName = "test-schedule-test"
	)

	BeforeEach(func() {
		now = time.Now()
	})

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

	Context("Schedule is new", func() {
		It("should acquire the lock", func() {
			s := NewSchedule(scheduleName, 0)
			err := AcquireLock(s, now)

			GinkgoT().Logf("UUID=%v", s.Uuid)
			Expect(err).NotTo(HaveOccurred())
			Expect(s.Uuid).NotTo(Equal(uuid.Nil))
		})
	})

	Describe("A schedule has been created", func() {
		Context("lock is held too long", func() {
			It("should preempt the lock", func() {

				By("Lock is held by another schedule")
				s := NewSchedule(scheduleName, 1)
				err := AcquireLock(s, now)
				GinkgoT().Logf("UUID=%v", s.Uuid)
				Expect(err).NotTo(HaveOccurred())

				By("Acquire lock from the stale task")
				ps := NewSchedule(scheduleName, 0)
				now = now.Add(2 * time.Second)
				err = AcquireLock(ps, now)
				GinkgoT().Logf("UUID=%v", ps.Uuid)
				Expect(err).NotTo(HaveOccurred())
				Expect(ps.Uuid).NotTo(Equal(uuid.Nil))
			})
		})

		Context("lock is held but cannot determine the timeout", func() {
			var (
				crashedScheduleName = "test-schedule-crash"
				crashedSchedule     = `
					INSERT INTO owl_schedule (sch_name, sch_lock, sch_modify_time)
					VALUES ('test-schedule-crash', 1, '2020-01-01 12:00:00')
				`
			)
			BeforeEach(inTx(crashedSchedule))

			It("should preempt the lock", func() {
				By("Acquire lock from the crashed task")
				s := NewSchedule(crashedScheduleName, 0)
				err := AcquireLock(s, now)
				GinkgoT().Logf("UUID=%v", s.Uuid)
				Expect(err).NotTo(HaveOccurred())
				Expect(s.Uuid).NotTo(Equal(uuid.Nil))
			})
		})

		Context("lock is just held", func() {
			It("should trigger error", func() {
				By("Lock is held")
				s := NewSchedule(scheduleName, 2)
				err := AcquireLock(s, now)
				Expect(err).NotTo(HaveOccurred())

				By("Acquire lock but get error")
				ps := NewSchedule(scheduleName, 0)
				err = AcquireLock(ps, now)
				GinkgoT().Logf("Err=%s", err)
				Expect(err).To(HaveOccurred())
				Expect(ps.Uuid).To(Equal(uuid.Nil))
			})
		})
	})

}))
