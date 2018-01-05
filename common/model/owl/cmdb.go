package owl

import (
	"fmt"
	"time"

	ojson "github.com/Cepave/open-falcon-backend/common/json"
)

type UnableToLockCmdbSyncJob struct {
	Uuid ojson.Uuid `json:"last_sync_id"`
}

func (e *UnableToLockCmdbSyncJob) Error() string {
	return fmt.Sprintf("Some task is running: [%s]", e.Uuid.ToUuid().String())
}

type SyncCmdbJobInfo struct {
	StartTime ojson.JsonTime `json:"start_time"`
	EndTime   ojson.JsonTime `json:"end_time"`
	Timeout   int            `json:"timeout"`
	Status    TaskStatus     `json:"status"`
}

func (info *SyncCmdbJobInfo) IsTimeout() bool {
	startTime := time.Time(info.StartTime)

	return info.Status == JobRunning &&
		time.Now().After(startTime.Add(time.Duration(info.Timeout)*time.Second))
}
