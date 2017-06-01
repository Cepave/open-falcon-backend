package rdb

import (
	"github.com/Cepave/open-falcon-backend/common/db"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
)

var updateSql = `
	UPDATE host
	SET ip = ?,
		agent_version = ?,
		plugin_version = ?,
		update_at = FROM_UNIXTIME(?)
	WHERE hostname = ?
		AND update_at < FROM_UNIXTIME(?)
`

func FalconAgentHeartbeat(agents []*model.FalconAgentHeartbeat, updateOnly bool) *model.FalconAgentHeartbeatResult {
	if updateOnly {
		return updateHost(agents)
	}
	return updateOrInsertHost(agents)
}

func updateOrInsertHost(agents []*model.FalconAgentHeartbeat) *model.FalconAgentHeartbeatResult {
	updateOrInsertHosts := &updateOrInsertHostsInTx{
		hosts: agents,
	}

	DbFacade.SqlxDbCtrl.InTx(updateOrInsertHosts)

	return &updateOrInsertHosts.result
}

func updateHost(agents []*model.FalconAgentHeartbeat) *model.FalconAgentHeartbeatResult {
	updateHosts := &updateHostsInTx{
		hosts: agents,
	}

	DbFacade.SqlxDbCtrl.InTx(updateHosts)

	return &updateHosts.result
}

type updateHostsInTx struct {
	hosts  []*model.FalconAgentHeartbeat
	result model.FalconAgentHeartbeatResult
}

func (uHost *updateHostsInTx) InTx(tx *sqlx.Tx) db.TxFinale {
	updateStmt := sqlxExt.ToTxExt(tx).Preparex(updateSql)

	for _, host := range uHost.hosts {
		uHost.result.RowsAffected += updateAndGetRowsAffected(updateStmt, host)
	}

	return db.TxCommit
}

func updateAndGetRowsAffected(updateStmt *sqlx.Stmt, agent *model.FalconAgentHeartbeat) int64 {
	r := updateStmt.MustExec(
		agent.IP,
		agent.AgentVersion,
		agent.PluginVersion,
		agent.UpdateTime,
		agent.Hostname,
		agent.UpdateTime,
	)
	return db.ToResultExt(r).RowsAffected()
}

type updateOrInsertHostsInTx struct {
	hosts  []*model.FalconAgentHeartbeat
	result model.FalconAgentHeartbeatResult
}

func (uoiHost *updateOrInsertHostsInTx) InTx(tx *sqlx.Tx) db.TxFinale {
	selectStmt := sqlxExt.ToStmtExt(sqlxExt.ToTxExt(tx).Preparex(`
		SELECT COUNT(id)
		FROM host
		WHERE hostname = ?
		FOR UPDATE
		`,
	))
	updateStmt := sqlxExt.ToTxExt(tx).Preparex(updateSql)
	insertStmt := sqlxExt.ToTxExt(tx).Preparex(`
		INSERT INTO host(
			hostname,
			ip, agent_version, plugin_version, update_at
		)
		VALUES (?, ?, ?, ?, FROM_UNIXTIME(?))
		ON DUPLICATE KEY UPDATE
			ip = VALUES(ip),
			agent_version = VALUES(agent_version),
			plugin_version = VALUES(plugin_version),
			update_at = IF(TIMESTAMPDIFF(SECOND, update_at, VALUES(update_at)) > 0, VALUES(update_at), update_at)
	`)

	for _, host := range uoiHost.hosts {
		if uoiHost.isHostExisting(selectStmt, host) {
			uoiHost.updateHost(updateStmt, host)
		} else {
			uoiHost.addHost(insertStmt, host)
		}
	}

	return db.TxCommit
}

func (uoiHost *updateOrInsertHostsInTx) isHostExisting(selectExt *sqlxExt.StmtExt, agent *model.FalconAgentHeartbeat) bool {
	var count int
	selectExt.Get(&count, agent.Hostname)
	return count >= 1
}

func (uoiHost *updateOrInsertHostsInTx) addHost(insertStmt *sqlx.Stmt, agent *model.FalconAgentHeartbeat) {
	r := insertStmt.MustExec(
		agent.Hostname,
		agent.IP,
		agent.AgentVersion,
		agent.PluginVersion,
		agent.UpdateTime,
	)
	uoiHost.result.RowsAffected += db.ToResultExt(r).RowsAffected()
}

func (uoiHost *updateOrInsertHostsInTx) updateHost(updateStmt *sqlx.Stmt, agent *model.FalconAgentHeartbeat) {
	uoiHost.result.RowsAffected += updateAndGetRowsAffected(updateStmt, agent)
}

func UpdateNqmAgentHeartbeat(reqs []*model.NqmAgentHeartbeatRequest) {
	updateTx := &updateNqmAgentHeartbeatTx{
		Reqs: reqs,
	}
	DbFacade.SqlxDbCtrl.InTx(updateTx)
}

type updateNqmAgentHeartbeatTx struct {
	Reqs []*model.NqmAgentHeartbeatRequest
}

func (t *updateNqmAgentHeartbeatTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	updateStmt := sqlxExt.ToTxExt(tx).Preparex(`
	UPDATE nqm_agent
	SET ag_hostname = ?,
		ag_ip_address = ?,
		ag_last_heartbeat = FROM_UNIXTIME(?)
	WHERE ag_connection_id = ?
		AND (ag_last_heartbeat < FROM_UNIXTIME(?) OR
				 ag_last_heartbeat is NULL)
	`)

	for _, e := range t.Reqs {
		updateStmt.MustExec(
			e.Hostname,
			e.IpAddress,
			e.Timestamp,
			e.ConnectionId,
			e.Timestamp,
		)
	}
	return commonDb.TxCommit
}

func SelectNqmAgentByConnId(connId string) *nqmModel.Agent {
	var selectAgent = DbFacade.GormDb.Model(&nqmModel.Agent{}).
		Select(`
			ag_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status, ag_comment, ag_last_heartbeat,
			COUNT(DISTINCT pt.pt_id) AS ag_num_of_enabled_pingtasks,
			isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value,
			GROUP_CONCAT(gt.gt_id ORDER BY gt_name ASC SEPARATOR ',') AS gt_ids,
			GROUP_CONCAT(gt.gt_name ORDER BY gt_name ASC SEPARATOR '\0') AS gt_names
		`).
		Joins(`
			LEFT JOIN
			nqm_agent_ping_task AS apt
			ON ag_id = apt.apt_ag_id
			LEFT JOIN
			nqm_ping_task AS pt
			ON apt.apt_pt_id = pt.pt_id AND pt.pt_enable=true
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
		Where("ag_connection_id = ?", connId).
		Group(`
			ag_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address, ag_status, ag_comment, ag_last_heartbeat,
			isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value
		`)

	var loadedAgent = &nqmModel.Agent{}
	selectAgent = selectAgent.Find(loadedAgent)

	if selectAgent.Error == gorm.ErrRecordNotFound {
		return nil
	}
	gormExt.ToDefaultGormDbExt(selectAgent).PanicIfError()

	loadedAgent.AfterLoad()

	if !loadedAgent.Status {
		return nil
	}

	return loadedAgent
}

func NotNewNqmAgent(connId string) bool {
	return DbFacade.SqlxDbCtrl.GetOrNoRow(
		new(int),
		`
			SELECT ag_id
			FROM nqm_agent
			WHERE ag_connection_id = ?
		`,
		connId,
	)
}

type insertNqmAgentByHeartbeatTx struct {
	Req *model.NqmAgentHeartbeatRequest
}

func InsertNqmAgentByHeartbeat(r *model.NqmAgentHeartbeatRequest) {
	insertTx := &insertNqmAgentByHeartbeatTx{
		Req: r,
	}
	DbFacade.SqlxDbCtrl.InTx(insertTx)
}

func (t *insertNqmAgentByHeartbeatTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	DbFacade.SqlxDb.MustExec(`
		INSERT INTO host(hostname, ip, agent_version, plugin_version)
		VALUES(?, ?, '', '')
		ON DUPLICATE KEY UPDATE
			ip = VALUES(ip)
	`,
		t.Req.Hostname,
		string(t.Req.IpAddress),
	)
	DbFacade.SqlxDb.MustExec(`
		INSERT INTO nqm_agent(ag_connection_id, ag_hostname, ag_ip_address, ag_last_heartbeat, ag_hs_id)
		SELECT ?, ?, ?, FROM_UNIXTIME(?), id
		FROM host
		WHERE hostname = ?
		ON DUPLICATE KEY UPDATE
			ag_hostname = VALUES(ag_hostname),
			ag_ip_address = VALUES(ag_ip_address),
			ag_last_heartbeat = VALUES(ag_last_heartbeat)
	`,
		t.Req.ConnectionId,
		t.Req.Hostname,
		t.Req.IpAddress,
		t.Req.Timestamp,
		t.Req.Hostname,
	)
	return commonDb.TxCommit
}
