package rdb

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
)

var _ = Describe("Tests AcquireLock(...)", itSkip.PrependBeforeEach(func() {

	const (
		defaultTimeout       = 2
		scheduleNameTemplate = "test-schedule-"
		deleteLockSql        = `
			DELETE FROM owl_schedule WHERE sch_name LIKE 'test-schedule-%'
		`
		deleteLogSql = `
			DELETE sl
			FROM owl_schedule sch
			LEFT JOIN owl_schedule_log sl
			ON sch.sch_id = sl.sl_sch_id
			WHERE sch_name LIKE 'test-schedule-%'
		`
	)

	var (
		scheduleName    string
		defaultSchedule *model.Schedule
		defaultNow      time.Time

		/**
		 * Helper function
		 */
		ExpectSuccessSchedule = func(testSchedule *model.Schedule, testError error) {
			GinkgoT().Logf("UUID=%v", testSchedule.Uuid)
			Expect(testError).NotTo(HaveOccurred())
			Expect(testSchedule.Uuid).NotTo(Equal(uuid.Nil))
		}

		ExpectLockAndLog = func(expSchedule *model.Schedule, expNow time.Time, expLogCount int) {
			const (
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
				timeThreshold = 500 * time.Millisecond
			)

			var (
				lockTable model.OwlSchedule
				logTable  model.OwlScheduleLog
			)

			By("Check lock")
			DbFacade.SqlxDbCtrl.Get(&lockTable, selectLockSql, expSchedule.Name)
			Expect(lockTable.IsLocked()).To(BeTrue())
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
		// :~)
	)

	BeforeEach(func() {
		scheduleName = scheduleNameTemplate + fmt.Sprint(rand.Int())
		GinkgoT().Logf("Name=%s", scheduleName)
		defaultSchedule = model.NewSchedule(scheduleName, defaultTimeout)
		defaultNow = time.Now()
	})

	AfterEach(func() {
		inTx(deleteLogSql, deleteLockSql)
	})

	Context("model.Schedule is new", func() {
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
				newCurrent := defaultNow.Add(time.Duration(thisTimeout) * time.Second)
				ps := model.NewSchedule(scheduleName, thisTimeout)
				err := AcquireLock(ps, newCurrent)

				ExpectSuccessSchedule(ps, err)
				ExpectLockAndLog(ps, newCurrent, 2)
			})
		})

		Context("lock is just held", func() {
			It("should trigger error", func() {
				thisTimeout := defaultTimeout + 1
				ps := model.NewSchedule(scheduleName, thisTimeout)
				err := AcquireLock(ps, defaultNow)

				Expect(err).To(HaveOccurred())
				Expect(ps.Uuid).To(Equal(uuid.Nil))

				ExpectLockAndLog(defaultSchedule, defaultNow, 1)
			})
		})

		Context("lock is held but cannot determine the timeout", func() {
			BeforeEach(func() {
				_ = AcquireLock(defaultSchedule, defaultNow)
			})

			JustBeforeEach(func() {
				inTx(deleteLogSql)
			})

			It("should preempt the lock", func() {
				By("Acquire lock from the crashed task")
				thisTimeout := defaultTimeout + 1
				newCurrent := defaultNow.Add(time.Duration(thisTimeout) * time.Second)
				sp := model.NewSchedule(scheduleName, thisTimeout)
				err := AcquireLock(sp, newCurrent)

				ExpectSuccessSchedule(sp, err)
				ExpectLockAndLog(sp, newCurrent, 1)
			})
		})

	})

}))
