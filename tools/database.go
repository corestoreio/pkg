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
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

const (
	TableNameSeparator string = "_"
	TableEavEntityType string = "eav_entity_type"
)

var (
	// TableEntityTypeSuffix e.g. for catalog_product_entity, customer_entity
	TableEntityTypeSuffix = "entity"
	// TableEntityTypeValueSuffixes defines all possible value type tables which an EAV model can have.
	TableEntityTypeValueSuffixes = ValueSuffixes{
		"datetime",
		"decimal",
		"int",
		"text",
		"varchar",
	}
)

type (
	ValueSuffixes      []string
	TypeCodeValueTable map[string]map[string]string // 1. key entity_type_code 2. key table name => value ValueSuffix

	// EntityTypeMap applies a JSON map to the Go type of EntityType struct
	EntityTypeMap struct {
		ImportPath                string `json:"import_path"`
		EntityTypeID              int64  `db:"entity_type_id"`
		EntityTypeCode            string `db:"entity_type_code"`
		EntityModel               string `json:"entity_model"`
		AttributeModel            string `json:"attribute_model"`
		EntityTable               string `json:"entity_table"`
		ValueTablePrefix          string `db:"value_table_prefix"`
		EntityIDField             string
		IsDataSharing             bool
		DataSharingKey            string
		DefaultAttributeSetID     int64
		IncrementModel            string `json:"increment_model"`
		IncrementPerStore         bool
		IncrementPadLength        int64
		IncrementPadChar          string
		AdditionalAttributeTable  string `json:"additional_attribute_table"`
		EntityAttributeCollection string `json:"entity_attribute_collection"`
	}
)

func (vs ValueSuffixes) contains(suffix string) bool {
	for _, v := range vs {
		if v == suffix {
			return true
		}
	}
	return false
}

func (vs ValueSuffixes) String() string {
	return strings.Join(vs, ", ")
}

func (m TypeCodeValueTable) Empty() bool {
	_, ok := m[""]
	return len(m) < 1 || ok
}

// GetTables returns all tables from a database which starts with a prefix. % wild card will be added
// automatically.
func GetTables(db *sql.DB, prefix string) ([]string, error) {

	var tableNames = make([]string, 0, 200)
	qry := "SHOW TABLES like '" + prefix + "%'"

	rows, err := db.Query(qry)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		tableNames = append(tableNames, tableName)
	}
	err = rows.Err()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return tableNames, nil
}

// GetEavValueTables returns a map of all custom and default EAV value tables for entity type codes.
// Despite value_table_prefix can have in Magento a different table name we treat it here
// as the table name itself. Not thread safe.
func GetEavValueTables(dbrConn *dbr.Connection, prefix string, entityTypeCodes []string) (TypeCodeValueTable, error) {

	typeCodeTables := make(TypeCodeValueTable, len(entityTypeCodes))

	for _, typeCode := range entityTypeCodes {

		vtp, err := dbrConn.NewSession(nil).
			Select("`value_table_prefix`").
			From(prefix+TableEavEntityType).
			Where("`value_table_prefix` IS NOT NULL").
			Where("`entity_type_code` = ?", typeCode).
			ReturnString()

		if err != nil && err != dbr.ErrNotFound {
			return nil, errgo.Mask(err)
		}
		if vtp == "" {
			vtp = typeCode + TableNameSeparator + TableEntityTypeSuffix + TableNameSeparator // e.g. catalog_product_entity_
		} else {
			vtp = vtp + TableNameSeparator
		}

		tableNames, err := GetTables(dbrConn.Db, prefix+vtp)
		if err != nil {
			return nil, errgo.Mask(err)
		}

		if _, ok := typeCodeTables[typeCode]; !ok {
			typeCodeTables[typeCode] = make(map[string]string, len(tableNames))
		}
		for _, t := range tableNames {
			valueSuffix := t[len(prefix+vtp):]
			if TableEntityTypeValueSuffixes.contains(valueSuffix) {
				/*
				   other tables like catalog_product_entity_gallery, catalog_product_entity_group_price,
				   catalog_product_entity_tier_price, etc are the backend model tables for different storage systems.
				   they are not part of the default EAV model.
				*/
				typeCodeTables[typeCode][t] = valueSuffix
			}

		}

	}

	return typeCodeTables, nil
}

type (
	// column internal struct for a database column
	column struct {
		Field, Type, Null, Key, Default, Extra sql.NullString
		GoType, GoName                         string
	}
)

// Comment creates a comment from a database column to be used in Go code
func (c column) Comment() string {
	sqlNull := "NOT NULL"
	if c.Null.String == "YES" {
		sqlNull = "NULL"
	}
	sqlDefault := ""
	if c.Default.String != "" {
		sqlDefault = "DEFAULT '" + c.Default.String + "'"
	}
	return "// " + c.Field.String + " " + c.Type.String + " " + sqlNull + " " + c.Key.String + " " + sqlDefault + " " + c.Extra.String
}

// GetColumns returns all columns from a table. It discards the column entity_type_id from some
// entity tables.
func GetColumns(db *sql.DB, table string) ([]*column, error) {
	var columns = make([]*column, 0, 200)
	rows, err := db.Query("SHOW COLUMNS FROM `" + table + "`")
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer rows.Close()

	// Drop unused column entity_type_id in customer__* and catalog_* tables
	isEntityTypeIdFree := strings.Index(table, "catalog_") >= 0 || strings.Index(table, "customer_") >= 0

	for rows.Next() {
		col := &column{}
		err := rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		if isEntityTypeIdFree && col.Field.String == "entity_type_id" {
			continue
		}
		columns = append(columns, col)
	}
	err = rows.Err()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return columns, nil
}

// isBoolColumn checks the name of a column if it contains bool values. Magento uses often smallint field types
// to store bool values and also to store other integer numbers.
func isBoolColumn(cp *column) bool {
	if len(cp.Field.String) < 3 {
		return false
	}
	switch cp.Field.String[:3] {
	case "is_", "has":
		return true
	}
	return strings.Index(cp.Field.String, "used_") == 0 || cp.Field.String == "increment_per_store"
}

// MapSQLToGoType maps a column to a GoType. This GoType is not a dbr.Null* struct. This function only updates
// the fields GoType and GoName of column struct. The 2nd argument ifm interface map replaces the primitive type
// with an interface type, the column name must be found as a key in the map.
func MapSQLToGoType(colums []*column, ifm map[string]string) error {
	for _, col := range colums {
		col.GoType = "undefined"
		col.GoName = Camelize(col.Field.String)
		if isBoolColumn(col) {
			col.GoType = "bool"
		} else if strings.Contains(col.Type.String, "int") {
			col.GoType = "int64"
		} else if strings.Contains(col.Type.String, "varchar") || strings.Contains(col.Type.String, "text") {
			col.GoType = "string"
		} else if strings.Contains(col.Type.String, "decimal") || strings.Contains(col.Type.String, "float") {
			col.GoType = "float64"
		} else if strings.Contains(col.Type.String, "timestamp") || strings.Contains(col.Type.String, "date") {
			col.GoType = "time.Time"
		}
		if val, ok := ifm[col.Field.String]; ok {
			col.GoType = val // Type is now an interface name
		}
	}
	return nil
}

// MapSQLToGoDBRType takes a slice of columns an sets the fields GoType and GoName to the correct value
// to create a Go struct. These generated structs are mainly used in a result from a SQL query
func MapSQLToGoDBRType(colums []*column) error {

	for _, col := range colums {
		col.GoType = "undefined"
		col.GoName = Camelize(col.Field.String)
		if isBoolColumn(col) {
			col.GoType = "bool"
			if col.Null.String == "YES" {
				col.GoType = "dbr.NullBool"
			}
		} else if strings.Contains(col.Type.String, "int") {
			col.GoType = "int64"
			if col.Null.String == "YES" {
				col.GoType = "dbr.NullInt64"
			}
		} else if strings.Contains(col.Type.String, "varchar") || strings.Contains(col.Type.String, "text") {
			col.GoType = "string"
			if col.Null.String == "YES" {
				col.GoType = "dbr.NullString"
			}
		} else if strings.Contains(col.Type.String, "decimal") || strings.Contains(col.Type.String, "float") {
			col.GoType = "float64"
			if col.Null.String == "YES" {
				col.GoType = "dbr.NullFloat64"
			}
		} else if strings.Contains(col.Type.String, "timestamp") || strings.Contains(col.Type.String, "date") {
			col.GoType = "time.Time"
			if col.Null.String == "YES" {
				col.GoType = "dbr.NullTime"
			}
		}
	}
	return nil
}

const tplQueryDBRStruct = `
type (
    // {{.Name | prepareVar}}Slice contains pointers to {{.Name | prepareVar}} types
    {{.Name | prepareVar}}Slice []*{{.Name | prepareVar}}
    // {{.Name | prepareVar}} a type for a MySQL Query
    {{.Name | prepareVar}} struct {
        {{ range .Columns }}{{.GoName}} {{.GoType}} {{ $.Tick }}db:"{{.Field.String}}"{{ $.Tick }} {{.Comment}}
        {{ end }} }
)
`

// SQLQueryToColumns generates from a SQL query an array containing all the column properties.
// dbSelect argument can be nil but then you must provide query strings which will be joined to the final query.
func SQLQueryToColumns(db *sql.DB, dbSelect *dbr.SelectBuilder, query ...string) ([]*column, error) {

	tableName := "tmp_" + randSeq(20)
	dropTable := func() {
		_, err := db.Exec("DROP TABLE IF EXISTS `" + tableName + "`")
		if err != nil {
			panic(err)
		}
	}
	dropTable()
	defer dropTable()

	qry := strings.Join(query, " ")
	var args []interface{}
	if qry == "" && dbSelect != nil {
		qry, args = dbSelect.ToSql()
	}
	_, err := db.Exec("CREATE TABLE `"+tableName+"` AS "+qry, args...)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	return GetColumns(db, tableName)
}

// ColumnsToStructCode generates Go code from a name and a slice of columns.
// If you don't like the template you can provide your own template as 3rd to n-th argument.
func ColumnsToStructCode(name string, cols []*column, templates ...string) ([]byte, error) {

	tplData := struct {
		Name    string
		Columns []*column
		Tick    string
	}{
		Name:    name,
		Columns: cols,
		Tick:    "`",
	}

	tpl := strings.Join(templates, "")
	if tpl == "" {
		tpl = tplQueryDBRStruct
	}

	return GenerateCode("", tpl, tplData)
}

// GetSQL executes a SELECT query and returns a slice containing columns names and its string values
func GetSQL(db *sql.DB, dbSelect *dbr.SelectBuilder, query ...string) ([]StringEntities, error) {

	qry := strings.Join(query, " ")
	var args []interface{}
	if qry == "" && dbSelect != nil {
		qry, args = dbSelect.ToSql()
	}

	rows, err := db.Query(qry, args...)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, errgo.Mask(err)
	}

	ret := make([]StringEntities, 0, 2000)
	rss := newRowTransformer(columnNames)
	for rows.Next() {

		if err := rows.Scan(rss.cp...); err != nil {
			return nil, errgo.Mask(err)
		}
		err := rss.toString()
		if err != nil {
			return nil, errgo.Mask(err)
		}
		rss.append(&ret)
	}
	return ret, nil
}

type rowTransformer struct {
	// cp are the column pointers
	cp []interface{}
	// row contains the final row result
	se       StringEntities
	colCount int
	colNames []string
}

func newRowTransformer(columnNames []string) *rowTransformer {
	lenCN := len(columnNames)
	s := &rowTransformer{
		cp:       make([]interface{}, lenCN),
		se:       make(StringEntities, len(columnNames)),
		colCount: lenCN,
		colNames: columnNames,
	}
	for i := 0; i < lenCN; i++ {
		s.cp[i] = new(sql.RawBytes)
	}
	return s
}

func (s *rowTransformer) toString() error {
	for i := 0; i < s.colCount; i++ {
		if rb, ok := s.cp[i].(*sql.RawBytes); ok {
			s.se[s.colNames[i]] = string(*rb)
			*rb = nil // reset pointer to discard current value to avoid a bug
		} else {
			return errors.New("Cannot convert index " + strconv.Itoa(i) + " column " + s.colNames[i] + " to type *sql.RawBytes")
		}
	}
	return nil
}

// append appends the current row to the ret return value and clears the row result
func (s *rowTransformer) append(ret *[]StringEntities) {
	*ret = append(*ret, s.se)
	s.se = make(StringEntities, len(s.colNames))
}

type (
	StringEntities map[string]string
)

//func (s *StringEntities) Column(c string) string {
//	if v, ok := s.row[c]; ok {
//		return v
//	}
//	return ""
//}
//
//func (s *StringEntities) Columns() []string {
//	return s.col
//}
