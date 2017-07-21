package database

import (
	log "github.com/sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/modules/query/g"
	"github.com/jinzhu/gorm"

	"github.com/Cepave/open-falcon-backend/common/db/facade"

	cdb "github.com/Cepave/open-falcon-backend/common/db"
	nqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var PortalDbFacade *facade.DbFacade

var (
	db  *gorm.DB
	err error
)

func DBConn() *gorm.DB {
	return db
}

func Init() {
	conf := g.Config()

	/**
	 * Use Db Facade to initialize related service
	 */
	PortalDbFacade = &facade.DbFacade{}
	err = PortalDbFacade.Open(
		&cdb.DbConfig{
			Dsn:     conf.Db.Addr,
			MaxIdle: conf.Db.Idle,
		},
	)

	if err != nil {
		log.Printf("%v\n", err)
	}

	owlDb.DbFacade = PortalDbFacade
	nqmDb.DbFacade = PortalDbFacade
	// :~)

	db = PortalDbFacade.GormDb
}
