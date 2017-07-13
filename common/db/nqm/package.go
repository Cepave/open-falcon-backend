package nqm

import (
	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	tb "github.com/Cepave/open-falcon-backend/common/textbuilder"
)

var DbFacade *f.DbFacade

var t = tb.Dsl

var logger = log.NewDefaultLogger("warn")
