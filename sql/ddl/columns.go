// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"hash/fnv"
	"strconv"
	"strings"
	"sync"

	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/slices"
	"github.com/corestoreio/errors"
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
	Field   string         //`COLUMN_NAME` varchar(64) NOT NULL DEFAULT '',
	Pos     uint64         //`ORDINAL_POSITION` bigint(21) unsigned NOT NULL DEFAULT '0',
	Default dml.NullString //`COLUMN_DEFAULT` longtext,
	Null    string         //`IS_NULLABLE` varchar(3) NOT NULL DEFAULT '',
	// DataType contains the basic type of a column like smallint, int, mediumblob,
	// float, double, etc... but always transformed to lower case.
	DataType      string        //`DATA_TYPE` varchar(64) NOT NULL DEFAULT '',
	CharMaxLength dml.NullInt64 //`CHARACTER_MAXIMUM_LENGTH` bigint(21) unsigned DEFAULT NULL,
	Precision     dml.NullInt64 //`NUMERIC_PRECISION` bigint(21) unsigned DEFAULT NULL,
	Scale         dml.NullInt64 //`NUMERIC_SCALE` bigint(21) unsigned DEFAULT NULL,
	// ColumnType full SQL string of the column type
	ColumnType string //`COLUMN_TYPE` longtext NOT NULL,
	// Key primary or unique or ...
	Key     string //`COLUMN_KEY` varchar(3) NOT NULL DEFAULT '',
	Extra   string //`EXTRA` varchar(30) NOT NULL DEFAULT '',
	Comment string //`COLUMN_COMMENT` varchar(1024) NOT NULL DEFAULT '',

	mu sync.RWMutex
	// DataTypeSimple contains the simplified data type of the field DataType.
	// Fo example bigint, smallint, tinyint will result in "int".
	dataTypeSimple string
}

// DMLLoadColumns specifies the data manipulation language for retrieving all
// columns in the current database for a specific table.
const selTablesColumns = `SELECT
	TABLE_NAME, COLUMN_NAME, ORDINAL_POSITION, COLUMN_DEFAULT, IS_NULLABLE,
		DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE,
		COLUMN_TYPE, COLUMN_KEY, EXTRA, COLUMN_COMMENT	
	 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME IN (?)
	 ORDER BY TABLE_NAME, ORDINAL_POSITION`

const selAllTablesColumns = `SELECT
	TABLE_NAME, COLUMN_NAME, ORDINAL_POSITION, COLUMN_DEFAULT, IS_NULLABLE,
		DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE,
		COLUMN_TYPE, COLUMN_KEY, EXTRA, COLUMN_COMMENT
	 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() ORDER BY TABLE_NAME, ORDINAL_POSITION`

// LoadColumns returns all columns from a list of table names in the current
// database. For now MySQL DSN must have set interpolateParams to true. Map key
// contains the table name. Returns a NotFound error if the table is not
// available. All columns from all tables gets selected when you don't provide
// the argument `tables`.
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
			return nil, errors.Wrapf(err, "[ddl] LoadColumns dml.Repeat for tables %v", tables)
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
			err = errors.Wrap(err2, "[dml] LoadColumns.Rows.Close")
		}
	}()

	tc := make(map[string]Columns)
	rc := new(dml.RowConvert)
	for rows.Next() {
		if err = rc.Scan(rows); err != nil {
			return nil, errors.Wrapf(err, "[ddl] Scan Query for tables: %v", tables)
		}
		c, tn, err := NewColumn(rc)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if _, ok := tc[tn]; !ok {
			tc[tn] = make(Columns, 0, 10)
		}

		c.DataType = strings.ToLower(c.DataType)
		tc[tn] = append(tc[tn], c)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "[ddl] rows.Err Query")
	}
	if len(tc) == 0 {
		return nil, errors.NewNotFoundf("[ddl] Tables %v not found", tables)
	}
	return tc, err
}

// Hash calculates a non-cryptographic, fast and efficient hash value from all
// columns. Current hash algorithm is fnv64a.
func (cs Columns) Hash() ([]byte, error) {
	// TODO use encoding/binary
	var tr byte = 't' // letter t for true
	var fl byte = 'f' // letter f for false
	var buf bytes.Buffer
	for _, c := range cs {
		_, _ = buf.WriteString(c.Field)

		_, _ = buf.WriteString(strconv.Itoa(int(c.Pos)))
		_, _ = buf.WriteString(c.Default.String)

		if c.IsNull() {
			_ = buf.WriteByte(tr)
		} else {
			_ = buf.WriteByte(fl)
		}

		_, _ = buf.WriteString(c.DataType)
		_, _ = buf.WriteString(strconv.Itoa(int(c.CharMaxLength.Int64)))
		_, _ = buf.WriteString(strconv.Itoa(int(c.Precision.Int64)))
		_, _ = buf.WriteString(strconv.Itoa(int(c.Scale.Int64)))
		_, _ = buf.WriteString(c.ColumnType)
		_, _ = buf.WriteString(c.Key)
		_, _ = buf.WriteString(c.Extra)
		_, _ = buf.WriteString(c.Comment)
	}

	f64 := fnv.New64a()
	if _, err := f64.Write(buf.Bytes()); err != nil {
		return nil, err
	}
	buf.Reset()
	return f64.Sum(buf.Bytes()), nil
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

// Map will run function f on all items in Columns and returns a copy of the
// slice and the item.
func (cs Columns) Map(f func(*Column) *Column) Columns {
	cols := make(Columns, cs.Len())
	for i, c := range cs {
		var c2 = new(Column)
		*c2 = *c
		// columns.go:161::error: assignment copies lock value to *c2: ddl.Column contains sync.RWMutex (vet)
		// hmmm ...
		cols[i] = f(c2)
	}
	return cols
}

// FieldNames returns all column names
func (cs Columns) FieldNames() []string {
	fieldNames := make([]string, 0, len(cs))
	for _, c := range cs {
		if c.Field != "" {
			fieldNames = append(fieldNames, c.Field)
		}
	}
	return fieldNames
}

func colIsPK(c *Column) bool {
	return c.IsPK()
}

// PrimaryKeys returns all primary key columns
func (cs Columns) PrimaryKeys() Columns {
	return cs.Filter(colIsPK)
}

func colIsUnique(c *Column) bool {
	return c.IsUnique()
}

// UniqueKeys returns all unique key columns
func (cs Columns) UniqueKeys() Columns {
	return cs.Filter(colIsUnique)
}

func colIsNotPK(c *Column) bool {
	return !c.IsPK()
}

// ColumnsNoPK returns all non primary key columns
func (cs Columns) ColumnsNoPK() Columns {
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
func NewColumn(rc *dml.RowConvert) (c *Column, tableName string, err error) {
	c = new(Column)
	for i, col := range rc.Columns {
		if rc.Alias != nil {
			// col might be the alias name used in the query. Lookup the
			// original upper case name as used in the switch mapping.
			if orgCol, ok := rc.Alias[col]; ok {
				col = orgCol
			}
		}

		rc = rc.Index(i)
		var err error
		switch col {
		case "TABLE_NAME":
			tableName, err = rc.Str()
		case "COLUMN_NAME":
			c.Field, err = rc.Str()
		case "ORDINAL_POSITION":
			c.Pos, err = rc.Uint64()
		case "COLUMN_DEFAULT":
			c.Default.NullString, err = rc.NullString()
		case "IS_NULLABLE":
			c.Null, err = rc.Str()
		case "DATA_TYPE":
			c.DataType, err = rc.Str()
		case "CHARACTER_MAXIMUM_LENGTH":
			c.CharMaxLength.NullInt64, err = rc.NullInt64()
		case "NUMERIC_PRECISION":
			c.Precision.NullInt64, err = rc.NullInt64()
		case "NUMERIC_SCALE":
			c.Scale.NullInt64, err = rc.NullInt64()
		case "COLUMN_TYPE":
			c.ColumnType, err = rc.Str()
		case "COLUMN_KEY":
			c.Key, err = rc.Str()
		case "EXTRA":
			c.Extra, err = rc.Str()
		case "COLUMN_COMMENT":
			c.Comment, err = rc.Str()
		default:
			return nil, "", errors.NewNotSupportedf("[ddl] Column %q not supported or alias not found")
		}
		if err != nil {
			return nil, "", errors.Wrapf(err, "[ddl] Failed to scan %q at row %d", col, rc.Count)
		}
	}
	return c, tableName, nil
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
	// fix tests if you change this layout of the returned string or rename columns.
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	_, _ = buf.WriteString("&ddl.Column{")
	fmt.Fprintf(buf, "Field: %q, ", c.Field)
	if c.Pos > 0 {
		fmt.Fprintf(buf, "Pos: %d, ", c.Pos)
	}
	if c.Default.Valid {
		fmt.Fprintf(buf, "Default: dml.MakeNullString(%q), ", c.Default.String)
	}
	if c.Null != "" {
		fmt.Fprintf(buf, "Null: %q, ", c.Null)
	}
	if c.DataType != "" {
		fmt.Fprintf(buf, "DataType: %q, ", c.DataType)
	}
	if c.CharMaxLength.Valid {
		fmt.Fprintf(buf, "CharMaxLength: dml.MakeNullInt64(%d), ", c.CharMaxLength.Int64)
	}
	if c.Precision.Valid {
		fmt.Fprintf(buf, "Precision: dml.MakeNullInt64(%d), ", c.Precision.Int64)
	}
	if c.Scale.Valid {
		fmt.Fprintf(buf, "Scale: dml.MakeNullInt64(%d), ", c.Scale.Int64)
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

const (
	colTypeBool   = "bool"
	colTypeByte   = "bytes"
	colTypeDate   = "date"
	colTypeFloat  = "float"
	colTypeInt    = "int"
	colTypeMoney  = "money"
	colTypeString = "string"
	colTypeTime   = "time"
)

// DataTypeSimple calculates the simplified data type of the field DataType. The
// calculated result will be cached. For example bigint, smallint, tinyint will
// result in "int". The returned string guarantees to be lower case. Available
// returned types are: bool, bytes, date, float, int, money, string, time. Data
// type money is special for the database schema. This function is thread safe.
func (c *Column) DataTypeSimple() string {
	c.mu.RLock()
	dts := c.dataTypeSimple
	c.mu.RUnlock()
	if dts != "" {
		return dts
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.dataTypeSimple = "undefnied"

	switch c.DataType {
	case "bigint", "int", "mediumint", "smallint", "tinyint":
		c.dataTypeSimple = colTypeInt
	case "longtext", "mediumtext", "text", "tinytext", "varchar", "enum", "char":
		c.dataTypeSimple = colTypeString
	case "longblob", "mediumblob", "blob", "varbinary", "binary":
		c.dataTypeSimple = colTypeByte
	case "date", "datetime", "timestamp":
		c.dataTypeSimple = colTypeDate
	case "time":
		c.dataTypeSimple = colTypeTime
	case "decimal", "float", "double":
		c.dataTypeSimple = colTypeFloat
	case "bit":
		c.dataTypeSimple = colTypeBool
	}

	switch {
	case columnTypes.byName.bool.ContainsReverse(c.Field):
		c.dataTypeSimple = colTypeBool
	case c.dataTypeSimple == colTypeFloat && c.isMoney():
		c.dataTypeSimple = colTypeMoney
	}
	return c.dataTypeSimple
}

// isMoney checks if a column contains a MySQL decimal or float type and if the
// column name has a special naming.
// This function needs a lot of care ...
func (c *Column) isMoney() bool {
	// needs more love
	switch {
	// could us a long list of || statements but switch looks nicer :-)
	case columnTypes.byName.moneyEqual.Contains(c.Field):
		return true
	case columnTypes.byName.money.ContainsReverse(c.Field):
		return true
	case columnTypes.byName.moneySW.StartsWithReverse(c.Field):
		return true
	case !c.IsNull() && c.Default.String == "0.0000":
		return true
	}
	return false
}

// GoPrimitive detects the Go type of a SQL table column as a non nullable type.
func (c *Column) GoPrimitive() string {
	return c.goPrimitive(false)
}

// GoPrimitiveNull detects the Go type of a SQL table column as a nullable type.
func (c *Column) GoPrimitiveNull() string {
	return c.goPrimitive(true)
}

func (c *Column) goPrimitive(useNullType bool) string {
	var goType = "undefined"
	isNull := c.IsNull() && useNullType
	switch c.DataTypeSimple() {
	case colTypeBool:
		goType = "bool"
		if isNull {
			goType = "dml.NullBool"
		}
	case colTypeInt:
		goType = "int64"
		if isNull {
			goType = "dml.NullInt64"
		}
	case colTypeString:
		goType = "string"
		if isNull {
			goType = "dml.NullString"
		}
	case colTypeFloat:
		goType = "float64"
		if isNull {
			goType = "dml.NullFloat64"
		}
	case colTypeDate:
		goType = "time.Time"
		if isNull {
			goType = "dml.NullTime"
		}
	case colTypeTime:
		goType = "time.Time"
		if isNull {
			goType = "dml.NullTime"
		}
	case colTypeMoney:
		goType = "money.Money"
	}
	return goType
}

// columnTypes looks ugly but ... refactor later
var columnTypes = struct { // the slices in this struct are only for reading. no mutex protection required
	byName struct {
		bool       slices.String
		money      slices.String
		moneySW    slices.String
		moneyEqual slices.String
	}
	byType struct {
		int     slices.String
		string  slices.String
		dateSW  slices.String
		floatSW slices.String
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
	struct {
		int     slices.String // contains
		string  slices.String // contains
		dateSW  slices.String // SW starts with
		floatSW slices.String // SW starts with
	}{
		slices.String{"int"},
		slices.String{"char", "text"},
		slices.String{"time", "date"},
		slices.String{"decimal", "float", "double"},
	},
}
