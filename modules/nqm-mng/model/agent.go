package model

import (
	"time"
)

type PingListLog struct {
	NumberOfTargets int32     `db:"apll_number_of_targets"`
	AccessTime      time.Time `db:"apll_time_access"`
	RefreshTime     time.Time `db:"apll_time_refresh"`
}
