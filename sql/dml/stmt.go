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
	"sync"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Stmter represents a prepared statement. It wraps a *sql.Stmt with a specific
// SQL query. To create a Stmter call the Prepare function of a type DML type.
// Stmter is not safe for concurrent use, you must call for each goroutine the Clone function, to create
// a dedicated instance. Don't forget to call Close!

// Stmt wraps a *sql.Stmt (a prepared statement) with a specific SQL query.
// To create a Stmt call the Prepare function of type Select. Stmt is
// not safe for concurrent use, despite the underlying *sql.Stmt is. Don't
// forget to call Close!
type Stmt struct {
	base builderCommon
	Stmt *sql.Stmt
}

func (st *Stmt) WithArgs(rawArgs ...interface{}) *Arguments {
	var args [defaultArgumentsCapacity]argument
	return &Arguments{
		base:      st.base,
		raw:       rawArgs,
		arguments: args[:0],
	}
}

// Closes closes the statement in the database and frees its resources.
func (st *Stmt) Close() error { return st.Stmt.Close() }

// More Load* functions can be added later

// Queueing https://github.com/eapache/queue/blob/master/queue.go but maybe a lock can enough.
func newReduxStmt(db *sql.DB, name string, qb QueryBuilder, idleTime time.Duration, l log.Logger) (*StmtRedux, error) {
	query, _, err := qb.ToSQL()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rs := &StmtRedux{
		name:           name,
		db:             db,
		query:          query,
		stopIdleDaemon: make(chan struct{}),
		idleDuration:   idleTime,
		lastUsed:       time.Now(),
	}
	rs.stmt.base.Log = l
	go rs.idleDaemon()
	return rs, nil
}

// StmtRedux represents a prepared statement which can resurrect/redux and is
// bound to a single connection. After an idle time the statement and the
// connection gets closed and resources on the DB server freed. Once the
// statement will be used again, it reduxes, gets re-prepared, via its dedicated
// connection. If there are no more connections available, it waits until it can
// connect. Due to the connection handling implementation of database/sql/DB
// object we must grab a dedicated connection. A later aim would that multiple
// prepared statements can share a single connection.
type StmtRedux struct {
	stmt  Stmt
	name  string
	query string
	db    *sql.DB // to get a new con

	// idleDuration defines the duration how long to wait until no query will be
	// executed and the prepared statement deallocated.
	idleDuration   time.Duration
	stopIdleDaemon chan struct{} // tells the ticker to stop and closes the Stmt and the con

	mu       sync.RWMutex
	con      *sql.Conn // Stmt bound to this con, gets recreated each time a query gets re-prepared
	lastUsed time.Time // time when the Stmt has last been used
	status   byte      // c=closed, p=prepared,   0 = nothing
	inUse    bool
}

func (rs *StmtRedux) Close() error {
	rs.mu.RLock()
	s := rs.status
	rs.mu.RUnlock()
	if s == 'c' {
		return nil
	}
	close(rs.stopIdleDaemon)
	return rs.closeStmtCon()
}

// WithArgs returns a new type to support multiple executions of the underlying
// SQL statement and reuse of memory allocations for the arguments. WithArgs
// builds the SQL string and sets the optional raw interfaced arguments for the
// later execution. It copies the underlying connection and settings from the
// current DML type (Delete, Insert, Select, Update, Union, With, etc.). The
// query executor can still be overwritten. Interpolation does not support the
// raw interfaces.
func (st *StmtRedux) WithArgs(rawArgs ...interface{}) *Arguments {
	// todo: correct implementation
	var args [defaultArgumentsCapacity]argument
	return &Arguments{
		base:      st.stmt.base,
		raw:       rawArgs,
		arguments: args[:0],
	}
}

func (rs *StmtRedux) closeStmtCon() (err error) {
	rs.mu.Lock()
	if rs.status == 'c' || rs.con == nil {
		rs.mu.Unlock()
		return nil
	}
	rs.status = 'c'
	defer func() {
		if err2 := rs.con.Close(); err == nil && err2 != nil {
			err = errors.Wrapf(err2, "[dml] StmtRedux.closeStmtCon.con name: %q", rs.name)
		}
		rs.con = nil
		if rs.stmt.base.Log != nil && rs.stmt.base.Log.IsDebug() {
			rs.stmt.base.Log.Debug("StmtRedux.closeStmtCon.con.close", log.String("name", rs.name), log.Time("last_used", rs.lastUsed))
		}
		if err == nil && rs.stmt.base.ärgErr != nil {
			err = rs.stmt.base.ärgErr
		}
		rs.mu.Unlock()
	}()
	if rs.stmt.Stmt != nil {
		if err = rs.stmt.Stmt.Close(); err != nil {
			err = errors.Wrapf(err, "[ddl] StmtRedux.Stmt.close name: %q", rs.name)
		}
		if rs.stmt.base.Log != nil && rs.stmt.base.Log.IsDebug() {
			rs.stmt.base.Log.Debug("StmtRedux.closeStmtCon.Stmt.close", log.String("name", rs.name), log.Time("last_used", rs.lastUsed))
		}
	}
	return err
}

// rePrepare protected by an outer mutex. As long as the mutex locks, the
// statement is in use.
func (rs *StmtRedux) rePrepare() (err error) {
	// optimization to remove the lock can be to rely on sync/atomic. store the
	// time as nano second and status aus an int value. Via CAS or Store/Load
	// operations performance should improve.
	rs.lastUsed = time.Now()
	if rs.status == 'p' {
		if rs.stmt.base.Log != nil && rs.stmt.base.Log.IsDebug() {
			rs.stmt.base.Log.Debug("StmtRedux.rePrepare.Stmt.prepared", log.String("name", rs.name), log.Time("last_used", rs.lastUsed))
		}
		return nil
	}
	ctx := context.Background() // for now, can be changed later
	// get a fresh connection; maybe use max tries and back off ... etc
	if rs.con, err = rs.db.Conn(ctx); err != nil {
		return errors.WithStack(err)
	}

	rs.stmt.Stmt, err = rs.con.PrepareContext(ctx, rs.query)
	if err != nil {
		return errors.WithStack(err)
	}
	rs.status = 'p'
	if rs.stmt.base.Log != nil && rs.stmt.base.Log.IsDebug() {
		rs.stmt.base.Log.Debug("StmtRedux.rePrepare.Stmt.preparing", log.String("name", rs.name), log.Time("last_used", rs.lastUsed), log.String("query", rs.query))
	}
	return nil
}

// idleDaemon runs in its own goroutine
func (rs *StmtRedux) idleDaemon() {
	ticker := time.NewTicker(rs.idleDuration)
	for {
		select {
		case t, ok := <-ticker.C:
			if !ok {
				return
			}
			t = t.Add(-rs.idleDuration)
			if rs.canClose(t) {
				// Stmt has not been used within the last duration. so close
				// the Stmt and release the resources in the DB. And also close
				// the dedicated connection!
				if err := rs.closeStmtCon(); err != nil && rs.stmt.base.Log != nil && rs.stmt.base.Log.IsInfo() {
					rs.stmt.base.Log.Info("dml.StmtRedux.close.error", log.Err(err), log.String("name", rs.name))
				}
				if rs.stmt.base.Log != nil && rs.stmt.base.Log.IsDebug() {
					rs.stmt.base.Log.Debug("StmtRedux.idleDaemon.Stmt.closing", log.String("name", rs.name), log.Time("last_used", rs.lastUsed), log.Time("current", t))
				}
			}
		case <-rs.stopIdleDaemon:
			if rs.stmt.base.Log != nil && rs.stmt.base.Log.IsDebug() {
				rs.stmt.base.Log.Debug("StmtRedux.idleDaemon.ticker.stopped", log.String("name", rs.name), log.Time("last_used", rs.lastUsed))
			}
			ticker.Stop()
			return
		}
	}
}

func (rs *StmtRedux) canClose(t time.Time) bool {
	rs.mu.RLock()
	ok := t.After(rs.lastUsed) && rs.status == 'p'
	rs.mu.RUnlock()
	return ok
}

func (rs *StmtRedux) lock() {
	rs.mu.Lock()
	rs.inUse = true
}

func (rs *StmtRedux) unlock() {
	rs.inUse = false
	rs.mu.Unlock()
}

func (rs *StmtRedux) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	rs.lock()
	defer rs.unlock()
	if err := rs.rePrepare(); err != nil {
		return nil, errors.WithStack(err)
	}

	res, err := rs.stmt.Stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "[dml] StmtRedux.Exec with name %q", rs.name)
	}
	return res, nil
}

// Query traditional way, allocation heavy.
func (rs *StmtRedux) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	rs.lock()
	defer rs.unlock()
	if err := rs.rePrepare(); err != nil {
		return nil, errors.WithStack(err)
	}
	rows, err := rs.stmt.Stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "[dml] StmtRedux.Query with name %q", rs.name)
	}
	return rows, nil
}

// QueryRow traditional way, allocation heavy.
func (rs *StmtRedux) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	rs.lock()
	defer rs.unlock()
	if err := rs.rePrepare(); err != nil {
		rs.stmt.base.ärgErr = errors.WithStack(err)
	}
	return rs.stmt.Stmt.QueryRowContext(ctx, args...)
}

//// Load loads data from a query into an object. Load can load a single row
//// or n-rows.
//func (rs *StmtRedux) Load(ctx context.Context, s ColumnMapper, args ...interface{}) (rowCount uint64, err error) {
//	rs.lock()
//	defer rs.unlock()
//	if err = rs.rePrepare(); err != nil {
//		return 0, errors.WithStack(err)
//	}
//	rowCount, err = rs.stmt.Stmt.Load(ctx, s, args...)
//	if err != nil {
//		err = errors.Wrapf(err, "[dml] StmtRedux.Load with name %q", rs.name)
//	}
//	return
//}
//
//// LoadInt64 executes the prepared statement and returns the value as an
//// int64. It returns a NotFound error if the query returns nothing.
//func (rs *StmtRedux) LoadInt64(ctx context.Context, args ...interface{}) (value int64, err error) {
//	rs.lock()
//	defer rs.unlock()
//	if err = rs.rePrepare(); err != nil {
//		return 0, errors.WithStack(err)
//	}
//	value, err = rs.Stmt.LoadInt64(ctx, args...)
//	if err != nil {
//		err = errors.Wrapf(err, "[dml] StmtRedux.LoadInt64 with name %q", rs.name)
//	}
//	return
//}
//
//// LoadInt64s executes the Select and returns the value as a slice of
//// int64s.
//func (rs *StmtRedux) LoadInt64s(ctx context.Context, args ...interface{}) (values []int64, err error) {
//	rs.lock()
//	defer rs.unlock()
//	if err := rs.rePrepare(); err != nil {
//		return nil, errors.WithStack(err)
//	}
//	values, err = rs.Stmt.LoadInt64s(ctx, args...)
//	if err != nil {
//		err = errors.Wrapf(err, "[dml] StmtRedux.LoadInt64s with name %q", rs.name)
//	}
//	return
//}
