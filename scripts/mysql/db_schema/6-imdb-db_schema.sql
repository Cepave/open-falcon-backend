USE `falcon_portal`;
ALTER TABLE `host`
 ADD resource_object_id int(10);

-- update resouce_object_id to exising table --
UPDATE `falcon_portal`.`host` SET resource_object_id = id;

ALTER TABLE `grp`
  ADD objects_search_term varchar(150),
  ADD object_group_type varchar(50),
  ADD auto_sync BOOLEAN;

CREATE DATABASE IF NOT EXISTS `imdb`;

USE `imdb`;

CREATE TABLE `resource_objects` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `object_type` int(6) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `tag_types` (
  `id` int(5) NOT NULL AUTO_INCREMENT,
  `type_name` varchar(30) NOT NULL,
  `db_table_name` varchar(30) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

INSERT INTO `tag_types` (id, type_name, db_table_name) VALUES (1, 'string', 'str_values');
INSERT INTO `tag_types` (id, type_name, db_table_name) VALUES (2, 'int', 'int_values');
INSERT INTO `tag_types` (id, type_name, db_table_name) VALUES (3, 'value_model', 'vmodel_values');
INSERT INTO `tag_types` (id, type_name, db_table_name) VALUES (4, 'description', 'description_values');

CREATE TABLE `tags` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `name` varchar(30) NOT NULL,
  `tag_type_id` int(5) NOT NULL,
  `description` varchar(100),
  `default` TINYINT(1) NOT NULL DEFAULT -1,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tag_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

ALTER TABLE `tags` add CONSTRAINT `tag_type_map_1` FOREIGN KEY (`tag_type_id`) REFERENCES `tag_types` (`id`) ON DELETE CASCADE;

-- 刪除object_tags會連帶一併刪除相對應的 value records
CREATE TABLE `object_tags` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tag_id` int(10) NOT NULL,
  `resource_object_id` int(10) NOT NULL,
  `value_id` int(10) NOT NULL DEFAULT -1,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_by` varchar(30),
  PRIMARY KEY (`id`),
  CONSTRAINT `tags_ibfk_1` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE,
  CONSTRAINT `resource_object_ibfk_1` FOREIGN KEY (`resource_object_id`) REFERENCES `resource_objects` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- tag_type_id = 1
CREATE TABLE `str_values` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tag_id` int(10) NOT NULL,
  `resource_object_id` int(10) NOT NULL,
  `object_tag_id` int(10) NOT NULL DEFAULT -1,
  `value` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_val_ti_roi_1` (`tag_id`, `resource_object_id`, `value`),
  CONSTRAINT `resource_object_strv_1` FOREIGN KEY (`object_tag_id`) REFERENCES `object_tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

create index index_str_value_1 on `str_values` (value);

-- tag_type_id = 2
CREATE TABLE `int_values` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tag_id` int(10) NOT NULL,
  `resource_object_id` int(10) NOT NULL,
  `object_tag_id` int(10) NOT NULL DEFAULT -1,
  `value` int(6) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_val_ti_roi_2` (`tag_id`, `resource_object_id`, `value`),
  CONSTRAINT `resource_object_intv_1` FOREIGN KEY (`object_tag_id`) REFERENCES `object_tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- reference for value_values
-- 刪除value_models時需要手動刪除相關的object_tags, 結果會一併刪除對應的`value_models`
CREATE TABLE `value_models` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tag_id` int(10) NOT NULL,
  `value` varchar(30) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tag_name` (`value`,`tag_id`),
  CONSTRAINT `tag_refer_vd_1` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- tag_type_id = 3
-- 刪除value_models時需要手動刪除相關的object_tags, 結果會一併刪除對應的`value_models`
CREATE TABLE `vmodel_values` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tag_id` int(10) NOT NULL,
  `resource_object_id` int(10) NOT NULL,
  `value_model_id` int(10) NOT NULL DEFAULT -1,
  `object_tag_id` int(10) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_val_ti_roi_3` (`tag_id`, `resource_object_id`, `value_model_id`),
  CONSTRAINT `resource_object_vmd_1` FOREIGN KEY (`object_tag_id`) REFERENCES `object_tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- tag_type_id = 4
CREATE TABLE `description_values` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tag_id` int(10) NOT NULL,
  `resource_object_id` int(10) NOT NULL,
  `object_tag_id` int(10) NOT NULL DEFAULT -1,
  `value` varchar(300) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_val_ti_roi_4` (`tag_id`, `resource_object_id`, `value`),
  CONSTRAINT `resource_object_desv_1` FOREIGN KEY (`object_tag_id`) REFERENCES `object_tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



