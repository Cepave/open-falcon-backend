SET NAMES 'utf8';
SET @@session.default_storage_engine = 'InnoDB';

ALTER DATABASE
DEFAULT CHARACTER SET = 'utf8';

CREATE TABLE IF NOT EXISTS owl_isp(
	isp_id SMALLINT PRIMARY KEY,
	isp_name VARCHAR(64) NOT NULL,
	isp_acronym VARCHAR(16) NOT NULL,
	CONSTRAINT UNIQUE INDEX unq_owl_isp__isp_acronym
		(isp_acronym ASC),
	INDEX ix_owl_isp__isp_name
		(isp_name ASC)
);

CREATE TABLE IF NOT EXISTS owl_province(
	pv_id SMALLINT PRIMARY KEY,
	pv_name VARCHAR(64) NOT NULL,
	CONSTRAINT UNIQUE INDEX unq_owl_province__pv_name
		(pv_name ASC)
);

CREATE TABLE IF NOT EXISTS owl_city(
	ct_id SMALLINT PRIMARY KEY,
	ct_pv_id SMALLINT NOT NULL,
	ct_name VARCHAR(64) NOT NULL,
	ct_post_code VARCHAR(16) NOT NULL,
	CONSTRAINT UNIQUE INDEX unq_owl_city__pv_id_ct_name
		(ct_pv_id ASC, ct_name ASC),
	CONSTRAINT fk_owl_city__owl_province FOREIGN KEY
		(ct_pv_id) REFERENCES owl_province(pv_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT UNIQUE INDEX unq_owl_city__ct_post_code
		(ct_post_code)
);

CREATE TABLE IF NOT EXISTS nqm_agent(
	ag_id INT PRIMARY KEY AUTO_INCREMENT,
	ag_connection_id VARCHAR(128) NOT NULL,
	ag_hostname VARCHAR(256) NOT NULL,
	ag_ip_address VARBINARY(16) NOT NULL,
	ag_isp_id SMALLINT NOT NULL DEFAULT -1,
	ag_pv_id SMALLINT NOT NULL DEFAULT -1,
	ag_ct_id SMALLINT NOT NULL DEFAULT -1,
	ag_status BIT(8) NOT NULL DEFAULT b'00000001',
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
		(ag_connection_id)
);

CREATE TABLE IF NOT EXISTS nqm_target(
	tg_id INT PRIMARY KEY AUTO_INCREMENT,
	tg_name VARCHAR(128) NOT NULL,
	tg_host VARCHAR(128) NOT NULL,
	tg_isp_id SMALLINT NOT NULL DEFAULT -1,
	tg_pv_id SMALLINT NOT NULL DEFAULT -1,
	tg_ct_id SMALLINT NOT NULL DEFAULT -1,
	tg_name_tag VARCHAR(64),
	tg_probed_by_all BOOLEAN NOT NULL DEFAULT false,
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
			ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS nqm_ping_task(
	pt_ag_id INT PRIMARY KEY,
	pt_period SMALLINT NOT NULL,
	pt_time_last_execute DATETIME,
	CONSTRAINT fk_nqm_ping_task__nqm_agent FOREIGN KEY
		(pt_ag_id) REFERENCES nqm_agent(ag_id)
			ON DELETE CASCADE
			ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS nqm_pt_target_filter_name_tag(
	tfnt_pt_ag_id INT,
	tfnt_name_tag VARCHAR(64),
	CONSTRAINT pk_nqm_pt_target_filter_isp PRIMARY KEY (tfnt_pt_ag_id, tfnt_name_tag),
	CONSTRAINT fk_nqm_pt_target_filter_nt__nqm_ping_task FOREIGN KEY
		(tfnt_pt_ag_id) REFERENCES nqm_ping_task(pt_ag_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS nqm_pt_target_filter_isp(
	tfisp_pt_ag_id INT,
	tfisp_isp_id SMALLINT,
	CONSTRAINT pk_nqm_pt_target_filter_isp PRIMARY KEY (tfisp_pt_ag_id, tfisp_isp_id),
	CONSTRAINT fk_nqm_pt_target_filter_isp__nqm_ping_task FOREIGN KEY
		(tfisp_pt_ag_id) REFERENCES nqm_ping_task(pt_ag_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_pt_target_filter_isp__owl_isp FOREIGN KEY
		(tfisp_isp_id) REFERENCES owl_isp(isp_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS nqm_pt_target_filter_province(
	tfpv_pt_ag_id INT,
	tfpv_pv_id SMALLINT,
	CONSTRAINT pk_nqm_pt_target_filter_province PRIMARY KEY (tfpv_pt_ag_id, tfpv_pv_id),
	CONSTRAINT fk_nqm_pt_target_filter_province__nqm_ping_task FOREIGN KEY
		(tfpv_pt_ag_id) REFERENCES nqm_ping_task(pt_ag_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_pt_target_filter_province__owl_province FOREIGN KEY
		(tfpv_pv_id) REFERENCES owl_province(pv_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS nqm_pt_target_filter_city(
	tfct_pt_ag_id INT,
	tfct_ct_id SMALLINT,
	CONSTRAINT pk_nqm_pt_target_filter_city PRIMARY KEY (tfct_pt_ag_id, tfct_ct_id),
	CONSTRAINT fk_nqm_pt_target_filter_city__nqm_ping_task FOREIGN KEY
		(tfct_pt_ag_id) REFERENCES nqm_ping_task(pt_ag_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_pt_target_filter_city_pv__owl_city FOREIGN KEY
		(tfct_ct_id) REFERENCES owl_city(ct_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT
);
