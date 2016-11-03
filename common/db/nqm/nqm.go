package nqm

import (
	"database/sql"
	"time"
	"fmt"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	utils "github.com/Cepave/open-falcon-backend/common/utils"
)

/**
 * Refresh agent or add a new one
 */
type refreshAgentProcessor struct {
	agent *nqmModel.NqmAgent
}
func (self *refreshAgentProcessor) BootCallback(tx *sql.Tx) bool {
	result := commonDb.ToTxExt(tx).Exec(
		`
		UPDATE nqm_agent
		SET ag_hostname = ?,
			ag_ip_address = ?,
			ag_last_heartbeat = ?
		WHERE ag_connection_id = ?
		`,
		self.agent.Hostname(),
		[]byte(self.agent.IpAddress),
		time.Now(),
		self.agent.ConnectionId(),
	)

	return commonDb.ToResultExt(result).RowsAffected() == 0
}
func (self *refreshAgentProcessor) IfTrue(tx *sql.Tx) {
	now := time.Now()
	ipAddress := []byte(self.agent.IpAddress)

	commonDb.ToTxExt(tx).Exec(
		`
		INSERT INTO host(hostname, ip, agent_version, plugin_version)
		VALUES(?, ?, '', '')
		ON DUPLICATE KEY UPDATE
			ip = ?
		`,
		self.agent.Hostname(),
		self.agent.IpAddress.String(),
		self.agent.IpAddress.String(),
	)
	commonDb.ToTxExt(tx).Exec(
		`
		INSERT INTO nqm_agent(ag_connection_id, ag_hostname, ag_ip_address, ag_last_heartbeat, ag_hs_id)
		VALUES(
			?, ?, ?, ?,
			(
				SELECT id
				FROM host
				WHERE hostname = ?
			)
		)
		ON DUPLICATE KEY UPDATE
			ag_hostname = ?,
			ag_ip_address = ?,
			ag_last_heartbeat = ?
		`,
		self.agent.ConnectionId(),
		self.agent.Hostname(),
		ipAddress,
		now,
		self.agent.Hostname(),
		self.agent.Hostname(),
		ipAddress,
		now,
	)
}
func (self *refreshAgentProcessor) ResultRow(row *sql.Row) {
	commonDb.ToRowExt(row).Scan(&self.agent.Id)
}

// Inserts a new agent or updates existing one
func RefreshAgentInfo(agent *nqmModel.NqmAgent) (dbError error) {
	dbCtrl := commonDb.NewDbController(DbFacade.SqlDb)
	dbCtrl.RegisterPanicHandler(commonDb.NewDbErrorCapture(&dbError))

	agentProcessor := &refreshAgentProcessor{ agent }

	dbCtrl.InTxForIf(agentProcessor)
	if dbError != nil {
		return
	}

	dbCtrl.QueryForRow(
		agentProcessor,
		`
		SELECT ag_id
		FROM nqm_agent
		WHERE ag_connection_id = ?
		`,
		agent.ConnectionId(),
	)

	return
}

// Gets the data of agent for RPC
//
// If there is no need to perform ping task, this method would return nil as result.
//
// Reasons for not doing ping task:
// 1) No ping task configuration
// 2) The period is overed yet
func GetAndRefreshNeedPingAgentForRpc(agentId int, checkedTime time.Time) (result *commonModel.NqmAgent, err error) {
	dbCtrl := commonDb.NewDbController(DbFacade.SqlDb)
	dbCtrl.RegisterPanicHandler(commonDb.NewDbErrorCapture(&err))

	result = &commonModel.NqmAgent{}

	dbCtrl.InTx(commonDb.TxCallbackFunc(
		func (tx *sql.Tx) {
			txRefreshAgent(tx, agentId, checkedTime, result)
		},
	))

	if err != nil {
		result = nil
	}

	return
}

func txRefreshAgent(tx *sql.Tx, agentId int, checkedTime time.Time, result *commonModel.NqmAgent) {
	var temporaryTableName = fmt.Sprintf("tmp_effective_ping_task_%d", agentId)

	tx.Exec(fmt.Sprintf(
		"DROP TEMPORARY TABLE %s IF EXISTS", temporaryTableName,
	));

	tx.Exec(
		/**
		 * Creates temporary table for matched ping task
		 *
		 * Must be matched with all of the following conditions
		 * 1) The agent is enabled
		 * 2) The ping task is enabled
		 * 3) The period of time is elapsed or the ping task is never executed
		 */
		fmt.Sprintf(`
			CREATE TEMPORARY TABLE IF NOT EXISTS %s
			AS
			SELECT apt.apt_pt_id AS pt_id
			FROM nqm_agent AS ag
				INNER JOIN
				nqm_agent_ping_task AS apt
				ON apt.apt_ag_id = ?
					AND ag.ag_id = apt.apt_ag_id
					AND ag.ag_status = TRUE # Agent must be enabled
				INNER JOIN
				nqm_ping_task AS pt
				ON apt.apt_pt_id = pt.pt_id
					AND pt.pt_enable = TRUE # Task must be enabled
					AND TIMESTAMPDIFF(
						MINUTE,
						IFNULL(apt.apt_time_last_execute, FROM_UNIXTIME(0)), /* Use the very first time */
						FROM_UNIXTIME(?)
					) >= pt.pt_period
			`,
			temporaryTableName,
		),
		// :~)
		agentId, checkedTime.Unix(),
	);

	var numberOfPingTasks int

	/**
	 * If There is no ping task, this function gives no data
	 */
	rowForNumberOfPingTask := tx.QueryRow(fmt.Sprintf(
		`
		SELECT COUNT(pt_id) AS number_of_ping_task
		FROM %s
		`,
		temporaryTableName,
	))

	commonDb.ToRowExt(rowForNumberOfPingTask).Scan(&numberOfPingTasks);

	if numberOfPingTasks == 0 {
		return
	}
	// :~)

	var dbAgentName sql.NullString
	var concatedIdsOfGroupTag sql.NullString

	result.Id = agentId

	rowOfAgent := tx.QueryRow(
		/**
		 * Gets data of agent if any of the ping task need to be executed
		 */
		`
		SELECT ag_name,
			isp_id, isp_name,
			pv_id, pv_name,
			ct_id, ct_name,
			ag_nt_id,
			GROUP_CONCAT(agt.agt_gt_id ORDER BY agt.agt_gt_id ASC SEPARATOR ',') AS gts
		FROM nqm_agent AS ag
			INNER JOIN
			owl_isp AS isp
			ON ag.ag_isp_id = isp.isp_id
			INNER JOIN
			owl_province AS pv
			ON ag.ag_pv_id = pv.pv_id
			INNER JOIN
			owl_city AS ct
			ON ag.ag_ct_id = ct.ct_id
			LEFT OUTER JOIN
			nqm_agent_group_tag AS agt
			ON ag.ag_id = agt.agt_ag_id
		WHERE ag.ag_id = ?
		GROUP BY ag_name, isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, ag_nt_id
		`,
		// :~)
		agentId,
	)

	commonDb.ToRowExt(rowOfAgent).Scan(
		&dbAgentName,
		&result.IspId, &result.IspName,
		&result.ProvinceId, &result.ProvinceName,
		&result.CityId, &result.CityName,
		&result.NameTagId, &concatedIdsOfGroupTag,
	)

	result.GroupTagIds = utils.IntTo32(
		commonDb.GroupedStringToIntArray(concatedIdsOfGroupTag, ","),
	)

	/**
	 * Loads name of agent
	 */
	if dbAgentName.Valid {
		result.Name = dbAgentName.String
	} else {
		result.Name = commonModel.UNDEFINED_STRING
	}
	// :~)

	/**
	 * Updates the time of last executed
	 */
	commonDb.ToTxExt(tx).Exec(
		fmt.Sprintf(
			`
			UPDATE nqm_agent_ping_task AS apt,
				%s AS effective_pt
			SET apt_time_last_execute = FROM_UNIXTIME(?)
			WHERE apt.apt_pt_id = effective_pt.pt_id
				AND apt_ag_id = ?
			`,
			temporaryTableName,
		),
		checkedTime.Unix(), agentId,
	)
	// :~)
}

const (
	NO_PING_TASK = 0
	HAS_PING_TASK = 1
	HAS_PING_TASK_MATCH_ANY_TARGET = 2
)

// Gets the targets(to be probed) for RPC
func GetTargetsByAgentForRpc(agentId int) (targets []commonModel.NqmTarget, err error) {
	var taskState int

	if taskState, err = getPingTaskState(agentId); err != nil {
		return
	}

	if taskState == NO_PING_TASK {
		targets = make([]commonModel.NqmTarget, 0)
		return
	}

	var rows *sql.Rows

	switch taskState {
	case HAS_PING_TASK_MATCH_ANY_TARGET:
		rows ,err = loadAllEnabledTargets()
	case HAS_PING_TASK:
		rows, err = loadTargetsByFilter(agentId)
	}

	if err != nil {
		return
	}
	defer rows.Close()

	/**
	 * Converts the data to NQM targets for RPC
	 */
	targets = make([]commonModel.NqmTarget, 0, 32)
	for rows.Next() {
		targets = append(targets, commonModel.NqmTarget{})

		var currentTarget *commonModel.NqmTarget = &targets[len(targets)-1]
		var concatedIdsOfGroupTags sql.NullString

		rows.Scan(
			&currentTarget.Id,
			&currentTarget.Host,
			&currentTarget.IspId,
			&currentTarget.IspName,
			&currentTarget.ProvinceId,
			&currentTarget.ProvinceName,
			&currentTarget.CityId,
			&currentTarget.CityName,
			&currentTarget.NameTagId,
			&currentTarget.NameTag,
			&concatedIdsOfGroupTags,
		)

		currentTarget.GroupTagIds = utils.IntTo32(
			commonDb.GroupedStringToIntArray(concatedIdsOfGroupTags, ","),
		)

		if currentTarget.NameTagId == commonModel.UNDEFINED_NAME_TAG_ID {
			currentTarget.NameTag = commonModel.UNDEFINED_STRING
		}
	}
	// :~)

	return
}

func loadAllEnabledTargets() (*sql.Rows, error) {
	return DbFacade.SqlDb.Query(
		`
		SELECT
			tg_id, tg_host,
			isp_id, isp_name,
			pv_id, pv_name,
			ct_id, ct_name,
			nt_id, nt_value,
			GROUP_CONCAT(tgt.tgt_gt_id ORDER BY tgt.tgt_gt_id ASC SEPARATOR ',') AS gts
		FROM nqm_target AS tg
			INNER JOIN
			owl_isp AS isp
			ON tg.tg_isp_id = isp.isp_id
			INNER JOIN
			owl_province AS pv
			ON tg.tg_pv_id = pv.pv_id
			INNER JOIN
			owl_city AS ct
			ON tg.tg_ct_id = ct.ct_id
			INNER JOIN
			owl_name_tag AS nt
			ON tg.tg_nt_id = nt.nt_id
			LEFT OUTER JOIN
			nqm_target_group_tag AS tgt
			ON tg.tg_id = tgt.tgt_tg_id
		WHERE tg.tg_status = TRUE
			AND tg.tg_available = TRUE
		GROUP BY tg_id, tg_host, isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value
		`,
	);
}

func loadTargetsByFilter(agentId int) (*sql.Rows, error) {
	return DbFacade.SqlDb.Query(
		`
		SELECT
			tg_id, tg_host,
			isp_id, isp_name,
			pv_id, pv_name,
			ct_id, ct_name,
			nt_id, nt_value,
			GROUP_CONCAT(tgt.tgt_gt_id ORDER BY tgt.tgt_gt_id ASC SEPARATOR ',') AS gts
		FROM (
				/* Matched target by ISP */
				SELECT tg_id, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
				FROM
					nqm_agent_ping_task AS apt
					INNER JOIN
					nqm_ping_task AS pt
					ON pt.pt_id = apt.apt_pt_id
						AND apt.apt_ag_id = ?
					INNER JOIN
					vw_enabled_targets_by_ping_task AS vw_tg
					ON pt.pt_id = vw_tg.tg_pt_id
				/* :~) */
				UNION
				/* Matched target which to be probed by all */
				SELECT tg_id, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
				FROM nqm_target tg
				WHERE tg_probed_by_all = true
					AND tg.tg_status = true
					AND tg.tg_available = true
				/* :~) */
			) AS tg
			INNER JOIN
			owl_isp AS isp
			ON tg.tg_isp_id = isp.isp_id
			INNER JOIN
			owl_province AS pv
			ON tg.tg_pv_id = pv.pv_id
			INNER JOIN
			owl_city AS ct
			ON tg.tg_ct_id = ct.ct_id
			INNER JOIN
			owl_name_tag AS nt
			ON tg.tg_nt_id = nt.nt_id
			LEFT OUTER JOIN
			nqm_target_group_tag AS tgt
			ON tg.tg_id = tgt.tgt_tg_id
		GROUP BY tg_id, tg_host, isp_id, isp_name, pv_id, pv_name, ct_id, ct_name, nt_id, nt_value
		`,
		agentId,
	)
}

func getPingTaskState(agentId int) (result int, err error) {
	result = NO_PING_TASK

	var numberOfViablePingTasks int
	var numberOfEmptyPingTasks int

	/**
	 * Checks if there is any PING TASK(enabled)
	 */
	if err = DbFacade.SqlDb.QueryRow(
		`
		SELECT
			COUNT(IF(pt_has_filter = 1, 1, NULL)) AS num_of_viable_ping_task,
			COUNT(IF(pt_has_filter = 0, 1, NULL)) AS num_of_empty_ping_task
		FROM (
			SELECT
				/**
				 * 1 - The ping task has at least one filter
				 * 0 - The ping task has no filter(matching all targets)
				 */
				IF(
					pt_number_of_name_tag_filters +
					pt_number_of_isp_filters +
					pt_number_of_province_filters +
					pt_number_of_city_filters +
					pt_number_of_group_tag_filters
					> 0,
					1, 0
				) AS pt_has_filter
				# //:~)
			FROM nqm_agent_ping_task AS apt
				INNER JOIN
				nqm_ping_task AS pt
				ON apt.apt_pt_id = pt.pt_id
					AND pt.pt_enable = TRUE
					AND apt_ag_id = ?
		) AS pt_filter_counter
		`,
		agentId,
	).Scan(&numberOfViablePingTasks, &numberOfEmptyPingTasks); err != nil {
		return
	}

	if numberOfEmptyPingTasks > 0 {
		result = HAS_PING_TASK_MATCH_ANY_TARGET
	} else if numberOfViablePingTasks > 0 {
		result = HAS_PING_TASK
	}
	// :~)

	return
}
