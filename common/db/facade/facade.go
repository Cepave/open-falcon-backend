package facade

import (
	"database/sql"
	"fmt"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	commonSqlx "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
)

// Gives facade interface supporting multiple object of db
//
// This facade supports:
// 	gorm - github.com/jinzhu/gorm
// 	sqlx - github.com/Cepave/open-falcon-backend/common/db/sqlx
// 	database/sql.DB
// 	dbCtrl
//
// Release resources
//
// In order to release resources in solid way, this facade provides "SetReleaseCallback(func())" to
// register callback function which gets called before this object releases the connections of database.
type DbFacade struct {
	SqlDb      *sql.DB
	SqlDbCtrl  *commonDb.DbController
	GormDb     *gorm.DB
	SqlxDb     *sqlx.DB
	SqlxDbCtrl *commonSqlx.DbController

	dbConfig commonDb.DbConfig

	releaseCallback func()

	initialized bool
}

// Open this facade with ping()
func (facade *DbFacade) Open(dbConfig *commonDb.DbConfig) (err error) {
	if facade.initialized {
		return
	}

	facade.dbConfig = *dbConfig

	/**
	 * Initialize Gorm(It would call ping())
	 */
	facade.GormDb, err = gorm.Open("mysql", dbConfig.Dsn)
	if err != nil {
		err = fmt.Errorf("Open Gorm error: %v", err)
	}
	// :~)

	/**
	 * Use the sql.DB object from Gorm and ping
	 */
	facade.SqlDb = facade.GormDb.DB()
	facade.SqlDb.SetMaxIdleConns(dbConfig.MaxIdle)
	// :~)

	facade.SqlxDb = sqlx.NewDb(facade.SqlDb, "mysql")
	facade.SqlxDbCtrl = commonSqlx.NewDbController(facade.SqlxDb)

	facade.SqlDbCtrl = commonDb.NewDbController(facade.SqlDb)
	facade.initialized = true

	return
}

// Gets the configuration of database.
//
// Warning: the information of password is revealed
func (facade *DbFacade) GetDbConfig() *commonDb.DbConfig {
	newConfig := facade.dbConfig
	return &newConfig
}

// Sets the callback used before releasing connections
func (facade *DbFacade) SetReleaseCallback(callback func()) {
	facade.releaseCallback = callback
}

// Close the database, release the resources
func (facade *DbFacade) Release() {
	if !facade.initialized {
		return
	}

	if facade.releaseCallback != nil {
		facade.releaseCallback()
		facade.releaseCallback = nil
	}

	facade.GormDb.Close()

	facade.SqlDb = nil
	facade.SqlDbCtrl = nil
	facade.GormDb = nil
	facade.SqlxDb = nil
	facade.SqlxDbCtrl = nil

	facade.dbConfig = commonDb.DbConfig{}

	facade.initialized = false
}

// Generates a new controller of sql.DB
func (facade *DbFacade) NewDbCtrl() *commonDb.DbController {
	return commonDb.NewDbController(facade.SqlDb)
}

// Generates a new controller of sqlx.DB
func (facade *DbFacade) NewSqlxDbCtrl() *commonSqlx.DbController {
	return commonSqlx.NewDbController(facade.SqlxDb)
}
