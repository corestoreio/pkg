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
	"sort"
	"strings"
	"sync"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
)

// @deprecated
const (
	MainTable       = "main_table"
	AdditionalTable = "additional_table"
	ScopeTable      = "scope_table"
)

// TableOption applies options and helper functions when creating a new table.
// For example loading column definitions.
type TableOption struct {
	// priority takes care that the options gets applied in the correct order.
	// e.g. column loading can only happen when a table is present.
	priority uint8
	fn       func(*Tables) error
}

// Tables handles all the tables defined for a package. Thread safe.
type Tables struct {
	// Schema represents the name of the database. Might be empty.
	Schema string
	mu     sync.RWMutex
	// ts uses int as the table index.
	// What is the reason to use int as the table index and not a name? Because
	// table names between M1 and M2 get renamed and in a Go SQL code generator
	// script of the CoreStore project, we can guarantee that the generated
	// index constant will always stay the same but the name of the table
	// differs.
	ts map[int]*Table
	// tn for faster access we use tn and also because ts might get removed
	tn map[string]*Table
}

// WithTableOrViewFromQuery creates the new view or table from the SELECT query and
// adds it to the internal table manager including all loaded column
// definitions. If providing true in the argument "dropIfExists" the view or
// table gets first dropped, if exists, and then created. Argument typ can be
// only `table` or `view`.
func WithTableOrViewFromQuery(ctx context.Context, db interface {
	dbr.Execer
	dbr.Querier
}, typ string, idx int, objectName string, query string, dropIfExists ...bool) TableOption {
	return TableOption{
		priority: 10,
		fn: func(tm *Tables) error {

			if err := IsValidIdentifier(objectName); err != nil {
				return errors.Wrapf(err, "[csdb] WithTableOrViewFromQuery.IsValidIdentifier")
			}

			var viewOrTable string
			switch typ {
			case "view":
				viewOrTable = "VIEW"
			case "table":
				viewOrTable = "TABLE"
			default:
				return errors.NewUnavailablef("[csdb] Option %q for variable typ not available. Only `view` or `table`", typ)
			}

			vnq := dbr.Quoter.Name(objectName)
			if len(dropIfExists) > 0 && dropIfExists[0] {
				if _, err := db.ExecContext(ctx, "DROP "+viewOrTable+" IF EXISTS "+vnq); err != nil {
					return errors.Wrapf(err, "[csdb] Drop view failed %q", objectName)
				}
			}

			_, err := db.ExecContext(ctx, "CREATE "+viewOrTable+" "+vnq+" AS "+query)
			if err != nil {
				return errors.Wrapf(err, "[csdb] Create view %q failed", objectName)
			}

			tc, err := LoadColumns(ctx, db, objectName)
			if err != nil {
				return errors.Wrapf(err, "[csdb] Load columns failed for %q", objectName)
			}

			if err := WithTable(idx, objectName, tc[objectName]...).fn(tm); err != nil {
				return errors.Wrapf(err, "[csdb] Failed to add new table %q", objectName)
			}

			tm.mu.Lock()
			defer tm.mu.Unlock()
			tm.ts[idx].IsView = viewOrTable == "VIEW"

			return nil
		},
	}
}

// WithTable inserts a new table to the Tables struct, identified by its index.
// You can optionally specify the columns. What is the reason to use int as the
// table index and not a name? Because table names between M1 and M2 get renamed
// and in a Go SQL code generator script of the CoreStore project, we can
// guarantee that the generated index constant will always stay the same but the
// name of the table differs.
func WithTable(idx int, tableName string, cols ...*Column) TableOption {
	return TableOption{
		fn: func(tm *Tables) error {
			if err := IsValidIdentifier(tableName); err != nil {
				return errors.Wrap(err, "[csdb] WithNewTable.IsValidIdentifier")
			}

			if err := tm.Upsert(idx, NewTable(tableName, cols...)); err != nil {
				return errors.Wrap(err, "[csdb] WithNewTable.Tables.Insert")
			}
			return nil
		},
	}
}

// WithTableLoadColumns inserts a new table to the Tables struct, identified by
// its index. What is the reason to use int as the table index and not a name?
// Because table names between M1 and M2 get renamed and in a Go SQL code
// generator script of the CoreStore project, we can guarantee that the
// generated index constant will always stay the same but the name of the table
// differs.
func WithTableLoadColumns(ctx context.Context, db dbr.Querier, idx int, tableName string) TableOption {
	return TableOption{
		fn: func(tm *Tables) error {
			if err := IsValidIdentifier(tableName); err != nil {
				return errors.Wrap(err, "[csdb] WithTableLoadColumns.IsValidIdentifier")
			}

			t := NewTable(tableName)
			t.Schema = tm.Schema
			if err := t.LoadColumns(ctx, db); err != nil {
				return errors.Wrap(err, "[csdb] WithTableLoadColumns.LoadColumns")
			}

			if err := tm.Upsert(idx, t); err != nil {
				return errors.Wrap(err, "[csdb] Tables.Insert")
			}
			return nil
		},
	}
}

// WithTableNames creates for each table name and its index a new table pointer.
// You should call afterwards the functional option WithLoadColumnDefinitions.
// This function returns an error if a table index already exists.
func WithTableNames(idx []int, tableName []string) TableOption {
	return TableOption{
		fn: func(tm *Tables) error {
			if len(idx) != len(tableName) {
				return errors.NewNotValidf("[csdb] Length of the index must be equal to the length of the table names: %d != %d", len(idx), len(tableName))
			}

			if err := IsValidIdentifier(tableName...); err != nil {
				return errors.Wrap(err, "[csdb] WithTable.IsValidIdentifier")
			}

			for i, tn := range tableName {
				if err := tm.Upsert(idx[i], NewTable(tn)); err != nil {
					return errors.Wrapf(err, "[csdb] Tables.Insert %q", tn)
				}
			}
			return nil
		},
	}
}

// WithLoadTableNames executes a query to load all available tables in the
// current database. Argument sql will be either appended to the SHOW TABLES
// statement or if it starts with SELECT then it replaces the SHOW TABLES
// statement.
func WithLoadTableNames(querier dbr.Querier, sql ...string) TableOption {
	qry := "SHOW TABLES"
	if len(sql) > 0 && sql[0] != "" {
		if false == dbr.Stmt.IsSelect(sql[0]) {
			qry = qry + " LIKE '" + strings.Replace(sql[0], "'", "", -1) + "'"
		} else {
			qry = sql[0]
		}
	}
	return TableOption{
		fn: func(tm *Tables) error {
			rows, err := querier.QueryContext(context.Background(), qry)
			if err != nil {
				return errors.Wrapf(err, "[csdb] Query %q failed", qry)
			}
			var tableName string

			i := 0
			for rows.Next() {
				if err := rows.Scan(&tableName); err != nil {
					return errors.Wrapf(err, "Scan Query %q", qry)
				}
				if err := tm.Upsert(i, NewTable(tableName)); err != nil {
					return errors.Wrapf(err, "[csdb] Tables.Insert Index %d with name %q", i, tableName)
				}
				i++
			}

			if err = rows.Err(); err != nil {
				return errors.Wrapf(err, "[csdb] Rows with query %q", qry)
			}
			return nil
		},
	}
}

// WithLoadColumnDefinitions loads the column definitions from the database for each
// table in the internal map. Thread safe.
func WithLoadColumnDefinitions(ctx context.Context, db dbr.Querier) TableOption {
	return TableOption{
		priority: 255, // must be the last element
		fn: func(tm *Tables) error {

			tc, err := LoadColumns(ctx, db, tm.Tables()...)
			if err != nil {
				return errors.Wrap(err, "[csdb] table.LoadColumns")
			}

			tm.mu.Lock()
			defer tm.mu.Unlock()
			for _, t := range tm.ts {
				if c, ok := tc[t.Name]; ok {
					t.Columns = c
					t.update()
				}
			}

			return nil
		},
	}
}

// WithTableDMLListeners adds event listeners to a table object. It doesn't
// matter if the table has already been set. If the table object gets set later,
// the events will be copied to the new object.
func WithTableDMLListeners(idx int, events ...*dbr.ListenerBucket) TableOption {
	return TableOption{
		priority: 254,
		fn: func(tm *Tables) error {
			tm.mu.Lock()
			defer tm.mu.Unlock()

			t, ok := tm.ts[idx]
			if !ok {
				return errors.NewNotFoundf("[csdb] Table at index %d not found", idx)
			}
			t.Listeners.Merge(events...)
			tm.ts[idx] = t

			return nil
		},
	}
}

// NewTables creates a new TableService satisfying interface Manager.
func NewTables(opts ...TableOption) (*Tables, error) {
	tm := &Tables{
		ts: make(map[int]*Table),
	}
	if err := tm.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[csdb] NewTables applied option error")
	}
	return tm, nil
}

// MustNewTables same as NewTableService but panics on error.
func MustNewTables(opts ...TableOption) *Tables {
	ts, err := NewTables(opts...)
	if err != nil {
		panic(err)
	}
	return ts
}

// MustInitTables helper function in init() statements to initialize the global
// table collection variable independent of knowing when this variable is nil.
// We cannot assume the correct order, how all init() invocations are executed,
// at least they don't run in parallel during packet initialization. Yes ... bad
// practice to rely on init ... but for now it works very well.
//
//		func init() {
//			TableCollection = csdb.MustInitTables(TableCollection,[Options])
// 		}
// TODO(CyS) rethink and refactor maybe.
func MustInitTables(ts *Tables, opts ...TableOption) *Tables {
	if ts == nil {
		var err error
		ts, err = NewTables()
		if err != nil {
			panic(err)
		}
	}
	if err := ts.Options(opts...); err != nil {
		panic(err)
	}
	return ts
}

// Options applies options to the Tables service.
func (tm *Tables) Options(opts ...TableOption) error {

	// SliceStable must be stable to maintain the order of all options where
	// priority is zero.
	sort.SliceStable(opts, func(i, j int) bool {
		return opts[i].priority < opts[j].priority
	})

	for _, to := range opts {
		if err := to.fn(tm); err != nil {
			return errors.Wrap(err, "[csdb] Applied option error")
		}
	}
	return nil
}

// Table returns the structure from a map m by a giving index i. What is the
// reason to use int as the table index and not a name? Because table names
// between M1 and M2 get renamed and in a Go SQL code generator script of the
// CoreStore project, we can guarantee that the generated index constant will
// always stay the same but the name of the table differs.
func (tm *Tables) Table(i int) (*Table, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if t, ok := tm.ts[i]; ok {
		return t, nil
	}
	return nil, errors.NewNotFoundf("[csdb] Table at index %d not found.", i)
}

// TableByName returns a table object via its table name. Case sensitive.
func (tm *Tables) TableByName(name string) (*Table, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	for _, t := range tm.ts {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, errors.NewNotFoundf("[csdb] Table %q not found.", name)
}

// Tables returns a list of all available table names.
func (tm *Tables) Tables() []string {
	// todo maybe use internal cache
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	ts := make([]string, 0, len(tm.ts))
	for _, table := range tm.ts {
		ts = append(ts, table.Name)
	}
	return ts
}

// MustTable same as Table function but panics when the table cannot be found or
// any other error occurs.
func (tm *Tables) MustTable(i int) *Table {
	t, err := tm.Table(i)
	if err != nil {
		panic(err)
	}
	return t
}

// Name is a short hand to return a table name by given index i. Does not return
// an error when the table can't be found but returns an empty string.
func (tm *Tables) Name(i int) string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if ts, ok := tm.ts[i]; ok && ts != nil {
		return ts.Name
	}
	return ""
}

// Len returns the number of all tables.
func (tm *Tables) Len() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return len(tm.ts)
}

// Upsert adds or updates a new table into the internal cache. If a table
// already exists, then the new table gets applied. The ListenerBuckets gets
// merged from the existing table to the new table, they will be appended to the
// new table buckets. Empty fields in the new table gets updated from the
// existing table.
func (tm *Tables) Upsert(i int, tNew *Table) error {
	_ = tNew.Name // let it panic as early as possible if *Table is nil

	tm.mu.Lock()
	defer tm.mu.Unlock()

	tOld, ok := tm.ts[i]
	if tOld == nil || !ok {
		tm.ts[i] = tNew
		return nil
	}

	// for now copy only the events from the existing table
	tNew.Listeners.Merge(&tOld.Listeners)

	if tNew.Schema == "" {
		tNew.Schema = tOld.Schema
	}
	if tNew.Name == "" {
		tNew.Name = tOld.Name
	}
	if len(tNew.Columns) == 0 {
		tNew.Columns = tOld.Columns
	}

	tm.ts[i] = tNew.update()
	return nil
}

// DeleteFromCache removes tables by their given indexes. If no index has been passed
// then all entries get removed and the map reinitialized.
func (tm *Tables) DeleteFromCache(idxs ...int) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	for _, idx := range idxs {
		delete(tm.ts, idx)
	}
}

// DeleteAllFromCache clears the internal table cache and resets the map.
func (tm *Tables) DeleteAllFromCache() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// maybe clear each pointer in the Table struct to avoid a memory leak
	tm.ts = make(map[int]*Table)
}
