SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS `store`;
DROP TABLE IF EXISTS `store_group`;
DROP TABLE IF EXISTS `store_website`;

CREATE TABLE `store_website` (
    `website_id`       SMALLINT(5) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Website ID',
    `code`             VARCHAR(64)                   DEFAULT NULL COMMENT 'Code',
    `name`             VARCHAR(128)                  DEFAULT NULL COMMENT 'Website Name',
    `sort_order`       SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Sort Order',
    `default_group_id` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Default Group ID',
    `is_default`       SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Defines Is Website Default',
    PRIMARY KEY (`website_id`),
    UNIQUE KEY `STORE_WEBSITE_CODE` (`code`),
    KEY `STORE_WEBSITE_SORT_ORDER` (`sort_order`),
    KEY `STORE_WEBSITE_DEFAULT_GROUP_ID` (`default_group_id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4` COMMENT ='Websites';

CREATE TABLE `store_group` (
    `group_id`         SMALLINT(5) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Group ID',
    `website_id`       SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Website ID',
    `name`             VARCHAR(255)         NOT NULL COMMENT 'Store Group Name',
    `root_category_id` INT(10) UNSIGNED     NOT NULL DEFAULT 0 COMMENT 'Root Category ID',
    `default_store_id` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Default Store ID',
    `code`             VARCHAR(64)                   DEFAULT NULL COMMENT 'Store group unique code',
    PRIMARY KEY (`group_id`),
    UNIQUE KEY `STORE_GROUP_CODE` (`code`),
    KEY `STORE_GROUP_WEBSITE_ID` (`website_id`),
    KEY `STORE_GROUP_DEFAULT_STORE_ID` (`default_store_id`),
    CONSTRAINT `STORE_GROUP_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID` FOREIGN KEY (`website_id`) REFERENCES `store_website` (`website_id`) ON DELETE CASCADE
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4` COMMENT ='Store Groups';

CREATE TABLE `store` (
    `store_id`   SMALLINT(5) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Store ID',
    `code`       VARCHAR(64)                   DEFAULT NULL COMMENT 'Code',
    `website_id` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Website ID',
    `group_id`   SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Group ID',
    `name`       VARCHAR(255)         NOT NULL COMMENT 'Store Name',
    `sort_order` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Store Sort Order',
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

SET FOREIGN_KEY_CHECKS = 1;
