package rdb

import (
	"database/sql"
	"time"

	cdb "github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

func FreeLock(
	scheduleLog *model.OwlScheduleLog,
	endStatus model.TaskStatus, endMsg string,
	endTime time.Time,
) {
	txProcessor := &txFreeLock{
		scheduleLog: scheduleLog,
		status:   endStatus,
		message:  endMsg,
		endTime:  endTime,
	}
	DbFacade.SqlxDbCtrl.InTx(txProcessor)
}

type txFreeLock struct {
	scheduleLog *model.OwlScheduleLog
	endTime  time.Time
	status   model.TaskStatus
	message  string
}

func (self *txFreeLock) InTx(tx *sqlx.Tx) cdb.TxFinale {
	/**
	 * Release the lock directly rather than check the lock holder
	 */
	nullableMessage := sql.NullString{ Valid: false }
	if self.status == model.FAIL {
		nullableMessage.String = self.message
		nullableMessage.Valid = true
	}

	uuid := cdb.DbUuid(self.scheduleLog.Uuid)
	tx.MustExec(
		`
		UPDATE owl_schedule_log
		SET sl_end_time = ?,
			sl_status = ?,
			sl_message = ?
		WHERE sl_uuid = ?
		`,
		self.endTime, self.status,
		nullableMessage, uuid,
	)
	tx.MustExec(
		`
		UPDATE owl_schedule
		SET sch_lock = 0,
			sch_modify_time = ?
		WHERE sch_id = ?
		`,
		self.endTime, self.scheduleLog.SchId,
	)
	// :~)

	return cdb.TxCommit
}

func AcquireLock(schedule *model.Schedule, startTime time.Time) (*model.OwlScheduleLog, error) {
	txProcessor := &txAcquireLock{
		schedule:  schedule,
		startTime: startTime,
	}

	DbFacade.SqlxDbCtrl.InTx(txProcessor)

	return txProcessor.scheduleLog, txProcessor.lockError
}

type txAcquireLock struct {
	schedule  *model.Schedule
	startTime time.Time

	scheduleLog *model.OwlScheduleLog
	lockError error
}

func (ack *txAcquireLock) InTx(tx *sqlx.Tx) cdb.TxFinale {
	/**
	 * Lock table
	 */
	scheduleData := ack.selectOrInsertLock(tx, ack.schedule.Name)

	// Builds error if the previous task is locked and is not timeout
	if scheduleData.IsLocked() {
		if existingLog := ack.getExistingLog(tx, scheduleData.Id);
			existingLog != nil && !existingLog.IsTimeout(ack.startTime) {
			ack.lockError = &model.UnableToLockSchedule{
				LastStartTime: existingLog.StartTime,
				AcquiredTime:  ack.startTime,
				Timeout:       existingLog.Timeout,
			}
			return cdb.TxCommit
		}
	}
	// :~)

	newLog := &model.OwlScheduleLog{
		Uuid: cdb.DbUuid(uuid.NewV4()),
		SchId: scheduleData.Id,
		StartTime: ack.startTime,
		Timeout: ack.schedule.Timeout,
		Status: model.RUN,
	}

	/**
	 * Log table
	 */
	ack.updateLockByName(tx, ack.schedule.Name)
	sqlxExt.ToTxExt(tx).NamedExec(
		`
		INSERT INTO owl_schedule_log(
			sl_uuid, sl_sch_id,
			sl_start_time, sl_timeout, sl_status
		)
		VALUES (:uuid, :schid, :starttime, :timeout, :status)
		`,
		map[string]interface{}{
			"uuid":      newLog.Uuid,
			"schid":     newLog.SchId,
			"starttime": newLog.StartTime,
			"timeout":   newLog.Timeout,
			"status":    newLog.Status,
		},
	)
	// :~)

	ack.scheduleLog = newLog

	return cdb.TxCommit
}

func (ack *txAcquireLock) selectOrInsertLock(tx *sqlx.Tx, name string) *model.OwlSchedule {
	scheduleData := &model.OwlSchedule{}

	existing := sqlxExt.ToTxExt(tx).GetOrNoRow(
		scheduleData,
		`
		SELECT *
		FROM owl_schedule
		WHERE sch_name = ?
		FOR UPDATE
		`,
		name,
	)

	if existing {
		return scheduleData
	}

	/**
	 * Re-read data after insertion
	 */
	tx.MustExec(
		`
		INSERT INTO owl_schedule(
			sch_lock, sch_name, sch_modify_time
		)
		VALUES (0, ?, ?)
		`,
		name, ack.startTime,
	)

	return ack.selectOrInsertLock(tx, name)
	// :~)
}

func (ack *txAcquireLock) updateLockByName(tx *sqlx.Tx, name string) {
	tx.MustExec(
		`
		UPDATE owl_schedule
		SET sch_lock = 1,
			sch_modify_time = ?
		WHERE sch_name = ?
		`,
		ack.startTime, name,
	)
}

func (ack *txAcquireLock) getExistingLog(tx *sqlx.Tx, scheduleId int32) *model.OwlScheduleLog {
	logRecord := &model.OwlScheduleLog{}
	existing := sqlxExt.ToTxExt(tx).GetOrNoRow(
		logRecord,
		`
		SELECT sl.sl_start_time, sl.sl_timeout
		FROM owl_schedule_log sl
		WHERE sl.sl_sch_id = ?
		ORDER BY sl.sl_start_time DESC
		LIMIT 1
		`,
		scheduleId,
	)

	if !existing {
		return nil
	}

	return logRecord
}
