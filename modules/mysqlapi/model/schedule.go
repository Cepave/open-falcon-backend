package model

import (
	"fmt"
	"time"
	"database/sql/driver"

	cdb "github.com/Cepave/open-falcon-backend/common/db"
	uuid "github.com/satori/go.uuid"
)

type Schedule struct {
	Name    string
	Timeout int
	Uuid    uuid.UUID
}

func (s *Schedule) GetUuidString() string {
	return s.Uuid.String()
}

func (s *Schedule) String() string {
	return fmt.Sprintf("Schedule<Name: %s, Timeout: %d, UUID: %s>",
		s.Name, s.Timeout, s.Uuid)
}

func NewSchedule(name string, timeout int) *Schedule {
	return &Schedule{
		Name:    name,
		Timeout: timeout,
	}
}

type OwlSchedule struct {
	Id             int       `db:"sch_id"`
	Name           string    `db:"sch_name"`
	Lock           LockStatus      `db:"sch_lock"`
	LastUpdateTime time.Time `db:"sch_modify_time"`
}

func (sch *OwlSchedule) IsLocked() bool {
	return sch.Lock == LOCKED
}

type OwlScheduleLog struct {
	Uuid      cdb.DbUuid `db:"sl_uuid"`
	SchId     int        `db:"sl_sch_id"`
	StartTime time.Time  `db:"sl_start_time"`
	EndTime   *time.Time `db:"sl_end_time"`
	Timeout   int        `db:"sl_timeout"`
	Status    TaskStatus       `db:"sl_status"`
	Message   *string    `db:"sl_message"`
}

type UnableToLockSchedule struct {
	AcquiredTime  time.Time
	LastStartTime time.Time
	Timeout       int
}

func (t *UnableToLockSchedule) Error() string {
	return fmt.Sprintf("Unable to lock schedule error: period between %s and %s should longer than %ds",
		t.LastStartTime.Format(time.RFC3339),
		t.AcquiredTime.Format(time.RFC3339),
		t.Timeout)
}


type LockStatus byte
func (s *LockStatus) Scan(src interface{}) error {
	*s = LockStatus(src.(int64))
	return nil
}
func (s LockStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

const (
	FREE LockStatus = iota
	LOCKED
)

type TaskStatus byte

func (s *TaskStatus) Scan(src interface{}) error {
	*s = TaskStatus(src.(int64))
	return nil
}
func (s TaskStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

const (
	DONE TaskStatus = iota
	RUN
	FAIL
	TIMEOUT
)
