ALTER TABLE nqm_agent
	ADD COLUMN ag_name VARCHAR(128)
		AFTER ag_id,
	ADD INDEX ix_nqm_agent__ag_name (ag_name)
;
