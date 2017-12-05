package rdb

import (
	"fmt"
	"database/sql"
	"math/rand"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

const (
	defaultTimeout       = 2
	scheduleNamePrefix = "test-schedule"
	timeThreshold        = 500 * time.Millisecond
	deleteLockSql        =
		`
		DELETE FROM owl_schedule
		WHERE sch_name LIKE 'test-schedule%'
		`
	deleteLogSql =
		`
		DELETE sl
		FROM owl_schedule sch
			LEFT JOIN owl_schedule_log sl
			ON sch.sch_id = sl.sl_sch_id
		WHERE sch_name LIKE 'test-schedule%'
		`
)

var _ = Describe("Tests AcquireLock(...)", itSkip.PrependBeforeEach(func() {
	var (
		scheduleName    string
		defaultSchedule *model.Schedule
		defaultNow      time.Time

		ExpectSuccessSchedule = func(testSchedule *model.Schedule, testError error) {
			Expect(testError).NotTo(HaveOccurred())
			Expect(testSchedule.Uuid).NotTo(Equal(uuid.Nil))
		}

		ExpectLockAndLog = func(expSchedule *model.Schedule, expTime time.Time, expLogCount int) {
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
			)

			var (
				lockTable model.OwlSchedule
				logTable  model.OwlScheduleLog
			)
			GinkgoT().Log(defaultSchedule)

			By("Check lock")
			DbFacade.SqlxDbCtrl.Get(&lockTable, selectLockSql, expSchedule.Name)
			Expect(lockTable.IsLocked()).To(BeTrue())
			Expect(lockTable.LastUpdateTime).To(BeTemporally("~", expTime, timeThreshold))

			By("Check time")
			DbFacade.SqlxDbCtrl.Get(&logTable, selectLogSql, lockTable.Id)
			Expect(logTable.Timeout).To(Equal(expSchedule.Timeout))
			Expect(logTable.StartTime).To(BeTemporally("~", expTime, timeThreshold))

			By("Check log count")
			var count int
			DbFacade.SqlxDbCtrl.Get(&count, countLogSql, lockTable.Id)
			Expect(count).To(Equal(expLogCount))
		}
	)

	BeforeEach(func() {
		scheduleName = fmt.Sprintf("%s-%d", scheduleNamePrefix, rand.Int())
		defaultSchedule = model.NewSchedule(scheduleName, defaultTimeout)
		defaultNow = time.Now()
	})

	AfterEach(func() {
		inTx(deleteLogSql, deleteLockSql)
	})

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

var _ = Describe("Tests FreeLock(...)", itSkip.PrependBeforeEach(func() {
	var (
		scheduleName    string
		defaultSchedule *model.Schedule
		startTime      time.Time
	)

	BeforeEach(func() {
		scheduleName = fmt.Sprintf("%s-%d", scheduleNamePrefix, rand.Int())
		defaultSchedule = model.NewSchedule(scheduleName, 300)
		startTime = time.Now().Add(-120 * time.Second)
	})

	AfterEach(func() {
		inTx(deleteLogSql, deleteLockSql)
	})

	JustBeforeEach(func() {
		AcquireLock(defaultSchedule, startTime)
		GinkgoT().Logf("Setup new schedule: %v", defaultSchedule)
	})

	DescribeTable("The updated data should be as expected",
		func(expStatus model.TaskStatus, expMsg string) {
			var (
				expSchedule = defaultSchedule
				expTime     = startTime.Add(10 * time.Second)
			)
			FreeLock(expSchedule, expStatus, expMsg, expTime)

			var checkedResult = &struct {
				Locked bool `db:"sch_lock"`
				ModifyTime time.Time `db:"sch_modify_time"`
				Status model.TaskStatus `db:"sl_status"`
				Message sql.NullString `db:"sl_message"`
				EndTime time.Time `db:"sl_end_time"`
			}{}

			expectedMessage := sql.NullString{ String: expMsg, Valid: false }
			if expStatus == model.FAIL {
				expectedMessage.Valid = true
			}

			DbFacade.SqlxDbCtrl.Get(
				checkedResult,
				`
				SELECT sch_lock, sch_modify_time, sl_status, sl_message, sl_end_time
				FROM owl_schedule
					INNER JOIN
					owl_schedule_log
					ON sch_id = sl_sch_id
				WHERE sch_name = ?
				`,
				expSchedule.Name,
			)
			Expect(checkedResult).To(PointTo(
				MatchAllFields(Fields{
					"Locked": BeFalse(),
					"ModifyTime": BeTemporally("~", expTime, time.Second),
					"Status": Equal(expStatus),
					"Message": Equal(expectedMessage),
					"EndTime": BeTemporally("~", expTime, time.Second),
				}),
			))
		},
		Entry("DONE", model.DONE, ""),
		Entry("FAIL", model.FAIL, "Sample error message"),
	)
}))
