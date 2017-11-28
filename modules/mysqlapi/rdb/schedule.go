package rdb

import (
	"time"

	cdb "github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

func FreeLock(schedule *model.Schedule,
	endStatus model.TaskStatus, endMsg *string, endTime time.Time) {
	txProcessor := &txFreeLock{
		schedule: schedule,
		status:   byte(endStatus),
		message:  endMsg,
		endTime:  endTime,
	}
	DbFacade.SqlxDbCtrl.InTx(txProcessor)
}

type txFreeLock struct {
	schedule *model.Schedule
	endTime  time.Time
	status   byte
	message  *string

	lockTable    model.OwlSchedule
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
		if free.lockTable.IsLocked() &&
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

func AcquireLock(schedule *model.Schedule, now time.Time) error {
	txProcessor := &txAcquireLock{
		schedule:  schedule,
		timeNow:   now,
		lockError: nil,
	}
	DbFacade.SqlxDbCtrl.InTx(txProcessor)
	return txProcessor.lockError
}

type txAcquireLock struct {
	schedule  *model.Schedule
	lockError *model.UnableToLockSchedule

	timeNow   time.Time
	lockTable model.OwlSchedule
	logTable  model.OwlScheduleLog
}

func (ack *txAcquireLock) InTx(tx *sqlx.Tx) cdb.TxFinale {

	/**
	 * Lock table
	 */
	ack.selectOrInsertLock(tx)
	// The previous task is not timeout()
	if ack.lockTable.IsLocked() && ack.notTimeout(tx) {
		ack.lockError = &model.UnableToLockSchedule{
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
			"status":    model.RUN,
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
		ack.lockTable.Lock = byte(model.FREE)
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
