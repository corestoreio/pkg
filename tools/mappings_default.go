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

package tools

// JSONMapEntityTypes provides a mapping for Mage1+2. A developer has the option to provide a custom map.
var JSONMapEntityTypes = []byte(`{
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
}`)

// EavAttributeModelMap contains default mappings for Mage1+2. A developer has the option to provide a custom map.
var EavAttributeModelMap = columnMap{
	EavAttributeFrontendModel: []byte(`{
        "catalog\/product_attribute_frontend_image": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "frontend_model": "catalog.Product().Attribute().Frontend().Image()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Frontend\\Image": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "frontend_model": "catalog.Product().Attribute().Frontend().Image()"
        },
        "eav\/entity_attribute_frontend_datetime": {
            "import_path": "github.com\/corestoreio\/csfw\/eav",
            "frontend_model": "eav.Entity().Attribute().Frontend().Datetime()"
        },
        "Magento\\Eav\\Model\\Entity\\Attribute\\Frontend\\Datetime": {
            "import_path": "github.com\/corestoreio\/csfw\/eav",
            "frontend_model": "eav.Entity().Attribute().Frontend().Datetime()"
        }
    }`),
	EavAttributeBackendModel: []byte(`{
        "catalog\/attribute_backend_customlayoutupdate": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Attribute().Backend().Customlayoutupdate()"
        },
        "Magento\\Catalog\\Model\\Attribute\\Backend\\Customlayoutupdate": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Attribute().Backend().Customlayoutupdate()"
        },
        "catalog\/category_attribute_backend_image": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Category().Attribute().Backend().Image()"
        },
        "Magento\\Catalog\\Model\\Category\\Attribute\\Backend\\Image": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Category().Attribute().Backend().Image()"
        },
        "catalog\/category_attribute_backend_sortby": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Category().Attribute().Backend().Sortby()"
        },
        "Magento\\Catalog\\Model\\Category\\Attribute\\Backend\\Sortby": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Category().Attribute().Backend().Sortby()"
        },
        "catalog\/category_attribute_backend_urlkey": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": ""
        },
        "catalog\/product_attribute_backend_boolean": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Boolean()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Boolean": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Boolean()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Category": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Category()"
        },
        "catalog\/product_attribute_backend_groupprice": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Groupprice()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\GroupPrice": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().GroupPrice()"
        },
        "catalog\/product_attribute_backend_media": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Media()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Media": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Media()"
        },
        "catalog\/product_attribute_backend_msrp": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Price()"
        },
        "catalog\/product_attribute_backend_price": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Price()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Price": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Price()"
        },
        "catalog\/product_attribute_backend_recurring": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Recurring()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Recurring": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Recurring()"
        },
        "catalog\/product_attribute_backend_sku": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Sku()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Sku": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Sku()"
        },
        "catalog\/product_attribute_backend_startdate": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Startdate()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Startdate": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Startdate()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Stock": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Stock()"
        },
        "catalog\/product_attribute_backend_tierprice": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Tierprice()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Tierprice": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Tierprice()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Weight": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": "catalog.Product().Attribute().Backend().Weight()"
        },
        "catalog\/product_attribute_backend_urlkey": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "backend_model": ""
        },
        "customer\/attribute_backend_data_boolean": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Attribute().Backend().Data().Boolean()"
        },
        "Magento\\Customer\\Model\\Attribute\\Backend\\Data\\Boolean": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Attribute().Backend().Data().Boolean()"
        },
        "customer\/customer_attribute_backend_billing": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Attribute().Backend().Billing()"
        },
        "Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Billing": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Attribute().Backend().Billing()"
        },
        "customer\/customer_attribute_backend_password": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Attribute().Backend().Password()"
        },
        "Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Password": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Attribute().Backend().Password()"
        },
        "customer\/customer_attribute_backend_shipping": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Attribute().Backend().Shipping()"
        },
        "Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Shipping": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Attribute().Backend().Shipping()"
        },
        "customer\/customer_attribute_backend_store": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Attribute().Backend().Store()"
        },
        "Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Store": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Attribute().Backend().Store()"
        },
        "customer\/customer_attribute_backend_website": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Attribute().Backend().Website()"
        },
        "Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Website": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Attribute().Backend().Website()"
        },
        "customer\/entity_address_attribute_backend_region": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Address().Attribute().Backend().Region()"
        },
        "Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Backend\\Region": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Address().Attribute().Backend().Region()"
        },
        "customer\/entity_address_attribute_backend_street": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "backend_model": "customer.Address().Attribute().Backend().Street()"
        },
        "Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\DefaultBackend": {
            "import_path": "github.com\/corestoreio\/csfw\/eav",
            "backend_model": "eav.Attribute().Backend().DefaultBackend()"
        },
        "eav\/entity_attribute_backend_datetime": {
            "import_path": "github.com\/corestoreio\/csfw\/eav",
            "backend_model": "eav.Attribute().Backend().Datetime()"
        },
        "Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Datetime": {
            "import_path": "github.com\/corestoreio\/csfw\/eav",
            "backend_model": "eav.Attribute().Backend().Datetime()"
        },
        "eav\/entity_attribute_backend_time_created": {
            "import_path": "github.com\/corestoreio\/csfw\/eav",
            "backend_model": "eav.Entity().Attribute().Backend().Time().Created()"
        },
        "Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Time\\Created": {
            "import_path": "github.com\/corestoreio\/csfw\/eav",
            "backend_model": "eav.Attribute().Backend().Time().Created()"
        },
        "eav\/entity_attribute_backend_time_updated": {
            "import_path": "github.com\/corestoreio\/csfw\/eav",
            "backend_model": "eav.Attribute().Backend().Time().Updated()"
        },
        "Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Time\\Updated": {
            "import_path": "github.com\/corestoreio\/csfw\/eav",
            "backend_model": "eav.Attribute().Backend().Time().Updated()"
        }
    }`),
	EavAttributeSourceModel: []byte(`{
        "bundle\/product_attribute_source_price_view": {
            "import_path": "github.com\/corestoreio\/csfw\/bundle",
            "source_model": "bundle.Product().Attribute().Source().Price().View()"
        },
        "Magento\\Bundle\\Model\\Product\\Attribute\\Source\\Price\\View": {
            "import_path": "github.com\/corestoreio\/csfw\/bundle",
            "source_model": "bundle.Product().Attribute().Source().Price().View()"
        },
        "Magento\\CatalogInventory\\Model\\Source\\Stock": {
            "import_path": "github.com\/corestoreio\/csfw\/cataloginventory",
            "source_model": "cataloginventory.Source().Stock()"
        },
        "catalog\/category_attribute_source_layout": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": ""
        },
        "Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Layout": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": ""
        },
        "catalog\/category_attribute_source_mode": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Category().Attribute().Source().Mode()"
        },
        "Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Mode": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Category().Attribute().Source().Mode()"
        },
        "catalog\/category_attribute_source_page": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Category().Attribute().Source().Page()"
        },
        "Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Page": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Category().Attribute().Source().Page()"
        },
        "catalog\/category_attribute_source_sortby": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Category().Attribute().Source().Sortby()"
        },
        "Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Sortby": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Category().Attribute().Source().Sortby()"
        },
        "catalog\/entity_product_attribute_design_options_container": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Product().Attribute().Design().Options().Container()"
        },
        "Magento\\Catalog\\Model\\Entity\\Product\\Attribute\\Design\\Options\\Container": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Product().Attribute().Design().Options().Container()"
        },
        "catalog\/product_attribute_source_countryofmanufacture": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Product().Attribute().Source().Countryofmanufacture()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Countryofmanufacture": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Product().Attribute().Source().Countryofmanufacture()"
        },
        "catalog\/product_attribute_source_layout": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Product().Attribute().Source().Layout()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Layout": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Product().Attribute().Source().Layout()"
        },
        "Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Status": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Product().Attribute().Source().Status()"
        },
        "catalog\/product_attribute_source_msrp_type_enabled": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Product().Attribute().Source().Msrp().Type().Enabled()"
        },
        "catalog\/product_attribute_source_msrp_type_price": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Product().Attribute().Source().Msrp().Type().Price()"
        },
        "catalog\/product_status": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Product().Status()"
        },
        "catalog\/product_visibility": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Product().Visibility()"
        },
        "Magento\\Catalog\\Model\\Product\\Visibility": {
            "import_path": "github.com\/corestoreio\/csfw\/catalog",
            "source_model": "catalog.Product().Visibility()"
        },
        "core\/design_source_design": {
            "import_path": "github.com\/corestoreio\/csfw\/core",
            "source_model": "core.Design().Source().Design()"
        },
        "customer\/customer_attribute_source_group": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "source_model": "customer.Attribute().Source().Group()"
        },
        "Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Group": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "source_model": "customer.Attribute().Source().Group()"
        },
        "customer\/customer_attribute_source_store": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "source_model": "customer.Attribute().Source().Store()"
        },
        "Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Store": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "source_model": "customer.Attribute().Source().Store()"
        },
        "customer\/customer_attribute_source_website": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "source_model": "customer.Attribute().Source().Website()"
        },
        "Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Website": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "source_model": "customer.Attribute().Source().Website()"
        },
        "customer\/entity_address_attribute_source_country": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "source_model": "customer.Address().Attribute().Source().Country()"
        },
        "Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Source\\Country": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "source_model": "customer.Address().Attribute().Source().Country()"
        },
        "customer\/entity_address_attribute_source_region": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "source_model": "customer.Address().Attribute().Source().Region()"
        },
        "Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Source\\Region": {
            "import_path": "github.com\/corestoreio\/csfw\/customer",
            "source_model": "customer.Address().Attribute().Source().Region()"
        },
        "eav\/entity_attribute_source_boolean": {
            "import_path": "github.com\/corestoreio\/csfw\/eav",
            "source_model": "eav.Attribute().Source().Boolean()"
        },
        "Magento\\Eav\\Model\\Entity\\Attribute\\Source\\Boolean": {
            "import_path": "github.com\/corestoreio\/csfw\/eav",
            "source_model": "eav.Attribute().Source().Boolean()"
        },
        "eav\/entity_attribute_source_table": {
            "import_path": "github.com\/corestoreio\/csfw\/eav",
            "source_model": "eav.Attribute().Source().Table()"
        },
        "Magento\\Eav\\Model\\Entity\\Attribute\\Source\\Table": {
            "import_path": "github.com\/corestoreio\/csfw\/eav",
            "source_model": "eav.Attribute().Source().Table()"
        },
        "Magento\\Msrp\\Model\\Product\\Attribute\\Source\\Type\\Price": {
            "import_path": "github.com\/corestoreio\/csfw\/msrp",
            "source_model": "msrp.Product().Attribute().Source().Type().Price()"
        },
        "tax\/class_source_product": {
            "import_path": "github.com\/corestoreio\/csfw\/tax",
            "source_model": "tax.Class().Source().Product()"
        },
        "Magento\\Tax\\Model\\TaxClass\\Source\\Product": {
            "import_path": "github.com\/corestoreio\/csfw\/tax",
            "source_model": "tax.TaxClass().Source().Product()"
        },
        "Magento\\Theme\\Model\\Theme\\Source\\Theme": {
            "import_path": "github.com\/corestoreio\/csfw\/theme",
            "source_model": ""
        }
    }`),
}
