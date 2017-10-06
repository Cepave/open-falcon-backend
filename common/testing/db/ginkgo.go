package db

import (
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	tflag "github.com/Cepave/open-falcon-backend/common/testing/flag"

	. "github.com/onsi/gomega"
)

// Initializes a configuration of database by DSN.
//
// If DSN is empty string, this function returns nil pointer.
//
// See "NewDbConfigByFlag(int)"
func NewDbConfigByDsn(dsn string) *commonDb.DbConfig {
	if dsn == "" {
		return nil
	}

	return &commonDb.DbConfig{
		Dsn:     dsn,
		MaxIdle: 2,
	}
}

// Initializes a configuration of database by test flags(OWL database)
//
// For example: flag.OWL_DB_PORTAL
//
// See godoc of "common/testing/flag"
func NewDbConfigByFlag(dbFlag int) *commonDb.DbConfig {
	return NewDbConfigByDsn(getTestFlags().GetMysqlOfOwlDb(dbFlag))
}

// Provides some utility functions to ease the construction of database connection for testing.
type GinkgoDb struct{}

// Constructs "*db/facade/DbFacade" object by configuration of flag.
//
// The checker object is used to trigger panic if the database cannot be opened.
//
// If the environment is not ready(flag is empty), this function returns "nil"
func (g *GinkgoDb) InitDbFacade() *f.DbFacade {
	return g.InitDbFacadeByDbConfig(g.GetDbConfig())
}

// Constructs "*db/facade/DbFacade" object by test flags(OWL database)
//
// For example: flag.OWL_DB_PORTAL
//
// See godoc of "common/testing/flag"
func (g *GinkgoDb) InitDbFacadeByFlag(flag int) *f.DbFacade {
	return g.InitDbFacadeByDbConfig(NewDbConfigByFlag(flag))
}

// Constructs "*db/facade/DbFacade" object by configuration object
//
// This function uses "Ginkgo" assertion to ensure that the database connection gets succeeded.
func (g *GinkgoDb) InitDbFacadeByDbConfig(dbConfig *commonDb.DbConfig) *f.DbFacade {
	if dbConfig == nil {
		return nil
	}

	var dbFacade = &f.DbFacade{}
	err := dbFacade.Open(dbConfig)
	Expect(err).To(BeNil())

	return dbFacade
}

// Constructs "*db/facade/DbFacade" object by default configuration object
//
// To support new property of "flag.OWL_DB_PORTAL", this function checks the flag first or
// use "flag.F_MySql" for compatibility of old flags.
//
// Deprecated: Since it is possible that a test environment needs multiple connections to different databases,
//	you should use "*GinkgoDb.InitDbFacadeByFlag" to initialize corresponding databases instead.
func (g *GinkgoDb) GetDbConfig() *commonDb.DbConfig {
	testFlags := getTestFlags()

	/**
	 * Supports new flag of OWL database
	 */
	if testFlags.HasMySqlOfOwlDb(tflag.OWL_DB_PORTAL) {
		return NewDbConfigByDsn(testFlags.GetMysqlOfOwlDb(tflag.OWL_DB_PORTAL))
	}
	// :~)

	return NewDbConfigByDsn(testFlags.GetMySql())
}

// Releases "*db/facade/DbFacade" object by configuration of flag.
//
// The checker object is used to trigger panic if the database cannot be released.
func (g *GinkgoDb) ReleaseDbFacade(dbFacade *f.DbFacade) {
	if dbFacade != nil {
		dbFacade.Release()
	}
}
