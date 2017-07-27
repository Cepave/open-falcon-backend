package sqlx

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Cepave/open-falcon-backend/common/db"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	or "github.com/Cepave/open-falcon-backend/common/runtime"
	"github.com/Cepave/open-falcon-backend/common/utils"
)

var buildError = utils.BuildErrorWithCaller

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
func (txExt *TxExt) SqlxTx() *sqlx.Tx {
	return (*sqlx.Tx)(txExt)
}
func (txExt *TxExt) QueryRowAndScan(query string, args []interface{}, dest []interface{}) {
	err := txExt.SqlxTx().QueryRow(query, args...).Scan(dest...)
	db.PanicIfError(buildError(err))
}

func (txExt *TxExt) QueryRowxAndMapScan(dest map[string]interface{}, query string, args ...interface{}) {
	row := txExt.SqlxTx().QueryRowx(query, args...)
	db.PanicIfError(buildError(row.MapScan(dest)))
}
func (txExt *TxExt) QueryRowxAndScan(query string, args []interface{}, dest ...interface{}) {
	row := txExt.SqlxTx().QueryRowx(query, args...)
	db.PanicIfError(buildError(row.Scan(dest...)))
}
func (txExt *TxExt) QueryRowxAndSliceScan(query string, args ...interface{}) []interface{} {
	row := txExt.SqlxTx().QueryRowx(query, args...)
	result, err := row.SliceScan()
	db.PanicIfError(buildError(err))

	return result
}
func (txExt *TxExt) QueryRowxAndStructScan(dest interface{}, query string, args ...interface{}) {
	row := txExt.SqlxTx().QueryRowx(query, args...)
	db.PanicIfError(buildError(row.StructScan(dest)))
}

func (txExt *TxExt) BindNamed(query string, arg interface{}) (string, []interface{}) {
	newQuery, newArgs, err := txExt.SqlxTx().BindNamed(query, arg)
	db.PanicIfError(buildError(err))

	return newQuery, newArgs
}
func (txExt *TxExt) NamedExec(query string, arg interface{}) sql.Result {
	result, err := txExt.SqlxTx().NamedExec(query, arg)
	db.PanicIfError(buildError(err))

	return result
}
func (txExt *TxExt) NamedQuery(query string, arg interface{}) *sqlx.Rows {
	rows, err := txExt.SqlxTx().NamedQuery(query, arg)
	db.PanicIfError(buildError(err))

	return rows
}

func (txExt *TxExt) Queryx(query string, args ...interface{}) *sqlx.Rows {
	rows, err := txExt.SqlxTx().Queryx(query, args...)
	db.PanicIfError(buildError(err))

	return rows
}
func (txExt *TxExt) QueryxExt(query string, args ...interface{}) *RowsExt {
	return ToRowsExt(txExt.Queryx(query, args...))
}
func (txExt *TxExt) QueryRowxExt(query string, args ...interface{}) *RowExt {
	return ToRowExt(txExt.SqlxTx().QueryRowx(query, args...))
}

func (txExt *TxExt) Get(dest interface{}, query string, args ...interface{}) {
	err := txExt.SqlxTx().Get(dest, query, args...)
	db.PanicIfError(buildError(err))
}
func (txExt *TxExt) GetOrNoRow(dest interface{}, query string, args ...interface{}) bool {
	err := txExt.SqlxTx().Get(dest, query, args...)
	if err == sql.ErrNoRows {
		return false
	}

	db.PanicIfError(buildError(err))
	return true
}

func (txExt *TxExt) Select(dest interface{}, query string, args ...interface{}) {
	err := txExt.SqlxTx().Select(dest, query, args...)
	db.PanicIfError(buildError(err))
}

func (txExt *TxExt) PrepareNamed(query string) *sqlx.NamedStmt {
	namedStmt, err := txExt.SqlxTx().PrepareNamed(query)
	db.PanicIfError(buildError(err))
	return namedStmt
}
func (txExt *TxExt) PrepareNamedExt(query string) *NamedStmtExt {
	defer utils.DeferCatchPanicWithCaller()()
	return ToNamedStmtExt(txExt.PrepareNamed(query))
}
func (txExt *TxExt) Preparex(query string) *sqlx.Stmt {
	stmt, err := txExt.SqlxTx().Preparex(query)
	db.PanicIfError(buildError(err))
	return stmt
}
func (txExt *TxExt) PreparexExt(query string) *StmtExt {
	defer utils.DeferCatchPanicWithCaller()()
	return ToStmtExt(txExt.Preparex(query))
}

func NewDbController(newSqlxDb *sqlx.DB) *DbController {
	return &DbController{
		sqlxDb: newSqlxDb,
	}
}

// Prepares transaction and feed it to callback function
func (ctrl *DbController) InTx(txCallback TxCallback) {
	defer utils.DeferCatchPanicWithCaller()()
	callerInfo := or.GetCallerInfo()
	tx := ctrl.sqlxDb.MustBegin()

	defer func() {
		p := recover()
		if p == nil {
			return
		}

		finalError := utils.BuildErrorWithCallerInfo(
			utils.SimpleErrorConverter(p), callerInfo,
		)

		rollbackError := tx.Rollback()
		if rollbackError != nil {
			finalError = utils.BuildErrorWithCallerInfo(
				fmt.Errorf("Transaction has error: %v. Rollback has error too: %v", finalError, rollbackError),
				callerInfo,
			)
		}

		db.PanicIfError(finalError)
	}()

	switch txCallback.InTx(tx) {
	case db.TxCommit:
		db.PanicIfError(tx.Commit())
	case db.TxRollback:
		db.PanicIfError(tx.Rollback())
	}
}

func (ctrl *DbController) SqlxDb() *sqlx.DB {
	return ctrl.sqlxDb
}
func (ctrl *DbController) BindNamed(query string, arg interface{}) (string, []interface{}) {
	r1, r2, err := ctrl.sqlxDb.BindNamed(query, arg)
	db.PanicIfError(buildError(err))

	return r1, r2
}
func (ctrl *DbController) NamedExec(query string, arg interface{}) sql.Result {
	result, err := ctrl.sqlxDb.NamedExec(query, arg)
	db.PanicIfError(buildError(err))

	return result
}
func (ctrl *DbController) NamedQuery(query string, arg interface{}) *sqlx.Rows {
	rows, err := ctrl.sqlxDb.NamedQuery(query, arg)
	db.PanicIfError(buildError(err))

	return rows
}
func (ctrl *DbController) Select(dest interface{}, query string, args ...interface{}) {
	err := ctrl.sqlxDb.Select(dest, query, args...)
	db.PanicIfError(buildError(err))
}

func (ctrl *DbController) Get(dest interface{}, query string, args ...interface{}) {
	err := ctrl.sqlxDb.Get(dest, query, args...)
	db.PanicIfError(buildError(err))
}
func (ctrl *DbController) GetOrNoRow(dest interface{}, query string, args ...interface{}) bool {
	err := ctrl.sqlxDb.Get(dest, query, args...)
	if err == sql.ErrNoRows {
		return false
	}

	db.PanicIfError(buildError(err))
	return true
}

func (ctrl *DbController) QueryRowxAndMapScan(dest map[string]interface{}, query string, args ...interface{}) {
	row := ctrl.sqlxDb.QueryRowx(query, args...)
	db.PanicIfError(buildError(row.MapScan(dest)))
}
func (ctrl *DbController) QueryRowxAndScan(query string, args []interface{}, dest ...interface{}) {
	row := ctrl.sqlxDb.QueryRowx(query, args...)
	db.PanicIfError(buildError(row.Scan(dest...)))
}
func (ctrl *DbController) QueryRowxAndSliceScan(query string, args ...interface{}) []interface{} {
	row := ctrl.sqlxDb.QueryRowx(query, args...)
	result, err := row.SliceScan()
	db.PanicIfError(buildError(err))

	return result
}
func (ctrl *DbController) QueryRowxAndStructScan(dest interface{}, query string, args ...interface{}) {
	row := ctrl.sqlxDb.QueryRowx(query, args...)
	db.PanicIfError(buildError(row.StructScan(dest)))
}

func (ctrl *DbController) Queryx(query string, args ...interface{}) *sqlx.Rows {
	rows, err := ctrl.sqlxDb.Queryx(query, args...)
	db.PanicIfError(buildError(err))

	return rows
}
func (ctrl *DbController) QueryxExt(query string, args ...interface{}) *RowsExt {
	return ToRowsExt(ctrl.Queryx(query, args...))
}
func (ctrl *DbController) QueryRowxExt(query string, args ...interface{}) *RowExt {
	defer utils.DeferCatchPanicWithCaller()()
	return ToRowExt(ctrl.sqlxDb.QueryRowx(query, args...))
}

func (ctrl *DbController) SelectWithFoundRows(txCallback TxCallback, paging *commonModel.Paging) {
	defer utils.DeferCatchPanicWithCaller()()
	finalFunc := func(sqlxTx *sqlx.Tx) db.TxFinale {
		defer utils.DeferCatchPanicWithCaller()()
		txFinale := txCallback.InTx(sqlxTx)

		var numOfRows int32
		ToTxExt(sqlxTx).Get(&numOfRows, "SELECT FOUND_ROWS()")
		paging.SetTotalCount(numOfRows)

		return txFinale
	}

	ctrl.InTx(TxCallbackFunc(finalFunc))
}

func (ctrl *DbController) PrepareNamed(query string) *sqlx.NamedStmt {
	stmt, err := ctrl.sqlxDb.PrepareNamed(query)
	db.PanicIfError(buildError(err))

	return stmt
}
func (ctrl *DbController) PrepareNamedExt(query string) *NamedStmtExt {
	defer utils.DeferCatchPanicWithCaller()()
	return ToNamedStmtExt(ctrl.PrepareNamed(query))
}
func (ctrl *DbController) Preparex(query string) *sqlx.Stmt {
	stmt, err := ctrl.sqlxDb.Preparex(query)
	db.PanicIfError(buildError(err))

	return stmt
}
func (ctrl *DbController) PreparexExt(query string) *StmtExt {
	defer utils.DeferCatchPanicWithCaller()()
	return ToStmtExt(ctrl.Preparex(query))
}

type RowExt sqlx.Row

func ToRowExt(row *sqlx.Row) *RowExt {
	return (*RowExt)(row)
}

func (r *RowExt) SqlxRow() *sqlx.Row {
	return (*sqlx.Row)(r)
}
func (r *RowExt) MapScanOrNoRow(dest map[string]interface{}) bool {
	err := r.SqlxRow().MapScan(dest)
	if err == sql.ErrNoRows {
		return false
	}

	db.PanicIfError(buildError(err))
	return true
}
func (r *RowExt) ScanOrNoRow(dest ...interface{}) bool {
	err := r.SqlxRow().Scan(dest...)
	if err == sql.ErrNoRows {
		return false
	}

	db.PanicIfError(buildError(err))
	return true
}
func (r *RowExt) SliceScanOrNoRow() ([]interface{}, bool) {
	result, err := r.SqlxRow().SliceScan()
	if err == sql.ErrNoRows {
		return result, false
	}

	db.PanicIfError(buildError(err))
	return result, true
}
func (r *RowExt) StructScanOrNoRow(dest interface{}) bool {
	err := r.SqlxRow().StructScan(dest)
	if err == sql.ErrNoRows {
		return false
	}

	db.PanicIfError(buildError(err))
	return true
}

type RowsExt sqlx.Rows

func ToRowsExt(row *sqlx.Rows) *RowsExt {
	return (*RowsExt)(row)
}

func (r *RowsExt) SqlxRows() *sqlx.Rows {
	return (*sqlx.Rows)(r)
}
func (r *RowsExt) MapScan(dest map[string]interface{}) {
	db.PanicIfError(buildError(
		r.SqlxRows().MapScan(dest),
	))
}
func (r *RowsExt) SliceScan() []interface{} {
	result, err := r.SqlxRows().SliceScan()
	db.PanicIfError(buildError(err))

	return result
}
func (r *RowsExt) StructScan(dest interface{}) {
	db.PanicIfError(buildError(
		r.SqlxRows().StructScan(dest),
	))
}

// Extensions for *sqlx.NamedStmt
type NamedStmtExt sqlx.NamedStmt

func ToNamedStmtExt(stmt *sqlx.NamedStmt) *NamedStmtExt {
	return (*NamedStmtExt)(stmt)
}

func (s *NamedStmtExt) SqlxNamedStmt() *sqlx.NamedStmt {
	return (*sqlx.NamedStmt)(s)

}
func (s *NamedStmtExt) Get(dest interface{}, arg interface{}) {
	err := s.SqlxNamedStmt().Get(dest, arg)
	db.PanicIfError(buildError(err))
}
func (s *NamedStmtExt) GetOrNoRow(dest interface{}, args interface{}) bool {
	err := s.SqlxNamedStmt().Get(dest, args)
	if err == sql.ErrNoRows {
		return false
	}

	db.PanicIfError(buildError(err))
	return true
}

func (s *NamedStmtExt) Select(dest interface{}, arg interface{}) {
	err := s.SqlxNamedStmt().Select(dest, arg)
	db.PanicIfError(buildError(err))
}

func (s *NamedStmtExt) Queryx(args interface{}) *sqlx.Rows {
	rows, err := s.SqlxNamedStmt().Queryx(args)
	db.PanicIfError(buildError(err))

	return rows
}
func (s *NamedStmtExt) QueryxExt(args interface{}) *RowsExt {
	return ToRowsExt(s.Queryx(args))
}
func (s *NamedStmtExt) QueryRowxExt(args interface{}) *RowExt {
	defer utils.DeferCatchPanicWithCaller()()
	return ToRowExt(s.SqlxNamedStmt().QueryRowx(args))
}

func (s *NamedStmtExt) QueryRowxAndMapScan(dest map[string]interface{}, args interface{}) {
	row := s.SqlxNamedStmt().QueryRowx(args)
	db.PanicIfError(buildError(row.MapScan(dest)))
}
func (s *NamedStmtExt) QueryRowxAndScan(args []interface{}, dest ...interface{}) {
	row := s.SqlxNamedStmt().QueryRowx(args)
	db.PanicIfError(buildError(row.Scan(dest...)))
}
func (s *NamedStmtExt) QueryRowxAndSliceScan(args interface{}) []interface{} {
	row := s.SqlxNamedStmt().QueryRowx(args)

	result, err := row.SliceScan()
	db.PanicIfError(buildError(err))

	return result
}
func (s *NamedStmtExt) QueryRowxAndStructScan(dest interface{}, args interface{}) {
	row := s.SqlxNamedStmt().QueryRowx(args)
	db.PanicIfError(buildError(row.StructScan(dest)))
}

// Extensions for *sqlx.Stmt
type StmtExt sqlx.Stmt

func ToStmtExt(stmt *sqlx.Stmt) *StmtExt {
	return (*StmtExt)(stmt)
}

func (s *StmtExt) SqlxStmt() *sqlx.Stmt {
	return (*sqlx.Stmt)(s)

}
func (s *StmtExt) Get(dest interface{}, args ...interface{}) {
	err := s.SqlxStmt().Get(dest, args...)
	db.PanicIfError(buildError(err))
}
func (s *StmtExt) GetOrNoRow(dest interface{}, args ...interface{}) bool {
	err := s.SqlxStmt().Get(dest, args...)
	if err == sql.ErrNoRows {
		return false
	}

	db.PanicIfError(buildError(err))
	return true
}

func (s *StmtExt) Select(dest interface{}, args ...interface{}) {
	err := s.SqlxStmt().Select(dest, args...)
	db.PanicIfError(buildError(err))
}

func (s *StmtExt) Queryx(args ...interface{}) *sqlx.Rows {
	rows, err := s.SqlxStmt().Queryx(args...)
	db.PanicIfError(buildError(err))

	return rows
}
func (s *StmtExt) QueryxExt(args ...interface{}) *RowsExt {
	return ToRowsExt(s.Queryx(args...))
}
func (s *StmtExt) QueryRowxExt(args ...interface{}) *RowExt {
	defer utils.DeferCatchPanicWithCaller()()
	return ToRowExt(s.SqlxStmt().QueryRowx(args...))
}

func (s *StmtExt) QueryRowxAndMapScan(dest map[string]interface{}, args ...interface{}) {
	row := s.SqlxStmt().QueryRowx(args...)
	db.PanicIfError(buildError(row.MapScan(dest)))
}
func (s *StmtExt) QueryRowxAndScan(args []interface{}, dest ...interface{}) {
	row := s.SqlxStmt().QueryRowx(args...)
	db.PanicIfError(buildError(row.Scan(dest...)))
}
func (s *StmtExt) QueryRowxAndSliceScan(args ...interface{}) []interface{} {
	row := s.SqlxStmt().QueryRowx(args...)

	result, err := row.SliceScan()
	db.PanicIfError(buildError(err))

	return result
}
func (s *StmtExt) QueryRowxAndStructScan(dest interface{}, args ...interface{}) {
	row := s.SqlxStmt().QueryRowx(args...)
	db.PanicIfError(buildError(row.StructScan(dest)))
}
