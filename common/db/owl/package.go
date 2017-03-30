package owl

import (
	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
)

var DbFacade *f.DbFacade
var logger = log.NewDefaultLogger("warn")
