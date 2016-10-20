package basis

import (
	"database/sql"
	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	log "github.com/Sirupsen/logrus"
)

func UpdateAgent(agentInfo *model.AgentUpdateInfo) error {
	var err error = nil

	/**
	 * 如果 config 的 Hosts 有值，
	 * 只更新 host，不新增 host
	 */
	if g.Config().Hosts == "" {
		err = refreshAgent(agentInfo)
	} else {
		err = updateAgent(agentInfo)
	}
	// :~)

	if err != nil {
		log.Errorf("Refresh host info(by agent) has error: %v.", err)
	}

	return err
}

/**
 * Refresh agent or add a new one
 */
type refreshHostProcessor model.AgentReportRequest
func (self *refreshHostProcessor) BootCallback(tx *sql.Tx) bool {
	result := commonDb.ToTxExt(tx).Exec(
		`
		UPDATE host
		SET ip = ?,
			agent_version = ?,
			plugin_version = ?,
			update_at = NOW()
		WHERE hostname = ?
		`,
		self.IP,
		self.AgentVersion,
		self.PluginVersion,
		self.Hostname,
	)

	return commonDb.ToResultExt(result).RowsAffected() == 0
}
func (self *refreshHostProcessor) IfTrue(tx *sql.Tx) {
	commonDb.ToTxExt(tx).Exec(
		`
		INSERT INTO host(
			hostname,
			ip, agent_version, plugin_version
		)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			ip = ?,
			agent_version = ?,
			plugin_version = ?
		`,
		self.Hostname,
		self.IP,
		self.AgentVersion,
		self.PluginVersion,
		self.IP,
		self.AgentVersion,
		self.PluginVersion,
	)
}

func refreshAgent(agentInfo *model.AgentUpdateInfo) (dbError error) {
	processor := refreshHostProcessor(*agentInfo.ReportRequest)

	DbFacade.SqlDbCtrl.InTxForIf(&processor)

	return
}
func updateAgent(agentInfo *model.AgentUpdateInfo) (dbError error) {
	DbFacade.SqlDbCtrl.Exec(
		`
		UPDATE host
		SET ip = ?,
			agent_version = ?,
			plugin_version = ?
		WHERE hostname = ?
		`,
		agentInfo.ReportRequest.IP,
		agentInfo.ReportRequest.AgentVersion,
		agentInfo.ReportRequest.PluginVersion,
		agentInfo.ReportRequest.Hostname,
	)

	return
}
