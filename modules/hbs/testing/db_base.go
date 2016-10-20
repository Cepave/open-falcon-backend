package testing

import (
	"flag"
	"fmt"
	"github.com/Cepave/open-falcon-backend/modules/hbs/db"
	check "gopkg.in/check.v1"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
)

// The base environment for RDB testing

var dsnMysql = flag.String("dsn_mysql", "", "DSN of MySql")

func init() {
	flag.Parse()
}

func InitDb(c *check.C) {
	if *dsnMysql == "" {
		c.Skip(fmt.Sprintf("Skip database testing. Needs \"-dsn_mysql=<MySQL DSN>\""))
		return
	}

	err := db.DbInit(
		&commonDb.DbConfig {
			Dsn: *dsnMysql,
			MaxIdle: 2,
		},
	)

	c.Assert(err, check.IsNil)
}
func ReleaseDb(c *check.C) {
	db.Release()
}

// Checks whether or not skipping testing by viable arguments
func HasDbEnvForMysqlOrSkip(c *check.C) bool {
	var hasMySqlDsn = *dsnMysql != ""

	if !hasMySqlDsn {
		c.Skip("Skip Mysql Test: -dsn_mysql=<dsn>")
	}

	return hasMySqlDsn
}
