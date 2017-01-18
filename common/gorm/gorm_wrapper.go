package gorm

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/Cepave/open-falcon-backend/common/db"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
)

type GormDbExt struct {
	// Error converter of Orm
	ConvertError ErrorConverter

	gormDb *gorm.DB
}

// Callback function for transaction
type TxCallback interface {
	InTx(*gorm.DB) db.TxFinale
}

// the function object delegates the TxCallback interface
type TxCallbackFunc func(*gorm.DB) db.TxFinale
func (callbackFunc TxCallbackFunc) InTx(gormDB *gorm.DB) db.TxFinale {
	return callbackFunc(gormDB)
}

// Converter of error
type ErrorConverter func(error) error

// Raise panic if the error is not nil
func (self ErrorConverter) PanicIfDbError(gormDb *gorm.DB) {
	self.PanicIfError(gormDb.Error)
}
// Raise panic if the error is not nil
func (self ErrorConverter) PanicIfError(err error) {
	if err != nil {
		panic(self(err))
	}
}

func sameError(err error) error {
	return err
}

// Converts gorm.DB to GormDbExt
func ToGormDbExt(gormDb *gorm.DB) *GormDbExt {
	return &GormDbExt {
		ConvertError: sameError,
		gormDb: gormDb,
	}
}

// Raise panic if the Gorm has error
func (self *GormDbExt) PanicIfError() {
	self.ConvertError.PanicIfDbError(self.gormDb)
}

func (self *GormDbExt) IsRecordNotFound() bool {
	return self.gormDb.Error == gorm.ErrRecordNotFound
}

// Gets the foundValue if there is no error.
// Or gets notFoundValue if the error is gorm.ErrRecordNotFound
//
// This function also calls PanicIfError(
func (self *GormDbExt) IfRecordNotFound(foundValue interface{}, notFoundValue interface{}) interface{} {
	if self.IsRecordNotFound() {
		return notFoundValue
	}

	self.PanicIfError()
	return foundValue
}

// Iterate rows(and close it) with callback
func (self *GormDbExt) IterateRows(
	rowsCallback db.RowsCallback,
) {
	rows := self.Rows()
	defer rows.Close()

	for rows.Next() {
		if rowsCallback.NextRow(rows) == db.IterateStop {
			break;
		}
	}
}

// Same as Rows() with panic instead of returned error
func (self *GormDbExt) Rows() *sql.Rows {
	rows, err := self.gormDb.Rows()
	self.ConvertError.PanicIfError(err)

	return rows
}

// Executes gorm in transaction
func (self *GormDbExt) InTx(txCallback TxCallback) {
	txGormDb := self.gormDb.Begin()

	defer func() {
		p := recover()
		if p != nil {
			txGormDb = txGormDb.Rollback()
			if txGormDb.Error != nil{
				p = fmt.Errorf("Transaction has error: %v. Rollback has error too: %v", p, txGormDb.Error)
			}

			panic(p)
		}
	}()

	switch txCallback.InTx(txGormDb) {
	case db.TxCommit:
		self.ConvertError.PanicIfDbError(txGormDb.Commit())
	case db.TxRollback:
		self.ConvertError.PanicIfDbError(txGormDb.Rollback())
	}
}

// Selects the query by callback and perform "SELECT FOUND_ROWS()" to gets the total number of matched rows
func (self *GormDbExt) SelectWithFoundRows(txCallback TxCallback, paging *commonModel.Paging) {
	var finalFunc TxCallbackFunc = func(txGormDb *gorm.DB) db.TxFinale {
		txFinale := txCallback.InTx(txGormDb)

		var selectFoundRows = txGormDb.Raw("SELECT FOUND_ROWS()")
		err := selectFoundRows.Row().Scan(&paging.TotalCount)
		self.ConvertError.PanicIfError(err)

		return txFinale
	}

	self.InTx(finalFunc)
}

// Same as ScanRows with panic instead of returned error
func (self *GormDbExt) ScanRows(rows *sql.Rows, result interface{}) {
	err := self.gormDb.ScanRows(rows, result)
	self.ConvertError.PanicIfError(err)
}
