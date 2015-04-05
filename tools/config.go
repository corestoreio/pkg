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

import (
	"go/build"
	"os"
)

const PS = string(os.PathSeparator)
const CSImportPath string = "github.com" + PS + "corestoreio" + PS + "csfw"

type (
	// TableToStructMap uses a string key as easy identifier for maybe later manipulation in config_user.go
	// and a pointer to a TableToStruct struct. Developers can later use the init() func in config_user.go to change
	// the value of variable ConfigTableToStruct.
	TableToStructMap map[string]*TableToStruct
	TableToStruct    struct {
		// Package defines the name of the target package
		Package string
		// OutputFile specifies the full path where to write the newly generated code
		OutputFile string
		// QueryString SQL query to filter all the tables which you desire, e.g. SHOW TABLES LIKE 'catalog\_%'
		// This query must specify all tables you need for a package.
		SQLQuery string
		// EntityTypeCodes If provided then eav_entity_type.value_table_prefix will be evaluated for further tables.
		EntityTypeCodes []string
	}

	// AttributeToStructMap is a map which points to the AttributeToStruct configuration. Default entries can be
	// overriden by your configuration.
	AttributeToStructMap map[string]*AttributeToStruct
	// AttributeToStruct contains the configuration to materialize all attributes belonging to one EAV model
	AttributeToStruct struct {
		// AttrPkgImp defines the package import path to use: possible ATM: custattr and catattr and your custom
		// EAV package.
		AttrPkgImp string
		// FuncCollection specifies the name of the set attribute collection function name within the AttrPkgImp.
		FuncCollection string
		// FuncGetter specifies the name of the set attribute getter function name within the AttrPkgImp.
		FuncGetter string
		// Package defines the name of the target package, must be external. Default is {{.AttrPkgImp}}_test but
		// for your project you must provide your package name.
		Package string
		// MyStruct is the optional name of your struct from your package. MyStruct must embed a pointer to the
		// generated private attribute struct. This is useful if you want to override the method receivers.
		MyStruct string
		// OutputFile specifies the full path where to write the newly generated code
		OutputFile string
	}

	// EntityTypeMap uses a string key as easy identifier, which must also exists in the table eav_EntityTable,
	// for maybe later manipulation in config_user.go
	// and as value a pointer to a EntityType struct. Developers can later use the init() func in config_user.go to change
	// the value of variable ConfigEntityType.
	// The values of EntityType will be uses for materialization in Go code of the eav_entity_type table data.
	EntityTypeMap map[string]*EntityType
	// EntityType is configuration struct which maps the PHP classes to Go types, interfaces and table names.
	EntityType struct {
		// ImportPath path to the package
		ImportPath string
		// EntityModel Go type which implements eav.EntityTypeModeller
		EntityModel string
		// AttributeModel Go type which implements eav.EntityTypeAttributeModeller
		AttributeModel string
		// EntityTable Go type which implements eav.EntityTypeTabler
		EntityTable string
		// IncrementModel Go type which implements eav.EntityTypeIncrementModeller
		IncrementModel string
		// AdditionalAttributeTable Go type which implements eav.EntityTypeAdditionalAttributeTabler
		AdditionalAttributeTable string
		// EntityAttributeCollection Go type which implements eav.EntityAttributeCollectioner
		EntityAttributeCollection string

		// TempAdditionalAttributeTable string which defines the existing table name
		// and specifies more attribute configuration options besides eav_attribute table.
		// This table name is used in a DB query while materializing attribute configuration to Go code.
		// Mage_Eav_Model_Resource_Attribute_Collection::_initSelect()
		TempAdditionalAttributeTable string

		// TempAdditionalAttributeTableWebsite string which defines the existing table name
		// and stores website-dependent attribute parameters.
		// If an EAV model doesn't demand this functionality, let this string empty.
		// This table name is used in a DB query while materializing attribute configuration to Go code.
		// Mage_Customer_Model_Resource_Attribute::_getEavWebsiteTable()
		// Mage_Eav_Model_Resource_Attribute_Collection::_getEavWebsiteTable()
		TempAdditionalAttributeTableWebsite string
	}

	// AttributeModelDefMap contains data to map the three eav_attribute columns
	// (backend|frontend|source)_model to the correct Go function and package.
	// It contains mappings for Magento 1 & 2. A developer has the option to to change/extend the value
	// using the file config_user.go with the init() func.
	// Rethink the Go code here ... because catalog.Product().Attribute().Frontend().Image() is pretty long ... BUT
	// developers coming from Magento are already familiar with this code base and naming ...
	// Def for Definition to avoid a naming conflict :-( Better name?
	AttributeModelDefMap map[string]*AttributeModelDef
	// AttributeModelDef defines which Go type/func has which import path
	AttributeModelDef struct {
		// Is used when you want to specify your own package
		ImportPath string
		// GoModel is a function string which implements, when later executed, one of the n interfaces
		// for (backend|frontend|source|data)_model
		GoModel string
	}
)

// Keys returns all keys from a EntityTypeMap
func (m EntityTypeMap) Keys() []string {
	ret := make([]string, len(m), len(m))
	i := 0
	for k, _ := range m {
		ret[i] = k
		i++
	}
	return ret
}

var myPath = build.Default.GOPATH + PS + "src" + PS + CSImportPath + PS

// EavAttributeColumnNameToInterface mapping column name to Go interface name. Do not add attribute_model
// as this column is unused in Magento 1+2. If you have custom column then add it here.
var EavAttributeColumnNameToInterface = map[string]string{
	"backend_model":           "eav.AttributeBackendModeller",
	"frontend_model":          "eav.AttributeFrontendModeller",
	"source_model":            "eav.AttributeSourceModeller",
	"frontend_input_renderer": "eav.FrontendInputRendererIFace",
	"data_model":              "eav.AttributeDataModeller",
}

// TablePrefix defines the global table name prefix. See Magento install tool. Can be overridden via func init()
var TablePrefix string = ""

// TableMapMagento1To2 provides mapping between table names in tableToStruct. If a table name is in
// the map then the struct name will be rewritten to that new Magneto2 compatible table name.
// Do not change entries in this map except you can always append.
// @see Magento2 dev/tools/Magento/Tools/Migration/factory_table_names/replace_ce.php
var TableMapMagento1To2 = map[string]string{
	"core_cache":             "cache",     // not needed but added in case
	"core_cache_tag":         "cache_tag", // not needed but added in case
	"core_config_data":       "core_config_data",
	"core_design_change":     "design_change", // not needed but added in case
	"core_directory_storage": "media_storage_directory_storage",
	"core_email_template":    "email_template",
	"core_file_storage":      "media_storage_file_storage",
	"core_flag":              "flag",          // not needed but added in case
	"core_layout_link":       "layout_link",   // not needed but added in case
	"core_layout_update":     "layout_update", // not needed but added in case
	"core_resource":          "setup_module",  // not needed but added in case
	"core_session":           "session",       // not needed but added in case
	"core_store":             "store",
	"core_store_group":       "store_group",
	"core_variable":          "variable",
	"core_variable_value":    "variable_value",
	"core_website":           "store_website",
}

// ConfigTableToStruct contains default configuration. Use the file config_user.go with the func init() to change/extend it.
var ConfigTableToStruct = TableToStructMap{
	"eav": &TableToStruct{
		Package:         "eav",
		OutputFile:      myPath + "eav" + PS + "generated_tables.go",
		SQLQuery:        `SHOW TABLES LIKE "{{tableprefix}}eav%"`,
		EntityTypeCodes: nil,
	},
	"store": &TableToStruct{
		Package:    "store",
		OutputFile: myPath + "store" + PS + "generated_tables.go",
		SQLQuery: `SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE
		    TABLE_SCHEMA = DATABASE() AND
		    TABLE_NAME IN (
		    	'{{tableprefix}}core_store','{{tableprefix}}store',
		    	'{{tableprefix}}core_store_group','{{tableprefix}}store_group',
		    	'{{tableprefix}}core_website','{{tableprefix}}website'
		    ) GROUP BY TABLE_NAME;`,
		EntityTypeCodes: nil,
	},
	"catalog": &TableToStruct{
		Package:    "catalog",
		OutputFile: myPath + "catalog" + PS + "generated_tables.go",
		SQLQuery: `SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE
		    TABLE_SCHEMA = DATABASE() AND
		    (TABLE_NAME LIKE '{{tableprefix}}catalog\_%' OR TABLE_NAME LIKE '{{tableprefix}}catalogindex%' ) AND
		    TABLE_NAME NOT LIKE '{{tableprefix}}%bundle%' AND
		    TABLE_NAME NOT LIKE '{{tableprefix}}%\_flat\_%' GROUP BY TABLE_NAME;`,
		EntityTypeCodes: []string{"catalog_category", "catalog_product"},
	},
	"customer": &TableToStruct{
		Package:         "customer",
		OutputFile:      myPath + "customer" + PS + "generated_tables.go",
		SQLQuery:        `SHOW TABLES LIKE "{{tableprefix}}customer%"`,
		EntityTypeCodes: []string{"customer", "customer_address"},
	},
}

// ConfigMaterializationEntityType configuration for materializeEntityType() to write the materialized entity types
// into a folder. Other fields of the struct TableToStruct are ignored. Use the file config_user.go with the
// func init() to change/extend it.
var ConfigMaterializationEntityType = &TableToStruct{
	Package:    "eav_test",
	OutputFile: myPath + "eav" + PS + "generated_entity_type_test.go",
}

// ConfigMaterializationStore configuration for materializeStore() to write the materialized store data into a folder.
// For using this in your project you must modify the package name and output file path
var ConfigMaterializationStore = &TableToStruct{
	Package:    "store_test",
	OutputFile: myPath + "store" + PS + "generated_store_test.go",
}

// ConfigEntityType contains default configuration o materialize the entity types.
// Use the file config_user.go with the func init() to change/extend it.
// Needed in materializeEntityType()
var ConfigEntityType = EntityTypeMap{
	"customer": &EntityType{
		ImportPath:                          "github.com/corestoreio/csfw/customer",
		EntityModel:                         "customer.Customer()",
		AttributeModel:                      "customer.Attribute()",
		EntityTable:                         "customer.Customer()",
		IncrementModel:                      "customer.Customer()",
		AdditionalAttributeTable:            "customer.Customer()",
		EntityAttributeCollection:           "customer.Customer()",
		TempAdditionalAttributeTable:        "{{tableprefix}}customer_eav_attribute",
		TempAdditionalAttributeTableWebsite: "{{tableprefix}}customer_eav_attribute_website",
	},
	"customer_address": &EntityType{
		ImportPath:                          "github.com/corestoreio/csfw/customer",
		EntityModel:                         "customer.Address()",
		AttributeModel:                      "customer.AddressAttribute()",
		EntityTable:                         "customer.Address()",
		AdditionalAttributeTable:            "customer.Address()",
		EntityAttributeCollection:           "customer.Address()",
		TempAdditionalAttributeTable:        "{{tableprefix}}customer_eav_attribute",
		TempAdditionalAttributeTableWebsite: "{{tableprefix}}customer_eav_attribute_website",
	},
	"catalog_category": &EntityType{
		ImportPath:                   "github.com/corestoreio/csfw/catalog",
		EntityModel:                  "catalog.Category()",
		AttributeModel:               "catalog.Attribute()",
		EntityTable:                  "catalog.Category()",
		AdditionalAttributeTable:     "catalog.Category()",
		EntityAttributeCollection:    "catalog.Category()",
		TempAdditionalAttributeTable: "{{tableprefix}}catalog_eav_attribute",
	},
	"catalog_product": &EntityType{
		ImportPath:                   "github.com/corestoreio/csfw/catalog",
		EntityModel:                  "catalog.Product()",
		AttributeModel:               "catalog.Attribute()",
		EntityTable:                  "catalog.Product()",
		AdditionalAttributeTable:     "catalog.Product()",
		EntityAttributeCollection:    "catalog.Product()",
		TempAdditionalAttributeTable: "{{tableprefix}}catalog_eav_attribute",
	},
	// @todo extend for all sales entities
}

// ConfigMaterializationAttributes contains the configuration to materialize all attributes for the defined
// EAV entity types.
var ConfigMaterializationAttributes = AttributeToStructMap{
	"customer": &AttributeToStruct{
		AttrPkgImp:     "github.com/corestoreio/csfw/customer/custattr",
		FuncCollection: "SetCustomerCollection",
		FuncGetter:     "SetCustomerGetter",
		MyStruct:       "",
		Package:        "customer_test", // external package name
		OutputFile:     myPath + "customer" + PS + "generated_customer_attribute_test.go",
	},
	"customer_address": &AttributeToStruct{
		AttrPkgImp:     "github.com/corestoreio/csfw/customer/custattr",
		FuncCollection: "SetAddressCollection",
		FuncGetter:     "SetAddressGetter",
		MyStruct:       "",
		Package:        "customer_test",
		OutputFile:     myPath + "customer" + PS + "generated_address_attribute_test.go",
	},
	"catalog_product": &AttributeToStruct{
		AttrPkgImp:     "github.com/corestoreio/csfw/catalog/catattr",
		FuncCollection: "SetProductCollection",
		FuncGetter:     "SetProductGetter",
		MyStruct:       "",
		Package:        "catalog_test",
		OutputFile:     myPath + "catalog" + PS + "generated_product_attribute_test.go",
	},
	"catalog_category": &AttributeToStruct{
		AttrPkgImp:     "github.com/corestoreio/csfw/catalog/catattr",
		FuncCollection: "SetCategoryCollection",
		FuncGetter:     "SetCategoryGetter",
		MyStruct:       "",
		Package:        "catalog_test",
		OutputFile:     myPath + "catalog" + PS + "generated_category_attribute_test.go",
	},
	// extend here for other EAV attributes (not sales* types)
}

// ConfigAttributeModel contains default configuration. Use the file config_user.go with the func init() to change/extend it.
var ConfigAttributeModel = AttributeModelDefMap{
	"catalog/product_attribute_frontend_image": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductFrontendImage().Config(eav.AttributeFrontendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Frontend\\Image": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductFrontendImage().Config(eav.AttributeFrontendIdx({{.AttributeIndex}}))",
	},
	"eav/entity_attribute_frontend_datetime": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.AttributeFrontendDatetime().Config(eav.AttributeFrontendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Frontend\\Datetime": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.AttributeFrontendDatetime().Config(eav.AttributeFrontendIdx({{.AttributeIndex}}))",
	},
	"catalog/attribute_backend_customlayoutupdate": &AttributeModelDef{
		ImportPath: "",
		GoModel:    "",
	},
	"Magento\\Catalog\\Model\\Attribute\\Backend\\Customlayoutupdate": &AttributeModelDef{
		ImportPath: "",
		GoModel:    "",
	},
	"catalog/category_attribute_backend_image": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.CategoryBackendImage().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Backend\\Image": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.CategoryBackendImage().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"catalog/category_attribute_backend_sortby": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.CategoryBackendSortby().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Backend\\Sortby": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.CategoryBackendSortby().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"catalog/category_attribute_backend_urlkey": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "",
	},
	"catalog/product_attribute_backend_boolean": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendBoolean().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Boolean": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendBoolean().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Category": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendCategory().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"catalog/product_attribute_backend_groupprice": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendGroupPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\GroupPrice": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendGroupPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"catalog/product_attribute_backend_media": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendMedia().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Media": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendMedia().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"catalog/product_attribute_backend_msrp": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"catalog/product_attribute_backend_price": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Price": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"catalog/product_attribute_backend_recurring": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendRecurring().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Recurring": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendRecurring().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"catalog/product_attribute_backend_sku": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendSku().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Sku": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendSku().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"catalog/product_attribute_backend_startdate": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendStartDate().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Startdate": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendStartDate().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Stock": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendStock().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"catalog/product_attribute_backend_tierprice": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendTierPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Tierprice": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendTierPrice().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Weight": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductBackendWeight().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"catalog/product_attribute_backend_urlkey": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "",
	},
	"customer/attribute_backend_data_boolean": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerBackendDataBoolean().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Customer\\Model\\Attribute\\Backend\\Data\\Boolean": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerBackendDataBoolean().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"customer/customer_attribute_backend_billing": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerBackendBilling().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Billing": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerBackendBilling().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"customer/customer_attribute_backend_password": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerBackendPassword().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Password": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerBackendPassword().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"customer/customer_attribute_backend_shipping": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerBackendShipping().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Shipping": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerBackendShipping().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"customer/customer_attribute_backend_store": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerBackendStore().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Store": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerBackendStore().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"customer/customer_attribute_backend_website": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerBackendWebsite().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Website": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerBackendWebsite().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"customer/entity_address_attribute_backend_region": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.AddressBackendRegion().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Backend\\Region": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.AddressBackendRegion().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"customer/entity_address_attribute_backend_street": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.AddressBackendStreet().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\DefaultBackend": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.AttributeBackendDefaultBackend().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"eav/entity_attribute_backend_datetime": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.AttributeBackendDatetime().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Datetime": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.AttributeBackendDatetime().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"eav/entity_attribute_backend_time_created": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.AttributeBackendTimeCreated().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Time\\Created": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.AttributeBackendTimeCreated().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"eav/entity_attribute_backend_time_updated": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.AttributeBackendTimeUpdated().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Time\\Updated": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.AttributeBackendTimeUpdated().Config(eav.AttributeBackendIdx({{.AttributeIndex}}))",
	},
	"bundle/product_attribute_source_price_view": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/bundle",
		GoModel:    "bundle.AttributeSourcePriceView().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Bundle\\Model\\Product\\Attribute\\Source\\Price\\View": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/bundle",
		GoModel:    "bundle.AttributeSourcePriceView().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\CatalogInventory\\Model\\Source\\Stock": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/cataloginventory",
		GoModel:    "cataloginventory.AttributeSourceStock().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"catalog/category_attribute_source_layout": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Layout": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "",
	},
	"catalog/category_attribute_source_mode": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.CategorySourceMode().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Mode": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.CategorySourceMode().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"catalog/category_attribute_source_page": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.CategorySourcePage().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Page": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.CategorySourcePage().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"catalog/category_attribute_source_sortby": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.CategorySourceSortby().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Sortby": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.CategorySourceSortby().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"catalog/entity_product_attribute_design_options_container": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductSourceDesignOptionsContainer().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Entity\\Product\\Attribute\\Design\\Options\\Container": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductSourceDesignOptionsContainer().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"catalog/product_attribute_source_countryofmanufacture": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductSourceCountryOfManufacture().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Countryofmanufacture": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductSourceCountryOfManufacture().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"catalog/product_attribute_source_layout": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductSourceLayout().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Layout": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductSourceLayout().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Status": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductSourceStatus().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"catalog/product_attribute_source_msrp_type_enabled": &AttributeModelDef{
		ImportPath: "",
		GoModel:    "",
	},
	"catalog/product_attribute_source_msrp_type_price": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/msrp",
		GoModel:    "msrp.NewAttributeSourcePrice().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Msrp\\Model\\Product\\Attribute\\Source\\Type\\Price": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/msrp",
		GoModel:    "msrp.NewAttributeSourcePrice().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"catalog/product_status": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductSourceStatus().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"catalog/product_visibility": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductSourceVisibility().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Visibility": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductSourceVisibility().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Visibility": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "catattr.ProductSourceVisibility().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"core/design_source_design": &AttributeModelDef{
		ImportPath: "",
		GoModel:    "",
	},
	"customer/customer_attribute_source_group": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerSourceGroup().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Group": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerSourceGroup().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"customer/customer_attribute_source_store": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerSourceStore().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Store": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerSourceStore().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"customer/customer_attribute_source_website": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerSourceWebsite().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Website": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.CustomerSourceWebsite().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"customer/entity_address_attribute_source_country": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.AddressSourceCountry().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Source\\Country": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.AddressSourceCountry().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"customer/entity_address_attribute_source_region": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.AddressSourceRegion().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Source\\Region": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.AddressSourceRegion().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"eav/entity_attribute_source_boolean": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.AttributeSourceBoolean().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Source\\Boolean": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.AttributeSourceBoolean().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"eav/entity_attribute_source_table": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.AttributeSourceTable().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Source\\Table": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.AttributeSourceTable().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"tax/class_source_product": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/tax",
		GoModel:    "tax.AttributeSourceTaxClassProduct().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Tax\\Model\\TaxClass\\Source\\Product": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/tax",
		GoModel:    "tax.AttributeSourceTaxClassProduct().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Tax\\Model\\TaxClass\\Source\\Customer": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/tax",
		GoModel:    "tax.AttributeSourceTaxClassCustomer().Config(eav.AttributeSourceIdx({{.AttributeIndex}}))",
	},
	"Magento\\Theme\\Model\\Theme\\Source\\Theme": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/theme",
		GoModel:    "",
	},
	"customer/attribute_data_postcode": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.AddressDataPostcode().Config(eav.AttributeDataIdx({{.AttributeIndex}}))",
	},
	"Magento\\Customer\\Model\\Attribute\\Data\\Postcode": &AttributeModelDef{
		ImportPath: "github.com/corestoreio/csfw/customer/custattr",
		GoModel:    "custattr.AddressDataPostcode().Config(eav.AttributeDataIdx({{.AttributeIndex}}))",
	},
}
