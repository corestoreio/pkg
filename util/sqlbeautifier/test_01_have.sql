SELECT COUNT(DISTINCT e.entity_id) FROM `catalog_product_flat_1` AS `e`
 INNER JOIN `catalog_category_product_index` AS `cat_index` ON cat_index.product_id=e.entity_id AND cat_index.store_id='1' AND cat_index.visibility IN(2, 4) AND cat_index.category_id='2'
 INNER JOIN `catalog_product_index_price` AS `price_index` ON price_index.entity_id = e.entity_id AND price_index.website_id = '1' AND price_index.customer_group_id = 0 WHERE (((IFNULL(`sku`, 0) IN ('WS12', 'WT09', 'MT07', 'MH07', '24-MB02', '24-WB04', '241-MB08', '240-LV05')) ))
