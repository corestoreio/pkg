select 
	max(
	date_format(
	period, 
	'%Y-%m-%d' 
) 
) as period, 
	sum(
	qty_ordered 
) as qty_ordered, 
	sales_bestsellers_aggregated_yearly.product_id, 
	max(
	product_name 
) as product_name, 
	max(
	product_price 
) as product_price 
 from 
	sales_bestsellers_aggregated_yearly
 where 
		(sales_bestsellers_aggregated_yearly.product_id is not null) and 
		(store_id in (0)) and 
		(store_id in (0))
 group by product_id

 order by 
	qty_ordered desc
 limit 5