package facade

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	commonSqlx "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

// Gives facade interface supporting multiple object of db
//
// This facade supports:
// 	gorm - github.com/jinzhu/gorm
// 	sqlx - github.com/Cepave/open-falcon-backend/common/db/sqlx
// 	database/sql.DB
// 	dbCtrl
type DbFacade struct {
	SqlDb *sql.DB
	SqlDbCtrl *commonDb.DbController
	GormDb *gorm.DB
	SqlxDb *sqlx.DB
	SqlxDbCtrl *commonSqlx.DbController

	initialized bool
}

// Open this facade with ping()
func (facade *DbFacade) Open(dbConfig *commonDb.DbConfig) (err error) {
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

	facade.SqlxDb = sqlx.NewDb(facade.SqlDb, "mysql")
	facade.SqlxDbCtrl = commonSqlx.NewDbController(facade.SqlxDb)

	facade.SqlDbCtrl = commonDb.NewDbController(facade.SqlDb)
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
	facade.SqlxDb = nil
	facade.SqlxDbCtrl = nil

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
