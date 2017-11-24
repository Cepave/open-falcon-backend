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
		defaultTimeout      = 2
		defaultNow          time.Time
		defaultScheduleName = "test-schedule-test"

		lockTable OwlSchedule
		logTable  OwlScheduleLog

		ExpectSuccessSchedule = func(testSchedule *Schedule, testError error) {
			GinkgoT().Logf("UUID=%v", testSchedule.Uuid)
			Expect(testError).NotTo(HaveOccurred())
			Expect(testSchedule.Uuid).NotTo(Equal(uuid.Nil))
		}

		ExpectLockAndLog = func(expSchedule *Schedule, expNow time.Time, expLogCount int) {
			selectLockSql := `
				SELECT *
				FROM owl_schedule
				WHERE sch_name = ?
			`
			selectLogSql := `
				SELECT *
				FROM owl_schedule_log
				WHERE sl_sch_id = ?
				ORDER BY sl_start_time DESC
				LIMIT 1
			`
			countLogSql := `
				SELECT COUNT(*)
				FROM owl_schedule_log
				WHERE sl_sch_id = ?
			`
			timeThreshold := 500 * time.Millisecond

			By("Check lock")
			DbFacade.SqlxDbCtrl.Get(&lockTable, selectLockSql, expSchedule.Name)
			Expect(lockTable.isLocked()).To(BeTrue())
			Expect(lockTable.LastUpdateTime).To(BeTemporally("~", expNow, timeThreshold))

			By("Check time")
			DbFacade.SqlxDbCtrl.Get(&logTable, selectLogSql, lockTable.Id)
			Expect(logTable.Timeout).To(Equal(expSchedule.Timeout))
			Expect(logTable.StartTime).To(BeTemporally("~", expNow, timeThreshold))

			By("Check log count")
			var count int
			DbFacade.SqlxDbCtrl.Get(&count, countLogSql, lockTable.Id)
			Expect(count).To(Equal(expLogCount))
		}
	)

	BeforeEach(func() {
		defaultSchedule = NewSchedule(defaultScheduleName, defaultTimeout)
		defaultNow = time.Now()
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
			err := AcquireLock(defaultSchedule, defaultNow)

			ExpectSuccessSchedule(defaultSchedule, err)
			ExpectLockAndLog(defaultSchedule, defaultNow, 1)
		})
	})

	Context("A schedule has been created", func() {
		JustBeforeEach(func() {
			_ = AcquireLock(defaultSchedule, defaultNow)
		})

		Context("lock is held too long", func() {
			It("should preempt the lock", func() {
				thisTimeout := defaultTimeout + 1
				ps := NewSchedule(defaultScheduleName, thisTimeout)
				newCurrent := defaultNow.Add(time.Duration(thisTimeout) * time.Second)
				err := AcquireLock(ps, newCurrent)

				ExpectSuccessSchedule(ps, err)
				ExpectLockAndLog(ps, newCurrent, 2)
			})
		})

		Context("lock is just held", func() {
			It("should trigger error", func() {
				thisTimeout := defaultTimeout + 1
				ps := NewSchedule(defaultScheduleName, thisTimeout)
				err := AcquireLock(ps, defaultNow)

				Expect(err).To(HaveOccurred())
				Expect(ps.Uuid).To(Equal(uuid.Nil))

				ExpectLockAndLog(defaultSchedule, defaultNow, 1)
			})
		})

		// Context("lock is held but cannot determine the timeout", func() {
		// 	var (
		// 		crasheddefaultScheduleName = "test-schedule-crash"
		// 		crashedSchedule            = `
		// 				INSERT INTO owl_schedule (sch_name, sch_lock, sch_modify_time)
		// 				VALUES ('test-schedule-crash', 1, '2020-01-01 12:00:00')
		// 			`
		// 	)
		// 	BeforeEach(inTx(crashedSchedule))

		// 	It("should preempt the lock", func() {
		// 		By("Acquire lock from the crashed task")
		// 		s := NewSchedule(crasheddefaultScheduleName, 0)
		// 		err := AcquireLock(s, defaultNow)
		// 		GinkgoT().Logf("UUID=%v", s.Uuid)
		// 		Expect(err).NotTo(HaveOccurred())
		// 		Expect(s.Uuid).NotTo(Equal(uuid.Nil))
		// 	})
		// })

	})

}))
