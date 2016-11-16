/**
 * 此 Patch 是重整 ID，下列兩條件成立時才可使用:
 * 1) 在 mike-21.sql, mike-22.sql 執行過之後才能使用，
 * 2) 任何 _old 欄位未移除才可以使用
 * 因為 HBS 的版本為更新，造成 ID 持續跳號的問題
 */

DELETE FROM nqm_agent
WHERE ag_hostname LIKE '%.cdn.fastweb.com.cn';

/**
 * 移除上一次重整的 ID 留存欄位
 */
ALTER TABLE nqm_agent_ping_task
DROP COLUMN apt_ag_id_old;

ALTER TABLE nqm_agent_group_tag
DROP COLUMN agt_ag_id_old;

ALTER TABLE nqm_agent
DROP COLUMN ag_id_old;

ALTER TABLE host
DROP COLUMN hs_id_old;

ALTER TABLE grp_host
DROP COLUMN gh_hs_id_old;
-- :~)

/**
 * 重整 nqm_agent.id
 */
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
-- :~)

/**
 * 重整 host.id
 */
ALTER TABLE grp_host
	DROP FOREIGN KEY fk_grp_host__host,
	DROP FOREIGN KEY fk_grp_host__grp,
	DROP KEY fk_grp_host__host,
	DROP PRIMARY KEY,
	CHANGE COLUMN host_id gh_hs_id_old INT NULL;

ALTER TABLE grp_host
	ADD COLUMN host_id INT AFTER grp_id;

ALTER TABLE host
	DROP PRIMARY KEY,
	CHANGE COLUMN `id` hs_id_old INT NULL,
	AUTO_INCREMENT = 1,
	ADD COLUMN `id` INT AUTO_INCREMENT FIRST,
	ADD CONSTRAINT PRIMARY KEY (`id`);

UPDATE host AS hs, grp_host AS gh
SET gh.host_id = hs.id
WHERE gh.gh_hs_id_old = hs.hs_id_old;

ALTER TABLE grp_host
	ADD PRIMARY KEY(grp_id, host_id),
	ADD CONSTRAINT fk_grp_host__host FOREIGN KEY(host_id)
		REFERENCES host(`id`)
		ON DELETE CASCADE
		ON UPDATE RESTRICT,
	ADD CONSTRAINT fk_grp_host__grp FOREIGN KEY(grp_id)
		REFERENCES grp(`id`)
		ON DELETE CASCADE
		ON UPDATE RESTRICT;
-- :~)
