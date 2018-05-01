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

package ccd

import (
	"context"
	"sync"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/store/scope"
)

// Options applies options to the DBStorage type.
type Options struct {
	TableName           string
	Log                 log.Logger
	QueryContextTimeout time.Duration
	// IdleAllKeys default 5s
	IdleAllKeys time.Duration
	// IdleRead default 20s
	IdleRead time.Duration
	// IdleWrite default 10s
	IdleWrite             time.Duration
	ContextTimeoutAllKeys time.Duration
	ContextTimeoutRead    time.Duration
	ContextTimeoutWrite   time.Duration
	SkipSchemaValidation  bool
	// TODO implement UseDedicatedDBConnection per prepared statement, bit complicated
	// UseDedicatedDBConnection *sql.DB
}

const (
	stateClosed = iota // must be zero
	stateOpen
	stateInUse
)

type stats struct {
	Open  uint64
	Close uint64
}

// DBStorage connects the MySQL DB with the config.Service type. Implements
// interface config.Storager.
type DBStorage struct {
	cfg Options

	sqlAll   *dml.Select
	sqlRead  *dml.Select
	sqlWrite *dml.Insert

	tickerDaemonStop chan struct{}
	tickerAll        *time.Ticker
	tickerRead       *time.Ticker
	tickerWrite      *time.Ticker

	muAll        sync.Mutex
	stmtAll      *dml.Artisan
	stmtAllState uint8
	stmtAllStat  stats

	muRead        sync.Mutex
	stmtRead      *dml.Artisan
	stmtReadState uint8
	stmtReadStat  stats

	muWrite        sync.Mutex
	stmtWrite      *dml.Artisan
	stmtWriteState uint8
	stmtWriteStat  stats
}

// NewDBStorage creates a new database backed storage service. It creates three
// prepared statements which are getting automatically closed after an idle time
// and once used again, they get re-prepared. The database schema gets
// validated, but can also be disabled via options struct. Implements interface
// config.Storager.
func NewDBStorage(tbls *ddl.Tables, o Options) (*DBStorage, error) {
	// TODO reconfigure itself once it is running to load the timeout values from the DB and apply it. Restart the go routines.
	// Then during restart, block Set and Value operations to let the caller wait until restart complete. might be an overhead.

	tn := o.TableName
	if tn == "" {
		tn = TableNameCoreConfigData
	}

	if !o.SkipSchemaValidation {
		to := o.ContextTimeoutRead
		if to == 0 {
			to = time.Second * 10
		}
		ctx, cancel := context.WithTimeout(context.Background(), to)
		defer cancel()
		if err := tbls.Validate(ctx); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	tbl, err := tbls.Table(tn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qryAll := tbl.Select("scope", "scope_id", "path").OrderBy("scope", "scope_id", "path")
	qryAll.Log = o.Log

	qryRead := tbl.Select("value").Where(
		dml.Column("scope").PlaceHolder(),
		dml.Column("scope_id").PlaceHolder(),
		dml.Column("path").PlaceHolder(),
	)
	qryRead.Log = o.Log

	qryWrite := tbl.Insert().BuildValues()
	qryWrite.OnDuplicateKeys = dml.Conditions{dml.Column("value")}
	qryWrite.Log = o.Log

	dbs := &DBStorage{
		cfg:              o,
		tickerDaemonStop: make(chan struct{}),
		sqlAll:           qryAll,
		sqlRead:          qryRead,
		sqlWrite:         qryWrite,
	}
	if dbs.cfg.IdleAllKeys == 0 {
		dbs.cfg.IdleAllKeys = time.Second * 5 // just a guess
	}
	if dbs.cfg.IdleRead == 0 {
		dbs.cfg.IdleRead = time.Second * 20 // just a guess
	}
	if dbs.cfg.IdleWrite == 0 {
		dbs.cfg.IdleWrite = time.Second * 10 // just a guess
	}
	if dbs.cfg.ContextTimeoutAllKeys == 0 {
		dbs.cfg.ContextTimeoutAllKeys = time.Second * 10 // just a guess
	}
	if dbs.cfg.ContextTimeoutRead == 0 {
		dbs.cfg.ContextTimeoutRead = dbs.cfg.ContextTimeoutAllKeys
	}
	if dbs.cfg.ContextTimeoutWrite == 0 {
		dbs.cfg.ContextTimeoutWrite = dbs.cfg.ContextTimeoutAllKeys
	}
	dbs.runStateCheckers()
	return dbs, nil
}

// MustNewDBStorage same as NewDBStorage but panics on error. Implements
// interface config.Storager.
func MustNewDBStorage(tbls *ddl.Tables, o Options) *DBStorage {
	s, err := NewDBStorage(tbls, o)
	if err != nil {
		panic(err)
	}
	return s
}

func (dbs *DBStorage) closeStmtAll(t time.Time) {
	dbs.muAll.Lock()
	if dbs.stmtAllState == stateOpen {
		if err := dbs.stmtAll.Close(); err != nil && dbs.cfg.Log != nil && dbs.cfg.Log.IsInfo() {
			dbs.cfg.Log.Info("ccd.DBStorage.stmtAll.Close", log.Stringer("ticker", t), log.Err(err))
		}
		dbs.stmtAllState = stateClosed
		dbs.stmtAllStat.Close++
	}
	dbs.muAll.Unlock()
}

func (dbs *DBStorage) closeStmtRead(t time.Time) {
	dbs.muRead.Lock()
	if dbs.stmtReadState == stateOpen {
		if err := dbs.stmtRead.Close(); err != nil && dbs.cfg.Log != nil && dbs.cfg.Log.IsInfo() {
			dbs.cfg.Log.Info("ccd.DBStorage.stmtRead.Close", log.Stringer("ticker", t), log.Err(err))
		}
		dbs.stmtReadState = stateClosed
		dbs.stmtReadStat.Close++
	}
	dbs.muRead.Unlock()
}

func (dbs *DBStorage) closeStmtWrite(t time.Time) {
	dbs.muWrite.Lock()
	if dbs.stmtWriteState == stateOpen {
		if err := dbs.stmtWrite.Close(); err != nil && dbs.cfg.Log != nil && dbs.cfg.Log.IsInfo() {
			dbs.cfg.Log.Info("ccd.DBStorage.stmtWrite.Close", log.Stringer("ticker", t), log.Err(err))
		}
		dbs.stmtWriteState = stateClosed
		dbs.stmtWriteStat.Close++
	}
	dbs.muWrite.Unlock()
}

func (dbs *DBStorage) runStateCheckers() {

	dbs.tickerAll = time.NewTicker(dbs.cfg.IdleAllKeys)
	dbs.tickerRead = time.NewTicker(dbs.cfg.IdleRead)
	dbs.tickerWrite = time.NewTicker(dbs.cfg.IdleWrite)

	go func() {
		for {
			select {
			case <-dbs.tickerDaemonStop:
				dbs.tickerAll.Stop()
				dbs.tickerRead.Stop()
				dbs.tickerWrite.Stop()
				t := time.Now()
				dbs.closeStmtAll(t)
				dbs.closeStmtRead(t)
				dbs.closeStmtWrite(t)
				return
			case t := <-dbs.tickerAll.C:
				dbs.closeStmtAll(t)
			case t := <-dbs.tickerRead.C:
				dbs.closeStmtRead(t)
			case t := <-dbs.tickerWrite.C:
				dbs.closeStmtWrite(t)
			}
		}
	}()
}

// Close terminates the prepared statements and internal go routines.
func (dbs *DBStorage) Close() error {
	dbs.tickerDaemonStop <- struct{}{} // no need to close this chan because we might restart later the goroutine.
	return nil
}

// Set sets a value with its key. Database errors get logged as Info message.
// Enabled debug level logs the insert ID or rows affected.
func (dbs *DBStorage) Set(scp scope.TypeID, path string, value []byte) error {
	dbs.muWrite.Lock()
	prevState := dbs.stmtWriteState
	dbs.stmtWriteState = stateInUse
	defer func() {
		dbs.stmtWriteState = stateOpen
		dbs.muWrite.Unlock()
	}()

	ctx := context.Background()
	if prevState == stateClosed {
		ctx2, cancel := context.WithTimeout(ctx, dbs.cfg.ContextTimeoutWrite)
		defer cancel()
		stmt, err := dbs.sqlWrite.Prepare(ctx2)
		if err != nil {
			return errors.WithStack(err)
		}
		if dbs.stmtWrite == nil {
			dbs.stmtWrite = stmt.WithArgs()
		} else {
			dbs.stmtWrite.WithStmt(stmt.Stmt)
		}
		dbs.stmtWriteStat.Open++
	}

	ctx, cancel := context.WithTimeout(ctx, dbs.cfg.ContextTimeoutWrite)
	defer cancel()
	res, err := dbs.stmtWrite.Uint64(scp.ToUint64()).String(path).Bytes(value).ExecContext(ctx)
	dbs.stmtWrite.Reset()
	if dbs.cfg.Log != nil && dbs.cfg.Log.IsDebug() {
		li, err1 := res.LastInsertId()
		ra, err2 := res.RowsAffected()
		dbs.cfg.Log.Debug(
			"config.DBStorage.Set.Write.Result",
			log.Int64("lastInsertID", li),
			log.ErrWithKey("lastInsertIDErr", err1),
			log.Int64("rowsAffected", ra),
			log.ErrWithKey("rowsAffectedErr", err2),
			log.String("scope", scp.String()),
			log.String("path", path),
			log.Int("value_len", len(value)),
		)
	}

	return err
}

// Value performs a read operation from the database and returns a value from
// the table. The `ok` return argument can be true even if byte slice `v` is
// nil, which means that the path and scope are stored in the database table.
func (dbs *DBStorage) Value(scp scope.TypeID, path string) (v []byte, ok bool, err error) {
	dbs.muRead.Lock()
	prevState := dbs.stmtReadState
	dbs.stmtReadState = stateInUse
	defer func() {
		dbs.stmtReadState = stateOpen
		dbs.muRead.Unlock()
	}()

	ctx := context.Background()
	if prevState == stateClosed {
		ctx2, cancel := context.WithTimeout(ctx, dbs.cfg.ContextTimeoutRead)
		defer cancel()
		stmt, err := dbs.sqlRead.Prepare(ctx2)
		if err != nil {
			return nil, false, errors.WithStack(err)
		}
		if dbs.stmtRead == nil {
			dbs.stmtRead = stmt.WithArgs()
		} else {
			dbs.stmtRead.WithStmt(stmt.Stmt)
		}
		dbs.stmtReadStat.Open++
	}

	ctx, cancel := context.WithTimeout(ctx, dbs.cfg.ContextTimeoutRead)
	defer cancel()
	s, id := scp.Unpack()
	nv, found, err := dbs.stmtRead.String(s.StrType()).Int64(id).String(path).LoadNullString(ctx)
	if err != nil {
		return nil, false, errors.Wrapf(err, "[ccd] Scope %q Path %q", scp.String(), path)
	}
	if !found {
		return nil, false, nil
	}
	var ret []byte
	if nv.Valid {
		ret = []byte(nv.String)
	}
	return ret, true, nil
}

// AllKeys returns all available keys. Bot slices have the same length where
// index i of slice `scps` matches index j of slice paths.
func (dbs *DBStorage) AllKeys() (scps scope.TypeIDs, paths []string, err error) {
	dbs.muAll.Lock()
	prevState := dbs.stmtAllState
	dbs.stmtAllState = stateInUse
	defer func() {
		dbs.stmtAllState = stateOpen
		dbs.muAll.Unlock()
	}()

	ctx := context.Background()
	if prevState == stateClosed {
		ctx2, cancel := context.WithTimeout(ctx, dbs.cfg.ContextTimeoutAllKeys)
		defer cancel()
		var stmt *dml.Stmt
		stmt, err = dbs.sqlAll.Prepare(ctx2)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		if dbs.stmtAll == nil {
			dbs.stmtAll = stmt.WithArgs()
		} else {
			dbs.stmtAll.WithStmt(stmt.Stmt)
		}
		dbs.stmtAllStat.Open++
	}

	scps = make(scope.TypeIDs, 0, 100)
	paths = make([]string, 0, 100)

	ctx, cancel := context.WithTimeout(ctx, dbs.cfg.ContextTimeoutAllKeys)
	defer cancel()

	err = dbs.stmtAll.IterateSerial(ctx, func(cm *dml.ColumnMap) error {
		var ccd TableCoreConfigData
		if err2 := ccd.MapColumns(cm); err2 != nil {
			return errors.Wrapf(err2, "[ccd] dbs.stmtAll.IterateSerial at row %d", cm.Count)
		}
		scps = append(scps, scope.MakeTypeID(scope.FromString(ccd.Scope), ccd.ScopeID))
		paths = append(paths, ccd.Path)
		return nil
	})

	return scps, paths, err
}

// Statistics returns live statistics about opening and closing prepared statements.
func (dbs *DBStorage) Statistics() (value stats, set stats, all stats) {
	dbs.muAll.Lock()
	all = dbs.stmtAllStat
	dbs.muAll.Unlock()

	dbs.muRead.Lock()
	value = dbs.stmtReadStat
	dbs.muRead.Unlock()

	dbs.muWrite.Lock()
	set = dbs.stmtWriteStat
	dbs.muWrite.Unlock()

	return
}
