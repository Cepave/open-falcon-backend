package rpc

import (
	"flag"
	"testing"
	"github.com/Cepave/open-falcon-backend/modules/hbs/db"
	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var DbFacade *f.DbFacade = db.DbFacade

func init() {
	flag.Parse()
}
