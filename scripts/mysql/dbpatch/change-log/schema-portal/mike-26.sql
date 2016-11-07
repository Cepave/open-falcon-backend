ALTER TABLE nqm_agent_group_tag
DROP FOREIGN KEY fk_nqm_agent_group_tag__owl_group_tag;

ALTER TABLE nqm_target_group_tag
DROP FOREIGN KEY fk_nqm_target_group_tag__owl_group_tag;

ALTER TABLE nqm_pt_target_filter_group_tag
DROP FOREIGN KEY fk_nqm_pt_target_filter_group_tag__owl_group_tag;

ALTER TABLE owl_group_tag
MODIFY COLUMN gt_id INTEGER AUTO_INCREMENT;

ALTER TABLE nqm_agent_group_tag
DROP COLUMN agt_ag_id_old,
ADD CONSTRAINT fk_nqm_agent_group_tag__owl_group_tag FOREIGN KEY
	(agt_gt_id) REFERENCES owl_group_tag(gt_id)
	ON DELETE RESTRICT
	ON UPDATE RESTRICT;

ALTER TABLE nqm_target_group_tag
ADD CONSTRAINT fk_nqm_target_group_tag__owl_group_tag FOREIGN KEY
	(tgt_gt_id) REFERENCES owl_group_tag(gt_id)
	ON DELETE RESTRICT
	ON UPDATE RESTRICT;

ALTER TABLE nqm_pt_target_filter_group_tag
ADD CONSTRAINT fk_nqm_pt_target_filter_group_tag__owl_group_tag FOREIGN KEY
	(tfgt_gt_id) REFERENCES owl_group_tag(gt_id)
	ON DELETE RESTRICT
	ON UPDATE RESTRICT;
