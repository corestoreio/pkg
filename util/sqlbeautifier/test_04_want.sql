select 
	t_d.entity_id, 
	t_d.attribute_id, 
	t_d.value as default_value, 
	t_s.value as store_value, 
	if(
	t_s.value_id is null, 
	t_d.value, 
	t_s.value 
) as value 
 from 
	catalog_product_entity_varchar as t_d left join 
	catalog_product_entity_varchar as t_s on t_s.attribute_id = t_d.attribute_id and t_s.entity_id = t_d.entity_id and t_s.store_id = 1

 where 
		(t_d.entity_id in (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)) and 
		(t_d.attribute_id in ('86', '70', '81', '83', '97', '101', '111', '115', '130')) and 
		(t_d.store_id = ifnull(
	t_s.store_id, 
	0 
))

 union all select 
	t_d.entity_id, 
	t_d.attribute_id, 
	t_d.value as default_value, 
	t_s.value as store_value, 
	if(
	t_s.value_id is null, 
	t_d.value, 
	t_s.value 
) as value 
 from 
	catalog_product_entity_decimal as t_d left join 
	catalog_product_entity_decimal as t_s on t_s.attribute_id = t_d.attribute_id and t_s.entity_id = t_d.entity_id and t_s.store_id = 1

 where 
		(t_d.entity_id in (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)) and 
		(t_d.attribute_id in ('74', '75', '78', '79', '117')) and 
		(t_d.store_id = ifnull(
	t_s.store_id, 
	0 
))

 union all select 
	t_d.entity_id, 
	t_d.attribute_id, 
	t_d.value as default_value, 
	t_s.value as store_value, 
	if(
	t_s.value_id is null, 
	t_d.value, 
	t_s.value 
) as value 
 from 
	catalog_product_entity_int as t_d left join 
	catalog_product_entity_int as t_s on t_s.attribute_id = t_d.attribute_id and t_s.entity_id = t_d.entity_id and t_s.store_id = 1

 where 
		(t_d.entity_id in (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)) and 
		(t_d.attribute_id in ('96', '94', '90', '129')) and 
		(t_d.store_id = ifnull(
	t_s.store_id, 
	0 
))

 union all select 
	t_d.entity_id, 
	t_d.attribute_id, 
	t_d.value as default_value, 
	t_s.value as store_value, 
	if(
	t_s.value_id is null, 
	t_d.value, 
	t_s.value 
) as value 
 from 
	catalog_product_entity_text as t_d left join 
	catalog_product_entity_text as t_s on t_s.attribute_id = t_d.attribute_id and t_s.entity_id = t_d.entity_id and t_s.store_id = 1

 where 
		(t_d.entity_id in (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)) and 
		(t_d.attribute_id in ('73', '82')) and 
		(t_d.store_id = ifnull(
	t_s.store_id, 
	0 
))

 union all select 
	t_d.entity_id, 
	t_d.attribute_id, 
	t_d.value as default_value, 
	t_s.value as store_value, 
	if(
	t_s.value_id is null, 
	t_d.value, 
	t_s.value 
) as value 
 from 
	catalog_product_entity_datetime as t_d left join 
	catalog_product_entity_datetime as t_s on t_s.attribute_id = t_d.attribute_id and t_s.entity_id = t_d.entity_id and t_s.store_id = 1

 where 
		(t_d.entity_id in (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)) and 
		(t_d.attribute_id in ('76', '77', '91', '92', '98', '99')) and 
		(t_d.store_id = ifnull(
	t_s.store_id, 
	0 
))

