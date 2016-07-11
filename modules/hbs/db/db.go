package db

import (
	"database/sql"
	"fmt"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB

func Init() {
	err := dbInit(g.Config().Database)

	if err != nil {
		log.Fatalln(err)
	}

	DB.SetMaxIdleConns(g.Config().MaxIdle)
}

func dbInit(dsn string) (err error) {
	if DB, err = sql.Open("mysql", dsn); err != nil {
		return fmt.Errorf("Open DB error: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("Ping DB error: %v", err)
	}

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
