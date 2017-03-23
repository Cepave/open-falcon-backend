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
	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	check "gopkg.in/check.v1"
)

var dsnMysql = flag.String("dsn_mysql", "", "DSN of MySql")

// This callback is used to setup a viable database configuration for testing.
type ViableDbConfigFunc func(config *commonDb.DbConfig)

// This function is used to:
//
// 	1. Check whether or not the configuration o "dsn_mysql" has been supplied
// 	2. If it does, supply the data of configuration to callback function
func SetupByViableDbConfig(c *check.C, configFunc ViableDbConfigFunc) bool {
	config := GetDbConfig(c)

	if config != nil {
		configFunc(config)
	}

	return config != nil
}

// Gets the database configuration or skip the testing(depends on "gopkg.in/check.v1").
//
// If the environment is not ready(flag is empty), this function returns "nil"
func GetDbConfig(c *check.C) *commonDb.DbConfig {
	if *dsnMysql == "" {
		c.Skip("Skip database testing. Needs \"-dsn_mysql=<MySQL DSN>\"")
		return nil
	}

	return &commonDb.DbConfig{
		Dsn:     *dsnMysql,
		MaxIdle: 2,
	}
}

// Constructs "*db/facade/DbFacade" object by configuration of flag.
//
// The checker object is used to trigger panic if the database cannot be opened.
//
// If the environment is not ready(flag is empty), this function returns "nil"
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

// Releases "*db/facade/DbFacade" object by configuration of flag.
//
// The checker object is used to trigger panic if the database cannot be released.
func ReleaseDbFacade(c *check.C, dbFacade *f.DbFacade) {
	if dbFacade != nil {
		dbFacade.Release()
	}
}

// Checks whether or not skipping testing by viable arguments.
//
// If the environment is not ready(flag is empty), this function returns false value.
func HasDbEnvForMysqlOrSkip(c *check.C) bool {
	var hasMySqlDsn = *dsnMysql != ""

	if !hasMySqlDsn {
		c.Skip("Skip Mysql Test: -dsn_mysql=<dsn>")
	}

	return hasMySqlDsn
}
