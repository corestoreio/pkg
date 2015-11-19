select 
	main_table.* 
 from 
	store as main_table
 where 
		(main_table.website_id in ('1'))

 order by 
	case when main_table.store_id = 0 then 0 else 1 end asc, 
	main_table.sort_order asc, 
	main_table.name asc
