ALTER TABLE nqm_agent_ping_task
	DROP PRIMARY KEY,
	DROP FOREIGN KEY fk_nqm_agent_ping_task__nqm_agent,
	CHANGE COLUMN apt_ag_id apt_ag_id_old INT NULL,
	ADD COLUMN apt_ag_id INT FIRST;

ALTER TABLE nqm_agent_group_tag
	DROP PRIMARY KEY,
	DROP FOREIGN KEY fk_nqm_agent_group_tag__nqm_agent,
	CHANGE COLUMN agt_ag_id agt_ag_id_old INT NULL,
	ADD COLUMN agt_ag_id INT FIRST;

ALTER TABLE nqm_agent
	DROP PRIMARY KEY,
	CHANGE COLUMN ag_id ag_id_old INT NULL,
	AUTO_INCREMENT = 1,
	ADD COLUMN ag_id INT AUTO_INCREMENT FIRST,
	ADD CONSTRAINT PRIMARY KEY(ag_id);

UPDATE nqm_agent_ping_task AS apt,
	nqm_agent AS ag
SET apt.apt_ag_id = ag.ag_id
WHERE apt.apt_ag_id_old = ag.ag_id_old;

UPDATE nqm_agent_group_tag AS agt,
	nqm_agent AS ag
SET agt.agt_ag_id = ag.ag_id
WHERE agt.agt_ag_id_old = ag.ag_id_old;

ALTER TABLE nqm_agent_ping_task
  ADD CONSTRAINT PRIMARY KEY(apt_ag_id, apt_pt_id),
	ADD CONSTRAINT fk_nqm_agent_ping_task__nqm_agent FOREIGN KEY(apt_ag_id)
		REFERENCES nqm_agent(ag_id)
		ON DELETE CASCADE
		ON UPDATE RESTRICT;

ALTER TABLE nqm_agent_group_tag
  ADD CONSTRAINT PRIMARY KEY(agt_ag_id, agt_gt_id),
	ADD CONSTRAINT fk_nqm_agent_group_tag__nqm_agent FOREIGN KEY(agt_ag_id)
		REFERENCES nqm_agent(ag_id)
		ON DELETE CASCADE
		ON UPDATE RESTRICT;
