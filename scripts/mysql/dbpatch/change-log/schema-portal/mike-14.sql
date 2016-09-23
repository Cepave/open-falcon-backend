/**
 * 移除 filter 的 PK 與 FK
 */
ALTER TABLE nqm_pt_target_filter_name_tag
	ADD COLUMN tfnt_pt_id INTEGER,
	DROP FOREIGN KEY fk_nqm_pt_target_filter_nt__nqm_ping_task,
	DROP PRIMARY KEY,
	MODIFY COLUMN tfnt_pt_ag_id INTEGER NULL;
ALTER TABLE nqm_pt_target_filter_isp
	ADD COLUMN tfisp_pt_id INTEGER,
	DROP FOREIGN KEY fk_nqm_pt_target_filter_isp__nqm_ping_task,
	DROP PRIMARY KEY,
	MODIFY COLUMN tfisp_pt_ag_id INTEGER NULL;
ALTER TABLE nqm_pt_target_filter_province
	ADD COLUMN tfpv_pt_id INTEGER,
	DROP FOREIGN KEY fk_nqm_pt_target_filter_province__nqm_ping_task,
	DROP PRIMARY KEY,
	MODIFY COLUMN tfpv_pt_ag_id INTEGER NULL;
ALTER TABLE nqm_pt_target_filter_city
	ADD COLUMN tfct_pt_id INTEGER,
	DROP FOREIGN KEY fk_nqm_pt_target_filter_city__nqm_ping_task,
	DROP PRIMARY KEY,
	MODIFY COLUMN tfct_pt_ag_id INTEGER NULL;

/**
 * 加入新 auto_increment 到 PK
 */
ALTER TABLE nqm_ping_task
	ADD COLUMN pt_id INTEGER AUTO_INCREMENT,
	ADD COLUMN pt_enable BOOLEAN NOT NULL DEFAULT TRUE,
	DROP FOREIGN KEY fk_nqm_ping_task__nqm_agent,
	DROP PRIMARY KEY,
	MODIFY COLUMN pt_ag_id INTEGER NULL,
	ADD CONSTRAINT pk_nqm_ping_task PRIMARY KEY(pt_id);

/**
 * 建立 agent 與 ping task 的表格
 */
CREATE TABLE IF NOT EXISTS nqm_agent_ping_task(
	apt_ag_id INTEGER,
	apt_pt_id INTEGER,
	apt_time_last_execute DATETIME,
	CONSTRAINT pk_nqm_agent_ping_task PRIMARY KEY(apt_ag_id, apt_pt_id),
	CONSTRAINT fk_nqm_agent_ping_task__nqm_agent FOREIGN KEY(apt_ag_id)
		REFERENCES nqm_agent(ag_id)
		ON UPDATE RESTRICT
		ON DELETE RESTRICT,
	CONSTRAINT fk_nqm_agent_ping_task__nqm_ping_task FOREIGN KEY(apt_pt_id)
		REFERENCES nqm_ping_task(pt_id)
		ON UPDATE RESTRICT
		ON DELETE RESTRICT
);

/**
 * 建立 agent 與 ping task 的關聯
 */
INSERT INTO nqm_agent_ping_task(apt_ag_id, apt_pt_id, apt_time_last_execute)
SELECT pt_ag_id, pt_id, pt_time_last_execute
FROM nqm_ping_task;

/**
 * 將新的 ping task id 設入到 filter
 */
UPDATE
	nqm_pt_target_filter_name_tag AS filter,
	nqm_ping_task AS pt
SET filter.tfnt_pt_id = pt.pt_id
WHERE filter.tfnt_pt_ag_id = pt.pt_ag_id;

UPDATE
	nqm_pt_target_filter_isp AS filter,
	nqm_ping_task AS pt
SET filter.tfisp_pt_id = pt.pt_id
WHERE filter.tfisp_pt_ag_id = pt.pt_ag_id;

UPDATE
	nqm_pt_target_filter_province AS filter,
	nqm_ping_task AS pt
SET filter.tfpv_pt_id = pt.pt_id
WHERE filter.tfpv_pt_ag_id = pt.pt_ag_id;

UPDATE
	nqm_pt_target_filter_city AS filter,
	nqm_ping_task AS pt
SET filter.tfct_pt_id = pt.pt_id
WHERE filter.tfct_pt_ag_id = pt.pt_ag_id;

/**
 * 將 filter 的 PK, FK 設回
 */
ALTER TABLE nqm_pt_target_filter_name_tag
	ADD CONSTRAINT pk_nqm_pt_target_filter_name_tag PRIMARY KEY
		(tfnt_pt_id, tfnt_nt_id),
	ADD CONSTRAINT fk_nqm_pt_target_filter_nt__nqm_ping_task FOREIGN KEY (tfnt_pt_id)
		REFERENCES nqm_ping_task(pt_id)
		ON UPDATE RESTRICT
		ON DELETE CASCADE;
ALTER TABLE nqm_pt_target_filter_isp
	ADD CONSTRAINT pk_nqm_pt_target_filter_isp PRIMARY KEY
		(tfisp_pt_id, tfisp_isp_id),
	ADD CONSTRAINT fk_nqm_pt_target_filter_isp__nqm_ping_task FOREIGN KEY (tfisp_pt_id)
		REFERENCES nqm_ping_task(pt_id)
		ON UPDATE RESTRICT
		ON DELETE CASCADE;
ALTER TABLE nqm_pt_target_filter_province
	ADD CONSTRAINT pk_nqm_pt_target_filter_province PRIMARY KEY
		(tfpv_pt_id, tfpv_pv_id),
	ADD CONSTRAINT fk_nqm_pt_target_filter_province__nqm_ping_task FOREIGN KEY (tfpv_pt_id)
		REFERENCES nqm_ping_task(pt_id)
		ON UPDATE RESTRICT
		ON DELETE CASCADE;
ALTER TABLE nqm_pt_target_filter_city
	ADD CONSTRAINT pk_nqm_pt_target_filter_city PRIMARY KEY
		(tfct_pt_id, tfct_ct_id),
	ADD CONSTRAINT fk_nqm_pt_target_filter_city__nqm_ping_task FOREIGN KEY (tfct_pt_id)
		REFERENCES nqm_ping_task(pt_id)
		ON UPDATE RESTRICT
		ON DELETE CASCADE;
