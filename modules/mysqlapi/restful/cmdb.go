package restful

import (
	"net/http"

	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	oJson "github.com/Cepave/open-falcon-backend/common/json"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/service/cmdb"

	"github.com/satori/go.uuid"
)

func addNewCmdbSync() mvc.OutputBody {
	scheduleLog, err := cmdb.SyncDataFromBoss()

	if err == nil {
		return mvc.JsonOutputBody(
			map[string]interface{}{
				"sync_id":    scheduleLog.GetUuidString(),
				"start_time": scheduleLog.StartTime.Unix(),
			},
		)
	}

	if cannotLockError, ok := err.(*model.UnableToLockSchedule); ok {
		return mvc.JsonOutputBody2(
			http.StatusBadRequest,
			map[string]interface{}{
				"error_code":    1,
				"error_message": cannotLockError.Error(),
				"last_sync_id":  cannotLockError.Uuid.String(),
			},
		)
	} else {
		panic(err)
	}
}

func getSyncTask(
	params *struct {
		Uuid string `mvc:"param[uuid]"`
	},
) mvc.OutputBody {
	targetUuid, err := uuid.FromString(params.Uuid)
	if err != nil {
		panic(err)
	}

	scheduleLog := rdb.GetScheduleLog(targetUuid)

	if scheduleLog == nil {
		return mvc.NotFoundOutputBody
	}

	return mvc.JsonOutputBody(
		map[string]interface{}{
			"status":     scheduleLog.Status,
			"start_time": scheduleLog.StartTime.Unix(),
			"end_time":   oJson.JsonTime(scheduleLog.EndTime),
			"timeout":    scheduleLog.Timeout,
		},
	)
}
