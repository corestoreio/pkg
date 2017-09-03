
SELECT     Date(Date_add(`o`.`created_at`, INTERVAL 3600 second)) AS `period`,
           `o`.`store_id`,
           `o`.`status` AS `order_status`,
           Count(o.entity_id) AS `orders_count`,
           Sum(oi.total_qty_ordered) AS `total_qty_ordered`,
           Sum(oi.total_qty_invoiced) AS `total_qty_invoiced`,
           Sum((Ifnull(o.base_grand_total, 0)    - Ifnull(o.base_total_canceled, 0)) * Ifnull(o.base_to_global_rate, 0)) AS `total_income_amount`,
           Sum((Ifnull(o.base_total_invoiced, 0) - Ifnull(o.base_tax_invoiced, 0) - Ifnull(o.base_shipping_invoiced, 0) - (Ifnull(o.base_total_refunded, 0) - Ifnull(o.base_tax_refunded, 0) - Ifnull(o.base_shipping_refunded, 0))) * Ifnull(o.base_to_global_rate, 0)) AS `total_revenue_amount`,
           Sum(((Ifnull(o.base_total_paid, 0) - Ifnull(o.base_total_refunded, 0)) - (Ifnull(o.base_tax_invoiced, 0) - Ifnull(o.base_tax_refunded, 0)) - (Ifnull(o.base_shipping_invoiced, 0) - Ifnull(o.base_shipping_refunded, 0)) - Ifnull(o.base_total_invoiced_cost, 0)) * Ifnull(o.base_to_global_rate, 0)) AS `total_profit_amount`,
           Sum(Ifnull(o.base_total_invoiced, 0) * Ifnull(o.base_to_global_rate, 0)) AS `total_invoiced_amount`,
           Sum(Ifnull(o.base_total_canceled, 0) * Ifnull(o.base_to_global_rate, 0)) AS `total_canceled_amount`,
           Sum(Ifnull(o.base_total_paid, 0)     * Ifnull(o.base_to_global_rate, 0)) AS `total_paid_amount`,
           Sum(Ifnull(o.base_total_refunded, 0) * Ifnull(o.base_to_global_rate, 0)) AS `total_refunded_amount`,
           Sum((Ifnull(o.base_tax_amount, 0)           - Ifnull(o.base_tax_canceled, 0)) * Ifnull(o.base_to_global_rate, 0)) AS `total_tax_amount`,
           Sum((Ifnull(o.base_tax_invoiced, 0)         -Ifnull(o.base_tax_refunded, 0)) * Ifnull(o.base_to_global_rate, 0)) AS `total_tax_amount_actual`,
           Sum((Ifnull(o.base_shipping_amount, 0)      - Ifnull(o.base_shipping_canceled, 0)) * Ifnull(o.base_to_global_rate, 0)) AS `total_shipping_amount`,
           Sum((Ifnull(o.base_shipping_invoiced, 0)    - Ifnull(o.base_shipping_refunded, 0)) * Ifnull(o.base_to_global_rate, 0)) AS `total_shipping_amount_actual`,
           Sum((Abs(Ifnull(o.base_discount_amount, 0)) - Ifnull(o.base_discount_canceled, 0)) * Ifnull(o.base_to_global_rate, 0)) AS `total_discount_amount`,
           Sum((Ifnull(o.base_discount_invoiced, 0)    - Ifnull(o.base_discount_refunded, 0)) * Ifnull(o.base_to_global_rate, 0)) AS `total_discount_amount_actual`
FROM       `sales_flat_order` AS `o`
INNER JOIN
           (
                    SELECT   `sales_flat_order_item`.`order_id`,
                             Sum(qty_ordered - Ifnull(qty_canceled, 0)) AS `total_qty_ordered`,
                             Sum(qty_invoiced) AS `total_qty_invoiced`
                    FROM     `sales_flat_order_item`
                    WHERE    (
                                      parent_item_id IS NULL)
                    GROUP BY `order_id`) AS `oi`
ON         oi.order_id = o.entity_id
WHERE      (
                      o.state NOT IN ('pending_payment',
                                      'new'))
GROUP BY   Date(Date_add(`o`.`created_at`, INTERVAL 3600 second)), `o`.`store_id`, `o`.`status`
HAVING     (period LIKE '2016-11-20' OR period LIKE '2016-11-21')

on duplicate KEY
UPDATE `period` = VALUES
       (
              `period`
       )
       ,
       `store_id` = VALUES
       (
              `store_id`
       )
       ,
       `order_status` = VALUES
       (
              `order_status`
       )
       ,
       `orders_count` = VALUES
       (
              `orders_count`
       )
       ,
       `total_qty_ordered` = VALUES
       (
              `total_qty_ordered`
       )
       ,
       `total_qty_invoiced` = VALUES
       (
              `total_qty_invoiced`
       )
       ,
       `total_income_amount` = VALUES
       (
              `total_income_amount`
       )
       ,
       `total_revenue_amount` = VALUES
       (
              `total_revenue_amount`
       )
       ,
       `total_profit_amount` = VALUES
       (
              `total_profit_amount`
       )
       ,
       `total_invoiced_amount` = VALUES
       (
              `total_invoiced_amount`
       )
       ,
       `total_canceled_amount` = VALUES
       (
              `total_canceled_amount`
       )
       ,
       `total_paid_amount` = VALUES
       (
              `total_paid_amount`
       )
       ,
       `total_refunded_amount` = VALUES
       (
              `total_refunded_amount`
       )
       ,
       `total_tax_amount` = VALUES
       (
              `total_tax_amount`
       )
       ,
       `total_tax_amount_actual` = VALUES
       (
              `total_tax_amount_actual`
       )
       ,
       `total_shipping_amount` = VALUES
       (
              `total_shipping_amount`
       )
       ,
       `total_shipping_amount_actual` = VALUES
       (
              `total_shipping_amount_actual`
       )
       ,
       `total_discount_amount` = VALUES
       (
              `total_discount_amount`
       )
       ,
       `total_discount_amount_actual` = VALUES
       (
              `total_discount_amount_actual`
       );
