package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	cdb "github.com/Cepave/open-falcon-backend/common/db"
)

type Schedule struct {
	Name    string
	Timeout int32
}

func (s *Schedule) String() string {
	return fmt.Sprintf("Schedule<Name: %s, Timeout: %d",
		s.Name, s.Timeout)
}

func NewSchedule(name string, timeout int) *Schedule {
	return &Schedule{
		Name:    name,
		Timeout: int32(timeout),
	}
}

type OwlSchedule struct {
	Id             int32      `db:"sch_id"`
	Name           string     `db:"sch_name"`
	Lock           LockStatus `db:"sch_lock"`
	LastUpdateTime time.Time  `db:"sch_modify_time"`
}

func (sch *OwlSchedule) IsLocked() bool {
	return sch.Lock == LOCKED
}

type OwlScheduleLog struct {
	Uuid      cdb.DbUuid `db:"sl_uuid"`
	SchId     int32      `db:"sl_sch_id"`
	StartTime time.Time  `db:"sl_start_time"`
	EndTime   cdb.DbTime `db:"sl_end_time"`
	Timeout   int32      `db:"sl_timeout"`
	Status    TaskStatus `db:"sl_status"`
	Message   *string    `db:"sl_message"`
}

func (log *OwlScheduleLog) IsTimeout(checkedTime time.Time) bool {
	return checkedTime.After(log.StartTime.Add(time.Duration(log.Timeout) * time.Second))
}
func (log *OwlScheduleLog) GetUuidString() string {
	return log.Uuid.ToUuid().String()
}

type UnableToLockSchedule struct {
	AcquiredTime  time.Time
	LastStartTime time.Time
	Timeout       int32
}

func (t *UnableToLockSchedule) Error() string {
	return fmt.Sprintf(
		"Unable to lock schedule error. Timeout[%d]. Period start time:[%s]. Acquired time:[%s]",
		t.Timeout,
		t.LastStartTime.Format(time.RFC3339),
		t.AcquiredTime.Format(time.RFC3339),
	)
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
