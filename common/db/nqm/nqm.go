package nqm

import (
	"database/sql"
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	osqlx "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	utils "github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/jmoiron/sqlx"
	"time"
)

// Inserts a new agent or updates existing one
//
// returns true if the agent is enabled
func RefreshAgentInfo(agent *nqmModel.NqmAgent, checkedTime time.Time) *commonModel.NqmAgent {
	refreshTx := &refreshAgentProcessor{
		agent:       agent,
		checkedTime: checkedTime,
	}

	// Builds or updates data of NQM agent(including host)
	DbFacade.SqlxDbCtrl.InTx(refreshTx)

	/**
	 * Loads id and status of NQM agent
	 */
	var isAgentEnable bool
	DbFacade.SqlxDbCtrl.QueryRowxAndScan(
		`
		SELECT ag_id, ag_status
		FROM nqm_agent
		WHERE ag_connection_id = ?
		`,
		[]interface{}{agent.ConnectionId()},
		&refreshTx.agent.Id,
		&isAgentEnable,
	)
	// :~)

	if !isAgentEnable {
		return nil
	}

	return loadAgentDetail(agent.Id)
}

/**
 * Refresh agent or add a new one
 */
type refreshAgentProcessor struct {
	agent       *nqmModel.NqmAgent
	checkedTime time.Time
}

func (self *refreshAgentProcessor) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	if self.BootCallback(tx) {
		self.IfTrue(tx)
	}

	return commonDb.TxCommit
}
func (self *refreshAgentProcessor) BootCallback(tx *sqlx.Tx) bool {
	agent := self.agent

	result := DbFacade.SqlxDb.MustExec(
		`
		UPDATE nqm_agent
		SET ag_hostname = ?,
			ag_ip_address = ?,
			ag_last_heartbeat = FROM_UNIXTIME(?)
		WHERE ag_connection_id = ?
		`,
		agent.Hostname(),
		[]byte(self.agent.IpAddress),
		self.checkedTime.Unix(),
		agent.ConnectionId(),
	)

	return commonDb.ToResultExt(result).RowsAffected() == 0
}
func (self *refreshAgentProcessor) IfTrue(tx *sqlx.Tx) {
	ipAddress := []byte(self.agent.IpAddress)

	tx.MustExec(
		`
		INSERT INTO host(hostname, ip, agent_version, plugin_version)
		VALUES(?, ?, '', '')
		ON DUPLICATE KEY UPDATE
			ip = VALUES(ip)
		`,
		self.agent.Hostname(),
		self.agent.IpAddress.String(),
	)

	tx.MustExec(
		`
		INSERT INTO nqm_agent(ag_connection_id, ag_hostname, ag_ip_address, ag_last_heartbeat, ag_hs_id)
		SELECT ?, ?, ?, FROM_UNIXTIME(?), id
		FROM host
		WHERE hostname = ?
		ON DUPLICATE KEY UPDATE
			ag_hostname = VALUES(ag_hostname),
			ag_ip_address = VALUES(ag_ip_address),
			ag_last_heartbeat = VALUES(ag_last_heartbeat)
		`,
		self.agent.ConnectionId(),
		self.agent.Hostname(),
		ipAddress,
		self.checkedTime.Unix(),
		self.agent.Hostname(),
	)
}

// Gets the ping list from cache
func GetPingListFromCache(agent *nqmModel.NqmAgent, checkedTime time.Time) ([]commonModel.NqmTarget, *nqmModel.PingListLog) {
	agentId := int32(agent.Id)

	pingListLog := getCacheLogOfPingList(agentId)

	/**
	 * Builds cache(synchronized) if there is no one
	 */
	hasCache := pingListLog != nil
	if !hasCache {
		pingListLog = BuildCacheOfPingList(agentId, checkedTime)
	}
	// :~)

	result := getPingList(agent, checkedTime)

	/**
	 * Refresh the cache for 1 hour live
	 */
	go utils.BuildPanicCapture(
		func() {
			logger.Debugf("Update Access Time. Agent Id: [%d]", agentId)
			updateAccessTime(agentId, checkedTime)
		},
		func(p interface{}) {
			logger.Errorf("Update access time of ping task has error. Agent Id: [%d]. %v", agent.Id, p)
		},
	)()
	// :~)

	return result, pingListLog
}

func updateAccessTime(agentId int32, accessTime time.Time) {
	DbFacade.SqlxDbCtrl.InTx(osqlx.TxCallbackFunc(
		func(tx *sqlx.Tx) commonDb.TxFinale {
			/**
			 * Updates the access time of log
			 */
			tx.MustExec(
				`
				UPDATE nqm_cache_agent_ping_list_log
				SET apll_time_access = FROM_UNIXTIME(?)
				WHERE apll_ag_id = ?
				`,
				accessTime.Unix(), agentId,
			)
			// :~)

			/**
			 * Updates the access time of target in cache
			 */
			tx.MustExec(
				`
				UPDATE nqm_cache_agent_ping_list
				SET apl_time_access = FROM_UNIXTIME(?)
				WHERE apl_apll_ag_id = ?
					AND TIMESTAMPDIFF(MINUTE, apl_time_access, FROM_UNIXTIME(?)) >= apl_min_period
				`,
				accessTime.Unix(), agentId, accessTime.Unix(),
			)
			// :~)

			return commonDb.TxCommit
		},
	))
}

type nqmTargetImpl struct {
	Id   int    `db:"apl_tg_id"`
	Host string `db:"tg_host"`

	IspId   int16  `db:"isp_id"`
	IspName string `db:"isp_name"`

	ProvinceId   int16  `db:"pv_id"`
	ProvinceName string `db:"pv_name"`

	CityId   int16  `db:"ct_id"`
	CityName string `db:"ct_name"`

	NameTagId    int16  `db:"nt_id"`
	NameTagValue string `db:"nt_value"`

	StreamOfGroupTagIds sql.NullString `db:"gts"`
}

func (t *nqmTargetImpl) toNqmTarget() commonModel.NqmTarget {
	return commonModel.NqmTarget{
		Id:           t.Id,
		Host:         t.Host,
		IspId:        t.IspId,
		IspName:      t.IspName,
		ProvinceId:   t.ProvinceId,
		ProvinceName: t.ProvinceName,
		CityId:       t.CityId,
		CityName:     t.CityName,
		NameTagId:    t.NameTagId,
		NameTag:      t.NameTagValue,
		GroupTagIds: utils.IntTo32(
			commonDb.GroupedStringToIntArray(t.StreamOfGroupTagIds, ","),
		),
	}
}

// 1. Update the access time
// 2. Retrieve the ping list
func getPingList(agent *nqmModel.NqmAgent, checkedTime time.Time) []commonModel.NqmTarget {
	implResult := make([]*nqmTargetImpl, 0)

	DbFacade.SqlxDbCtrl.Select(
		&implResult,
		`
		SELECT apl_tg_id, tg_host,
			isp_id, isp_name,
			pv_id, pv_name,
			ct_id, ct_name,
			nt_id, nt_value,
			GROUP_CONCAT(tgt.tgt_gt_id ORDER BY tgt.tgt_gt_id ASC SEPARATOR ',') AS gts
		FROM nqm_cache_agent_ping_list_log
			INNER JOIN
			nqm_cache_agent_ping_list
			ON apll_ag_id = apl_apll_ag_id
				AND apll_ag_id = ?
				AND TIMESTAMPDIFF(MINUTE, apl_time_access, FROM_UNIXTIME(?)) >= apl_min_period
			INNER JOIN
			nqm_target AS tg
			ON apl_tg_id = tg.tg_id
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
		GROUP BY apl_tg_id, tg_host,
			isp_id, isp_name, pv_id, pv_name, ct_id, ct_name,
			nt_id, nt_value
		ORDER BY apl_tg_id ASC
		`,
		agent.Id, checkedTime.Unix(),
	)

	/**
	 * Converts the data of table to target
	 */
	result := make([]commonModel.NqmTarget, 0)
	for _, targetImpl := range implResult {
		/**
		 * Skips the same ip address of agent for target
		 */
		if targetImpl.Host == agent.IpAddress.String() {
			continue
		}
		// :~)

		result = append(result, targetImpl.toNqmTarget())
	}
	// :~)

	return result
}
func getCacheLogOfPingList(agentId int32) *nqmModel.PingListLog {
	listLog := &nqmModel.PingListLog{}

	hasFound := DbFacade.SqlxDbCtrl.GetOrNoRow(
		listLog,
		`
		SELECT apll_number_of_targets,
			apll_time_access, apll_time_refresh
		FROM nqm_cache_agent_ping_list_log
		WHERE apll_ag_id = ?
		`,
		agentId,
	)

	if !hasFound {
		return nil
	}

	return listLog
}

// 1. Update(or INSERT) the refresh time
// 2. Re-build the list of targets
func BuildCacheOfPingList(agentId int32, checkedTime time.Time) *nqmModel.PingListLog {
	logger.Debugf("Agent Id[%d] -> Set to be delete flags...", agentId)
	DbFacade.SqlxDbCtrl.InTx(
		&toBeDeletedTargets{agentId: agentId},
	)
	logger.Debugf("Agent Id[%d] -> Finish", agentId)

	logger.Debugf("Agent Id[%d] -> Refresh targets...", agentId)
	DbFacade.SqlxDbCtrl.InTx(
		&refreshTargets{agentId: agentId},
	)
	logger.Debugf("Agent Id[%d] -> Finish", agentId)

	logger.Debugf("Agent Id[%d] -> Remove to-be-deleted targets...", agentId)
	DbFacade.SqlxDbCtrl.InTx(
		&removeToBeDeletedTargets{agentId: agentId},
	)
	logger.Debugf("Agent Id[%d] -> Finish", agentId)

	logger.Debugf("Agent Id[%d] -> Update final status of cache log...", agentId)
	resultTx := &updateRefreshTime{
		agentId:     agentId,
		checkedTime: checkedTime,
		logObject:   &nqmModel.PingListLog{},
	}
	DbFacade.SqlxDbCtrl.InTx(resultTx)
	logger.Debugf("Agent Id[%d] -> Finish", agentId)

	return resultTx.logObject
}

type toBeDeletedTargets struct {
	agentId int32
}

func (t *toBeDeletedTargets) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	/**
	 * Sets the status to `1`
	 */
	tx.MustExec(
		`
		INSERT INTO nqm_cache_agent_ping_list_log(
			apll_ag_id, apll_number_of_targets, apll_status,
			apll_time_access, apll_time_refresh
		)
		VALUES(?, -1, 1, FROM_UNIXTIME(0), NOW())
		ON DUPLICATE KEY UPDATE
			apll_number_of_targets = VALUES(apll_number_of_targets),
			apll_status = VALUES(apll_status),
			apll_time_refresh = VALUES(apll_time_refresh)
		`,
		t.agentId,
	)
	// :~)

	/**
	 * Marks the targets in cache as `2`(to be deleted)
	 */
	tx.MustExec(
		`
		UPDATE nqm_cache_agent_ping_list
		SET apl_build_flag = 2
		WHERE apl_apll_ag_id = ?
		`,
		t.agentId,
	)
	// :~)

	return commonDb.TxCommit
}

type refreshTargets struct {
	agentId int32
}

func (t *refreshTargets) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	/**
	 * Sets the status of cache to `2`
	 */
	tx.MustExec(
		`
		UPDATE nqm_cache_agent_ping_list_log
		SET apll_status = 2,
			apll_time_refresh = NOW()
		WHERE apll_ag_id = ?
		`,
		t.agentId,
	)
	// :~)

	/**
	 * Builds the list of targets
	 *
	 * For existing target in cache:
	 * 1) Marks it's flag to `1`
	 * 2) Update the period
	 */
	tx.MustExec(
		`
		INSERT INTO nqm_cache_agent_ping_list(
			apl_apll_ag_id, apl_tg_id, apl_min_period, apl_time_access
		)
		SELECT ?, tg_id, MIN(tg.pt_period), FROM_UNIXTIME(0) -- Use the very first time as access time
		FROM (
				/**
				 * Filter targets with:
				 *
				 * 1. Empty ping task - all enabled targets
				 * 2. Viable ping task - matched targets by filters
				 */
				SELECT tg_id, pt.pt_period
				FROM
					nqm_agent_ping_task AS apt
					INNER JOIN
					nqm_ping_task AS pt
					ON pt.pt_id = apt.apt_pt_id
						AND apt.apt_ag_id = ?
						AND pt.pt_enable = TRUE
					INNER JOIN
					vw_enabled_targets_by_ping_task AS vw_tg
					ON pt.pt_id = vw_tg.tg_pt_id
				/* :~) */
				UNION ALL
				/**
				 * Targets to be probed by all
				 *
				 * Even the agent has no ping tasks
				 */
				SELECT tg_id, -1
				FROM nqm_target tg
				WHERE tg_probed_by_all = TRUE
					AND tg.tg_status = TRUE
					AND tg.tg_available = TRUE
				/* :~) */
			) AS tg
		GROUP BY tg_id
		ON DUPLICATE KEY UPDATE
			apl_build_flag = 1,
			apl_min_period = VALUES(apl_min_period)
		`,
		t.agentId, t.agentId,
	)
	// :~)

	return commonDb.TxCommit
}

type removeToBeDeletedTargets struct {
	agentId int32
}

func (t *removeToBeDeletedTargets) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	/**
	 * Sets the status of cache to `3`
	 */
	tx.MustExec(
		`
		UPDATE nqm_cache_agent_ping_list_log
		SET apll_status = 3,
			apll_time_refresh = NOW()
		WHERE apll_ag_id = ?
		`,
		t.agentId,
	)
	// :~)

	/**
	 * Deletes the flag with value of `2`(to be deleted)
	 */
	tx.MustExec(
		`
		DELETE FROM nqm_cache_agent_ping_list
		WHERE apl_apll_ag_id = ?
			AND apl_build_flag = 2
		`,
		t.agentId,
	)
	// :~)

	return commonDb.TxCommit
}

type updateRefreshTime struct {
	agentId     int32
	checkedTime time.Time
	logObject   *nqmModel.PingListLog
}

func (t *updateRefreshTime) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	/**
	 * Updates the number of targets in cache
	 *
	 * If there is no matched target, remove the log
	 */
	tx.MustExec(
		`
		UPDATE nqm_cache_agent_ping_list_log
		SET apll_status = 0,
			apll_time_refresh = FROM_UNIXTIME(?),
			apll_number_of_targets = (
				SELECT COUNT(apl_tg_id)
				FROM nqm_cache_agent_ping_list
				WHERE apl_apll_ag_id = ?
			)
		WHERE apll_ag_id = ?
		`,
		t.checkedTime.Unix(),
		t.agentId, t.agentId,
	)
	// :~)

	/**
	 * Loads data of log object
	 */
	t.logObject = &nqmModel.PingListLog{}
	txExt := osqlx.ToTxExt(tx)
	txExt.Get(
		t.logObject,
		`
		SELECT apll_number_of_targets, apll_time_access, apll_time_refresh
		FROM nqm_cache_agent_ping_list_log
		WHERE apll_ag_id = ?
		`,
		t.agentId,
	)
	// :~)

	return commonDb.TxCommit
}

func loadAgentDetail(agentId int) *commonModel.NqmAgent {
	var dbAgentName sql.NullString
	var concatedIdsOfGroupTag sql.NullString

	loadedAgent := &commonModel.NqmAgent{}
	DbFacade.SqlDbCtrl.QueryForRow(
		commonDb.RowCallbackFunc(func(row *sql.Row) {
			loadedAgent.Id = agentId
			commonDb.ToRowExt(row).Scan(
				&dbAgentName,
				&loadedAgent.IspId, &loadedAgent.IspName,
				&loadedAgent.ProvinceId, &loadedAgent.ProvinceName,
				&loadedAgent.CityId, &loadedAgent.CityName,
				&loadedAgent.NameTagId, &concatedIdsOfGroupTag,
			)

			loadedAgent.GroupTagIds = utils.IntTo32(
				commonDb.GroupedStringToIntArray(concatedIdsOfGroupTag, ","),
			)

			/**
			 * Loads name of agent
			 */
			if dbAgentName.Valid {
				loadedAgent.Name = dbAgentName.String
			} else {
				loadedAgent.Name = commonModel.UNDEFINED_STRING
			}
			// :~)
		}),
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

	return loadedAgent
}
