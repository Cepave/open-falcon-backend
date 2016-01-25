package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/hbs/g"
	"log"
)

var DB *sql.DB

func Init() {
	err := dbInit(g.Config().Database)

	if err != nil {
		log.Fatalf("open db fail: %v", err)
	}

	DB.SetMaxIdleConns(g.Config().MaxIdle)
}

func dbInit(dsn string) (err error) {
	if DB, err = sql.Open("mysql", dsn)
		err != nil {
		return
	}

	if err = DB.Ping()
		err != nil {
		return
	}

	return
}

// Convenient IoC for transaction processing
func inTx(txCallback func(tx *sql.Tx) error) (err error) {
	var tx *sql.Tx

	if tx, err = DB.Begin()
		err != nil {
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
