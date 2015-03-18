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
package tools

var JSONMappingEntityTypes = []byte(`{
  "customer": {
    "import_path": "github.com/corestoreio/csfw/customer",
    "entity_model": "customer.Customer()",
    "attribute_model": "customer.Attribute()",
    "entity_table": "customer.Customer()",
    "increment_model": "customer.Customer()",
    "additional_attribute_table": "customer.Customer()",
    "entity_attribute_collection": "customer.Customer()"
  },
  "customer_address": {
    "import_path": "github.com/corestoreio/csfw/customer",
    "entity_model": "customer.Address()",
    "attribute_model": "customer.AddressAttribute()",
    "entity_table": "customer.Address()",
    "additional_attribute_table": "customer.Address()",
    "entity_attribute_collection": "customer.Address()"
  },
  "catalog_category": {
    "import_path": "github.com/corestoreio/csfw/catalog",
    "entity_model": "catalog.Category()",
    "attribute_model": "catalog.Attribute()",
    "entity_table": "catalog.Category()",
    "additional_attribute_table": "catalog.Category()",
    "entity_attribute_collection": "catalog.Category()"
  },
  "catalog_product": {
    "import_path": "github.com/corestoreio/csfw/catalog",
    "entity_model": "catalog.Product()",
    "attribute_model": "catalog.Attribute()",
    "entity_table": "catalog.Product()",
    "additional_attribute_table": "catalog.Product()",
    "entity_attribute_collection": "catalog.Product()"
  }
}
`)

var JSONMappingEAVAttributeModels = []byte(`{
{
    "eav/entity_attribute_backend_datetime": {
        "import_path": "github.com/corestoreio/csfw/eav",
        "backend_model": "eav.Attribute().Backend().Datetime()",
    },
    "catalog/product_attribute_backend_price": {
        "import_path": "github.com/corestoreio/csfw/catalog",
        "backend_model": "catalog.Attribute().Backend().Price()",
    }
}
`)
