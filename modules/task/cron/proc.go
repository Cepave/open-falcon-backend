package cron

import (
	"github.com/Cepave/open-falcon-backend/modules/task/database"
)

func buildProcOfVacuumQueryObjects(forDays int) func() {
	return func() {
		logger.Infof("[Start] Vacuum query objects. For days: %d", forDays)

		result := database.QueryObjectService.VacuumQueryObjects(forDays)

		logger.Infof("[Finish] Vacuum [%d] query objects. Before time: [%s]", result.AffectedRows, result.GetBeforeTime())
	}
}
