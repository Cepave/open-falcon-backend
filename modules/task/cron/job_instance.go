package cron

import (
	"time"

	"github.com/juju/errors"
	"github.com/satori/go.uuid"

	owlModel "github.com/Cepave/open-falcon-backend/common/model/owl"

	"github.com/Cepave/open-falcon-backend/modules/task/database"
)

type syncCmdbFromBoss struct {
	currentUuid uuid.UUID
	runningInfo *owlModel.SyncCmdbJobInfo
}

func (s *syncCmdbFromBoss) AfterJobStopped() {}
func (s *syncCmdbFromBoss) BeforeJobStart()  {}
func (s *syncCmdbFromBoss) BeforeJobStop() {
	if s.runningInfo != nil {
		logger.Warn("A sync job is running: %#v", s.runningInfo)
	}
}
func (s *syncCmdbFromBoss) Do() error {
	/**
	 * Starts a job
	 */
	if s.runningInfo == nil {
		newUuid, err := database.CmdbService.StartSyncJob()
		if err != nil {
			logger.Warnf("Start synchronization job of CMDB(from BOSS) has error: %v", err)
			return err
		}

		s.currentUuid = newUuid
		s.runningInfo = database.CmdbService.GetJobStatus(newUuid)
	}
	// :~)

	defer func() {
		s.currentUuid = uuid.Nil
		s.runningInfo = nil
	}()

	/**
	 * Pooling the status of a running job(not-timeout yet)
	 */
	for s.runningInfo.Status == owlModel.JobRunning && !s.runningInfo.IsTimeout() {
		time.Sleep(3 * time.Second)
		s.runningInfo = database.CmdbService.GetJobStatus(s.currentUuid)
	}
	// :~)

	switch s.runningInfo.Status {
	case owlModel.JobRunning:
		logger.Warnf("Job [%s] is timeout", s.currentUuid)
		return errors.Errorf("Job [%s] is timeout", s.currentUuid)
	case owlModel.JobFailed:
		logger.Warnf("Job [%s] is failed", s.currentUuid)
		return errors.Errorf("Job [%s] is failed", s.currentUuid)
	}

	return nil
}
