TRUNCATE TABLE `hosts`;
TRUNCATE TABLE `contacts`;
TRUNCATE TABLE `platforms`;

ALTER TABLE `hosts` ADD COLUMN `platforms` varchar(150) CHARACTER SET utf8 DEFAULT NULL;
ALTER TABLE `hosts` ADD COLUMN `remark` varchar(256) CHARACTER SET utf8 DEFAULT NULL;
ALTER TABLE `hosts` MODIFY COLUMN `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP;

ALTER TABLE `contacts` MODIFY COLUMN `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP;

ALTER TABLE `platforms` MODIFY COLUMN `contacts` varchar(80) CHARACTER SET utf8 DEFAULT NULL;
ALTER TABLE `platforms` ADD COLUMN `principal` varchar(30) CHARACTER SET utf8 DEFAULT NULL;
ALTER TABLE `platforms` ADD COLUMN `deputy` varchar(30) CHARACTER SET utf8 DEFAULT NULL;
ALTER TABLE `platforms` ADD COLUMN `upgrader` varchar(30) CHARACTER SET utf8 DEFAULT NULL;
ALTER TABLE `platforms` MODIFY COLUMN `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP;

CREATE TABLE `idcs` (
  `id` int(6) NOT NULL AUTO_INCREMENT,
  `popid` int(6) NOT NULL,
  `idc` varchar(20) CHARACTER SET utf8 NOT NULL,
  `bandwidth` int(10) DEFAULT NULL,
  `count` int(6) DEFAULT NULL,
  `area` varchar(10) CHARACTER SET utf8 NOT NULL,
  `province` varchar(10) CHARACTER SET utf8 NOT NULL,
  `city` varchar(20) CHARACTER SET utf8 NOT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

CREATE TABLE `ips` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `ip` varchar(20) CHARACTER SET utf8 NOT NULL,
  `exist` boolean DEFAULT NULL,
  `status` boolean DEFAULT NULL,
  `hostname` varchar(30) CHARACTER SET utf8 NOT NULL,
  `platform` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

CREATE TABLE `nodes` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `node` varchar(40) CHARACTER SET utf8 NOT NULL UNIQUE,
  `ip` varchar(20) CHARACTER SET utf8 NOT NULL UNIQUE,
  `area` varchar(25) CHARACTER SET utf8 NOT NULL,
  `province` varchar(10) CHARACTER SET utf8 NOT NULL,
  `city` varchar(15) CHARACTER SET utf8 NOT NULL,
  `idc` varchar(50) CHARACTER SET utf8 NOT NULL,
  `isp` varchar(15) CHARACTER SET utf8 NOT NULL,
  `ping` FLOAT DEFAULT NULL,
  `loss` FLOAT DEFAULT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

