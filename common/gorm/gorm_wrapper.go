package gorm

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/Cepave/open-falcon-backend/common/db"
)

type GormDbExt struct {
	// Error converter of Orm
	ConvertError ErrorConverter

	gormDb *gorm.DB
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

// Same as ScanRows with panic instead of returned error
func (self *GormDbExt) ScanRows(rows *sql.Rows, result interface{}) {
	err := self.gormDb.ScanRows(rows, result)
	self.ConvertError.PanicIfError(err)
}
