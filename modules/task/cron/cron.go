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

func NewCronServices(config *TaskCronConfig) *TaskCronService {
	cronImpl := cron.New()

	/**
	 * Constructs job for vacuum of query objects
	 */
	if config.VacuumQueryObjects.Enable {
		vacuumConfig := config.VacuumQueryObjects

		logger.Infof("Enable vacuum of query objects. %s", vacuumConfig)

		err := errors.Annotatef(
			cronImpl.AddFunc(vacuumConfig.Cron, buildProcOfVacuumQueryObjects(vacuumConfig.ForDays)),
			"Cannot add cron job for vacuum of query objects. Config: %s", vacuumConfig,
		)
		if err != nil {
			panic(errors.Details(err))
		}
	}
	// :~)

	return &TaskCronService{cronImpl}
}

type TaskCronService struct {
	cronImpl *cron.Cron
}

func (s *TaskCronService) Start() {
	logger.Info("Start cron(by robfig/cron) service.")
	s.cronImpl.Start()
}
func (s *TaskCronService) Stop() {
	logger.Info("Stop cron(by robfig/cron) service.")
	s.cronImpl.Stop()
}
