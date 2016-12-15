// Connect and execute patch
package sql

import (
	dbsql "database/sql"
	"github.com/jmoiron/sqlx"
	osqlx "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	"fmt"
)

// Represents the connection information to database
type DatabaseConfig struct {
	SqlxDb *sqlx.DB
	SqlxDbCtrl *osqlx.DbController

	driverName string
	dsn string
	db *dbsql.DB
}

// Initialize a new configuration to database
// This function also tries to ping the database
func NewDatabaseConfig(driverName string, dsn string) (dbConfig *DatabaseConfig, err error) {
	/**
	 * Opens the connection to database
	 */
	var openedDb *dbsql.DB
	if openedDb, err = dbsql.Open(driverName, dsn)
		err != nil {
		return
	}

	if err = openedDb.Ping()
		err != nil {
		return
	}
	// :~)

	sqlxDb := sqlx.NewDb(openedDb, driverName)
	dbConfig = &DatabaseConfig{
		SqlxDb: sqlxDb,
		SqlxDbCtrl: osqlx.NewDbController(sqlxDb),
		driverName: driverName,
		dsn: dsn,
		db: openedDb,
	}
	return
}

// Close the db resource hold by the DatabaseConfig
func (databaseConfig *DatabaseConfig) Close() error {
	var dbObject = databaseConfig.db
	databaseConfig.db = nil
	databaseConfig.SqlxDbCtrl = nil
	databaseConfig.SqlxDb = nil

	return dbObject.Close()
}

// Execute a callback which accepts "database/sql.DB" object
func (databaseConfig *DatabaseConfig) Execute(
	dbCallback func(db *dbsql.DB) error,
) (err error) {

	if databaseConfig.db == nil {
		return fmt.Errorf("Need open connection of database")
	}

	return dbCallback(databaseConfig.db)
}

func (c *DatabaseConfig) GetDatabaseName() string {
	var dbName string

	c.SqlxDbCtrl.QueryRowxAndScan(
		`SELECT DATABASE()`, nil, &dbName,
	)

	return dbName
}
