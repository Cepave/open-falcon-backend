package owl

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
)

var _ = Describe("Tests AcquireLock(...)", itSkip.PrependBeforeEach(func() {

	var (
		defaultSchedule     *Schedule
		defaultTimeout      = 0
		now                 time.Time
		defaultScheduleName = "test-schedule-test"

		lockTable     OwlSchedule
		logTable      OwlScheduleLog
		selectLockSql = `
			SELECT *
			FROM owl_schedule
			WHERE sch_name = ?
		`
		selectLogSql = `
			SELECT *
			FROM owl_schedule_log
			WHERE sl_sch_id = ?
			ORDER BY sl_start_time DESC
			LIMIT 1
		`
		countLogSql = `
			SELECT COUNT(*)
			FROM owl_schedule_log
			WHERE sl_sch_id = ?
		`
	)

	BeforeEach(func() {
		defaultSchedule = NewSchedule(defaultScheduleName, defaultTimeout)
		now = time.Now()
		lockTable = OwlSchedule{}
		logTable = OwlScheduleLog{}
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
			err := AcquireLock(defaultSchedule, now)

			GinkgoT().Logf("UUID=%v", defaultSchedule.Uuid)
			Expect(err).NotTo(HaveOccurred())
			Expect(defaultSchedule.Uuid).NotTo(Equal(uuid.Nil))

			By("Check lock")
			DbFacade.SqlxDbCtrl.Get(&lockTable, selectLockSql, defaultScheduleName)
			Expect(lockTable.isLocked()).To(BeTrue())

			By("Check time")
			DbFacade.SqlxDbCtrl.Get(&logTable, selectLogSql, lockTable.Id)
			Expect(logTable.StartTime).To(Equal(lockTable.LastUpdateTime))

			By("Check log count")
			var count int
			DbFacade.SqlxDbCtrl.Get(&count, countLogSql, lockTable.Id)
			Expect(count).To(Equal(1))
		})
	})

	Describe("A schedule has been created", func() {
		Context("lock is held too long", func() {
			It("should preempt the lock", func() {

				By("Lock is held by another schedule")
				s := NewSchedule(defaultScheduleName, 1)
				err := AcquireLock(s, now)
				GinkgoT().Logf("UUID=%v", s.Uuid)
				Expect(err).NotTo(HaveOccurred())

				By("Acquire lock from the stale task")
				ps := NewSchedule(defaultScheduleName, 0)
				now = now.Add(2 * time.Second)
				err = AcquireLock(ps, now)
				GinkgoT().Logf("UUID=%v", ps.Uuid)
				Expect(err).NotTo(HaveOccurred())
				Expect(ps.Uuid).NotTo(Equal(uuid.Nil))
			})
		})

		Context("lock is held but cannot determine the timeout", func() {
			var (
				crasheddefaultScheduleName = "test-schedule-crash"
				crashedSchedule            = `
					INSERT INTO owl_schedule (sch_name, sch_lock, sch_modify_time)
					VALUES ('test-schedule-crash', 1, '2020-01-01 12:00:00')
				`
			)
			BeforeEach(inTx(crashedSchedule))

			It("should preempt the lock", func() {
				By("Acquire lock from the crashed task")
				s := NewSchedule(crasheddefaultScheduleName, 0)
				err := AcquireLock(s, now)
				GinkgoT().Logf("UUID=%v", s.Uuid)
				Expect(err).NotTo(HaveOccurred())
				Expect(s.Uuid).NotTo(Equal(uuid.Nil))
			})
		})

		Context("lock is just held", func() {
			It("should trigger error", func() {
				By("Lock is held")
				s := NewSchedule(defaultScheduleName, 2)
				err := AcquireLock(s, now)
				Expect(err).NotTo(HaveOccurred())

				By("Acquire lock but get error")
				ps := NewSchedule(defaultScheduleName, 0)
				err = AcquireLock(ps, now)
				GinkgoT().Logf("Err=%s", err)
				Expect(err).To(HaveOccurred())
				Expect(ps.Uuid).To(Equal(uuid.Nil))
			})
		})
	})

}))
