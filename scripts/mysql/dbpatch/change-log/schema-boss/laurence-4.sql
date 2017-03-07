TRUNCATE TABLE `hosts`;
ALTER TABLE `boss`.`hosts`
  CHANGE COLUMN bonding bonding int(3) NULL DEFAULT '-1' AFTER status,
  CHANGE COLUMN status status varchar(20) NULL AFTER city,
  CHANGE COLUMN province province varchar(10) NULL AFTER isp,
  CHANGE COLUMN remark remark varchar(256) NULL AFTER speed,
  CHANGE COLUMN isp isp varchar(10) NOT NULL AFTER ip,
  CHANGE COLUMN city city varchar(20) NULL AFTER province,
  CHANGE COLUMN ip ip varchar(20) NOT NULL AFTER idc,
  CHANGE COLUMN updated updated datetime NULL on update CURRENT_TIMESTAMP AFTER remark,
  CHANGE COLUMN idc idc varchar(30) NULL AFTER platforms,
  CHANGE COLUMN platforms platforms varchar(150) NULL AFTER platform,
  CHANGE COLUMN speed speed int(8) unsigned NULL AFTER bonding;

TRUNCATE TABLE `ips`;
ALTER TABLE `boss`.`ips`
  ADD COLUMN type varchar(10) NULL AFTER status,
  CHANGE COLUMN updated updated datetime NULL on update CURRENT_TIMESTAMP AFTER platform,
  CHANGE COLUMN hostname hostname varchar(30) NOT NULL,
  CHANGE COLUMN platform platform varchar(30) NULL AFTER hostname;

TRUNCATE TABLE `nodes`;
ALTER TABLE `boss`.`nodes`
  CHANGE COLUMN loss loss float(6,2) unsigned NULL AFTER ping,
  CHANGE COLUMN ping ping float(6,2) unsigned NULL,
  ADD COLUMN receive tinyint(3) unsigned NULL,
  ADD COLUMN send tinyint(3) unsigned NULL AFTER isp,
  CHANGE COLUMN updated updated datetime NULL on update CURRENT_TIMESTAMP AFTER loss;

TRUNCATE TABLE `platforms`;
ALTER TABLE `boss`.`platforms`
  ADD COLUMN department varchar(30) NULL AFTER count,
  ADD COLUMN visible tinyint(1) unsigned NULL,
  ADD COLUMN type varchar(20) NULL AFTER platform,
  ADD COLUMN description varchar(200) NULL,
  ADD COLUMN team varchar(30) NULL,
  CHANGE COLUMN count count int(6) NULL AFTER upgrader,
  CHANGE COLUMN updated updated datetime NULL on update CURRENT_TIMESTAMP,
  CHANGE COLUMN contacts contacts varchar(80) NULL;
