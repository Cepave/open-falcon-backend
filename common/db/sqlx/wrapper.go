package sqlx

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/Cepave/open-falcon-backend/common/db"
)

// The interface of transaction callback for sqlx package
type TxCallback interface {
	InTx(tx *sqlx.Tx)
}

// The function object delegates the TxCallback interface
type TxCallbackFunc func(*sqlx.Tx)
func (callbackFunc TxCallbackFunc) InTx(tx *sqlx.Tx) {
	callbackFunc(tx)
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

	txCallback.InTx(tx)

	tx.Commit()
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
