CREATE TABLE IF NOT EXISTS `resource_objects` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `object_type` int(6) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `tag_types` (
  `id` int(5) NOT NULL AUTO_INCREMENT,
  `type_name` varchar(30) NOT NULL,
  `db_table_name` varchar(30) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

REPLACE INTO `tag_types` (id, type_name, db_table_name) VALUES (1, 'string', 'str_values');
REPLACE INTO `tag_types` (id, type_name, db_table_name) VALUES (2, 'int', 'int_values');
REPLACE INTO `tag_types` (id, type_name, db_table_name) VALUES (3, 'value_model', 'vmodel_values');
REPLACE INTO `tag_types` (id, type_name, db_table_name) VALUES (4, 'description', 'description_values');


CREATE TABLE IF NOT EXISTS `tags` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `name` varchar(30) NOT NULL,
  `tag_type_id` int(5) NOT NULL,
  `description` varchar(100),
  `default` TINYINT(1) NOT NULL DEFAULT -1,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tag_name` (`name`),
  CONSTRAINT `tag_type_map_1` FOREIGN KEY (`tag_type_id`) REFERENCES `tag_types` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- 刪除object_tags會連帶一併刪除相對應的 value records
CREATE TABLE IF NOT EXISTS  `object_tags` (
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


CREATE TABLE IF NOT EXISTS `str_values` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tag_id` int(10) NOT NULL,
  `resource_object_id` int(10) NOT NULL,
  `object_tag_id` int(10) NOT NULL DEFAULT -1,
  `value` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_val_ti_roi_1` (`tag_id`, `resource_object_id`, `value`),
  CONSTRAINT `resource_object_strv_1` FOREIGN KEY (`object_tag_id`) REFERENCES `object_tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `int_values` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tag_id` int(10) NOT NULL,
  `resource_object_id` int(10) NOT NULL,
  `object_tag_id` int(10) NOT NULL DEFAULT -1,
  `value` int(6) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_val_ti_roi_2` (`tag_id`, `resource_object_id`, `value`),
  CONSTRAINT `resource_object_intv_1` FOREIGN KEY (`object_tag_id`) REFERENCES `object_tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `value_models` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tag_id` int(10) NOT NULL,
  `value` varchar(30) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tag_name` (`value`,`tag_id`),
  CONSTRAINT `tag_refer_vd_1` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS  `vmodel_values` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tag_id` int(10) NOT NULL,
  `resource_object_id` int(10) NOT NULL,
  `value_model_id` int(10) NOT NULL DEFAULT -1,
  `object_tag_id` int(10) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_val_ti_roi_3` (`tag_id`, `resource_object_id`, `value_model_id`),
  CONSTRAINT `resource_object_vmd_1` FOREIGN KEY (`object_tag_id`) REFERENCES `object_tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `description_values` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `tag_id` int(10) NOT NULL,
  `resource_object_id` int(10) NOT NULL,
  `object_tag_id` int(10) NOT NULL DEFAULT -1,
  `value` varchar(300) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_val_ti_roi_4` (`tag_id`, `resource_object_id`, `value`),
  CONSTRAINT `resource_object_desv_1` FOREIGN KEY (`object_tag_id`) REFERENCES `object_tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- insert resource_objects for hosts
INSERT `resource_objects`
 (
     id,
     object_type
 )
 select resource_object_id, '1' as object_type from `falcon_portal`.`host`;
