TRUNCATE TABLE `ips`;
TRUNCATE TABLE `hosts`;
TRUNCATE TABLE `nodes`;

ALTER TABLE `boss`.`ips`
  ADD COLUMN type varchar(10) NULL AFTER platform;

ALTER TABLE `boss`.`hosts`
  CHANGE COLUMN remark remark varchar(512) NULL;

ALTER TABLE `boss`.`nodes`
  CHANGE COLUMN loss loss float(6,2) unsigned NULL AFTER ping,
  CHANGE COLUMN ping ping float(6,2) unsigned NULL,
  ADD COLUMN receive tinyint(3) unsigned NULL,
  ADD COLUMN send tinyint(3) unsigned NULL AFTER isp;

