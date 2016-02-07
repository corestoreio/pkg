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
	"errors"
	"sync"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util"
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
	ErrTableNotFound          = errors.New("Table not found")
	ErrTableServiceInitReload = errors.New("You cannot force reload when the init process were not able to run.")
)

// Index defines the table index within the TableService
type Index uint

type (
	TableManager interface {
		// Structure returns the TableStructure from a read-only map m by a giving index i.
		Structure(Index) (*Table, error)
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
		// Len returns the length of the underlying map
		Len() Index
		// Append adds a table. Overrides silently existing entries.
		Append(Index, *Table) error
		// Init reloads the internal table structs from the database
		Init(dbrSess dbr.SessionRunner, reInit ...bool) error
	}

	// ManagerOption applies options to the TableService
	ManagerOption func(*TableService)

	// TableService implements interface Manager
	TableService struct {
		initDone bool
		errs     []error
		mu       sync.RWMutex
		ts       map[Index]*Table
	}
)

// WithTable adds a database table to the TableService by the table name and index.
// You can optionally specify the columns to skip the Init() function.
func WithTable(idx Index, name string, cols ...Column) ManagerOption {
	return func(tm *TableService) {

		if err := IsValidIdentifier(name); err != nil {
			tm.errs = append(tm.errs, err)
		}

		if len(tm.errs) > 0 {
			return
		}

		if err := tm.Append(idx, NewTable(name, cols...)); err != nil {
			tm.errs = append(tm.errs, err)
		}
	}
}

// NewTableService creates a new TableService satisfying interface Manager.
func NewTableService(opts ...ManagerOption) (*TableService, error) {
	tm := &TableService{
		mu: sync.RWMutex{},
		ts: make(map[Index]*Table),
	}
	for _, o := range opts {
		o(tm)
	}
	if len(tm.errs) > 0 {
		return nil, tm
	}
	return tm, nil
}

// MustNewTableService same as NewTableService but panics on error.
func MustNewTableService(opts ...ManagerOption) *TableService {
	ts, err := NewTableService(opts...)
	if err != nil {
		panic(err)
	}
	return ts
}

// Structure returns the TableStructure from a read-only map m by a giving index i.
func (tm *TableService) Structure(i Index) (*Table, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if ts, ok := tm.ts[i]; ok && ts != nil {
		return ts, nil
	}
	return nil, ErrTableNotFound
}

// Name is a short hand to return a table name by given index i. Does not return an error
// when the table can't be found. Returns an empty string
func (tm *TableService) Name(i Index) string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if ts, ok := tm.ts[i]; ok && ts != nil {
		return ts.Name
	}
	return ""
}

// Len returns the length of the slice data
func (tm *TableService) Len() Index {
	return Index(len(tm.ts))
}

// Next iterator function where i is the current index starting with zero.
// Example:
//	for i := Index(0); tableMap.Next(i); i++ {
//		table, err := tableMap.Structure(i)
//		...
//	}
func (tm *TableService) Next(i Index) bool {
	return i < tm.Len()
}

// Append adds a table. Overrides silently existing entries.
func (tm *TableService) Append(i Index, ts *Table) error {
	if ts == nil {
		return errgo.Newf("Table pointer cannot be nil for Index %d", i)
	}
	tm.mu.Lock()
	tm.ts[i] = ts
	tm.mu.Unlock() // use defer once there are multiple returns
	return nil
}

// Init loads the column definitions from the database for each table. Set reInit
// to true to allow reloading otherwise it loads only once.
func (tm *TableService) Init(dbrSess dbr.SessionRunner, reInit ...bool) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	reLoad := false
	if len(reInit) > 0 {
		reLoad = reInit[0]
	}
	if true == tm.initDone && false == reLoad {
		return nil
	}
	if false == tm.initDone && true == reLoad {
		return errgo.Mask(ErrTableServiceInitReload)
	}
	tm.initDone = true

	for _, table := range tm.ts {
		if err := table.LoadColumns(dbrSess); err != nil {
			if PkgLog.IsDebug() {
				PkgLog.Debug("csdb.TableService.Init.LoadColumns", "err", err, "table", table)
			}
			return errgo.Mask(err)
		}
	}

	return nil
}

// Error implements error interface
func (tm *TableService) Error() string {
	return util.Errors(tm.errs...)
}
