INSERT INTO `sales_bestsellers_aggregated_monthly`
(
  `period`,
  `store_id`,
  `product_id`,
  `product_name`,
  `product_price`,
  `qty_ordered`,
  `rating_pos`
)
  SELECT
    Date_format(t.period, '%Y-%m-01') AS `period`,
    `t`.`store_id`,
    `t`.`product_id`,
    `t`.`product_name`,
    `t`.`product_price`,
    `t`.`qty_ordered`,
    `t`.`rating_pos`
  FROM (
         SELECT
           `t`.`period`,
           `t`.`store_id`,
           `t`.`product_id`,
           `t`.`product_name`,
           `t`.`product_price`,
           `t`.`total_qty`                                                                AS `qty_ordered`,
           (@pos := IF(t.`store_id` <> @prevstoreid
                       OR Date_format(t.period, '%Y-%m-01') <> @prevperiod, 1, @pos + 1)) AS `rating_pos`,
           (@prevstoreid := t.`store_id`)                                                 AS `prevstoreid`,
           (@prevperiod := Date_format(t.period, '%Y-%m-01'))                             AS `prevperiod`
         FROM
           (SELECT
              `t`.`period`,
              `t`.`store_id`,
              `t`.`product_id`,
              `t`.`product_name`,
              `t`.`product_price`,
              Sum(t.qty_ordered) AS `total_qty`
            FROM `sales_bestsellers_aggregated_daily` AS `t`
            GROUP BY `t`.`store_id`,
              Date_format(t.period, '%Y-%m-01'),
              `t`.`product_id`
            ORDER BY `t`.`store_id` ASC,
              Date_format(t.period, '%Y-%m-01'),
              `total_qty` DESC
           ) AS `t`
       ) AS `t`

ON DUPLICATE KEY
UPDATE `period`   = VALUES(`period`),
  `store_id`      = VALUES(`store_id`),
  `product_id`    = VALUES(`product_id`),
  `product_name`  = VALUES(`product_name`),
  `product_price` = VALUES(`product_price`),
  `qty_ordered`   = VALUES(`qty_ordered`),
  `rating_pos`    = VALUES(`rating_pos`);
