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
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
)

const (
	// PrefixView  is an anti-pattern I've seen many such systems where at some point a view will become a table.
	// deprecated
	PrefixView      = "view_" // If identifier starts with this, it is considered a view.
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
	dcp *dml.ConnPool // database connection pool
	// Schema represents the name of the database. Might be empty.
	Schema string

	mu sync.RWMutex
	// tm a map where key = table name and value the table pointer
	tm map[string]*Table
}

// WithDB sets the DB object to the Tables and all sub Table types to handle the
// database connections. It must be set if other options get used to access the
// DB.
func WithDB(db *sql.DB, opts ...dml.ConnPoolOption) TableOption {
	return TableOption{
		sortOrder: 1,
		fn: func(tm *Tables) error {
			p, err := dml.NewConnPool(append(opts, dml.WithDB(db))...)
			if err != nil {
				return errors.WithStack(err)
			}
			return WithConnPool(p).fn(tm)
		},
	}
}

// WithConnPool sets the connection pool to the Tables and each of it Table
// type. This function has precedence over WithDB.
func WithConnPool(db *dml.ConnPool) TableOption {
	return TableOption{
		sortOrder: 2,
		fn: func(tm *Tables) error {
			tm.mu.Lock()
			defer tm.mu.Unlock()
			tm.dcp = db
			for _, t := range tm.tm {
				t.dcp = db
			}
			return nil
		},
	}
}

// WithTable inserts a new table to the Tables struct. You can optionally
// specify the columns. Without columns the call to load the columns from the
// INFORMATION_SCHEMA must be added.
func WithTable(tableName string, cols ...*Column) TableOption {
	return TableOption{
		sortOrder: 10,
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

// WithCreateTable upserts tables to the current `Tables` object. Either it adds
// a new table/view or overwrites existing entries. Argument
// `identifierCreateSyntax` must be balanced slice where index i is the
// table/view name and i+1 can be either empty or contain the SQL CREATE
// statement. In case a SQL CREATE statement has been supplied, it gets executed
// otherwise ignored. After table initialization the create syntax and the
// column specifications are getting loaded but only if a connection has been
// set beforehand. Write the SQL CREATE statement in upper case.
//		WithCreateTable(
//			"sales_order_history", "CREATE TABLE `sales_order_history` ( ... )", // table created if not exists
//			"sales_order_stat", "CREATE VIEW `sales_order_stat` AS SELECT ...", // table created if not exists
//			"sales_order", "", // table/view already exists and gets loaded, NOT dropped.
//		)
func WithCreateTable(ctx context.Context, identifierCreateSyntax ...string) TableOption {
	return TableOption{
		sortOrder: 50,
		fn: func(tm *Tables) error {
			tm.mu.Lock()
			defer tm.mu.Unlock()

			lenIDCS := len(identifierCreateSyntax)
			if lenIDCS%2 == 1 {
				return errors.NotValid.Newf("[ddl] WithCreateTable expects a balanced slice, but got %d items.", lenIDCS)
			}

			tvNames := make([]string, 0, lenIDCS/2)
			for i := 0; i < lenIDCS; i = i + 2 {
				// tv = table or view
				tvName := identifierCreateSyntax[i]
				tvCreate := identifierCreateSyntax[i+1]

				if err := dml.IsValidIdentifier(tvName); err != nil {
					return errors.WithStack(err)
				}

				tvNames = append(tvNames, tvName)
				t := NewTable(tvName)
				tm.tm[tvName] = t

				if isCreateStmt(tvName, tvCreate) {
					t.IsView = strings.Contains(tvCreate, " VIEW ") || strings.HasPrefix(tvName, PrefixView)
					if _, err := tm.dcp.DB.ExecContext(ctx, tvCreate); err != nil {
						return errors.Wrapf(err, "[ddl] WithCreateTable failed to run for table %q the query: %q", tvName, tvCreate)
					}
				}
			}
			if tm.dcp == nil {
				return nil
			}
			tc, err := LoadColumns(ctx, tm.dcp.DB, tvNames...)
			if err != nil {
				return errors.WithStack(err)
			}
			for _, n := range tvNames {
				t := tm.tm[n]
				t.Schema = tm.Schema
				t.Columns = tc[n]
				t.update()
			}
			return nil
		},
	}
}

var regexpCreateTable = regexp.MustCompile(`CREATE\s+(VIEW|TABLE)\s*(?:IF\s+NOT\s+EXISTS)?\s+`)

func isCreateStmt(idName, stmt string) bool {
	return regexpCreateTable.MatchString(stmt) && strings.Contains(stmt, idName)
}

func isCreateStmtBytes(idName, stmt []byte) bool {
	return regexpCreateTable.Match(stmt) && bytes.Contains(stmt, idName)
}

// WithCreateTableFromFile creates the defined tables from the loaded *.sql
// files.
func WithCreateTableFromFile(ctx context.Context, globPattern string, tableNames ...string) TableOption {
	return TableOption{
		sortOrder: 60,
		fn: func(tm *Tables) error {
			matches, err := filepath.Glob(globPattern)
			if err != nil {
				return errors.Wrapf(err, "[ddl] WithCreateTableFromFile and pattern %q", globPattern)
			}
			identifierCreateSyntax, err := loadSQLFiles(matches, tableNames)
			if err != nil {
				return errors.WithStack(err)
			}
			return WithCreateTable(ctx, identifierCreateSyntax...).fn(tm)
		},
	}
}

func loadSQLFiles(fileNames, tableNames []string) ([]string, error) {
	ret := make([]string, 0, len(tableNames)*2)
	var notFound []string
	for _, tn := range tableNames {
		found := false
		for _, fn := range fileNames {
			if strings.Contains(fn, tn) {
				data, err := ioutil.ReadFile(fn)
				if err != nil {
					return nil, errors.ReadFailed.New(err, "[ddl] WithCreateTableFromFile failed to file %q for table %q", fn, tn)
				}
				if !isCreateStmtBytes([]byte(tn), data) { // drop all comments
					return nil, errors.NotAllowed.Newf("[ddl] WithCreateTableFromFile allows only CREATE TABLE|VIEW statements, got %q", data)
				}
				ret = append(ret, tn, string(data))
				found = true
			}
		}
		if !found {
			notFound = append(notFound, tn)
		}
	}
	if len(notFound) > 0 {
		return nil, errors.Mismatch.Newf("[dd] WithCreateTableFromFile cannot load the files for tables: %v", notFound)
	}
	return ret, nil
}

// WithDropTable drops the tables or views listed in argument `tableViewNames`.
// If argument `option` contains the string "DISABLE_FOREIGN_KEY_CHECKS", then foreign keys get disabled
// and at the end re-enabled.
func WithDropTable(ctx context.Context, option string, tableViewNames ...string) TableOption {
	return TableOption{
		sortOrder: 11,
		fn: func(tm *Tables) (err error) {
			tm.mu.Lock()
			defer tm.mu.Unlock()

			if option != "" && strings.Contains(strings.ToUpper(option), "DISABLE_FOREIGN_KEY_CHECKS") {
				return tm.dcp.WithDisabledForeignKeyChecks(ctx, func(conn *dml.Conn) error {
					return withDropTable(ctx, tm, conn.DB, tableViewNames)
				})
			}
			return withDropTable(ctx, tm, tm.dcp.DB, tableViewNames)
		},
	}
}

func withDropTable(ctx context.Context, tm *Tables, db dml.QueryExecPreparer, tableViewNames []string) (err error) {
	for _, name := range tableViewNames {
		if t, ok := tm.tm[name]; ok {
			t.customDB = db
			if err = t.Drop(ctx); err != nil {
				t.customDB = nil
				return errors.WithStack(err)
			}
			t.customDB = nil
			continue
		}

		if err := dml.IsValidIdentifier(name); err != nil {
			return errors.WithStack(err)
		}
		typ := "TABLE"
		if strings.HasPrefix(name, PrefixView) {
			typ = "VIEW"
		}
		if _, err = db.ExecContext(ctx, "DROP "+typ+" IF EXISTS "+dml.Quoter.Name(name)); err != nil {
			return errors.Wrapf(err, "[ddl] Failed to drop %q", name)
		}
	}
	return nil
}

// NewTables creates a new TableService satisfying interface Manager.
func NewTables(opts ...TableOption) (*Tables, error) {
	tm := &Tables{
		tm: make(map[string]*Table),
	}
	if err := tm.Options(opts...); err != nil {
		return nil, errors.WithStack(err)
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
			return errors.WithStack(err)
		}
	}
	tm.mu.Lock()
	for _, tbl := range tm.tm {
		if tbl.dcp != tm.dcp {
			tbl.dcp = tm.dcp
		}
	}
	tm.mu.Unlock()
	return nil
}

// errTableNotFound provides a custom error behaviour with not capturing the
// stack trace and hence less allocs.
type errTableNotFound string

func (t errTableNotFound) ErrorKind() errors.Kind { return errors.NotFound }
func (t errTableNotFound) Error() string {
	return fmt.Sprintf("[ddl] Table %q not found or not yet added.", string(t))
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
	return nil, errTableNotFound(name)
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

	tMap, err := LoadColumns(ctx, tm.dcp.DB, tblNames...)
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

// Truncate force truncates all tables by also disabling foreign keys.
func (tm *Tables) Truncate(ctx context.Context) error {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return DisableForeignKeys(ctx, tm.dcp.DB, func() error {
		for _, t := range tm.tm {
			if err := t.Truncate(ctx); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})
}
