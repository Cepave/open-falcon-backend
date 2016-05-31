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
