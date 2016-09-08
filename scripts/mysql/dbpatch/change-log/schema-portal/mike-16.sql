ALTER TABLE nqm_agent
	ADD COLUMN ag_comment VARCHAR(2048);

ALTER TABLE nqm_target
	ADD COLUMN tg_comment VARCHAR(2048);

ALTER TABLE nqm_ping_task
	ADD COLUMN pt_comment VARCHAR(2048);
