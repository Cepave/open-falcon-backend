package db

import (
	"database/sql"
	"fmt"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/jinzhu/gorm"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
var GormDb *gorm.DB

// Initialize the resource for RDB
func Init() {
	err := dbInit(g.Config().Database)

	if err != nil {
		log.Fatalln(err)
	}
}
// Initialize the resource for RDB
func Release() {
	if GormDb != nil {
		GormDb.Close()
	}
}

func dbInit(dsn string) (err error) {
	/**
	 * Initialize Gorm(It would call ping())
	 */
	GormDb, err = gorm.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("Open Gorm error: %v", err)
	}
	// :~)

	/**
	 * Use the sql.DB object from Gorm and ping
	 */
	DB = GormDb.DB()
	DB.SetMaxIdleConns(g.Config().MaxIdle)
	// :~)

	return
}

// Convenient IoC for transaction processing
func inTx(txCallback func(tx *sql.Tx) error) (err error) {
	var tx *sql.Tx

	if tx, err = DB.Begin(); err != nil {
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
