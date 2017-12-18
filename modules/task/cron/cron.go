package cron

import (
	log "github.com/Cepave/open-falcon-backend/common/logruslog"

	ocron "github.com/Cepave/open-falcon-backend/common/cron"

	"github.com/juju/errors"
	"github.com/robfig/cron"
)

var logger = log.NewDefaultLogger("info")

type cronJobConfig interface {
	isEnable() bool
	getSchedule() string
	buildJob() func()
}

type intervalJobConfig interface {
	taskName() string
	isEnable() bool
	getIntervalConfig() *ocron.IntervalConfig
	buildJob() ocron.Job
}

type TaskCronConfig struct {
	VacuumQueryObjects  *VacuumQueryObjectsConf
	VacuumGraphIndex    *VacuumGraphIndexConf
	ClearTaskLogEntries *ClearTaskLogEntriesConf
	SyncCmdbFromBoss    *SyncCmdbFromBossConf
}

func NewCronServices(cronConfig *TaskCronConfig) *TaskCronService {
	cronServ := &TaskCronService{
		cronImpl:     cron.New(),
		intervalImpl: ocron.NewIntervalService(),
	}

	cronServ.addCronJob(cronConfig.VacuumQueryObjects)
	cronServ.addCronJob(cronConfig.VacuumGraphIndex)
	cronServ.addCronJob(cronConfig.ClearTaskLogEntries)

	cronServ.addIntervalJob(cronConfig.SyncCmdbFromBoss)

	return cronServ
}

type TaskCronService struct {
	cronImpl     *cron.Cron
	intervalImpl *ocron.IntervalService
}

func (s *TaskCronService) Start() {
	logger.Info("Start cron/interval services.")

	s.cronImpl.Start()
	s.intervalImpl.Start()
}
func (s *TaskCronService) Stop() {
	logger.Info("Stop cron/interval services.")

	s.cronImpl.Stop()
	s.intervalImpl.Stop()
}

func (s *TaskCronService) addIntervalJob(jobConfig intervalJobConfig) {
	if !jobConfig.isEnable() {
		logger.Infof("Interval job is disabled: %s", jobConfig)
		return
	}

	logger.Infof("Interval job is enabled: %s", jobConfig)

	s.intervalImpl.Add(jobConfig.taskName(), jobConfig.getIntervalConfig(), jobConfig.buildJob())
}

func (s *TaskCronService) addCronJob(jobConfig cronJobConfig) {
	if !jobConfig.isEnable() {
		logger.Infof("Cron job is disabled: %s", jobConfig)
		return
	}

	logger.Infof("Cron job is enabled: %s", jobConfig)

	if err := errors.Annotate(
		s.cronImpl.AddFunc(jobConfig.getSchedule(), jobConfig.buildJob()), "Cannot add cron job",
	); err != nil {
		panic(errors.Details(err))
	}
}
