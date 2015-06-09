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
	"errors"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

const (
	// MainTable used in SQL to refer to the main table as an alias
	MainTable = "main_table"
	// AddTable additional table
	AdditionalTable = "additional_table"
	// ScopeTable used in SQl to refer to a website scope table as an alias
	ScopeTable = "scope_table"
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
		// IDFieldNames contains only primary keys
		IDFieldNames []string
		// Columns all other columns which are not primary keys
		Columns []string
	}

	DbrSelectCb func(*dbr.SelectBuilder) *dbr.SelectBuilder
)

var _ TableStructurer = (*TableStructureSlice)(nil)

func NewTableStructure(n string, IDs, c []string) *TableStructure {
	return &TableStructure{
		Name:         n,
		IDFieldNames: IDs,
		Columns:      c,
	}
}

// remove this once the ALIAS via []string is implemented in DBR
func (ts *TableStructure) TableAliasQuote(alias string) string {
	return "`" + ts.Name + "` AS `" + alias + "`"
}

// ColumnAliasQuote prefixes non-id columns with an alias and puts quotes around them. Returns a copy.
func (ts *TableStructure) ColumnAliasQuote(alias string) []string {
	return dbr.TableColumnQuote(alias, append([]string(nil), ts.Columns...)...)
}

// AllColumnAliasQuote prefixes all columns with an alias and puts quotes around them. Returns a copy.
func (ts *TableStructure) AllColumnAliasQuote(alias string) []string {
	c := append([]string(nil), ts.IDFieldNames...)
	return dbr.TableColumnQuote(alias, append(c, ts.Columns...)...)
}

// In checks if column name n is a column of this table
func (ts *TableStructure) In(n string) bool {
	for _, c := range ts.IDFieldNames {
		if c == n {
			return true
		}
	}
	for _, c := range ts.Columns {
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
