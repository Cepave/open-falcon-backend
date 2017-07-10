package rdb

import (
	"time"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	osqlx "github.com/Cepave/open-falcon-backend/common/db/sqlx"
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"
	"github.com/jmoiron/sqlx"
)

// Gets the ping list from cache
func GetPingListFromCache(agentID int32, checkedTime time.Time) ([]*model.NqmTarget, *model.PingListLog) {
	pingListLog := getCacheLogOfPingList(agentID)

	/**
	 * Builds cache(synchronized) if there is no one
	 */
	hasCache := pingListLog != nil
	if !hasCache {
		pingListLog = BuildCacheOfPingList(agentID, checkedTime)
	}
	// :~)

	result := getPingList(agentID, checkedTime)

	/**
	 * Refresh the cache for 1 hour live
	 */
	go utils.BuildPanicCapture(
		func() {
			logger.Debugf("Update Access Time. Agent Id: [%d]", agentID)
			updateAccessTime(agentID, checkedTime)
		},
		func(p interface{}) {
			logger.Errorf("Update access time of ping task has error. Agent Id: [%d]. %v", agentID, p)
		},
	)()
	// :~)

	return result, pingListLog
}

func getCacheLogOfPingList(agentID int32) *model.PingListLog {
	listLog := &model.PingListLog{}

	hasFound := DbFacade.SqlxDbCtrl.GetOrNoRow(
		listLog,
		`
		SELECT apll_number_of_targets,
			apll_time_access, apll_time_refresh
		FROM nqm_cache_agent_ping_list_log
		WHERE apll_ag_id = ?
		`,
		agentID,
	)

	if !hasFound {
		return nil
	}

	return listLog
}

// 1. Update(or INSERT) the refresh time
// 2. Re-build the list of targets
func BuildCacheOfPingList(agentID int32, checkedTime time.Time) *model.PingListLog {
	logger.Debugf("Agent Id[%d] -> Set to be delete flags...", agentID)
	DbFacade.SqlxDbCtrl.InTx(
		&toBeDeletedTargets{agentID: agentID},
	)
	logger.Debugf("Agent Id[%d] -> Finish", agentID)

	logger.Debugf("Agent Id[%d] -> Refresh targets...", agentID)
	DbFacade.SqlxDbCtrl.InTx(
		&refreshTargets{agentID: agentID},
	)
	logger.Debugf("Agent Id[%d] -> Finish", agentID)

	logger.Debugf("Agent Id[%d] -> Remove to-be-deleted targets...", agentID)
	DbFacade.SqlxDbCtrl.InTx(
		&removeToBeDeletedTargets{agentID: agentID},
	)
	logger.Debugf("Agent Id[%d] -> Finish", agentID)

	logger.Debugf("Agent Id[%d] -> Update final status of cache log...", agentID)
	resultTx := &updateRefreshTime{
		agentID:     agentID,
		checkedTime: checkedTime,
		logObject:   &model.PingListLog{},
	}
	DbFacade.SqlxDbCtrl.InTx(resultTx)
	logger.Debugf("Agent Id[%d] -> Finish", agentID)

	return resultTx.logObject
}

type toBeDeletedTargets struct {
	agentID int32
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
		t.agentID,
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
		t.agentID,
	)
	// :~)

	return commonDb.TxCommit
}

type refreshTargets struct {
	agentID int32
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
		t.agentID,
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
		t.agentID, t.agentID,
	)
	// :~)
	tx.MustExec(
		`
		DELETE apl
		FROM nqm_cache_agent_ping_list AS apl
		    INNER JOIN
		    nqm_target
		    ON apl.apl_tg_id = tg_id
		        AND apl.apl_apll_ag_id = ?
		WHERE tg_host = (
		    SELECT INET_NTOA(CONV(HEX(ag_ip_address), 16, 10))
		    FROM nqm_agent
		    WHERE ag_id = ?
		)
		`,
		t.agentID, t.agentID,
	)

	return commonDb.TxCommit
}

type removeToBeDeletedTargets struct {
	agentID int32
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
		t.agentID,
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
		t.agentID,
	)
	// :~)

	return commonDb.TxCommit
}

type updateRefreshTime struct {
	agentID     int32
	checkedTime time.Time
	logObject   *model.PingListLog
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
		t.agentID, t.agentID,
	)
	// :~)

	/**
	 * Loads data of log object
	 */
	t.logObject = &model.PingListLog{}
	txExt := osqlx.ToTxExt(tx)
	txExt.Get(
		t.logObject,
		`
		SELECT apll_number_of_targets, apll_time_access, apll_time_refresh
		FROM nqm_cache_agent_ping_list_log
		WHERE apll_ag_id = ?
		`,
		t.agentID,
	)
	// :~)

	return commonDb.TxCommit
}

// 1. Update the access time
// 2. Retrieve the ping list
func getPingList(agentID int32, checkedTime time.Time) []*model.NqmTarget {
	var result []*model.NqmTarget

	var dbListTargets = DbFacade.GormDb.Model(&model.NqmTarget{}).
		Select(`
			apl_tg_id, tg_host,
			isp_id, isp_name,
			pv_id, pv_name,
			ct_id, ct_name,
			nt_id, nt_value,
			GROUP_CONCAT(tgt.tgt_gt_id ORDER BY tgt.tgt_gt_id ASC SEPARATOR ',') AS gts
		`).
		Joins(`
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
		`,
			agentID, checkedTime.Unix(),
		).
		Group(`
			apl_tg_id, tg_host,
			isp_id, isp_name, pv_id, pv_name, ct_id, ct_name,
			nt_id, nt_value
		`).
		Order(`
			apl_tg_id ASC
		`)

	selectNqmTargets := dbListTargets.Find(&result)
	gormExt.ToDefaultGormDbExt(selectNqmTargets).PanicIfError()

	/**
	 * Load group tags
	 */
	for _, target := range result {
		target.AfterLoad()
	}
	// :~)

	return result
}

func updateAccessTime(agentID int32, accessTime time.Time) {
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
				accessTime.Unix(), agentID,
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
				accessTime.Unix(), agentID, accessTime.Unix(),
			)
			// :~)

			return commonDb.TxCommit
		},
	))
}
