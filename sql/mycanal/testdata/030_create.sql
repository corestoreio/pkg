/*
 * Copyright © Magento, Inc. All rights reserved.
 * Licensed under OSL 3.0
 * http://opensource.org/licenses/osl-3.0.php Open Software License (OSL 3.0)
*/
SET NAMES utf8mb4;
SET foreign_key_checks = 0;
DROP TABLE IF EXISTS `catalog_category_entity`;

CREATE TABLE `catalog_category_entity` (
  `entity_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Entity Id',
  `attribute_set_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Attribute Set ID',
  `parent_id` int(10) unsigned NOT NULL DEFAULT 0 COMMENT 'Parent Category ID',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() COMMENT 'Creation Time',
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT 'Update Time',
  `path` varchar(255) NOT NULL COMMENT 'Tree Path',
  `position` int(11) NOT NULL COMMENT 'Position',
  `level` int(11) NOT NULL DEFAULT 0 COMMENT 'Tree Level',
  `children_count` int(11) NOT NULL COMMENT 'Child Count',
  PRIMARY KEY (`entity_id`),
  KEY `CATALOG_CATEGORY_ENTITY_LEVEL` (`level`),
  KEY `CATALOG_CATEGORY_ENTITY_PATH` (`path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Catalog Category Table';

LOCK TABLES `catalog_category_entity` WRITE;
/*!40000 ALTER TABLE `catalog_category_entity` DISABLE KEYS */;

INSERT INTO `catalog_category_entity` (`entity_id`, `attribute_set_id`, `parent_id`, `created_at`, `updated_at`, `path`, `position`, `level`, `children_count`)
VALUES
	(1,3,0,'2018-04-17 21:33:50','2018-04-17 21:44:22','1',0,0,39),
	(2,3,1,'2018-04-17 21:33:50','2018-04-17 21:44:22','1/2',1,1,38),
	(3,3,2,'2018-04-17 21:41:51','2018-04-17 21:41:51','1/2/3',4,2,3),
	(4,3,3,'2018-04-17 21:41:51','2018-04-17 21:41:51','1/2/3/4',1,3,0),
	(5,3,3,'2018-04-17 21:41:51','2018-04-17 21:41:51','1/2/3/5',2,3,0),
	(6,3,3,'2018-04-17 21:41:51','2018-04-18 07:40:23','1/2/3/6',3,3,0),
	(7,3,2,'2018-04-17 21:41:51','2018-04-17 21:44:22','1/2/7',5,2,6),
	(8,3,7,'2018-04-17 21:41:51','2018-04-17 21:41:51','1/2/7/8',1,3,0),
	(9,3,2,'2018-04-17 21:42:02','2018-04-17 21:42:02','1/2/9',5,2,1),
	(10,3,9,'2018-04-17 21:42:02','2018-04-17 21:42:03','1/2/9/10',1,3,0),
	(11,3,2,'2018-04-17 21:42:07','2018-04-17 21:42:08','1/2/11',3,2,8),
	(12,3,11,'2018-04-17 21:42:07','2018-04-17 21:42:08','1/2/11/12',1,3,4),
	(13,3,11,'2018-04-17 21:42:07','2018-04-17 21:42:08','1/2/11/13',2,3,2),
	(14,3,12,'2018-04-17 21:42:07','2018-04-17 21:42:07','1/2/11/12/14',1,4,0),
	(15,3,12,'2018-04-17 21:42:08','2018-04-17 21:42:08','1/2/11/12/15',2,4,0),
	(16,3,12,'2018-04-17 21:42:08','2018-04-17 21:42:08','1/2/11/12/16',3,4,0),
	(17,3,12,'2018-04-17 21:42:08','2018-04-17 21:42:08','1/2/11/12/17',4,4,0),
	(18,3,13,'2018-04-17 21:42:08','2018-04-17 21:42:08','1/2/11/13/18',1,4,0),
	(19,3,13,'2018-04-17 21:42:08','2018-04-17 21:42:08','1/2/11/13/19',2,4,0),
	(20,3,2,'2018-04-17 21:42:08','2018-04-17 21:42:09','1/2/20',2,2,8),
	(21,3,20,'2018-04-17 21:42:09','2018-04-17 21:42:09','1/2/20/21',1,3,4),
	(22,3,20,'2018-04-17 21:42:09','2018-04-17 21:42:09','1/2/20/22',2,3,2),
	(23,3,21,'2018-04-17 21:42:09','2018-04-17 21:42:09','1/2/20/21/23',1,4,0),
	(24,3,21,'2018-04-17 21:42:09','2018-04-17 21:42:09','1/2/20/21/24',2,4,0),
	(25,3,21,'2018-04-17 21:42:09','2018-04-17 21:42:09','1/2/20/21/25',3,4,0),
	(26,3,21,'2018-04-17 21:42:09','2018-04-17 21:42:09','1/2/20/21/26',4,4,0),
	(27,3,22,'2018-04-17 21:42:09','2018-04-17 21:42:09','1/2/20/22/27',1,4,0),
	(28,3,22,'2018-04-17 21:42:09','2018-04-17 21:42:10','1/2/20/22/28',2,4,0),
	(29,3,2,'2018-04-17 21:42:10','2018-04-17 21:42:10','1/2/29',6,2,4),
	(30,3,29,'2018-04-17 21:42:10','2018-04-17 21:42:10','1/2/29/30',1,3,0),
	(31,3,29,'2018-04-17 21:42:10','2018-04-17 21:42:10','1/2/29/31',2,3,0),
	(32,3,29,'2018-04-17 21:42:10','2018-04-17 21:42:10','1/2/29/32',3,3,0),
	(33,3,29,'2018-04-17 21:42:10','2018-04-17 21:42:10','1/2/29/33',4,3,0),
	(34,3,7,'2018-04-17 21:42:10','2018-04-17 21:42:10','1/2/7/34',2,3,0),
	(35,3,7,'2018-04-17 21:42:10','2018-04-17 21:42:10','1/2/7/35',3,3,0),
	(36,3,7,'2018-04-17 21:42:11','2018-04-17 21:42:11','1/2/7/36',4,3,0),
	(37,3,2,'2018-04-17 21:44:22','2018-04-17 21:44:22','1/2/37',6,2,0),
	(38,3,2,'2018-04-17 21:44:22','2018-04-17 21:44:22','1/2/38',1,2,0),
	(39,3,7,'2018-04-17 21:44:22','2018-04-17 21:44:22','1/2/7/39',5,3,0),
	(40,3,7,'2018-04-17 21:44:22','2018-04-17 21:44:22','1/2/7/40',6,3,0);

/*!40000 ALTER TABLE `catalog_category_entity` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table catalog_category_entity_datetime
# ------------------------------------------------------------

DROP TABLE IF EXISTS `catalog_category_entity_datetime`;

CREATE TABLE `catalog_category_entity_datetime` (
  `value_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Value ID',
  `attribute_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Attribute ID',
  `store_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Store ID',
  `entity_id` int(10) unsigned NOT NULL DEFAULT 0,
  `value` datetime DEFAULT NULL COMMENT 'Value',
  PRIMARY KEY (`value_id`),
  UNIQUE KEY `CATALOG_CATEGORY_ENTITY_DATETIME_ENTITY_ID_ATTRIBUTE_ID_STORE_ID` (`entity_id`,`attribute_id`,`store_id`),
  KEY `CATALOG_CATEGORY_ENTITY_DATETIME_ENTITY_ID` (`entity_id`),
  KEY `CATALOG_CATEGORY_ENTITY_DATETIME_ATTRIBUTE_ID` (`attribute_id`),
  KEY `CATALOG_CATEGORY_ENTITY_DATETIME_STORE_ID` (`store_id`),
  CONSTRAINT `CAT_CTGR_ENTT_DTIME_ENTT_ID_CAT_CTGR_ENTT_ENTT_ID` FOREIGN KEY (`entity_id`) REFERENCES `catalog_category_entity` (`entity_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Catalog Category Datetime Attribute Backend Table';

LOCK TABLES `catalog_category_entity_datetime` WRITE;
/*!40000 ALTER TABLE `catalog_category_entity_datetime` DISABLE KEYS */;

INSERT INTO `catalog_category_entity_datetime` (`value_id`, `attribute_id`, `store_id`, `entity_id`, `value`)
VALUES
	(1,61,2,6,NULL),
	(2,62,2,6,NULL);

/*!40000 ALTER TABLE `catalog_category_entity_datetime` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table catalog_category_entity_decimal
# ------------------------------------------------------------

DROP TABLE IF EXISTS `catalog_category_entity_decimal`;

CREATE TABLE `catalog_category_entity_decimal` (
  `value_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Value ID',
  `attribute_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Attribute ID',
  `store_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Store ID',
  `entity_id` int(10) unsigned NOT NULL DEFAULT 0,
  `value` decimal(12,4) DEFAULT NULL COMMENT 'Value',
  PRIMARY KEY (`value_id`),
  UNIQUE KEY `CATALOG_CATEGORY_ENTITY_DECIMAL_ENTITY_ID_ATTRIBUTE_ID_STORE_ID` (`entity_id`,`attribute_id`,`store_id`),
  KEY `CATALOG_CATEGORY_ENTITY_DECIMAL_ENTITY_ID` (`entity_id`),
  KEY `CATALOG_CATEGORY_ENTITY_DECIMAL_ATTRIBUTE_ID` (`attribute_id`),
  KEY `CATALOG_CATEGORY_ENTITY_DECIMAL_STORE_ID` (`store_id`),
  CONSTRAINT `CAT_CTGR_ENTT_DEC_ENTT_ID_CAT_CTGR_ENTT_ENTT_ID` FOREIGN KEY (`entity_id`) REFERENCES `catalog_category_entity` (`entity_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Catalog Category Decimal Attribute Backend Table';



# Dump of table catalog_category_entity_int
# ------------------------------------------------------------

DROP TABLE IF EXISTS `catalog_category_entity_int`;

CREATE TABLE `catalog_category_entity_int` (
  `value_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Value ID',
  `attribute_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Attribute ID',
  `store_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Store ID',
  `entity_id` int(10) unsigned NOT NULL DEFAULT 0,
  `value` int(11) DEFAULT NULL COMMENT 'Value',
  PRIMARY KEY (`value_id`),
  UNIQUE KEY `CATALOG_CATEGORY_ENTITY_INT_ENTITY_ID_ATTRIBUTE_ID_STORE_ID` (`entity_id`,`attribute_id`,`store_id`),
  KEY `CATALOG_CATEGORY_ENTITY_INT_ENTITY_ID` (`entity_id`),
  KEY `CATALOG_CATEGORY_ENTITY_INT_ATTRIBUTE_ID` (`attribute_id`),
  KEY `CATALOG_CATEGORY_ENTITY_INT_STORE_ID` (`store_id`),
  CONSTRAINT `CAT_CTGR_ENTT_INT_ENTT_ID_CAT_CTGR_ENTT_ENTT_ID` FOREIGN KEY (`entity_id`) REFERENCES `catalog_category_entity` (`entity_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Catalog Category Integer Attribute Backend Table';

LOCK TABLES `catalog_category_entity_int` WRITE;
/*!40000 ALTER TABLE `catalog_category_entity_int` DISABLE KEYS */;

INSERT INTO `catalog_category_entity_int` (`value_id`, `attribute_id`, `store_id`, `entity_id`, `value`)
VALUES
	(1,69,0,1,1),
	(2,46,0,2,1),
	(3,69,0,2,1),
	(4,46,0,3,1),
	(5,54,0,3,0),
	(6,69,0,3,1),
	(7,46,0,4,1),
	(8,54,0,4,1),
	(9,69,0,4,1),
	(10,46,0,5,1),
	(11,54,0,5,1),
	(12,69,0,5,1),
	(13,46,0,6,1),
	(14,54,0,6,1),
	(15,69,0,6,1),
	(16,46,0,7,0),
	(17,54,0,7,0),
	(18,69,0,7,0),
	(19,46,0,8,1),
	(20,54,0,8,1),
	(21,69,0,8,0),
	(22,46,0,9,1),
	(23,54,0,9,0),
	(24,69,0,9,1),
	(25,46,0,10,1),
	(26,54,0,10,1),
	(27,69,0,10,1),
	(28,46,0,11,1),
	(29,54,0,11,0),
	(30,69,0,11,1),
	(31,46,0,12,1),
	(32,54,0,12,1),
	(33,69,0,12,1),
	(34,46,0,13,1),
	(35,54,0,13,1),
	(36,69,0,13,1),
	(37,46,0,14,1),
	(38,54,0,14,1),
	(39,69,0,14,1),
	(40,46,0,15,1),
	(41,54,0,15,1),
	(42,69,0,15,1),
	(43,46,0,16,1),
	(44,54,0,16,1),
	(45,69,0,16,1),
	(46,46,0,17,1),
	(47,54,0,17,1),
	(48,69,0,17,1),
	(49,46,0,18,1),
	(50,54,0,18,1),
	(51,69,0,18,1),
	(52,46,0,19,1),
	(53,54,0,19,1),
	(54,69,0,19,1),
	(55,46,0,20,1),
	(56,54,0,20,0),
	(57,69,0,20,1),
	(58,46,0,21,1),
	(59,54,0,21,1),
	(60,69,0,21,1),
	(61,46,0,22,1),
	(62,54,0,22,1),
	(63,69,0,22,1),
	(64,46,0,23,1),
	(65,54,0,23,1),
	(66,69,0,23,1),
	(67,46,0,24,1),
	(68,54,0,24,1),
	(69,69,0,24,1),
	(70,46,0,25,1),
	(71,54,0,25,1),
	(72,69,0,25,1),
	(73,46,0,26,1),
	(74,54,0,26,1),
	(75,69,0,26,1),
	(76,46,0,27,1),
	(77,54,0,27,1),
	(78,69,0,27,1),
	(79,46,0,28,1),
	(80,54,0,28,1),
	(81,69,0,28,1),
	(82,46,0,29,0),
	(83,54,0,29,0),
	(84,69,0,29,0),
	(85,46,0,30,1),
	(86,54,0,30,1),
	(87,69,0,30,0),
	(88,46,0,31,1),
	(89,54,0,31,1),
	(90,69,0,31,0),
	(91,46,0,32,1),
	(92,54,0,32,1),
	(93,69,0,32,0),
	(94,46,0,33,1),
	(95,54,0,33,1),
	(96,69,0,33,0),
	(97,46,0,34,1),
	(98,54,0,34,1),
	(99,69,0,34,0),
	(100,46,0,35,1),
	(101,54,0,35,1),
	(102,69,0,35,0),
	(103,46,0,36,1),
	(104,54,0,36,1),
	(105,69,0,36,0),
	(106,46,0,37,1),
	(107,54,0,37,0),
	(108,69,0,37,1),
	(109,46,0,38,1),
	(110,54,0,38,0),
	(111,69,0,38,1),
	(112,46,0,39,1),
	(113,54,0,39,0),
	(114,69,0,39,0),
	(115,46,0,40,1),
	(116,54,0,40,0),
	(117,69,0,40,0),
	(118,71,2,6,0);

/*!40000 ALTER TABLE `catalog_category_entity_int` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table catalog_category_entity_text
# ------------------------------------------------------------

DROP TABLE IF EXISTS `catalog_category_entity_text`;

CREATE TABLE `catalog_category_entity_text` (
  `value_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Value ID',
  `attribute_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Attribute ID',
  `store_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Store ID',
  `entity_id` int(10) unsigned NOT NULL DEFAULT 0,
  `value` text DEFAULT NULL COMMENT 'Value',
  PRIMARY KEY (`value_id`),
  UNIQUE KEY `CATALOG_CATEGORY_ENTITY_TEXT_ENTITY_ID_ATTRIBUTE_ID_STORE_ID` (`entity_id`,`attribute_id`,`store_id`),
  KEY `CATALOG_CATEGORY_ENTITY_TEXT_ENTITY_ID` (`entity_id`),
  KEY `CATALOG_CATEGORY_ENTITY_TEXT_ATTRIBUTE_ID` (`attribute_id`),
  KEY `CATALOG_CATEGORY_ENTITY_TEXT_STORE_ID` (`store_id`),
  CONSTRAINT `CAT_CTGR_ENTT_TEXT_ENTT_ID_CAT_CTGR_ENTT_ENTT_ID` FOREIGN KEY (`entity_id`) REFERENCES `catalog_category_entity` (`entity_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Catalog Category Text Attribute Backend Table';

LOCK TABLES `catalog_category_entity_text` WRITE;
/*!40000 ALTER TABLE `catalog_category_entity_text` DISABLE KEYS */;

INSERT INTO `catalog_category_entity_text` (`value_id`, `attribute_id`, `store_id`, `entity_id`, `value`)
VALUES
	(1,64,0,3,'<ref name=\"catalog.leftnav\" remove=\"true\"/>'),
	(2,64,0,9,'<ref name=\"catalog.leftnav\" remove=\"true\"/>'),
	(3,64,0,11,'<ref name=\"catalog.leftnav\" remove=\"true\"/>'),
	(4,64,0,20,'<ref name=\"catalog.leftnav\" remove=\"true\"/>'),
	(5,64,2,6,NULL);

/*!40000 ALTER TABLE `catalog_category_entity_text` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table catalog_category_entity_varchar
# ------------------------------------------------------------

DROP TABLE IF EXISTS `catalog_category_entity_varchar`;

CREATE TABLE `catalog_category_entity_varchar` (
  `value_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Value ID',
  `attribute_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Attribute ID',
  `store_id` smallint(5) unsigned NOT NULL DEFAULT 0 COMMENT 'Store ID',
  `entity_id` int(10) unsigned NOT NULL DEFAULT 0,
  `value` varchar(255) DEFAULT NULL COMMENT 'Value',
  PRIMARY KEY (`value_id`),
  UNIQUE KEY `CATALOG_CATEGORY_ENTITY_VARCHAR_ENTITY_ID_ATTRIBUTE_ID_STORE_ID` (`entity_id`,`attribute_id`,`store_id`),
  KEY `CATALOG_CATEGORY_ENTITY_VARCHAR_ENTITY_ID` (`entity_id`),
  KEY `CATALOG_CATEGORY_ENTITY_VARCHAR_ATTRIBUTE_ID` (`attribute_id`),
  KEY `CATALOG_CATEGORY_ENTITY_VARCHAR_STORE_ID` (`store_id`),
  CONSTRAINT `CAT_CTGR_ENTT_VCHR_ENTT_ID_CAT_CTGR_ENTT_ENTT_ID` FOREIGN KEY (`entity_id`) REFERENCES `catalog_category_entity` (`entity_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Catalog Category Varchar Attribute Backend Table';

LOCK TABLES `catalog_category_entity_varchar` WRITE;
/*!40000 ALTER TABLE `catalog_category_entity_varchar` DISABLE KEYS */;

INSERT INTO `catalog_category_entity_varchar` (`value_id`, `attribute_id`, `store_id`, `entity_id`, `value`)
VALUES
	(1,45,0,1,'Root Catalog'),
	(2,45,0,2,'Default Category'),
	(3,52,0,2,'PRODUCTS'),
	(4,45,0,3,'Gear'),
	(5,52,0,3,'PAGE'),
	(6,117,0,3,'gear'),
	(7,118,0,3,'gear'),
	(8,45,0,4,'Bags'),
	(9,117,0,4,'bags'),
	(10,118,0,4,'gear/bags'),
	(11,45,0,5,'Fitness Equipment'),
	(12,117,0,5,'fitness-equipment'),
	(13,118,0,5,'gear/fitness-equipment'),
	(14,45,0,6,'Watches'),
	(15,117,0,6,'watches'),
	(16,118,0,6,'gear/watches'),
	(17,45,0,7,'Collections'),
	(18,52,0,7,'PAGE'),
	(19,117,0,7,'collections'),
	(20,118,0,7,'collections'),
	(21,45,0,8,'New Luma Yoga Collection'),
	(22,117,0,8,'yoga-new'),
	(23,118,0,8,'collections/yoga-new'),
	(24,45,0,9,'Training'),
	(25,52,0,9,'PAGE'),
	(26,117,0,9,'training'),
	(27,118,0,9,'training'),
	(28,45,0,10,'Video Download'),
	(29,117,0,10,'training-video'),
	(30,118,0,10,'training/training-video'),
	(31,45,0,11,'Men'),
	(32,52,0,11,'PAGE'),
	(33,117,0,11,'men'),
	(34,118,0,11,'men'),
	(35,45,0,12,'Tops'),
	(36,117,0,12,'tops-men'),
	(37,118,0,12,'men/tops-men'),
	(38,45,0,13,'Bottoms'),
	(39,117,0,13,'bottoms-men'),
	(40,118,0,13,'men/bottoms-men'),
	(41,45,0,14,'Jackets'),
	(42,117,0,14,'jackets-men'),
	(43,118,0,14,'men/tops-men/jackets-men'),
	(44,45,0,15,'Hoodies & Sweatshirts'),
	(45,117,0,15,'hoodies-and-sweatshirts-men'),
	(46,118,0,15,'men/tops-men/hoodies-and-sweatshirts-men'),
	(47,45,0,16,'Tees'),
	(48,117,0,16,'tees-men'),
	(49,118,0,16,'men/tops-men/tees-men'),
	(50,45,0,17,'Tanks'),
	(51,117,0,17,'tanks-men'),
	(52,118,0,17,'men/tops-men/tanks-men'),
	(53,45,0,18,'Pants'),
	(54,117,0,18,'pants-men'),
	(55,118,0,18,'men/bottoms-men/pants-men'),
	(56,45,0,19,'Shorts'),
	(57,117,0,19,'shorts-men'),
	(58,118,0,19,'men/bottoms-men/shorts-men'),
	(59,45,0,20,'Women'),
	(60,52,0,20,'PAGE'),
	(61,117,0,20,'women'),
	(62,118,0,20,'women'),
	(63,45,0,21,'Tops'),
	(64,117,0,21,'tops-women'),
	(65,118,0,21,'women/tops-women'),
	(66,45,0,22,'Bottoms'),
	(67,117,0,22,'bottoms-women'),
	(68,118,0,22,'women/bottoms-women'),
	(69,45,0,23,'Jackets'),
	(70,117,0,23,'jackets-women'),
	(71,118,0,23,'women/tops-women/jackets-women'),
	(72,45,0,24,'Hoodies & Sweatshirts'),
	(73,117,0,24,'hoodies-and-sweatshirts-women'),
	(74,118,0,24,'women/tops-women/hoodies-and-sweatshirts-women'),
	(75,45,0,25,'Tees'),
	(76,117,0,25,'tees-women'),
	(77,118,0,25,'women/tops-women/tees-women'),
	(78,45,0,26,'Bras & Tanks'),
	(79,117,0,26,'tanks-women'),
	(80,118,0,26,'women/tops-women/tanks-women'),
	(81,45,0,27,'Pants'),
	(82,117,0,27,'pants-women'),
	(83,118,0,27,'women/bottoms-women/pants-women'),
	(84,45,0,28,'Shorts'),
	(85,117,0,28,'shorts-women'),
	(86,118,0,28,'women/bottoms-women/shorts-women'),
	(87,45,0,29,'Promotions'),
	(88,52,0,29,'PAGE'),
	(89,117,0,29,'promotions'),
	(90,118,0,29,'promotions'),
	(91,45,0,30,'Women Sale'),
	(92,117,0,30,'women-sale'),
	(93,118,0,30,'promotions/women-sale'),
	(94,45,0,31,'Men Sale'),
	(95,117,0,31,'men-sale'),
	(96,118,0,31,'promotions/men-sale'),
	(97,45,0,32,'Pants'),
	(98,117,0,32,'pants-all'),
	(99,118,0,32,'promotions/pants-all'),
	(100,45,0,33,'Tees'),
	(101,117,0,33,'tees-all'),
	(102,118,0,33,'promotions/tees-all'),
	(103,45,0,34,'Erin Recommends'),
	(104,117,0,34,'erin-recommends'),
	(105,118,0,34,'collections/erin-recommends'),
	(106,45,0,35,'Performance Fabrics'),
	(107,117,0,35,'performance-fabrics'),
	(108,118,0,35,'collections/performance-fabrics'),
	(109,45,0,36,'Eco Friendly'),
	(110,117,0,36,'eco-friendly'),
	(111,118,0,36,'collections/eco-friendly'),
	(112,45,0,37,'Sale'),
	(113,52,0,37,'PAGE'),
	(114,117,0,37,'sale'),
	(115,118,0,37,'sale'),
	(116,45,0,38,'What\'s New'),
	(117,52,0,38,'PAGE'),
	(118,117,0,38,'what-is-new'),
	(119,118,0,38,'what-is-new'),
	(120,45,0,39,'Performance Sportswear New'),
	(121,52,0,39,'PAGE'),
	(122,63,0,39,'1column'),
	(123,117,0,39,'performance-new'),
	(124,118,0,39,'collections/performance-new'),
	(125,45,0,40,'Eco Collection New'),
	(126,52,0,40,'PAGE'),
	(127,63,0,40,'1column'),
	(128,117,0,40,'eco-new'),
	(129,118,0,40,'collections/eco-new'),
	(130,48,2,6,NULL),
	(131,60,2,6,NULL),
	(132,63,2,6,NULL),
	(133,118,2,6,'gear/watches'),
	(134,45,1,6,'Watches StoreID_1');

/*!40000 ALTER TABLE `catalog_category_entity_varchar` ENABLE KEYS */;
UNLOCK TABLES;
SET foreign_key_checks = 1;
