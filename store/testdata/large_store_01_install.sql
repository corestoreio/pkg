SET sql_mode='NO_AUTO_VALUE_ON_ZERO,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';
SET NAMES utf8mb4;

SET foreign_key_checks = 0;
DROP TABLE IF EXISTS `store`;
DROP TABLE IF EXISTS `store_group`;
DROP TABLE IF EXISTS `store_website`;
SET foreign_key_checks = 1;

# Dump of table store_website
# ------------------------------------------------------------
CREATE TABLE `store_website` (
    `website_id`       SMALLINT(5) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Website Id',
    `code`             VARCHAR(32)                   NOT NULL COMMENT 'CODE',
    `name`             VARCHAR(64)                   DEFAULT NULL COMMENT 'Website NAME',
    `sort_order`       SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Sort ORDER',
    `default_group_id` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'DEFAULT GROUP Id',
    `is_default`       SMALLINT(5) UNSIGNED          DEFAULT 0 COMMENT 'Defines IS Website DEFAULT',
    PRIMARY KEY (`website_id`),
    UNIQUE KEY `STORE_WEBSITE_CODE` (`code`),
    KEY `STORE_WEBSITE_SORT_ORDER` (`sort_order`),
    KEY `STORE_WEBSITE_DEFAULT_GROUP_ID` (`default_group_id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4` COMMENT ='Websites';

LOCK TABLES `store_website` WRITE;
/*!40000 ALTER TABLE `store_website`
    DISABLE KEYS */;

INSERT INTO `store_website` (`website_id`, `code`, `name`, `sort_order`, `default_group_id`, `is_default`)
VALUES (0,'ADMIN','ADMIN', 0, 0, 0),
       (1,'world','Corporate Website', 0, 1, 1),
       (2,'us','USA Website', 1, 14, 0);

/*!40000 ALTER TABLE `store_website`
    ENABLE KEYS */;
UNLOCK TABLES;

# Dump of table store_group
# ------------------------------------------------------------
CREATE TABLE `store_group` (
    `group_id`         SMALLINT(5) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'GROUP Id',
    `website_id`       SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Website Id',
    `code`             VARCHAR(32)                   NOT NULL COMMENT 'Store GROUP UNIQUE CODE',
    `name`             VARCHAR(255)         NOT NULL COMMENT 'Store GROUP NAME',
    `root_category_id` INT(10) UNSIGNED     NOT NULL DEFAULT 0 COMMENT 'Root Category Id',
    `default_store_id` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'DEFAULT `Store` Id',
    PRIMARY KEY (`group_id`),
    UNIQUE KEY `STORE_GROUP_CODE` (`code`),
    KEY `STORE_GROUP_WEBSITE_ID` (`website_id`),
    KEY `STORE_GROUP_DEFAULT_STORE_ID` (`default_store_id`),
    CONSTRAINT `STORE_GROUP_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID` FOREIGN KEY (`website_id`) REFERENCES `store_website` (`website_id`) ON DELETE CASCADE
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4` COMMENT ='Store GROUPS';

LOCK TABLES `store_group` WRITE;
/*!40000 ALTER TABLE `store_group`
    DISABLE KEYS */;

INSERT INTO `store_group` (`group_id`, `website_id`, `code`, `name`, `root_category_id`, `default_store_id`)
VALUES (0, 0,'DEFAULT','DEFAULT', 0, 0),
       (1, 1,'english_international','English (International)', 2, 1),
       (3, 1,'french','French', 2, 2),
       (4, 1,'german','German', 2, 3),
       (5, 1,'spanish','Spanish', 2, 4),
       (6, 1,'italian','Italian', 2, 5),
       (7, 1,'portuguese','Portuguese', 2, 6),
       (8, 1,'russian','Russian', 2, 7),
       (9, 1,'japanese','Japanese', 2, 8),
       (10, 1,'simplified_chinese','Simplified Chinese', 2, 9),
       (11, 1,'traditional_chinese_hk','TRADITIONAL Chinese HK', 2, 10),
       (12, 1,'traditional_chinese_tw','TRADITIONAL Chinese TW', 2, 11),
       (13, 1,'korean','Korean', 2, 12),
       (14, 2,'english_usa','English (USA)', 2, 14);

/*!40000 ALTER TABLE `store_group`
    ENABLE KEYS */;
UNLOCK TABLES;

# Dump of table store
# ------------------------------------------------------------
CREATE TABLE `store` (
    `store_id`   SMALLINT(5) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Store Id',
    `code`       VARCHAR(32)                   NOT NULL COMMENT 'CODE',
    `website_id` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Website Id',
    `group_id`   SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'GROUP Id',
    `name`       VARCHAR(255)         NOT NULL COMMENT 'Store NAME',
    `sort_order` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Store Sort ORDER',
    `is_active`  SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Store Activity',
    PRIMARY KEY (`store_id`),
    UNIQUE KEY `STORE_CODE` (`code`),
    KEY `STORE_WEBSITE_ID` (`website_id`),
    KEY `STORE_IS_ACTIVE_SORT_ORDER` (`is_active`, `sort_order`),
    KEY `STORE_GROUP_ID` (`group_id`),
    CONSTRAINT `STORE_GROUP_ID_STORE_GROUP_GROUP_ID` FOREIGN KEY (`group_id`) REFERENCES `store_group` (`group_id`) ON DELETE CASCADE,
    CONSTRAINT `STORE_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID` FOREIGN KEY (`website_id`) REFERENCES `store_website` (`website_id`) ON DELETE CASCADE
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4` COMMENT ='Stores';

LOCK TABLES `store` WRITE;
/*!40000 ALTER TABLE `store`
    DISABLE KEYS */;

INSERT INTO `store` (`store_id`, `code`, `website_id`, `group_id`, `name`, `sort_order`, `is_active`)
VALUES (0,'ADMIN', 0, 0,'ADMIN', 0, 1),
       (1,'world_en', 1, 1,'English (International)', 0, 1),
       (2,'world_fr', 1, 3,'French', 0, 1),
       (3,'world_de', 1, 4,'German', 0, 1),
       (4,'world_es', 1, 5,'Spanish', 0, 1),
       (5,'world_it', 1, 6,'Italian', 0, 1),
       (6,'world_pt', 1, 7,'Portuguese', 0, 1),
       (7,'world_ru', 1, 8,'Russian', 0, 1),
       (8,'world_jp', 1, 9,'Japanese', 0, 1),
       (9,'world_cn', 1, 10,'Simplified Chinese', 0, 1),
       (10,'world_cn_zh', 1, 11,'TRADITIONAL Chinese HK', 0, 1),
       (11,'world_tw', 1, 12,'TRADITIONAL Chinese TW', 0, 1),
       (12,'world_ko', 1, 13,'Korean', 0, 1),
       (14,'us_en', 2, 14,'English (USA)', 0, 1);

/*!40000 ALTER TABLE `store`
    ENABLE KEYS */;
UNLOCK TABLES;
