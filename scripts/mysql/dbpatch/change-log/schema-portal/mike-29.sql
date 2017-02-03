CREATE TABLE nqm_cache_agent_ping_list_log(
	apll_ag_id INT PRIMARY KEY,
	apll_number_of_targets INT NOT NULL,
	apll_time_access DATETIME NOT NULL,
	apll_time_refresh DATETIME NOT NULL,
	apll_status TINYINT NOT NULL DEFAULT 0,
	CONSTRAINT FOREIGN KEY fk_nqm_cache_agent_ping_list_log__nqm_agent(apll_ag_id)
		REFERENCES nqm_agent(ag_id)
		ON DELETE CASCADE
		ON UPDATE RESTRICT
);

CREATE TABLE nqm_cache_agent_ping_list(
	apl_apll_ag_id INT,
	apl_tg_id INT,
	apl_min_period SMALLINT NOT NULL,
	apl_time_access DATETIME NOT NULL,
	apl_build_flag TINYINT NOT NULL DEFAULT 1,
	CONSTRAINT PRIMARY KEY(apl_apll_ag_id, apl_tg_id),
	CONSTRAINT FOREIGN KEY fk_nqm_cache_agent_ping_list__nqm_cache_agent_ping_list_log(apl_apll_ag_id)
		REFERENCES nqm_cache_agent_ping_list_log(apll_ag_id)
		ON DELETE CASCADE
		ON UPDATE RESTRICT,
	CONSTRAINT FOREIGN KEY fk_nqm_cache_agent_ping_list__nqm_target (apl_tg_id)
		REFERENCES nqm_target(tg_id)
		ON DELETE CASCADE
		ON UPDATE RESTRICT,
	INDEX ix_nqm_cache_agent_ping_list__apl_apll_ag_id_apl_min_period(apl_apll_ag_id, apl_min_period)
);

/**
 * Filters the enabled targets with ping tasks(enabled).
 *
 * 1) This view includes the empty ping tasks(without any filter).
 * 2) This view doesn't include the targets which are probed by all(nqm_target.tg_probed_by_all).
 * 3) A ping task might include mutiple of same targets
 */
CREATE OR REPLACE VIEW vw_enabled_targets_by_ping_task(
	tg_pt_id, tg_id
)
AS
SELECT ptc.pt_id, tg_id
FROM
	/* Enabled targets */
	(
		SELECT tg_id,
			tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id,
			tgt.tgt_gt_id
		FROM nqm_target AS tg
			LEFT OUTER JOIN
			nqm_target_group_tag AS tgt
			ON tg.tg_id = tgt.tgt_tg_id
		WHERE tg.tg_status = TRUE
			AND tg.tg_available = TRUE
	) AS tg
	INNER JOIN
	-- :~)
	/* Expension of ping tasks */
	(
		SELECT pt.pt_id,
			tfisp_isp_id AS isp_id,
			tfpv_pv_id AS pv_id,
			tfct_ct_id AS ct_id,
			tfnt_nt_id AS nt_id,
			tfgt_gt_id AS gt_id
		FROM nqm_ping_task AS pt
			LEFT OUTER JOIN
			nqm_pt_target_filter_isp AS tfisp
			ON pt.pt_id = tfisp.tfisp_pt_id
			LEFT OUTER JOIN
			nqm_pt_target_filter_province AS tfpv
			ON pt.pt_id = tfpv.tfpv_pt_id
			LEFT OUTER JOIN
			nqm_pt_target_filter_city AS tfct
			ON pt.pt_id = tfct.tfct_pt_id
			LEFT OUTER JOIN
			nqm_pt_target_filter_name_tag AS tfnt
			ON pt.pt_id = tfnt.tfnt_pt_id
			LEFT OUTER JOIN
			nqm_pt_target_filter_group_tag AS tfgt
			ON pt.pt_id = tfgt.tfgt_pt_id
		WHERE pt.pt_enable = 1
	) AS ptc
	-- :~)
	ON tg.tg_isp_id = IFNULL(ptc.isp_id, tg.tg_isp_id)
		AND tg.tg_pv_id = IFNULL(ptc.pv_id, tg.tg_pv_id)
		AND tg.tg_ct_id = IFNULL(ptc.ct_id, tg.tg_ct_id)
		AND tg.tg_nt_id = IFNULL(ptc.nt_id, tg.tg_nt_id)
		AND IFNULL(tg.tgt_gt_id, -1) = IFNULL(ptc.gt_id, IFNULL(tg.tgt_gt_id, -1));
