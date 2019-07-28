SET foreign_key_checks = 0;
SET NAMES utf8mb4;

DROP TABLE IF EXISTS `store_website`;
CREATE TABLE `store_website` (
    `website_id`       SMALLINT(5) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Website Id',
    `code`             VARCHAR(32)                   DEFAULT NULL COMMENT 'Code',
    `name`             VARCHAR(64)                   DEFAULT NULL COMMENT 'Website Name',
    `sort_order`       SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Sort Order',
    `default_group_id` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Default Group Id',
    `is_default`       SMALLINT(5) UNSIGNED          DEFAULT 0 COMMENT 'Defines Is Website Default',
    PRIMARY KEY (`website_id`),
    UNIQUE KEY `STORE_WEBSITE_CODE` (`code`),
    KEY `STORE_WEBSITE_SORT_ORDER` (`sort_order`),
    KEY `STORE_WEBSITE_DEFAULT_GROUP_ID` (`default_group_id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4` COMMENT ='Websites';

DROP TABLE IF EXISTS `store_group`;
CREATE TABLE `store_group` (
    `group_id`         SMALLINT(5) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Group Id',
    `website_id`       SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Website Id',
    `name`             VARCHAR(255)         NOT NULL COMMENT 'Store Group Name',
    `root_category_id` INT(10) UNSIGNED     NOT NULL DEFAULT 0 COMMENT 'Root Category Id',
    `default_store_id` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Default Store Id',
    `code`             VARCHAR(32)                   DEFAULT NULL COMMENT 'Store group unique code',
    PRIMARY KEY (`group_id`),
    UNIQUE KEY `STORE_GROUP_CODE` (`code`),
    KEY `STORE_GROUP_WEBSITE_ID` (`website_id`),
    KEY `STORE_GROUP_DEFAULT_STORE_ID` (`default_store_id`),
    CONSTRAINT `STORE_GROUP_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID` FOREIGN KEY (`website_id`) REFERENCES `store_website` (`website_id`) ON DELETE CASCADE
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4` COMMENT ='Store Groups';

DROP TABLE IF EXISTS `store`;
CREATE TABLE `store` (
    `store_id`   SMALLINT(5) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Store Id',
    `code`       VARCHAR(32)                   DEFAULT NULL COMMENT 'Code',
    `website_id` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Website Id',
    `group_id`   SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Group Id',
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


DROP TABLE IF EXISTS `x910cms_block`;
CREATE TABLE `x910cms_block` (
    `block_id` SMALLINT(6) NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (`block_id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4`;

DROP TABLE IF EXISTS `x910cms_block_store`;
CREATE TABLE `x910cms_block_store` (
    `block_id` SMALLINT(6)          NOT NULL,
    `store_id` SMALLINT(5) UNSIGNED NOT NULL,
    PRIMARY KEY (`block_id`, `store_id`),
    KEY `CMS_BLOCK_STORE_STORE_ID` (`store_id`),
    CONSTRAINT `CMS_BLOCK_STORE_BLOCK_ID_CMS_BLOCK_BLOCK_ID` FOREIGN KEY (`block_id`) REFERENCES `x910cms_block` (`block_id`) ON DELETE CASCADE,
    CONSTRAINT `CMS_BLOCK_STORE_STORE_ID_STORE_STORE_ID` FOREIGN KEY (`store_id`) REFERENCES `store` (`store_id`) ON DELETE CASCADE
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4`;


DROP TABLE IF EXISTS `x910cms_page`;
CREATE TABLE `x910cms_page` (
    `page_id` SMALLINT(6) NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (`page_id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4`;

DROP TABLE IF EXISTS `x910cms_page_store`;
CREATE TABLE `x910cms_page_store` (
    `page_id`  SMALLINT(6)          NOT NULL,
    `store_id` SMALLINT(5) UNSIGNED NOT NULL,
    PRIMARY KEY (`page_id`, `store_id`),
    KEY `CMS_PAGE_STORE_STORE_ID` (`store_id`),
    CONSTRAINT `CMS_PAGE_STORE_PAGE_ID_CMS_PAGE_PAGE_ID` FOREIGN KEY (`page_id`) REFERENCES `x910cms_page` (`page_id`) ON DELETE CASCADE,
    CONSTRAINT `CMS_PAGE_STORE_STORE_ID_STORE_STORE_ID` FOREIGN KEY (`store_id`) REFERENCES `store` (`store_id`) ON DELETE CASCADE
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4`;

DROP TABLE IF EXISTS `x910catalog_eav_attribute`;
CREATE TABLE `x910catalog_eav_attribute` (
    `attribute_id` SMALLINT(5) UNSIGNED NOT NULL COMMENT 'Attribute ID',
    PRIMARY KEY (`attribute_id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4`;


DROP TABLE IF EXISTS `sequence_catalog_category`;
CREATE TABLE `sequence_catalog_category` (
    `sequence_value` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (`sequence_value`)
)
    ENGINE = InnoDB
    AUTO_INCREMENT = 42
    DEFAULT CHARSET = `utf8`;

-- catalog_category_entity has a 1:1 to table sequence_catalog_category and cannot be reversed
DROP TABLE IF EXISTS `catalog_category_entity`;
CREATE TABLE `catalog_category_entity` (
    `entity_id` INT(10) UNSIGNED NOT NULL COMMENT 'Entity Id',
    `row_id`    INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Version Id',
    PRIMARY KEY (`row_id`),
    KEY `CATALOG_CATEGORY_ENTITY_ENTITY_ID` (`entity_id`),
    CONSTRAINT `CAT_CTGR_ENTT_ENTT_ID_SEQUENCE_CAT_CTGR_SEQUENCE_VAL` FOREIGN KEY (`entity_id`) REFERENCES `sequence_catalog_category` (`sequence_value`) ON DELETE CASCADE
)
    ENGINE = InnoDB
    AUTO_INCREMENT = 42
    DEFAULT CHARSET = `utf8` COMMENT ='Catalog Category Table';


DROP TABLE IF EXISTS `x859admin_user`;

CREATE TABLE `x859admin_user` (
    `user_id`  INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'User ID',
    `email`    VARCHAR(128) DEFAULT NULL COMMENT 'User Email',
    `username` VARCHAR(40)  DEFAULT NULL COMMENT 'User Login',
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `ADMIN_USER_USERNAME` (`username`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4`;

DROP TABLE IF EXISTS `x859admin_passwords`;

CREATE TABLE `x859admin_passwords` (
    `password_id`   INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Password Id',
    `user_id`       INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'User Id',
    `password_hash` VARCHAR(100)              DEFAULT NULL COMMENT 'Password Hash',
    `expires`       INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Deprecated',
    `last_updated`  INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Last Updated',
    PRIMARY KEY (`password_id`),
    KEY `ADMIN_PASSWORDS_USER_ID` (`user_id`),
    CONSTRAINT `ADMIN_PASSWORDS_USER_ID_ADMIN_USER_USER_ID` FOREIGN KEY (`user_id`) REFERENCES `x859admin_user` (`user_id`) ON DELETE CASCADE
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4`;

DROP TABLE IF EXISTS `athlete_team_member`;
CREATE TABLE `athlete_team_member` (
    `id`         INT(12) UNSIGNED NOT NULL AUTO_INCREMENT,
    `team_id`    INT(10) UNSIGNED NOT NULL COMMENT 'Athlete Team ID or AID',
    `athlete_id` INT(10) UNSIGNED NOT NULL COMMENT 'Athlete ID or AGID',
    PRIMARY KEY (`id`),
    UNIQUE KEY `UNQ_ATHLETE_TEAM_MEMBER_TEAM_ATHLETE` (`team_id`, `athlete_id`),
    CONSTRAINT `FK_ATHLETE_TEAM_MEMBER_TEAM_ID` FOREIGN KEY (`team_id`) REFERENCES `athlete_team` (`team_id`) ON DELETE CASCADE,
    CONSTRAINT `FK_ATHLETE_TEAM_MEMBER_ATHLETE_ID` FOREIGN KEY (`athlete_id`) REFERENCES `athlete` (`athlete_id`) ON DELETE CASCADE
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4` COMMENT ='Athlete Team Members';

DROP TABLE IF EXISTS `athlete_team`;
CREATE TABLE `athlete_team` (
    `team_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name`    VARCHAR(340)     NOT NULL COMMENT 'Team name',
    PRIMARY KEY (`team_id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4` COMMENT ='Athlete Team';

DROP TABLE IF EXISTS `athlete`;
CREATE TABLE `athlete` (
    `athlete_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Athlete ID',
    `firstname`  VARCHAR(340) COMMENT 'First Name',
    `lastname`   VARCHAR(340) COMMENT 'Last Name',
    PRIMARY KEY (`athlete_id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4` COMMENT ='Athletes';

SET foreign_key_checks = 1;
