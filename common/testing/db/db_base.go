package db

import (
	"flag"
	"fmt"
	check "gopkg.in/check.v1"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	f "github.com/Cepave/open-falcon-backend/common/db/facade"
)

// The base environment for RDB testing

var dsnMysql = flag.String("dsn_mysql", "", "DSN of MySql")

// This callback is used to setup a viable database configuration while testing
type ViableDbConfigFunc func(config *commonDb.DbConfig)

// This function is used to:
//
// 1) Check whether or not the configuration o "dsn_mysql" has been supplied
// 2) If it does, supply the data of configuration to callback function
func SetupByViableDbConfig(c *check.C, configFunc ViableDbConfigFunc) bool {
	config := GetDbConfig(c)

	if config != nil {
		configFunc(config)
	}

	return config != nil
}

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
func InitDbFacade(c *check.C) *f.DbFacade {
	var dbFacade = &f.DbFacade{}
	dbConfig := GetDbConfig(c)

	if dbConfig == nil {
		return nil
	}

	err := dbFacade.Open(dbConfig)
	c.Assert(err, check.IsNil)

	return dbFacade
}
func ReleaseDbFacade(c *check.C, dbFacade *f.DbFacade) {
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
