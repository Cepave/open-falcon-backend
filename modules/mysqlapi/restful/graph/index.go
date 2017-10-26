package graph

import (
	"time"

	"github.com/Cepave/open-falcon-backend/common/gin/mvc"

	db "github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/graph"
)

type relativeTimeParams struct {
	ForDays    int `mvc:"query[for_days] default[-1]" validate:"min=1|eq=-1"`
	ForMinutes int `mvc:"query[for_minutes] default[-1]" validate:"min=1|eq=-1"`
}

func vacuumEndpointIndex(
	q *relativeTimeParams,
) mvc.OutputBody {
	beforeTime := buildBeforeTime(time.Now(), q)

	result := db.VacuumEndpointIndex(beforeTime)

	return mvc.JsonOutputBody(
		map[string]interface{}{
			"before_time": beforeTime.Unix(),
			"affected_rows": map[string]interface{}{
				"endpoints": result.CountOfVacuumedEndpoints,
				"tags":      result.CountOfVacuumedTags,
				"counters":  result.CountOfVacuumedCounters,
			},
		},
	)
}

func buildBeforeTime(now time.Time, q *relativeTimeParams) time.Time {
	beforeTime := now

	if q.ForDays > 0 {
		beforeTime = beforeTime.Add(time.Duration(-q.ForDays) * time.Hour * 24)
		return beforeTime
	}

	if q.ForMinutes > 0 {
		beforeTime = beforeTime.Add(time.Duration(-q.ForMinutes) * time.Minute)
		return beforeTime
	}

	beforeTime = beforeTime.Add(-14 * time.Hour * 24)
	return beforeTime
}
