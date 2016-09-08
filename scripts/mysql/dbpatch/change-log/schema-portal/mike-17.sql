CREATE TABLE IF NOT EXISTS owl_group_tag (
	gt_id INTEGER,
	gt_name VARCHAR(64) NOT NULL,
	CONSTRAINT pk_owl_group_tag PRIMARY KEY(gt_id),
	CONSTRAINT unq_owl_group_tag__gt_name UNIQUE INDEX(gt_name)
);

CREATE TABLE IF NOT EXISTS nqm_agent_group_tag(
	agt_ag_id INTEGER,
	agt_gt_id INTEGER,
	CONSTRAINT pk_nqm_agent_group_tag PRIMARY KEY(agt_ag_id, agt_gt_id),
	CONSTRAINT fk_nqm_agent_group_tag__nqm_agent FOREIGN KEY
		(agt_ag_id) REFERENCES nqm_agent(ag_id)
		ON DELETE CASCADE
		ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_agent_group_tag__owl_group_tag FOREIGN KEY
		(agt_gt_id) REFERENCES owl_group_tag(gt_id)
		ON DELETE RESTRICT
		ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS nqm_target_group_tag(
	tgt_tg_id INTEGER,
	tgt_gt_id INTEGER,
	CONSTRAINT pk_nqm_target_group_tag PRIMARY KEY(tgt_tg_id, tgt_gt_id),
	CONSTRAINT fk_nqm_target_group_tag__nqm_target FOREIGN KEY
		(tgt_tg_id) REFERENCES nqm_target(tg_id)
		ON DELETE CASCADE
		ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_target_group_tag__owl_group_tag FOREIGN KEY
		(tgt_gt_id) REFERENCES owl_group_tag(gt_id)
		ON DELETE RESTRICT
		ON UPDATE RESTRICT
);
