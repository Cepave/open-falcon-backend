SET NAMES 'utf8';

ALTER TABLE event
    ADD COLUMN closed_note VARCHAR(250),
    ADD COLUMN tpl_creator VARCHAR(64),
    DROP FOREIGN KEY event_ibfk_1,
    DROP FOREIGN KEY event_ibfk_2,
    ADD INDEX (strategy_id, template_id);
RENAME TABLE event TO event_cases;

CREATE TABLE IF NOT EXISTS event_cases(
		id VARCHAR(50),
		endpoint VARCHAR(100) NOT NULL,
		metric VARCHAR(200) NOT NULL,
		func VARCHAR(50),
		cond VARCHAR(200) NOT NULL,
		note VARCHAR(200),
		max_step int(10) unsigned,
		current_step int(10) unsigned,
		priority INT(6) NOT NULL,
		status VARCHAR(20) NOT NULL,
		timestamp Timestamp NOT NULL,
		update_at Timestamp NULL DEFAULT NULL,
		closed_at Timestamp NULL DEFAULT NULL,
		closed_note VARCHAR(250),
		user_modified int(10) unsigned,
		tpl_creator VARCHAR(64),
		expression_id int(10) unsigned,
		strategy_id int(10) unsigned,
		template_id int(10) unsigned,
		PRIMARY KEY (id),
		INDEX (endpoint, strategy_id, template_id)
)
	ENGINE =InnoDB
	DEFAULT CHARSET =utf8;

UPDATE event_cases AS e, tpl AS t
SET e.tpl_creator = t.create_user
WHERE e.template_id = t.id;

CREATE TABLE IF NOT EXISTS events (
		id MEDIUMINT NOT NULL AUTO_INCREMENT,
		event_caseId VARCHAR(50),
		step int(10) unsigned,
		cond VARCHAR(200) NOT NULL,
		timestamp Timestamp,
		PRIMARY KEY (id),
		INDEX(event_caseId),
		FOREIGN KEY (event_caseId) REFERENCES event_cases(id)
			ON DELETE CASCADE
			ON UPDATE CASCADE
)
	ENGINE =InnoDB
	DEFAULT CHARSET =utf8;


