package nqm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
)

// Lists the agents by query condition
func ListAgents(query *nqmModel.AgentQuery, paging commonModel.Paging) ([]*nqmModel.Agent, *commonModel.Paging) {
	var result []*nqmModel.Agent

	var funcTxLoader gormExt.TxCallbackFunc = func(txGormDb *gorm.DB) {
		/**
		 * Retrieves the page of data
		 */
		var selectAgent = txGormDb.Model(&nqmModel.Agent{}).
			Select(`SQL_CALC_FOUND_ROWS
				ag_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status, ag_comment, ag_last_heartbeat,
				isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value,
				COUNT(gt.gt_id) AS gt_number,
				GROUP_CONCAT(gt.gt_id ORDER BY gt_name ASC SEPARATOR ',') AS gt_ids,
				GROUP_CONCAT(gt.gt_name ORDER BY gt_name ASC SEPARATOR '\0') AS gt_names
			`).
			Joins(`
				INNER JOIN
				owl_isp AS isp
				ON ag_isp_id = isp.isp_id
				INNER JOIN
				owl_province AS pv
				ON ag_pv_id = pv.pv_id
				INNER JOIN
				owl_city AS ct
				ON ag_ct_id = ct.ct_id
				INNER JOIN
				owl_name_tag AS nt
				ON ag_nt_id = nt.nt_id
				LEFT OUTER JOIN
				nqm_agent_group_tag AS agt
				ON ag_id = agt.agt_ag_id
				LEFT OUTER JOIN
				owl_group_tag AS gt
				ON agt.agt_gt_id = gt.gt_id
			`).
			Limit(paging.Size).
			Group(`
				ag_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status, ag_comment, ag_last_heartbeat,
				isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value
			`).
			Order(buildSortingClauseOfAgents(&paging)).
			Offset(paging.GetOffset())

		if query.Name != "" {
			selectAgent = selectAgent.Where("ag_name LIKE ?", query.Name + "%")
		}
		if query.ConnectionId != "" {
			selectAgent = selectAgent.Where("ag_connection_id LIKE ?", query.ConnectionId + "%")
		}
		if query.Hostname != "" {
			selectAgent = selectAgent.Where("ag_hostname LIKE ?", query.Hostname + "%")
		}
		if query.HasIspId {
			selectAgent = selectAgent.Where("ag_isp_id = ?", query.IspId)
		}
		if query.HasStatusCondition {
			selectAgent = selectAgent.Where("ag_status = ?", query.Status)
		}
		if query.IpAddress != "" {
			selectAgent = selectAgent.Where("ag_ip_address LIKE ?", query.GetIpForLikeCondition())
		}
		// :~)

		gormExt.ToDefaultGormDbExt(selectAgent.Find(&result)).PanicIfError()
	}

	gormExt.ToDefaultGormDbExt(DbFacade.GormDb).SelectWithFoundRows(
		funcTxLoader, &paging,
	)

	/**
	 * Loads group tags
	 */
	for _, agent := range result {
		agent.AfterLoad()
	}
	// :~)

	return result, &paging
}

var orderByDialect = commonModel.NewSqlOrderByDialect(
	map[string]string {
		"status": "ag_status",
		"name": "ag_name",
		"connection_id": "ag_connection_id",
		"comment": "ag_comment",
		"province": "pv_name",
		"city": "ct_name",
		"last_heartbeat_time": "ag_last_heartbeat",
		"name_tag": "nt_value",
	},
)
func init() {
	originFunc := orderByDialect.FuncEntityToSyntax
	orderByDialect.FuncEntityToSyntax = func(entity *commonModel.OrderByEntity) (string, error) {
		switch entity.Expr {
		case "group_tag":
			var dirOfGroupTags = "DESC"
			if entity.Direction == commonModel.Ascending {
				dirOfGroupTags = "ASC"
			}
			return fmt.Sprintf("gt_number %s, gt_names %s", dirOfGroupTags, dirOfGroupTags), nil
		}

		return originFunc(entity)
	}
}

func buildSortingClauseOfAgents(paging *commonModel.Paging) string {
	if len(paging.OrderBy) == 0 {
		paging.OrderBy = append(paging.OrderBy, &commonModel.OrderByEntity{ "status", commonModel.Descending })
	}

	if len(paging.OrderBy) == 1 {
		switch paging.OrderBy[0].Expr {
		case "province":
			paging.OrderBy = append(paging.OrderBy, &commonModel.OrderByEntity{ "city", commonModel.Ascending })
		}
	}

	if paging.OrderBy[len(paging.OrderBy) - 1].Expr != "last_heartbeat_time" {
		paging.OrderBy = append(paging.OrderBy, &commonModel.OrderByEntity{ "last_heartbeat_time", commonModel.Descending })
	}

	querySyntax, err := orderByDialect.ToQuerySyntax(paging.OrderBy)
	gormExt.DefaultGormErrorConverter.PanicIfError(err)

	return querySyntax
}
