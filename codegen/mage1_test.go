// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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
// limitations under the License.

// +build mage1

package codegen

func getDataGetEavValueTables() []dataGetEavValueTables {
	return []dataGetEavValueTables{
		{
			haveETC: []string{"catalog_category", "catalog_product"},
			wantErr: false,
			wantCVMap: TypeCodeValueTable{"catalog_category": map[string]string{"catalog_category_entity_datetime": "datetime", "catalog_category_entity_decimal": "decimal", "catalog_category_entity_int": "int", "catalog_category_entity_text": "text", "catalog_category_entity_varchar": "varchar"},
				"catalog_product": map[string]string{"catalog_product_entity_datetime": "datetime", "catalog_product_entity_decimal": "decimal", "catalog_product_entity_int": "int", "catalog_product_entity_text": "text", "catalog_product_entity_varchar": "varchar"}},
		},
		{
			haveETC: []string{"customer_address", "customer"},
			wantErr: false,
			wantCVMap: TypeCodeValueTable{"customer_address": map[string]string{"customer_address_entity_text": "text", "customer_address_entity_varchar": "varchar", "customer_address_entity_datetime": "datetime", "customer_address_entity_decimal": "decimal", "customer_address_entity_int": "int"},
				"customer": map[string]string{"csCustomer_value_decimal": "decimal", "csCustomer_value_int": "int", "csCustomer_value_text": "text", "csCustomer_value_varchar": "varchar", "csCustomer_value_datetime": "datetime"}},
		},
		{
			haveETC:   []string{"catalog_address"},
			wantErr:   false,
			wantCVMap: TypeCodeValueTable{"catalog_address": map[string]string{}},
		},
	}
}
