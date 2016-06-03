SET NAMES 'utf8';

CREATE TABLE IF NOT EXISTS event_note (
  id MEDIUMINT NOT NULL AUTO_INCREMENT,
  event_caseId VARCHAR(50),
  note    VARCHAR(300),
  case_id VARCHAR(20),
  status VARCHAR(15),
  timestamp Timestamp,
  user_id int(10) unsigned,
  PRIMARY KEY (id),
  INDEX (event_caseId),
  FOREIGN KEY (event_caseId) REFERENCES event_cases(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  FOREIGN KEY (user_id) REFERENCES uic.user(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

ALTER TABLE falcon_portal.event_cases
  DROP COLUMN closed_at,
  DROP COLUMN closed_note,
  DROP COLUMN user_modified;

ALTER TABLE falcon_portal.event_cases
  ADD COLUMN process_note MEDIUMINT,
  ADD COLUMN process_status VARCHAR(20)
