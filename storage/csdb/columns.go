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
	"strings"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util"
	"github.com/corestoreio/csfw/util/errors"
)

const (
	ColumnPrimary       = "PRI"
	ColumnUnique        = "UNI"
	ColumnNull          = "YES"
	ColumnNotNull       = "NO"
	ColumnAutoIncrement = "auto_increment"
)

// Columns contains a slice of column types
type Columns []Column

// Column contains info about one database column retrieved from `SHOW COLUMNS FROM table`
type Column struct {
	Field, Type, Null, Key, Default, Extra dbr.NullString
}

// new idea and use information_schema.columns instead of SHOW COLUMNs query ...
//
//type DataType uint8
//type ColumnKey uint8
//type ColumnExtra uint8
//
//const (
//	DataTypeVarChar DataType = 1 << iota
//	DataTypeChar
//	DataTypeInt
//	DataTypeSmallInt
//	DataTypeDecimal
//)
//
//const (
//	ColumnKeyPRI ColumnKey = 1 << iota
//	ColumnKeyUNI
//	ColumnKeyMUL
//)
//
//const (
//	ColumnExtraAutoIncrement ColumnExtra = 1 << iota
//	ColumnExtraOnUpdateCurrentTimestamp
//)
//
//// Derived from information_schema.columns
//type Column2 struct {
//	Field            string         `db:"COLUMN_NAME"`              //`COLUMN_NAME` varchar(64) NOT NULL DEFAULT '',
//	Pos              int            `db:"ORDINAL_POSITION"`         //`ORDINAL_POSITION` bigint(21) unsigned NOT NULL DEFAULT '0',
//	Default          dbr.NullString `db:"COLUMN_DEFAULT"`           //`COLUMN_DEFAULT` longtext,
//	Null             bool           `db:"IS_NULLABLE"`              //`IS_NULLABLE` varchar(3) NOT NULL DEFAULT '',
//	DataType         DataType       `db:"DATA_TYPE"`                //`DATA_TYPE` varchar(64) NOT NULL DEFAULT '',
//	CharMaxLength    dbr.NullInt64  `db:"CHARACTER_MAXIMUM_LENGTH"` //`CHARACTER_MAXIMUM_LENGTH` bigint(21) unsigned DEFAULT NULL,
//	NumericPrecision dbr.NullInt64  `db:"NUMERIC_PRECISION"`        //`NUMERIC_PRECISION` bigint(21) unsigned DEFAULT NULL,
//	NumericScale     dbr.NullInt64  `db:"NUMERIC_SCALE"`            //`NUMERIC_SCALE` bigint(21) unsigned DEFAULT NULL,
//	Type             string         `db:"COLUMN_TYPE"`              //`COLUMN_TYPE` longtext NOT NULL,
//	ColumnKey        ColumnKey      `db:"COLUMN_KEY"`               //`COLUMN_KEY` varchar(3) NOT NULL DEFAULT '',
//	ColumnExtra      ColumnExtra    `db:"EXTRA"`                    //`EXTRA` varchar(30) NOT NULL DEFAULT '',
//	Comment          string         `db:"COLUMN_COMMENT"`           //`COLUMN_COMMENT` varchar(1024) NOT NULL DEFAULT '',
//}

// GetColumns returns all columns from a table. It discards the column entity_type_id from some
// entity tables.
func GetColumns(dbrSess dbr.SessionRunner, table string) (Columns, error) {
	var cols = make(Columns, 0, 100)

	sel := dbrSess.SelectBySql("SHOW COLUMNS FROM " + dbr.Quoter.QuoteAs(table))

	selSql, selArg, err := sel.ToSql()
	if err != nil {
		return Columns{}, errors.Wrap(err, "[csdb] ToSql")
	}

	rows, err := sel.Query(selSql, selArg...)
	if err != nil {
		return nil, errors.Wrapf(err, "[csdb] Query: %q Args: %#v", selSql, selArg)
	}
	defer rows.Close()

	col := Column{}
	for rows.Next() {
		err := rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra)
		if err != nil {
			return nil, errors.Wrapf(err, "[csdb] Scan Query: %q Args: %#v", selSql, selArg)
		}
		cols = append(cols, col)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "[csdb] rows.Err Query: %q Args: %#v", selSql, selArg)
	}
	return cols, nil
}

// Hash calculates a non-cryptographic, fast and efficient hash value from all columns.
// Current hash algorithm is fnv64.
func (cs Columns) Hash() ([]byte, error) {
	var tr byte = 0x74 // letter t for true
	var buf bytes.Buffer
	for _, c := range cs {
		buf.WriteString(c.Field.String)
		if c.Field.Valid {
			buf.WriteByte(tr)
		}
		buf.WriteString(c.Type.String)
		if c.Type.Valid {
			buf.WriteByte(tr)
		}
		buf.WriteString(c.Null.String)
		if c.Null.Valid {
			buf.WriteByte(tr)
		}
		buf.WriteString(c.Key.String)
		if c.Key.Valid {
			buf.WriteByte(tr)
		}
		buf.WriteString(c.Default.String)
		if c.Default.Valid {
			buf.WriteByte(tr)
		}
		buf.WriteString(c.Extra.String)
		if c.Extra.Valid {
			buf.WriteByte(tr)
		}
	}
	f64 := fnv.New64()
	if _, err := f64.Write(buf.Bytes()); err != nil {
		return nil, err
	}
	buf.Reset()
	ret := buf.Bytes()
	return f64.Sum(ret), nil
}

// Filter returns a new slice filtered by predicate f
func (cs Columns) Filter(f func(Column) bool) (cols Columns) {
	for _, c := range cs {
		if f(c) {
			cols = append(cols, c)
		}
	}
	return
}

// Map will run function f on all items in Columns and returns a copy of the slice.
func (cs Columns) Map(f func(Column) Column) Columns {
	cols := make(Columns, cs.Len())
	for i, c := range cs {
		cols[i] = f(c)
	}
	return cols
}

// FieldNames returns all column names
func (cs Columns) FieldNames() (fieldNames []string) {
	for _, c := range cs {
		if c.Field.Valid {
			fieldNames = append(fieldNames, c.Field.String)
		}
	}
	return
}

// PrimaryKeys returns all primary key columns
func (cs Columns) PrimaryKeys() Columns {
	return cs.Filter(func(c Column) bool {
		return c.IsPK()
	})
}

// UniqueKeys returns all unique key columns
func (cs Columns) UniqueKeys() Columns {
	return cs.Filter(func(c Column) bool {
		return c.IsUnique()
	})
}

// ColumnsNoPK returns all non primary key columns
func (cs Columns) ColumnsNoPK() Columns {
	return cs.Filter(func(c Column) bool {
		return !c.IsPK()
	})
}

// Len returns the length
func (cs Columns) Len() int {
	return len(cs)
}

// ByName finds a column by its name
func (cs Columns) ByName(fieldName string) Column {
	for _, c := range cs {
		if c.Field.Valid && c.Field.String == fieldName {
			return c
		}
	}
	return Column{}
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
	lcs := len(cs)
	for i, c := range cs {
		fmt.Fprintf(&buf, "%#v", c)
		if i+1 < lcs {
			buf.WriteString(",\n")
		}
	}
	return buf.String()
}

// First returns the first column from the Columns slice
func (cs Columns) First() Column {
	if len(cs) > 0 {
		return cs[0]
	}
	return Column{}
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

// Name returns the name of the column, a helper function.
func (c Column) Name() string {
	return c.Field.String
}

// IsPK checks if column is a primary key
func (c Column) IsPK() bool {
	return c.Field.Valid && c.Key.Valid && c.Key.String == ColumnPrimary
}

// IsPK checks if column is a unique key
func (c Column) IsUnique() bool {
	return c.Field.Valid && c.Key.Valid && c.Key.String == ColumnUnique
}

// IsAutoIncrement checks if column has an auto increment property
func (c Column) IsAutoIncrement() bool {
	return c.Field.Valid && c.Extra.Valid && c.Extra.String == ColumnAutoIncrement
}

// IsNull checks if column can have null values
func (c Column) IsNull() bool {
	return c.Field.Valid && c.Null.Valid && c.Null.String == ColumnNull
}

// IsBool checks the name of a column if it contains bool values. Magento uses
// often smallint field types to store bool values.
func (c Column) IsBool() bool {
	if len(c.Field.String) < 3 {
		return false
	}
	return columnTypes.byName.bool.ContainsReverse(c.Field.String)
}

// IsInt checks if a column contains a MySQL int type, independent from its length.
func (c Column) IsInt() bool {
	return strings.Contains(c.Type.String, "int")
}

// IsString checks if a column contains a MySQL varchar or text type.
func (c Column) IsString() bool {
	return columnTypes.byType.string.ContainsReverse(c.Type.String)
}

// IsDate checks if a column contains a MySQL timestamp or date type.
func (c Column) IsDate() bool {
	return columnTypes.byType.dateSW.StartsWithReverse(c.Type.String)
}

// IsFloat checks if a column contains a MySQL decimal or float type.
func (c Column) IsFloat() bool {
	return columnTypes.byType.floatSW.StartsWithReverse(c.Type.String)
}

// IsMoney checks if a column contains a MySQL decimal or float type and the
// column name.
// This function needs a lot of care ...
func (c Column) IsMoney() bool {
	// needs more love
	if false == c.IsFloat() {
		return false
	}
	var ret bool
	switch {
	// could us a long list of || statements but switch looks nicer :-)
	case columnTypes.byName.moneyEqual.Contains(c.Field.String):
		ret = true
	case columnTypes.byName.money.ContainsReverse(c.Field.String):
		ret = true
	case columnTypes.byName.moneySW.StartsWithReverse(c.Field.String):
		ret = true
	case false == c.IsNull() && c.Default.String == "0.0000":
		ret = true
	}
	return ret
}

// GetGoPrimitive detects the Go type of a SQL table column.
func (c Column) GetGoPrimitive(useNullType bool) string {

	var goType = "undefined"
	isNull := c.IsNull() && useNullType
	switch {
	case c.IsBool() && isNull:
		goType = "dbr.NullBool"
	case c.IsBool():
		goType = "bool"
	case c.IsInt() && isNull:
		goType = "dbr.NullInt64"
	case c.IsInt():
		goType = "int64" // rethink if it is worth to introduce uint64 because of some unsigned columns
	case c.IsString() && isNull:
		goType = "dbr.NullString"
	case c.IsString():
		goType = "string"
	case c.IsMoney():
		goType = "money.Money"
	case c.IsFloat() && isNull:
		goType = "dbr.NullFloat64"
	case c.IsFloat():
		goType = "float64"
	case c.IsDate() && isNull:
		goType = "dbr.NullTime"
	case c.IsDate():
		goType = "time.Time"
	}
	return goType
}

// columnTypes looks ugly but ... refactor later
var columnTypes = struct { // the slices in this struct are only for reading. no mutex protection required
	byName struct {
		bool       util.StringSlice
		money      util.StringSlice
		moneySW    util.StringSlice
		moneyEqual util.StringSlice
	}
	byType struct {
		int     util.StringSlice
		string  util.StringSlice
		dateSW  util.StringSlice
		floatSW util.StringSlice
	}
}{
	struct {
		bool       util.StringSlice // contains
		money      util.StringSlice // contains
		moneySW    util.StringSlice // sw == starts with
		moneyEqual util.StringSlice
	}{
		util.StringSlice{"used_", "is_", "has_", "increment_per_store"},
		util.StringSlice{"price", "_tax", "tax_", "_amount", "amount_", "total", "adjustment", "discount"},
		util.StringSlice{"base_", "grand_"},
		util.StringSlice{"value", "price", "cost", "msrp"},
	},
	struct {
		int     util.StringSlice // contains
		string  util.StringSlice // contains
		dateSW  util.StringSlice // SW starts with
		floatSW util.StringSlice // SW starts with
	}{
		util.StringSlice{"int"},
		util.StringSlice{"char", "text"},
		util.StringSlice{"time", "date"},
		util.StringSlice{"decimal", "float", "double"},
	},
}
