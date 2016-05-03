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

package codegen

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/log"
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
	// ValueSuffixes contains the suffixes for an entity type table, e.g. datetime, decimal, int ... then a
	// table value name would be catalog_product_entity_datetime, catalog_product_entity_decimal, ...
	ValueSuffixes []string
	// TypeCodeValueTable 2 dimensional map. 1. key entity_type_code 2. key table name => value ValueSuffix
	TypeCodeValueTable map[string]map[string]string
)

func (vs ValueSuffixes) contains(suffix string) bool {
	for _, v := range vs {
		if v == suffix {
			return true
		}
	}
	return false
}

// String joins the slice of strings separated by a comma. Only for debug.
func (vs ValueSuffixes) String() string {
	return strings.Join(vs, ", ")
}

// Empty checks if the map is empty or has an empty "" entry.
func (m TypeCodeValueTable) Empty() bool {
	_, ok := m[""]
	return len(m) < 1 || ok
}

// GetTables returns all tables from a database. AndWhere can be optionally applied.
// Only first index (0) will be added.
func GetTables(dbrSess dbr.SessionRunner, sql ...string) ([]string, error) {

	qry := "SHOW TABLES"
	if len(sql) > 0 && sql[0] != "" {
		if false == dbr.Stmt.IsSelect(sql[0]) {
			qry = qry + " LIKE '" + sql[0] + "'"
		} else {
			qry = sql[0]
		}
	}
	if PkgLog.IsDebug() { // this if reduces 9 allocs ...
		defer log.WhenDone(PkgLog).Debug("Stats", "Package", "codegen", "Step", "GetTables", "query", qry)
	}

	sb := dbrSess.SelectBySql(qry)
	query, args, err := sb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "ToSql")
	}
	rows, err := sb.Query(query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "Query %q", query)
	}
	defer rows.Close()

	var tableName string
	var tableNames = make([]string, 0, 200)
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			return nil, errors.Wrapf(err, "Scan Query %q", query)
		}
		tableNames = append(tableNames, tableName)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "Rows")
	}
	return tableNames, nil
}

// GetEavValueTables returns a map of all custom and default EAV value tables for entity type codes.
// Despite value_table_prefix can have in Magento a different table name we treat it here
// as the table name itself. Not thread safe.
func GetEavValueTables(dbrConn *dbr.Connection, entityTypeCodes []string) (TypeCodeValueTable, error) {

	typeCodeTables := make(TypeCodeValueTable, len(entityTypeCodes))

	for _, typeCode := range entityTypeCodes {

		vtp, err := dbrConn.NewSession().
			Select("`value_table_prefix`").
			From(TablePrefix+TableEavEntityType).
			Where(dbr.ConditionRaw("`value_table_prefix` IS NOT NULL"), dbr.ConditionRaw("`entity_type_code` = ?", typeCode)).
			ReturnString()

		if err != nil && err != dbr.ErrNotFound {
			return nil, errors.Wrapf(err, "Select Error. Code %q", typeCode)
		}
		if vtp == "" {
			vtp = typeCode + TableNameSeparator + TableEntityTypeSuffix + TableNameSeparator // e.g. catalog_product_entity_
		} else {
			vtp = vtp + TableNameSeparator
		}

		tableNames, err := GetTables(dbrConn.NewSession(), vtp+`%`)
		if err != nil {
			return nil, errors.Wrapf(err, "GetTables %q", typeCode)
		}

		if _, ok := typeCodeTables[typeCode]; !ok {
			typeCodeTables[typeCode] = make(map[string]string, len(tableNames))
		}
		for _, t := range tableNames {
			valueSuffix := t[len(vtp):]
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
	// Columns contains a slice to pointer column types. A column has the fields: Field, Type, Null, Key, Default and
	// Extra of type sql.NullString and GoType, GoName of type string.
	Columns []column
	// column contains info about one database column retrieve from SHOW COLUMNS FROM tbl`
	column struct {
		csdb.Column
		GoType, GoName string
	}
)

// Comment creates a comment from a database column to be used in Go code
func (c column) Comment() string {
	sqlNull := "NOT NULL"
	if c.IsNull() {
		sqlNull = "NULL"
	}
	sqlDefault := ""
	if c.Default.String != "" {
		sqlDefault = "DEFAULT '" + c.Default.String + "'"
	}
	return "// " + c.Name() + " " + c.Type.String + " " + sqlNull + " " + c.Key.String + " " + sqlDefault + " " + c.Extra.String
}

func (c column) updateGoPrimitive(useSQL bool) column {
	c.GoName = util.UnderscoreCamelize(c.Field.String)
	c.GoType = c.GetGoPrimitive(useSQL)
	return c
}

// CopyToCSDB copies the underlying slice to a csdb columns type
func (cc Columns) CopyToCSDB() csdb.Columns {
	ret := make(csdb.Columns, len(cc))
	for i, c := range cc {
		ret[i] = c.Column
	}
	return ret
}

// GetByName returns a column from Columns slice by a give name
func (cc Columns) GetByName(name string) column {
	for _, c := range cc {
		if c.Field.String == name {
			return c
		}
	}
	return column{}
}

// MapSQLToGoDBRType takes a slice of Columns and sets the fields GoType and GoName to the correct value
// to create a Go struct. These generated structs are mainly used in a result from a SQL query. The field GoType
// will contain dbr.Null* types.
func (cc Columns) MapSQLToGoDBRType() error {
	for i, col := range cc {
		cc[i] = col.updateGoPrimitive(true)
	}
	return nil
}

// MapSQLToGoType maps a column to a GoType. This GoType is not a dbr.Null* struct. This function only updates
// the fields GoType and GoName of column struct. The 2nd argument ifm interface map replaces the primitive type
// with an interface type, the column name must be found as a key in the map.
func (cc Columns) MapSQLToGoType(ifm map[string]string) error {
	for i, col := range cc {
		cc[i] = col.updateGoPrimitive(false)
		if val, ok := ifm[col.Field.String]; col.Field.Valid && ok {
			cc[i].GoType = val // Type is now an interface name
		}
	}
	return nil
}

// GetFieldNames returns from a Columns slice the column names. If pkOnly is true then only the
// primary key columns will be returned.
func (cc Columns) GetFieldNames(pkOnly bool) []string {
	ret := make([]string, 0, len(cc))
	for _, col := range cc {
		isPk := col.Key.String == "PRI"
		if pkOnly && isPk {
			ret = append(ret, col.Field.String)
		}
		if !pkOnly && !isPk {
			ret = append(ret, col.Field.String)
		}
	}
	return ret
}

// isIgnoredColumn Drop unused column entity_type_id in customer__* and catalog_* tables
func isIgnoredColumn(table, column string) bool {
	const etid = "entity_type_id"
	switch {
	case strings.Contains(table, "catalog_") && column == etid:
		return true
	case strings.Contains(table, "customer_") && column == etid:
		return true
	case strings.Contains(table, "eav_attribute") && column == "attribute_model":
		return true
	}
	return false
}

// GetColumns returns all columns from a table. It discards the column
// entity_type_id from some entity tables. The column attribute_model will also
// be dropped from table eav_attribute
func GetColumns(db *sql.DB, table string) (Columns, error) {
	var cols = make(Columns, 0, 200)
	rows, err := db.Query("SHOW COLUMNS FROM `" + table + "`")
	if err != nil {
		return nil, errors.Wrapf(err, "Query Table %q", table)
	}
	defer rows.Close()

	col := column{}
	for rows.Next() {

		err := rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra)
		if err != nil {
			return nil, errors.Wrapf(err, "Query Scan Table %q", table)
		}
		if isIgnoredColumn(table, col.Field.String) {
			continue
		}
		cols = append(cols, col)
	}
	return cols, errors.Wrapf(rows.Err(), "Rows Error Table %q", table)
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
func SQLQueryToColumns(db *sql.DB, dbSelect *dbr.SelectBuilder, query ...string) (Columns, error) {

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
		var err error
		qry, args, err = dbSelect.ToSql()
		if err != nil {
			return nil, errors.Wrap(err, "dbSelect.ToSql")
		}
	}
	_, err := db.Exec("CREATE TABLE `"+tableName+"` AS "+qry, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "Create Table %q", tableName)
	}

	return GetColumns(db, tableName)
}

// ColumnsToStructCode generates Go code from a name and a slice of columns.
// Make sure that the fields GoType and GoName has been setup
// If you don't like the template you can provide your own template as 3rd to n-th argument.
func ColumnsToStructCode(tplData map[string]interface{}, name string, cols Columns, templates ...string) ([]byte, error) {

	if nil == tplData {
		tplData = make(map[string]interface{})
	}

	tplData["Name"] = name
	tplData["Columns"] = cols
	tplData["Tick"] = "`"

	tpl := strings.Join(templates, "")
	if tpl == "" {
		tpl = tplQueryDBRStruct
	}

	return GenerateCode("", tpl, tplData, nil)
}

// LoadStringEntities executes a SELECT query and returns a slice containing columns names and its string values
func LoadStringEntities(db *sql.DB, dbSelect *dbr.SelectBuilder) ([]StringEntities, error) {

	qry, args, err := dbSelect.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "ToSQL")
	}

	rows, err := db.Query(qry, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "Query %q", qry)
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrap(err, "rows.Columns")
	}

	ret := make([]StringEntities, 0, 2000)
	rss := newRowTransformer(columnNames)
	for rows.Next() {

		if err := rows.Scan(rss.cp...); err != nil {
			return nil, errors.Wrap(err, "Scan")
		}
		err := rss.toString()
		if err != nil {
			return nil, errors.Wrap(err, "toString")
		}
		rss.append(&ret)
	}
	return ret, errors.Wrap(rows.Err(), "Rows.err")
}

type (
	// StringEntities contains as key the column name and value the string value from the column.
	// sql.RawBytes are converted to a string.
	StringEntities map[string]string
	rowTransformer struct {
		// cp are the column pointers
		cp []interface{}
		// row contains the final row result
		se       StringEntities
		colCount int
		colNames []string
	}
)

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
			return errors.NewFatalf("Cannot convert index %d column %q to type *sql.RawBytes", i, s.colNames[i])
		}
	}
	return nil
}

// append appends the current row to the ret return value and clears the row result
func (s *rowTransformer) append(ret *[]StringEntities) {
	*ret = append(*ret, s.se)
	s.se = make(StringEntities, len(s.colNames))
}

func validImportPath(ad *AttributeModelDef, ip []string, targetPkg string) bool {
	if ad.GoFunc == "" {
		return false
	}
	if len(targetPkg) > 0 && ad.Import()[len(ad.Import())-len(targetPkg):] == targetPkg {
		return false
	}
	add := true
	for i := 0; i < len(ip); i++ {
		if ip[i] == ad.Import() { // check for duplicates
			add = false
		}
	}
	if add {
		return true
	}
	return false
}

// PrepareForTemplate uses the Columns slice to transform the rows so that correct Go code can be printed.
// int/Float values won't be touched. Bools or IntBools will be converted to true/false. Strings will be quoted.
// And if there is an entry in the AttributeModelMap then the Go code from the map will be used.
// Returns a slice containing all the import paths. Import paths which are equal to pkg will be filtered out.
func PrepareForTemplate(cols Columns, rows []StringEntities, amm AttributeModelDefMap, targetPkg string) []string {
	ip := make([]string, 0, 10) // import_path container
	for _, row := range rows {
		for colName, colValue := range row {
			var c = cols.GetByName(colName)

			if false == c.Field.Valid {
				continue
			}

			goType, hasModel := amm[colValue]
			_, isAllowedInterfaceChange := EavAttributeColumnNameToInterface[colName]
			switch {
			case hasModel:
				row[colName] = "nil"
				if goType.GoFunc != "" {
					row[colName] = goType.Func()
					if validImportPath(goType, ip, targetPkg) {
						ip = append(ip, goType.Import())
					}
				}
				break
			case isAllowedInterfaceChange:
				// if there is no defined model but column is (backend|frontend|data|source)_model then nil it
				row[colName] = "nil"
				break
			case c.IsBool():
				row[colName] = "false"
				if colValue == "1" {
					row[colName] = "true"
				}
				break
			case c.IsInt():
				if colValue == "" {
					row[colName] = "0"
				}
				break
			case c.IsString():
				row[colName] = strconv.Quote(colValue)
				break
			case c.IsFloat():
				if colValue == "" {
					row[colName] = "0.0"
				}
				break
			case c.IsDate():
				if colValue == "" {
					row[colName] = "nil"
				} else {
					row[colName] = "time.Parse(`2006-01-02 15:04:05`," + strconv.Quote(colValue) + ")" // @todo timezone
				}
				break
			default:
				panic(fmt.Sprintf("\nERROR cannot detect SQL type: %s -> %s\n%#v\n", colName, colValue, c))
			}

		}
	}
	return ip
}
