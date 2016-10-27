package testing

import (
	"github.com/Cepave/open-falcon-backend/modules/hbs/db"
	check "gopkg.in/check.v1"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
)

// The base environment for RDB testing
func InitDb(c *check.C) {
	db.DbInit(dbTest.GetDbConfig(c))
}
func ReleaseDb(c *check.C) {
	db.Release()
}
