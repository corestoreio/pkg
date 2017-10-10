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
	"strconv"
	"strings"

	"hash"

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
	rc := new(dml.ColumnMap)
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
// columns.
func (cs Columns) Hash(h hash.Hash) ([]byte, error) {
	// TODO use encoding/binary
	const tr = 't' // letter t for true
	const fl = 'f' // letter f for false
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

	if _, err := h.Write(buf.Bytes()); err != nil {
		return nil, err
	}
	buf.Reset()
	return h.Sum(buf.Bytes()), nil
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
		default:
			return nil, "", errors.NewNotSupportedf("[ddl] Column %q not supported or alias not found", col)
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
func (c *Column) IsBool() bool {
	var isInt bool
	switch c.DataType {
	case "int", "tinyint", "smallint", "bigint":
		isInt = true
	case "bit":
		return true
	}
	return isInt && columnTypes.byName.bool.ContainsReverse(c.Field)
}

// columnTypes looks ugly but ... refactor later
var columnTypes = struct { // the slices in this struct are only for reading. no mutex protection required
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
