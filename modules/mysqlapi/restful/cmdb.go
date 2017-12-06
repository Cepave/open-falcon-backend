package restful

import (
	"net/http"
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"

	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/cmdb"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/service"
)

func addNewCmdbSync(
	sourceData *model.SyncForAdding,
) mvc.OutputBody {
	task := model.NewSchedule("import.imdb", 300)

	scheduleLog, err := service.ScheduleService.Execute(
		task,
		func() error {
			cmdb.SyncForHosts(sourceData)
			return nil
		},
	)

	if err != nil {
		if cannotLockError, ok := err.(*model.UnableToLockSchedule); ok {
			return mvc.JsonOutputBody2(
				http.StatusBadRequest,
				map[string]interface{} {
					"error_code": 1,
					"error_message": cannotLockError.Error(),
				},
			)
		} else {
			panic(err)
		}
	}

	return mvc.JsonOutputBody(
		map[string]interface{} {
			"sync_id": scheduleLog.GetUuidString(),
			"start_time": scheduleLog.StartTime.Unix(),
		},
	)
}
