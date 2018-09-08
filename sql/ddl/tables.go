// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package ddl

import (
	"context"
	"sort"
	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
)

const (
	MainTable       = "main_table"
	AdditionalTable = "additional_table"
	ScopeTable      = "scope_table"
)

// TableOption applies options and helper functions when creating a new table.
// For example loading column definitions.
type TableOption struct {
	// sortOrder takes care that the options gets applied in the correct order.
	// e.g. column loading can only happen when a table is present.
	sortOrder uint8
	fn        func(*Tables) error
}

// Tables handles all the tables defined for a package. Thread safe.
type Tables struct {
	DB dml.QueryExecPreparer
	// Schema represents the name of the database. Might be empty.
	Schema        string
	previousTable string // the table which has been scanned beforehand
	mu            sync.RWMutex
	// tm a map where key = table name and value the table pointer
	tm map[string]*Table
}

// WithDB sets the DB object to the Tables and all sub Table types to handle the
// database connections. It must be set if other options get used to access the
// DB.
func WithDB(db dml.QueryExecPreparer) TableOption {
	return TableOption{
		sortOrder: 2,
		fn: func(tm *Tables) error {
			tm.DB = db
			tm.mu.Lock()
			defer tm.mu.Unlock()
			for _, t := range tm.tm {
				t.DB = db
			}
			return nil
		},
	}
}

// WithTableOrViewFromQuery creates the new view or table from the SELECT query and
// adds it to the internal table manager including all loaded column
// definitions. If providing true in the argument "dropIfExists" the view or
// table gets first dropped, if exists, and then created. Argument typ can be
// only `table` or `view`.
func WithTableOrViewFromQuery(ctx context.Context, db dml.QueryExecPreparer, typ string, objectName string, query string, dropIfExists ...bool) TableOption {
	return TableOption{
		sortOrder: 10,
		fn: func(tm *Tables) error {

			if err := dml.IsValidIdentifier(objectName); err != nil {
				return errors.WithStack(err)
			}

			var viewOrTable string
			switch typ {
			case "view":
				viewOrTable = "VIEW"
			case "table":
				viewOrTable = "TABLE"
			default:
				return errors.Unavailable.Newf("[ddl] Option %q for variable typ not available. Only `view` or `table`", typ)
			}

			vnq := dml.Quoter.Name(objectName)
			if len(dropIfExists) > 0 && dropIfExists[0] {
				if _, err := db.ExecContext(ctx, "DROP "+viewOrTable+" IF EXISTS "+vnq); err != nil {
					return errors.Wrapf(err, "[ddl] Drop view failed %q", objectName)
				}
			}

			_, err := db.ExecContext(ctx, "CREATE "+viewOrTable+" "+vnq+" AS "+query)
			if err != nil {
				return errors.Wrapf(err, "[ddl] Create view %q failed", objectName)
			}

			tc, err := LoadColumns(ctx, db, objectName)
			if err != nil {
				return errors.Wrapf(err, "[ddl] Load columns failed for %q", objectName)
			}

			if err := WithTable(objectName, tc[objectName]...).fn(tm); err != nil {
				return errors.Wrapf(err, "[ddl] Failed to add new table %q", objectName)
			}

			tm.mu.Lock()
			defer tm.mu.Unlock()
			tm.tm[objectName].IsView = viewOrTable == "VIEW"

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
func WithTable(tableName string, cols ...*Column) TableOption {
	return TableOption{
		fn: func(tm *Tables) error {
			if err := dml.IsValidIdentifier(tableName); err != nil {
				return errors.WithStack(err)
			}

			if err := tm.Upsert(NewTable(tableName, cols...)); err != nil {
				return errors.Wrap(err, "[ddl] WithNewTable.Tables.Insert")
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
func WithTableLoadColumns(ctx context.Context, db dml.Querier, names ...string) TableOption {
	return TableOption{
		fn: func(tm *Tables) error {
			for _, n := range names {
				if err := dml.IsValidIdentifier(n); err != nil {
					return errors.WithStack(err)
				}
			}

			tc, err := LoadColumns(ctx, db, names...)
			if err != nil {
				return errors.WithStack(err)
			}

			for _, n := range names {
				t := NewTable(n)
				t.Schema = tm.Schema

				t.Columns = tc[n]
				if err := tm.Upsert(t); err != nil {
					return errors.Wrapf(err, "[ddl] Tables.Insert for %q", n)
				}
			}
			return nil
		},
	}
}

// WithTableNames creates for each table name and its index a new table pointer.
// You should call afterwards the functional option WithLoadColumnDefinitions.
// This function returns an error if a table index already exists.
func WithTableNames(names ...string) TableOption {
	return TableOption{
		fn: func(tm *Tables) error {
			for _, name := range names {
				if err := dml.IsValidIdentifier(name); err != nil {
					return errors.WithStack(err)
				}
			}

			for _, tn := range names {
				if err := tm.Upsert(NewTable(tn)); err != nil {
					return errors.Wrapf(err, "[ddl] Tables.Insert %q", tn)
				}
			}
			return nil
		},
	}
}

// WithTableDMLListeners adds event listeners to a table object. It doesn't
// matter if the table has already been set. If the table object gets set later,
// the events will be copied to the new object.
func WithTableDMLListeners(tableName string, events ...*dml.ListenerBucket) TableOption {
	return TableOption{
		sortOrder: 253,
		fn: func(tm *Tables) error {
			tm.mu.Lock()
			defer tm.mu.Unlock()

			t, ok := tm.tm[tableName]
			if !ok {
				return errors.NotFound.Newf("[ddl] Table %q not found", tableName)
			}
			t.Listeners.Merge(events...)
			tm.tm[tableName] = t

			return nil
		},
	}
}

// NewTables creates a new TableService satisfying interface Manager.
func NewTables(opts ...TableOption) (*Tables, error) {
	tm := &Tables{
		tm: make(map[string]*Table),
	}
	if err := tm.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[ddl] NewTables applied option error")
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

// Options applies options to the Tables service.
func (tm *Tables) Options(opts ...TableOption) error {
	// SliceStable must be stable to maintain the order of all options where
	// sortOrder is zero.
	sort.SliceStable(opts, func(i, j int) bool {
		return opts[i].sortOrder < opts[j].sortOrder
	})

	for _, to := range opts {
		if err := to.fn(tm); err != nil {
			return errors.Wrap(err, "[ddl] Applied option error")
		}
	}
	tm.mu.Lock()
	for _, tbl := range tm.tm {
		if tbl.DB != tm.DB {
			tbl.DB = tm.DB
		}
	}
	tm.mu.Unlock()
	return nil
}

// Table returns the structure from a map m by a giving index i. What is the
// reason to use int as the table index and not a name? Because table names
// between M1 and M2 get renamed and in a Go SQL code generator script of the
// CoreStore project, we can guarantee that the generated index constant will
// always stay the same but the name of the table differs.
func (tm *Tables) Table(name string) (*Table, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if t, ok := tm.tm[name]; ok {
		return t, nil
	}
	return nil, errors.NotFound.Newf("[ddl] Table %q not found.", name)
}

// MustTable same as Table function but panics when the table cannot be found or
// any other error occurs.
func (tm *Tables) MustTable(name string) *Table {
	t, err := tm.Table(name)
	if err != nil {
		panic(err)
	}
	return t
}

// Tables returns a random list of all available table names. It can append the
// names to the argument slice.
func (tm *Tables) Tables(ret ...string) []string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if len(tm.tm) == 0 {
		return ret
	}
	if ret == nil {
		ret = make([]string, 0, len(tm.tm))
	}
	for tn := range tm.tm {
		ret = append(ret, tn)
	}
	return ret
}

// Len returns the number of all tables.
func (tm *Tables) Len() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return len(tm.tm)
}

// Upsert adds or updates a new table into the internal cache. If a table
// already exists, then the new table gets applied. The ListenerBuckets gets
// merged from the existing table to the new table, they will be appended to the
// new table buckets. Empty columns in the new table gets updated from the
// existing table.
func (tm *Tables) Upsert(tNew *Table) error {

	tm.mu.Lock()
	defer tm.mu.Unlock()

	tOld, ok := tm.tm[tNew.Name]
	if tOld == nil || !ok {
		tm.tm[tNew.Name] = tNew
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

	tm.tm[tNew.Name] = tNew.update()
	return nil
}

// DeleteFromCache removes tables by their given indexes. If no index has been passed
// then all entries get removed and the map reinitialized.
func (tm *Tables) DeleteFromCache(tableNames ...string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	for _, tn := range tableNames {
		delete(tm.tm, tn)
	}
}

// DeleteAllFromCache clears the internal table cache and resets the map.
func (tm *Tables) DeleteAllFromCache() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	// maybe clear each pointer in the Table struct to avoid a memory leak
	tm.tm = make(map[string]*Table)
}

// MapColumns scans a row from a database. It creates automatically a new Table
// object for non-existing ones. Existing tables gets reset their columns slice
// and it refreshes them.
func (tm *Tables) MapColumns(rc *dml.ColumnMap) error {
	if rc.Count == 0 {
		tm.mu.Lock()
	}

	c, tableName, err := NewColumn(rc)
	if err != nil {
		return errors.WithStack(err)
	}

	t, ok := tm.tm[tableName]
	if !ok {
		t = NewTable(tableName)
		tm.tm[tableName] = t
	}

	if tm.previousTable != tableName {
		tm.previousTable = tableName
		t.resetColumns()
	}

	t.Columns = append(t.Columns, c)
	return nil
}

// Close implements io.Closer interface used in dml.Load. It unlocks the
// internal mutex.
func (tm *Tables) Close() error {
	tm.mu.Unlock()
	return nil
}

// ToSQL returns the SQL string for loading the column definitions of either all
// tables or of the already created Table objects.
func (tm *Tables) ToSQL() (string, []interface{}, error) {
	if tn := tm.Tables(); len(tn) > 0 {
		return dml.Interpolate(selTablesColumns).Strs(tn...).ToSQL()
	}
	return dml.QuerySQL(selAllTablesColumns).ToSQL()
}

// Validate validates the table names and their column against the current
// database schema. The context is used to maybe cancel the "Load Columns"
// query.
func (tm *Tables) Validate(ctx context.Context) error {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tblNames := make([]string, 0, len(tm.tm))
	for tn := range tm.tm {
		tblNames = append(tblNames, tn)
	}

	tMap, err := LoadColumns(ctx, tm.DB, tblNames...)
	if err != nil {
		return errors.WithStack(err)
	}
	if have, want := len(tMap), len(tm.tm); have != want {
		return errors.Mismatch.Newf("[ddl] Tables count %d does not match table count %d in database.", want, have)
	}
	dbTableNames := make([]string, 0, len(tMap))
	for tn := range tMap {
		dbTableNames = append(dbTableNames, tn)
	}
	sort.Strings(dbTableNames)

	// TODO compare it that way, that the DB table is the master and Go objects must be updated
	// once they do not match the database version.
	for tn, tbl := range tm.tm {
		dbTblCols, ok := tMap[tn]
		if !ok {
			return errors.NotFound.Newf("[ddl] Table %q not found in database. Available tables: %v", tn, dbTableNames)
		}
		if want, have := len(tbl.Columns), len(dbTblCols); want > have {
			return errors.Mismatch.Newf("[ddl] Table %q has more columns (count %d) than its object (column count %d) in the database.", tn, want, have)
		}
		for idx, c := range tbl.Columns {
			dbCol := dbTblCols[idx]
			if c.Field != dbCol.Field {
				return errors.Mismatch.Newf("[ddl] Table %q with column name %q at index %d does not match database column name %q",
					tn, c.Field, idx, dbCol.Field,
				)
			}
			if c.ColumnType != dbCol.ColumnType {
				return errors.Mismatch.Newf("[ddl] Table %q with Go column name %q does not match MySQL column type. MySQL: %q Go: %q.",
					tn, c.Field, dbCol.ColumnType, c.ColumnType,
				)
			}
			if c.Null != dbCol.Null {
				return errors.Mismatch.Newf("[ddl] Table %q with column name %q does not match MySQL null types. MySQL: %q Go: %q",
					tn, c.Field, dbCol.Null, c.Null,
				)
			}
			// maybe more comparisons
		}
	}

	return nil
}
