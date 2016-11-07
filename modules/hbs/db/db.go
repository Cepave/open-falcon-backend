package db

import (
	"database/sql"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	dbNqm "github.com/Cepave/open-falcon-backend/common/db/nqm"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
var DbFacade = &f.DbFacade{}

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

func DbInit(dbConfig *commonDb.DbConfig) (err error) {
	err = DbFacade.Open(dbConfig)
	if err != nil {
		return
	}

	DB = DbFacade.SqlDb
	dbNqm.DbFacade = DbFacade

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
