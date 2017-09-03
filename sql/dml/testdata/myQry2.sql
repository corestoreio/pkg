INSERT INTO `catalog_product_index_eav_idx`
  SELECT
    `pid`.`entity_id`,
    `pid`.`attribute_id`,
    `pid`.`store_id`,
    Ifnull(pis.value, pid.value) AS `value`
  FROM (SELECT
          `s`.`store_id`,
          `s`.`website_id`,
          `d`.`entity_id`,
          `d`.`attribute_id`,
          `d`.`value`
        FROM `core_store` AS `s`
          LEFT JOIN `catalog_product_entity_int` AS `d`
            ON 1 = 1
               AND d.store_id = 0
          INNER JOIN `catalog_product_entity_int` AS `tad_status`
            ON tad_status.entity_id = d.entity_id
               AND tad_status.attribute_id = 89
               AND tad_status.store_id = 0
          LEFT JOIN `catalog_product_entity_int` AS `tas_status`
            ON tas_status.entity_id = d.entity_id
               AND tas_status.attribute_id = 89
               AND tas_status.store_id = s.store_id
        WHERE (s.store_id != 0)
              AND (IF(Ifnull(tas_status.value_id, -1) > 0, tas_status.value,
                      tad_status.value) =
                   1)) AS `pid`
    LEFT JOIN `catalog_product_entity_int` AS `pis`
      ON pis.entity_id = pid.entity_id
         AND pis.attribute_id = pid.attribute_id
         AND pis.store_id = pid.store_id
    INNER JOIN `cataloginventory_stock_status` AS `ciss`
      ON ciss.product_id = pid.entity_id
         AND ciss.website_id = pid.website_id
  WHERE (pid.attribute_id IN ('75', '115', '171', '186',
                              '187', '188', '210', '211'))
        AND (Ifnull(pis.value, pid.value) IS NOT NULL)
        AND (ciss.stock_status = 1);