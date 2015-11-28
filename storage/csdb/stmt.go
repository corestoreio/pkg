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
	"sync"
	"time"

	"database/sql"

	"github.com/corestoreio/csfw/utils/log"
	"github.com/juju/errgo"
)

// DefaultResurrectStmtIdleTime is the global idle time when you create a new
// PersistentStmt. If no query will be executed within this idle time the
// statement gets closed.
var DefaultResurrectStmtIdleTime = time.Second * 10

// NewResurrectStmt creates a new resurrected statement via a DB connection
// to prepare the stmt and a SQL query string. Default idle time is defined
// in DefaultResurrectStmtIdleTime. Default logger: PkgLog.
func NewResurrectStmt(db *sql.DB, SQL string) *ResurrectStmt {
	return &ResurrectStmt{
		DB:       db,
		SQL:      SQL,
		Idle:     DefaultResurrectStmtIdleTime,
		Log:      PkgLog,
		stop:     make(chan struct{}),
		closeErr: make(chan error),
		closed:   true,
	}
}

// ResurrectStmt creates a long living sql.Stmt in the database but closes it
// if within an idle time no query will be executed. Once there is a new
// query the statement gets resurrected. The ResurrectStmt type is safe for
// concurrent use with every of its function.
//
// A full working implementation can be found in package config with type DBStorage.
type ResurrectStmt struct {
	DB *sql.DB
	// SQL is any prepareable SQL command, use ? for argument placeholders
	SQL string
	// Idle defines the duration how to wait until no query will be executed.
	Idle time.Duration
	// Log default logger is PkgLof
	Log              log.Logger
	stop             chan struct{} // tells the ticker to stop and close
	idleCheckStarted bool

	mu       sync.Mutex // protects the last fields
	stmt     *sql.Stmt
	closed   bool       // stmt is closed and can be reopened
	closeErr chan error // only available when Stop() has been called
	lastUsed time.Time  // time when the stmt has last been used
	inUse    bool       // stmt is currently in use by Set or Get
}

func (su *ResurrectStmt) close(retErr bool) {
	// retErr returns only then the error when the main go routine of the ticker
	// has been stopped. otherwise close errors will only be logged.
	su.mu.Lock()
	defer su.mu.Unlock()
	if su.stmt == nil {
		su.closeErr <- nil
		return
	}
	err := errgo.Mask(su.stmt.Close())
	if err != nil {
		su.Log.Info("csdb.ResurrectStmt.stmt.Close.error", "err", err, "SQL", su.SQL)
	} else {
		su.closed = true
	}
	if retErr {
		su.closeErr <- err
	}
	if su.Log.IsDebug() {
		su.Log.Debug("csdb.ResurrectStmt.stmt.Close", "SQL", su.SQL)
	}
}

func (su *ResurrectStmt) canClose(t time.Time) bool {
	su.mu.Lock()
	defer su.mu.Unlock()
	return t.After(su.lastUsed) && !su.closed && !su.inUse
}

// StartIdleChecker starts the internal goroutine which checks the idle time.
// You can only start it once. sql.Stmt.Close() errors gets logged to Info. Those
// errors will only be returned if you stop the idle checker goroutine.
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
	return <-su.closeErr
}

func (su *ResurrectStmt) checkIdle() {
	ticker := time.NewTicker(su.Idle)
	for {
		// maybe squeeze all three go routines into one. for each statement one select case.
		select {
		case t, ok := <-ticker.C:
			if !ok {
				// todo maybe debug log?
				return
			}

			if su.canClose(t) {
				// stmt has not been used within the last x seconds.
				// so close the stmt and release the resources in the DB.
				su.close(false)
			}
		case <-su.stop:
			ticker.Stop()
			su.close(true)
			return
		}
	}
}

// Stmt returns a prepared statement or an error. The statement gets
// automatically re-opened once it's closed after an idle time.
func (su *ResurrectStmt) Stmt() (*sql.Stmt, error) {
	su.mu.Lock()
	defer su.mu.Unlock()

	if false == su.closed {
		return su.stmt, nil
	}

	var err error
	su.stmt, err = su.DB.Prepare(su.SQL)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	if su.Log.IsDebug() {
		su.Log.Debug("csdb.ResurrectStmt.stmt.Prepare", "SQL", su.SQL)
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
