DROP DATABASE IF EXISTS `apollo`;
CREATE DATABASE apollo
  DEFAULT CHARACTER SET utf8
  DEFAULT COLLATE utf8_general_ci;
USE apollo;
SET NAMES utf8;

DROP TABLE IF EXISTS `deviations`;
CREATE TABLE `deviations` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `date` datetime NOT NULL,
  `platform` varchar(30) CHARACTER SET utf8 NOT NULL,
  `metric` tinyint (1) NOT NULL,
  `samples` int(3) NOT NULL,
  `mean` bigint UNSIGNED DEFAULT NULL,
  `deviation` bigint UNSIGNED DEFAULT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `net`;
CREATE TABLE `net` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `date` datetime NOT NULL,
  `hour` int(2) NOT NULL,
  `minute` int(2) NOT NULL,
  `platform` varchar(30) CHARACTER SET utf8 NOT NULL,
  `metric` tinyint (1) NOT NULL,
  `count` int(5) UNSIGNED DEFAULT NULL,
  `bits` bigint UNSIGNED DEFAULT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `remarks`;
CREATE TABLE `remarks` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `hostname` varchar(30) CHARACTER SET utf8 NOT NULL,
  `remark` varchar(512) CHARACTER SET utf8 DEFAULT NULL,
  `userid` int(10) UNSIGNED DEFAULT NULL,
  `account` varchar(64) CHARACTER SET utf8 NOT NULL,
  `name` varchar(128) CHARACTER SET utf8 DEFAULT NULL,
  `email` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
