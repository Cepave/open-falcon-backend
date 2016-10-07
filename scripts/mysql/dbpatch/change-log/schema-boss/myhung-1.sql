SET NAMES utf8;

CREATE TABLE IF NOT EXISTS `platforms` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `platform` varchar(30) CHARACTER SET utf8 NOT NULL UNIQUE,
  `contacts` varchar(50) CHARACTER SET utf8 DEFAULT NULL,
  `count` int(6) DEFAULT NULL,
  `updated` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

CREATE TABLE IF NOT EXISTS `hosts` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `hostname` varchar(30) CHARACTER SET utf8 NOT NULL UNIQUE,
  `exist` boolean DEFAULT NULL,
  `activate` boolean DEFAULT NULL,
  `platform` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  `idc` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  `ip` varchar(20) CHARACTER SET utf8 NOT NULL,
  `isp` varchar(10) CHARACTER SET utf8 NOT NULL,
  `province` varchar(10) CHARACTER SET utf8 DEFAULT NULL,
  `city` varchar(20) CHARACTER SET utf8 DEFAULT NULL,
  `status` varchar(20) CHARACTER SET utf8 DEFAULT NULL,
  `updated` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;


CREATE TABLE IF NOT EXISTS `contacts` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `name` varchar(10) CHARACTER SET utf8 NOT NULL UNIQUE,
  `phone` varchar(20) CHARACTER SET utf8 NOT NULL,
  `email` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  `updated` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
