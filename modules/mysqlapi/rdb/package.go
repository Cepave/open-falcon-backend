package rdb

import (
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	nqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
)

var logger = log.NewDefaultLogger("INFO")

var DbFacade = &f.DbFacade{}
var DbConfig *commonDb.DbConfig = &commonDb.DbConfig{}

func InitRdb(dbConfig *commonDb.DbConfig) {
	logger.Infof("Open RDB: %s ...", dbConfig)

	err := DbFacade.Open(dbConfig)
	if err != nil {
		logger.Warnf("Open database error: %v", err)
	}

	nqmDb.DbFacade = DbFacade
	owlDb.DbFacade = DbFacade

	*DbConfig = *dbConfig

	logger.Info("[FINISH] Open RDB.")
}
func ReleaseRdb() {
	logger.Info("Release RDB resources...")

	DbFacade.Release()

	nqmDb.DbFacade = nil
	owlDb.DbFacade = nil

	logger.Info("[FINISH] Release RDB resources.")
}
