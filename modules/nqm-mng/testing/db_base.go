package testing

import (
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/rdb"
	check "gopkg.in/check.v1"
)

// The base environment for RDB testing
func InitRdb(c *check.C) {
	dbTest.SetupByViableDbConfig(c, rdb.InitRdb)
}
func ReleaseRdb(c *check.C) {
	rdb.ReleaseRdb()
}
