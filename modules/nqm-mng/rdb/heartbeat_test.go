package rdb

import (
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
)

type TestUpdateOrInsertSuite struct{}

var _ = Suite(&TestUpdateOrInsertSuite{})

func (suite *TestUpdateOrInsertSuite) TestAddHost(c *C) {
	testCases := []*struct{}
}

func (suite *TestUpdateOrInsertSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
	owlDb.DbFacade = DbFacade
}

func (suite *TestUpdateOrInsertSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
	owlDb.DbFacade = nil
}
