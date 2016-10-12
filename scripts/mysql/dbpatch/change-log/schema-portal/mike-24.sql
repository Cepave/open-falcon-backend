/**
 * OWL-924
 */
ALTER TABLE nqm_target
DROP COLUMN tg_name_tag;

ALTER TABLE nqm_pt_target_filter_name_tag
DROP COLUMN tfnt_name_tag;
-- OWL-924 :~)

/**
 * OWL-933
 */
ALTER TABLE nqm_ping_task
DROP COLUMN pt_ag_id,
DROP COLUMN pt_time_last_execute;

ALTER TABLE nqm_pt_target_filter_name_tag
DROP COLUMN tfnt_pt_ag_id;

ALTER TABLE nqm_pt_target_filter_isp
DROP COLUMN tfisp_pt_ag_id;

ALTER TABLE nqm_pt_target_filter_province
DROP COLUMN tfpv_pt_ag_id;

ALTER TABLE nqm_pt_target_filter_city
DROP COLUMN tfct_pt_ag_id;
-- OWL-933 :~)

/**
 * OWL-1138
 */
ALTER TABLE host
DROP COLUMN hs_id_old;

ALTER TABLE grp_host
DROP COLUMN gh_hs_id_old;

ALTER TABLE nqm_agent
DROP COLUMN ag_id_old;

ALTER TABLE nqm_agent_ping_task
DROP COLUMN apt_ag_id_old;
-- OWL-933 :~)
