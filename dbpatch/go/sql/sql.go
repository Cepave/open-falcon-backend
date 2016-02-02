// Connect and execute patch
package sql

import (
	dbsql "database/sql"
)

// Represents the connection information to database
type DatabaseConfig struct {
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

	dbConfig = &DatabaseConfig{
		driverName: driverName,
		dsn: dsn,
		db: openedDb,
	}
	return
}

// Close the db resource hold by the DatabaseConfig
func (databaseConfig *DatabaseConfig) Close() error {
	var oldDb = databaseConfig.db
	databaseConfig.db = nil

	return oldDb.Close()
}

// Execute a callback which accepts "database/sql.DB" object
func (databaseConfig *DatabaseConfig) Execute(
	dbCallback func(db *dbsql.DB) error,
) (err error) {

	var db *dbsql.DB = nil

	/**
	 * Opens the connection to database
	 */
	if db, err = dbsql.Open(
		databaseConfig.driverName,
		databaseConfig.dsn,
	)
		err != nil {
		return
	}
	// :~)

	return dbCallback(db)
}
