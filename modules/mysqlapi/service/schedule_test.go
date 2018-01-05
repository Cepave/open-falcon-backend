package service

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"

	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Tests ScheduleService", itSkip.PrependBeforeEach(func() {
	const schedulePrefix = "ts-GP221"

	var (
		waitToken = &sync.WaitGroup{}
	)

	BeforeEach(func() {
		waitToken.Add(1)
	})

	AfterEach(func() {
		inTx(
			`
			DELETE sl
			FROM owl_schedule sch
				LEFT JOIN owl_schedule_log sl
				ON sch.sch_id = sl.sl_sch_id
			WHERE sch_name LIKE 'ts-GP221%'
			`,
			`
			DELETE FROM owl_schedule
			WHERE sch_name LIKE 'ts-GP221%'
			`,
		)
	})

	DescribeTable("Free lock & record log",
		func(
			name string, callback ScheduleCallback,
			expectedStatus owlModel.TaskStatus, matchMessage string,
		) {
			newSchedule := model.NewSchedule(fmt.Sprintf("%s-%s", schedulePrefix, name), 300)

			scheduleLog, err := ScheduleService.Execute(newSchedule, callback)
			Expect(err).To(Succeed())
			Expect(scheduleLog).To(PointTo(
				MatchFields(
					IgnoreExtras,
					Fields{
						"Status":    BeEquivalentTo(owlModel.JobRunning),
						"StartTime": BeTemporally("<=", time.Now(), time.Millisecond),
						"Timeout":   BeEquivalentTo(300),
						"SchId":     BeNumerically(">", 0),
					},
				),
			))

			waitToken.Wait()
			time.Sleep(time.Second)

			testedLog := getContentOfScheduleLog(newSchedule.Name)

			/**
			 * Asserts the status and message
			 */
			Expect(testedLog.Status).To(BeEquivalentTo(expectedStatus))

			messageMatcher := Equal("")
			if expectedStatus == owlModel.JobFailed {
				messageMatcher = MatchRegexp(matchMessage)
			}
			Expect(testedLog.Message.String).To(messageMatcher)
			// :~)
		},
		Entry("Callback runs successfully",
			"s1",
			func() error {
				defer waitToken.Done()
				defer GinkgoT().Logf("Finish execution of run successfully")
				return nil
			}, owlModel.JobDone, "",
		),
		Entry("Callback returns error",
			"e1",
			func() error {
				defer waitToken.Done()
				defer GinkgoT().Logf("Finish execution of run with error returned")
				return errors.New("Normal error")
			}, owlModel.JobFailed, "Error from.*",
		),
		Entry("Callback invokes panic",
			"p1",
			func() error {
				defer waitToken.Done()
				defer GinkgoT().Logf("Finish execution of run with PANIC!!")
				panic("Go Panic")
			}, owlModel.JobFailed, "Panic from.*",
		),
	)
}))

type scheduleLog struct {
	Message sql.NullString   `db:"sl_message"`
	Status  model.TaskStatus `db:"sl_status"`
}

func getContentOfScheduleLog(name string) *scheduleLog {
	newLog := &scheduleLog{}

	rdb.DbFacade.SqlxDbCtrl.Get(
		newLog,
		`
		SELECT sl_message, sl_status
		FROM owl_schedule
			INNER JOIN
			owl_schedule_log
			ON sch_id = sl_sch_id
				AND sch_name = ?
		`,
		name,
	)

	return newLog
}
