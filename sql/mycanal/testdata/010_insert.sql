/*
 * Copyright © Magento, Inc. All rights reserved.
 * Licensed under OSL 3.0
 * http://opensource.org/licenses/osl-3.0.php Open Software License (OSL 3.0)
*/

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

# Dump of table catalog_product_entity
# ------------------------------------------------------------

LOCK TABLES `catalog_product_entity` WRITE;
/*!40000 ALTER TABLE `catalog_product_entity` DISABLE KEYS */;

INSERT INTO `catalog_product_entity` (`entity_id`, `attribute_set_id`, `type_id`, `sku`, `has_options`, `required_options`, `created_at`, `updated_at`)
VALUES
	(44,11,'simple','24-WG02',0,0,'2018-04-17 21:42:01','2018-04-17 21:42:01'),
	(45,11,'bundle','24-WG080',1,1,'2018-04-17 21:42:02','2018-04-17 21:42:02'),
	(46,14,'downloadable','240-LV04',0,0,'2018-04-17 21:42:03','2018-04-17 21:42:03'),
	(52,9,'simple','MH01-XS-Black',0,0,'2018-04-17 21:42:20','2018-04-17 21:42:20'),
	(53,9,'simple','MH01-XS-Gray',0,0,'2018-04-17 21:42:20','2018-04-17 21:42:20'),
	(54,9,'simple','MH01-XS-Orange',0,0,'2018-04-17 21:42:21','2018-04-17 21:42:21'),
	(55,9,'simple','MH01-S-Black',0,0,'2018-04-17 21:42:21','2018-04-17 21:42:21'),
	(56,9,'simple','MH01-S-Gray',0,0,'2018-04-17 21:42:21','2018-04-17 21:42:21'),
	(57,9,'simple','MH01-S-Orange',0,0,'2018-04-17 21:42:21','2018-04-17 21:42:21'),
	(58,9,'simple','MH01-M-Black',0,0,'2018-04-17 21:42:21','2018-04-17 21:42:21'),
	(59,9,'simple','MH01-M-Gray',0,0,'2018-04-17 21:42:21','2018-04-17 21:42:21'),
	(60,9,'simple','MH01-M-Orange',0,0,'2018-04-17 21:42:21','2018-04-17 21:42:21'),
	(61,9,'simple','MH01-L-Black',0,0,'2018-04-17 21:42:21','2018-04-17 21:42:21'),
	(62,9,'simple','MH01-L-Gray',0,0,'2018-04-17 21:42:21','2018-04-17 21:42:21'),
	(63,9,'simple','MH01-L-Orange',0,0,'2018-04-17 21:42:21','2018-04-17 21:42:21'),
	(64,9,'simple','MH01-XL-Black',0,0,'2018-04-17 21:42:21','2018-04-17 21:42:21'),
	(65,9,'simple','MH01-XL-Gray',0,0,'2018-04-17 21:42:21','2018-04-17 21:42:21'),
	(66,9,'simple','MH01-XL-Orange',0,0,'2018-04-17 21:42:21','2018-04-17 21:42:21'),
	(67,9,'configurable','MH01',1,0,'2018-04-17 21:42:21','2018-04-17 21:42:21');

/*!40000 ALTER TABLE `catalog_product_entity` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table catalog_product_entity_datetime
# ------------------------------------------------------------



# Dump of table catalog_product_entity_decimal
# ------------------------------------------------------------

LOCK TABLES `catalog_product_entity_decimal` WRITE;
/*!40000 ALTER TABLE `catalog_product_entity_decimal` DISABLE KEYS */;

INSERT INTO `catalog_product_entity_decimal` (`value_id`, `attribute_id`, `store_id`, `entity_id`, `value`)
VALUES
	(50,77,0,44,92.0000),
	(51,77,0,46,6.0000),
	(57,77,0,52,52.0000),
	(58,82,0,52,1.0000),
	(59,77,0,53,52.0000),
	(60,82,0,53,1.0000),
	(61,77,0,54,52.0000),
	(62,82,0,54,1.0000),
	(63,77,0,55,52.0000),
	(64,82,0,55,1.0000),
	(65,77,0,56,52.0000),
	(66,82,0,56,1.0000),
	(67,77,0,57,52.0000),
	(68,82,0,57,1.0000),
	(69,77,0,58,52.0000),
	(70,82,0,58,1.0000),
	(71,77,0,59,52.0000),
	(72,82,0,59,1.0000),
	(73,77,0,60,52.0000),
	(74,82,0,60,1.0000),
	(75,77,0,61,52.0000),
	(76,82,0,61,1.0000),
	(77,77,0,62,52.0000),
	(78,82,0,62,1.0000),
	(79,77,0,63,52.0000),
	(80,82,0,63,1.0000),
	(81,77,0,64,52.0000),
	(82,82,0,64,1.0000),
	(83,77,0,65,52.0000),
	(84,82,0,65,1.0000),
	(85,77,0,66,52.0000),
	(86,82,0,66,1.0000),
	(87,77,0,67,52.0000);

/*!40000 ALTER TABLE `catalog_product_entity_decimal` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table catalog_product_entity_int
# ------------------------------------------------------------

LOCK TABLES `catalog_product_entity_int` WRITE;
/*!40000 ALTER TABLE `catalog_product_entity_int` DISABLE KEYS */;

INSERT INTO `catalog_product_entity_int` (`value_id`, `attribute_id`, `store_id`, `entity_id`, `value`)
VALUES
	(211,97,0,44,1),
	(212,99,0,44,4),
	(213,115,0,44,1),
	(214,134,0,44,2),
	(215,146,0,44,1),
	(216,97,0,45,1),
	(217,99,0,45,4),
	(218,115,0,45,1),
	(219,121,0,45,0),
	(220,122,0,45,0),
	(221,123,0,45,0),
	(222,124,0,45,0),
	(223,125,0,45,0),
	(224,134,0,45,2),
	(225,97,0,46,1),
	(226,99,0,46,4),
	(227,115,0,46,1),
	(228,126,0,46,0),
	(229,134,0,46,2),
	(230,148,0,46,102),
	(261,97,0,52,1),
	(262,134,0,52,2),
	(263,115,0,52,1),
	(264,99,0,52,1),
	(265,93,0,52,49),
	(266,142,0,52,167),
	(267,97,0,53,1),
	(268,134,0,53,2),
	(269,115,0,53,1),
	(270,99,0,53,1),
	(271,93,0,53,52),
	(272,142,0,53,167),
	(273,97,0,54,1),
	(274,134,0,54,2),
	(275,115,0,54,1),
	(276,99,0,54,1),
	(277,93,0,54,56),
	(278,142,0,54,167),
	(279,97,0,55,1),
	(280,134,0,55,2),
	(281,115,0,55,1),
	(282,99,0,55,1),
	(283,93,0,55,49),
	(284,142,0,55,168),
	(285,97,0,56,1),
	(286,134,0,56,2),
	(287,115,0,56,1),
	(288,99,0,56,1),
	(289,93,0,56,52),
	(290,142,0,56,168),
	(291,97,0,57,1),
	(292,134,0,57,2),
	(293,115,0,57,1),
	(294,99,0,57,1),
	(295,93,0,57,56),
	(296,142,0,57,168),
	(297,97,0,58,1),
	(298,134,0,58,2),
	(299,115,0,58,1),
	(300,99,0,58,1),
	(301,93,0,58,49),
	(302,142,0,58,169),
	(303,97,0,59,1),
	(304,134,0,59,2),
	(305,115,0,59,1),
	(306,99,0,59,1),
	(307,93,0,59,52),
	(308,142,0,59,169),
	(309,97,0,60,1),
	(310,134,0,60,2),
	(311,115,0,60,1),
	(312,99,0,60,1),
	(313,93,0,60,56),
	(314,142,0,60,169),
	(315,97,0,61,1),
	(316,134,0,61,2),
	(317,115,0,61,1),
	(318,99,0,61,1),
	(319,93,0,61,49),
	(320,142,0,61,170),
	(321,97,0,62,1),
	(322,134,0,62,2),
	(323,115,0,62,1),
	(324,99,0,62,1),
	(325,93,0,62,52),
	(326,142,0,62,170),
	(327,97,0,63,1),
	(328,134,0,63,2),
	(329,115,0,63,1),
	(330,99,0,63,1),
	(331,93,0,63,56),
	(332,142,0,63,170),
	(333,97,0,64,1),
	(334,134,0,64,2),
	(335,115,0,64,1),
	(336,99,0,64,1),
	(337,93,0,64,49),
	(338,142,0,64,171),
	(339,97,0,65,1),
	(340,134,0,65,2),
	(341,115,0,65,1),
	(342,99,0,65,1),
	(343,93,0,65,52),
	(344,142,0,65,171),
	(345,97,0,66,1),
	(346,134,0,66,2),
	(347,115,0,66,1),
	(348,99,0,66,1),
	(349,93,0,66,56),
	(350,142,0,66,171),
	(351,97,0,67,1),
	(352,134,0,67,2),
	(353,115,0,67,1),
	(354,99,0,67,4),
	(355,143,0,67,1),
	(356,144,0,67,0),
	(357,145,0,67,0),
	(358,146,0,67,0),
	(359,147,0,67,1);

/*!40000 ALTER TABLE `catalog_product_entity_int` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table catalog_product_entity_text
# ------------------------------------------------------------

LOCK TABLES `catalog_product_entity_text` WRITE;
/*!40000 ALTER TABLE `catalog_product_entity_text` DISABLE KEYS */;

INSERT INTO `catalog_product_entity_text` (`value_id`, `attribute_id`, `store_id`, `entity_id`, `value`)
VALUES
	(44,75,0,44,'<p>The Didi Sport Watch helps you keep your workout plan down to the second.</li>\n</ul>'),
	(45,75,0,45,'<p>A well-rounded yoga workout takes more than a mat.\n</ul>'),
	(46,75,0,46,'<p>Beginner\'s Yoga starts you down the path toward strength, balance and mental focus.</li>\n</ul>'),
	(47,76,0,46,'<p>\nThe most difficult yoga poses to master are the ones learned incorrectly as a beginner.\n</p>'),
	(58,75,0,52,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(59,75,0,53,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(60,75,0,54,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(61,75,0,55,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(62,75,0,56,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(63,75,0,57,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(64,75,0,58,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(65,75,0,59,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(66,75,0,60,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(67,75,0,61,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(68,75,0,62,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(69,75,0,63,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(70,75,0,64,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(71,75,0,65,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(72,75,0,66,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>'),
	(73,75,0,67,'<p>Ideal for cold-weather training or work outdoors, the Chaz Hoodie promises superior warmth with every wear. Thick material blocks out the wind as ribbed cuffs and bottom band seal in body heat.</p>\n<p>&bull; Two-tone gray heather hoodie.<br />&bull; Drawstring-adjustable hood. <br />&bull; Machine wash/dry.</p>');

/*!40000 ALTER TABLE `catalog_product_entity_text` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table catalog_product_entity_varchar
# ------------------------------------------------------------

LOCK TABLES `catalog_product_entity_varchar` WRITE;
/*!40000 ALTER TABLE `catalog_product_entity_varchar` DISABLE KEYS */;

INSERT INTO `catalog_product_entity_varchar` (`value_id`, `attribute_id`, `store_id`, `entity_id`, `value`)
VALUES
	(436,73,0,44,'Didi Sport Watch'),
	(437,106,0,44,'container2'),
	(438,119,0,44,'didi-sport-watch'),
	(439,135,0,44,'16,11'),
	(440,137,0,44,'43,45,48'),
	(441,140,0,44,'81'),
	(442,141,0,44,'86,87,90'),
	(443,87,0,44,'/w/g/wg02-bk-0.jpg'),
	(444,88,0,44,'/w/g/wg02-bk-0.jpg'),
	(445,89,0,44,'/w/g/wg02-bk-0.jpg'),
	(446,73,0,45,'Sprite Yoga Companion Kit'),
	(447,106,0,45,'container2'),
	(448,119,0,45,'sprite-yoga-companion-kit'),
	(449,135,0,45,'8,11'),
	(450,140,0,45,'80,81,84'),
	(451,141,0,45,'87'),
	(452,87,0,45,'/l/u/luma-yoga-kit-2.jpg'),
	(453,88,0,45,'/l/u/luma-yoga-kit-2.jpg'),
	(454,89,0,45,'/l/u/luma-yoga-kit-2.jpg'),
	(455,73,0,46,'Beginner\'s Yoga'),
	(456,106,0,46,'container2'),
	(457,119,0,46,'beginner-s-yoga'),
	(458,127,0,46,'Trailers'),
	(459,128,0,46,'Downloads'),
	(460,135,0,46,'8,16,17,5,11'),
	(461,87,0,46,'/l/t/lt01.jpg'),
	(462,88,0,46,'/l/t/lt01.jpg'),
	(463,89,0,46,'/l/t/lt01.jpg'),
	(510,106,0,52,'container2'),
	(511,132,0,52,'0'),
	(512,119,0,52,'chaz-kangeroo-hoodie-xs-black'),
	(513,87,0,52,'/m/h/mh01-black_main_1.jpg'),
	(514,88,0,52,'/m/h/mh01-black_main_1.jpg'),
	(515,89,0,52,'/m/h/mh01-black_main_1.jpg'),
	(516,73,0,52,'Chaz Kangeroo Hoodie-XS-Black'),
	(517,106,0,53,'container2'),
	(518,132,0,53,'0'),
	(519,119,0,53,'chaz-kangeroo-hoodie-xs-gray'),
	(520,87,0,53,'/m/h/mh01-gray_main_1.jpg'),
	(521,88,0,53,'/m/h/mh01-gray_main_1.jpg'),
	(522,89,0,53,'/m/h/mh01-gray_main_1.jpg'),
	(523,73,0,53,'Chaz Kangeroo Hoodie-XS-Gray'),
	(524,106,0,54,'container2'),
	(525,132,0,54,'0'),
	(526,119,0,54,'chaz-kangeroo-hoodie-xs-orange'),
	(527,87,0,54,'/m/h/mh01-orange_main_1.jpg'),
	(528,88,0,54,'/m/h/mh01-orange_main_1.jpg'),
	(529,89,0,54,'/m/h/mh01-orange_main_1.jpg'),
	(530,73,0,54,'Chaz Kangeroo Hoodie-XS-Orange'),
	(531,106,0,55,'container2'),
	(532,132,0,55,'0'),
	(533,119,0,55,'chaz-kangeroo-hoodie-s-black'),
	(534,87,0,55,'/m/h/mh01-black_main_1.jpg'),
	(535,88,0,55,'/m/h/mh01-black_main_1.jpg'),
	(536,89,0,55,'/m/h/mh01-black_main_1.jpg'),
	(537,73,0,55,'Chaz Kangeroo Hoodie-S-Black'),
	(538,106,0,56,'container2'),
	(539,132,0,56,'0'),
	(540,119,0,56,'chaz-kangeroo-hoodie-s-gray'),
	(541,87,0,56,'/m/h/mh01-gray_main_1.jpg'),
	(542,88,0,56,'/m/h/mh01-gray_main_1.jpg'),
	(543,89,0,56,'/m/h/mh01-gray_main_1.jpg'),
	(544,73,0,56,'Chaz Kangeroo Hoodie-S-Gray'),
	(545,106,0,57,'container2'),
	(546,132,0,57,'0'),
	(547,119,0,57,'chaz-kangeroo-hoodie-s-orange'),
	(548,87,0,57,'/m/h/mh01-orange_main_1.jpg'),
	(549,88,0,57,'/m/h/mh01-orange_main_1.jpg'),
	(550,89,0,57,'/m/h/mh01-orange_main_1.jpg'),
	(551,73,0,57,'Chaz Kangeroo Hoodie-S-Orange'),
	(552,106,0,58,'container2'),
	(553,132,0,58,'0'),
	(554,119,0,58,'chaz-kangeroo-hoodie-m-black'),
	(555,87,0,58,'/m/h/mh01-black_main_1.jpg'),
	(556,88,0,58,'/m/h/mh01-black_main_1.jpg'),
	(557,89,0,58,'/m/h/mh01-black_main_1.jpg'),
	(558,73,0,58,'Chaz Kangeroo Hoodie-M-Black'),
	(559,106,0,59,'container2'),
	(560,132,0,59,'0'),
	(561,119,0,59,'chaz-kangeroo-hoodie-m-gray'),
	(562,87,0,59,'/m/h/mh01-gray_main_1.jpg'),
	(563,88,0,59,'/m/h/mh01-gray_main_1.jpg'),
	(564,89,0,59,'/m/h/mh01-gray_main_1.jpg'),
	(565,73,0,59,'Chaz Kangeroo Hoodie-M-Gray'),
	(566,106,0,60,'container2'),
	(567,132,0,60,'0'),
	(568,119,0,60,'chaz-kangeroo-hoodie-m-orange'),
	(569,87,0,60,'/m/h/mh01-orange_main_1.jpg'),
	(570,88,0,60,'/m/h/mh01-orange_main_1.jpg'),
	(571,89,0,60,'/m/h/mh01-orange_main_1.jpg'),
	(572,73,0,60,'Chaz Kangeroo Hoodie-M-Orange'),
	(573,106,0,61,'container2'),
	(574,132,0,61,'0'),
	(575,119,0,61,'chaz-kangeroo-hoodie-l-black'),
	(576,87,0,61,'/m/h/mh01-black_main_1.jpg'),
	(577,88,0,61,'/m/h/mh01-black_main_1.jpg'),
	(578,89,0,61,'/m/h/mh01-black_main_1.jpg'),
	(579,73,0,61,'Chaz Kangeroo Hoodie-L-Black'),
	(580,106,0,62,'container2'),
	(581,132,0,62,'0'),
	(582,119,0,62,'chaz-kangeroo-hoodie-l-gray'),
	(583,87,0,62,'/m/h/mh01-gray_main_1.jpg'),
	(584,88,0,62,'/m/h/mh01-gray_main_1.jpg'),
	(585,89,0,62,'/m/h/mh01-gray_main_1.jpg'),
	(586,73,0,62,'Chaz Kangeroo Hoodie-L-Gray'),
	(587,106,0,63,'container2'),
	(588,132,0,63,'0'),
	(589,119,0,63,'chaz-kangeroo-hoodie-l-orange'),
	(590,87,0,63,'/m/h/mh01-orange_main_1.jpg'),
	(591,88,0,63,'/m/h/mh01-orange_main_1.jpg'),
	(592,89,0,63,'/m/h/mh01-orange_main_1.jpg'),
	(593,73,0,63,'Chaz Kangeroo Hoodie-L-Orange'),
	(594,106,0,64,'container2'),
	(595,132,0,64,'0'),
	(596,119,0,64,'chaz-kangeroo-hoodie-xl-black'),
	(597,87,0,64,'/m/h/mh01-black_main_1.jpg'),
	(598,88,0,64,'/m/h/mh01-black_main_1.jpg'),
	(599,89,0,64,'/m/h/mh01-black_main_1.jpg'),
	(600,73,0,64,'Chaz Kangeroo Hoodie-XL-Black'),
	(601,106,0,65,'container2'),
	(602,132,0,65,'0'),
	(603,119,0,65,'chaz-kangeroo-hoodie-xl-gray'),
	(604,87,0,65,'/m/h/mh01-gray_main_1.jpg'),
	(605,88,0,65,'/m/h/mh01-gray_main_1.jpg'),
	(606,89,0,65,'/m/h/mh01-gray_main_1.jpg'),
	(607,73,0,65,'Chaz Kangeroo Hoodie-XL-Gray'),
	(608,106,0,66,'container2'),
	(609,132,0,66,'0'),
	(610,119,0,66,'chaz-kangeroo-hoodie-xl-orange'),
	(611,87,0,66,'/m/h/mh01-orange_main_1.jpg'),
	(612,88,0,66,'/m/h/mh01-orange_main_1.jpg'),
	(613,89,0,66,'/m/h/mh01-orange_main_1.jpg'),
	(614,73,0,66,'Chaz Kangeroo Hoodie-XL-Orange'),
	(615,106,0,67,'container2'),
	(616,132,0,67,'0'),
	(617,119,0,67,'chaz-kangeroo-hoodie'),
	(618,87,0,67,'/m/h/mh01-gray_main_1.jpg'),
	(619,88,0,67,'/m/h/mh01-gray_main_1.jpg'),
	(620,89,0,67,'/m/h/mh01-gray_main_1.jpg'),
	(621,73,0,67,'Chaz Kangeroo Hoodie'),
	(622,137,0,67,'159'),
	(623,153,0,67,'195'),
	(624,154,0,67,'202,204,205,208,210');

/*!40000 ALTER TABLE `catalog_product_entity_varchar` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
