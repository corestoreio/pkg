SET foreign_key_checks = 0;

DROP TABLE IF EXISTS `dmlgen_types`;
DROP TABLE IF EXISTS `store`;
DROP TABLE IF EXISTS `store_group`;
DROP TABLE IF EXISTS `store_website`;
DROP TABLE IF EXISTS `customer_entity`;
DROP TABLE IF EXISTS `customer_address_entity`;
DROP TABLE IF EXISTS `core_configuration`;
DROP VIEW IF EXISTS `view_customer_no_auto_increment`;
DROP VIEW IF EXISTS `view_customer_auto_increment`;
DROP TABLE IF EXISTS `catalog_product_index_eav_decimal_idx`;
DROP TABLE IF EXISTS `sales_order_status_state`;

CREATE TABLE `dmlgen_types`
(
    `id`              INT(11)              NOT NULL AUTO_INCREMENT,
    col_bigint_1      BIGINT(20)           NULL,
    col_bigint_2      BIGINT(20)           NOT NULL DEFAULT 0,
    col_bigint_3      BIGINT(20) UNSIGNED  NULL,
    col_bigint_4      BIGINT(20) UNSIGNED  NOT NULL DEFAULT 0,
    col_blob          BLOB                          DEFAULT NULL,
    col_date_1        DATE                          DEFAULT NULL,
    col_date_2        DATE                 NOT NULL DEFAULT '0000-00-00',
    col_datetime_1    DATETIME                      DEFAULT NULL,
    col_datetime_2    DATETIME             NOT NULL DEFAULT '0000-00-00 00:00:00',
    col_decimal_10_1  DECIMAL(10, 1) UNSIGNED       DEFAULT NULL,
    col_decimal_12_4  DECIMAL(12, 4)                DEFAULT NULL,
    price_a_12_4      DECIMAL(12, 4)                DEFAULT NULL,
    price_b_12_4      DECIMAL(12, 4)       NOT NULL DEFAULT 0,
    col_decimal_12_3  DECIMAL(12, 3)       NOT NULL DEFAULT 0,
    col_decimal_20_6  DECIMAL(20, 6)       NOT NULL DEFAULT 0.000000,
    col_decimal_24_12 DECIMAL(24, 12)      NOT NULL DEFAULT 0.000000000000,
    /* col_float                 FLOAT                NOT NULL           DEFAULT 1, hard to compare, so skipping */
    col_int_1         INT(10)              NULL,
    col_int_2         INT(10)              NOT NULL DEFAULT 0,
    col_int_3         INT(10) UNSIGNED     NULL,
    col_int_4         INT(10) UNSIGNED     NOT NULL DEFAULT 0,
    col_longtext_1    LONGTEXT                      DEFAULT NULL,
    col_longtext_2    LONGTEXT             NOT NULL DEFAULT '',
    col_mediumblob    MEDIUMBLOB                    DEFAULT NULL,
    col_mediumtext_1  MEDIUMTEXT                    DEFAULT NULL,
    col_mediumtext_2  MEDIUMTEXT           NOT NULL DEFAULT '',
    col_smallint_1    SMALLINT(5)          NULL,
    col_smallint_2    SMALLINT(5)          NOT NULL DEFAULT 0,
    col_smallint_3    SMALLINT(5) UNSIGNED NULL,
    col_smallint_4    SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0,
    has_smallint_5    SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0,
    is_smallint_5     SMALLINT(5)          NULL,
    col_text          TEXT                          DEFAULT NULL,
    col_timestamp_1   TIMESTAMP            NOT NULL DEFAULT current_timestamp(),
    col_timestamp_2   TIMESTAMP            NULL,
    col_tinyint_1     TINYINT(1)           NOT NULL DEFAULT 0,
    col_varchar_1     VARCHAR(1)           NOT NULL DEFAULT '0',
    col_varchar_100   VARCHAR(100)                  DEFAULT NULL,
    col_varchar_16    VARCHAR(16)          NOT NULL DEFAULT 'de_DE',
    col_char_1        CHAR(21)                      DEFAULT NULL,
    col_char_2        CHAR(17)             NOT NULL DEFAULT 'xchar',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

-- Create syntax for TABLE 'store'
CREATE TABLE `store`
(
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
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='Stores';

-- Create syntax for TABLE 'store_group'
CREATE TABLE `store_group`
(
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
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='Store Groups';

-- Create syntax for TABLE 'store_website'
CREATE TABLE `store_website`
(
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
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='Websites';

-- Create syntax for TABLE 'customer_entity'
CREATE TABLE `customer_entity`
(
    `entity_id`           INT(10) UNSIGNED     NOT NULL AUTO_INCREMENT COMMENT 'Entity ID',
    `website_id`          SMALLINT(5) UNSIGNED          DEFAULT NULL COMMENT 'Website ID',
    `email`               VARCHAR(255)                  DEFAULT NULL COMMENT 'Email',
    `group_id`            SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Group ID',
    `store_id`            SMALLINT(5) UNSIGNED          DEFAULT 0 COMMENT 'Store ID',
    `created_at`          TIMESTAMP            NOT NULL DEFAULT current_timestamp() COMMENT 'Created At',
    `updated_at`          TIMESTAMP            NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Updated At',
    `is_active`           SMALLINT(5) UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Is Active',
    `created_in`          VARCHAR(255)                  DEFAULT NULL COMMENT 'Created From',
    `firstname`           VARCHAR(255)                  DEFAULT NULL COMMENT 'First Name',
    `lastname`            VARCHAR(255)                  DEFAULT NULL COMMENT 'Last Name',
    `dob`                 DATE                          DEFAULT NULL COMMENT 'Date of Birth',
    `password_hash`       VARCHAR(128)                  DEFAULT NULL COMMENT 'Password_hash',
    `rp_token`            VARCHAR(128)                  DEFAULT NULL COMMENT 'Reset password token',
    `rp_token_created_at` DATETIME                      DEFAULT NULL COMMENT 'Reset password token creation time',
    `default_billing`     INT(10) UNSIGNED              DEFAULT NULL COMMENT 'Default Billing Address',
    `default_shipping`    INT(10) UNSIGNED              DEFAULT NULL COMMENT 'Default Shipping Address',
    `gender`              SMALLINT(5) UNSIGNED          DEFAULT NULL COMMENT 'Gender',
    PRIMARY KEY (`entity_id`),
    UNIQUE KEY `CUSTOMER_ENTITY_EMAIL_WEBSITE_ID` (`email`, `website_id`),
    KEY `CUSTOMER_ENTITY_STORE_ID` (`store_id`),
    KEY `CUSTOMER_ENTITY_WEBSITE_ID` (`website_id`),
    KEY `CUSTOMER_ENTITY_FIRSTNAME` (`firstname`),
    KEY `CUSTOMER_ENTITY_LASTNAME` (`lastname`),
    CONSTRAINT `CUSTOMER_ENTITY_STORE_ID_STORE_STORE_ID` FOREIGN KEY (`store_id`) REFERENCES `store` (`store_id`) ON DELETE SET NULL,
    CONSTRAINT `CUSTOMER_ENTITY_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID` FOREIGN KEY (`website_id`) REFERENCES `store_website` (`website_id`) ON DELETE SET NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='Customer Entity';

-- Create syntax for TABLE 'customer_address_entity'
CREATE TABLE `customer_address_entity`
(
    `entity_id`    INT(10) UNSIGNED     NOT NULL AUTO_INCREMENT COMMENT 'Entity ID',
    `increment_id` VARCHAR(50)                   DEFAULT NULL COMMENT 'Increment Id',
    `parent_id`    INT(10) UNSIGNED              DEFAULT NULL COMMENT 'Parent ID',
    `created_at`   TIMESTAMP            NOT NULL DEFAULT current_timestamp() COMMENT 'Created At',
    `updated_at`   TIMESTAMP            NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Updated At',
    `is_active`    SMALLINT(5) UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Is Active',
    `city`         VARCHAR(255)         NOT NULL COMMENT 'City',
    `company`      VARCHAR(255)                  DEFAULT NULL COMMENT 'Company',
    `country_id`   VARCHAR(255)         NOT NULL COMMENT 'Country',
    `firstname`    VARCHAR(255)         NOT NULL COMMENT 'First Name',
    `lastname`     VARCHAR(255)         NOT NULL COMMENT 'Last Name',
    `postcode`     VARCHAR(255)                  DEFAULT NULL COMMENT 'Zip/Postal Code',
    `region`       VARCHAR(255)                  DEFAULT NULL COMMENT 'State/Province',
    `street`       TEXT                 NOT NULL COMMENT 'Street Address',
    PRIMARY KEY (`entity_id`),
    KEY `CUSTOMER_ADDRESS_ENTITY_PARENT_ID` (`parent_id`),
    CONSTRAINT `CUSTOMER_ADDRESS_ENTITY_PARENT_ID_CUSTOMER_ENTITY_ENTITY_ID` FOREIGN KEY (`parent_id`) REFERENCES `customer_entity` (`entity_id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='Customer Address Entity';

CREATE TABLE `customer_entity_varchar`
(
    `value_id`     INT(11)              NOT NULL AUTO_INCREMENT COMMENT 'Value ID',
    `attribute_id` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Attribute ID',
    `entity_id`    INT(10) UNSIGNED     NOT NULL DEFAULT 0 COMMENT 'Entity ID',
    `value`        VARCHAR(255)                  DEFAULT NULL COMMENT 'Value',
    PRIMARY KEY (`value_id`),
    UNIQUE KEY `CUSTOMER_ENTITY_VARCHAR_ENTITY_ID_ATTRIBUTE_ID` (`entity_id`, `attribute_id`),
    CONSTRAINT `CUSTOMER_ENTITY_VARCHAR_ENTITY_ID_CUSTOMER_ENTITY_ENTITY_ID` FOREIGN KEY (`entity_id`) REFERENCES `customer_entity` (`entity_id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='Customer Entity Varchar';

CREATE TABLE `customer_entity_int`
(
    `value_id`     INT(11)              NOT NULL AUTO_INCREMENT COMMENT 'Value ID',
    `attribute_id` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Attribute ID',
    `entity_id`    INT(10) UNSIGNED     NOT NULL DEFAULT 0 COMMENT 'Entity ID',
    `value`        INT(11)              NOT NULL DEFAULT 0 COMMENT 'Value',
    PRIMARY KEY (`value_id`),
    UNIQUE KEY `CUSTOMER_ENTITY_INT_ENTITY_ID_ATTRIBUTE_ID` (`entity_id`, `attribute_id`),
    CONSTRAINT `CUSTOMER_ENTITY_INT_ENTITY_ID_CUSTOMER_ENTITY_ENTITY_ID` FOREIGN KEY (`entity_id`) REFERENCES `customer_entity` (`entity_id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='Customer Entity Int';

CREATE VIEW `view_customer_no_auto_increment` AS
    -- this join is logically wrong because demo data needed as parent_id column in
    -- customer_address_entity is empty.
SELECT `ce`.`email`,
       `cae`.`firstname`,
       `cae`.`lastname`,
       `cae`.`city`
FROM `customer_entity` `ce`
         JOIN `customer_address_entity` `cae` ON `ce`.`entity_id` = `cae`.`entity_id`
;

CREATE VIEW `view_customer_auto_increment` AS
    -- this join is logically wrong because demo data needed as parent_id column in
    -- customer_address_entity is empty.
SELECT `ce`.`entity_id` AS ce_entity_id,
       `ce`.`email`,
       `cae`.`firstname`,
       `cae`.`lastname`,
       `cae`.`city`
FROM `customer_entity` `ce`
         JOIN `customer_address_entity` `cae` ON `ce`.`entity_id` = `cae`.`entity_id`
;

CREATE TABLE `core_configuration`
(
    `config_id`  INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Id',
    `scope`      VARCHAR(8)       NOT NULL DEFAULT 'default' COMMENT 'Scope',
    `scope_id`   INT(11)          NOT NULL DEFAULT 0 COMMENT 'Scope Id',
    `expires`    DATETIME         NULL COMMENT 'Value expiration time',
    `path`       VARCHAR(255)     NOT NULL COMMENT 'Path',
    `value`      TEXT                      DEFAULT NULL COMMENT 'Value',
    `version_ts` TIMESTAMP(6) GENERATED ALWAYS AS ROW START COMMENT 'Timestamp Start Versioning',
    `version_te` TIMESTAMP(6) GENERATED ALWAYS AS ROW END COMMENT 'Timestamp End Versioning',
    PERIOD FOR SYSTEM_TIME(`version_ts`, `version_te`),
    PRIMARY KEY (`config_id`),
    UNIQUE KEY `CORE_CONFIG_DATA_SCOPE_SCOPE_ID_PATH` (`scope`, `scope_id`, `path`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='Config Data'
  WITH SYSTEM VERSIONING
    PARTITION BY SYSTEM_TIME (
        PARTITION p_hist HISTORY,
        PARTITION p_cur CURRENT
        );

INSERT INTO `store` (`store_id`, `code`, `website_id`, `group_id`, `name`, `sort_order`, `is_active`)
VALUES (1, 'admin', 0, 0, 'Admin', 0, 1);

INSERT INTO `store_group` (`group_id`, `website_id`, `name`, `root_category_id`, `default_store_id`, `code`)
VALUES (1, 0, 'Admin', 0, 0, 'admin');

INSERT INTO `store_website` (`website_id`, `code`, `name`, `sort_order`, `default_group_id`, `is_default`)
VALUES (1, 'admin', 'Admin', 0, 0, 1);

CREATE TABLE `catalog_product_index_eav_decimal_idx`
(
    `entity_id`    INT(10) UNSIGNED     NOT NULL COMMENT 'Entity ID',
    `attribute_id` SMALLINT(5) UNSIGNED NOT NULL COMMENT 'Attribute ID',
    `store_id`     SMALLINT(5) UNSIGNED NOT NULL COMMENT 'Store ID',
    `source_id`    INT(10) UNSIGNED     NOT NULL DEFAULT 0 COMMENT 'Original entity Id for attribute value',
    `value`        DECIMAL(12, 4)       NOT NULL COMMENT 'Value',
    PRIMARY KEY (`entity_id`, `attribute_id`, `store_id`, `source_id`),
    KEY `CATALOG_PRODUCT_INDEX_EAV_DECIMAL_IDX_ATTRIBUTE_ID` (`attribute_id`),
    KEY `CATALOG_PRODUCT_INDEX_EAV_DECIMAL_IDX_STORE_ID` (`store_id`),
    KEY `CATALOG_PRODUCT_INDEX_EAV_DECIMAL_IDX_VALUE` (`value`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='Catalog Product EAV Decimal Indexer Index Table';

INSERT INTO `catalog_product_index_eav_decimal_idx` (`entity_id`, `attribute_id`, `store_id`, `source_id`, `value`)
VALUES (1, 1, 1, 1, 1.0000),
       (2, 1, 1, 2, 2.2150);

CREATE TABLE `sales_order_status_state`
(
    `status`           VARCHAR(32)          NOT NULL COMMENT 'Status',
    `state`            VARCHAR(32)          NOT NULL COMMENT 'Label',
    `is_default`       SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Is Default',
    `visible_on_front` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Visible on front',
    PRIMARY KEY (`status`, `state`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='Sales Order Status Table';

INSERT INTO `sales_order_status_state` (`status`, `state`, `is_default`, `visible_on_front`)
VALUES ('canceled', 'canceled', 1, 1),
       ('closed', 'closed', 1, 1),
       ('complete', 'complete', 1, 1),
       ('fraud', 'payment_review', 0, 1),
       ('fraud', 'processing', 0, 1),
       ('holded', 'holded', 1, 1),
       ('payment_review', 'payment_review', 1, 1),
       ('pending', 'new', 1, 1),
       ('pending_payment', 'pending_payment', 1, 0),
       ('processing', 'processing', 1, 1);

DROP TABLE IF EXISTS `sequence_catalog_category`;
CREATE TABLE `sequence_catalog_category`
(
    `sequence_value` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (`sequence_value`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 42
  DEFAULT CHARSET = utf8;

-- catalog_category_entity has a 1:1 to table sequence_catalog_category and cannot be reversed
DROP TABLE IF EXISTS `catalog_category_entity`;
CREATE TABLE `catalog_category_entity`
(
    `entity_id` INT(10) UNSIGNED NOT NULL COMMENT 'Entity Id',
    `row_id`    INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Version Id',
    PRIMARY KEY (`row_id`),
    KEY `CATALOG_CATEGORY_ENTITY_ENTITY_ID` (`entity_id`),
    CONSTRAINT `CAT_CTGR_ENTT_ENTT_ID_SEQUENCE_CAT_CTGR_SEQUENCE_VAL` FOREIGN KEY (`entity_id`) REFERENCES `sequence_catalog_category` (`sequence_value`) ON DELETE CASCADE
) ENGINE = InnoDB
  AUTO_INCREMENT = 42
  DEFAULT CHARSET = utf8 COMMENT ='Catalog Category Table';

DROP TABLE IF EXISTS `athlete_team_member`;
CREATE TABLE `athlete_team_member`
(
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
CREATE TABLE `athlete_team`
(
    `team_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name`    VARCHAR(340)     NOT NULL COMMENT 'Team name',
    PRIMARY KEY (`team_id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4` COMMENT ='Athlete Team';

DROP TABLE IF EXISTS `athlete`;
CREATE TABLE `athlete`
(
    `athlete_id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Athlete ID',
    `firstname`  VARCHAR(340) COMMENT 'First Name',
    `lastname`   VARCHAR(340) COMMENT 'Last Name',
    PRIMARY KEY (`athlete_id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = `utf8mb4` COMMENT ='Athletes';


SET foreign_key_checks = 1;
