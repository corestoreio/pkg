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
	"context"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/errors"
)

// Table represents a table from a database.
type Table struct {
	// Name of the table
	Name string
	// Columns all table columns
	Columns Columns
	// CountPK number of primary keys. Auto updated.
	CountPK int
	// CountUnique number of unique keys. Auto updated.
	CountUnique int

	// internal caches
	fieldsPK  []string // all PK column field
	fieldsUNI []string // all unique key column field
	fields    []string // all other non-pk column field
}

// NewTable initializes a new table structure
func NewTable(tableName string, cs ...*Column) *Table {
	ts := &Table{
		Name:    tableName,
		Columns: Columns(cs),
	}
	return ts.update()
}

// update recalculates the internal cached fields
func (ts *Table) update() *Table {
	if len(ts.Columns) == 0 {
		return ts
	}
	ts.fieldsPK = ts.Columns.PrimaryKeys().FieldNames()
	ts.fieldsUNI = ts.Columns.UniqueKeys().FieldNames()
	ts.fields = ts.Columns.ColumnsNoPK().FieldNames()
	ts.CountPK = ts.Columns.PrimaryKeys().Len()
	ts.CountUnique = ts.Columns.UniqueKeys().Len()
	return ts
}

// LoadColumns reads the column information from the DB.
func (ts *Table) LoadColumns(ctx context.Context, db Querier) (err error) {
	ts.Columns, err = LoadColumns(ctx, db, ts.Name)
	ts.update()
	return errors.Wrapf(err, "[csdb] table.LoadColumns. Table %q", ts.Name)
}

// TableAliasQuote returns a table name with the alias. catalog_product_entity
// with alias e would become `catalog_product_entity` AS `e`.
func (ts *Table) TableAliasQuote(alias string) string {
	return dbr.Quoter.QuoteAs(ts.Name, alias)
}

// ColumnAliasQuote prefixes non-id columns with an alias and puts quotes around
// them. Returns a copy.
func (ts *Table) ColumnAliasQuote(alias string) []string {
	sl := make([]string, len(ts.fields))
	copy(sl, ts.fields)
	return dbr.Quoter.TableColumnAlias(alias, sl...)
}

// AllColumnAliasQuote prefixes all columns with an alias and puts quotes around
// them. Returns a copy.
func (ts *Table) AllColumnAliasQuote(alias string) []string {
	sl := make([]string, len(ts.fieldsPK)+len(ts.fields))
	n := copy(sl, ts.fieldsPK)
	copy(sl[n:], ts.fields)
	return dbr.Quoter.TableColumnAlias(alias, sl...)
}

// In checks if column name n is a column of this table. Case sensitive.
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
func (t *Table) Select(dbrSess dbr.SessionRunner) (*dbr.SelectBuilder, error) {
	if t == nil {
		return nil, errors.NewFatalf("[csdb] Table cannot be nil")
	}
	return dbrSess.
		Select(t.AllColumnAliasQuote(MainTable)...).
		From(t.Name, MainTable), nil
}

// LoadSlice performs a SELECT * FROM `tableName` query and puts the results
// into the pointer slice `dest`. Returns the number of loaded rows and nil or 0
// and an error. The variadic thrid arguments can modify the SQL query.
func (t *Table) LoadSlice(dbrSess dbr.SessionRunner, dest interface{}, cbs ...dbr.SelectCb) (int, error) {
	sb, err := t.Select(dbrSess)
	if err != nil {
		return 0, errors.Wrap(err, "[csdb] LoadSlice.Select")
	}

	for _, cb := range cbs {
		sb = cb(sb)
	}
	return sb.LoadStructs(dest)
}

//func (ts *Table) Update() {}
//func (ts *Table) Delete() {}
//func (ts *Table) Insert() {}
//func (ts *Table) Alter()  {}
//func (ts *Table) Drop()   {}
//func (ts *Table) Create() {}
