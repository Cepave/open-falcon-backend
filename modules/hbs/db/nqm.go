package db

import (
	"database/sql"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/hbs/model"
	log "github.com/Sirupsen/logrus"
	"time"
)

// Inserts a new agent or updates existing one
func RefreshAgentInfo(agent *model.NqmAgent) (err error) {
	/**
	 * Update the data
	 */
	if _, err = DB.Exec(
		`
		INSERT INTO nqm_agent(ag_connection_id, ag_hostname, ag_ip_address, ag_last_heartbeat)
		VALUES(?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			ag_hostname = VALUES(ag_hostname),
			ag_ip_address = VALUES(ag_ip_address),
			ag_last_heartbeat = VALUES(ag_last_heartbeat)
		`,
		agent.ConnectionId(),
		agent.Hostname(),
		[]byte(agent.IpAddress),
		time.Now(),
	); err != nil {

		log.Printf("Cannot refresh agent: [%v]. Error: %v", agent.ConnectionId(), err)
		return
	}
	// :~)

	/**
	 * Loads id from auto-generated PK
	 */
	if err = DB.QueryRow(
		`
		SELECT ag_id
		FROM nqm_agent
		WHERE ag_connection_id = ?
		`,
		agent.ConnectionId(),
	).Scan(&agent.Id); err != nil {
		return
	}
	// :~)

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
	result = &commonModel.NqmAgent{
		Id: agentId,
	}

	err = inTx(func(tx *sql.Tx) (err error) {
		var dbAgentName sql.NullString

		if err = tx.QueryRow(
			/**
			 * Gets one row if the executing of ping task is needed
			 *
			 * Gets no row if:
			 * 1) No ping task configuration
			 * 2) The period is overed yet
			 * 3) The agent is disabled
			 */
			`
			SELECT ag_name,
				isp_id, isp_name,
				pv_id, pv_name,
				ct_id, ct_name,
				ag_nt_id
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
				INNER JOIN
				owl_isp AS isp
				ON ag.ag_isp_id = isp.isp_id
				INNER JOIN
				owl_province AS pv
				ON ag.ag_pv_id = pv.pv_id
				INNER JOIN
				owl_city AS ct
				ON ag.ag_ct_id = ct.ct_id
			`,
			// :~)
			agentId, checkedTime.Unix(),
		).Scan(
			&dbAgentName,
			&result.IspId, &result.IspName,
			&result.ProvinceId, &result.ProvinceName,
			&result.CityId, &result.CityName,
			&result.NameTagId,
		); err != nil {

			/**
			 * There is no need to perform ping task
			 */
			result = nil
			if err == sql.ErrNoRows {
				err = nil
			}
			// :~)

			return
		}

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
		 * Updates the last execute
		 */
		if _, err = DB.Exec(
			`
			UPDATE nqm_agent_ping_task
			SET apt_time_last_execute = FROM_UNIXTIME(?)
			WHERE apt_ag_id = ?
			`,
			checkedTime.Unix(), agentId,
		); err != nil {

			result = nil
		}
		// :~)

		return
	})

	return
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
		)

		if currentTarget.NameTagId == commonModel.UNDEFINED_NAME_TAG_ID {
			currentTarget.NameTag = commonModel.UNDEFINED_STRING
		}
	}
	// :~)

	return
}

func loadAllEnabledTargets() (*sql.Rows, error) {
	return DB.Query(
		`
		SELECT
			tg_id, tg_host,
			isp_id, isp_name,
			pv_id, pv_name,
			ct_id, ct_name,
			nt_id, nt_value
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
		WHERE tg.tg_status = true
			AND tg.tg_available = true
		`,
	);
}

func loadTargetsByFilter(agentId int) (*sql.Rows, error) {
	return DB.Query(
		`
		SELECT
			tg_id, tg_host,
			isp_id, isp_name,
			pv_id, pv_name,
			ct_id, ct_name,
			nt_id, nt_value
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
	if err = DB.QueryRow(
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
					pt_number_of_city_filters
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
