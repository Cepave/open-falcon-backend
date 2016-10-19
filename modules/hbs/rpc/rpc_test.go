package rpc

import (
	"testing"

	"github.com/Cepave/open-falcon-backend/modules/hbs/db"
	hbstesting "github.com/Cepave/open-falcon-backend/modules/hbs/testing"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestRpcSuite struct{}

var _ = Suite(&TestRpcSuite{})

func (s *TestRpcSuite) SetUpSuite(c *C) {
	if !hbstesting.HasDbEnvForMysqlOrSkip(c) {
		return
	}

	hbstesting.DoInitDb(db.DbInit)
}

func (s *TestRpcSuite) TearDownSuite(c *C) {
	hbstesting.DoReleaseDb(db.Release)
}
