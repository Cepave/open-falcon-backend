package db

import (
	"flag"
	"fmt"
	check "gopkg.in/check.v1"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
)

// The base environment for RDB testing

var dsnMysql = flag.String("dsn_mysql", "", "DSN of MySql")

func GetDbConfig(c *check.C) *commonDb.DbConfig {
	if *dsnMysql == "" {
		c.Skip(fmt.Sprintf("Skip database testing. Needs \"-dsn_mysql=<MySQL DSN>\""))
		return nil
	}

	return &commonDb.DbConfig {
		Dsn: *dsnMysql,
		MaxIdle: 2,
	}
}
func InitDbFacade(c *check.C) *commonDb.DbFacade {
	var dbFacade = &commonDb.DbFacade{}
	dbConfig := GetDbConfig(c)

	if dbConfig == nil {
		return nil
	}

	err := dbFacade.Open(dbConfig)
	c.Assert(err, check.IsNil)

	return dbFacade
}
func ReleaseDbFacade(c *check.C, dbFacade *commonDb.DbFacade) {
	if dbFacade != nil {
		dbFacade.Release()
	}
}

// Checks whether or not skipping testing by viable arguments
func HasDbEnvForMysqlOrSkip(c *check.C) bool {
	var hasMySqlDsn = *dsnMysql != ""

	if !hasMySqlDsn {
		c.Skip("Skip Mysql Test: -dsn_mysql=<dsn>")
	}

	return hasMySqlDsn
}
