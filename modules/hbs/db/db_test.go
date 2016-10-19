package db

import (
	hbstesting "github.com/Cepave/open-falcon-backend/modules/hbs/testing"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type TestDbSuite struct{}

var _ = Suite(&TestDbSuite{})

func (s *TestDbSuite) SetUpSuite(c *C) {
	if !hbstesting.HasDbEnvForMysqlOrSkip(c) {
		return
	}

	hbstesting.DoInitDb(DbInit)
}

func (s *TestDbSuite) TearDownSuite(c *C) {
	hbstesting.DoReleaseDb(Release)
}
