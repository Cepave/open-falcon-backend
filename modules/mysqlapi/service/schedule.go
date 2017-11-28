package service

import (
	"fmt"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"
)

type ScheduleCallback func() error

func Execute(schedule *model.Schedule, callback ScheduleCallback) error {
	err := rdb.AcquireLock(schedule, time.Now())
	if err != nil {
		return err
	}

	go func() {
		var err error

		/**
		 * Free lock after callback is finished
		 */
		defer func() {
			var (
				errMsg    string
				endStatus model.TaskStatus
			)

			if p := recover(); p != nil {
				endStatus = model.FAIL
				errMsg = fmt.Sprint(p)
			} else if err != nil {
				endStatus = model.FAIL
				errMsg = err.Error()
			} else {
				endStatus = model.DONE
			}

			rdb.FreeLock(schedule, endStatus, errMsg, time.Now())
		}()
		// :~)

		err = callback()
	}()

	return nil
}
