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
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

// Table represents a table from the database
type Table struct {
	// Name is the table name
	Name string
	// Columns all table columns
	Columns Columns
	// CountPK number of primary keys
	CountPK int
	// CountUnique number of unique keys
	CountUnique int

	// internal caches
	fieldsPK  []string // all PK column field
	fieldsUNI []string // all unique key column field
	fields    []string // all other non-pk column field
}

// NewTable initializes a new table structure
func NewTable(n string, cs ...Column) *Table {
	ts := &Table{
		Name:    n,
		Columns: Columns(cs),
	}
	return ts.update()
}

// update recalculates the internal cached fields
func (ts *Table) update() *Table {
	ts.fieldsPK = ts.Columns.PrimaryKeys().FieldNames()
	ts.fieldsUNI = ts.Columns.UniqueKeys().FieldNames()
	ts.fields = ts.Columns.ColumnsNoPK().FieldNames()
	ts.CountPK = ts.Columns.PrimaryKeys().Len()
	ts.CountUnique = ts.Columns.UniqueKeys().Len()
	return ts
}

// Load reads the column information from the DB. @todo
func (ts *Table) LoadColumns(dbrSess dbr.SessionRunner) (err error) {
	ts.Columns, err = GetColumns(dbrSess, ts.Name)
	ts.update()
	return errgo.Mask(err)
}

// TableAliasQuote returns a table name with the alias.
// catalog_product_entity with alias e would become `catalog_product_entity` AS `e`.
func (ts *Table) TableAliasQuote(alias string) string {
	return dbr.Quoter.QuoteAs(ts.Name, alias)
}

// ColumnAliasQuote prefixes non-id columns with an alias and puts quotes around them. Returns a copy.
func (ts *Table) ColumnAliasQuote(alias string) []string {
	return dbr.Quoter.TableColumnAlias(alias, append([]string(nil), ts.fields...)...)
}

// AllColumnAliasQuote prefixes all columns with an alias and puts quotes around them. Returns a copy.
func (ts *Table) AllColumnAliasQuote(alias string) []string {
	c := append([]string(nil), ts.fieldsPK...)
	return dbr.Quoter.TableColumnAlias(alias, append(c, ts.fields...)...)
}

// In checks if column name n is a column of this table
func (ts *Table) In(n string) bool {
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
func (ts *Table) Select(dbrSess dbr.SessionRunner) (*dbr.SelectBuilder, error) {
	if ts == nil {
		return nil, ErrTableNotFound
	}
	return dbrSess.
		Select(ts.AllColumnAliasQuote(MainTable)...).
		From(ts.Name, MainTable), nil
}

func (ts *Table) Update() {}
func (ts *Table) Delete() {}
func (ts *Table) Insert() {}
func (ts *Table) Alter()  {}
func (ts *Table) Drop()   {}
func (ts *Table) Create() {}
