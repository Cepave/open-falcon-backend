SET NAMES 'utf8';
SET @@session.default_storage_engine = 'InnoDB';

CREATE TABLE IF NOT EXISTS owl_name_tag
(
	nt_id SMALLINT PRIMARY KEY AUTO_INCREMENT,
	nt_value VARCHAR(64) NOT NULL,
	CONSTRAINT unq_owl_name_tag__nt_value UNIQUE INDEX (nt_value)
);

INSERT INTO owl_name_tag(nt_id, nt_value)
VALUES(-1, '<UNDEFINED>')
ON DUPLICATE KEY UPDATE
	nt_value = VALUES(nt_value);

/**
 * Builds data of name tag
 */
INSERT INTO owl_name_tag(nt_value)
SELECT DISTINCT tg_name_tag
FROM nqm_target
WHERE tg_name_tag IS NOT NULL
ON DUPLICATE KEY UPDATE
	nt_value = VALUES(nt_value);

INSERT INTO owl_name_tag(nt_value)
SELECT DISTINCT tfnt_name_tag
FROM nqm_pt_target_filter_name_tag
WHERE tfnt_name_tag IS NOT NULL
ON DUPLICATE KEY UPDATE
	nt_value = VALUES(nt_value);

ALTER TABLE nqm_agent
ADD COLUMN ag_nt_id SMALLINT NOT NULL DEFAULT -1,
ADD CONSTRAINT fk_nqm_agent__owl_name_tag
	FOREIGN KEY (ag_nt_id)
	REFERENCES owl_name_tag(nt_id)
	ON DELETE RESTRICT
	ON UPDATE RESTRICT;

ALTER TABLE nqm_target
ADD tg_nt_id SMALLINT NOT NULL DEFAULT -1,
ADD CONSTRAINT fk_nqm_target__owl_name_tag
	FOREIGN KEY (tg_nt_id)
	REFERENCES owl_name_tag(nt_id)
	ON DELETE RESTRICT
	ON UPDATE RESTRICT;

ALTER TABLE nqm_pt_target_filter_name_tag
DROP PRIMARY KEY,
MODIFY COLUMN tfnt_name_tag VARCHAR(64) NULL,
ADD tfnt_nt_id SMALLINT NOT NULL DEFAULT -1,
ADD CONSTRAINT fk_nqm_pt_target_filter_name_tag__owl_name_tag
	FOREIGN KEY (tfnt_nt_id)
	REFERENCES owl_name_tag(nt_id)
	ON DELETE RESTRICT
	ON UPDATE RESTRICT,
ADD PRIMARY KEY (tfnt_pt_ag_id, tfnt_nt_id);

UPDATE
	nqm_target AS tg,
	owl_name_tag AS nt
SET tg.tg_nt_id = nt.nt_id
WHERE tg.tg_name_tag = nt.nt_value;

UPDATE
	nqm_pt_target_filter_name_tag AS tfnt,
	owl_name_tag AS nt
SET tfnt.tfnt_nt_id = nt.nt_id
WHERE tfnt.tfnt_name_tag = nt.nt_value;
