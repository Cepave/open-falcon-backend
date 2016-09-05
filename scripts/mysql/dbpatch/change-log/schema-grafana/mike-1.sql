/*!40101 SET NAMES utf8 */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40101 SET character_set_client = utf8 */;

CREATE TABLE IF NOT EXISTS `city` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `city` varchar(20) CHARACTER SET utf8 NOT NULL,
  `province` varchar(10) CHARACTER SET utf8 NOT NULL,
  `count` int(6) DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

CREATE TABLE IF NOT EXISTS `idc` (
  `id` int(6) NOT NULL AUTO_INCREMENT,
  `pop_id` int(6) NOT NULL,
  `name` varchar(20) CHARACTER SET utf8 NOT NULL,
  `count` int(6) DEFAULT NULL,
  `area` varchar(10) CHARACTER SET utf8 NOT NULL,
  `province` varchar(10) CHARACTER SET utf8 NOT NULL,
  `city` varchar(20) CHARACTER SET utf8 NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

CREATE TABLE IF NOT EXISTS `province` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `province` varchar(10) CHARACTER SET utf8 NOT NULL,
  `count` int(6) DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
