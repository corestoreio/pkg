// Copyright 2015 CoreStore Authors
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
)

var (
	ErrTableNotFound = errors.New("Table not found")
)

type (
	Index               int
	TableStructureSlice []*TableStructure

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

func (ts *TableStructure) ColumnAliasQuote(alias string) []string {
	ret := make([]string, len(ts.Columns))
	for i, c := range ts.Columns {
		ret[i] = "`" + alias + "`.`" + c + "`"
	}
	return ret
}

func (ts *TableStructure) AllColumnAliasQuote(alias string) []string {
	ret := make([]string, len(ts.IDFieldNames), len(ts.Columns))
	for i, c := range ts.IDFieldNames {
		ret[i] = "`" + alias + "`.`" + c + "`"
	}
	return append(ret, ts.ColumnAliasQuote(alias)...)
}

// Structure returns the TableStructure from a read-only map m by a giving index i.
func (m TableStructureSlice) Structure(i Index) (*TableStructure, error) {
	return m[i], nil
}

// Name is a short hand to return a table name by given index i. Does not return an error
// when the table can't be found.
func (m TableStructureSlice) Name(i Index) string {
	return m[i].Name
}
