/**
 * ping task with additional column for number of filter type
 *
 * e.x. If a ping task has isp filter and province filter,
 * 		the value of "pt_number_of_filter_types" would be "2"
 */
CREATE OR REPLACE VIEW vw_ping_task
AS
SELECT pt_id, pt_name, pt_period, pt_enable,
	pt_number_of_name_tag_filters, pt_number_of_isp_filters,
	pt_number_of_province_filters, pt_number_of_city_filters,
	pt_number_of_group_tag_filters,
	pt_comment,
	(
		IF(pt_number_of_name_tag_filters > 0, 1, 0) +
		IF(pt_number_of_isp_filters > 0, 1, 0) +
		IF(pt_number_of_province_filters > 0, 1, 0) +
		IF(pt_number_of_city_filters > 0, 1, 0) +
		IF(pt_number_of_group_tag_filters > 0, 1, 0)
	) pt_number_of_filter_types
FROM nqm_ping_task AS pt
;

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
SELECT pts.pt_id, tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
FROM (
	/* Matched target by ISP */
	SELECT pt_id, pt_number_of_filter_types,
		tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
	FROM nqm_target tg
		INNER JOIN
		nqm_pt_target_filter_isp AS tfisp
		ON tg.tg_isp_id = tfisp.tfisp_isp_id
			AND tg.tg_status = TRUE
			AND tg.tg_available = TRUE
		INNER JOIN
		vw_ping_task AS pt
		ON pt.pt_id = tfisp.tfisp_pt_id
			AND pt.pt_enable = TRUE
	/* :~) */
	UNION ALL
	/* Matched target by province */
	SELECT pt_id, pt_number_of_filter_types,
		tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
	FROM nqm_target tg
		INNER JOIN
		nqm_pt_target_filter_province AS tfpv
		ON tg.tg_pv_id = tfpv.tfpv_pv_id
			AND tg.tg_status = TRUE
			AND tg.tg_available = TRUE
		INNER JOIN
		vw_ping_task AS pt
		ON pt.pt_id = tfpv.tfpv_pt_id
			AND pt.pt_enable = TRUE
	/* :~) */
	UNION ALL
	/* Matched target by city */
	SELECT pt_id, pt_number_of_filter_types,
		tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
	FROM nqm_target tg
		INNER JOIN
		nqm_pt_target_filter_city AS tfct
		ON tg.tg_ct_id = tfct.tfct_ct_id
			AND tg.tg_status = TRUE
			AND tg.tg_available = TRUE
		INNER JOIN
		vw_ping_task AS pt
		ON pt.pt_id = tfct.tfct_pt_id
			AND pt.pt_enable = TRUE
	/* :~) */
	UNION ALL
	/* Matched target by name tag */
	SELECT pt_id, pt_number_of_filter_types,
		tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
	FROM nqm_target tg
		INNER JOIN
		nqm_pt_target_filter_name_tag AS tfnt
		ON tg.tg_nt_id = tfnt.tfnt_nt_id
			AND tg.tg_status = TRUE
			AND tg.tg_available = TRUE
		INNER JOIN
		vw_ping_task AS pt
		ON pt.pt_id = tfnt.tfnt_pt_id
			AND pt.pt_enable = TRUE
	/* :~) */
	UNION ALL
	/* Matched target by group tag */
	SELECT DISTINCT pt_id, pt_number_of_filter_types,
		tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
	FROM nqm_target tg
		INNER JOIN
		nqm_target_group_tag AS tgt
		ON tg.tg_id = tgt.tgt_tg_id
			AND tg.tg_status = TRUE
			AND tg.tg_available = TRUE
		INNER JOIN
		nqm_pt_target_filter_group_tag AS tfgt
		ON tgt.tgt_gt_id = tfgt.tfgt_gt_id
		INNER JOIN
		vw_ping_task AS pt
		ON pt.pt_id = tfgt.tfgt_pt_id
			AND pt.pt_enable = TRUE
	/* :~) */
) AS pts
GROUP BY pts.pt_id, pts.pt_number_of_filter_types,
	tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
HAVING COUNT(pts.tg_id) = pts.pt_number_of_filter_types
;
