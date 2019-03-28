CREATE TABLE `core_configuration` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `scope` varchar(8) NOT NULL DEFAULT 'default' COMMENT 'Scope',
  `scope_id` int(11) NOT NULL DEFAULT 0 COMMENT 'Scope ID',
  `expires` DATETIME NULL COMMENT 'Value expiration time',
  `path` varchar(255) NOT NULL COMMENT 'Path',
  `value` text DEFAULT NULL COMMENT 'Value',
  `version_ts` TIMESTAMP(6) GENERATED ALWAYS AS ROW START INVISIBLE COMMENT 'Timestamp Start Versioning',
  `version_te` TIMESTAMP(6) GENERATED ALWAYS AS ROW END INVISIBLE COMMENT 'Timestamp End Versioning',
  PERIOD FOR SYSTEM_TIME(`version_ts`, `version_te`),
  PRIMARY KEY (`id`),
  UNIQUE KEY `CORE_CONFIG_DATA_SCOPE_SCOPE_ID_PATH` (`scope`,`scope_id`,`expires`,`path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Config Data'
  WITH SYSTEM VERSIONING
  PARTITION BY SYSTEM_TIME (
    PARTITION p_hist HISTORY,
    PARTITION p_cur CURRENT
  );
