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
	"errors"
	"fmt"

	"strings"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
	"github.com/kr/pretty"
)

const (
	// MainTable used in SQL to refer to the main table as an alias
	MainTable = "main_table"
	// AddTable additional table
	AdditionalTable = "additional_table"
	// ScopeTable used in SQl to refer to a website scope table as an alias
	ScopeTable = "scope_table"

	ColumnPrimary       = "PRI"
	ColumnUnique        = "UNI"
	ColumnNull          = "YES"
	ColumnNotNull       = "NO"
	ColumnAutoIncrement = "auto_increment"
)

var (
	ErrTableNotFound = errors.New("Table not found")
)

type (
	Index int
	// TableStructureSlice implements interface TableStructurer
	TableStructureSlice []*TableStructure

	TableStructurer interface {
		// Structure returns the TableStructure from a read-only map m by a giving index i.
		Structure(Index) (*TableStructure, error)
		// Name is a short hand to return a table name by given index i. Does not return an error
		// when the table can't be found.
		Name(Index) string

		// Next iterator function where i is the current index starting with zero.
		// Example:
		//	for i := Index(0); tableMap.Next(i); i++ {
		//		table, err := tableMap.Structure(i)
		//		...
		//	}
		Next(Index) bool
		// Len returns the length of the underlying slice
		Len() Index
	}

	// temporary place
	TableStructure struct {
		// Name is the table name
		Name string
		// Columns all table columns
		Columns Columns
		// CountPK number of primary keys
		CountPK int
		// CountUni number of unique keys
		CountUni int

		// internal caches
		fieldsPK  []string // all PK column field
		fieldsUNI []string // all unique key column field
		fields    []string // all other non-pk column field
	}

	DbrSelectCb func(*dbr.SelectBuilder) *dbr.SelectBuilder

	// Columns contains a slice of column types
	Columns []Column
	// Column contains info about one database column retrieved from `SHOW COLUMNS FROM table`
	Column struct {
		Field, Type, Null, Key, Default, Extra sql.NullString
	}
)

// Load reads the column information from the DB. @todo
func (cs Columns) Load(dbrSess dbr.Session, tableName string) (cols Columns) {
	// @todo
	return
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

var _ TableStructurer = (*TableStructureSlice)(nil)

// NewTableStructure initializes a new table structure
func NewTableStructure(n string, cs ...Column) *TableStructure {
	ts := &TableStructure{
		Name:    n,
		Columns: Columns(cs),
	}
	ts.fieldsPK = ts.Columns.PrimaryKeys().FieldNames()
	ts.fieldsUNI = ts.Columns.UniqueKeys().FieldNames()
	ts.fields = ts.Columns.ColumnsNoPK().FieldNames()
	ts.CountPK = ts.Columns.PrimaryKeys().Len()
	ts.CountUni = ts.Columns.UniqueKeys().Len()

	return ts
}

// remove this once the ALIAS via []string is implemented in DBR
func (ts *TableStructure) TableAliasQuote(alias string) string {
	return "`" + ts.Name + "` AS `" + alias + "`"
}

// ColumnAliasQuote prefixes non-id columns with an alias and puts quotes around them. Returns a copy.
func (ts *TableStructure) ColumnAliasQuote(alias string) []string {
	return dbr.TableColumnQuote(alias, append([]string(nil), ts.fields...)...)
}

// AllColumnAliasQuote prefixes all columns with an alias and puts quotes around them. Returns a copy.
func (ts *TableStructure) AllColumnAliasQuote(alias string) []string {
	c := append([]string(nil), ts.fieldsPK...)
	return dbr.TableColumnQuote(alias, append(c, ts.fields...)...)
}

// In checks if column name n is a column of this table
func (ts *TableStructure) In(n string) bool {
	for _, c := range ts.fieldsPK {
		if c == n {
			return true
		}
	}
	for _, c := range ts.fields {
		if c == n {
			return true
		}
	}
	return false
}

// Select generates a SELECT * FROM tableName statement
func (ts *TableStructure) Select(dbrSess dbr.SessionRunner) (*dbr.SelectBuilder, error) {
	if ts == nil {
		return nil, ErrTableNotFound
	}
	return dbrSess.
		Select(ts.AllColumnAliasQuote("main_table")...).
		From(ts.Name, "main_table"), nil
}

// Structure returns the TableStructure from a read-only map m by a giving index i.
func (m TableStructureSlice) Structure(i Index) (*TableStructure, error) {
	if i < m.Len() {
		return m[i], nil
	}
	return nil, ErrTableNotFound
}

// Name is a short hand to return a table name by given index i. Does not return an error
// when the table can't be found.
func (m TableStructureSlice) Name(i Index) string {
	if i < m.Len() {
		return m[i].Name
	}
	return ""
}

// Len returns the length of the slice data
func (m TableStructureSlice) Len() Index {
	return Index(len(m))
}

// Next iterator function where i is the current index starting with zero.
// Example:
//	for i := Index(0); tableMap.Next(i); i++ {
//		table, err := tableMap.Structure(i)
//		...
//	}
func (m TableStructureSlice) Next(i Index) bool {
	return i < m.Len()
}

// LoadSlice loads the slice dest with the table structure from tsr TableStructurer and table index ti.
// Returns the number of loaded rows and nil or 0 and an error. Slice must be a pointer to structs.
func LoadSlice(dbrSess dbr.SessionRunner, tsr TableStructurer, ti Index, dest interface{}, cbs ...DbrSelectCb) (int, error) {
	ts, err := tsr.Structure(ti)
	if err != nil {
		return 0, errgo.Mask(err)
	}

	sb, err := ts.Select(dbrSess)
	if err != nil {
		return 0, errgo.Mask(err)
	}

	for _, cb := range cbs {
		sb = cb(sb)
	}
	return sb.LoadStructs(dest)
}
