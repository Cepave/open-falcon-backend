package boss

import (
	model "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
)

// Get the Data from boss.hosts
func GetSyncData() []*model.BossHost {
	var hosts []*model.BossHost
	DbFacade.SqlxDbCtrl.Select(
		&hosts,
		`
		SELECT hostname, ip, activate, platform
		FROM hosts
		WHERE exist = 1
		`,
	)

	return hosts
}
