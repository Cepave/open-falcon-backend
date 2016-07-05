package db

import (
	"testing"
	hbstesting "github.com/Cepave/open-falcon-backend/modules/hbs/testing"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestDbSuite struct{}

var _ = Suite(&TestDbSuite{})

func (s *TestDbSuite) SetUpSuite(c *C) {
	hbstesting.InitDb()
	DB = hbstesting.DbForTest
}

func (s *TestDbSuite) TearDownSuite(c *C) {
	hbstesting.ReleaseDb()
	DB = nil
}
