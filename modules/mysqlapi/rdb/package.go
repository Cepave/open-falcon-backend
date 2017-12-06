package rdb

import (
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	f "github.com/Cepave/open-falcon-backend/common/db/facade"
	commonNqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	commonOwlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	apiModel "github.com/Cepave/open-falcon-backend/common/model/mysqlapi"

	graphdb "github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/graph"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/hbsdb"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/cmdb"
	apiOwlDb "github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/owl"
)

const (
	DB_PORTAL = "portal"
	DB_GRAPH  = "graph"
)

type DbHolder struct {
	facades map[string]*f.DbFacade
}

func (self *DbHolder) setDb(dbname string, facade *f.DbFacade) {
	self.facades[dbname] = facade
}
func (self *DbHolder) releaseDb(dbname string) {
	if facade, ok := self.facades[dbname]; ok {
		facade.Release()
		delete(self.facades, dbname)
	}
}
func (self *DbHolder) Diagnose(dbname string) *apiModel.Rdb {
	facade, ok := self.facades[dbname]

	if !ok {
		return nil
	}

	return DiagnoseRdb(facade.GetDbConfig().Dsn, facade.SqlDb)
}

var GlobalDbHolder *DbHolder = &DbHolder{
	facades: make(map[string]*f.DbFacade),
}

var logger = log.NewDefaultLogger("INFO")

var DbFacade = &f.DbFacade{}

func InitPortalRdb(dbConfig *commonDb.DbConfig) {
	GlobalDbHolder.setDb(DB_PORTAL, DbFacade)

	logger.Infof("Open RDB: %s ...", dbConfig)

	err := DbFacade.Open(dbConfig)
	if err != nil {
		logger.Warnf("Open database error: %v", err)
	}

	DbFacade.SetReleaseCallback(func() {
		commonNqmDb.DbFacade = nil
		commonOwlDb.DbFacade = nil
		apiOwlDb.DbFacade = nil

		hbsdb.DbFacade = nil
		hbsdb.DB = nil
	})

	/**
	 * Protal database
	 */
	commonNqmDb.DbFacade = DbFacade
	commonOwlDb.DbFacade = DbFacade
	apiOwlDb.DbFacade = DbFacade
	cmdb.DbFacade = DbFacade

	hbsdb.DbFacade = DbFacade
	hbsdb.DB = DbFacade.SqlDb
	// :~)

	logger.Info("[FINISH] Open RDB.")
}
func InitGraphRdb(dbConfig *commonDb.DbConfig) {
	graphDbFacade := &f.DbFacade{}
	GlobalDbHolder.setDb(DB_GRAPH, graphDbFacade)

	logger.Infof("Open RDB: %s ...", dbConfig)

	err := graphDbFacade.Open(dbConfig)
	if err != nil {
		logger.Warnf("Open database error: %v", err)
	}

	DbFacade.SetReleaseCallback(func() {
		graphdb.DbFacade = nil
	})

	graphdb.DbFacade = graphDbFacade
}

func ReleaseAllRdb() {
	logger.Info("Release RDB resources...")

	GlobalDbHolder.releaseDb(DB_PORTAL)
	GlobalDbHolder.releaseDb(DB_GRAPH)

	logger.Info("[FINISH] Release RDB resources.")
}
