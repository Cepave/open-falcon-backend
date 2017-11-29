package cron

import (
	"fmt"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"

	"github.com/juju/errors"
	"github.com/robfig/cron"
)

var logger = log.NewDefaultLogger("info")

type commonJobConfig interface {
	isEnable() bool
	getSchedule() string
	buildJob() func()
}

type TaskCronConfig struct {
	VacuumQueryObjects  *VacuumQueryObjectsConf
	VacuumGraphIndex    *VacuumGraphIndexConf
	ClearTaskLogEntries *ClearTaskLogEntriesConf
}

// Configurations for vacuum of query objects
type VacuumQueryObjectsConf struct {
	Cron    string
	ForDays int
	Enable  bool
}

func (v *VacuumQueryObjectsConf) String() string {
	return fmt.Sprintf("[Vacuum Query Object] For days: %d. Schedule: [%s].", v.ForDays, v.Cron)
}
func (v *VacuumQueryObjectsConf) isEnable() bool {
	return v.Enable
}
func (v *VacuumQueryObjectsConf) getSchedule() string {
	return v.Cron
}
func (v *VacuumQueryObjectsConf) buildJob() func() {
	return buildProcOfVacuumQueryObjects(v.ForDays)
}

// Configurations for vacuum of graph index
type VacuumGraphIndexConf struct {
	Cron    string
	ForDays int
	Enable  bool
}

func (v *VacuumGraphIndexConf) String() string {
	return fmt.Sprintf("[Vacuum Graph Index] For days: %d. Schedule: [%s].", v.ForDays, v.Cron)
}
func (v *VacuumGraphIndexConf) isEnable() bool {
	return v.Enable
}
func (v *VacuumGraphIndexConf) getSchedule() string {
	return v.Cron
}
func (v *VacuumGraphIndexConf) buildJob() func() {
	return buildProcOfVacuumGraphIndex(v.ForDays)
}

// Configurations for vacuum of query objects
type ClearTaskLogEntriesConf struct {
	Cron    string
	ForDays int
	Enable  bool
}

func (c *ClearTaskLogEntriesConf) String() string {
	return fmt.Sprintf("[Clear Task Log Entries] For days: %d. Schedule: [%s].", c.ForDays, c.Cron)
}
func (c *ClearTaskLogEntriesConf) isEnable() bool {
	return c.Enable
}
func (c *ClearTaskLogEntriesConf) getSchedule() string {
	return c.Cron
}
func (c *ClearTaskLogEntriesConf) buildJob() func() {
	return clearLogs(c.ForDays)
}

func NewCronServices(config *TaskCronConfig) *TaskCronService {
	cronServ := &TaskCronService{cron.New()}

	cronServ.addFunc(config.VacuumQueryObjects)
	cronServ.addFunc(config.VacuumGraphIndex)
	cronServ.addFunc(config.ClearTaskLogEntries)

	return cronServ
}

type TaskCronService struct {
	cronImpl *cron.Cron
}

func (s *TaskCronService) Start() {
	logger.Info("Start cron(by \"robfig/cron\") service.")
	s.cronImpl.Start()
}
func (s *TaskCronService) Stop() {
	logger.Info("Stop cron(by \"robfig/cron\") service.")
	s.cronImpl.Stop()
}

func (s *TaskCronService) addFunc(jobConfig commonJobConfig) {
	if !jobConfig.isEnable() {
		logger.Infof("Job is disabled: %s", jobConfig)
		return
	}

	logger.Infof("Job is enabled: %s", jobConfig)

	if err := errors.Annotate(
		s.cronImpl.AddFunc(jobConfig.getSchedule(), jobConfig.buildJob()), "Cannot add cron job",
	); err != nil {
		panic(errors.Details(err))
	}
}
