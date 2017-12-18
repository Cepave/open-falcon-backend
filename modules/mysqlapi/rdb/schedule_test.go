package rdb

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"

	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

const (
	_DEFAULT_TIMEOUT = 60
	_SCHEDULE_PREFIX = "test-schedule"
)

var (
	deleteLockSql = fmt.Sprintf(
		`
		DELETE FROM owl_schedule
		WHERE sch_name LIKE '%s%%'
		`,
		_SCHEDULE_PREFIX,
	)
	deleteLogSql = fmt.Sprintf(
		`
		DELETE sl
		FROM owl_schedule sch
			LEFT JOIN owl_schedule_log sl
			ON sch.sch_id = sl.sl_sch_id
		WHERE sch_name LIKE '%s%%'
		`,
		_SCHEDULE_PREFIX,
	)
)

var _ = Describe("Tests AcquireLock(...)", itSkip.PrependBeforeEach(func() {
	AfterEach(func() {
		inTx(deleteLogSql, deleteLockSql)
	})

	Context("New schedule(not-used before)", func() {
		It("Successful locking", func() {
			sampleSchedule := model.NewSchedule(randomScheduleName(), _DEFAULT_TIMEOUT)
			now := time.Now()

			log, err := AcquireLock(sampleSchedule, now)

			assertSuccessScheduleLog(log, err)
			assertLockedSchedule(sampleSchedule, now, 1)
		})
	})

	Context("Existing schedule(used before)", func() {
		var scheduleName string
		var lastLockTime time.Time
		var lastJob *model.OwlScheduleLog

		BeforeEach(func() {
			scheduleName = randomScheduleName()
		})

		JustBeforeEach(func() {
			var err error

			lastJob, err = AcquireLock(model.NewSchedule(scheduleName, _DEFAULT_TIMEOUT), lastLockTime)
			Expect(err).To(Succeed())
		})

		Context("Lock is held too long", func() {
			BeforeEach(func() {
				lastLockTime = time.Now().Add(-(_DEFAULT_TIMEOUT + 2) * time.Second)
			})

			It("Should preempt the lock", func() {
				schedule := model.NewSchedule(scheduleName, _DEFAULT_TIMEOUT)
				now := time.Now()
				log, err := AcquireLock(schedule, now)

				assertSuccessScheduleLog(log, err)
				assertLockedSchedule(schedule, now, 2)
			})
		})

		Context("Lock is still held", func() {
			BeforeEach(func() {
				lastLockTime = time.Now()
			})

			It("should trigger error", func() {
				schedule := model.NewSchedule(scheduleName, _DEFAULT_TIMEOUT)
				log, err := AcquireLock(schedule, lastLockTime)

				Expect(err).To(HaveOccurred())
				Expect(err.(*model.UnableToLockSchedule).Uuid).To(BeEquivalentTo(lastJob.Uuid))
				Expect(log).To(BeNil())

				assertLockedSchedule(schedule, lastLockTime, 1)
			})
		})

		Context("Lock is held but the log is missing", func() {
			BeforeEach(func() {
				lastLockTime = time.Now().Add(-10 * time.Second)
			})

			JustBeforeEach(func() {
				inTx(deleteLogSql)
			})

			It("should preempt the lock", func() {
				schedule := model.NewSchedule(scheduleName, _DEFAULT_TIMEOUT)
				now := time.Now()
				log, err := AcquireLock(schedule, now)

				assertSuccessScheduleLog(log, err)
				assertLockedSchedule(schedule, now, 1)
			})
		})
	})
}))

var _ = Describe("Tests FreeLock(...)", itSkip.PrependBeforeEach(func() {
	var (
		sampleScheduleLog *model.OwlScheduleLog
		startTime         time.Time = time.Now().Add(-120 * time.Second)
		logTime                     = startTime.Add(10 * time.Second)
	)

	BeforeEach(func() {
		sampleSchedule := model.NewSchedule(randomScheduleName(), _DEFAULT_TIMEOUT)

		var err error
		sampleScheduleLog, err = AcquireLock(sampleSchedule, startTime)
		Expect(err).To(Succeed())

		GinkgoT().Logf("Setup new schedule: %v", sampleSchedule)
	})

	AfterEach(func() {
		inTx(deleteLogSql, deleteLockSql)
	})

	DescribeTable("The updated data should be as expected",
		func(expStatus owlModel.TaskStatus, expMsg string) {
			FreeLock(sampleScheduleLog, model.TaskStatus(expStatus), expMsg, logTime)

			var checkedResult = &struct {
				Locked     bool             `db:"sch_lock"`
				ModifyTime time.Time        `db:"sch_modify_time"`
				Status     model.TaskStatus `db:"sl_status"`
				Message    sql.NullString   `db:"sl_message"`
				EndTime    time.Time        `db:"sl_end_time"`
			}{}

			expectedMessage := sql.NullString{String: expMsg, Valid: false}
			if expStatus == owlModel.JobFailed {
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
				WHERE sch_id = ?
				`,
				sampleScheduleLog.SchId,
			)
			Expect(checkedResult).To(PointTo(
				MatchAllFields(Fields{
					"Locked":     BeFalse(),
					"ModifyTime": BeTemporally("~", logTime, time.Second),
					"Status":     Equal(model.TaskStatus(expStatus)),
					"Message":    Equal(expectedMessage),
					"EndTime":    BeTemporally("~", logTime, time.Second),
				}),
			))
		},
		Entry("DONE", owlModel.JobDone, ""),
		Entry("FAIL", owlModel.JobFailed, "Sample error message"),
	)
}))

var _ = Describe("Tests getting of schedule log", func() {
	BeforeEach(func() {
		inTx(
			`
			INSERT INTO owl_schedule(sch_id, sch_name, sch_lock, sch_modify_time)
			VALUES(302, 'gp.009', 0, NOW())
			`,
			`
			INSERT INTO owl_schedule_log(
				sl_sch_id, sl_uuid, sl_status,
				sl_timeout, sl_start_time, sl_end_time
			)
			VALUES(302, x'93a01bfc55ebcedf9c6d962ed51f496d', 0, 300, '2017-06-12 10:10:20', NOW()),
				(302, x'21a01bfc55ebcedf9c6d962ed51f496d', 1, 300, '2017-06-12 18:40:20', NULL)
			`,
		)
	})
	AfterEach(func() {
		inTx(
			`
			DELETE FROM owl_schedule_log
			WHERE sl_sch_id = 302
			`,
			`
			DELETE FROM owl_schedule
			WHERE sch_id = 302
			`,
		)
	})

	DescribeTable("Getting schedule log by UUID",
		func(sampleUuidString string, existing bool, hasEndTime bool) {
			sampleUuid, _ := uuid.FromString(sampleUuidString)

			testedLog := GetScheduleLog(sampleUuid)

			/**
			 * Asserts existing
			 */
			nilMatcher := BeNil()
			if existing {
				nilMatcher = Not(nilMatcher)
			}
			Expect(testedLog).To(nilMatcher)
			// :~)

			if !existing {
				return
			}

			/**
			 * Asserts end time
			 */
			endTimeMatcher := BeTrue()
			if !hasEndTime {
				endTimeMatcher = BeFalse()
			}
			Expect(testedLog.EndTime.IsNil()).NotTo(endTimeMatcher)
			// :~)
		},
		Entry("Existing log", "93a01bfc-55eb-cedf-9c6d-962ed51f496d", true, true),
		Entry("Existing log(null end time)", "21a01bfc-55eb-cedf-9c6d-962ed51f496d", true, false),
		Entry("Non-existing log", "41301bfc-55eb-cedf-9c6d-962ed51f496d", false, false),
	)
})

func randomScheduleName() string {
	return fmt.Sprintf("%s-%d", _SCHEDULE_PREFIX, rand.Int())
}

func assertSuccessScheduleLog(scheduleLog *model.OwlScheduleLog, testError error) {
	Expect(testError).To(Succeed())
	Expect(scheduleLog.Uuid).NotTo(Equal(uuid.Nil))
}

func assertLockedSchedule(scheduleContent *model.Schedule, expTime time.Time, expLogCount int) {
	type scheduleData struct {
		Lock           model.LockStatus `db:"sch_lock"`
		LastUpdateTime time.Time        `db:"sch_modify_time"`
		Timeout        int32            `db:"sl_timeout"`
		StartTime      time.Time        `db:"sl_start_time"`
		TotalCount     int              `db:"count_log"`
	}

	var testedData = &scheduleData{}

	DbFacade.SqlxDbCtrl.Get(
		testedData,
		`
		SELECT sch_lock, sch_modify_time, lm.sl_timeout, lm.sl_start_time,
			ll.count_log
		FROM owl_schedule
			LEFT OUTER JOIN
			owl_schedule_log AS lm
			ON sch_id = lm.sl_sch_id,
			(
				SELECT COUNT(sl_uuid) AS count_log
				FROM owl_schedule
					INNER JOIN
					owl_schedule_log
					ON sch_id = sl_sch_id
						AND sch_name = ?
			) AS ll
		WHERE sch_name = ?
		ORDER BY lm.sl_start_time DESC
		LIMIT 1
		`,
		scheduleContent.Name, scheduleContent.Name,
	)
	Expect(testedData).To(PointTo(
		MatchFields(IgnoreExtras, Fields{
			"Lock":           Equal(model.LOCKED),
			"LastUpdateTime": BeTemporally("~", expTime, time.Second),
			"Timeout":        Equal(scheduleContent.Timeout),
			"StartTime":      BeTemporally("~", expTime, time.Second),
			"TotalCount":     Equal(expLogCount),
		}),
	))
}
