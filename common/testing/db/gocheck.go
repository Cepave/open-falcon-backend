package db

import (
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	check "gopkg.in/check.v1"
)

// This function is used to:
//
// 	1. Checks whether or not the property of "mysql" has been supplied
// 	2. If it does, supply the data of configuration to callback function
//
// Deprecated: Try to use "flag.SkipFactory"(common/testing/flag) instead
func SetupByViableDbConfig(c *check.C, configFunc ViableDbConfigFunc) bool {
	config := getDbConfig(c)

	if config != nil {
		configFunc(config)
	}

	return config != nil
}

// Constructs "*db/facade/DbFacade" object by configuration of flag.
//
// The checker object is used to trigger panic if the database cannot be opened.
//
// If the environment is not ready(flag is empty), this function returns "nil"
//
// Deprecated: Try to use "flag.SkipFactory"(common/testing/flag) instead
func InitDbFacade(c *check.C) *f.DbFacade {
	var dbFacade = &f.DbFacade{}
	dbConfig := getDbConfig(c)

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

// Gets the database configuration or skip the testing(depends on "gopkg.in/check.v1").
//
// If the environment is not ready(flag is empty), this function returns "nil"
//
// Deprecated: Try to use "flag.SkipFactory"
//
// See "common/testing/flag"
func getDbConfig(c *check.C) *commonDb.DbConfig {
	if !getTestFlags().HasMySql() {
		c.Skip(flagMessage)
		return nil
	}

	return &commonDb.DbConfig{
		Dsn:     getTestFlags().GetMySql(),
		MaxIdle: 2,
	}
}
