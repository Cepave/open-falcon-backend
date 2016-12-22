package sqlx

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/Cepave/open-falcon-backend/common/db"
)

// The interface of transaction callback for sqlx package
type TxCallback interface {
	InTx(tx *sqlx.Tx) db.TxFinale
}

// The function object delegates the TxCallback interface
type TxCallbackFunc func(*sqlx.Tx) db.TxFinale
func (callbackFunc TxCallbackFunc) InTx(tx *sqlx.Tx) db.TxFinale {
	return callbackFunc(tx)
}

type DbController struct {
	sqlxDb *sqlx.DB
}

// Extentsion for *sqlx.Tx
type TxExt sqlx.Tx

func ToTxExt(tx *sqlx.Tx) *TxExt {
	return (*TxExt)(tx)
}

// Query row and scan the result with panic
func (txExt *TxExt) QueryRowAndScan(query string, args []interface{}, dest []interface{}) {
	tx := (*sqlx.Tx)(txExt)

	err := tx.QueryRow(query, args...).Scan(dest...)
	db.PanicIfError(err)
}

func (txExt *TxExt) QueryRowxAndMapScan(dest map[string]interface{}, query string, args ...interface{}) {
	tx := (*sqlx.Tx)(txExt)

	row := tx.QueryRowx(query, args...)
	db.PanicIfError(row.MapScan(dest))
}
func (txExt *TxExt) QueryRowxAndScan(query string, args []interface{}, dest ...interface{}) {
	tx := (*sqlx.Tx)(txExt)

	row := tx.QueryRowx(query, args...)
	db.PanicIfError(row.Scan(dest...))
}
func (txExt *TxExt) QueryRowxAndSliceScan(query string, args ...interface{}) []interface{} {
	tx := (*sqlx.Tx)(txExt)

	row := tx.QueryRowx(query, args...)
	result, err := row.SliceScan()
	db.PanicIfError(err)

	return result
}
func (txExt *TxExt) QueryRowxAndStructScan(dest interface{}, query string, args ...interface{}) {
	tx := (*sqlx.Tx)(txExt)

	row := tx.QueryRowx(query, args...)
	db.PanicIfError(row.StructScan(dest))
}

func (txExt *TxExt) QueryRowxExt(query string, args ...interface{}) *RowExt {
	tx := (*sqlx.Tx)(txExt)
	return ToRowExt(tx.QueryRowx(query, args...))
}

func (txExt *TxExt) NamedExec(query string, arg interface{}) sql.Result {
	tx := (*sqlx.Tx)(txExt)

	result, err := tx.NamedExec(query, arg)
	db.PanicIfError(err)

	return result
}
func (txExt *TxExt) NamedQuery(query string, arg interface{}) *sqlx.Rows {
	tx := (*sqlx.Tx)(txExt)

	rows, err := tx.NamedQuery(query, arg)
	db.PanicIfError(err)

	return rows
}

func (txExt *TxExt) Queryx(query string, args ...interface{}) *sqlx.Rows {
	tx := (*sqlx.Tx)(txExt)

	rows, err := tx.Queryx(query, args...)
	db.PanicIfError(err)

	return rows
}

func (txExt *TxExt) Get(dest interface{}, query string, args ...interface{}) {
	tx := (*sqlx.Tx)(txExt)

	err := tx.Get(dest, query, args...)
	db.PanicIfError(err)
}
func (txExt *TxExt) GetOrNoRow(dest interface{}, query string, args ...interface{}) bool {
	tx := (*sqlx.Tx)(txExt)

	err := tx.Get(dest, query, args...)
	if err == sql.ErrNoRows {
		return false
	}

	db.PanicIfError(err)
	return true
}

func (txExt *TxExt) Select(dest interface{}, query string, args ...interface{}) {
	tx := (*sqlx.Tx)(txExt)

	err := tx.Select(dest, query, args...)
	db.PanicIfError(err)
}

func (txExt *TxExt) BindNamed(query string, arg interface{}) (string, []interface{}) {
	tx := (*sqlx.Tx)(txExt)

	newQuery, newArgs, err := tx.BindNamed(query, arg)
	db.PanicIfError(err)

	return newQuery, newArgs
}

func NewDbController(newSqlxDb *sqlx.DB) *DbController {
	return &DbController{
		sqlxDb: newSqlxDb,
	}
}

// Prepares transaction and feed it to callback function
func (ctrl *DbController) InTx(txCallback TxCallback) {
	tx := ctrl.sqlxDb.MustBegin()

	defer func() {
		p := recover()
		if p == nil {
			return
		}

		rollbackError := tx.Rollback()
		if rollbackError != nil {
			p = fmt.Errorf("Transaction has error: %v. Rollback has error too: %v", p, rollbackError)
		}

		panic(p)
	}()

	switch txCallback.InTx(tx) {
	case db.TxCommit:
		db.PanicIfError(tx.Commit())
	case db.TxRollback:
		db.PanicIfError(tx.Rollback())
	}
}

func (ctrl *DbController) BindNamed(query string, arg interface{}) (string, []interface{}) {
	r1, r2, err := ctrl.sqlxDb.BindNamed(query, arg)
	db.PanicIfError(err)

	return r1, r2
}
func (ctrl *DbController) Get(dest interface{}, query string, args ...interface{}) {
	err := ctrl.sqlxDb.Get(dest, query, args...)
	db.PanicIfError(err)
}
func (ctrl *DbController) GetOrNoRow(dest interface{}, query string, args ...interface{}) bool {
	err := ctrl.sqlxDb.Get(dest, query, args...)
	if err == sql.ErrNoRows {
		return false
	}

	db.PanicIfError(err)
	return true
}
func (ctrl *DbController) NamedExec(query string, arg interface{}) sql.Result {
	result, err := ctrl.sqlxDb.NamedExec(query, arg)
	db.PanicIfError(err)

	return result
}
func (ctrl *DbController) NamedQuery(query string, arg interface{}) *sqlx.Rows {
	rows, err := ctrl.sqlxDb.NamedQuery(query, arg)
	db.PanicIfError(err)

	return rows
}
func (ctrl *DbController) PrepareNamed(query string) *sqlx.NamedStmt {
	stmt, err := ctrl.sqlxDb.PrepareNamed(query)
	db.PanicIfError(err)

	return stmt
}
func (ctrl *DbController) Preparex(query string) *sqlx.Stmt {
	stmt, err := ctrl.sqlxDb.Preparex(query)
	db.PanicIfError(err)

	return stmt
}
func (ctrl *DbController) Queryx(query string, args ...interface{}) *sqlx.Rows {
	rows, err := ctrl.sqlxDb.Queryx(query, args...)
	db.PanicIfError(err)

	return rows
}
func (ctrl *DbController) Select(dest interface{}, query string, args ...interface{}) {
	err := ctrl.sqlxDb.Select(dest, query, args...)
	db.PanicIfError(err)
}

func (ctrl *DbController) QueryRowxAndMapScan(dest map[string]interface{}, query string, args ...interface{}) {
	row := ctrl.sqlxDb.QueryRowx(query, args...)
	db.PanicIfError(row.MapScan(dest))
}
func (ctrl *DbController) QueryRowxAndScan(query string, args []interface{}, dest ...interface{}) {
	row := ctrl.sqlxDb.QueryRowx(query, args...)
	db.PanicIfError(row.Scan(dest...))
}
func (ctrl *DbController) QueryRowxAndSliceScan(query string, args ...interface{}) []interface{} {
	row := ctrl.sqlxDb.QueryRowx(query, args...)
	result, err := row.SliceScan()
	db.PanicIfError(err)

	return result
}
func (ctrl *DbController) QueryRowxAndStructScan(dest interface{}, query string, args ...interface{}) {
	row := ctrl.sqlxDb.QueryRowx(query, args...)
	db.PanicIfError(row.StructScan(dest))
}

func (ctrl *DbController) QueryRowxExt(query string, args ...interface{}) *RowExt {
	return ToRowExt(ctrl.sqlxDb.QueryRowx(query, args...))
}

type RowExt sqlx.Row
func ToRowExt(row *sqlx.Row) *RowExt {
	return (*RowExt)(row)
}

func (r *RowExt) MapScanOrNoRow(dest map[string]interface{}) bool {
	row := (*sqlx.Row)(r)

	err := row.MapScan(dest)
	if err == sql.ErrNoRows {
		return false
	}

	db.PanicIfError(err)
	return true
}
func (r *RowExt) ScanOrNoRow(dest ...interface{}) bool {
	row := (*sqlx.Row)(r)

	err := row.Scan(dest...)
	if err == sql.ErrNoRows {
		return false
	}

	db.PanicIfError(err)
	return true
}
func (r *RowExt) SliceScanOrNoRow() ([]interface{}, bool) {
	row := (*sqlx.Row)(r)

	result, err := row.SliceScan()
	if err == sql.ErrNoRows {
		return result, false
	}

	db.PanicIfError(err)
	return result, true
}
func (r *RowExt) StructScanOrNoRow(dest interface{}) bool {
	row := (*sqlx.Row)(r)

	err := row.StructScan(dest)
	if err == sql.ErrNoRows {
		return false
	}

	db.PanicIfError(err)
	return true
}

type RowsExt sqlx.Rows

func ToRowsExt(row *sqlx.Rows) *RowsExt {
	return (*RowsExt)(row)
}

func (r *RowsExt) MapScan(dest map[string]interface{}) {
	rows := (*sqlx.Rows)(r)
	db.PanicIfError(rows.MapScan(dest))
}
func (r *RowsExt) SliceScan() []interface{} {
	rows := (*sqlx.Rows)(r)

	result, err := rows.SliceScan()
	db.PanicIfError(err)

	return result
}
func (r *RowsExt) StructScan(dest interface{}) {
	rows := (*sqlx.Rows)(r)
	db.PanicIfError(rows.StructScan(dest))
}
