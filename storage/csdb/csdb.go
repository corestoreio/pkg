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
	"github.com/juju/errgo"
)

var (
	ErrTableNotFound = errors.New("Table not found")
)

type (
	Index    int
	TableMap map[Index]*TableStructure

	// temporary place
	TableStructure struct {
		Name         string
		IDFieldNames []string
		Columns      []string
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

func (ts *TableStructure) ColumnAliasQuote(alias string) []string {
	ret := make([]string, len(ts.Columns))
	for i, c := range ts.Columns {
		ret[i] = "`" + alias + "`.`" + c + "`"
	}
	return ret
}

// Structure returns the TableStructure from a read-only map m by a giving index i.
func (m TableMap) Structure(i Index) (*TableStructure, error) {
	if t, ok := m[i]; ok {
		return t, nil
	}
	return nil, errgo.Mask(ErrTableNotFound)
}

// Name is a short hand to return a table name by given index i. Does not return an error
// when the table can't be found.
func (m TableMap) Name(i Index) string {
	if t, ok := m[i]; ok {
		return t.Name
	}
	return ""
}
