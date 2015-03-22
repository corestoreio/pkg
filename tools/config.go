package tools

type (
	// TableToStructMap uses a string key as easy identifier for maybe later manipulation in config_user.go
	// and a pointer to a TableToStruct struct. Developers can later use the init() func in config_user.go to change
	// the value of ConfigTableToStruct.
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
)

// TablePrefix defines the global table name prefix. See Magento install tool. Can be override via func init()
var TablePrefix string = ""

// ConfigTableToStruct contains base main configuration. Use the file config_user.go with the func init() to change/extend it.
var ConfigTableToStruct = TableToStructMap{
	"eav": &TableToStruct{
		Package:         "eav",
		OutputFile:      "eav/generated_tables.go",
		QueryString:     `SHOW TABLES LIKE "eav%"`,
		EntityTypeCodes: nil,
	},
	"catalog": &TableToStruct{
		Package:    "catalog",
		OutputFile: "catalog/generated_tables.go",
		QueryString: `SELECT TABLE_NAME FROM information_schema.COLUMNS WHERE
		    TABLE_SCHEMA = DATABASE() AND
		    (TABLE_NAME LIKE 'catalog\_%' OR TABLE_NAME LIKE 'catalogindex%' ) AND
		    TABLE_NAME NOT LIKE '%bundle%' AND
		    TABLE_NAME NOT LIKE '%\_flat\_%' GROUP BY TABLE_NAME;`,
		EntityTypeCodes: []string{"catalog_category", "catalog_product"},
	},
	"customer": &TableToStruct{
		Package:         "customer",
		OutputFile:      "customer/generated_tables.go",
		QueryString:     `SHOW TABLES LIKE "customer%"`,
		EntityTypeCodes: []string{"customer", "customer_address"},
	},
	// @todo extend for all sales_* tables
}
