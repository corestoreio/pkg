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
	"time"

	"database/sql"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/store/scope"
	"sync"
)

// Options applies options to the DBStorage type.
type Options struct {
	TableName             string
	Log                   log.Logger
	QueryContextTimeout   time.Duration
	IdleAllKeys           time.Duration
	IdleRead              time.Duration
	IdleWrite             time.Duration
	ContextTimeoutAllKeys time.Duration
	ContextTimeoutRead    time.Duration
	ContextTimeoutWrite   time.Duration
}

const (
	stateClosed = iota // must be zero
	stateOpen
	stateInUse
)

// DBStorage connects the MySQL DB with the config.Service type. Implements
// interface config.Storager.
type DBStorage struct {
	cfg Options

	sqlAll   *dml.Select
	sqlRead  *dml.Select
	sqlWrite *dml.Insert

	tickerAll   *time.Ticker
	tickerRead  *time.Ticker
	tickerWrite *time.Ticker

	muAll        sync.Mutex
	stmtAll      *dml.Artisan
	stmtAllState uint8

	muRead        sync.Mutex
	stmtRead      *dml.Artisan
	stmtReadState uint8

	muWrite        sync.Mutex
	stmtWrite      *dml.Artisan
	stmtWriteState uint8
}

// NewDBStorage creates a new database backed storage service. It creates three
// prepared statements which are getting automatically closed after an idle time
// and once used again they get re-prepared.
//
// All has an idle time of 15s. Read an idle time of 10s. Write an idle time of
// 30s. Implements interface config.Storager.
func NewDBStorage(tbls *ddl.Tables, o Options) (*DBStorage, error) {

	tn := o.TableName
	if tn == "" {
		tn = TableNameCoreConfigData
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

	qryWrite := tbl.Insert()
	qryWrite.OnDuplicateKeys = dml.Conditions{
		dml.Column("value"),
	}
	qryWrite.Log = o.Log

	dbs := &DBStorage{
		cfg:      o,
		sqlAll:   qryAll,
		sqlRead:  qryRead,
		sqlWrite: qryWrite,
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
	dbs.tickerAll = time.NewTicker(dbs.cfg.IdleAllKeys)
	dbs.tickerRead = time.NewTicker(dbs.cfg.IdleRead)
	dbs.tickerWrite = time.NewTicker(dbs.cfg.IdleWrite)

	go func() {
		for {
			select {
			case t, ok := <-dbs.tickerAll.C:
				if !ok {
					return
				}
				dbs.muAll.Lock()
				if dbs.stmtAllState == stateOpen {
					if err := dbs.stmtAll.Close(); err != nil && dbs.cfg.Log != nil && dbs.cfg.Log.IsInfo() {
						dbs.cfg.Log.Info("ccd.DBStorage.stmtAll.Close", log.Stringer("ticker", t), log.Err(err))
					}
					dbs.stmtAllState = stateClosed
				}
				dbs.muAll.Unlock()
			case t, ok := <-dbs.tickerRead.C:
				if !ok {
					return
				}
				dbs.muRead.Lock()
				if dbs.stmtReadState == stateOpen {
					if err := dbs.stmtRead.Close(); err != nil && dbs.cfg.Log != nil && dbs.cfg.Log.IsInfo() {
						dbs.cfg.Log.Info("ccd.DBStorage.stmtRead.Close", log.Stringer("ticker", t), log.Err(err))
					}
					dbs.stmtReadState = stateClosed
				}
				dbs.muRead.Unlock()
			case t, ok := <-dbs.tickerWrite.C:
				if !ok {
					return
				}
				dbs.muWrite.Lock()
				if dbs.stmtWriteState == stateOpen {
					if err := dbs.stmtWrite.Close(); err != nil && dbs.cfg.Log != nil && dbs.cfg.Log.IsInfo() {
						dbs.cfg.Log.Info("ccd.DBStorage.stmtWrite.Close", log.Stringer("ticker", t), log.Err(err))
					}
					dbs.stmtWriteState = stateClosed
				}
				dbs.muWrite.Unlock()
			}
		}
	}()

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

// Close terminates the prepared statements and internal go routines.
func (dbs *DBStorage) Close() error {
	dbs.tickerAll.Stop()
	dbs.tickerRead.Stop()
	dbs.tickerWrite.Stop()
	return nil
}

// Set sets a value with its key. Database errors get logged as Info message.
// Enabled debug level logs the insert ID or rows affected.
func (dbs *DBStorage) Set(scp scope.TypeID, path string, value []byte) (err error) {
	dbs.muAll.Lock()
	prevState := dbs.stmtWriteState
	dbs.stmtWriteState = stateInUse
	defer func() {
		dbs.stmtWriteState = stateOpen
		dbs.muAll.Unlock()
	}()

	if prevState == stateClosed {
		ctx, cancel := context.WithTimeout(context.Background(), dbs.cfg.ContextTimeoutWrite)
		defer cancel()
		var stmt *dml.Stmt
		stmt, err = dbs.sqlWrite.Prepare(ctx)
		if err != nil {
			return errors.WithStack(err)
		}
		if dbs.stmtWrite == nil {
			dbs.stmtWrite = stmt.WithArgs()
		} else {
			dbs.stmtWrite.WithStmt(stmt.Stmt)
		}
	}

	var res sql.Result
	res, err = dbs.stmtWrite.Uint64(scp.ToUint64()).String(path).Bytes(value).ExecContext(context.Background())

	//if dbs.log.IsDebug() {
	//	li, err1 := result.LastInsertId()
	//	ra, err2 := result.RowsAffected()
	//	dbs.log.Debug(
	//		"config.DBStorage.Set.Write.Result",
	//		log.Int64("lastInsertID", li),
	//		log.ErrWithKey("lastInsertIDErr", err1),
	//		log.Int64("rowsAffected", ra),
	//		log.ErrWithKey("rowsAffectedErr", err2),
	//		log.String("SQL", dbs.Write.sqlRaw),
	//		log.Stringer("key", key),
	//		log.Object("value", value),
	//	)
	//}

	return nil
}

// Get returns a value from the database by its key. It is guaranteed that the
// type in the empty interface is a string. It returns nil on error but errors
// get logged as info message. Error behaviour: NotFound
func (dbs *DBStorage) Value(scp scope.TypeID, path string) (v []byte, ok bool, err error) {
	return

	// update lastUsed at the end because there might be the slight chance that
	// a statement gets closed despite we're waiting for the result from the
	// server.
	//dbs.Read.StartStmtUse()
	//defer dbs.Read.StopStmtUse()
	//
	//stmt, err := dbs.Read.Stmt(context.TODO())
	//if err != nil {
	//	return nil, errors.Wrapf(err, "[ccd] Get.Read.Stmt. SQL: %q Key: %q", dbs.Read.sqlRaw, key)
	//}
	//
	//pl, err := key.Level(-1)
	//if err != nil {
	//	return nil, errors.Wrapf(err, "[ccd] Get.key.Level. SQL: %q Key: %q", dbs.Read.sqlRaw, key)
	//}
	//
	//var data null.String
	//scp, id := key.ScopeID.Unpack()
	//err = stmt.QueryRow(scp.StrType(), id, pl).Scan(&data)
	//if err != nil {
	//	return nil, errors.Wrapf(err, "[ccd] Get.QueryRow. SQL: %q Key: %q PathLevel: %q", dbs.Read.sqlRaw, key, pl)
	//}
	//if data.Valid {
	//	return data.String, nil
	//}
	//return nil, errKeyNotFound
}

// AllKeys returns all available keys. Database errors get logged as info message.
func (dbs *DBStorage) AllKeys() (scps scope.TypeIDs, paths []string, err error) {
	dbs.muAll.Lock()
	prevState := dbs.stmtAllState
	dbs.stmtAllState = stateInUse
	defer func() {
		dbs.stmtAllState = stateOpen
		dbs.muAll.Unlock()
	}()

	if prevState == stateClosed {
		ctx, cancel := context.WithTimeout(context.Background(), dbs.cfg.ContextTimeoutAllKeys)
		defer cancel()
		var stmt *dml.Stmt
		stmt, err = dbs.sqlAll.Prepare(ctx)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		if dbs.stmtAll == nil {
			dbs.stmtAll = stmt.WithArgs()
		} else {
			dbs.stmtAll.WithStmt(stmt.Stmt)
		}
	}

	scps = make(scope.TypeIDs, 0, 100)
	paths = make([]string, 0, 100)

	dbs.stmtAll.IterateSerial(context.Background(), func(cm *dml.ColumnMap) error {
		var ccd TableCoreConfigData
		if err := ccd.MapColumns(cm); err != nil {
			return errors.Wrapf(err, "[ccd] dbs.stmtAll.IterateSerial at row %d", cm.Count)
		}
		scps = append(scps, scope.MakeTypeID(scope.FromString(ccd.Scope), ccd.ScopeID))
		paths = append(paths, ccd.Path)
		return nil
	})

	return scps, paths, nil
}
