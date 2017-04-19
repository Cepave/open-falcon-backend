package rdb

import (
	"github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
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

func AgentHeartbeat(agents []*model.AgentHeartbeat, updateOnly bool) *model.AgentHeartbeatResult {
	if updateOnly {
		return updateHost(agents)
	}
	return updateOrInsertHost(agents)
}

func updateOrInsertHost(agents []*model.AgentHeartbeat) *model.AgentHeartbeatResult {
	updateOrInsertHosts := &updateOrInsertHostsInTx{
		hosts: agents,
	}

	DbFacade.SqlxDbCtrl.InTx(updateOrInsertHosts)

	return &updateOrInsertHosts.result
}

func updateHost(agents []*model.AgentHeartbeat) *model.AgentHeartbeatResult {
	updateHosts := &updateHostsInTx{
		hosts: agents,
	}

	DbFacade.SqlxDbCtrl.InTx(updateHosts)

	return &updateHosts.result
}

type updateHostsInTx struct {
	hosts  []*model.AgentHeartbeat
	result model.AgentHeartbeatResult
}

func (uHost *updateHostsInTx) InTx(tx *sqlx.Tx) db.TxFinale {
	updateStmt := sqlxExt.ToTxExt(tx).Preparex(updateSql)

	for _, host := range uHost.hosts {
		uHost.result.RowsAffected += updateAndGetRowsAffected(updateStmt, host)
	}

	return db.TxCommit
}

func updateAndGetRowsAffected(updateStmt *sqlx.Stmt, agent *model.AgentHeartbeat) int64 {
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
	hosts  []*model.AgentHeartbeat
	result model.AgentHeartbeatResult
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

func (uoiHost *updateOrInsertHostsInTx) isHostExisting(selectExt *sqlxExt.StmtExt, agent *model.AgentHeartbeat) bool {
	var count int
	selectExt.Get(&count, agent.Hostname)
	return count >= 1
}

func (uoiHost *updateOrInsertHostsInTx) addHost(insertStmt *sqlx.Stmt, agent *model.AgentHeartbeat) {
	r := insertStmt.MustExec(
		agent.Hostname,
		agent.IP,
		agent.AgentVersion,
		agent.PluginVersion,
		agent.UpdateTime,
	)
	uoiHost.result.RowsAffected += db.ToResultExt(r).RowsAffected()
}

func (uoiHost *updateOrInsertHostsInTx) updateHost(updateStmt *sqlx.Stmt, agent *model.AgentHeartbeat) {
	uoiHost.result.RowsAffected += updateAndGetRowsAffected(updateStmt, agent)
}
