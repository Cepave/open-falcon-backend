SET NAMES 'utf8';

CREATE TABLE IF NOT EXISTS nqm_target_class(
	tc_id SMALLINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
	tc_weight SMALLINT UNSIGNED NOT NULL,
	tc_name VARCHAR(50) NOT NULL,
	CONSTRAINT UNIQUE INDEX unq_nqm_target_class__tc_weight
		(tc_weight)
)
	DEFAULT CHARSET =utf8
	COLLATE =utf8_unicode_ci;

INSERT INTO owl_isp(isp_id, isp_name, isp_acronym)
VALUES
  (25, '台湾大电讯', 'TWM'),
	(26, '世纪互联', '21VIANET'),
	(27, '中电华通', 'CHINACOMM'),
	(28, '境外其它', 'others_offshore'),
	(29, '大陆其它', 'others_domestic'),
	(30, '科技网', 'CSTNET'),
	(31, '网通', 'CNC')
ON DUPLICATE KEY UPDATE
    isp_name = VALUES(isp_name),
    isp_acronym = VALUES(isp_acronym);

INSERT INTO owl_province(pv_id, pv_name)
VALUES
	(34, '台湾'),
	(35, '中国其它'),
	(36, '国外其它')
ON DUPLICATE KEY UPDATE
  pv_name = VALUES(pv_name);

INSERT INTO nqm_target_class(tc_weight, tc_name)
VALUES
	(0, "<UNDEFINED>"),
	(1000, "普通主机IP"),
	(3000, "普通路由IP"),
	(5000, "IDC机房主机IP"),
	(7000, "IDC机房路由IP"),
	(8000, "基调IDC主机IP"),
	(8500, "基调IDC未分类IP"),
	(9000, "基调IDC路由IP"),
	(10000, "骨干网路由IP")
ON DUPLICATE KEY UPDATE
		tc_weight = VALUES(tc_weight),
		tc_name = VALUES(tc_name);

ALTER TABLE nqm_ping_task DROP FOREIGN KEY fk_nqm_ping_task__nqm_agent;
ALTER TABLE nqm_pt_target_filter_name_tag DROP FOREIGN KEY fk_nqm_pt_target_filter_nt__nqm_ping_task;
ALTER TABLE nqm_pt_target_filter_isp DROP FOREIGN KEY fk_nqm_pt_target_filter_isp__nqm_ping_task;
ALTER TABLE nqm_pt_target_filter_isp DROP FOREIGN KEY fk_nqm_pt_target_filter_isp__owl_isp;
ALTER TABLE nqm_pt_target_filter_province DROP FOREIGN KEY fk_nqm_pt_target_filter_province__nqm_ping_task;
ALTER TABLE nqm_pt_target_filter_province DROP FOREIGN KEY fk_nqm_pt_target_filter_province__owl_province;
ALTER TABLE nqm_pt_target_filter_city DROP FOREIGN KEY fk_nqm_pt_target_filter_city__nqm_ping_task;
ALTER TABLE nqm_pt_target_filter_city DROP FOREIGN KEY fk_nqm_pt_target_filter_city_pv__owl_city;

TRUNCATE TABLE nqm_ping_task;
TRUNCATE TABLE nqm_pt_target_filter_name_tag;
TRUNCATE TABLE nqm_pt_target_filter_isp;
TRUNCATE TABLE nqm_pt_target_filter_province;
TRUNCATE TABLE nqm_pt_target_filter_city;

DROP TABLE nqm_agent;
DROP TABLE nqm_target;

CREATE TABLE IF NOT EXISTS nqm_agent(
	ag_id INT PRIMARY KEY AUTO_INCREMENT,
	ag_name VARCHAR(128),
	ag_connection_id VARCHAR(128) NOT NULL,
	ag_hostname VARCHAR(256) NOT NULL,
	ag_ip_address VARBINARY(16) NOT NULL,
	ag_isp_id SMALLINT NOT NULL DEFAULT -1,
	ag_pv_id SMALLINT NOT NULL DEFAULT -1,
	ag_ct_id SMALLINT NOT NULL DEFAULT -1,
	ag_status BOOLEAN NOT NULL DEFAULT true,
	ag_last_heartbeat DATETIME,
	CONSTRAINT fk_nqm_agent__owl_isp FOREIGN KEY
		(ag_isp_id) REFERENCES owl_isp(isp_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_agent__owl_province FOREIGN KEY
		(ag_pv_id) REFERENCES owl_province(pv_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_agent__owl_city FOREIGN KEY
		(ag_ct_id) REFERENCES owl_city(ct_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT UNIQUE INDEX unq_nqm_agent__ag_connection_id
		(ag_connection_id),
	INDEX ix_nqm_agent__ag_name(ag_name)
)
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS nqm_target(
	tg_id INT PRIMARY KEY AUTO_INCREMENT,
	tg_name VARCHAR(128) NOT NULL,
	tg_host VARCHAR(128) NOT NULL,
	tg_isp_id SMALLINT NOT NULL DEFAULT -1,
	tg_pv_id SMALLINT NOT NULL DEFAULT -1,
	tg_ct_id SMALLINT NOT NULL DEFAULT -1,
	tg_name_tag VARCHAR(64),
	tg_probed_by_all BOOLEAN NOT NULL DEFAULT false,
	tg_class_id SMALLINT UNSIGNED NOT NULL DEFAULT 1,
	tg_available BOOLEAN NOT NULL DEFAULT false,
	tg_last_result BOOLEAN NOT NULL DEFAULT false,
	tg_status BOOLEAN NOT NULL DEFAULT false,
	tg_last_probed_ts DATETIME,
	tg_created_ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT UNIQUE INDEX unq_nqm_target__tg_host
		(tg_host),
	INDEX ix_nqm_target__tg_name_tag
		(tg_name_tag),
	INDEX ix_nqm_target__tg_probed_by_all
		(tg_probed_by_all),
	CONSTRAINT fk_nqm_target__owl_isp FOREIGN KEY
		(tg_isp_id) REFERENCES owl_isp(isp_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_target__owl_province FOREIGN KEY
		(tg_pv_id) REFERENCES owl_province(pv_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_target__owl_city FOREIGN KEY
		(tg_ct_id) REFERENCES owl_city(ct_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_target__nqm_target_class FOREIGN KEY
		(tg_class_id) REFERENCES nqm_target_class(tc_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT
)
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

ALTER TABLE nqm_ping_task ADD
CONSTRAINT fk_nqm_ping_task__nqm_agent FOREIGN KEY
  (pt_ag_id) REFERENCES nqm_agent(ag_id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

ALTER TABLE nqm_pt_target_filter_name_tag ADD
CONSTRAINT fk_nqm_pt_target_filter_nt__nqm_ping_task FOREIGN KEY
  (tfnt_pt_ag_id) REFERENCES nqm_ping_task(pt_ag_id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT;

ALTER TABLE nqm_pt_target_filter_isp ADD
CONSTRAINT fk_nqm_pt_target_filter_isp__nqm_ping_task FOREIGN KEY
  (tfisp_pt_ag_id) REFERENCES nqm_ping_task(pt_ag_id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT;
ALTER TABLE nqm_pt_target_filter_isp ADD
CONSTRAINT fk_nqm_pt_target_filter_isp__owl_isp FOREIGN KEY
  (tfisp_isp_id) REFERENCES owl_isp(isp_id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT;

ALTER TABLE nqm_pt_target_filter_province ADD
CONSTRAINT fk_nqm_pt_target_filter_province__nqm_ping_task FOREIGN KEY
  (tfpv_pt_ag_id) REFERENCES nqm_ping_task(pt_ag_id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT;
ALTER TABLE nqm_pt_target_filter_province ADD
CONSTRAINT fk_nqm_pt_target_filter_province__owl_province FOREIGN KEY
  (tfpv_pv_id) REFERENCES owl_province(pv_id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT;

ALTER TABLE nqm_pt_target_filter_city ADD
CONSTRAINT fk_nqm_pt_target_filter_city__nqm_ping_task FOREIGN KEY
  (tfct_pt_ag_id) REFERENCES nqm_ping_task(pt_ag_id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT;
ALTER TABLE nqm_pt_target_filter_city ADD
CONSTRAINT fk_nqm_pt_target_filter_city_pv__owl_city FOREIGN KEY
  (tfct_ct_id) REFERENCES owl_city(ct_id)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT;
