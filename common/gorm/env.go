package gorm

import (
	"github.com/Cepave/open-falcon-backend/common/db"
	"github.com/jinzhu/gorm"
)

// This function converts the error to default database error
//
// See ToGormDbExt
var DefaultGormErrorConverter ErrorConverter = func(err error) error {
	return db.NewDatabaseError(err)
}

// Converts gormDb to GormDbExt with convertion of DbError
func ToDefaultGormDbExt(gormDb *gorm.DB) *GormDbExt {
	gormDbExt := ToGormDbExt(gormDb)
	gormDbExt.ConvertError = DefaultGormErrorConverter
	return gormDbExt
}
