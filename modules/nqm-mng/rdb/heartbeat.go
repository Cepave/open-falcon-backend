package rdb

import (
	"github.com/Cepave/open-falcon-backend/common/db"
	sqlxExt "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/jmoiron/sqlx"
)

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
	for _, host := range uHost.hosts {
		uHost.result.RowsAffected += updateAndGetRowsAffected(tx, host)
	}

	return db.TxCommit
}

func updateAndGetRowsAffected(tx *sqlx.Tx, agent *model.AgentHeartbeat) int64 {
	r := tx.MustExec(`
		UPDATE host
		SET ip = ?,
			agent_version = ?,
			plugin_version = ?,
			update_at = FROM_UNIXTIME(?)
		WHERE hostname = ?
			AND update_at < FROM_UNIXTIME(?)
		`,
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
	for _, host := range uoiHost.hosts {
		if uoiHost.isHostExisting(tx, host) {
			uoiHost.updateHost(tx, host)
		} else {
			uoiHost.addHost(tx, host)
		}
	}

	return db.TxCommit
}

func (uoiHost *updateOrInsertHostsInTx) isHostExisting(tx *sqlx.Tx, agent *model.AgentHeartbeat) bool {
	var count int
	sqlxExt.ToTxExt(tx).Get(
		&count,
		`
		SELECT COUNT("hostname")
		FROM host
		WHERE hostname = ?
		`,
		agent.Hostname,
	)

	return count >= 1
}

func (uoiHost *updateOrInsertHostsInTx) addHost(tx *sqlx.Tx, agent *model.AgentHeartbeat) {
	r := tx.MustExec(`
		INSERT INTO host(
			hostname,
			ip, agent_version, plugin_version, update_at
		)
		VALUES (?, ?, ?, ?, FROM_UNIXTIME(?))
		`,
		agent.Hostname,
		agent.IP,
		agent.AgentVersion,
		agent.PluginVersion,
		agent.UpdateTime,
	)
	uoiHost.result.RowsAffected += db.ToResultExt(r).RowsAffected()
}

func (uoiHost *updateOrInsertHostsInTx) updateHost(tx *sqlx.Tx, agent *model.AgentHeartbeat) {
	uoiHost.result.RowsAffected += updateAndGetRowsAffected(tx, agent)
}
