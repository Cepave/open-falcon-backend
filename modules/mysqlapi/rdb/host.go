package rdb

import (
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
)

// ListHosts returns the info of hosts by host ID
func ListHosts(paging commonModel.Paging) ([]*model.HostsResult, *commonModel.Paging) {
	var result []*model.HostsResult

	var funcTxLoader gormExt.TxCallbackFunc = func(txGormDb *gorm.DB) commonDb.TxFinale {
		var dbListHosts = txGormDb.Model(&model.HostsResult{}).
			Select(`SQL_CALC_FOUND_ROWS
				host.hostname,
				host.id,
				GROUP_CONCAT(g.id ORDER BY g.id ASC SEPARATOR ',') AS gt_ids,
				GROUP_CONCAT(g.grp_name ORDER BY g.id ASC SEPARATOR '\0') AS gt_names
			`).
			Joins(`
				LEFT JOIN grp_host gh
				ON host.id = gh.host_id
				LEFT JOIN grp g
				ON gh.grp_id = g.id`).
			Group(`host.id, host.hostname`).
			Limit(paging.Size).
			Order(buildSortingClauseOfHosts(&paging)).
			Offset(paging.GetOffset())

		selectHost := dbListHosts.Find(&result)
		gormExt.ToDefaultGormDbExt(selectHost).PanicIfError()

		return commonDb.TxCommit
	}

	gormExt.ToDefaultGormDbExt(DbFacade.GormDb).SelectWithFoundRows(
		funcTxLoader, &paging,
	)

	/**
	 * Loads group tags
	 */
	for _, host := range result {
		host.AfterLoad()
	}

	return result, &paging
}

var orderByDialectForHosts = commonModel.NewSqlOrderByDialect(
	map[string]string{
		"id":   "id",
		"name": "hostname",
	},
)

func buildSortingClauseOfHosts(paging *commonModel.Paging) string {
	querySyntax, err := orderByDialectForHosts.ToQuerySyntax(paging.OrderBy)
	gormExt.DefaultGormErrorConverter.PanicIfError(
		errors.Annotate(err, "Order by to query syntax has error"),
	)

	return querySyntax
}

// ListHostgroups returns the info of hostgroups by group ID
func ListHostgroups(paging commonModel.Paging) ([]*model.HostgroupsResult, *commonModel.Paging) {
	var result []*model.HostgroupsResult

	var funcTxLoader gormExt.TxCallbackFunc = func(txGormDb *gorm.DB) commonDb.TxFinale {
		var dbListHosts = txGormDb.Model(&model.HostgroupsResult{}).
			Select(`SQL_CALC_FOUND_ROWS
				grp.id,
				grp.grp_name,
				GROUP_CONCAT(pd.id ORDER BY pd.id ASC SEPARATOR ',') AS gt_ids,
				GROUP_CONCAT(pd.dir ORDER BY pd.id ASC SEPARATOR '\0') AS gt_names
			`).
			Joins(`
				LEFT JOIN plugin_dir pd
				ON grp.id = pd.grp_id
			`).
			Group(`grp.id, grp.grp_name`).
			Limit(paging.Size).
			Order(buildSortingClauseOfHostgroups(&paging)).
			Offset(paging.GetOffset())

		selectHost := dbListHosts.Find(&result)
		gormExt.ToDefaultGormDbExt(selectHost).PanicIfError()

		return commonDb.TxCommit
	}

	gormExt.ToDefaultGormDbExt(DbFacade.GormDb).SelectWithFoundRows(
		funcTxLoader, &paging,
	)

	/**
	 * Loads group tags
	 */
	for _, hostgroup := range result {
		hostgroup.AfterLoad()
	}

	return result, &paging
}

var orderByDialectForHostgroups = commonModel.NewSqlOrderByDialect(
	map[string]string{
		"id":   "id",
		"name": "grp_name",
	},
)

func buildSortingClauseOfHostgroups(paging *commonModel.Paging) string {
	querySyntax, err := orderByDialectForHostgroups.ToQuerySyntax(paging.OrderBy)
	gormExt.DefaultGormErrorConverter.PanicIfError(
		errors.Annotate(err, "Order by to query syntax has error"),
	)

	return querySyntax
}
