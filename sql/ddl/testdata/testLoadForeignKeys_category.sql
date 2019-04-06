SET foreign_key_checks = 0;
SET NAMES utf8mb4;

DROP TABLE IF EXISTS `sequence_catalog_category`;
CREATE TABLE `sequence_catalog_category` (
    `sequence_value` int(10) unsigned NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (`sequence_value`)
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8;

-- catalog_category_entity has a 1:1 to table sequence_catalog_category and cannot be reversed
DROP TABLE IF EXISTS `catalog_category_entity`;
CREATE TABLE `catalog_category_entity` (
    `entity_id` int(10) unsigned NOT NULL COMMENT 'Entity Id',
    `row_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Version Id',
    PRIMARY KEY (`row_id`),
    KEY `CATALOG_CATEGORY_ENTITY_ENTITY_ID` (`entity_id`),
    CONSTRAINT `CAT_CTGR_ENTT_ENTT_ID_SEQUENCE_CAT_CTGR_SEQUENCE_VAL` FOREIGN KEY (`entity_id`) REFERENCES `sequence_catalog_category` (`sequence_value`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8 COMMENT='Catalog Category Table';

SET foreign_key_checks = 1;
