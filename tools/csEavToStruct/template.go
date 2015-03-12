// Copyright 2015 CoreStore Authors
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

// Generates code for all EAV types
package main

/*
   Data will be "carved in stone" because it only changes during development.
   - entity_type with translation of some columns to the Go type
   - attribute_set related tables: eav_attribute_set, eav_entity_attribute, eav_attribute_group, etc
   - label and option tables will not be hard coded
   - eav_attribute full config and from that one the flat table structure
*/

const tplEav = `// Copyright 2015 CoreStore Authors
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

// auto generated via csEavToStruct
package {{ .Package }}
import (
    "github.com/corestoreio/csfw/catalog"
    "github.com/corestoreio/csfw/customer"
)

var (
    CSEntityTypeCollection = CSEntityTypeSlice{
        {{ range .ETypeData }} &CSEntityType {
            EntityTypeId: {{ .EntityTypeId }},
            EntityTypeCode: "{{ .EntityTypeCode }}",
            EntityModel: {{ .EntityModel }},
            AttributeModel: {{ .AttributeModel }},
            EntityTable: {{ .EntityTable }},
            ValueTablePrefix: "{{ .ValueTablePrefix }}",
            IsDataSharing: {{ .IsDataSharing }},
            DataSharingKey: "{{ .DataSharingKey }}",
            DefaultAttributeSetId: {{ .DefaultAttributeSetId }},
            {{ if ne "" .IncrementModel }}IncrementModel: {{ .IncrementModel }},{{ end }}
            IncrementPerStore: {{ .IncrementPerStore }},
            IncrementPadLength: {{ .IncrementPadLength }},
            IncrementPadChar: "{{ .IncrementPadChar }}",
            AdditionalAttributeTable: {{ .AdditionalAttributeTable }},
            EntityAttributeCollection: {{ .EntityAttributeCollection }},
        },
        {{ end }}
    }
)
`

var defaultMapping = []byte(`{
  "customer": {
    "entity_model": "customer.Customer()",
    "attribute_model": "customer.Attribute()",
    "entity_table": "customer.EntityTable",
    "increment_model": "customer.Increment()",
    "additional_attribute_table": "customer.EavAttributeTable",
    "entity_attribute_collection": "customer.AttributeCollection()"
  },
  "customer_address": {
    "entity_model": "customer.Address()",
    "attribute_model": "customer.AddressAttribute()",
    "entity_table": "customer.EntityAddressTable",
    "additional_attribute_table": "customer.EavAttributeTable",
    "entity_attribute_collection": "customer.AddressAttributeCollection()"
  },
  "catalog_category": {
    "entity_model": "catalog.Category()",
    "attribute_model": "catalog.Attribute()",
    "entity_table": "catalog.EntityCategoryTable",
    "additional_attribute_table": "catalog.EavAttributeTable",
    "entity_attribute_collection": "catalog.CategoryAttributeCollection()"
  },
  "catalog_product": {
    "entity_model": "catalog.Product()",
    "attribute_model": "catalog.Attribute()",
    "entity_table": "catalog.EntityProductTable",
    "additional_attribute_table": "catalog.EavAttributeTable",
    "entity_attribute_collection": "catalog.ProductAttributeCollection()"
  }
}`)
