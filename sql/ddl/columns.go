// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package ddl

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/bufferpool"
	"github.com/corestoreio/pkg/util/slices"
)

// Helper constants to detect certain features of a table and or column.
const (
	columnPrimary          = "PRI"
	columnUnique           = "UNI"
	columnNull             = "YES"
	columnAutoIncrement    = "auto_increment"
	columnUnsigned         = "unsigned"
	columnCurrentTimestamp = "CURRENT_TIMESTAMP"
)

// Columns contains a slice of column types
type Columns []*Column

// Column contains information about one database table column retrieved from
// information_schema.COLUMNS
type Column struct {
	Field   string      // `COLUMN_NAME` varchar(64) NOT NULL DEFAULT '',
	Pos     uint64      // `ORDINAL_POSITION` bigint(21) unsigned NOT NULL DEFAULT '0',
	Default null.String // `COLUMN_DEFAULT` longtext,
	Null    string      // `IS_NULLABLE` varchar(3) NOT NULL DEFAULT '',
	// DataType contains the basic type of a column like smallint, int, mediumblob,
	// float, double, etc... but always transformed to lower case.
	DataType      string     // `DATA_TYPE` varchar(64) NOT NULL DEFAULT '',
	CharMaxLength null.Int64 // `CHARACTER_MAXIMUM_LENGTH` bigint(21) unsigned DEFAULT NULL,
	Precision     null.Int64 // `NUMERIC_PRECISION` bigint(21) unsigned DEFAULT NULL,
	Scale         null.Int64 // `NUMERIC_SCALE` bigint(21) unsigned DEFAULT NULL,
	// ColumnType full SQL string of the column type
	ColumnType string // `COLUMN_TYPE` longtext NOT NULL,
	// Key primary or unique or ...
	Key                  string      // `COLUMN_KEY` varchar(3) NOT NULL DEFAULT '',
	Extra                string      // `EXTRA` varchar(30) NOT NULL DEFAULT '',
	Comment              string      // `COLUMN_COMMENT` varchar(1024) NOT NULL DEFAULT '',
	Generated            string      // `IS_GENERATED` varchar(6) NOT NULL DEFAULT '', MariaDB only https://mariadb.com/kb/en/library/information-schema-columns-table/
	GenerationExpression null.String // `GENERATION_EXPRESSION` longtext DEFAULT NULL, MariaDB only https://mariadb.com/kb/en/library/information-schema-columns-table/
	// Aliases specifies different names used for this column. Mainly used when
	// generating code for interface dml.ColumnMapper. For example
	// customer_entity.entity_id can also be sales_order.customer_id. The alias
	// would be just: entity_id:[]string{"customer_id"}.
	Aliases []string
	// Uniquified used when generating code to uniquify the values in a
	// collection when the column is not a primary or unique key. The values get
	// returned in its own primitive slice.
	Uniquified bool
	// StructTag  used in code generation and applies a custom struct tag.
	StructTag string
}

// TODO check DB flavor: if MySQL or MariaDB, first one does not have column IS_GENERATED
const (
	selBaseSelect = `SELECT
	TABLE_NAME, COLUMN_NAME, ORDINAL_POSITION, COLUMN_DEFAULT, IS_NULLABLE,
		DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE,
		COLUMN_TYPE, COLUMN_KEY, EXTRA, COLUMN_COMMENT, IS_GENERATED, GENERATION_EXPRESSION
	 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE()`
	// DMLLoadColumns specifies the data manipulation language for retrieving
	// all columns in the current database for a specific table. TABLE_NAME is
	// always lower case.
	selTablesColumns    = selBaseSelect + ` AND TABLE_NAME IN ? ORDER BY TABLE_NAME, ORDINAL_POSITION`
	selAllTablesColumns = selBaseSelect + ` ORDER BY TABLE_NAME, ORDINAL_POSITION`
)

// LoadColumns returns all columns from a list of table names in the current
// database. Map key contains the table name. Returns a NotFound error if the
// table is not available. All columns from all tables gets selected when you
// don't provide the argument `tables`.
func LoadColumns(ctx context.Context, db dml.Querier, tables ...string) (map[string]Columns, error) {
	var rows *sql.Rows

	if len(tables) == 0 {
		var err error
		rows, err = db.QueryContext(ctx, selAllTablesColumns)
		if err != nil {
			return nil, errors.Wrapf(err, "[ddl] LoadColumns QueryContext for tables %v", tables)
		}
	} else {
		sqlStr, _, err := dml.Interpolate(selTablesColumns).Strs(tables...).ToSQL()
		if err != nil {
			return nil, errors.Wrapf(err, "[ddl] LoadColumns dml.ExpandPlaceHolders for tables %v", tables)
		}
		rows, err = db.QueryContext(ctx, sqlStr)
		if err != nil {
			return nil, errors.Wrapf(err, "[ddl] LoadColumns QueryContext for tables %v with WHERE clause", tables)
		}
	}
	var err error
	defer func() {
		// Not testable with the sqlmock package :-(
		if err2 := rows.Close(); err2 != nil && err == nil {
			err = errors.WithStack(err2)
		}
	}()

	tc := make(map[string]Columns)
	rc := new(dml.ColumnMap)
	for rows.Next() {
		if err = rc.Scan(rows); err != nil {
			return nil, errors.Wrapf(err, "[ddl] Scan Query for tables: %v", tables)
		}
		var c *Column
		var tableName string
		c, tableName, err = NewColumn(rc)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if _, ok := tc[tableName]; !ok {
			tc[tableName] = make(Columns, 0, 10)
		}

		c.DataType = strings.ToLower(c.DataType)
		tc[tableName] = append(tc[tableName], c)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	if len(tc) == 0 {
		return nil, errors.NotFound.Newf("[ddl] Tables %v not found", tables)
	}
	return tc, err
}

// Filter filters the columns by predicate f and appends the column pointers to
// the optional argument `cols`.
func (cs Columns) Filter(f func(*Column) bool, cols ...*Column) Columns {
	for _, c := range cs {
		if f(c) {
			cols = append(cols, c)
		}
	}
	return cols
}

// Each applies function f to all elements.
func (cs Columns) Each(f func(*Column)) Columns {
	for _, c := range cs {
		f(c)
	}
	return cs
}

// FieldNames returns all column names and appends it to `fn`, if provided.
func (cs Columns) FieldNames(fn ...string) []string {
	if fn == nil {
		fn = make([]string, 0, len(cs))
	}
	for _, c := range cs {
		if c.Field != "" {
			fn = append(fn, c.Field)
		}
	}
	return fn
}

func colIsPK(c *Column) bool {
	return c.IsPK() && !c.IsSystemVersioned()
}

func colIsUnique(c *Column) bool {
	return c.IsUnique() && !c.IsSystemVersioned()
}

func colIsNotSysVers(c *Column) bool {
	return !c.IsSystemVersioned()
}

func colIsNotGeneratedNonPK(c *Column) bool {
	return !c.IsGenerated() && !c.IsSystemVersioned() && c.Extra != "auto_increment"
}

// PrimaryKeys returns all primary key columns. It may append the columns to the
// provided argument slice.
func (cs Columns) PrimaryKeys(cols ...*Column) Columns {
	return cs.Filter(colIsPK, cols...)
}

// UniqueKeys returns all unique key columns. It may append the columns to the
// provided argument slice.
func (cs Columns) UniqueKeys(cols ...*Column) Columns {
	return cs.Filter(colIsUnique, cols...)
}

// UniqueColumns returns all columns which are either a single primary key or a
// single unique key. If a PK or UK consists of more than one column, then they
// won't be included in the returned Columns slice. The result might be appended
// to argument `cols`, if provided.
func (cs Columns) UniqueColumns(cols ...*Column) Columns {
	if cols == nil {
		cols = make(Columns, 0, 3) // 3 is just a guess
	}
	pkCount, ukCount := 0, 0
	for _, c := range cs {
		if c.IsPK() {
			pkCount++
		}
		if c.IsUnique() {
			ukCount++
		}
	}
	if pkCount >= 1 {
		cols = cs.PrimaryKeys(cols...)
	}
	if ukCount >= 1 {
		cols = cs.UniqueKeys(cols...)
	}
	return cols
}

func colIsNotPK(c *Column) bool {
	return !c.IsPK() && !c.IsUnique() && c.GenerationExpression.String == ""
}

// NonPrimaryColumns returns all non primary key and non-unique key columns.
func (cs Columns) NonPrimaryColumns() Columns {
	return cs.Filter(colIsNotPK)
}

// Len returns the length
func (cs Columns) Len() int {
	return len(cs)
}

// Less compares via the Pos field.
func (cs Columns) Less(i, j int) bool { return cs[i].Pos < cs[j].Pos }

// Swap changes the position
func (cs Columns) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }

// Contains returns true if fieldName is contained in slice Columns.
func (cs Columns) Contains(fieldName string) bool {
	for _, c := range cs {
		if c.Field == fieldName {
			return true
		}
	}
	return false
}

func colIsNotUniquified(c *Column) bool {
	return c.Uniquified && c.GenerationExpression.String == ""
}

// UniquifiedColumns returns all columns which have the flag Uniquified set to
// true. The result might be appended to argument `cols`, if provided.
func (cs Columns) UniquifiedColumns(cols ...*Column) Columns {
	return cs.Filter(colIsNotUniquified, cols...)
}

// ByField finds a column by its field name. Case sensitive. Guaranteed to
// return a non-nil return value.
func (cs Columns) ByField(fieldName string) *Column {
	for _, c := range cs {
		if c.Field == fieldName {
			return c
		}
	}
	return new(Column)
}

// @todo add maybe more ByNull(), ByType(), ByKey(), ByDefault(), ByExtra()

// String same as GoString()
func (cs Columns) String() string {
	return cs.GoString()
}

// GoString returns the Go types representation. See interface fmt.GoStringer
func (cs Columns) GoString() string {
	// fix tests if you change this layout of the returned string
	var buf bytes.Buffer
	_, _ = buf.WriteString("ddl.Columns{\n")
	for _, c := range cs {
		_, _ = fmt.Fprintf(&buf, "%#v,\n", c)
	}
	_ = buf.WriteByte('}')
	return buf.String()
}

// First returns the first column from the slice. Guaranteed to a non-nil return
// value.
func (cs Columns) First() *Column {
	if len(cs) > 0 {
		return cs[0]
	}
	return new(Column)
}

// JoinFields joins the field names into a string, separated by the provided
// separator.
func (cs Columns) JoinFields(sep string) string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	i := 0
	for _, c := range cs {
		if c.Field != "" {
			if i > 0 {
				buf.WriteString(sep)
			}
			buf.WriteString(c.Field)
			i++
		}
	}
	return buf.String()
}

// NewColumn creates a new column pointer and maps it from a raw database row
// its bytes into the type Column.
func NewColumn(rc *dml.ColumnMap) (c *Column, tableName string, err error) {
	c = new(Column)
	for rc.Next() {
		switch col := rc.Column(); col {
		case "TABLE_NAME":
			rc.String(&tableName)
		case "COLUMN_NAME":
			rc.String(&c.Field)
		case "ORDINAL_POSITION":
			rc.Uint64(&c.Pos)
		case "COLUMN_DEFAULT":
			rc.NullString(&c.Default)
		case "IS_NULLABLE":
			rc.String(&c.Null)
		case "DATA_TYPE":
			rc.String(&c.DataType)
		case "CHARACTER_MAXIMUM_LENGTH":
			rc.NullInt64(&c.CharMaxLength)
		case "NUMERIC_PRECISION":
			rc.NullInt64(&c.Precision)
		case "NUMERIC_SCALE":
			rc.NullInt64(&c.Scale)
		case "COLUMN_TYPE":
			rc.String(&c.ColumnType)
		case "COLUMN_KEY":
			rc.String(&c.Key)
		case "EXTRA":
			rc.String(&c.Extra)
		case "COLUMN_COMMENT":
			rc.String(&c.Comment)
		case "IS_GENERATED":
			rc.String(&c.Generated)
		case "GENERATION_EXPRESSION":
			rc.NullString(&c.GenerationExpression)
		case "aliases":
			// TODO the query must be extendable for all three columns to attach any table from any DB.
			if aliases := ""; rc.Mode() == dml.ColumnMapScan {
				rc.String(&aliases)
				c.Aliases = strings.Split(aliases, ",")
			} else {
				aliases = strings.Join(c.Aliases, ",")
				rc.String(&aliases)
			}
		case "uniquified":
			rc.Bool(&c.Uniquified)
		case "struct_tag":
			rc.String(&c.StructTag)
		default:
			return nil, "", errors.NotSupported.Newf("[ddl] Column %q not supported or alias not found", col)
		}
	}
	return c, tableName, errors.WithStack(rc.Err())
}

// GoComment creates a comment from a database column to be used in Go code
func (c *Column) GoComment() string {
	sqlNull := "NOT NULL"
	if c.IsNull() {
		sqlNull = "NULL"
	}
	sqlDefault := ""
	if c.Default.Valid {
		sqlDefault = "DEFAULT '" + c.Default.String + "'"
	}
	return fmt.Sprintf("// %s %s %s %s %s %s %q",
		c.Field, c.ColumnType, sqlNull, c.Key, sqlDefault, c.Extra, c.Comment,
	)
}

// GoString returns the Go types representation. See interface fmt.GoStringer
func (c *Column) GoString() string {
	// mauybe this can be removed ...
	// fix tests if you change this layout of the returned string or rename columns.
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	_, _ = buf.WriteString("&ddl.Column{")
	fmt.Fprintf(buf, "Field: %q, ", c.Field)
	if c.Pos > 0 {
		fmt.Fprintf(buf, "Pos: %d, ", c.Pos)
	}
	if c.Default.Valid {
		fmt.Fprintf(buf, "Default: null.MakeString(%q), ", c.Default.String)
	}
	if c.Null != "" {
		fmt.Fprintf(buf, "Null: %q, ", c.Null)
	}
	if c.DataType != "" {
		fmt.Fprintf(buf, "DataType: %q, ", c.DataType)
	}
	if c.CharMaxLength.Valid {
		fmt.Fprintf(buf, "CharMaxLength: null.MakeInt64(%d), ", c.CharMaxLength.Int64)
	}
	if c.Precision.Valid {
		fmt.Fprintf(buf, "Precision: null.MakeInt64(%d), ", c.Precision.Int64)
	}
	if c.Scale.Valid {
		fmt.Fprintf(buf, "Scale: null.MakeInt64(%d), ", c.Scale.Int64)
	}
	if c.ColumnType != "" {
		fmt.Fprintf(buf, "ColumnType: %q, ", c.ColumnType)
	}
	if c.Key != "" {
		fmt.Fprintf(buf, "Key: %q, ", c.Key)
	}
	if c.Extra != "" {
		fmt.Fprintf(buf, "Extra: %q, ", c.Extra)
	}
	if c.Comment != "" {
		fmt.Fprintf(buf, "Comment: %q, ", c.Comment)
	}
	if len(c.Aliases) > 0 {
		fmt.Fprintf(buf, "Aliases: %#v, ", c.Aliases)
	}
	if c.Uniquified {
		fmt.Fprintf(buf, "Uniquified: %t, ", c.Uniquified)
	}
	if c.StructTag != "" {
		fmt.Fprintf(buf, "StructTag: %q, ", c.StructTag)
	}
	if c.Generated != "" {
		fmt.Fprintf(buf, "Generated: %q, ", c.Generated)
	}
	if c.GenerationExpression.Valid {
		fmt.Fprintf(buf, "GenerationExpression: %q, ", c.GenerationExpression.String)
	}
	_ = buf.WriteByte('}')
	return buf.String()
}

// IsNull checks if column can have null values
func (c *Column) IsNull() bool {
	return c.Null == columnNull
}

// IsPK checks if column is a primary key
func (c *Column) IsPK() bool {
	return c.Field != "" && c.Key == columnPrimary
}

// IsUnique checks if column is a unique key
func (c *Column) IsUnique() bool {
	return c.Field != "" && c.Key == columnUnique
}

// IsAutoIncrement checks if column has an auto increment property
func (c *Column) IsAutoIncrement() bool {
	return c.Field != "" && c.Extra == columnAutoIncrement
}

// IsUnsigned checks if field TypeRaw contains the word unsigned.
func (c *Column) IsUnsigned() bool {
	return strings.Contains(c.ColumnType, columnUnsigned)
}

// IsCurrentTimestamp checks if the Default field is a current timestamp
func (c *Column) IsCurrentTimestamp() bool {
	return c.Default.String == columnCurrentTimestamp
}

// IsGenerated returns true if the column is a virtual generated column.
func (c *Column) IsGenerated() bool {
	return c.Generated == "ALWAYS" || c.GenerationExpression.Valid
}

// IsSystemVersioned returns true if the column gets used for system versioning.
// https://mariadb.com/kb/en/library/system-versioned-tables/
func (c *Column) IsSystemVersioned() bool {
	return c.GenerationExpression.Valid && (c.GenerationExpression.String == "ROW START" || c.GenerationExpression.String == "ROW END")
}

// IsFloat returns true if a column is of one of the types: decimal, double or
// float.
func (c *Column) IsFloat() bool {
	switch c.DataType {
	case "decimal", "double", "float":
		return true
	}
	return false
}

// IsMoney checks if a column contains a MySQL decimal or float type and if the
// column name has a special naming.
// This function needs a lot of care ...
func (c *Column) IsMoney() bool {
	// needs more love
	switch {
	// could us a long list of || statements but switch looks nicer :-)
	case columnTypes.byName.moneyEqual.Contains(c.Field):
		return true
	case columnTypes.byName.money.ContainsReverse(c.Field):
		return true
	case columnTypes.byName.moneySW.StartsWithReverse(c.Field):
		return true
	}
	return false
}

// IsBool returns true if column is of type `int` and its name starts with a
// special string like: `used_`, `is_`, `has_`.
func (c *Column) IsBool() (ok bool) {
	switch c.DataType {
	case "int", "tinyint", "smallint", "bigint":
		ok = true
	case "bit":
		return true
	}
	return ok && columnTypes.byName.bool.ContainsReverse(c.Field)
}

// IsString returns true if the column can contain a string or byte values.
func (c *Column) IsString() bool {
	return c.CharMaxLength.Valid && c.CharMaxLength.Int64 > 0
}

// IsBlobDataType returns true if the columns data type is neither blob,
// text, binary nor json. It doesn't matter if tiny, long or small has been
// prefixed.
func (c *Column) IsBlobDataType() bool {
	dt := strings.ToLower(c.DataType)
	return strings.Contains(dt, "blob") || strings.Contains(dt, "text") ||
		strings.Contains(dt, "binary") || strings.Contains(dt, "json")
}

// columnTypes looks ugly but ... refactor later.
// the slices in this struct are only for reading. no mutex protection required.
// which partial column name triggers a specific type in Go or MySQL.
var columnTypes = struct {
	byName struct {
		bool       slices.String
		money      slices.String
		moneySW    slices.String
		moneyEqual slices.String
	}
}{
	struct {
		bool       slices.String // contains
		money      slices.String // contains
		moneySW    slices.String // sw == starts with
		moneyEqual slices.String
	}{
		slices.String{"used_", "is_", "has_", "increment_per_store"},
		slices.String{"price", "_tax", "tax_", "_amount", "amount_", "total", "adjustment", "discount"},
		slices.String{"base_", "grand_"},
		slices.String{"value", "price", "cost", "msrp"},
	},
}
