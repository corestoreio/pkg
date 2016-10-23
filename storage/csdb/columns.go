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
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/null"
	"github.com/corestoreio/csfw/util/slices"
)

// Helper constants to detect certain features of a table and or column.
const (
	ColumnPrimary          = "PRI"
	ColumnUnique           = "UNI"
	ColumnNull             = "YES"
	ColumnNotNull          = "NO"
	ColumnAutoIncrement    = "auto_increment"
	ColumnUnsigned         = "unsigned"
	ColumnCurrentTimestamp = "CURRENT_TIMESTAMP"
)

// Columns contains a slice of column types
type Columns []*Column

// Column contains information about one database column retrieved from
// information_schema.COLUMNS
type Column struct {
	Field   string      `db:"COLUMN_NAME"`      //`COLUMN_NAME` varchar(64) NOT NULL DEFAULT '',
	Pos     int64       `db:"ORDINAL_POSITION"` //`ORDINAL_POSITION` bigint(21) unsigned NOT NULL DEFAULT '0',
	Default null.String `db:"COLUMN_DEFAULT"`   //`COLUMN_DEFAULT` longtext,
	Null    string      `db:"IS_NULLABLE"`      //`IS_NULLABLE` varchar(3) NOT NULL DEFAULT '',
	// DataType contains the basic type of a column like smallint, int, mediumblob, float, double, etc...
	DataType      string     `db:"DATA_TYPE"`                //`DATA_TYPE` varchar(64) NOT NULL DEFAULT '',
	CharMaxLength null.Int64 `db:"CHARACTER_MAXIMUM_LENGTH"` //`CHARACTER_MAXIMUM_LENGTH` bigint(21) unsigned DEFAULT NULL,
	Precision     null.Int64 `db:"NUMERIC_PRECISION"`        //`NUMERIC_PRECISION` bigint(21) unsigned DEFAULT NULL,
	Scale         null.Int64 `db:"NUMERIC_SCALE"`            //`NUMERIC_SCALE` bigint(21) unsigned DEFAULT NULL,
	// TypeRaw SQL string of the column type
	TypeRaw string `db:"COLUMN_TYPE"` //`COLUMN_TYPE` longtext NOT NULL,
	// Key primary or unique or ...
	Key     string `db:"COLUMN_KEY"`     //`COLUMN_KEY` varchar(3) NOT NULL DEFAULT '',
	Extra   string `db:"EXTRA"`          //`EXTRA` varchar(30) NOT NULL DEFAULT '',
	Comment string `db:"COLUMN_COMMENT"` //`COLUMN_COMMENT` varchar(1024) NOT NULL DEFAULT '',
}

// LoadColumns returns all columns from a table in the current database to which
// dbr.SessionRunner has been bound to.
func LoadColumns(dbrSess dbr.SessionRunner, table string) (Columns, error) {
	sel := dbrSess.SelectBySql(`SELECT
		 COLUMN_NAME,ORDINAL_POSITION,COLUMN_DEFAULT,
		 IS_NULLABLE,DATA_TYPE,CHARACTER_MAXIMUM_LENGTH,NUMERIC_PRECISION,
		 NUMERIC_SCALE,COLUMN_TYPE,COLUMN_KEY,EXTRA,COLUMN_COMMENT
	 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME=?`, table)

	selSql, selArg, err := sel.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "[csdb] ToSql")
	}

	rows, err := sel.Query(selSql, selArg...)
	if err != nil {
		return nil, errors.Wrapf(err, "[csdb] Query: %q Args: %#v", selSql, selArg)
	}
	defer rows.Close()

	var cs = make(Columns, 0, 10)
	for rows.Next() {
		c := new(Column)
		err := rows.Scan(&c.Field, &c.Pos, &c.Default, &c.Null, &c.DataType, &c.CharMaxLength, &c.Precision, &c.Scale, &c.TypeRaw, &c.Key, &c.Extra, &c.Comment)
		if err != nil {
			return nil, errors.Wrapf(err, "[csdb] Scan Query: %q Args: %#v", selSql, selArg)
		}
		cs = append(cs, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "[csdb] rows.Err Query: %q Args: %#v", selSql, selArg)
	}
	if len(cs) == 0 {
		return nil, errors.NewNotFoundf("[csdb] Table %q not found in current database connection.", table)
	}
	return cs, nil
}

// Hash calculates a non-cryptographic, fast and efficient hash value from all
// columns. Current hash algorithm is fnv64a.
func (cs Columns) Hash() ([]byte, error) {
	var tr byte = 't' // letter t for true
	var fl byte = 'f' // letter f for false
	var buf bytes.Buffer
	for _, c := range cs {
		buf.WriteString(c.Field)

		buf.WriteString(strconv.Itoa(int(c.Pos)))
		buf.WriteString(c.Default.String)

		if c.IsNull() {
			buf.WriteByte(tr)
		} else {
			buf.WriteByte(fl)
		}

		buf.WriteString(c.DataType)
		buf.WriteString(strconv.Itoa(int(c.CharMaxLength.Int64)))
		buf.WriteString(strconv.Itoa(int(c.Precision.Int64)))
		buf.WriteString(strconv.Itoa(int(c.Scale.Int64)))
		buf.WriteString(c.TypeRaw)
		buf.WriteString(c.Key)
		buf.WriteString(c.Extra)
		buf.WriteString(c.Comment)

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

// ByName finds a column by its name. Case sensitive. Guaranteed to a non-nil
// return value.
func (cs Columns) ByName(fieldName string) *Column {
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
	buf.WriteString("csdb.Columns{\n")
	for _, c := range cs {
		fmt.Fprintf(&buf, "%#v,\n", c)
	}
	buf.WriteByte('}')
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

// IsNull checks if column can have null values
func (c *Column) IsNull() bool {
	return c.Null == ColumnNull
}

// IsPK checks if column is a primary key
func (c *Column) IsPK() bool {
	return c.Field != "" && c.Key == ColumnPrimary
}

// IsPK checks if column is a unique key
func (c *Column) IsUnique() bool {
	return c.Field != "" && c.Key == ColumnUnique
}

// IsAutoIncrement checks if column has an auto increment property
func (c *Column) IsAutoIncrement() bool {
	return c.Field != "" && c.Extra == ColumnAutoIncrement
}

// IsBool checks the name of a column if it contains bool values. Magento uses
// often smallint field types to store bool values.
func (c *Column) IsBool() bool {
	if len(c.Field) < 3 {
		return false
	}
	return columnTypes.byName.bool.ContainsReverse(c.Field)
}

// IsInt checks if a column contains a MySQL int type, independent from its length.
func (c *Column) IsInt() bool {
	switch c.DataType {
	case "bigint", "int", "mediumint", "smallint", "tinyint":
		return true
	}
	return false
}

// IsString checks if a column contains a MySQL varchar or text type.
func (c *Column) IsString() bool {
	switch c.DataType {
	case "longtext", "mediumtext", "text", "tinytext", "varchar", "enum", "char":
		return true
	}
	return false
}

// IsDate checks if a column contains a MySQL timestamp or date type.
func (c *Column) IsDate() bool {
	switch c.DataType {
	case "date", "datetime", "timestamp":
		return true
	}
	return false
}

// IsFloat checks if a column contains a MySQL decimal or float type.
func (c *Column) IsFloat() bool {
	switch c.DataType {
	case "decimal", "float", "double":
		return true
	}
	return false
}

// IsMoney checks if a column contains a MySQL decimal or float type and the
// column name.
// This function needs a lot of care ...
func (c *Column) IsMoney() bool {
	// needs more love
	if !c.IsFloat() {
		return false
	}
	var ret bool
	switch {
	// could us a long list of || statements but switch looks nicer :-)
	case columnTypes.byName.moneyEqual.Contains(c.Field):
		ret = true
	case columnTypes.byName.money.ContainsReverse(c.Field):
		ret = true
	case columnTypes.byName.moneySW.StartsWithReverse(c.Field):
		ret = true
	case !c.IsNull() && c.Default.String == "0.0000":
		ret = true
	}
	return ret
}

// IsUnsigned checks if field TypeRaw contains the word unsigned.
func (c *Column) IsUnsigned() bool {
	return strings.Contains(c.TypeRaw, ColumnUnsigned)
}

// IsCurrentTimestamp checks if the Default field is a current timestamp
func (c *Column) IsCurrentTimestamp() bool {
	return c.Default.String == ColumnCurrentTimestamp
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
	switch {
	case c.IsBool() && isNull:
		goType = "null.Bool"
	case c.IsBool():
		goType = "bool"
	case c.IsInt() && isNull:
		goType = "null.Int64"
	case c.IsInt():
		goType = "int64" // rethink if it is worth to introduce uint64 because of some unsigned columns
	case c.IsString() && isNull:
		goType = "null.String"
	case c.IsString():
		goType = "string"
	case c.IsMoney():
		goType = "money.Money"
	case c.IsFloat() && isNull:
		goType = "null.Float64"
	case c.IsFloat():
		goType = "float64"
	case c.IsDate() && isNull:
		goType = "null.Time"
	case c.IsDate():
		goType = "time.Time"
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
