package nqm

import (
	qtest "github.com/Cepave/open-falcon-backend/modules/query/test"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type dbTestSuite struct{}

func (s *dbTestSuite) SetUpSuite(c *C) {
	qtest.InitDb(c)
	initServices()
}
func (s *dbTestSuite) TearDownSuite(c *C) {
	qtest.ReleaseDb(c)
}
