package db

import (
	"database/sql"
	"fmt"

	or "github.com/Cepave/open-falcon-backend/common/runtime"
	"github.com/Cepave/open-falcon-backend/common/utils"
)

type TxFinale byte

const (
	TxCommit   TxFinale = 1
	TxRollback TxFinale = 2
)

// Configuration of database
type DbConfig struct {
	Dsn     string
	MaxIdle int
}

func (config *DbConfig) String() string {
	return fmt.Sprintf("DSN: [%s]. Max Idle: [%d]", config.Dsn, config.MaxIdle)
}

// The main functions of this file is to gives IoC(Inverse of Control) of database(RDB) objects.
//
// For exception handling, all callback method should use panic() or log.Panicf() to release the error object.
//
// You may use PanicIfError to ease your process of Error object.

// Main controller of database
type DbController struct {
	dbObject      *sql.DB
	panicHandlers []utils.PanicHandler
}

// The interface of DB callback for sql package
type DbCallback interface {
	OnDb(db *sql.DB)
}

// The function object delegates the DbCallback interface
type DbCallbackFunc func(*sql.DB)

func (f DbCallbackFunc) OnDb(db *sql.DB) {
	f(db)
}

// The interface of rows callback for sql package
type RowsCallback interface {
	NextRow(row *sql.Rows) IterateControl
}

// The function object delegates the RowsCallback interface
type RowsCallbackFunc func(*sql.Rows) IterateControl

func (callbackFunc RowsCallbackFunc) NextRow(rows *sql.Rows) IterateControl {
	return callbackFunc(rows)
}

// The interface of row callback for sql package
type RowCallback interface {
	ResultRow(row *sql.Row)
}

// The function object delegates the RowCallback interface
type RowCallbackFunc func(*sql.Row)

func (callbackFunc RowCallbackFunc) ResultRow(row *sql.Row) {
	callbackFunc(row)
}

// The interface of transaction callback for sql package
type TxCallback interface {
	InTx(tx *sql.Tx) TxFinale
}

// The function object delegates the TxCallback interface
type TxCallbackFunc func(*sql.Tx) TxFinale

func (callbackFunc TxCallbackFunc) InTx(tx *sql.Tx) TxFinale {
	return callbackFunc(tx)
}

// BuildTxForSqls builds function for exeuction of multiple SQLs
func BuildTxForSqls(queries ...string) TxCallback {
	return TxCallbackFunc(func(tx *sql.Tx) TxFinale {
		txExt := ToTxExt(tx)

		for _, v := range queries {
			txExt.Exec(v)
		}

		return TxCommit
	})
}

// Executes callbacks in transaction if the boot callback has true value
type ExecuteIfByTx interface {
	// First calling of database for boolean result
	BootCallback(tx *sql.Tx) bool
	// If the boot callback has true result, this callback would get called
	IfTrue(tx *sql.Tx)
}

// Extension for sql.Rows
type RowsExt sql.Rows

// converts the sql.Rows to RowsExt
func ToRowsExt(rows *sql.Rows) *RowsExt {
	return ((*RowsExt)(rows))
}

// Gets columns, with panic instead of returned error
func (rowsExt *RowsExt) Columns() []string {
	columns, err := ((*sql.Rows)(rowsExt)).Columns()
	PanicIfError(utils.BuildErrorWithCaller(err))

	return columns
}

// Scans the values of row into variables, with panic instead of returned error
func (rowsExt *RowsExt) Scan(dest ...interface{}) {
	err := ((*sql.Rows)(rowsExt)).Scan(dest...)
	PanicIfError(utils.BuildErrorWithCaller(err))
}

// Extension for sql.Row
type RowExt sql.Row

// Converts the sql.Row to RowExt
func ToRowExt(row *sql.Row) *RowExt {
	return ((*RowExt)(row))
}

// Scans the values of row into variables, with panic instead of returned error
func (rowExt *RowExt) Scan(dest ...interface{}) {
	err := ((*sql.Row)(rowExt)).Scan(dest...)
	PanicIfError(utils.BuildErrorWithCaller(err))
}

// Extension for sql.Stmt
type StmtExt sql.Stmt

// Converts sql.Stmt to StmtExt
func ToStmtExt(stmt *sql.Stmt) *StmtExt {
	return ((*StmtExt)(stmt))
}

// Exec with panic instead of error
func (stmtExt *StmtExt) Exec(args ...interface{}) sql.Result {
	result, err := ((*sql.Stmt)(stmtExt)).Exec(args...)
	PanicIfError(utils.BuildErrorWithCaller(err))

	return result
}

// Query with panic instead of error
func (stmtExt *StmtExt) Query(args ...interface{}) *sql.Rows {
	rows, err := ((*sql.Stmt)(stmtExt)).Query(args...)
	PanicIfError(utils.BuildErrorWithCaller(err))

	return rows
}

// Extnesion for sql.Tx
type TxExt sql.Tx

// Converts sql.Tx to TxExt
func ToTxExt(tx *sql.Tx) *TxExt {
	return ((*TxExt)(tx))
}

// Commit with panic instead of returned error
func (txExt *TxExt) Commit() {
	err := ((*sql.Tx)(txExt)).Commit()
	PanicIfError(utils.BuildErrorWithCaller(err))
}

// Commit with panic instead of returned error
func (txExt *TxExt) Exec(query string, args ...interface{}) sql.Result {
	result, err := ((*sql.Tx)(txExt)).Exec(query, args...)
	PanicIfError(utils.BuildErrorWithCaller(err))

	return result
}

// Prepare with panic instead of returned error
func (txExt *TxExt) Prepare(query string) *sql.Stmt {
	stmt, err := ((*sql.Tx)(txExt)).Prepare(query)
	PanicIfError(utils.BuildErrorWithCaller(err))

	return stmt
}

// Query with panic instead of returned error
func (txExt *TxExt) Query(query string, args ...interface{}) *sql.Rows {
	rows, err := ((*sql.Tx)(txExt)).Query(query)
	PanicIfError(utils.BuildErrorWithCaller(err))

	return rows
}

// Rollback with panic instead of returned error
func (txExt *TxExt) Rollback() {
	err := ((*sql.Tx)(txExt)).Rollback()
	PanicIfError(utils.BuildErrorWithCaller(err))
}

// Extension for sql.Result
type ResultExt struct {
	sqlResult sql.Result
}

// Converts sql.Result to ResultExt
func ToResultExt(result sql.Result) *ResultExt {
	return &ResultExt{result}
}

// Gets last id of insert with panic instead of returned error
func (resultExt *ResultExt) LastInsertId() int64 {
	insertId, err := resultExt.sqlResult.LastInsertId()
	PanicIfError(utils.BuildErrorWithCaller(err))

	return insertId
}

// Gets last number of affected rows with panic instead of returned error
func (resultExt *ResultExt) RowsAffected() int64 {
	numberOfRowsAffected, err := resultExt.sqlResult.RowsAffected()
	PanicIfError(utils.BuildErrorWithCaller(err))

	return numberOfRowsAffected
}

// The control of iterating
type IterateControl byte

const (
	IterateContinue = IterateControl(1)
	IterateStop     = IterateControl(0)
)

// Initialize a controller for database
//
// Without RegisterPanicHandler() any PanicHandler,
// The raised panic would be re-paniced.
func NewDbController(newDbObject *sql.DB) *DbController {
	if newDbObject == nil {
		PanicIfError(utils.BuildErrorWithCaller(
			fmt.Errorf("Need viable DB object(sql.DB)"),
		))
	}

	return &DbController{
		dbObject:      newDbObject,
		panicHandlers: make([]utils.PanicHandler, 0),
	}
}

// Registers a handler while a panic is raised
//
// This object may register multiple handlers for panic
func (dbController *DbController) RegisterPanicHandler(panicHandler utils.PanicHandler) {
	dbController.panicHandlers = append(dbController.panicHandlers, panicHandler)
}

// Operate on database
func (dbController *DbController) OperateOnDb(dbCallback DbCallback) {
	dbController.needInitializedOrPanic()
	defer dbController.handlePanic()
	defer utils.DeferCatchPanicWithCaller()()

	dbCallback.OnDb(dbController.dbObject)
}

// Executes the query string or panic
func (dbController *DbController) Exec(query string, args ...interface{}) sql.Result {
	callerInfo := or.GetCallerInfo()

	var finalResult sql.Result
	var dbFunc DbCallbackFunc = func(db *sql.DB) {
		r, err := db.Exec(query, args...)
		PanicIfError(utils.BuildErrorWithCallerInfo(err, callerInfo))

		finalResult = r
	}

	dbController.OperateOnDb(dbFunc)
	return finalResult
}

// Query for rows and get called of rows with Next()
func (dbController *DbController) QueryForRows(
	rowsCallback RowsCallback,
	sqlQuery string, args ...interface{},
) (numberOfRows uint) {
	defer utils.DeferCatchPanicWithCaller()()

	var dbFunc DbCallbackFunc = func(db *sql.DB) {
		rows, err := db.Query(
			sqlQuery, args...,
		)

		if err != nil {
			err := fmt.Errorf(
				"Query SQL with exception: %v. SQL: \"%s\" Params: %#v",
				err, sqlQuery, args,
			)
			PanicIfError(utils.BuildErrorWithCaller(err))
		}

		defer rows.Close()
		for rows.Next() {
			numberOfRows++

			if rowsCallback.NextRow(rows) == IterateStop {
				break
			}
		}
	}

	dbController.OperateOnDb(dbFunc)

	return
}

// Query for a row and get called if the query is not failed
func (dbController *DbController) QueryForRow(
	rowCallback RowCallback,
	sqlQuery string, args ...interface{},
) {
	defer utils.DeferCatchPanicWithCaller()()

	var dbFunc DbCallbackFunc = func(db *sql.DB) {
		row := db.QueryRow(
			sqlQuery, args...,
		)

		rowCallback.ResultRow(row)
	}

	dbController.OperateOnDb(dbFunc)
}

// Executes in transaction.
//
// This method would commit the transaction if there is no raised panic,
// rollback it otherwise.
func (dbController *DbController) InTx(txCallback TxCallback) {
	defer utils.DeferCatchPanicWithCaller()()

	var dbFunc DbCallbackFunc = func(db *sql.DB) {
		callerInfo := or.GetCallerInfo()

		tx, err := db.Begin()
		PanicIfError(utils.BuildErrorWithCallerInfo(err, callerInfo))

		/**
		 * Rollback the transaction when panic is raised
		 */
		defer func() {
			p := recover()
			if p == nil {
				return
			}

			var finalError = utils.BuildErrorWithCallerInfo(
				utils.SimpleErrorConverter(p), callerInfo,
			)

			/**
			 * Rollback the transaction
			 */
			rollbackError := tx.Rollback()
			if rollbackError != nil {
				finalError = utils.BuildErrorWithCallerInfo(
					fmt.Errorf("Rollback has error: %v. Cause Error: %v", rollbackError, finalError), callerInfo,
				)
			}
			// :~)

			PanicIfError(finalError)
		}()
		// :~)

		txExt := ToTxExt(tx)
		switch txCallback.InTx(tx) {
		case TxCommit:
			txExt.Commit()
		case TxRollback:
			txExt.Rollback()
		}
	}

	dbController.OperateOnDb(dbFunc)
}

// Executes the complex statement in transaction
func (dbController *DbController) InTxForIf(ifCallbacks ExecuteIfByTx) {
	defer utils.DeferCatchPanicWithCaller()()

	var txFunc TxCallbackFunc = func(tx *sql.Tx) TxFinale {
		if ifCallbacks.BootCallback(tx) {
			ifCallbacks.IfTrue(tx)
		}

		return TxCommit
	}

	dbController.InTx(txFunc)
}

// Executes in transaction
func (dbController *DbController) ExecQueriesInTx(queries ...string) {
	defer utils.DeferCatchPanicWithCaller()()
	dbController.InTx(BuildTxForSqls(queries...))
}

// Releases the database object under this object
//
// As of service application(web, daemon...), this method is rarely get called
func (dbController *DbController) Release() {
	dbController.needInitializedOrPanic()
	defer dbController.handlePanic()

	err := dbController.dbObject.Close()

	if err != nil {
		PanicIfError(utils.BuildErrorWithCaller(
			fmt.Errorf("Release database connection error. %v", err),
		))
	}

	dbController.dbObject = nil
}

func (dbController *DbController) needInitializedOrPanic() {
	if dbController.dbObject != nil {
		return
	}

	PanicIfError(utils.BuildErrorWithCallerInfo(
		fmt.Errorf("The controller is not initialized"),
		or.GetCallerInfoWithDepth(1),
	))
}

func (dbController *DbController) handlePanic() {
	p := recover()
	if p == nil {
		return
	}

	if len(dbController.panicHandlers) == 0 {
		panic(p)
	}

	for _, handler := range dbController.panicHandlers {
		handler(p)
	}
}
