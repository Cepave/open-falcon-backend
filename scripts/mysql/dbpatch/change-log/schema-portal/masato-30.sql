SET NAMES 'utf8';

use falcon_portal

CREATE TABLE `alarm_types` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `internal_data` TINYINT DEFAULT 1,
  `description` varchar(255) NOT NULL DEFAULT '',
  `color` varchar(20) NOT NULL DEFAULT 'black';
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8;

INSERT INTO `alarm_types` (id, name, internal_data, description) VALUES (1, 'owl', 1, 'blue', 'default type of owl');

ALTER TABLE `event_cases`
  ADD COLUMN alarm_type_id int(10) unsigned NOT NULL DEFAULT 1,
  ADD COLUMN ip VARCHAR(50),
  ADD COLUMN idc VARCHAR(50),
  ADD COLUMN platform VARCHAR(50),
  ADD COLUMN contact VARCHAR(20),
  ADD COLUMN extended_blob VARCHAR(255);

ALTER TABLE `event_cases`
  ADD INDEX alarm_id_index(`alarm_type_id`),
  ADD CONSTRAINT FK_alarm_types
  FOREIGN KEY (alarm_type_id) REFERENCES alarm_types(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;
