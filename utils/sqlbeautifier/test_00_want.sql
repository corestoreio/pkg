select 
	main_table.attribute_id, 
	main_table.entity_type_id, 
	main_table.attribute_code, 
	main_table.attribute_model, 
	main_table.backend_model, 
	main_table.backend_type, 
	main_table.backend_table, 
	main_table.frontend_model, 
	main_table.frontend_input, 
	main_table.frontend_label, 
	main_table.frontend_class, 
	main_table.source_model, 
	main_table.is_user_defined, 
	main_table.is_unique, 
	main_table.note, 
	additional_table.input_filter, 
	additional_table.validate_rules, 
	additional_table.is_system, 
	additional_table.sort_order, 
	additional_table.data_model, 
	ifnull(
	scope_table.is_visible, 
	additional_table.is_visible 
) as is_visible, 
	ifnull(
	scope_table.is_required, 
	main_table.is_required 
) as is_required, 
	ifnull(
	scope_table.default_value, 
	main_table.default_value 
) as default_value, 
	ifnull(
	scope_table.multiline_count, 
	additional_table.multiline_count 
) as multiline_count 
 from 
	eav_attribute as main_table join 
	customer_eav_attribute as additional_table on 
		(additional_table.attribute_id = main_table.attribute_id) and 
		(main_table.entity_type_id = :v1)
 left join 
	customer_eav_attribute_website as scope_table on 
		(scope_table.attribute_id = main_table.attribute_id) and 
		(scope_table.website_id = :v2)

 where multiline_count > 0

 order by 
	main_table.attribute_id asc
 limit 12, 13