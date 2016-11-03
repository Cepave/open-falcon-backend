package testing

import (
	"github.com/Cepave/open-falcon-backend/modules/nqm-conf-mng/rdb"
	check "gopkg.in/check.v1"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
)

// The base environment for RDB testing
func InitRdb(c *check.C) {
	dbTest.SetupByViableDbConfig(c, rdb.InitRdb)
}
func ReleaseRdb(c *check.C) {
	rdb.ReleaseRdb()
}

