SET foreign_key_checks = 0;

DROP TABLE IF EXISTS `dmlgen_types`;
DROP TABLE IF EXISTS `store`;
DROP TABLE IF EXISTS `store_group`;
DROP TABLE IF EXISTS `store_website`;
DROP TABLE IF EXISTS `customer_entity`;
DROP TABLE IF EXISTS `customer_address_entity`;
DROP TABLE IF EXISTS `core_config_data`;

CREATE TABLE `dmlgen_types` (
  `id`                      INT(11)              NOT NULL           AUTO_INCREMENT,
  col_bigint_1              BIGINT(20)           NULL,
  col_bigint_2              BIGINT(20)           NOT NULL           DEFAULT 0,
  col_bigint_3              BIGINT(20) UNSIGNED  NULL,
  col_bigint_4              BIGINT(20) UNSIGNED  NOT NULL           DEFAULT 0,
  col_blob                  BLOB                                    DEFAULT NULL,
  col_date_1                DATE                                    DEFAULT NULL,
  col_date_2                DATE                 NOT NULL           DEFAULT '0000-00-00',
  col_datetime_1            DATETIME                                DEFAULT NULL,
  col_datetime_2            DATETIME             NOT NULL           DEFAULT '0000-00-00 00:00:00',
  col_decimal_10_0          DECIMAL(10, 0) UNSIGNED                 DEFAULT NULL,
  col_decimal_12_4          DECIMAL(12, 4)                          DEFAULT NULL,
  price_12_4a               DECIMAL(12, 4)                          DEFAULT NULL,
  price_12_4b               DECIMAL(12, 4)       NOT NULL           DEFAULT 0,
  col_decimal_12_3          DECIMAL(12, 3)       NOT NULL           DEFAULT 0,
  col_decimal_20_6          DECIMAL(20, 6)       NOT NULL           DEFAULT 0.000000,
  col_decimal_24_12         DECIMAL(24, 12)      NOT NULL           DEFAULT 0.000000000000,
  col_float                 FLOAT                NOT NULL           DEFAULT 1,
  col_int_1                 INT(10)              NULL,
  col_int_2                 INT(10)              NOT NULL           DEFAULT 0,
  col_int_3                 INT(10) UNSIGNED     NULL,
  col_int_4                 INT(10) UNSIGNED     NOT NULL           DEFAULT 0,
  col_longtext_1            LONGTEXT                                DEFAULT NULL,
  col_longtext_2            LONGTEXT             NOT NULL           DEFAULT '',
  col_mediumblob            MEDIUMBLOB                              DEFAULT NULL,
  col_mediumtext_1          MEDIUMTEXT                              DEFAULT NULL,
  col_mediumtext_2          MEDIUMTEXT           NOT NULL           DEFAULT '',
  col_smallint_1            SMALLINT(5)          NULL,
  col_smallint_2            SMALLINT(5)          NOT NULL           DEFAULT 0,
  col_smallint_3            SMALLINT(5) UNSIGNED NULL,
  col_smallint_4            SMALLINT(5) UNSIGNED NOT NULL           DEFAULT 0,
  has_smallint_5            SMALLINT(5) UNSIGNED NOT NULL           DEFAULT 0,
  is_smallint_5             SMALLINT(5)          NULL,
  col_text                  TEXT                                    DEFAULT NULL,
  col_timestamp_1           TIMESTAMP            NOT NULL           DEFAULT current_timestamp(),
  col_timestamp_2           TIMESTAMP            NULL,
  col_tinyint_1             TINYINT(1)           NOT NULL           DEFAULT 0,
  col_varchar_1             VARCHAR(1)           NOT NULL           DEFAULT '0',
  col_varchar_100           VARCHAR(100)                            DEFAULT NULL,
  col_varchar_16            VARCHAR(16)          NOT NULL           DEFAULT 'de_DE',
  col_char_1                char(21)                                DEFAULT NULL,
  col_char_2                char(17)             NOT NULL DEFAULT 'xchar',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create syntax for TABLE 'store'
CREATE TABLE `store` (
  `store_id` smallint(5) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Store Id',
  `code` varchar(32) DEFAULT NULL COMMENT 'Code',
  `website_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Website Id',
  `group_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Group Id',
  `name` varchar(255) NOT NULL COMMENT 'Store Name',
  `sort_order` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Store Sort Order',
  `is_active` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Store Activity',
  PRIMARY KEY (`store_id`),
  UNIQUE KEY `STORE_CODE` (`code`),
  KEY `STORE_WEBSITE_ID` (`website_id`),
  KEY `STORE_IS_ACTIVE_SORT_ORDER` (`is_active`,`sort_order`),
  KEY `STORE_GROUP_ID` (`group_id`),
  CONSTRAINT `STORE_GROUP_ID_STORE_GROUP_GROUP_ID` FOREIGN KEY (`group_id`) REFERENCES `store_group` (`group_id`) ON DELETE CASCADE,
  CONSTRAINT `STORE_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID` FOREIGN KEY (`website_id`) REFERENCES `store_website` (`website_id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='Stores';

-- Create syntax for TABLE 'store_group'
CREATE TABLE `store_group` (
  `group_id` smallint(5) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Group Id',
  `website_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Website Id',
  `name` varchar(255) NOT NULL COMMENT 'Store Group Name',
  `root_category_id` int(10) unsigned NOT NULL DEFAULT 0 COMMENT 'Root Category Id',
  `default_store_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Default Store Id',
  `code` varchar(32) DEFAULT NULL COMMENT 'Store group unique code',
  PRIMARY KEY (`group_id`),
  UNIQUE KEY `STORE_GROUP_CODE` (`code`),
  KEY `STORE_GROUP_WEBSITE_ID` (`website_id`),
  KEY `STORE_GROUP_DEFAULT_STORE_ID` (`default_store_id`),
  CONSTRAINT `STORE_GROUP_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID` FOREIGN KEY (`website_id`) REFERENCES `store_website` (`website_id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='Store Groups';

-- Create syntax for TABLE 'store_website'
CREATE TABLE `store_website` (
  `website_id` smallint(5) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Website Id',
  `code` varchar(32) DEFAULT NULL COMMENT 'Code',
  `name` varchar(64) DEFAULT NULL COMMENT 'Website Name',
  `sort_order` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Sort Order',
  `default_group_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Default Group Id',
  `is_default` smallint(5) unsigned DEFAULT 0 COMMENT 'Defines Is Website Default',
  PRIMARY KEY (`website_id`),
  UNIQUE KEY `STORE_WEBSITE_CODE` (`code`),
  KEY `STORE_WEBSITE_SORT_ORDER` (`sort_order`),
  KEY `STORE_WEBSITE_DEFAULT_GROUP_ID` (`default_group_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='Websites';

-- Create syntax for TABLE 'customer_entity'
CREATE TABLE `customer_entity` (
  `entity_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Entity ID',
  `website_id` smallint(5) unsigned DEFAULT NULL COMMENT 'Website ID',
  `email` varchar(255) DEFAULT NULL COMMENT 'Email',
  `group_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Group ID',
  `increment_id` varchar(50) DEFAULT NULL COMMENT 'Increment Id',
  `store_id` smallint(5) unsigned DEFAULT 0 COMMENT 'Store ID',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Created At',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Updated At',
  `is_active` smallint(5) unsigned NOT NULL DEFAULT 1 COMMENT 'Is Active',
  `disable_auto_group_change` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Disable automatic group change based on VAT ID',
  `created_in` varchar(255) DEFAULT NULL COMMENT 'Created From',
  `prefix` varchar(40) DEFAULT NULL COMMENT 'Name Prefix',
  `firstname` varchar(255) DEFAULT NULL COMMENT 'First Name',
  `middlename` varchar(255) DEFAULT NULL COMMENT 'Middle Name/Initial',
  `lastname` varchar(255) DEFAULT NULL COMMENT 'Last Name',
  `suffix` varchar(40) DEFAULT NULL COMMENT 'Name Suffix',
  `dob` date DEFAULT NULL COMMENT 'Date of Birth',
  `password_hash` varchar(128) DEFAULT NULL COMMENT 'Password_hash',
  `rp_token` varchar(128) DEFAULT NULL COMMENT 'Reset password token',
  `rp_token_created_at` datetime DEFAULT NULL COMMENT 'Reset password token creation time',
  `default_billing` int(10) unsigned DEFAULT NULL COMMENT 'Default Billing Address',
  `default_shipping` int(10) unsigned DEFAULT NULL COMMENT 'Default Shipping Address',
  `taxvat` varchar(50) DEFAULT NULL COMMENT 'Tax/VAT Number',
  `confirmation` varchar(64) DEFAULT NULL COMMENT 'Is Confirmed',
  `gender` smallint(5) unsigned DEFAULT NULL COMMENT 'Gender',
  `failures_num` smallint(6) DEFAULT 0 COMMENT 'Failure Number',
  `first_failure` timestamp NULL DEFAULT NULL COMMENT 'First Failure',
  `lock_expires` timestamp NULL DEFAULT NULL COMMENT 'Lock Expiration Date',
  PRIMARY KEY (`entity_id`),
  UNIQUE KEY `CUSTOMER_ENTITY_EMAIL_WEBSITE_ID` (`email`,`website_id`),
  KEY `CUSTOMER_ENTITY_STORE_ID` (`store_id`),
  KEY `CUSTOMER_ENTITY_WEBSITE_ID` (`website_id`),
  KEY `CUSTOMER_ENTITY_FIRSTNAME` (`firstname`),
  KEY `CUSTOMER_ENTITY_LASTNAME` (`lastname`),
  CONSTRAINT `CUSTOMER_ENTITY_STORE_ID_STORE_STORE_ID` FOREIGN KEY (`store_id`) REFERENCES `store` (`store_id`) ON DELETE SET NULL,
  CONSTRAINT `CUSTOMER_ENTITY_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID` FOREIGN KEY (`website_id`) REFERENCES `store_website` (`website_id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='Customer Entity';

-- Create syntax for TABLE 'customer_address_entity'
CREATE TABLE `customer_address_entity` (
  `entity_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Entity ID',
  `increment_id` varchar(50) DEFAULT NULL COMMENT 'Increment Id',
  `parent_id` int(10) unsigned DEFAULT NULL COMMENT 'Parent ID',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Created At',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Updated At',
  `is_active` smallint(5) unsigned NOT NULL DEFAULT 1 COMMENT 'Is Active',
  `city` varchar(255) NOT NULL COMMENT 'City',
  `company` varchar(255) DEFAULT NULL COMMENT 'Company',
  `country_id` varchar(255) NOT NULL COMMENT 'Country',
  `fax` varchar(255) DEFAULT NULL COMMENT 'Fax',
  `firstname` varchar(255) NOT NULL COMMENT 'First Name',
  `lastname` varchar(255) NOT NULL COMMENT 'Last Name',
  `middlename` varchar(255) DEFAULT NULL COMMENT 'Middle Name',
  `postcode` varchar(255) DEFAULT NULL COMMENT 'Zip/Postal Code',
  `prefix` varchar(40) DEFAULT NULL COMMENT 'Name Prefix',
  `region` varchar(255) DEFAULT NULL COMMENT 'State/Province',
  `region_id` int(10) unsigned DEFAULT NULL COMMENT 'State/Province',
  `street` text NOT NULL COMMENT 'Street Address',
  `suffix` varchar(40) DEFAULT NULL COMMENT 'Name Suffix',
  `telephone` varchar(255) NOT NULL COMMENT 'Phone Number',
  `vat_id` varchar(255) DEFAULT NULL COMMENT 'VAT number',
  `vat_is_valid` int(10) unsigned DEFAULT NULL COMMENT 'VAT number validity',
  `vat_request_date` varchar(255) DEFAULT NULL COMMENT 'VAT number validation request date',
  `vat_request_id` varchar(255) DEFAULT NULL COMMENT 'VAT number validation request ID',
  `vat_request_success` int(10) unsigned DEFAULT NULL COMMENT 'VAT number validation request success',
  PRIMARY KEY (`entity_id`),
  KEY `CUSTOMER_ADDRESS_ENTITY_PARENT_ID` (`parent_id`),
  CONSTRAINT `CUSTOMER_ADDRESS_ENTITY_PARENT_ID_CUSTOMER_ENTITY_ENTITY_ID` FOREIGN KEY (`parent_id`) REFERENCES `customer_entity` (`entity_id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='Customer Address Entity';

CREATE TABLE `core_config_data` (
  `config_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Id',
  `scope` varchar(8) NOT NULL DEFAULT 'default' COMMENT 'Scope',
  `scope_id` int(11) NOT NULL DEFAULT 0 COMMENT 'Scope Id',
  `expires` DATETIME NULL COMMENT 'Value expiration time',
  `path` varchar(255) NOT NULL COMMENT 'Path',
  `value` text DEFAULT NULL COMMENT 'Value',
  `version_ts` TIMESTAMP(6) GENERATED ALWAYS AS ROW START COMMENT 'Timestamp Start Versioning',
  `version_te` TIMESTAMP(6) GENERATED ALWAYS AS ROW END COMMENT 'Timestamp End Versioning',
  PERIOD FOR SYSTEM_TIME(`version_ts`, `version_te`),
  PRIMARY KEY (`config_id`),
  UNIQUE KEY `CORE_CONFIG_DATA_SCOPE_SCOPE_ID_PATH` (`scope`,`scope_id`,`expires`,`path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Config Data'
  WITH SYSTEM VERSIONING
    PARTITION BY SYSTEM_TIME (
    PARTITION p_hist HISTORY,
    PARTITION p_cur CURRENT
    );


SET foreign_key_checks = 1;
