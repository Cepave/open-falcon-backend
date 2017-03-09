DROP DATABASE IF EXISTS `boss`;
CREATE DATABASE boss
  DEFAULT CHARACTER SET utf8
  DEFAULT COLLATE utf8_general_ci;
USE boss;
SET NAMES utf8;

DROP TABLE IF EXISTS `contacts`;
CREATE TABLE `contacts` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `name` varchar(10) CHARACTER SET utf8 NOT NULL UNIQUE,
  `phone` varchar(20) CHARACTER SET utf8 NOT NULL,
  `email` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `hosts`;
CREATE TABLE `hosts` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `hostname` varchar(30) CHARACTER SET utf8 NOT NULL UNIQUE,
  `exist` boolean DEFAULT NULL,
  `activate` boolean DEFAULT NULL,
  `platform` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  `platforms` varchar(150) CHARACTER SET utf8 DEFAULT NULL,
  `idc` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  `ip` varchar(20) CHARACTER SET utf8 NOT NULL,
  `isp` varchar(10) CHARACTER SET utf8 NOT NULL,
  `province` varchar(10) CHARACTER SET utf8 DEFAULT NULL,
  `city` varchar(20) CHARACTER SET utf8 DEFAULT NULL,
  `status` varchar(20) CHARACTER SET utf8 DEFAULT NULL,
  `bonding` int(3) DEFAULT -1,
  `speed` int(8) UNSIGNED DEFAULT NULL,
  `remark` varchar(256) CHARACTER SET utf8 DEFAULT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `idcs`;
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

DROP TABLE IF EXISTS `ips`;
CREATE TABLE `ips` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `ip` varchar(20) CHARACTER SET utf8 NOT NULL,
  `exist` boolean DEFAULT NULL,
  `status` boolean DEFAULT NULL,
  `type` varchar(10) CHARACTER SET utf8 DEFAULT NULL,
  `hostname` varchar(30) CHARACTER SET utf8 NOT NULL,
  `platform` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `nodes`;
CREATE TABLE `nodes` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `node` varchar(40) CHARACTER SET utf8 NOT NULL UNIQUE,
  `ip` varchar(20) CHARACTER SET utf8 NOT NULL UNIQUE,
  `area` varchar(25) CHARACTER SET utf8 NOT NULL,
  `province` varchar(10) CHARACTER SET utf8 NOT NULL,
  `city` varchar(15) CHARACTER SET utf8 NOT NULL,
  `idc` varchar(50) CHARACTER SET utf8 NOT NULL,
  `isp` varchar(15) CHARACTER SET utf8 NOT NULL,
  `send` tinyint(3) UNSIGNED DEFAULT NULL,
  `receive` tinyint(3) UNSIGNED DEFAULT NULL,
  `ping` float(6,2) UNSIGNED DEFAULT NULL,
  `loss` float(6,2) UNSIGNED DEFAULT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

DROP TABLE IF EXISTS `platforms`;
CREATE TABLE `platforms` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `platform` varchar(30) CHARACTER SET utf8 NOT NULL UNIQUE,
  `type` varchar(20) CHARACTER SET utf8 DEFAULT NULL,
  `visible` tinyint(1) UNSIGNED DEFAULT NULL,
  `contacts` varchar(80) CHARACTER SET utf8 DEFAULT NULL,
  `principal` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  `deputy` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  `upgrader` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  `count` int(6) DEFAULT NULL,
  `department` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  `team` varchar(30) CHARACTER SET utf8 DEFAULT NULL,
  `description` varchar(200) CHARACTER SET utf8 DEFAULT NULL,
  `updated` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
