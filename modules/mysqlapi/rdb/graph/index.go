package graph

import (
	"time"

	"github.com/Cepave/open-falcon-backend/common/db"
)

type VacuumEndpointResult struct {
	CountOfVacuumedEndpoints int64
	CountOfVacuumedCounters  int64
	CountOfVacuumedTags      int64
}

func VacuumEndpointIndex(beforeTime time.Time) *VacuumEndpointResult {
	result := &VacuumEndpointResult{}

	result.CountOfVacuumedCounters = db.ToResultExt(DbFacade.SqlxDbCtrl.NamedExec(
		`
		DELETE FROM endpoint_counter
		WHERE ts < :time_value
		`,
		map[string]interface{}{
			"time_value": beforeTime.Unix(),
		},
	)).RowsAffected()
	result.CountOfVacuumedTags = db.ToResultExt(DbFacade.SqlxDbCtrl.NamedExec(
		`
		DELETE FROM tag_endpoint
		WHERE ts < :time_value
		`,
		map[string]interface{}{
			"time_value": beforeTime.Unix(),
		},
	)).RowsAffected()
	result.CountOfVacuumedEndpoints = db.ToResultExt(DbFacade.SqlxDbCtrl.NamedExec(
		`
		DELETE FROM endpoint
		WHERE ts < :time_value
		`,
		map[string]interface{}{
			"time_value": beforeTime.Unix(),
		},
	)).RowsAffected()

	return result
}
