-- INSERT INTO `catalog_product_index_price_final_idx`
SELECT `e`.`entity_id`,
       `cg`.`customer_group_id`,
       `cw`.`website_id`,
       IF(IFNULL(tas_tax_class_id.value_id, -1) > 0, tas_tax_class_id.value, tad_tax_class_id.value) AS `tax_class_id`,
       IF(IFNULL(tas_price.value_id, -1) > 0, tas_price.value, tad_price.value) AS `orig_price`,
       IF(IF(gp.price IS NULL, IF(IFNULL(tas_price.value_id, -1) > 0, tas_price.value, tad_price.value), gp.price) < IF(IF(IF(IFNULL(tas_special_from_date.value_id, -1) > 0,
             tas_special_from_date.value,
             tad_special_from_date.value) IS NULL, 1, IF(
             Date(IF(IFNULL(tas_special_from_date.value_id, -1) > 0,
             tas_special_from_date.value,
             tad_special_from_date.value)) <= Date
             (cwd.website_date), 1, 0)) > 0
             AND IF(IF(IFNULL(tas_special_to_date.value_id, -1) > 0,
                    tas_special_to_date.value,
                           tad_special_to_date.value) IS NULL, 1, IF(Date(
                     IF(IFNULL(tas_special_to_date.value_id, -1) > 0,
             tas_special_to_date.value,
             tad_special_to_date.value)) >= Date(cwd.website_date), 1, 0)) > 0
             AND IF(IFNULL(tas_special_price.value_id, -1) > 0,
                 tas_special_price.value,
                     tad_special_price.value) <
                 IF(IFNULL(tas_price.value_id, -1) > 0, tas_price.value,
                 tad_price.value), IF(IFNULL(tas_special_price.value_id, -1) > 0, tas_special_price.value,
          tad_special_price.value), IF(IFNULL(tas_price.value_id, -1) > 0,
                                    tas_price.value, tad_price.value)), IF(
       gp.price IS NULL, IF(IFNULL(tas_price.value_id, -1) > 0, tas_price.value,
       tad_price.value), gp.price), IF(IF(
       IF(IFNULL(tas_special_from_date.value_id, -1) > 0,
                                       tas_special_from_date.value,
                                          tad_special_from_date.value) IS NULL, 1 ,
                                    IF(Date(IF(IFNULL( tas_special_from_date.value_id, -1) > 0, tas_special_from_date.value, tad_special_from_date.value)) <= Date( cwd.website_date), 1, 0)) > 0
                                       AND IF(IF(IFNULL( tas_special_to_date.value_id, -1) > 0,
                                              tas_special_to_date.value,tad_special_to_date.value)
                                              IS NULL , 1, IF(Date( IF(IFNULL( tas_special_to_date.value_id,
                                                  -1) > 0,
                                       tas_special_to_date.value, tad_special_to_date.value)) >= Date( cwd.website_date), 1, 0)) > 0
                                       AND IF(IFNULL(tas_special_price.value_id,
                                              -1) > 0, tas_special_price.value,
                                               tad_special_price.value) <
                                           IF(IFNULL(tas_price.value_id, -1) > 0, tas_price.value, tad_price.value), IF
        (IFNULL(tas_special_price.value_id, -1) > 0,tas_special_price.value,tad_special_price.value),
            IF(IFNULL(tas_price.value_id, -1) > 0, tas_price.value, tad_price.value)))
       AS `price`,

       IF(IF(gp.price IS NULL, IF(IFNULL(tas_price.value_id, -1) > 0,
                               tas_price.value,
                                  tad_price.value), gp.price) <
             IF(IF(IF(IFNULL(tas_special_from_date.value_id, -1) > 0,
             tas_special_from_date.value,
             tad_special_from_date.value) IS NULL, 1, IF(
             Date(IF(IFNULL(tas_special_from_date.value_id, -1) > 0,
             tas_special_from_date.value,
             tad_special_from_date.value)) <= Date
             (cwd.website_date), 1, 0)) > 0
             AND IF(IF(IFNULL(tas_special_to_date.value_id, -1) > 0,
                    tas_special_to_date.value,
                           tad_special_to_date.value) IS NULL, 1, IF(Date(
                     IF(IFNULL(tas_special_to_date.value_id, -1) > 0,
             tas_special_to_date.value,
             tad_special_to_date.value)) >= Date(cwd.website_date), 1, 0)) > 0
             AND IF(IFNULL(tas_special_price.value_id, -1) > 0,
                 tas_special_price.value,
                     tad_special_price.value) <
                 IF(
                 IFNULL(tas_price.value_id, -1) > 0,
                                                tas_price.value,
                 tad_price.value), IF(
             IFNULL
             (tas_special_price.value_id, -1) > 0, tas_special_price.value,
          tad_special_price.value), IF(IFNULL(tas_price.value_id, -1) > 0,
                                    tas_price.value, tad_price.value)), IF(
       gp.price IS NULL, IF(IFNULL(tas_price.value_id, -1) > 0, tas_price.value,
       tad_price.value), gp.price), IF(IF(
       IF(IFNULL(tas_special_from_date.value_id, -1) > 0,
                                       tas_special_from_date.value,
                                          tad_special_from_date.value) IS NULL, 1 ,
                                    IF(Date(IF(IFNULL(
                                       tas_special_from_date.value_id, -1) > 0,
                                       tas_special_from_date.value,
                                       tad_special_from_date.value)) <= Date(
                                       cwd.website_date), 1, 0)) > 0
                                       AND IF(IF(IFNULL( tas_special_to_date.value_id,
                                                 -1) > 0, tas_special_to_date.value, tad_special_to_date.value)
                                              IS NULL , 1, IF(Date( IF(IFNULL( tas_special_to_date.value_id, -1) > 0,
                                       tas_special_to_date.value,
                                       tad_special_to_date.value)) >= Date(
                                                   cwd.website_date), 1, 0)) > 0
                                       AND IF(IFNULL(tas_special_price.value_id,
                                              -1) >
                                              0, tas_special_price.value,
                                               tad_special_price.value) <
                                           IF(IFNULL(tas_price.value_id, -1) > 0
                                           ,
                                           tas_price.value, tad_price.value), IF
                                    (
                                                                 IFNULL(
                                    tas_special_price.value_id, -1) > 0,
       tas_special_price.value,
       tad_special_price.value), IF(IFNULL(tas_price.value_id, -1) > 0,
       tas_price.value, tad_price.value)))
       AS `min_price`,

       IF(IF(gp.price IS NULL, IF(IFNULL(tas_price.value_id, -1) > 0,
                               tas_price.value,
                                  tad_price.value), gp.price) <
             IF(IF(IF(IFNULL(tas_special_from_date.value_id, -1) > 0,
             tas_special_from_date.value,
             tad_special_from_date.value) IS NULL, 1, IF(
             Date(IF(IFNULL(tas_special_from_date.value_id, -1) > 0,
             tas_special_from_date.value,
             tad_special_from_date.value)) <= Date
             (cwd.website_date), 1, 0)) > 0
             AND IF(IF(IFNULL(tas_special_to_date.value_id, -1) > 0,
                    tas_special_to_date.value,
                           tad_special_to_date.value) IS NULL, 1, IF(Date(
                     IF(IFNULL(tas_special_to_date.value_id, -1) > 0,
             tas_special_to_date.value,
             tad_special_to_date.value)) >= Date(cwd.website_date), 1, 0)) > 0
             AND IF(IFNULL(tas_special_price.value_id, -1) > 0,
                 tas_special_price.value,
                     tad_special_price.value) <
                 IF(
                 IFNULL(tas_price.value_id, -1) > 0,
                                                tas_price.value,
                 tad_price.value), IF(
             IFNULL
             (tas_special_price.value_id, -1) > 0, tas_special_price.value,
          tad_special_price.value), IF(IFNULL(tas_price.value_id, -1) > 0,
                                    tas_price.value, tad_price.value)), IF(
       gp.price IS NULL, IF(IFNULL(tas_price.value_id, -1) > 0, tas_price.value,
       tad_price.value), gp.price), IF(IF(
       IF(IFNULL(tas_special_from_date.value_id, -1) > 0,
                                       tas_special_from_date.value,
                                          tad_special_from_date.value) IS NULL,
                                       1
                                 ,
                                    IF(Date(IF(IFNULL(
                                       tas_special_from_date.value_id, -1) > 0,
                                       tas_special_from_date.value,
                                       tad_special_from_date.value)) <= Date(
                                       cwd.website_date), 1, 0)) > 0
                                       AND IF(IF(IFNULL(
tas_special_to_date.value_id,
                                                 -1) > 0,
                                              tas_special_to_date.value,
                                                     tad_special_to_date.value)
                                              IS NULL
                                           , 1, IF(Date(
                                               IF(IFNULL(
tas_special_to_date.value_id,
                                                  -1) > 0,
                                       tas_special_to_date.value,
                                       tad_special_to_date.value)) >= Date(
                                                   cwd.website_date), 1, 0)) > 0
                                       AND IF(IFNULL(tas_special_price.value_id,
                                              -1) >
                                              0, tas_special_price.value,
                                               tad_special_price.value) <
                                           IF(IFNULL(tas_price.value_id, -1) > 0
                                           ,
                                           tas_price.value, tad_price.value), IF
                                    (
                                                                 IFNULL(
                                    tas_special_price.value_id, -1) > 0,
       tas_special_price.value,
       tad_special_price.value), IF(IFNULL(tas_price.value_id, -1) > 0,
       tas_price.value, tad_price.value))) AS `max_price`,

       tp.min_price AS `tier_price`,
       tp.min_price AS `base_tier`,
       gp.price AS `group_price`,
       gp.price AS `base_group_price`

FROM   `catalog_product_entity` AS `e`
       CROSS JOIN `customer_group` AS `cg`
       CROSS JOIN `core_website` AS `cw`
       INNER JOIN `catalog_product_index_website` AS `cwd`
               ON cw.website_id = cwd.website_id
       INNER JOIN `core_store_group` AS `csg`
               ON csg.website_id = cw.website_id
                  AND cw.default_group_id = csg.group_id
       INNER JOIN `core_store` AS `cs`
               ON csg.default_store_id = cs.store_id
                  AND cs.store_id != 0
       INNER JOIN `catalog_product_website` AS `pw`
               ON pw.product_id = e.entity_id
                  AND pw.website_id = cw.website_id
       LEFT JOIN `catalog_product_index_tier_price` AS `tp`
              ON tp.entity_id = e.entity_id
                 AND tp.website_id = cw.website_id
                 AND tp.customer_group_id = cg.customer_group_id
       LEFT JOIN `catalog_product_index_group_price` AS `gp`
              ON gp.entity_id = e.entity_id
                 AND gp.website_id = cw.website_id
                 AND gp.customer_group_id = cg.customer_group_id
       INNER JOIN `catalog_product_entity_int` AS `tad_status`
               ON tad_status.entity_id = e.entity_id
                  AND tad_status.attribute_id = 89
                  AND tad_status.store_id = 0
       LEFT JOIN `catalog_product_entity_int` AS `tas_status`
              ON tas_status.entity_id = e.entity_id
                 AND tas_status.attribute_id = 89
                 AND tas_status.store_id = cs.store_id
       LEFT JOIN `catalog_product_entity_int` AS `tad_tax_class_id`
              ON tad_tax_class_id.entity_id = e.entity_id
                 AND tad_tax_class_id.attribute_id = 115
                 AND tad_tax_class_id.store_id = 0
       LEFT JOIN `catalog_product_entity_int` AS `tas_tax_class_id`
              ON tas_tax_class_id.entity_id = e.entity_id
                 AND tas_tax_class_id.attribute_id = 115
                 AND tas_tax_class_id.store_id = cs.store_id
       LEFT JOIN `catalog_product_entity_decimal` AS `tad_price`
              ON tad_price.entity_id = e.entity_id
                 AND tad_price.attribute_id = 69
                 AND tad_price.store_id = 0
       LEFT JOIN `catalog_product_entity_decimal` AS `tas_price`
              ON tas_price.entity_id = e.entity_id
                 AND tas_price.attribute_id = 69
                 AND tas_price.store_id = cs.store_id
       LEFT JOIN `catalog_product_entity_decimal` AS `tad_special_price`
              ON tad_special_price.entity_id = e.entity_id
                 AND tad_special_price.attribute_id = 70
                 AND tad_special_price.store_id = 0
       LEFT JOIN `catalog_product_entity_decimal` AS `tas_special_price`
              ON tas_special_price.entity_id = e.entity_id
                 AND tas_special_price.attribute_id = 70
                 AND tas_special_price.store_id = cs.store_id
       LEFT JOIN `catalog_product_entity_datetime` AS `tad_special_from_date`
              ON tad_special_from_date.entity_id = e.entity_id
                 AND tad_special_from_date.attribute_id = 71
                 AND tad_special_from_date.store_id = 0
       LEFT JOIN `catalog_product_entity_datetime` AS `tas_special_from_date`
              ON tas_special_from_date.entity_id = e.entity_id
                 AND tas_special_from_date.attribute_id = 71
                 AND tas_special_from_date.store_id = cs.store_id
       LEFT JOIN `catalog_product_entity_datetime` AS `tad_special_to_date`
              ON tad_special_to_date.entity_id = e.entity_id
                 AND tad_special_to_date.attribute_id = 72
                 AND tad_special_to_date.store_id = 0
       LEFT JOIN `catalog_product_entity_datetime` AS `tas_special_to_date`
              ON tas_special_to_date.entity_id = e.entity_id
                 AND tas_special_to_date.attribute_id = 72
                 AND tas_special_to_date.store_id = cs.store_id
       INNER JOIN `cataloginventory_stock_status` AS `ciss`
               ON ciss.product_id = e.entity_id
                  AND ciss.website_id = cw.website_id
WHERE  ( e.type_id = 'simple' )
       AND ( IF(IFNULL(tas_status.value_id, -1) > 0, tas_status.value, tad_status.value) = 1 )
       AND ( ciss.stock_status = 1 );
