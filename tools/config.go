package tools

type (
	// TableToStructMap uses a string key as easy identifier for maybe later manipulation in config_user.go
	// and a pointer to a TableToStruct struct. Developers can later use the init() func in config_user.go to change
	// the value of variable ConfigTableToStruct.
	TableToStructMap map[string]*TableToStruct
	TableToStruct    struct {
		// Package defines the name of the target package
		Package string
		// OutputFile specifies the path where to write the newly generated code
		OutputFile string
		// QueryString SQL query to filter all the tables which you desire, e.g. SHOW TABLES LIKE 'catalog\_%'
		QueryString string
		// EntityTypeCodes If provided then eav_entity_type.value_table_prefix will be evaluated for further tables.
		EntityTypeCodes []string
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
		TempAdditionalAttributeTable string

		// TempAdditionalAttributeTableWebsite string which defines the existing table name
		// and stores website-dependent attribute parameters.
		// If an EAV model doesn't demand this functionality, let this string empty.
		// This table name is used in a DB query while materializing attribute configuration to Go code.
		TempAdditionalAttributeTableWebsite string
	}

	// AttributeModelMap contains data provided via JSON to map the three eav_attribute columns
	// (backend|frontend|source)_model to the correct Go function and package
	// DefaultAttributeModelMap contains default mappings for Mage1+2. A developer has the option to provide a custom map.
	// Rethink the Go code here ... because catalog.Product().Attribute().Frontend().Image() is pretty long ... BUT
	// developers coming from Magento are already familiar with this code base and naming ...
	AttributeModelMap map[string]*AttributeModel
	AttributeModel    struct {
		ImportPath string
		GoModel    string
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

// EavAttributeColumnNameToInterface mapping column name to Go interface name
var EavAttributeColumnNameToInterface = map[string]string{
	"backend_model":  "eav.AttributeBackendModeller",
	"frontend_model": "eav.AttributeFrontendModeller",
	"source_model":   "eav.AttributeSourceModeller",
}

// TablePrefix defines the global table name prefix. See Magento install tool. Can be override via func init()
var TablePrefix string = ""

// ConfigTableToStruct contains default configuration. Use the file config_user.go with the func init() to change/extend it.
var ConfigTableToStruct = TableToStructMap{
	"eav": &TableToStruct{
		Package:         "eav",
		OutputFile:      "eav/generated_tables.go",
		QueryString:     `SHOW TABLES LIKE "` + TablePrefix + `eav%"`,
		EntityTypeCodes: nil,
	},
	"catalog": &TableToStruct{
		Package:    "catalog",
		OutputFile: "catalog/generated_tables.go",
		QueryString: `SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE
		    TABLE_SCHEMA = DATABASE() AND
		    (TABLE_NAME LIKE '` + TablePrefix + `catalog\_%' OR TABLE_NAME LIKE '` + TablePrefix + `catalogindex%' ) AND
		    TABLE_NAME NOT LIKE '` + TablePrefix + `%bundle%' AND
		    TABLE_NAME NOT LIKE '` + TablePrefix + `%\_flat\_%' GROUP BY TABLE_NAME;`,
		EntityTypeCodes: []string{"catalog_category", "catalog_product"},
	},
	"customer": &TableToStruct{
		Package:         "customer",
		OutputFile:      "customer/generated_tables.go",
		QueryString:     `SHOW TABLES LIKE "` + TablePrefix + `customer%"`,
		EntityTypeCodes: []string{"customer", "customer_address"},
	},
	// @todo extend for all sales_* tables
}

var ConfigEntityTypeMaterialization = &TableToStruct{
	Package:    "materialized",
	OutputFile: "materialized/generated_eav_entity.go",
}

// ConfigEntityType contains default configuration. Use the file config_user.go with the func init() to change/extend it.
var ConfigEntityType = EntityTypeMap{
	"customer": &EntityType{
		ImportPath:                          "github.com/corestoreio/csfw/customer",
		EntityModel:                         "customer.Customer()",
		AttributeModel:                      "customer.Attribute()",
		EntityTable:                         "customer.Customer()",
		IncrementModel:                      "customer.Customer()",
		AdditionalAttributeTable:            "customer.Customer()",
		EntityAttributeCollection:           "customer.Customer()",
		TempAdditionalAttributeTable:        TablePrefix + "customer_eav_attribute",
		TempAdditionalAttributeTableWebsite: TablePrefix + "customer_eav_attribute_website",
	},
	"customer_address": &EntityType{
		ImportPath:                          "github.com/corestoreio/csfw/customer",
		EntityModel:                         "customer.Address()",
		AttributeModel:                      "customer.AddressAttribute()",
		EntityTable:                         "customer.Address()",
		AdditionalAttributeTable:            "customer.Address()",
		EntityAttributeCollection:           "customer.Address()",
		TempAdditionalAttributeTable:        TablePrefix + "customer_eav_attribute",
		TempAdditionalAttributeTableWebsite: TablePrefix + "customer_eav_attribute_website",
	},
	"catalog_category": &EntityType{
		ImportPath:                   "github.com/corestoreio/csfw/catalog",
		EntityModel:                  "catalog.Category()",
		AttributeModel:               "catalog.Attribute()",
		EntityTable:                  "catalog.Category()",
		AdditionalAttributeTable:     "catalog.Category()",
		EntityAttributeCollection:    "catalog.Category()",
		TempAdditionalAttributeTable: TablePrefix + "catalog_eav_attribute",
	},
	"catalog_product": &EntityType{
		ImportPath:                   "github.com/corestoreio/csfw/catalog",
		EntityModel:                  "catalog.Product()",
		AttributeModel:               "catalog.Attribute()",
		EntityTable:                  "catalog.Product()",
		AdditionalAttributeTable:     "catalog.Product()",
		EntityAttributeCollection:    "catalog.Product()",
		TempAdditionalAttributeTable: TablePrefix + "catalog_eav_attribute",
	},
	// @todo extend for all sales entities
}

var ConfigAttributeModel = AttributeModelMap{
	"catalog/product_attribute_frontend_image": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Frontend().Image()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Frontend\\Image": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Frontend().Image()",
	},
	"eav/entity_attribute_frontend_datetime": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Entity().Attribute().Frontend().Datetime()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Frontend\\Datetime": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Entity().Attribute().Frontend().Datetime()",
	},
	"catalog/attribute_backend_customlayoutupdate": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Attribute().Backend().Customlayoutupdate()",
	},
	"Magento\\Catalog\\Model\\Attribute\\Backend\\Customlayoutupdate": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Attribute().Backend().Customlayoutupdate()",
	},
	"catalog/category_attribute_backend_image": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Backend().Image()",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Backend\\Image": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Backend().Image()",
	},
	"catalog/category_attribute_backend_sortby": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Backend().Sortby()",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Backend\\Sortby": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Backend().Sortby()",
	},
	"catalog/category_attribute_backend_urlkey": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "",
	},
	"catalog/product_attribute_backend_boolean": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Boolean()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Boolean": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Boolean()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Category": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Category()",
	},
	"catalog/product_attribute_backend_groupprice": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Groupprice()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\GroupPrice": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().GroupPrice()",
	},
	"catalog/product_attribute_backend_media": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Media()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Media": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Media()",
	},
	"catalog/product_attribute_backend_msrp": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Price()",
	},
	"catalog/product_attribute_backend_price": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Price()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Price": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Price()",
	},
	"catalog/product_attribute_backend_recurring": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Recurring()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Recurring": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Recurring()",
	},
	"catalog/product_attribute_backend_sku": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Sku()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Sku": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Sku()",
	},
	"catalog/product_attribute_backend_startdate": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Startdate()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Startdate": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Startdate()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Stock": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Stock()",
	},
	"catalog/product_attribute_backend_tierprice": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Tierprice()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Tierprice": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Tierprice()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Backend\\Weight": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Backend().Weight()",
	},
	"catalog/product_attribute_backend_urlkey": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "",
	},
	"customer/attribute_backend_data_boolean": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Data().Boolean()",
	},
	"Magento\\Customer\\Model\\Attribute\\Backend\\Data\\Boolean": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Data().Boolean()",
	},
	"customer/customer_attribute_backend_billing": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Billing()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Billing": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Billing()",
	},
	"customer/customer_attribute_backend_password": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Password()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Password": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Password()",
	},
	"customer/customer_attribute_backend_shipping": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Shipping()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Shipping": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Shipping()",
	},
	"customer/customer_attribute_backend_store": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Store()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Store": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Store()",
	},
	"customer/customer_attribute_backend_website": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Website()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Backend\\Website": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Backend().Website()",
	},
	"customer/entity_address_attribute_backend_region": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Backend().Region()",
	},
	"Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Backend\\Region": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Backend().Region()",
	},
	"customer/entity_address_attribute_backend_street": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Backend().Street()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\DefaultBackend": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Backend().DefaultBackend()",
	},
	"eav/entity_attribute_backend_datetime": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Backend().Datetime()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Datetime": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Backend().Datetime()",
	},
	"eav/entity_attribute_backend_time_created": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Entity().Attribute().Backend().Time().Created()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Time\\Created": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Backend().Time().Created()",
	},
	"eav/entity_attribute_backend_time_updated": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Backend().Time().Updated()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Backend\\Time\\Updated": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Backend().Time().Updated()",
	},
	"bundle/product_attribute_source_price_view": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/bundle",
		GoModel:    "bundle.Product().Attribute().Source().Price().View()",
	},
	"Magento\\Bundle\\Model\\Product\\Attribute\\Source\\Price\\View": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/bundle",
		GoModel:    "bundle.Product().Attribute().Source().Price().View()",
	},
	"Magento\\CatalogInventory\\Model\\Source\\Stock": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/cataloginventory",
		GoModel:    "cataloginventory.Source().Stock()",
	},
	"catalog/category_attribute_source_layout": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Layout": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "",
	},
	"catalog/category_attribute_source_mode": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Source().Mode()",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Mode": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Source().Mode()",
	},
	"catalog/category_attribute_source_page": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Source().Page()",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Page": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Source().Page()",
	},
	"catalog/category_attribute_source_sortby": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Source().Sortby()",
	},
	"Magento\\Catalog\\Model\\Category\\Attribute\\Source\\Sortby": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Category().Attribute().Source().Sortby()",
	},
	"catalog/entity_product_attribute_design_options_container": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Design().Options().Container()",
	},
	"Magento\\Catalog\\Model\\Entity\\Product\\Attribute\\Design\\Options\\Container": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Design().Options().Container()",
	},
	"catalog/product_attribute_source_countryofmanufacture": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Countryofmanufacture()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Countryofmanufacture": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Countryofmanufacture()",
	},
	"catalog/product_attribute_source_layout": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Layout()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Layout": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Layout()",
	},
	"Magento\\Catalog\\Model\\Product\\Attribute\\Source\\Status": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Status()",
	},
	"catalog/product_attribute_source_msrp_type_enabled": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Msrp().Type().Enabled()",
	},
	"catalog/product_attribute_source_msrp_type_price": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Attribute().Source().Msrp().Type().Price()",
	},
	"catalog/product_status": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Status()",
	},
	"catalog/product_visibility": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Visibility()",
	},
	"Magento\\Catalog\\Model\\Product\\Visibility": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/catalog",
		GoModel:    "Product().Visibility()",
	},
	"core/design_source_design": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/core",
		GoModel:    "core.Design().Source().Design()",
	},
	"customer/customer_attribute_source_group": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Source().Group()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Group": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Source().Group()",
	},
	"customer/customer_attribute_source_store": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Source().Store()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Store": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Source().Store()",
	},
	"customer/customer_attribute_source_website": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Source().Website()",
	},
	"Magento\\Customer\\Model\\Customer\\Attribute\\Source\\Website": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Source().Website()",
	},
	"customer/entity_address_attribute_source_country": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Source().Country()",
	},
	"Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Source\\Country": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Source().Country()",
	},
	"customer/entity_address_attribute_source_region": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Source().Region()",
	},
	"Magento\\Customer\\Model\\Resource\\Address\\Attribute\\Source\\Region": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Address().Attribute().Source().Region()",
	},
	"eav/entity_attribute_source_boolean": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Source().Boolean()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Source\\Boolean": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Source().Boolean()",
	},
	"eav/entity_attribute_source_table": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Source().Table()",
	},
	"Magento\\Eav\\Model\\Entity\\Attribute\\Source\\Table": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/eav",
		GoModel:    "eav.Attribute().Source().Table()",
	},
	"Magento\\Msrp\\Model\\Product\\Attribute\\Source\\Type\\Price": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/msrp",
		GoModel:    "msrp.Product().Attribute().Source().Type().Price()",
	},
	"tax/class_source_product": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/tax",
		GoModel:    "tax.Class().Source().Product()",
	},
	"Magento\\Tax\\Model\\TaxClass\\Source\\Product": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/tax",
		GoModel:    "tax.TaxClass().Source().Product()",
	},
	"Magento\\Theme\\Model\\Theme\\Source\\Theme": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/theme",
		GoModel:    "",
	},
	"customer/attribute_data_postcode": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Data.Postcode()",
	},
	"Magento\\Customer\\Model\\Attribute\\Data\\Postcode": &AttributeModel{
		ImportPath: "github.com/corestoreio/csfw/customer",
		GoModel:    "Attribute().Data.Postcode()",
	},
}
