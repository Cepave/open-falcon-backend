package nqm

import (
	"testing"
	qtest "github.com/Cepave/open-falcon-backend/modules/query/test"
	. "gopkg.in/check.v1"
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
