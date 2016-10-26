package db

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
var DbFacade = &commonDb.DbFacade{}

// Initialize the resource for RDB
func Init() {
	err := DbInit(
		&commonDb.DbConfig {
			Dsn: g.Config().Database,
			MaxIdle: g.Config().MaxIdle,
		},
	)

	if err != nil {
		log.Fatalln(err)
	}
}

// Initialize the resource for RDB
func Release() {
	DbFacade.Release()
	DB = DbFacade.SqlDb
}

// This function converts the error to default database error
//
// See ToGormDbExt
var DefaultGormErrorConverter gormExt.ErrorConverter = func(err error) error {
	return commonDb.NewDatabaseError(err)
}

// Converts gormDb to GormDbExt with convertion of DbError
func ToGormDbExt(gormDb *gorm.DB) *gormExt.GormDbExt {
	gormDbExt := gormExt.ToGormDbExt(gormDb)
	gormDbExt.ConvertError = DefaultGormErrorConverter
	return gormDbExt
}

func DbInit(dbConfig *commonDb.DbConfig) (err error) {
	err = DbFacade.Open(dbConfig)
	if err != nil {
		return
	}

	DB = DbFacade.SqlDb

	return
}

// Convenient IoC for transaction processing
func inTx(txCallback func(tx *sql.Tx) error) (err error) {
	var tx *sql.Tx

	if tx, err = DbFacade.SqlDb.Begin(); err != nil {
		return
	}

	/**
	 * The transaction result by whether or not the callback has error
	 */
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	// :~)

	err = txCallback(tx)

	return
}
