package hbsdb

import (
	"database/sql"

	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
var DbFacade = &f.DbFacade{}

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
