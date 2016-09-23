CREATE TABLE nqm_pt_target_filter_group_tag(
	tfgt_pt_id INT,
	tfgt_gt_id INT,
	CONSTRAINT pk_nqm_pt_target_filter_group_tag PRIMARY KEY(tfgt_pt_id, tfgt_gt_id),
	CONSTRAINT fk_nqm_pt_target_filter_group_tag__nqm_ping_task FOREIGN KEY
		(tfgt_pt_id) REFERENCES nqm_ping_task(pt_id)
		ON DELETE CASCADE
		ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_pt_target_filter_group_tag__owl_group_tag FOREIGN KEY
		(tfgt_gt_id) REFERENCES owl_group_tag(gt_id)
		ON DELETE RESTRICT
		ON UPDATE RESTRICT
);

/**
 * Adds column for number of filters
 */
ALTER TABLE nqm_ping_task
	ADD COLUMN pt_number_of_group_tag_filters SMALLINT NOT NULL DEFAULT 0;

/**
 * Get called by various trigger
 */
DROP PROCEDURE IF EXISTS proc_ping_task_refresh_number_of_filters;
CREATE PROCEDURE proc_ping_task_refresh_number_of_filters(
	IN ping_task_id INTEGER
)
BEGIN
	UPDATE nqm_ping_task AS pt
	SET pt_number_of_name_tag_filters = (
			SELECT COUNT(tfnt_nt_id)
			FROM nqm_pt_target_filter_name_tag
			WHERE tfnt_pt_id = ping_task_id
		),
		pt_number_of_isp_filters = (
			SELECT COUNT(tfisp_isp_id)
			FROM nqm_pt_target_filter_isp
			WHERE tfisp_pt_id = ping_task_id
		),
		pt_number_of_province_filters = (
			SELECT COUNT(tfpv_pv_id)
			FROM nqm_pt_target_filter_province
			WHERE tfpv_pt_id = ping_task_id
		),
		pt_number_of_city_filters = (
			SELECT COUNT(tfct_ct_id)
			FROM nqm_pt_target_filter_city
			WHERE tfct_pt_id = ping_task_id
		),
		pt_number_of_group_tag_filters = (
			SELECT COUNT(tfgt_gt_id)
			FROM nqm_pt_target_filter_group_tag
			WHERE tfgt_pt_id = ping_task_id
		)
	WHERE pt.pt_id = ping_task_id;;
END;

CREATE TRIGGER tri_after_insert__nqm_pt_target_filter_group_tag
AFTER INSERT on nqm_pt_target_filter_group_tag
FOR EACH ROW
BEGIN
	CALL proc_ping_task_refresh_number_of_filters(NEW.tfgt_pt_id);;
END;

CREATE TRIGGER tri_after_delete__nqm_pt_target_filter_group_tag
AFTER DELETE on nqm_pt_target_filter_group_tag
FOR EACH ROW
BEGIN
	CALL proc_ping_task_refresh_number_of_filters(OLD.tfgt_pt_id);;
END;

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
/* Matched target by ISP */
SELECT pt_id, tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
FROM nqm_target tg
	INNER JOIN
	nqm_pt_target_filter_isp AS tfisp
	ON tg.tg_isp_id = tfisp.tfisp_isp_id
		AND tg.tg_status = TRUE
		AND tg.tg_available = TRUE
	INNER JOIN
	nqm_ping_task AS pt
	ON pt.pt_id = tfisp.tfisp_pt_id
		AND pt.pt_enable = TRUE
/* :~) */
UNION
/* Matched target by province */
SELECT pt_id, tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
FROM nqm_target tg
	INNER JOIN
	nqm_pt_target_filter_province AS tfpv
	ON tg.tg_pv_id = tfpv.tfpv_pv_id
		AND tg.tg_status = TRUE
		AND tg.tg_available = TRUE
	INNER JOIN
	nqm_ping_task AS pt
	ON pt.pt_id = tfpv.tfpv_pt_id
		AND pt.pt_enable = TRUE
/* :~) */
UNION
/* Matched target by city */
SELECT pt_id, tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
FROM nqm_target tg
	INNER JOIN
	nqm_pt_target_filter_city AS tfct
	ON tg.tg_ct_id = tfct.tfct_ct_id
		AND tg.tg_status = TRUE
		AND tg.tg_available = TRUE
	INNER JOIN
	nqm_ping_task AS pt
	ON pt.pt_id = tfct.tfct_pt_id
		AND pt.pt_enable = TRUE
/* :~) */
UNION
/* Matched target by name tag */
SELECT pt_id, tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
FROM nqm_target tg
	INNER JOIN
	nqm_pt_target_filter_name_tag AS tfnt
	ON tg.tg_nt_id = tfnt.tfnt_nt_id
		AND tg.tg_status = TRUE
		AND tg.tg_available = TRUE
	INNER JOIN
	nqm_ping_task AS pt
	ON pt.pt_id = tfnt.tfnt_pt_id
		AND pt.pt_enable = TRUE
/* :~) */
UNION
/* Matched target by group tag */
SELECT pt_id, tg_id, tg_name, tg_host, tg_isp_id, tg_pv_id, tg_ct_id, tg_nt_id
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
	nqm_ping_task AS pt
	ON pt.pt_id = tfgt.tfgt_pt_id
		AND pt.pt_enable = TRUE
/* :~) */
;
