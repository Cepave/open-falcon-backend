CREATE DATABASE falcon_portal;
USE falcon_portal;
SET NAMES 'utf8';

/**
 * 这里的机器是从机器管理系统中同步过来的
 * 系统拿出来单独部署需要为hbs增加功能，心跳上来的机器写入host表
 */
DROP TABLE IF EXISTS host;
CREATE TABLE host
(
  id             INT UNSIGNED NOT NULL AUTO_INCREMENT,
  hostname       VARCHAR(255) NOT NULL DEFAULT '',
  ip             VARCHAR(16)  NOT NULL DEFAULT '',
  agent_version  VARCHAR(16)  NOT NULL DEFAULT '',
  plugin_version VARCHAR(128) NOT NULL DEFAULT '',
  maintain_begin INT UNSIGNED NOT NULL DEFAULT 0,
  maintain_end   INT UNSIGNED NOT NULL DEFAULT 0,
  update_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY idx_host_hostname (hostname)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;


/**
 * 机器分组信息
 * come_from 0: 从机器管理同步过来的；1: 从页面创建的
 */
DROP TABLE IF EXISTS grp;
CREATE TABLE `grp` (
  id          INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  grp_name    VARCHAR(255)     NOT NULL DEFAULT '',
  create_user VARCHAR(64)      NOT NULL DEFAULT '',
  create_at   TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  come_from   TINYINT(4)       NOT NULL DEFAULT '0',
  PRIMARY KEY (id),
  UNIQUE KEY idx_host_grp_grp_name (grp_name)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;


DROP TABLE IF EXISTS grp_host;
CREATE TABLE grp_host
(
  grp_id  INT UNSIGNED NOT NULL,
  host_id INT UNSIGNED NOT NULL,
  KEY idx_grp_host_grp_id (grp_id),
  KEY idx_grp_host_host_id (host_id)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;


/**
 * 监控策略模板
 * tpl_name全局唯一，命名的时候可以适当带上一些前缀，比如：sa.falcon.base
 */
DROP TABLE IF EXISTS tpl;
CREATE TABLE tpl
(
  id          INT UNSIGNED NOT NULL AUTO_INCREMENT,
  tpl_name    VARCHAR(255) NOT NULL DEFAULT '',
  parent_id   INT UNSIGNED NOT NULL DEFAULT 0,
  action_id   INT UNSIGNED NOT NULL DEFAULT 0,
  create_user VARCHAR(64)  NOT NULL DEFAULT '',
  create_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY idx_tpl_name (tpl_name),
  KEY idx_tpl_create_user (create_user)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;


DROP TABLE IF EXISTS strategy;
CREATE TABLE `strategy` (
  `id`          INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `metric`      VARCHAR(128)     NOT NULL DEFAULT '',
  `tags`        VARCHAR(256)     NOT NULL DEFAULT '',
  `max_step`    INT(11)          NOT NULL DEFAULT '1',
  `priority`    TINYINT(4)       NOT NULL DEFAULT '0',
  `func`        VARCHAR(16)      NOT NULL DEFAULT 'all(#1)',
  `op`          VARCHAR(8)       NOT NULL DEFAULT '',
  `right_value` VARCHAR(64)      NOT NULL,
  `note`        VARCHAR(128)     NOT NULL DEFAULT '',
  `run_begin`   VARCHAR(16)      NOT NULL DEFAULT '',
  `run_end`     VARCHAR(16)      NOT NULL DEFAULT '',
  `tpl_id`      INT(10) UNSIGNED NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_strategy_tpl_id` (`tpl_id`)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;


DROP TABLE IF EXISTS expression;
CREATE TABLE `expression` (
  `id`          INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `expression`  VARCHAR(1024)    NOT NULL,
  `func`        VARCHAR(16)      NOT NULL DEFAULT 'all(#1)',
  `op`          VARCHAR(8)       NOT NULL DEFAULT '',
  `right_value` VARCHAR(16)      NOT NULL DEFAULT '',
  `max_step`    INT(11)          NOT NULL DEFAULT '1',
  `priority`    TINYINT(4)       NOT NULL DEFAULT '0',
  `note`        VARCHAR(1024)    NOT NULL DEFAULT '',
  `action_id`   INT(10) UNSIGNED NOT NULL DEFAULT '0',
  `create_user` VARCHAR(64)      NOT NULL DEFAULT '',
  `pause`       TINYINT(1)       NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;


DROP TABLE IF EXISTS grp_tpl;
CREATE TABLE `grp_tpl` (
  `grp_id`    INT(10) UNSIGNED NOT NULL,
  `tpl_id`    INT(10) UNSIGNED NOT NULL,
  `bind_user` VARCHAR(64)      NOT NULL DEFAULT '',
  KEY `idx_grp_tpl_grp_id` (`grp_id`),
  KEY `idx_grp_tpl_tpl_id` (`tpl_id`)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

CREATE TABLE `plugin_dir` (
  `id`          INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `grp_id`      INT(10) UNSIGNED NOT NULL,
  `dir`         VARCHAR(255)     NOT NULL,
  `create_user` VARCHAR(64)      NOT NULL DEFAULT '',
  `create_at`   TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_plugin_dir_grp_id` (`grp_id`)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;


DROP TABLE IF EXISTS action;
CREATE TABLE `action` (
  `id`                   INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `uic`                  VARCHAR(255)     NOT NULL DEFAULT '',
  `url`                  VARCHAR(255)     NOT NULL DEFAULT '',
  `callback`             TINYINT(4)       NOT NULL DEFAULT '0',
  `before_callback_sms`  TINYINT(4)       NOT NULL DEFAULT '0',
  `before_callback_mail` TINYINT(4)       NOT NULL DEFAULT '0',
  `after_callback_sms`   TINYINT(4)       NOT NULL DEFAULT '0',
  `after_callback_mail`  TINYINT(4)       NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

/**
 * nodata mock config
 */
DROP TABLE IF EXISTS `mockcfg`;
CREATE TABLE `mockcfg` (
  `id`       BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`     VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'name of mockcfg, used for uuid',
  `obj`      VARCHAR(10240) NOT NULL DEFAULT '' COMMENT 'desc of object',
  `obj_type` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'type of object, host or group or other',
  `metric`   VARCHAR(128) NOT NULL DEFAULT '',
  `tags`     VARCHAR(1024) NOT NULL DEFAULT '',
  `dstype`   VARCHAR(32)  NOT NULL DEFAULT 'GAUGE',
  `step`     INT(11) UNSIGNED  NOT NULL DEFAULT 60,
  `mock`     DOUBLE  NOT NULL DEFAULT 0  COMMENT 'mocked value when nodata occurs',
  `creator`  VARCHAR(64)  NOT NULL DEFAULT '',
  `t_create` DATETIME NOT NULL COMMENT 'create time',
  `t_modify` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS cluster;
CREATE TABLE cluster
(
  id          INT UNSIGNED   NOT NULL AUTO_INCREMENT,
  grp_id      INT            NOT NULL,
  numerator   VARCHAR(10240) NOT NULL,
  denominator VARCHAR(10240) NOT NULL,
  endpoint    VARCHAR(255)   NOT NULL,
  metric      VARCHAR(255)   NOT NULL,
  tags        VARCHAR(255)   NOT NULL,
  ds_type     VARCHAR(255)   NOT NULL,
  step        INT            NOT NULL,
  last_update TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  creator     VARCHAR(255)   NOT NULL,
  PRIMARY KEY (id)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8;

SET NAMES 'utf8';
SET @@session.default_storage_engine = 'InnoDB';

ALTER DATABASE
DEFAULT CHARACTER SET = 'utf8';

CREATE TABLE IF NOT EXISTS owl_isp(
	isp_id SMALLINT PRIMARY KEY,
	isp_name VARCHAR(64) NOT NULL,
	isp_acronym VARCHAR(16) NOT NULL,
	CONSTRAINT UNIQUE INDEX unq_owl_isp__isp_acronym
		(isp_acronym ASC),
	INDEX ix_owl_isp__isp_name
		(isp_name ASC)
)
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS owl_province(
	pv_id SMALLINT PRIMARY KEY,
	pv_name VARCHAR(64) NOT NULL,
	CONSTRAINT UNIQUE INDEX unq_owl_province__pv_name
		(pv_name ASC)
)
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS owl_city(
	ct_id SMALLINT PRIMARY KEY,
	ct_pv_id SMALLINT NOT NULL,
	ct_name VARCHAR(64) NOT NULL,
	ct_post_code VARCHAR(16) NOT NULL,
	CONSTRAINT UNIQUE INDEX unq_owl_city__pv_id_ct_name
		(ct_pv_id ASC, ct_name ASC),
	CONSTRAINT fk_owl_city__owl_province FOREIGN KEY
		(ct_pv_id) REFERENCES owl_province(pv_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT UNIQUE INDEX unq_owl_city__ct_post_code
		(ct_post_code)
)
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS nqm_agent(
	ag_id INT PRIMARY KEY AUTO_INCREMENT,
	ag_name VARCHAR(128),
	ag_connection_id VARCHAR(128) NOT NULL,
	ag_hostname VARCHAR(256) NOT NULL,
	ag_ip_address VARBINARY(16) NOT NULL,
	ag_isp_id SMALLINT NOT NULL DEFAULT -1,
	ag_pv_id SMALLINT NOT NULL DEFAULT -1,
	ag_ct_id SMALLINT NOT NULL DEFAULT -1,
	ag_status BIT(8) NOT NULL DEFAULT b'00000001',
	CONSTRAINT fk_nqm_agent__owl_isp FOREIGN KEY
		(ag_isp_id) REFERENCES owl_isp(isp_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_agent__owl_province FOREIGN KEY
		(ag_pv_id) REFERENCES owl_province(pv_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_agent__owl_city FOREIGN KEY
		(ag_ct_id) REFERENCES owl_city(ct_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT UNIQUE INDEX unq_nqm_agent__ag_connection_id
		(ag_connection_id),
	INDEX ix_nqm_agent__ag_name(ag_name)
)
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS nqm_target(
	tg_id INT PRIMARY KEY AUTO_INCREMENT,
	tg_name VARCHAR(128) NOT NULL,
	tg_host VARCHAR(128) NOT NULL,
	tg_isp_id SMALLINT NOT NULL DEFAULT -1,
	tg_pv_id SMALLINT NOT NULL DEFAULT -1,
	tg_ct_id SMALLINT NOT NULL DEFAULT -1,
	tg_name_tag VARCHAR(64),
	tg_probed_by_all BOOLEAN NOT NULL DEFAULT false,
	CONSTRAINT UNIQUE INDEX unq_nqm_target__tg_host
		(tg_host),
	INDEX ix_nqm_target__tg_name_tag
		(tg_name_tag),
	INDEX ix_nqm_target__tg_probed_by_all
		(tg_probed_by_all),
	CONSTRAINT fk_nqm_target__owl_isp FOREIGN KEY
		(tg_isp_id) REFERENCES owl_isp(isp_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_target__owl_province FOREIGN KEY
		(tg_pv_id) REFERENCES owl_province(pv_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_target__owl_city FOREIGN KEY
		(tg_ct_id) REFERENCES owl_city(ct_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT
)
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS nqm_ping_task(
	pt_ag_id INT PRIMARY KEY,
	pt_period SMALLINT NOT NULL,
	pt_time_last_execute DATETIME,
	CONSTRAINT fk_nqm_ping_task__nqm_agent FOREIGN KEY
		(pt_ag_id) REFERENCES nqm_agent(ag_id)
			ON DELETE CASCADE
			ON UPDATE CASCADE
)
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS nqm_pt_target_filter_name_tag(
	tfnt_pt_ag_id INT,
	tfnt_name_tag VARCHAR(64),
	CONSTRAINT pk_nqm_pt_target_filter_isp PRIMARY KEY (tfnt_pt_ag_id, tfnt_name_tag),
	CONSTRAINT fk_nqm_pt_target_filter_nt__nqm_ping_task FOREIGN KEY
		(tfnt_pt_ag_id) REFERENCES nqm_ping_task(pt_ag_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT
)
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS nqm_pt_target_filter_isp(
	tfisp_pt_ag_id INT,
	tfisp_isp_id SMALLINT,
	CONSTRAINT pk_nqm_pt_target_filter_isp PRIMARY KEY (tfisp_pt_ag_id, tfisp_isp_id),
	CONSTRAINT fk_nqm_pt_target_filter_isp__nqm_ping_task FOREIGN KEY
		(tfisp_pt_ag_id) REFERENCES nqm_ping_task(pt_ag_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_pt_target_filter_isp__owl_isp FOREIGN KEY
		(tfisp_isp_id) REFERENCES owl_isp(isp_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT
)
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS nqm_pt_target_filter_province(
	tfpv_pt_ag_id INT,
	tfpv_pv_id SMALLINT,
	CONSTRAINT pk_nqm_pt_target_filter_province PRIMARY KEY (tfpv_pt_ag_id, tfpv_pv_id),
	CONSTRAINT fk_nqm_pt_target_filter_province__nqm_ping_task FOREIGN KEY
		(tfpv_pt_ag_id) REFERENCES nqm_ping_task(pt_ag_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_pt_target_filter_province__owl_province FOREIGN KEY
		(tfpv_pv_id) REFERENCES owl_province(pv_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT
)
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS nqm_pt_target_filter_city(
	tfct_pt_ag_id INT,
	tfct_ct_id SMALLINT,
	CONSTRAINT pk_nqm_pt_target_filter_city PRIMARY KEY (tfct_pt_ag_id, tfct_ct_id),
	CONSTRAINT fk_nqm_pt_target_filter_city__nqm_ping_task FOREIGN KEY
		(tfct_pt_ag_id) REFERENCES nqm_ping_task(pt_ag_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT,
	CONSTRAINT fk_nqm_pt_target_filter_city_pv__owl_city FOREIGN KEY
		(tfct_ct_id) REFERENCES owl_city(ct_id)
			ON DELETE RESTRICT
			ON UPDATE RESTRICT
)
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

INSERT INTO owl_isp(isp_id, isp_name, isp_acronym)
VALUES
    (-1, '<UNDEFINED>', '<UNDEFINED>'),
    (1, 'BEIJINGSANXINSHIDAI', 'BJCIII'),
    (2, 'JIAOYUWANG', 'CERNET'),
    (3, 'YIDONG', 'CMCC'),
    (4, 'TIETONG', 'CRTC'),
    (5, 'DIANXIN', 'CTC'),
    (6, 'LIANTONG', 'CUC'),
    (7, 'DIANXINTONG', 'DXT'),
    (8, 'FANGZHENGKUANDAI', 'FBN'),
    (9, 'SHIXUNKUANDAI', 'GDCATV'),
    (10, 'GEHUAYOUXIAN', 'GHBN'),
    (11, '[ZHANG CHANG]CHENGKUANDAI', 'GWBN'),
    (12, 'GUANGZHOUeJIAKUAN', 'GZEJK'),
    (13, 'BEIJINGKUANJIE', 'KJNET'),
    (14, 'YOUXIANTONG', 'OCN'),
    (15, 'RISHENGKEJI', 'SIF'),
    (16, 'GUANGHUANXINWANG', 'SINNET'),
    (17, 'GUANGDIANWANG', 'SYCATV'),
    (18, 'TIANWEISHIXUN', 'TWSX'),
    (19, 'HUASHUKUANDAI', 'WASU'),
    (20, 'YOUTIANKUANDAI', 'YTBN'),
    (21, 'ZHONGXINWANGLUO', 'ZXNET'),
    (22, 'RIBENDIANXINDIANHUA', 'NTT'),
    (23, 'TAIWANZHONGHUADIANXIN', 'CHT'),
    (24, 'XIANGGANG[HE HUO]JITONGXUN', 'HGC'),
    (25, 'TAIWAN[DA DAI]DIANXUN', 'TWM')
ON DUPLICATE KEY UPDATE
    isp_name = VALUES(isp_name),
    isp_acronym = VALUES(isp_acronym);

INSERT INTO owl_province(pv_id, pv_name)
VALUES
    (-1, '<UNDEFINED>'),
    (1, 'NEIMENGGU'),
    (2, 'SHANXI'),
    (3, 'HEBEI'),
    (4, 'BEIJING'),
    (5, 'LIAONING'),
    (6, 'JILIN'),
    (7, 'HEILONGJIANG'),
    (8, 'SHANGHAI'),
    (9, 'JIANGSU'),
    (10, 'ANHUI'),
    (11, 'SHANDONG'),
    (12, 'TIANJIN'),
    (13, 'ZHEJIANG'),
    (14, 'JIANGXI'),
    (15, 'FUJIAN'),
    (16, '[ZHONG CHONG]QING'),
    (17, 'HUNAN'),
    (18, 'HUBEI'),
    (19, 'HENAN'),
    (20, 'GUANGDONG'),
    (21, 'GUANGXI'),
    (22, 'GUIZHOU'),
    (23, 'HAINAN'),
    (24, 'SICHUAN'),
    (25, 'YUNNAN'),
    (26, 'SHAANXI'),
    (27, 'GANSU'),
    (28, 'NINGXIA'),
    (29, 'QINGHAI'),
    (30, 'XINJIANG'),
    (31, 'XICANG'),
    (32, 'XIANGGANG'),
    (33, 'AOMEN'),
    (34, 'TAIWAN')
ON DUPLICATE KEY UPDATE
    pv_name = VALUES(pv_name);

INSERT INTO owl_city(ct_id, ct_pv_id, ct_name, ct_post_code)
VALUES
    (-1, -1, '<UNDEFINED>', '<UNDEFINED>'),
    (1, 4, 'BEIJINGSHI', '100000'),
    (2, 12, 'TIANJINSHI', '300000'),
    (3, 8, 'SHANGHAISHI', '200000'),
    (4, 16, '[ZHONG CHONG]QINGSHI', '400000'),
    (5, 20, 'GUANGZHOUSHI', '510000'),
    (6, 20, 'SHENZHENSHI', '518000'),
    (7, 20, 'DONG[GUAN WAN]SHI', '511700'),
    (8, 20, 'ZHUHAISHI', '519000'),
    (9, 20, 'SHANTOUSHI', '515000'),
    (10, 20, 'FUSHANSHI', '528000'),
    (11, 20, 'SHAOGUANSHI', '512000'),
    (12, 20, 'HEYUANSHI', '517000'),
    (13, 20, 'MEIZHOUSHI', '514000'),
    (14, 20, 'HUIZHOUSHI', '516000'),
    (15, 20, 'SHANWEISHI', '516600'),
    (16, 20, 'ZHONGSHANSHI', '528400'),
    (17, 20, 'JIANGMENSHI', '529000'),
    (18, 20, 'YANGJIANGSHI', '529500'),
    (19, 20, 'ZHANJIANGSHI', '524000'),
    (20, 20, 'MAOMINGSHI', '525000'),
    (21, 20, 'ZHAOQINGSHI', '526000'),
    (22, 20, 'QINGYUANSHI', '511500'),
    (23, 20, 'CHAOZHOUSHI', '521000'),
    (24, 20, 'JIEYANGSHI', '522000'),
    (25, 20, 'YUNFUSHI', '527300'),
    (26, 17, '[ZHANG CHANG]SHASHI', '410000'),
    (27, 17, 'YUEYANGSHI', '414000'),
    (28, 17, 'ZHANGJIAJIESHI', '427000'),
    (29, 17, 'CHANGDESHI', '415000'),
    (30, 17, 'YIYANGSHI', '413000'),
    (31, 17, 'XIANGTANSHI', '411100'),
    (32, 17, 'ZHUZHOUSHI', '412000'),
    (33, 17, 'LOUDISHI', '417000'),
    (34, 17, 'HUAIHUASHI', '418000'),
    (35, 17, 'SHAOYANGSHI', '422000'),
    (36, 17, 'HENGYANGSHI', '421000'),
    (37, 17, 'YONGZHOUSHI', '425000'),
    (38, 17, 'CHENZHOUSHI', '423000'),
    (39, 18, 'WUHANSHI', '430000'),
    (40, 18, 'SHIYANSHI', '442000'),
    (41, 18, 'XIANGFANSHI', '441000'),
    (42, 18, 'SUIZHOUSHI', '441300'),
    (43, 18, 'JINGMENSHI', '448000'),
    (44, 18, 'XIAOGANSHI', '432100'),
    (45, 18, 'YICHANGSHI', '443000'),
    (46, 18, 'HUANGGANGSHI', '438000'),
    (47, 18, 'EZHOUSHI', '436000'),
    (48, 18, 'JINGZHOUSHI', '434100'),
    (49, 18, 'HUANGSHISHI', '435000'),
    (50, 18, 'XIANNINGSHI', '437000'),
    (51, 14, 'NANCHANGSHI', '330000'),
    (52, 14, 'JIUJIANGSHI', '332000'),
    (53, 14, 'JINGDEZHENSHI', '333000'),
    (54, 14, 'SHANGRAOSHI', '334000'),
    (55, 14, 'YINGTANSHI', '335000'),
    (56, 14, 'FUZHOUSHI', '344000'),
    (57, 14, 'XINYUSHI', '336500'),
    (58, 14, 'YICHUNSHI', '336000'),
    (59, 14, 'PINGXIANGSHI', '337000'),
    (60, 14, 'JIANSHI', '343000'),
    (61, 14, 'GANZHOUSHI', '341000'),
    (62, 15, 'FUZHOUSHI', '350000'),
    (63, 15, 'NINGDESHI', '352100'),
    (64, 15, 'NANPINGSHI', '353000'),
    (65, 15, 'SANMINGSHI', '365000'),
    (66, 15, 'PUTIANSHI', '351100'),
    (67, 15, 'LONGYANSHI', '364000'),
    (68, 15, 'QUANZHOUSHI', '362000'),
    (69, 15, 'ZHANGZHOUSHI', '363000'),
    (70, 15, 'SHAMENSHI', '361000'),
    (71, 23, 'HAIKOUSHI', '570000'),
    (72, 23, 'SANYASHI', '572000'),
    (73, 24, 'CHENG[DOU DU]SHI', '610000'),
    (74, 24, 'GUANGYUANSHI', '628000'),
    (75, 24, 'BAZHONGSHI', '636000'),
    (76, 24, 'MIANYANGSHI', '621000'),
    (77, 24, 'DEYANGSHI', '618000'),
    (78, 24, 'DAZHOUSHI', '635000'),
    (79, 24, 'NANCHONGSHI', '637000'),
    (80, 24, 'SUININGSHI', '629000'),
    (81, 24, 'GUANGANSHI', '638000'),
    (82, 24, 'ZIYANGSHI', '641300'),
    (83, 24, 'MEISHANSHI', '620010'),
    (84, 24, 'YAANSHI', '625000'),
    (85, 24, 'NEIJIANGSHI', '641000'),
    (86, 24, '[LE YUE]SHANSHI', '614000'),
    (87, 24, 'ZIGONGSHI', '643000'),
    (88, 24, 'LUZHOUSHI', '646000'),
    (89, 24, 'YIBINSHI', '644000'),
    (90, 24, 'PANZHIHUASHI', '617000'),
    (91, 22, 'GUIYANGSHI', '550000'),
    (92, 22, 'ZUNYISHI', '563000'),
    (93, 22, 'LIUPANSHUISHI', '553000'),
    (94, 22, 'ANSHUNSHI', '561000'),
    (95, 25, 'KUNMINGSHI', '650000'),
    (96, 25, 'ZHAOTONGSHI', '657000'),
    (97, 25, 'LIJIANGSHI', '674100'),
    (98, 25, 'QUJINGSHI', '655000'),
    (99, 25, 'BAOSHANSHI', '678000'),
    (100, 25, 'YUXISHI', '653100'),
    (101, 25, 'LINCANGSHI', '677000'),
    (102, 25, 'PUERSHI', '665000'),
    (103, 21, 'NANNINGSHI', '530000'),
    (104, 21, 'GUILINSHI', '541000'),
    (105, 21, 'HECHISHI', '547000'),
    (106, 21, 'HEZHOUSHI', '542800'),
    (107, 21, 'LIUZHOUSHI', '545000'),
    (108, 21, 'BAISESHI', '533000'),
    (109, 21, 'LAIBINSHI', '546100'),
    (110, 21, 'WUZHOUSHI', '543000'),
    (111, 21, 'GUIGANGSHI', '537100'),
    (112, 21, 'YULINSHI', '537000'),
    (113, 21, 'CHONGZUOSHI', '532200'),
    (114, 21, 'QINZHOUSHI', '535000'),
    (115, 21, 'FANGCHENGGANGSHI', '538000'),
    (116, 21, 'BEIHAISHI', '536000'),
    (117, 31, 'LASASHI', '850000'),
    (118, 31, 'RIKAZESHI', '857000'),
    (119, 9, 'NANJINGSHI', '210000'),
    (120, 9, 'LIANYUNGANGSHI', '222000'),
    (121, 9, 'XUZHOUSHI', '221000'),
    (122, 9, '[SU XIU]QIANSHI', '223800'),
    (123, 9, 'HUAIANSHI', '223200'),
    (124, 9, 'YANCHENGSHI', '224000'),
    (125, 9, 'TAIZHOUSHI', '225300'),
    (126, 9, 'YANGZHOUSHI', '225000'),
    (127, 9, 'ZHENJIANGSHI', '212000'),
    (128, 9, 'NANTONGSHI', '226000'),
    (129, 9, 'CHANGZHOUSHI', '213000'),
    (130, 9, 'WUXISHI', '214000'),
    (131, 9, 'SUZHOUSHI', '215000'),
    (132, 13, 'HANGZHOUSHI', '310000'),
    (133, 13, 'HUZHOUSHI', '313000'),
    (134, 13, 'JIAXINGSHI', '314000'),
    (135, 13, 'SHAOXINGSHI', '312000'),
    (136, 13, 'ZHOUSHANSHI', '316000'),
    (137, 13, 'NINGBOSHI', '315000'),
    (138, 13, 'JINHUASHI', '321000'),
    (139, 13, 'QUZHOUSHI', '324000'),
    (140, 13, 'TAIZHOUSHI', '318000'),
    (141, 13, 'LISHUISHI', '323000'),
    (142, 13, 'WENZHOUSHI', '325000'),
    (143, 10, 'HEFEISHI', '230000'),
    (144, 10, 'HUAIBEISHI', '235000'),
    (145, 10, 'BOZHOUSHI', '236800'),
    (146, 10, '[SU XIU]ZHOUSHI', '234000'),
    (147, 10, '[BANG BENG]BUSHI', '233000'),
    (148, 10, 'FUYANGSHI', '236000'),
    (149, 10, 'HUAINANSHI', '232000'),
    (150, 10, 'CHUZHOUSHI', '239000'),
    (151, 10, 'LIUANSHI', '237000'),
    (152, 10, 'MAANSHANSHI', '243000'),
    (153, 10, 'CHAOHUSHI', '238000'),
    (154, 10, 'WUHUSHI', '241000'),
    (155, 10, 'XUANCHENGSHI', '242000'),
    (156, 10, 'TONGLINGSHI', '244000'),
    (157, 10, 'CHIZHOUSHI', '247000'),
    (158, 10, 'ANQINGSHI', '246000'),
    (159, 10, 'HUANGSHANSHI', '245000'),
    (160, 19, 'ZHENGZHOUSHI', '450000'),
    (161, 19, 'ANYANGSHI', '455000'),
    (162, 19, 'HEBISHI', '458000'),
    (163, 19, 'PUYANGSHI', '457000'),
    (164, 19, 'XINXIANGSHI', '453000'),
    (165, 19, 'JIAOZUOSHI', '454000'),
    (166, 19, 'SANMENXIASHI', '472000'),
    (167, 19, 'KAIFENGSHI', '475000'),
    (168, 19, 'LUOYANGSHI', '471000'),
    (169, 19, 'SHANGQIUSHI', '476000'),
    (170, 19, 'XUCHANGSHI', '461000'),
    (171, 19, 'PINGDINGSHANSHI', '467000'),
    (172, 19, 'ZHOUKOUSHI', '466000'),
    (173, 19, '[LUO TA]HESHI', '462000'),
    (174, 19, 'NANYANGSHI', '473000'),
    (175, 19, 'ZHUMADIANSHI', '463000'),
    (176, 19, 'XINYANGSHI', '464000'),
    (177, 2, 'TAIYUANSHI', '030000'),
    (178, 2, '[DA DAI]TONGSHI', '037000'),
    (179, 2, 'SHUOZHOUSHI', '038500'),
    (180, 2, 'XINZHOUSHI', '034000'),
    (181, 2, 'YANGQUANSHI', '045000'),
    (182, 2, 'JINZHONGSHI', '030600'),
    (183, 2, 'LULIANGSHI', '033000'),
    (184, 2, '[ZHANG CHANG]ZHISHI', '046000'),
    (185, 2, 'LINFENSHI', '041000'),
    (186, 2, 'JINCHENGSHI', '048000'),
    (187, 2, 'YUNCHENGSHI', '044000'),
    (188, 26, 'XIANSHI', '710000'),
    (189, 26, 'YULINSHI', '719000'),
    (190, 26, 'YANANSHI', '716000'),
    (191, 26, 'TONGCHUANSHI', '727000'),
    (192, 26, 'WEINANSHI', '714000'),
    (193, 26, 'BAOJISHI', '721000'),
    (194, 26, 'XIANYANGSHI', '712000'),
    (195, 26, 'SHANGLUOSHI', '726000'),
    (196, 26, 'HANZHONGSHI', '723000'),
    (197, 26, 'ANKANGSHI', '725000'),
    (198, 27, 'LANZHOUSHI', '730000'),
    (199, 27, 'JIAYUGUANSHI', '735100'),
    (200, 27, 'JIUQUANSHI', '735000'),
    (201, 27, 'ZHANGYESHI', '734000'),
    (202, 27, 'JINCHANGSHI', '737100'),
    (203, 27, 'WUWEISHI', '733000'),
    (204, 27, 'BAIYINSHI', '730900'),
    (205, 27, 'QINGYANGSHI', '745000'),
    (206, 27, 'PINGLIANGSHI', '744000'),
    (207, 27, 'DINGXISHI', '743000'),
    (208, 27, 'TIANSHUISHI', '741000'),
    (209, 27, 'LONGNANSHI', '746000'),
    (210, 29, 'XININGSHI', '810000'),
    (211, 28, 'YINCHUANSHI', '750000'),
    (212, 28, 'SHIZUISHANSHI', '753200'),
    (213, 28, 'WUZHONGSHI', '751100'),
    (214, 28, 'ZHONGWEISHI', '755000'),
    (215, 28, 'GUYUANSHI', '756000'),
    (216, 30, 'WULUMUQISHI', '830000'),
    (217, 30, 'KELAMAYISHI', '834000'),
    (218, 1, 'HU[HE HUO]HAOTESHI', '010000'),
    (219, 1, 'BAOTOUSHI', '014000'),
    (220, 1, 'WUHAISHI', '016000'),
    (221, 1, 'CHIFENGSHI', '024000'),
    (222, 1, 'TONGLIAOSHI', '028000'),
    (223, 1, 'EERDUOSISHI', '150600'),
    (224, 1, 'HULUNBEIERSHI', '021000'),
    (225, 1, 'BAYANNAOERSHI', '015000'),
    (226, 1, 'WULANCHABUSHI', '012000'),
    (227, 11, 'JINANSHI', '250000'),
    (228, 11, 'DEZHOUSHI', '253000'),
    (229, 11, 'BINZHOUSHI', '256600'),
    (230, 11, 'DONGYINGSHI', '257000'),
    (231, 11, 'YANTAISHI', '264000'),
    (232, 11, 'WEIHAISHI', '264200'),
    (233, 11, 'ZIBOSHI', '255000'),
    (234, 11, 'WEIFANGSHI', '261000'),
    (235, 11, 'LIAOCHENGSHI', '252000'),
    (236, 11, 'TAIANSHI', '271000'),
    (237, 11, 'LAIWUSHI', '271100'),
    (238, 11, 'QINGDAOSHI', '266000'),
    (239, 11, 'RIZHAOSHI', '276800'),
    (240, 11, 'JININGSHI', '272100'),
    (241, 11, 'HEZESHI', '274000'),
    (242, 11, 'LINYISHI', '276000'),
    (243, 11, 'ZAOZHUANGSHI', '277100'),
    (244, 3, 'SHIJIAZHUANGSHI', '050000'),
    (245, 3, 'ZHANGJIAKOUSHI', '075000'),
    (246, 3, 'CHENGDESHI', '067000'),
    (247, 3, 'TANGSHANSHI', '063000'),
    (248, 3, 'QINHUANGDAOSHI', '066000'),
    (249, 3, 'LANGFANGSHI', '065000'),
    (250, 3, 'BAODINGSHI', '071000'),
    (251, 3, 'CANGZHOUSHI', '061000'),
    (252, 3, 'HENGSHUISHI', '053000'),
    (253, 3, 'XINGTAISHI', '054000'),
    (254, 3, 'HANDANSHI', '056000'),
    (255, 7, 'HAERBINSHI', '150000'),
    (256, 7, 'HEIHESHI', '164300'),
    (257, 7, 'YICHUNSHI', '153000'),
    (258, 7, 'QIQIHAERSHI', '161000'),
    (259, 7, 'HEGANGSHI', '154100'),
    (260, 7, 'JIAMUSISHI', '154000'),
    (261, 7, 'SHUANGYASHANSHI', '155100'),
    (262, 7, 'SUIHUASHI', '152000'),
    (263, 7, '[DA DAI]QINGSHI', '163000'),
    (264, 7, 'QITAIHESHI', '154600'),
    (265, 7, 'JIXISHI', '158100'),
    (266, 7, 'MUDANJIANGSHI', '157000'),
    (267, 6, '[ZHANG CHANG]CHUNSHI', '130000'),
    (268, 6, 'BAICHENGSHI', '137000'),
    (269, 6, 'SONGYUANSHI', '138000'),
    (270, 6, 'JILINSHI', '132000'),
    (271, 6, 'SIPINGSHI', '136000'),
    (272, 6, 'LIAOYUANSHI', '136200'),
    (273, 6, 'BAISHANSHI', '134300'),
    (274, 6, 'TONGHUASHI', '134000'),
    (275, 5, '[CHEN SHEN]YANGSHI', '110000'),
    (276, 5, 'TIELINGSHI', '112000'),
    (277, 5, 'FUXINSHI', '123000'),
    (278, 5, 'FUSHUNSHI', '113000'),
    (279, 5, '[CHAO ZHAO]YANGSHI', '122000'),
    (280, 5, 'BENXISHI', '117000'),
    (281, 5, 'LIAOYANGSHI', '111000'),
    (282, 5, 'ANSHANSHI', '114000'),
    (283, 5, 'PANJINSHI', '124000'),
    (284, 5, 'JINZHOUSHI', '121000'),
    (285, 5, 'HULUDAOSHI', '125000'),
    (286, 5, 'YINGKOUSHI', '115000'),
    (287, 5, 'DANDONGSHI', '118000'),
    (288, 5, '[DA DAI]LIANSHI', '116000'),
    (289, 32, 'XIANGGANG', '999077'),
    (290, 33, 'AOMEN', '999078'),
    (291, 34, 'TAIWAN', '999079'),
    (292, 6, 'YANJISHI', '133000')
ON DUPLICATE KEY UPDATE
    ct_pv_id = VALUES(ct_pv_id),
    ct_name = VALUES(ct_name),
    ct_post_code = VALUES(ct_post_code);
