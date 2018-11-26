SET foreign_key_checks = 0;

DROP TABLE IF EXISTS `core_config_data_generated`;

CREATE TABLE `core_config_data_generated` (
  `config_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `type_id` bigint UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Type ID',
  `scope` varchar(8) AS (CASE type_id>>24 WHEN 2 THEN 'website' WHEN 3 THEN 'group'  WHEN 4 THEN 'store' ELSE 'default' END) VIRTUAL COMMENT 'Scope',
  `scope_id` bigint UNSIGNED AS ( type_id ^ ((type_id>>24)<<24) ) PERSISTENT COMMENT 'Scope ID',
  `expires` DATETIME NULL COMMENT 'Value expiration time',
  `path` varchar(255) NOT NULL COMMENT 'Path',
  `value` text DEFAULT NULL COMMENT 'Value',
  `version_ts` TIMESTAMP(6) GENERATED ALWAYS AS ROW START COMMENT 'Timestamp Start Versioning',
  `version_te` TIMESTAMP(6) GENERATED ALWAYS AS ROW END COMMENT 'Timestamp End Versioning',
  PERIOD FOR SYSTEM_TIME(`version_ts`, `version_te`),
  PRIMARY KEY (`config_id`),
  UNIQUE KEY `CORE_CONFIG_DATA_TYPE_ID_EXPIRES_PATH` (`type_id`,`expires`,`path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Config Data'
  WITH SYSTEM VERSIONING
    PARTITION BY SYSTEM_TIME (
    PARTITION p_hist HISTORY,
    PARTITION p_cur CURRENT
    );

INSERT INTO core_config_data_generated (type_id,path,value) VALUES (16777216,'aa/bb/cc','xvalue1');
INSERT INTO core_config_data_generated (type_id,path,value) VALUES (67108865,'aa/bb/cc','xvalue2');
INSERT INTO core_config_data_generated (type_id,path,value) VALUES (50331652,'aa/bb/cc','xvalue3');

SET foreign_key_checks = 1;
