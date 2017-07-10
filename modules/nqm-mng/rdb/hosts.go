package rdb

import (
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/jinzhu/gorm"
)

// ListHosts returns the info of hosts by host ID
func ListHosts(paging commonModel.Paging) ([]*model.HostsResult, *commonModel.Paging) {
	var result []*model.HostsResult

	var funcTxLoader gormExt.TxCallbackFunc = func(txGormDb *gorm.DB) commonDb.TxFinale {
		var dbListHosts = txGormDb.Model(&model.HostsResult{}).
			Select(`SQL_CALC_FOUND_ROWS
				host.hostname,
				host.id,
				GROUP_CONCAT(g.id ORDER BY g.id ASC SEPARATOR ',') AS gid,
				GROUP_CONCAT(g.grp_name ORDER BY g.id ASC SEPARATOR '\0') AS gname
			`).
			Joins(`
				LEFT JOIN grp_host gh
				ON host.id = gh.host_id
				LEFT JOIN grp g
				ON gh.grp_id = g.id`).
			Group(`host.id, host.hostname`).
			Limit(paging.Size).
			Order(`host.id ASC`).
			Offset(paging.GetOffset())

		selectHost := dbListHosts.Find(&result)
		gormExt.ToDefaultGormDbExt(selectHost).PanicIfError()

		return commonDb.TxCommit
	}

	gormExt.ToDefaultGormDbExt(DbFacade.GormDb).SelectWithFoundRows(
		funcTxLoader, &paging,
	)

	return result, &paging
}
