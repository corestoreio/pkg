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

var JSONMapFrontendModel = []byte(`{
    "catalog\/product_attribute_frontend_image": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "frontend_model": "catalog.Product().Attribute().Frontend().Image()"
    },
    "eav\/entity_attribute_frontend_datetime": {
        "import_path": "github.com\/corestoreio\/csfw\/eav",
        "frontend_model": "eav.Entity().Attribute().Frontend().Datetime()"
    }
}`)

var JSONMapBackendModel = []byte(`{
    "catalog\/attribute_backend_customlayoutupdate": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Attribute().Backend().Customlayoutupdate()"
    },
    "catalog\/category_attribute_backend_image": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Category().Attribute().Backend().Image()"
    },
    "catalog\/category_attribute_backend_sortby": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Category().Attribute().Backend().Sortby()"
    },
    "catalog\/category_attribute_backend_urlkey": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Category().Attribute().Backend().Urlkey()"
    },
    "catalog\/product_attribute_backend_boolean": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Product().Attribute().Backend().Boolean()"
    },
    "catalog\/product_attribute_backend_groupprice": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Product().Attribute().Backend().Groupprice()"
    },
    "catalog\/product_attribute_backend_media": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Product().Attribute().Backend().Media()"
    },
    "catalog\/product_attribute_backend_msrp": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Product().Attribute().Backend().Msrp()"
    },
    "catalog\/product_attribute_backend_price": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Product().Attribute().Backend().Price()"
    },
    "catalog\/product_attribute_backend_recurring": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Product().Attribute().Backend().Recurring()"
    },
    "catalog\/product_attribute_backend_sku": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Product().Attribute().Backend().Sku()"
    },
    "catalog\/product_attribute_backend_startdate": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Product().Attribute().Backend().Startdate()"
    },
    "catalog\/product_attribute_backend_tierprice": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Product().Attribute().Backend().Tierprice()"
    },
    "catalog\/product_attribute_backend_urlkey": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "backend_model": "catalog.Product().Attribute().Backend().Urlkey()"
    },
    "customer\/attribute_backend_data_boolean": {
        "import_path": "github.com\/corestoreio\/csfw\/customer",
        "backend_model": "customer.Attribute().Backend().Data().Boolean()"
    },
    "customer\/customer_attribute_backend_billing": {
        "import_path": "github.com\/corestoreio\/csfw\/customer",
        "backend_model": "customer.Attribute().Backend().Billing()"
    },
    "customer\/customer_attribute_backend_password": {
        "import_path": "github.com\/corestoreio\/csfw\/customer",
        "backend_model": "customer.Attribute().Backend().Password()"
    },
    "customer\/customer_attribute_backend_shipping": {
        "import_path": "github.com\/corestoreio\/csfw\/customer",
        "backend_model": "customer.Attribute().Backend().Shipping()"
    },
    "customer\/customer_attribute_backend_store": {
        "import_path": "github.com\/corestoreio\/csfw\/customer",
        "backend_model": "customer.Attribute().Backend().Store()"
    },
    "customer\/customer_attribute_backend_website": {
        "import_path": "github.com\/corestoreio\/csfw\/customer",
        "backend_model": "customer.Attribute().Backend().Website()"
    },
    "customer\/entity_address_attribute_backend_region": {
        "import_path": "github.com\/corestoreio\/csfw\/customer",
        "backend_model": "customer.Entity().Address().Attribute().Backend().Region()"
    },
    "customer\/entity_address_attribute_backend_street": {
        "import_path": "github.com\/corestoreio\/csfw\/customer",
        "backend_model": "customer.Entity().Address().Attribute().Backend().Street()"
    },
    "eav\/entity_attribute_backend_datetime": {
        "import_path": "github.com\/corestoreio\/csfw\/eav",
        "backend_model": "eav.Entity().Attribute().Backend().Datetime()"
    },
    "eav\/entity_attribute_backend_time_created": {
        "import_path": "github.com\/corestoreio\/csfw\/eav",
        "backend_model": "eav.Entity().Attribute().Backend().Time().Created()"
    },
    "eav\/entity_attribute_backend_time_updated": {
        "import_path": "github.com\/corestoreio\/csfw\/eav",
        "backend_model": "eav.Entity().Attribute().Backend().Time().Updated()"
    }
}`)

var JSONMapSourceModel = []byte(`{
    "bundle\/product_attribute_source_price_view": {
        "import_path": "github.com\/corestoreio\/csfw\/bundle",
        "source_model": "bundle.Product().Attribute().Source().Price().View()"
    },
    "catalog\/category_attribute_source_layout": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "source_model": "catalog.Category().Attribute().Source().Layout()"
    },
    "catalog\/category_attribute_source_mode": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "source_model": "catalog.Category().Attribute().Source().Mode()"
    },
    "catalog\/category_attribute_source_page": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "source_model": "catalog.Category().Attribute().Source().Page()"
    },
    "catalog\/category_attribute_source_sortby": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "source_model": "catalog.Category().Attribute().Source().Sortby()"
    },
    "catalog\/entity_product_attribute_design_options_container": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "source_model": "catalog.Entity().Product().Attribute().Design().Options().Container()"
    },
    "catalog\/product_attribute_source_countryofmanufacture": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "source_model": "catalog.Product().Attribute().Source().Countryofmanufacture()"
    },
    "catalog\/product_attribute_source_layout": {
        "import_path": "github.com\/corestoreio\/csfw\/catalog",
        "source_model": "catalog.Product().Attribute().Source().Layout()"
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
    "core\/design_source_design": {
        "import_path": "github.com\/corestoreio\/csfw\/core",
        "source_model": "core.Design().Source().Design()"
    },
    "customer\/customer_attribute_source_group": {
        "import_path": "github.com\/corestoreio\/csfw\/customer",
        "source_model": "customer.Attribute().Source().Group()"
    },
    "customer\/customer_attribute_source_store": {
        "import_path": "github.com\/corestoreio\/csfw\/customer",
        "source_model": "customer.Attribute().Source().Store()"
    },
    "customer\/customer_attribute_source_website": {
        "import_path": "github.com\/corestoreio\/csfw\/customer",
        "source_model": "customer.Attribute().Source().Website()"
    },
    "customer\/entity_address_attribute_source_country": {
        "import_path": "github.com\/corestoreio\/csfw\/customer",
        "source_model": "customer.Entity().Address().Attribute().Source().Country()"
    },
    "customer\/entity_address_attribute_source_region": {
        "import_path": "github.com\/corestoreio\/csfw\/customer",
        "source_model": "customer.Entity().Address().Attribute().Source().Region()"
    },
    "eav\/entity_attribute_source_boolean": {
        "import_path": "github.com\/corestoreio\/csfw\/eav",
        "source_model": "eav.Entity().Attribute().Source().Boolean()"
    },
    "eav\/entity_attribute_source_table": {
        "import_path": "github.com\/corestoreio\/csfw\/eav",
        "source_model": "eav.Entity().Attribute().Source().Table()"
    },
    "tax\/class_source_product": {
        "import_path": "github.com\/corestoreio\/csfw\/tax",
        "source_model": "tax.Class().Source().Product()"
    }
}`)
