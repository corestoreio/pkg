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
