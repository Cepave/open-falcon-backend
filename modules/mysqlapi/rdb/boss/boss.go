package boss

import (
	bossModel "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
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
