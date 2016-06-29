package rpc

import (
	"testing"

	"github.com/Cepave/hbs/db"
	hbstesting "github.com/Cepave/hbs/testing"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestRpcSuite struct{}

var _ = Suite(&TestRpcSuite{})

func (s *TestRpcSuite) SetUpSuite(c *C) {
	hbstesting.InitDb()
	db.DB = hbstesting.DbForTest
}

func (s *TestRpcSuite) TearDownSuite(c *C) {
	hbstesting.ReleaseDb()
	db.DB = nil
}
