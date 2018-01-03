package cron

import (
	"fmt"
	"time"

	ocron "github.com/Cepave/open-falcon-backend/common/cron"
)

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
	return buildProcOfClearTaskLogs(c.ForDays)
}

// Configuration for synchronized job of CMDB from BOSS database
type SyncCmdbFromBossConf struct {
	InitialDelayInSeconds int
	FixedDelayInSeconds   int
	ErrorDelayInSeconds   int

	Enable bool
}

func (c *SyncCmdbFromBossConf) isEnable() bool {
	return c.Enable
}
func (c *SyncCmdbFromBossConf) getIntervalConfig() *ocron.IntervalConfig {
	return &ocron.IntervalConfig{
		InitialDelay: time.Duration(c.InitialDelayInSeconds) * time.Second,
		FixedDelay:   time.Duration(c.FixedDelayInSeconds) * time.Second,
		ErrorDelay:   time.Duration(c.ErrorDelayInSeconds) * time.Second,
	}
}
func (c *SyncCmdbFromBossConf) buildJob() ocron.Job {
	return &syncCmdbFromBoss{}
}
func (c *SyncCmdbFromBossConf) taskName() string {
	return "Cmdb.Sync"
}

func (c *SyncCmdbFromBossConf) String() string {
	return fmt.Sprintf("Synchronization of CMDB from boss: %#v", c)
}
