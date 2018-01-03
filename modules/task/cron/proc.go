package cron

import (
	"github.com/Cepave/open-falcon-backend/modules/task/database"
	srv "github.com/Cepave/open-falcon-backend/modules/task/service"
)

func buildProcOfVacuumQueryObjects(forDays int) func() {
	return func() {
		logger.Infof("[Start] Vacuum query objects. For days: %d", forDays)

		result := database.QueryObjectService.VacuumQueryObjects(forDays)

		logger.Infof("[Finish] Vacuum [%d] query objects. Before time: [%s]", result.AffectedRows, result.GetBeforeTime())
	}
}

func buildProcOfVacuumGraphIndex(beforeDays int) func() {
	return func() {
		logger.Infof("[Start] Vacuum index of graph. For days: %d", beforeDays)

		result := srv.VacuumGraphIndex(beforeDays)

		logger.Infof("[Finish] Vacuum: %s", result)
	}
}

func buildProcOfClearTaskLogs(forDays int) func() {
	return func() {
		logger.Infof("[Start] Clear task log entries. For days: %d", forDays)

		result := database.ClearTaskLogEntryService.ClearLogEntries(forDays)

		logger.Infof("[Finish] Remove [%d] task log objects. Before time: [%s]", result.AffectedRows, result.GetBeforeTime())
	}
}
