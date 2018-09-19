/*
// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License. */

-- This CTE makes several columns and computations in `catalog_category_entity`
-- superfluous, also the attribute `url_path`.

WITH RECURSIVE categoryTreeCTE AS
(
	SELECT
		entity_id
		,parent_id
		,0 as `level`
		,`cat_name` AS name_path
		,`url_key` AS url_key_path
		, position
	FROM categoryFlatCTE WHERE parent_id = 1
	UNION ALL
	SELECT
		cat1.entity_id
		,cat1.parent_id
		,categoryTreeCTE.level+1 as `level`
		,concat(categoryTreeCTE.name_path,' -> ',cat1.cat_name) as name_path
		,concat(categoryTreeCTE.url_key_path,'/',cat1.url_key) as url_key_path
		,cat1.position
	FROM categoryFlatCTE cat1
	INNER JOIN categoryTreeCTE ON categoryTreeCTE.entity_id = cat1.parent_id
),
categoryFlatCTE AS (
	SELECT ce.entity_id,
		   ce.parent_id,
		   ce.position,
		   IF(cev45s1.value_id IS NULL,cev45s0.value,cev45s1.value) as `cat_name`,
		   IF(cev117s1.value_id IS NULL,IFNULL(cev117s0.value,''),cev117s1.value) as `url_key`
	FROM   catalog_category_entity ce
       LEFT JOIN catalog_category_entity_varchar cev45s0
              ON ce.entity_id = cev45s0.entity_id
                 AND cev45s0.attribute_id = 45
                 AND cev45s0.store_id = 0
       LEFT JOIN catalog_category_entity_varchar cev45s1
              ON ce.entity_id = cev45s1.entity_id
                 AND cev45s1.attribute_id = 45
                 AND cev45s1.store_id = 1

       LEFT JOIN catalog_category_entity_varchar cev117s0
              ON ce.entity_id = cev117s0.entity_id
                 AND cev117s0.attribute_id = 117
                 AND cev117s0.store_id = 0
       LEFT JOIN catalog_category_entity_varchar cev117s1
              ON ce.entity_id = cev117s1.entity_id
                 AND cev117s1.attribute_id = 117
                 AND cev117s1.store_id = 1

)
select * from categoryTreeCTE order by parent_id,position;
