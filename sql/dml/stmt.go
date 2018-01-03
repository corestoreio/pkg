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

package dml

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Stmter represents a prepared statement. It wraps a *sql.Stmt with a specific
// SQL query. To create a Stmter call the Prepare function of a type DML type.
// For now Stmter is not safe for concurrent use, despite the underlying
// *sql.Stmt is. Don't forget to call Close!
type Stmter interface {
	// Stmt returns the raw prepared statement for your convenience.
	Stmt() *sql.Stmt
	// Close closes the underlying prepared statement.
	Close() error
	// WithArgs sets the interfaced arguments for the execution with Query+. It
	// internally resets previously applied arguments.
	WithArgs(args ...interface{}) Stmter
	// WithArguments sets the arguments for the execution with Query+. It
	// internally resets previously applied arguments.
	WithArguments(args Arguments) Stmter
	// WithRecords sets the records for the execution with Query+. It internally
	// resets previously applied arguments.
	WithRecords(records ...QualifiedRecord) Stmter
	// Exec supports both either the traditional way or passing arguments or
	// in combination with the previously called WithArguments, WithRecords or
	// WithArgs functions. If you want to call it multiple times with the same
	// arguments, do not use the `args` variable, instead use the With+ functions.
	// Calling any of the With+ function and additionally setting the `args`, will
	// append the `args` at the end to the previously set or generated arguments.
	// This function is not thread safe.
	Exec(ctx context.Context, args ...interface{}) (sql.Result, error)
	// Query traditional way, allocation heavy.
	Query(ctx context.Context, args ...interface{}) (*sql.Rows, error)
	// QueryRow traditional way, allocation heavy.
	QueryRow(ctx context.Context, args ...interface{}) *sql.Row
	// Load loads data from a query into an object. Load can load a single row
	// or n-rows.
	Load(ctx context.Context, s ColumnMapper) (rowCount uint64, err error)
	// LoadInt64 executes the prepared statement and returns the value as an
	// int64. It returns a NotFound error if the query returns nothing.
	LoadInt64(ctx context.Context) (int64, error)
	// LoadInt64s executes the Select and returns the value as a slice of
	// int64s.
	LoadInt64s(ctx context.Context) (ret []int64, err error)
}

// stmtBase wraps a *sql.Stmt (a prepared statement) with a specific SQL query.
// To create a stmtBase call the Prepare function of type Select. stmtBase is
// not safe for concurrent use, despite the underlying *sql.Stmt is. Don't
// forget to call Close!
type stmtBase struct {
	source byte
	builderCommon
	stmt *sql.Stmt
}

func (st *stmtBase) Stmt() *sql.Stmt                               { return st.stmt }
func (st *stmtBase) Close() error                                  { return st.stmt.Close() }
func (st *stmtBase) WithArgs(args ...interface{}) Stmter           { st.withArgs(args); return st }
func (st *stmtBase) WithArguments(args Arguments) Stmter           { st.withArguments(args); return st }
func (st *stmtBase) WithRecords(records ...QualifiedRecord) Stmter { st.withRecords(records); return st }

func (st *stmtBase) resetArgs() {
	st.argsArgs = st.argsArgs[:0]
	st.argsRaw = st.argsRaw[:0]
	st.argsRecords = st.argsRecords[:0]
}

func (st *stmtBase) withArgs(args []interface{}) {
	st.resetArgs()
	st.argsRaw = args
	st.isWithInterfaces = true
}

func (st *stmtBase) withArguments(args Arguments) {
	st.resetArgs()
	st.argsArgs = args
	st.isWithInterfaces = false
}

// withRecords sets the records for the execution with Query or Exec. It
// internally resets previously applied arguments.
func (st *stmtBase) withRecords(records []QualifiedRecord) {
	st.resetArgs()
	st.argsRecords = records
	st.isWithInterfaces = false
}

// prepareArgs transforms mainly the Arguments into []interface{} but also
// appends the `args` from the Exec+ or Query+ function.
// All method receivers are not thread safe.
func (st *stmtBase) prepareArgs(args ...interface{}) error {
	if st.ärgErr != nil {
		return st.ärgErr
	}

	if !st.isWithInterfaces {
		st.argsRaw = st.argsRaw[:0]
	}

	argsArgs, err := st.convertRecordsToArguments()
	st.argsRaw = append(st.argsRaw, argsArgs.Interfaces()...)
	st.argsRaw = append(st.argsRaw, args...)
	return err
}

func (st *stmtBase) Exec(ctx context.Context, args ...interface{}) (sql.Result, error) {
	if err := st.prepareArgs(args...); err != nil {
		return nil, errors.WithStack(err)
	}
	if st.Log != nil && st.Log.IsDebug() {
		defer log.WhenDone(st.Log).Debug("Exec", log.Int("arg_len", len(st.argsRaw)))
	}
	return st.stmt.ExecContext(ctx, st.argsRaw...)
}

func (st *stmtBase) Query(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	if err := st.prepareArgs(args...); err != nil {
		return nil, errors.WithStack(err)
	}
	if st.Log != nil && st.Log.IsDebug() {
		defer log.WhenDone(st.Log).Debug("Query", log.Int("arg_len", len(st.argsRaw)))
	}
	return st.stmt.QueryContext(ctx, st.argsRaw...)
}

func (st *stmtBase) QueryRow(ctx context.Context, args ...interface{}) *sql.Row {
	if err := st.prepareArgs(args...); err != nil {
		_ = err
		// Hmmm what should happen here?
	}
	if st.Log != nil && st.Log.IsDebug() {
		defer log.WhenDone(st.Log).Debug("QueryRow", log.Int("arg_len", len(st.argsRaw)))
	}
	return st.stmt.QueryRowContext(ctx, st.argsRaw...)
}

func (st *stmtBase) Load(ctx context.Context, s ColumnMapper) (rowCount uint64, err error) {
	if st.Log != nil && st.Log.IsDebug() {
		defer log.WhenDone(st.Log).Debug("Load", log.Uint64("row_count", rowCount), log.String("object_type", fmt.Sprintf("%T", s)), log.Err(err))
	}
	r, err := st.Query(ctx)
	rowCount, err = load(r, err, s, &st.ColumnMap)
	return rowCount, errors.WithStack(err)
}

func (st *stmtBase) LoadInt64(ctx context.Context) (int64, error) {
	if st.Log != nil && st.Log.IsDebug() {
		defer log.WhenDone(st.Log).Debug("LoadInt64")
	}
	return loadInt64(st.Query(ctx))
}

func (st *stmtBase) LoadInt64s(ctx context.Context) (ret []int64, err error) {
	if st.Log != nil && st.Log.IsDebug() {
		// do not use fullSQL because we might log sensitive data
		defer log.WhenDone(st.Log).Debug("LoadInt64s", log.Int("row_count", len(ret)), log.Err(err))
	}
	ret, err = loadInt64s(st.Query(ctx))
	// Do not simplify it because we need ret in the defer. we don't log errors
	// because they get handled.
	return ret, err
}

// More Load* functions can be added later

func newReduxStmt(db *sql.DB, name string, qb QueryBuilder, idleTime time.Duration, l log.Logger) (*reduxStmt, error) {
	rs := &reduxStmt{
		name:           name,
		db:             db,
		qb:             qb,
		stopIdleDaemon: make(chan struct{}),
		idleDuration:   idleTime,
		lastUsed:       time.Now(),
	}
	rs.stmtBase.Log = l
	go rs.idleDaemon()
	return rs, nil
}

// reduxStmt represents a prepared statement which can resurrect/redux and is
// bound to a single connection. After an idle time the statement and the
// connection gets closed and resources on the DB server freed. Once the
// statement will be used again, it reduxes, gets re-prepared, via its dedicated
// connection. If there are no more connections available, it waits until it can
// connect. Due to the connection handling implementation of database/sql/DB
// object we must grab a dedicated connection. A later aim would that multiple
// prepared statements can share a single connection.
type reduxStmt struct {
	stmtBase
	name string
	db   *sql.DB // to get a new con
	qb   QueryBuilder

	// idleDuration defines the duration how long to wait until no query will be
	// executed and the prepared statement deallocated.
	idleDuration   time.Duration
	stopIdleDaemon chan struct{} // tells the ticker to stop and closes the stmt and the con

	mu       sync.RWMutex
	con      *sql.Conn // stmt bound to this con, gets recreated each time a query gets re-prepared
	lastUsed time.Time // time when the stmt has last been used
	status   byte      // c=closed, p=prepared, 0 = nothing
}

func (rs *reduxStmt) Stmt() *sql.Stmt {
	panic(errors.NotSupported.Newf("[dml] returning the raw sql.Stmt is not yet supported (%q", rs.name))
}

func (rs *reduxStmt) Close() error {
	rs.mu.RLock()
	s := rs.status
	rs.mu.RUnlock()
	if s == 'c' {
		return nil
	}
	close(rs.stopIdleDaemon)
	return rs.closeStmtCon()
}

func (rs *reduxStmt) closeStmtCon() (err error) {
	rs.mu.Lock()
	if rs.status == 'c' {
		rs.mu.Unlock()
		return nil
	}
	rs.status = 'c'
	defer func() {
		if err2 := rs.con.Close(); err == nil && err2 != nil {
			err = errors.Wrapf(err2, "[dml] reduxStmt.closeStmtCon.con name: %q", rs.name)
		}
		rs.con = nil
		if rs.Log != nil && rs.Log.IsDebug() {
			rs.Log.Debug("reduxStmt.closeStmtCon.con.close", log.String("name", rs.name), log.Time("last_used", rs.lastUsed))
		}
		if err == nil && rs.ärgErr != nil {
			err = rs.ärgErr
		}
		rs.mu.Unlock()
	}()
	if rs.stmt != nil {
		if err = rs.stmt.Close(); err != nil {
			err = errors.Wrapf(err, "[ddl] reduxStmt.stmt.close name: %q", rs.name)
		}
		if rs.Log != nil && rs.Log.IsDebug() {
			rs.Log.Debug("reduxStmt.closeStmtCon.stmt.close", log.String("name", rs.name), log.Time("last_used", rs.lastUsed))
		}
	}
	return err
}

func (rs *reduxStmt) rePrepare() (err error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	// optimization to remove the lock can be to rely on sync/atomic. store the
	// time as nano second and status aus an int value. Via CAS or Store/Load
	// operations performance should improve.
	rs.lastUsed = time.Now()
	if rs.status == 'p' {
		if rs.Log != nil && rs.Log.IsDebug() {
			rs.Log.Debug("reduxStmt.rePrepare.stmt.prepared", log.String("name", rs.name), log.Time("last_used", rs.lastUsed))
		}
		return nil
	}
	ctx := context.Background() // for now, can be changed later
	// get a fresh connection; maybe use max tries and back off ... etc
	if rs.con, err = rs.db.Conn(ctx); err != nil {
		return errors.WithStack(err)
	}

	qry, _, err := rs.qb.ToSQL()
	if err != nil {
		return errors.WithStack(err)
	}
	rs.stmt, err = rs.con.PrepareContext(ctx, qry)
	if err != nil {
		return errors.WithStack(err)
	}
	rs.status = 'p'
	if rs.Log != nil && rs.Log.IsDebug() {
		rs.Log.Debug("reduxStmt.rePrepare.stmt.preparing", log.String("name", rs.name), log.Time("last_used", rs.lastUsed), log.String("query", qry))
	}
	return nil
}

func (rs *reduxStmt) idleDaemon() {
	ticker := time.NewTicker(rs.idleDuration)
	for {
		select {
		case t, ok := <-ticker.C:
			if !ok {
				return
			}
			t = t.Add(-rs.idleDuration)
			if rs.canClose(t) {
				// stmt has not been used within the last duration. so close
				// the stmt and release the resources in the DB. And also close
				// the dedicated connection!
				if err := rs.closeStmtCon(); err != nil && rs.Log != nil && rs.Log.IsInfo() {
					rs.Log.Info("dml.reduxStmt.close.error", log.Err(err), log.String("name", rs.name))
				}
				if rs.Log != nil && rs.Log.IsDebug() {
					rs.Log.Debug("reduxStmt.idleDaemon.stmt.closing", log.String("name", rs.name), log.Time("last_used", rs.lastUsed), log.Time("current", t))
				}
			}
		case <-rs.stopIdleDaemon:
			if rs.Log != nil && rs.Log.IsDebug() {
				rs.Log.Debug("reduxStmt.idleDaemon.ticker.stopped", log.String("name", rs.name), log.Time("last_used", rs.lastUsed))
			}
			ticker.Stop()
			return
		}
	}
}

func (rs *reduxStmt) canClose(t time.Time) bool {
	rs.mu.RLock()
	ok := t.After(rs.lastUsed) && rs.status == 'p'
	rs.mu.RUnlock()
	return ok
}

func (rs *reduxStmt) Exec(ctx context.Context, args ...interface{}) (sql.Result, error) {
	if err := rs.rePrepare(); err != nil {
		return nil, errors.WithStack(err)
	}
	res, err := rs.stmtBase.Exec(ctx, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "[dml] reduxStmt.Exec with name %q", rs.name)
	}
	return res, nil
}

// Query traditional way, allocation heavy.
func (rs *reduxStmt) Query(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	if err := rs.rePrepare(); err != nil {
		return nil, errors.WithStack(err)
	}
	rows, err := rs.stmtBase.Query(ctx, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "[dml] reduxStmt.Query with name %q", rs.name)
	}
	return rows, nil
}

// QueryRow traditional way, allocation heavy.
func (rs *reduxStmt) QueryRow(ctx context.Context, args ...interface{}) *sql.Row {
	if err := rs.rePrepare(); err != nil {
		rs.ärgErr = errors.WithStack(err)
	}
	return rs.stmtBase.QueryRow(ctx, args...)
}

// Load loads data from a query into an object. Load can load a single row
// or n-rows.
func (rs *reduxStmt) Load(ctx context.Context, s ColumnMapper) (rowCount uint64, err error) {
	if err = rs.rePrepare(); err != nil {
		return 0, errors.WithStack(err)
	}
	rowCount, err = rs.stmtBase.Load(ctx, s)
	if err != nil {
		err = errors.Wrapf(err, "[dml] reduxStmt.Load with name %q", rs.name)
	}
	return
}

// LoadInt64 executes the prepared statement and returns the value as an
// int64. It returns a NotFound error if the query returns nothing.
func (rs *reduxStmt) LoadInt64(ctx context.Context) (value int64, err error) {
	if err = rs.rePrepare(); err != nil {
		return 0, errors.WithStack(err)
	}
	value, err = rs.stmtBase.LoadInt64(ctx)
	if err != nil {
		err = errors.Wrapf(err, "[dml] reduxStmt.LoadInt64 with name %q", rs.name)
	}
	return
}

// LoadInt64s executes the Select and returns the value as a slice of
// int64s.
func (rs *reduxStmt) LoadInt64s(ctx context.Context) (values []int64, err error) {
	if err := rs.rePrepare(); err != nil {
		return nil, errors.WithStack(err)
	}
	values, err = rs.stmtBase.LoadInt64s(ctx)
	if err != nil {
		err = errors.Wrapf(err, "[dml] reduxStmt.LoadInt64s with name %q", rs.name)
	}
	return
}
