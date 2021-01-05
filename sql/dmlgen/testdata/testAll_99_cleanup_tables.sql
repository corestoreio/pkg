SET foreign_key_checks = 0;

DROP TABLE IF EXISTS `dmlgen_types`;
DROP TABLE IF EXISTS `store`;
DROP TABLE IF EXISTS `store_group`;
DROP TABLE IF EXISTS `store_website`;
DROP TABLE IF EXISTS `customer_entity`;
DROP TABLE IF EXISTS `customer_address_entity`;
DROP TABLE IF EXISTS `customer_entity_varchar`;
DROP TABLE IF EXISTS `customer_entity_int`;
DROP TABLE IF EXISTS `core_configuration`;

DROP VIEW IF EXISTS `view_customer_no_auto_increment`;
DROP VIEW IF EXISTS `view_customer_auto_increment`;
DROP TABLE IF EXISTS `catalog_product_index_eav_decimal_idx`;
DROP TABLE IF EXISTS `sales_order_status_state`;

DROP TABLE IF EXISTS `catalog_category_entity`;
DROP TABLE IF EXISTS `sequence_catalog_category`;

DROP TABLE IF EXISTS `athlete_team_member`;
DROP TABLE IF EXISTS `athlete_team`;
DROP TABLE IF EXISTS `athlete`;

SET foreign_key_checks = 1;
