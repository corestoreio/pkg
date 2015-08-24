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
	"sync"

	"bytes"
	"fmt"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils/log"
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
	ErrTableNotFound         = errors.New("Table not found")
	ErrManagerIncorrectValue = errors.New("NewTableManager: Incorrect value for idx or name")
	ErrManagerInitReload     = errors.New("You cannot force reload when the init process were not able to run.")
)

type (
	Index int

	Manager interface {
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

	// ManagerOption applies options to the TableManager
	ManagerOption func(*TableManager)

	// TableManager implements interface Manager
	TableManager struct {
		initDone bool
		errs     []error
		mu       sync.RWMutex
		ts       map[Index]*Table
	}

	DbrSelectCb func(*dbr.SelectBuilder) *dbr.SelectBuilder
)

var _ Manager = (*TableManager)(nil)

// AddTableByName adds a database table to the TableManager by the table name
// and index.
func AddTableByName(idx Index, name string) ManagerOption {
	return func(tm *TableManager) {
		if idx < Index(0) || name == "" {
			tm.appendErr(errgo.Mask(ErrManagerIncorrectValue))
			return
		}
		if err := tm.Append(idx, NewTable(name)); err != nil {
			tm.appendErr(err)
		}
	}
}

// NewTableManager creates a new TableManager satisfying interface Manager.
// Panics if an error occurs in an option function. Errors will be logged.
// This function is only used in generated codes.
func NewTableManager(opts ...ManagerOption) *TableManager {
	tm := &TableManager{
		mu: sync.RWMutex{},
		ts: make(map[Index]*Table),
	}
	for _, o := range opts {
		o(tm)
	}
	if len(tm.errs) > 0 {
		panic(tm.errFlush())
	}
	return tm
}

// Structure returns the TableStructure from a read-only map m by a giving index i.
func (tm *TableManager) Structure(i Index) (*Table, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if ts, ok := tm.ts[i]; ok && ts != nil {
		return ts, nil
	}
	return nil, ErrTableNotFound
}

// Name is a short hand to return a table name by given index i. Does not return an error
// when the table can't be found. Returns an empty string
func (tm *TableManager) Name(i Index) string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if ts, ok := tm.ts[i]; ok && ts != nil {
		return ts.Name
	}
	return ""
}

// Len returns the length of the slice data
func (tm *TableManager) Len() Index {
	return Index(len(tm.ts))
}

// Next iterator function where i is the current index starting with zero.
// Example:
//	for i := Index(0); tableMap.Next(i); i++ {
//		table, err := tableMap.Structure(i)
//		...
//	}
func (tm *TableManager) Next(i Index) bool {
	return i < tm.Len()
}

// Append adds a table. Overrides silently existing entries.
func (tm *TableManager) Append(i Index, ts *Table) error {
	if ts == nil {
		return log.Error("csdb.TableManager.Init", "err", errgo.Newf("Table pointer cannot be nil for Index %d", i))
	}
	tm.mu.Lock()
	tm.ts[i] = ts
	tm.mu.Unlock() // use defer once there are multiple returns
	return nil
}

// Init loads the column definitions from the database for each table. Set reInit
// to true to allow reloading otherwise it loads only once.
func (tm *TableManager) Init(dbrSess dbr.SessionRunner, reInit ...bool) error {
	reLoad := false
	if len(reInit) > 0 {
		reLoad = reInit[0]
	}
	if true == tm.initDone && false == reLoad {
		return nil
	}
	if false == tm.initDone && true == reLoad {
		return log.Error("csdb.TableManager.Init", "err", ErrManagerInitReload)
	}
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.initDone = true

	for _, table := range tm.ts {
		if err := table.LoadColumns(dbrSess); err != nil {
			return log.Error("csdb.TableManager.Init.LoadColumns", "err", err, "table", table)
		}
	}

	return nil
}

func (tm *TableManager) appendErr(err error) {
	tm.errs = append(tm.errs, err)
}

func (tm *TableManager) errFlush() string {
	var buf bytes.Buffer
	for i, e := range tm.errs {
		log.Error("csdb.NewTableManager.errs", "err", e, "tablemanager", tm)
		fmt.Fprintf(&buf, "%02d: %s", i, e.Error())
		if l, ok := e.(errgo.Locationer); ok {
			buf.WriteString("\nLocation: " + l.Location().String())
		}
	}
	tm.errs = nil
	return buf.String()
}
