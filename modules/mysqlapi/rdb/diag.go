package rdb

import (
	"database/sql"
	"fmt"

	apiModel "github.com/Cepave/open-falcon-backend/common/model/mysqlapi"
	"github.com/go-sql-driver/mysql"
)

// Performs the diagnosis to RDB
func DiagnoseRdb(dsn string, db *sql.DB) *apiModel.Rdb {
	var v int
	err := db.QueryRow("SELECT 0 FROM DUAL").Scan(&v)

	pingResult := 0
	pingMessage := ""
	if err != nil {
		pingResult = 1
		pingMessage = err.Error()
	}

	return &apiModel.Rdb{
		Dsn:             hidePasswordOfDsn(dsn),
		OpenConnections: db.Stats().OpenConnections,
		PingResult:      pingResult,
		PingMessage:     pingMessage,
	}
}

func hidePasswordOfDsn(dsn string) string {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return fmt.Sprintf("Cannot parse DSN:[%s]. %v", dsn, err)
	}

	return fmt.Sprintf(
		"%s:!hidden password!@%s(%s)/%s",
		config.User, config.Net, config.Addr,
		config.DBName,
	)
}
