package rpc

import (
	"testing"
	"github.com/Cepave/open-falcon-backend/modules/hbs/db"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var DbFacade *commonDb.DbFacade = db.DbFacade
