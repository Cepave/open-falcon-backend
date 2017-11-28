package service

import (
	"fmt"
	"time"

	cutil "github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"
)

// ScheduleService is designed to be a Execute function with namespace.
var ScheduleService = &scheduleService{
	Execute: ScheduleExecutor,
}

type ScheduleCallback func() error

type scheduleService struct {
	Execute func(*model.Schedule, ScheduleCallback) error
}

func ScheduleExecutor(schedule *model.Schedule, callback ScheduleCallback) error {
	err := rdb.AcquireLock(schedule, time.Now())
	if err != nil {
		return err
	}

	var (
		errMsg    *string
		endStatus model.TaskStatus

		/**
		 * Because the panic in DB execution layer cannot be recovered, simply log it.
		 */
		freeLockTarget = func() {
			rdb.FreeLock(schedule, endStatus, errMsg, time.Now())
		}
		freeLockHandler = func(p interface{}) {
			logger.Errorf("During free lock of %s. Panic: %v", schedule, p)
		}
		// :~)
	)

	/**
	 * Since this go routine runs asynchronously, panic won't be captured by Gin.
	 * We should capture the panic by ourselves.
	 */
	callbackTarget := func() {
		err := callback()
		if err != nil {
			panic(err.Error())
		}
		errMsg = nil
		endStatus = model.DONE
		cutil.BuildPanicCapture(freeLockTarget, freeLockHandler)
	}

	go cutil.BuildPanicCapture(callbackTarget, func(p interface{}) {
		tmp := fmt.Sprint(p)
		errMsg = &tmp
		endStatus = model.FAIL
		cutil.BuildPanicCapture(freeLockTarget, freeLockHandler)
	})
	// :~)

	return nil
}
