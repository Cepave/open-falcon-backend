package boss

import (
	"github.com/jmoiron/sqlx"

	bossModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
)

// Get the Data from boss.hosts
func GetSyncData() []*bossModel.BossHost {
	var hosts []*bossModel.BossHost
	DbFacade.NewSqlxDbCtrl().Select(
		&hosts,
		`
		SELECT hostname, ip, activate, platform FROM hosts
			WHERE exist = 1
		`,
	)

	return hosts
}
