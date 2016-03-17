DROP TABLE IF EXISTS tags;
CREATE TABLE `tags` (
  `strategy_id` INT(10) UNSIGNED NOT NULL,
  `name`     VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'tag name',
  `value`     VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'tag value',
  `create_at` DATETIME NOT NULL COMMENT 'create time',
  `update_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
  PRIMARY KEY (`strategy_id`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
