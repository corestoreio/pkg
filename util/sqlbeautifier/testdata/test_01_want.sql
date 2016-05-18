select 
	count(distinct 
	e.entity_id 
) 
 from 
	catalog_product_flat_1 as e join 
	catalog_category_product_index as cat_index on cat_index.product_id = e.entity_id and cat_index.store_id = '1' and cat_index.visibility in (2, 4) and cat_index.category_id = '2'
 join 
	catalog_product_index_price as price_index on price_index.entity_id = e.entity_id and price_index.website_id = '1' and price_index.customer_group_id = 0

 where 
		(
		(
		(ifnull(
	sku, 
	0 
) in ('WS12', 'WT09', 'MT07', 'MH07', '24-MB02', '24-WB04', '241-MB08', '240-LV05'))))

