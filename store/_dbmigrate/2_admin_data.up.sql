SET FOREIGN_KEY_CHECKS=0;

SET NAMES 'utf8mb4';

LOCK TABLES `store` WRITE;
/*!40000 ALTER TABLE `store`
  DISABLE KEYS */;

INSERT INTO `store` (`store_id`, `code`, `website_id`, `group_id`, `name`, `sort_order`, `is_active`)
VALUES (1, 'admin', 0, 0, 'Admin', 0, 1);

/*!40000 ALTER TABLE `store`
  ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table store_group
# ------------------------------------------------------------

LOCK TABLES `store_group` WRITE;
/*!40000 ALTER TABLE `store_group`
  DISABLE KEYS */;

INSERT INTO `store_group` (`group_id`, `website_id`, `name`, `root_category_id`, `default_store_id`, `code`)
VALUES (1, 0, 'Admin', 0, 0, 'admin');

/*!40000 ALTER TABLE `store_group`
  ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table store_website
# ------------------------------------------------------------

LOCK TABLES `store_website` WRITE;
/*!40000 ALTER TABLE `store_website`
  DISABLE KEYS */;

INSERT INTO `store_website` (`website_id`, `code`, `name`, `sort_order`, `default_group_id`, `is_default`)
VALUES (1, 'admin', 'Admin', 0, 0, 1);

/*!40000 ALTER TABLE `store_website`
  ENABLE KEYS */;
UNLOCK TABLES;

SET FOREIGN_KEY_CHECKS=1;
