// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"database/sql"
	"fmt"
	"strings"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils"
	"github.com/juju/errgo"
	"github.com/kr/pretty"
)

const (
	ColumnPrimary       = "PRI"
	ColumnUnique        = "UNI"
	ColumnNull          = "YES"
	ColumnNotNull       = "NO"
	ColumnAutoIncrement = "auto_increment"
)

type (

	// Columns contains a slice of column types
	Columns []Column
	// Column contains info about one database column retrieved from `SHOW COLUMNS FROM table`
	Column struct {
		Field, Type, Null, Key, Default, Extra sql.NullString
	}
)

// MoneyTypeColumnNames part of the function IsMoney() to detect if a column
// type is a Go money.Currency type.
var MoneyTypeColumnNames = utils.StringSlice{
	"value", "price", "cost", "msrp",
}

// GetColumns returns all columns from a table. It discards the column entity_type_id from some
// entity tables.
func GetColumns(dbrSess dbr.SessionRunner, table string) (Columns, error) {
	var cols = make(Columns, 0, 100)

	sel := dbrSess.SelectBySql("SHOW COLUMNS FROM " + dbr.Quoter.Table(table))
	selSql, selArg := sel.ToSql()
	rows, err := sel.Query(selSql, selArg...)

	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer rows.Close()

	for rows.Next() {
		col := Column{}
		err := rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		cols = append(cols, col)
	}
	err = rows.Err()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return cols, nil
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

// String pretty print
func (cs Columns) String() string {
	// fix tests if you change this layout of the returned string
	var ret = make([]string, len(cs))
	for i, c := range cs {
		ret[i] = fmt.Sprintf("%# v", pretty.Formatter(c))
	}
	return strings.Join(ret, ",\n")
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
	case MoneyTypeColumnNames.Include(c.Field.String):
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
		goType = "money.Currency"
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
		bool    utils.StringSlice
		money   utils.StringSlice
		moneySW utils.StringSlice
	}
	byType struct {
		int     utils.StringSlice
		string  utils.StringSlice
		dateSW  utils.StringSlice
		floatSW utils.StringSlice
	}
}{
	struct {
		bool    utils.StringSlice // contains
		money   utils.StringSlice // contains
		moneySW utils.StringSlice // sw == starts with
	}{
		utils.StringSlice{"used_", "is_", "has_", "increment_per_store"},
		utils.StringSlice{"price", "_tax", "tax_", "_amount", "amount_", "total", "adjustment", "discount"},
		utils.StringSlice{"base_", "grand_"},
	},
	struct {
		int     utils.StringSlice // contains
		string  utils.StringSlice // contains
		dateSW  utils.StringSlice // SW starts with
		floatSW utils.StringSlice // SW starts with
	}{
		utils.StringSlice{"int"},
		utils.StringSlice{"char", "text"},
		utils.StringSlice{"time", "date"},
		utils.StringSlice{"decimal", "float", "double"},
	},
}
