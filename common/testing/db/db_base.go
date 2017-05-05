//
// The base environment for RDB testing
//
// Flags
//
// This package has pre-defined flags of command:
//
// 	-dsn_mysql - MySQL DSN used to intialize configuration of mysql connection
package db

import (
	"flag"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
)

var dsnMysql = flag.String("dsn_mysql", "", "DSN of MySql")

// This callback is used to setup a viable database configuration for testing.
type ViableDbConfigFunc func(config *commonDb.DbConfig)
