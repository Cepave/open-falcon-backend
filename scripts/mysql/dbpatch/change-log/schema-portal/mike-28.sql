CREATE INDEX ix_nqm_target__tg_status_tg_available
ON nqm_target(tg_status DESC, tg_available DESC);

/**
 * Filters the enabled targets with ping tasks(enabled).
 *
 * 1) This view ignores the empty ping tasks(without any filter).
 * 2) This view doesn't include the targets which are probed by all(nqm_target.tg_probed_by_all).
 */
CREATE OR REPLACE VIEW vw_enabled_targets_by_ping_task(
	tg_pt_id, tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
)
AS
SELECT DISTINCT
	ptc.pt_id,
	tg_id, tg_name, tg_host,
	tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
FROM
	/* Enabled targets */
	(
		SELECT tg_id, tg_name, tg_host,
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
