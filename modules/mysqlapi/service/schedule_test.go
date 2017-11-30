package service

import (
	"fmt"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tests ScheduleService", itSkip.PrependBeforeEach(func() {
	var (
		defaultSchedule = model.NewSchedule("test-scheduleService-", 0)
		defaultErrMsg   = "Default error message."
		defaultPanicMsg = "Default panic message."
	)

	const (
		deleteLockSql = `
			DELETE FROM owl_schedule WHERE sch_name LIKE 'test-scheduleService-%'
		`
		deleteLogSql = `
			DELETE sl
			FROM owl_schedule sch
			LEFT JOIN owl_schedule_log sl
			ON sch.sch_id = sl.sl_sch_id
			WHERE sch_name LIKE 'test-scheduleService-%'
		`
	)

	AfterEach(func() {
		inTx(deleteLogSql, deleteLockSql)
		time.Sleep(time.Second)
	})

	DescribeTable("Free lock & record log",
		func(callback ScheduleCallback, expStatus model.TaskStatus, errMsg *string) {
			err := ScheduleService.Execute(defaultSchedule, callback)
			Expect(err).NotTo(HaveOccurred())
		},
		Entry("Callback runs successfully",
			func() error {
				fmt.Println("Done")
				return nil
			}, model.DONE, nil,
		),
		Entry("Callback returns error",
			func() error {
				fmt.Println("Error")
				return fmt.Errorf(defaultErrMsg)
			}, model.FAIL, &defaultErrMsg,
		),
		Entry("Callback invokes panic",
			func() error {
				fmt.Println("Panic")
				panic(defaultPanicMsg)
				return nil
			}, model.FAIL, &defaultPanicMsg,
		),
	)
}))
