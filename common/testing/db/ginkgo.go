package db

import (
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	f "github.com/Cepave/open-falcon-backend/common/db/facade"

	. "github.com/onsi/gomega"
)

type GinkgoDb struct{}

// Constructs "*db/facade/DbFacade" object by configuration of flag.
//
// The checker object is used to trigger panic if the database cannot be opened.
//
// If the environment is not ready(flag is empty), this function returns "nil"
func (g *GinkgoDb) InitDbFacade() *f.DbFacade {
	var dbFacade = &f.DbFacade{}
	dbConfig := g.GetDbConfig()

	if dbConfig == nil {
		return nil
	}

	err := dbFacade.Open(dbConfig)
	Expect(err).To(BeNil())

	return dbFacade
}

func (g *GinkgoDb) GetDbConfig() *commonDb.DbConfig {
	if !getTestFlags().HasMySql() {
		return nil
	}

	return &commonDb.DbConfig{
		Dsn:     getTestFlags().GetMySql(),
		MaxIdle: 2,
	}
}

// Releases "*db/facade/DbFacade" object by configuration of flag.
//
// The checker object is used to trigger panic if the database cannot be released.
func (g *GinkgoDb) ReleaseDbFacade(dbFacade *f.DbFacade) {
	if dbFacade != nil {
		dbFacade.Release()
	}
}
