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

package csdb

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"sync"

	"github.com/corestoreio/csfw/storage/dbr"
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
	Field   string         `db:"COLUMN_NAME"`      //`COLUMN_NAME` varchar(64) NOT NULL DEFAULT '',
	Pos     int64          `db:"ORDINAL_POSITION"` //`ORDINAL_POSITION` bigint(21) unsigned NOT NULL DEFAULT '0',
	Default dbr.NullString `db:"COLUMN_DEFAULT"`   //`COLUMN_DEFAULT` longtext,
	Null    string         `db:"IS_NULLABLE"`      //`IS_NULLABLE` varchar(3) NOT NULL DEFAULT '',
	// DataType contains the basic type of a column like smallint, int, mediumblob,
	// float, double, etc... but always transformed to lower case.
	DataType      string        `db:"DATA_TYPE"`                //`DATA_TYPE` varchar(64) NOT NULL DEFAULT '',
	CharMaxLength dbr.NullInt64 `db:"CHARACTER_MAXIMUM_LENGTH"` //`CHARACTER_MAXIMUM_LENGTH` bigint(21) unsigned DEFAULT NULL,
	Precision     dbr.NullInt64 `db:"NUMERIC_PRECISION"`        //`NUMERIC_PRECISION` bigint(21) unsigned DEFAULT NULL,
	Scale         dbr.NullInt64 `db:"NUMERIC_SCALE"`            //`NUMERIC_SCALE` bigint(21) unsigned DEFAULT NULL,
	// ColumnType full SQL string of the column type
	ColumnType string `db:"COLUMN_TYPE"` //`COLUMN_TYPE` longtext NOT NULL,
	// Key primary or unique or ...
	Key     string `db:"COLUMN_KEY"`     //`COLUMN_KEY` varchar(3) NOT NULL DEFAULT '',
	Extra   string `db:"EXTRA"`          //`EXTRA` varchar(30) NOT NULL DEFAULT '',
	Comment string `db:"COLUMN_COMMENT"` //`COLUMN_COMMENT` varchar(1024) NOT NULL DEFAULT '',

	mu sync.RWMutex
	// DataTypeSimple contains the simplified data type of the field DataType.
	// Fo example bigint, smallint, tinyiny will result in "int".
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
func LoadColumns(ctx context.Context, db dbr.Querier, tables ...string) (map[string]Columns, error) {
	var rows *sql.Rows

	if len(tables) == 0 {
		var err error
		rows, err = db.QueryContext(ctx, selAllTablesColumns)
		if err != nil {
			return nil, errors.Wrapf(err, "[csdb] LoadColumns QueryContext for tables %v", tables)
		}
	} else {
		sqlStr, args, err := dbr.Repeat(selTablesColumns, dbr.ArgString(tables...))
		if err != nil {
			return nil, errors.Wrapf(err, "[csdb] LoadColumns dbr.Repeat for tables %v", tables)
		}
		rows, err = db.QueryContext(ctx, sqlStr, args...)
		if err != nil {
			return nil, errors.Wrapf(err, "[csdb] LoadColumns QueryContext for tables %v", tables)
		}
	}
	defer rows.Close()

	tc := make(map[string]Columns)

	var tn string
	for rows.Next() {
		c := new(Column)
		if err := rows.Scan(&tn, &c.Field, &c.Pos, &c.Default, &c.Null, &c.DataType, &c.CharMaxLength, &c.Precision, &c.Scale, &c.ColumnType, &c.Key, &c.Extra, &c.Comment); err != nil {
			return nil, errors.Wrap(err, "[csdb] Scan Query")
		}

		if _, ok := tc[tn]; !ok {
			tc[tn] = make(Columns, 0, 10)
		}
		c.DataType = strings.ToLower(c.DataType)
		tc[tn] = append(tc[tn], c)
		tn = ""
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "[csdb] rows.Err Query")
	}
	if len(tc) == 0 {
		return nil, errors.NewNotFoundf("[csdb] Tables %v not found", tables)
	}
	return tc, nil
}

// Hash calculates a non-cryptographic, fast and efficient hash value from all
// columns. Current hash algorithm is fnv64a.
func (cs Columns) Hash() ([]byte, error) {
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

// Filter returns a new slice filtered by predicate f
func (cs Columns) Filter(f func(*Column) bool) (cols Columns) {
	for _, c := range cs {
		if f(c) {
			cols = append(cols, c)
		}
	}
	return
}

// Map will run function f on all items in Columns and returns a copy of the
// slice and the item.
func (cs Columns) Map(f func(*Column) *Column) Columns {
	cols := make(Columns, cs.Len())
	for i, c := range cs {
		var c2 = new(Column)
		*c2 = *c
		// columns.go:161::error: assignment copies lock value to *c2: csdb.Column contains sync.RWMutex (vet)
		// hmmm ...
		cols[i] = f(c2)
	}
	return cols
}

// FieldNames returns all column names
func (cs Columns) FieldNames() (fieldNames []string) {
	for _, c := range cs {
		if c.Field != "" {
			fieldNames = append(fieldNames, c.Field)
		}
	}
	return
}

// PrimaryKeys returns all primary key columns
func (cs Columns) PrimaryKeys() Columns {
	return cs.Filter(func(c *Column) bool {
		return c.IsPK()
	})
}

// UniqueKeys returns all unique key columns
func (cs Columns) UniqueKeys() Columns {
	return cs.Filter(func(c *Column) bool {
		return c.IsUnique()
	})
}

// ColumnsNoPK returns all non primary key columns
func (cs Columns) ColumnsNoPK() Columns {
	return cs.Filter(func(c *Column) bool {
		return !c.IsPK()
	})
}

// Len returns the length
func (cs Columns) Len() int {
	return len(cs)
}

// Less compares via the Pos field.
func (cs Columns) Less(i, j int) bool { return cs[i].Pos < cs[j].Pos }

// Swap changes the position
func (cs Columns) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }

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
	_, _ = buf.WriteString("csdb.Columns{\n")
	for _, c := range cs {
		_, _ = fmt.Fprintf(&buf, "%#v,\n", c)
	}
	_ = buf.WriteByte('}')
	return buf.String()
}

// First returns the first column from the slice. Guaranteed to a non-nil return value.
func (cs Columns) First() *Column {
	if len(cs) > 0 {
		return cs[0]
	}
	return new(Column)
}

// JoinFields joins the field names into a string, separated by the optional
// first argument.
func (cs Columns) JoinFields(sep ...string) string {
	aSep := ""
	if len(sep) > 0 {
		aSep = sep[0]
	}
	return strings.Join(cs.FieldNames(), aSep)
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
	// fix tests if you change this layout of the returned string or rename fields.
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	_, _ = buf.WriteString("&csdb.Column{")
	fmt.Fprintf(buf, "Field: %q, ", c.Field)
	if c.Pos > 0 {
		fmt.Fprintf(buf, "Pos: %d, ", c.Pos)
	}
	if c.Default.Valid {
		fmt.Fprintf(buf, "Default: dbr.MakeNullString(%q), ", c.Default.String)
	}
	if c.Null != "" {
		fmt.Fprintf(buf, "Null: %q, ", c.Null)
	}
	if c.DataType != "" {
		fmt.Fprintf(buf, "DataType: %q, ", c.DataType)
	}
	if c.CharMaxLength.Valid {
		fmt.Fprintf(buf, "CharMaxLength: dbr.MakeNullInt64(%d), ", c.CharMaxLength.Int64)
	}
	if c.Precision.Valid {
		fmt.Fprintf(buf, "Precision: dbr.MakeNullInt64(%d), ", c.Precision.Int64)
	}
	if c.Scale.Valid {
		fmt.Fprintf(buf, "Scale: dbr.MakeNullInt64(%d), ", c.Scale.Int64)
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
			goType = "dbr.NullBool"
		}
	case colTypeInt:
		goType = "int64"
		if isNull {
			goType = "dbr.NullInt64"
		}
	case colTypeString:
		goType = "string"
		if isNull {
			goType = "dbr.NullString"
		}
	case colTypeFloat:
		goType = "float64"
		if isNull {
			goType = "dbr.NullFloat64"
		}
	case colTypeDate:
		goType = "time.Time"
		if isNull {
			goType = "dbr.NullTime"
		}
	case colTypeTime:
		goType = "time.Time"
		if isNull {
			goType = "dbr.NullTime"
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
