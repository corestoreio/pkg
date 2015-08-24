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
