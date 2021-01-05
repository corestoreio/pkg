SET foreign_key_checks = 0;

DROP TABLE IF EXISTS `customer_entity`;
DROP TABLE IF EXISTS `customer_address_entity`;
DROP TABLE IF EXISTS `customer_entity_varchar`;
DROP TABLE IF EXISTS `customer_entity_int`;

-- Create syntax for TABLE 'customer_entity'
CREATE TABLE `customer_entity` (
  `entity_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Entity ID',
  `website_id` smallint(5) unsigned DEFAULT NULL COMMENT 'Website ID',
  `email` varchar(255) DEFAULT NULL COMMENT 'Email',
  PRIMARY KEY (`entity_id`),
  UNIQUE KEY `CUSTOMER_ENTITY_EMAIL_WEBSITE_ID` (`email`,`website_id`),
  KEY `CUSTOMER_ENTITY_WEBSITE_ID` (`website_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Customer Entity';

-- Create syntax for TABLE 'customer_address_entity'
CREATE TABLE `customer_address_entity` (
  `entity_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Entity ID',
  `parent_id` int(10) unsigned DEFAULT NULL COMMENT 'Parent ID',
  `city` varchar(255) NOT NULL COMMENT 'City',
  `company` varchar(255) DEFAULT NULL COMMENT 'Company',
  `firstname` varchar(255) NOT NULL COMMENT 'First Name',
  `lastname` varchar(255) NOT NULL COMMENT 'Last Name',
  PRIMARY KEY (`entity_id`),
  KEY `CUSTOMER_ADDRESS_ENTITY_PARENT_ID` (`parent_id`),
  CONSTRAINT `CUSTOMER_ADDRESS_ENTITY_PARENT_ID_CUSTOMER_ENTITY_ENTITY_ID` FOREIGN KEY (`parent_id`) REFERENCES `customer_entity` (`entity_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Customer Address Entity';

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

INSERT INTO `customer_address_entity` (`entity_id`, `parent_id`, `city`, `company`, `firstname`, `lastname`)
VALUES
	(1,1,'Ransbach-Baumbach','Beck & Co. Einzelunternehmen','Luis','Kirschner'),
	(2,1,'Herzogenaurach','IBM PartG','Jonathan','Schaefer'),
	(3,1,'Runkel','China Construction Bank VVaG','Elif','Kurth'),
	(4,1,'Neckargemünd','Audi UG (haftungsbeschränkt)','Amalia','Stenzel');

INSERT INTO `customer_entity` (`entity_id`, `website_id`, `email`)
VALUES
	(1,2,'janiceschott@procter--gamble.prudential');

INSERT INTO `customer_entity` (`entity_id`, `website_id`, `email`)
VALUES
	(2,1,'info@paldi.de');

SET foreign_key_checks = 1;
