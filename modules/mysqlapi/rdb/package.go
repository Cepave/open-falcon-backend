package rdb

import (
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	nqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"

	graphdb "github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/graph"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/hbsdb"
	owldb "github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/owl"
)

var logger = log.NewDefaultLogger("INFO")

var DbFacade = &f.DbFacade{}
var DbConfig *commonDb.DbConfig = &commonDb.DbConfig{}

func InitPortalRdb(dbConfig *commonDb.DbConfig) {
	logger.Infof("Open RDB: %s ...", dbConfig)

	err := DbFacade.Open(dbConfig)
	if err != nil {
		logger.Warnf("Open database error: %v", err)
	}

	/**
	 * Protal database
	 */
	nqmDb.DbFacade = DbFacade
	owldb.DbFacade = DbFacade

	hbsdb.DbFacade = DbFacade
	hbsdb.DB = DbFacade.SqlDb
	// :~)

	*DbConfig = *dbConfig

	logger.Info("[FINISH] Open RDB.")
}
func InitGraphRdb(dbConfig *commonDb.DbConfig) {
	graphDbFacade := &f.DbFacade{}

	logger.Infof("Open RDB: %s ...", dbConfig)

	err := graphDbFacade.Open(dbConfig)
	if err != nil {
		logger.Warnf("Open database error: %v", err)
	}

	graphdb.DbFacade = graphDbFacade
}

func ReleaseAllRdb() {
	logger.Info("Release RDB resources...")

	/**
	 * Protal database
	 */
	DbFacade.Release()

	nqmDb.DbFacade = nil
	owlDb.DbFacade = nil

	hbsdb.DbFacade = nil
	hbsdb.DB = nil
	// :~)

	/**
	 * Graph database
	 */
	graphdb.DbFacade.Release()
	graphdb.DbFacade = nil
	// :~)

	logger.Info("[FINISH] Release RDB resources.")
}
