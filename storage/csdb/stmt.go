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
	"database/sql"
	"sync"
	"time"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// TODO(CyS) consider: http://dev.mysql.com/doc/refman/5.7/en/server-system-variables.html#sysvar_max_prepared_stmt_count
// TODO(CyS) analyze: Prepared Stmt and Query Cache http://dev.mysql.com/doc/refman/5.7/en/query-cache-operation.html
// TODO(CyS) implement to handle multiple statements, identified by a returned unique ID, so that you can grab that stmt.

// DefaultResurrectStmtIdleTime is the global idle time when you create a new
// PersistentStmt. If no query will be executed within this idle time the
// statement gets closed.
var DefaultResurrectStmtIdleTime = time.Second * 10

// ResurrectStmt creates a long living sql.Stmt in the database but closes it
// if within an idle time no query will be executed. Once there is a new
// query the statement gets resurrected. The ResurrectStmt type is safe for
// concurrent use with every of its function.
//
// A full working implementation can be found in package config with type DBStorage.
type ResurrectStmt struct {
	// DB contains for now only the prepare() function for a new statement
	// may be extended in the far future.
	db dbr.Preparer
	// qb is any prepareable SQL command, use ? for argument placeholders
	qb dbr.QueryBuilder
	// Idle defines the duration how to wait until no query will be executed.
	Idle time.Duration
	// Log default logger is PkgLof
	Log              log.Logger
	stop             chan struct{} // tells the ticker to stop and close
	idleCheckStarted bool

	mu       sync.Mutex // protects the last fields
	stmt     *sql.Stmt
	closed   bool      // stmt is closed and can be reopened
	lastUsed time.Time // time when the stmt has last been used
	inUse    bool      // stmt is currently in use by Set or Get
}

// NewResurrectStmt creates a new resurrected statement via a DB connection
// to prepare the stmt and a SQL query string. Default idle time is defined
// in DefaultResurrectStmtIdleTime. Default logger: PkgLog.
func NewResurrectStmt(p dbr.Preparer, qb dbr.QueryBuilder) *ResurrectStmt {
	// the overall question here is if the Stmt() function should
	// return an error once the ticker has been stopped or is not running.

	return &ResurrectStmt{
		db:     p,
		qb:     qb,
		Idle:   DefaultResurrectStmtIdleTime,
		Log:    log.BlackHole{},
		stop:   make(chan struct{}),
		closed: true,
	}
}

// StartIdleChecker starts the internal goroutine which checks the idle time.
// You can only start it once. sql.Stmt.Close() errors gets logged to Info. Those
// errors will only be returned if you stop the idle checker goroutine.
// Starting the idle checker is recommended because otherwise you might have
// a very long lived prepared statement.
func (su *ResurrectStmt) StartIdleChecker() {
	if su.idleCheckStarted {
		return
	}
	go su.checkIdle()
	su.idleCheckStarted = true
}

// StopIdleChecker stops the internal goroutine if it's started. Returns
// the sql.Stmt.Close error.
func (su *ResurrectStmt) StopIdleChecker() error {
	if su.idleCheckStarted {
		su.stop <- struct{}{}
	}
	su.idleCheckStarted = false
	return su.close()
}

// IsIdle returns true if the statement has been closed.
func (su *ResurrectStmt) IsIdle() bool {
	su.mu.Lock()
	defer su.mu.Unlock()
	return su.closed
}

func (su *ResurrectStmt) close() error {
	su.mu.Lock()
	defer func() {
		su.closed = true
		su.mu.Unlock()
	}()

	if su.Log.IsDebug() {
		sqlStr, _, err := su.qb.ToSQL()
		su.Log.Debug("csdb.ResurrectStmt.stmt.Close", log.String("SQL", sqlStr), log.ErrWithKey("ToSQL", err))
	}
	if su.stmt == nil {
		// statement has not been opened or is unused.
		return nil
	}
	return errors.Wrap(su.stmt.Close(), "[csdb] ResurrectStmt.close")
}

func (su *ResurrectStmt) checkIdle() {
	ticker := time.NewTicker(su.Idle)
	for {
		// maybe squeeze all three go routines into one. for each statement one select case.
		select {
		case t, ok := <-ticker.C:
			if !ok {
				return
			}

			if su.canClose(t) {
				// stmt has not been used within the last x seconds.
				// so close the stmt and release the resources in the DB.
				if err := su.close(); err != nil {
					sqlStr, _, err := su.qb.ToSQL()
					su.Log.Info("csdb.ResurrectStmt.stmt.Close.error", log.Err(err), log.String("SQL", sqlStr), log.ErrWithKey("ToSQL", err))
				}
			}
		case <-su.stop:
			ticker.Stop()
			return
		}
	}
}

func (su *ResurrectStmt) canClose(t time.Time) bool {
	su.mu.Lock()
	defer su.mu.Unlock()
	return t.After(su.lastUsed) && !su.closed && !su.inUse
}

// Stmt returns a prepared statement or an error. The statement gets
// automatically re-opened once it is closed after an idle time.
func (su *ResurrectStmt) Stmt() (*sql.Stmt, error) {
	su.mu.Lock()
	defer su.mu.Unlock()

	if !su.closed {
		return su.stmt, nil
	}
	sqlStr, _, err := su.qb.ToSQL()
	if err != nil {
		return nil, errors.Wrapf(err, "[csdb] DB.Prepare.ToSQL %q", sqlStr)
	}

	su.stmt, err = su.db.PrepareContext(context.Background(), sqlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "[csdb] DB.PrepareContext. %q", sqlStr)
	}
	if su.Log.IsDebug() {
		su.Log.Debug("csdb.ResurrectStmt.stmt.Prepare", log.String("SQL", sqlStr))
	}
	su.closed = false
	return su.stmt, nil
}

// StartStmtUse tells the ResurrectStmt type that Stmt() will be used.
func (su *ResurrectStmt) StartStmtUse() {
	su.mu.Lock()
	su.lastUsed = time.Now()
	su.inUse = true
	su.mu.Unlock()
}

// StopStmtUse tells the ResurrectStmt type that the Stmt() has been used.
func (su *ResurrectStmt) StopStmtUse() {
	su.mu.Lock()
	su.lastUsed = time.Now()
	su.inUse = false
	su.mu.Unlock()
}
