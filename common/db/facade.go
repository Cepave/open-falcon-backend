package db

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

// Gives facade interface supporting multiple object of db
//
// This facade supports:
// 	gorm - github.com/jinzhu/gorm
// 	database/sql.DB
// 	dbCtrl
type DbFacade struct {
	SqlDb *sql.DB
	SqlDbCtrl *DbController
	GormDb *gorm.DB

	initialized bool
}

// Configuration of database
type DbConfig struct {
	Dsn string
	MaxIdle int
}

func (config *DbConfig) String() string {
	return fmt.Sprintf("DSN: [%s]. Max Idle: [%d]", config.Dsn, config.MaxIdle)
}

// Open this facade with ping()
func (facade *DbFacade) Open(dbConfig *DbConfig) (err error) {
	if facade.initialized {
		return
	}

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

	facade.SqlDbCtrl = NewDbController(facade.SqlDb)
	facade.initialized = true

	return
}
// Close the database, release the resources
func (facade *DbFacade) Release() {
	if !facade.initialized {
		return
	}

	facade.GormDb.Close()

	facade.SqlDb = nil
	facade.SqlDbCtrl = nil
	facade.GormDb = nil

	facade.initialized = false
}
// Generates a new controller of sql.DB
func (facade *DbFacade) NewDbCtrl() *DbController {
	return NewDbController(facade.SqlDb)
}
