package nqm

import (
	tb "github.com/Cepave/open-falcon-backend/common/textbuilder"
	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
)

var DbFacade *f.DbFacade

var t = tb.Dsl

var logger = log.NewDefaultLogger("warn")
