package rdb

import (
	"fmt"
	"time"

	cdb "github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

type Schedule struct {
	Name    string
	Timeout int
	Uuid    uuid.UUID
}

func (s *Schedule) GetUuidString() string {
	return s.Uuid.String()
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
	Lock           byte      `db:"sch_lock"`
	LastUpdateTime time.Time `db:"sch_modify_time"`
}

func (sch *OwlSchedule) isLocked() bool {
	return sch.Lock == byte(LOCKED)
}

type OwlScheduleLog struct {
	Uuid      cdb.DbUuid `db:"sl_uuid"`
	SchId     int        `db:"sl_sch_id"`
	StartTime time.Time  `db:"sl_start_time"`
	EndTime   *time.Time `db:"sl_end_time"`
	Timeout   int        `db:"sl_timeout"`
	Status    byte       `db:"sl_status"`
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

const (
	FREE LockStatus = iota
	LOCKED
)

type TaskStatus byte

const (
	DONE TaskStatus = iota
	RUN
	FAIL
	TIMEOUT
)

type ScheduleCallback func() error

func Execute(schedule *Schedule, callback ScheduleCallback) error {
	err := AcquireLock(schedule, time.Now())
	if err != nil {
		return err
	}

	go func() {
		var err error

		/**
		 * Free lock after callback is finished
		 */
		defer func() {
			var (
				errMsg    string
				endStatus TaskStatus
			)

			if p := recover(); p != nil {
				endStatus = FAIL
				errMsg = fmt.Sprint(p)
			} else if err != nil {
				endStatus = FAIL
				errMsg = err.Error()
			} else {
				endStatus = DONE
			}

			FreeLock(schedule, endStatus, errMsg, time.Now())
		}()
		// :~)

		err = callback()
	}()

	return nil
}

func FreeLock(schedule *Schedule, endStatus TaskStatus, endMsg string, endTime time.Time) {
	txProcessor := &txFreeLock{
		schedule: schedule,
		status:   byte(endStatus),
		message:  endMsg,
		endTime:  endTime,
	}
	DbFacade.SqlxDbCtrl.InTx(txProcessor)
}

type txFreeLock struct {
	schedule *Schedule
	endTime  time.Time
	status   byte
	message  string

	lockTable    OwlSchedule
	logStartTime time.Time
}

func (free *txFreeLock) InTx(tx *sqlx.Tx) cdb.TxFinale {

	/**
	 * Lock & fetch table
	 */
	sqlxExt.ToTxExt(tx).Get(&free.lockTable, `
		SELECT sch_lock, sch_modify_time
		FROM owl_schedule
		WHERE sch_name = ?
		FOR UPDATE
	`, free.schedule.Name)

	exist := sqlxExt.ToTxExt(tx).GetOrNoRow(&free.logStartTime, `
		SELECT sl_start_time
		FROM owl_schedule_log
		WHERE sl_uuid = ?
	`, free.schedule.Uuid)
	// :~)

	/**
	 * Update table
	 */
	if exist {
		_ = tx.MustExec(`
				UPDATE owl_schedule_log
				SET sl_end_time = ?
				    sl_status = ?
					sl_message = ?
				WHERE sl_uuid = ?
			`, free.endTime, free.status, free.message, free.schedule.Uuid)
		// Release lock iff it is held by this task
		if free.lockTable.isLocked() &&
			free.lockTable.LastUpdateTime.Equal(free.logStartTime) {
			_ = tx.MustExec(`
					UPDATE owl_schedule
					SET sch_lock = 0
						sch_modify_time = ?
					WHERE sch_name = ?
				`, free.endTime, free.schedule.Name)
		}
	}
	// :~)

	return cdb.TxCommit
}

func AcquireLock(schedule *Schedule, now time.Time) error {
	txProcessor := &txAcquireLock{
		schedule:  schedule,
		timeNow:   now,
		lockError: nil,
	}
	DbFacade.SqlxDbCtrl.InTx(txProcessor)
	return txProcessor.lockError
}

type txAcquireLock struct {
	schedule  *Schedule
	lockError *UnableToLockSchedule

	timeNow   time.Time
	lockTable OwlSchedule
	logTable  OwlScheduleLog
}

func (ack *txAcquireLock) InTx(tx *sqlx.Tx) cdb.TxFinale {

	/**
	 * Lock table
	 */
	ack.selectOrInsertLock(tx)
	// The previous task is not timeout()
	if ack.lockTable.isLocked() && ack.notTimeout(tx) {
		ack.lockError = &UnableToLockSchedule{
			LastStartTime: ack.logTable.StartTime,
			AcquiredTime:  ack.timeNow,
			Timeout:       ack.logTable.Timeout,
		}
		return cdb.TxCommit
	}

	ack.updateLockByName(tx)
	// :~)

	/**
	 * Log table
	 */
	generatedUuid := uuid.NewV4()
	_ = sqlxExt.ToTxExt(tx).NamedExec(`
			INSERT INTO owl_schedule_log(
				sl_uuid, sl_sch_id,
				sl_start_time, sl_timeout, sl_status
			)
			VALUES (:uuid, :schid, :starttime, :timeout, :status)
		`,
		map[string]interface{}{
			"uuid":      cdb.DbUuid(generatedUuid),
			"schid":     ack.lockTable.Id,
			"starttime": ack.timeNow,
			"timeout":   ack.schedule.Timeout,
			"status":    RUN,
		},
	)
	ack.schedule.Uuid = generatedUuid
	// :~)

	return cdb.TxCommit
}

func (ack *txAcquireLock) selectOrInsertLock(tx *sqlx.Tx) {
	name := ack.schedule.Name
	exist := sqlxExt.ToTxExt(tx).GetOrNoRow(&ack.lockTable, `
		SELECT sch_id, sch_lock
		FROM owl_schedule
		WHERE sch_name = ?
		FOR UPDATE
	`, name)

	if !exist {
		r := tx.MustExec(`
			INSERT INTO owl_schedule(
				sch_name,
				sch_lock, sch_modify_time
			)
			VALUES (?, 0, ?)
		`, name, ack.timeNow)
		ack.lockTable.Id = int(cdb.ToResultExt(r).LastInsertId())
		ack.lockTable.Lock = byte(FREE)
	}
}

func (ack *txAcquireLock) updateLockByName(tx *sqlx.Tx) {
	_ = tx.MustExec(`
		UPDATE owl_schedule
		SET sch_lock = 1,
			sch_modify_time = ?
		WHERE sch_name = ?
	`, ack.timeNow, ack.schedule.Name)
}

func (ack *txAcquireLock) notTimeout(tx *sqlx.Tx) bool {
	ret := &ack.logTable
	exist := sqlxExt.ToTxExt(tx).GetOrNoRow(ret, `
		SELECT sl.sl_start_time, sl.sl_timeout
		FROM owl_schedule_log sl
		WHERE sl.sl_sch_id = ?
		ORDER BY sl.sl_start_time DESC
		LIMIT 1
	`, ack.lockTable.Id)

	// Check timeout iff row exists
	return exist && (ack.timeNow.Sub(ret.StartTime) <= time.Duration(ret.Timeout)*time.Second)
}
