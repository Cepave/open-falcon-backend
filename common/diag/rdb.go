package diag

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

type RdbDiagnosis struct {
	Dsn             string `json:"dsn"`
	OpenConnections int    `json:"open_connections"`
	PingResult      int    `json:"ping_result"`
	PingMessage     string `json:"ping_message"`
}

// Performs the diagnosis to RDB
func DiagnoseRdb(dsn string, db *sql.DB) *RdbDiagnosis {
	var v int
	err := db.QueryRow("SELECT 0 FROM DUAL").Scan(&v)

	pingResult := 0
	pingMessage := ""
	if err != nil {
		pingResult = 1
		pingMessage = err.Error()
	}

	return &RdbDiagnosis{
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
		"%s:!hide password!@%s(%s)/%s",
		config.User, config.Net, config.Addr,
		config.DBName,
	)
}
