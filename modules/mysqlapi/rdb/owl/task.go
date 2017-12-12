package owl

import (
	"time"

	"github.com/Cepave/open-falcon-backend/common/db"
)

func RemoveOldScheduleLogs(t time.Time) int64 {
	sqlResult := DbFacade.SqlxDbCtrl.NamedExec(
		`
		DELETE FROM owl_schedule_log
		WHERE sl_start_time < :timing
		`,
		map[string]interface{}{
			"timing": t,
		},
	)

	return db.ToResultExt(sqlResult).RowsAffected()
}
