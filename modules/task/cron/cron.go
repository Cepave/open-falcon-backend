package cron

import (
	"fmt"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"

	"github.com/juju/errors"
	"github.com/robfig/cron"
)

var logger = log.NewDefaultLogger("info")

type TaskCronConfig struct {
	VacuumQueryObjects *VacuumQueryObjectsConf
	VacuumGraphIndex   *VacuumGraphIndexConf
}

// Configurations for vacuum of query objects
type VacuumQueryObjectsConf struct {
	Cron    string
	ForDays int
	Enable  bool
}

func (v *VacuumQueryObjectsConf) String() string {
	return fmt.Sprintf("For days: %d. Schedule: [%s].", v.ForDays, v.Cron)
}

type VacuumGraphIndexConf struct {
	Cron    string
	ForDays int
	Enable  bool
}

func NewCronServices(config *TaskCronConfig) *TaskCronService {
	cronServ := &TaskCronService{cron.New()}

	cronServ.addFunc(
		config.VacuumQueryObjects.Enable,
		config.VacuumQueryObjects.Cron,
		func() func() {
			return buildProcOfVacuumQueryObjects(config.VacuumQueryObjects.ForDays)
		},
		config.VacuumQueryObjects,
	)

	cronServ.addFunc(
		config.VacuumGraphIndex.Enable,
		config.VacuumGraphIndex.Cron,
		func() func() {
			return buildProcOfVacuumGraphIndex(config.VacuumGraphIndex.ForDays)
		},
		config.VacuumGraphIndex,
	)

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

func (s *TaskCronService) addFunc(enabled bool, cron string, procBuilder func() func(), configObject interface{}) {
	if !enabled {
		return
	}

	if err := errors.Annotatef(
		s.cronImpl.AddFunc(cron, procBuilder()),
		"Cannot add cron job: %v", configObject,
	); err != nil {
		panic(errors.Details(err))
	}
}
