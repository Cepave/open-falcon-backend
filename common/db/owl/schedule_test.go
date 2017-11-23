package owl

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
)

var _ = Describe("Tests AcquireLock(...)", itSkip.PrependBeforeEach(func() {

	var (
		scheduleName    = "test-schedule-3"
		crashedSchedule = `
			INSERT INTO owl_schedule (sch_name, sch_lock, sch_modify_time)
			VALUES ('test-schedule-3', 1, '2020-01-01 12:00:00')
		`
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

	Context("Schedule is new", func() {
		It("should acquire the lock", func() {
			s := &Schedule{
				Name: scheduleName,
			}
			err := AcquireLock(s)

			GinkgoT().Logf("UUID=%v", s.Uuid)
			Expect(err).NotTo(HaveOccurred())
			Expect(s.Uuid).NotTo(Equal(uuid.Nil))
		})
	})

	Describe("A schedule has been created", func() {
		Context("lock is held too long", func() {
			It("should preempt the lock", func() {
				s := &Schedule{
					Name:    scheduleName,
					Timeout: 0,
				}

				By("Lock is held by another schedule")
				err := AcquireLock(s)
				GinkgoT().Logf("UUID=%v", s.Uuid)
				Expect(err).NotTo(HaveOccurred())

				By("Acquire lock from the stale task")
				time.Sleep(time.Second)
				uuidPrev := s.Uuid
				err = AcquireLock(s)
				GinkgoT().Logf("UUID=%v", s.Uuid)
				Expect(err).NotTo(HaveOccurred())
				Expect(s.Uuid).NotTo(Equal(uuidPrev))
			})
		})

		Context("lock is held but cannot determine the timeout", func() {
			BeforeEach(inTx(crashedSchedule))

			It("should preempt the lock", func() {
				s := &Schedule{
					Name: scheduleName,
				}

				By("Acquire lock from the crashed task")
				err := AcquireLock(s)
				GinkgoT().Logf("UUID=%v", s.Uuid)
				Expect(err).NotTo(HaveOccurred())
				Expect(s.Uuid).NotTo(Equal(uuid.Nil))
			})
		})

		Context("lock is just held", func() {
			It("should trigger error", func() {
				s := &Schedule{
					Name:    scheduleName,
					Timeout: 2,
				}

				By("Lock is held")
				err := AcquireLock(s)
				Expect(err).NotTo(HaveOccurred())
				Expect(s.Uuid).NotTo(Equal(uuid.Nil))

				By("Acquire lock but get error")
				err = AcquireLock(s)
				GinkgoT().Logf("Err=%s", err)
				Expect(err).To(HaveOccurred())
			})
		})
	})

}))
