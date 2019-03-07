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

// +build csall db

//go:generate go run db_client_main.go

// TODO: only an idea to create program dmlgen
// go_generate dmlgen -tags "csall db" -filename $GOFILE -pkg $GOPACKAGE core_config_data

package storage

import (
	"context"
	"sync"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/store/scope"
)

// TODO https://mariadb.com/kb/en/library/system-versioned-tables/

// DBOptions applies options to the `DB` type.
type DBOptions struct {
	// TableName if set, specifies the alternate table name, default:
	// `core_config_data` aka constant TableNameCoreConfigData.
	TableName           string
	Log                 log.Logger
	QueryContextTimeout time.Duration
	// IdleRead default 20s
	IdleRead time.Duration
	// IdleWrite default 10s
	IdleWrite           time.Duration
	ContextTimeoutRead  time.Duration
	ContextTimeoutWrite time.Duration
	// SkipSchemaValidation disables the validation of the DB schema compared
	// with the schema stored in Go source files.
	SkipSchemaValidation bool
	// TODO implement UseDedicatedDBConnection per prepared statement, bit complicated
	// UseDedicatedDBConnection *sql.DB
}

const (
	stateClosed = iota // must be zero
	stateOpen
	stateInUse
)

type dbStats struct {
	Open  uint64
	Close uint64
}

// Service connects the MySQL/MariaDB with the config.Service type. Implements
// interface config.Storager.
type DB struct {
	cfg DBOptions

	sqlRead  *dml.Select
	sqlWrite *dml.Insert

	tickerDaemonStop chan struct{}
	tickerRead       *time.Ticker
	tickerWrite      *time.Ticker

	muRead        sync.Mutex
	stmtRead      *dml.Artisan
	stmtReadState uint8
	stmtReadStat  dbStats

	muWrite        sync.Mutex
	stmtWrite      *dml.Artisan
	stmtWriteState uint8
	stmtWriteStat  dbStats
}

// NewDB creates a new database backed storage service. It creates three
// prepared statements which are getting automatically closed after an idle time
// and once used again, they get re-prepared. The database schema gets
// validated, but can also be disabled via options struct. Implements interface
// config.Storager.
// NewDB uses the MySQL/MariaDB based table `core_config_data`
// for reading and writing configuration paths, scopes and values.
//
// It also provides an option function to load data from core_config_data into
// a storage service.
func NewDB(tbls *ddl.Tables, o DBOptions) (*DB, error) {
	// TODO reconfigure itself once it is running to load the timeout values
	// from the DB and apply it. Restart the go routines. Then during restart,
	// block Set and Value operations to let the caller wait until restart
	// complete. might be an overhead.

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

	dbs := &DB{
		cfg:              o,
		tickerDaemonStop: make(chan struct{}),
		sqlRead:          qryRead,
		sqlWrite:         qryWrite,
	}
	if dbs.cfg.IdleRead == 0 {
		dbs.cfg.IdleRead = time.Second * 20 // just a guess
	}
	if dbs.cfg.IdleWrite == 0 {
		dbs.cfg.IdleWrite = time.Second * 10 // just a guess
	}
	if dbs.cfg.ContextTimeoutRead == 0 {
		dbs.cfg.ContextTimeoutRead = time.Second * 10 // just a guess
	}
	if dbs.cfg.ContextTimeoutWrite == 0 {
		dbs.cfg.ContextTimeoutWrite = time.Second * 10 // just a guess
	}
	dbs.runStateCheckers()
	return dbs, nil
}

// MustNewDB same as NewDB but panics on error. Implements
// interface config.Storager.
func MustNewDB(tbls *ddl.Tables, o DBOptions) *DB {
	s, err := NewDB(tbls, o)
	if err != nil {
		panic(err)
	}
	return s
}

func (dbs *DB) closeStmtRead(t time.Time) {
	dbs.muRead.Lock()
	if dbs.stmtReadState == stateOpen {
		if err := dbs.stmtRead.Close(); err != nil && dbs.cfg.Log != nil && dbs.cfg.Log.IsInfo() {
			dbs.cfg.Log.Info("config.storage.DB.stmtRead.Close", log.Stringer("ticker", t), log.Err(err))
		}
		dbs.stmtReadState = stateClosed
		dbs.stmtReadStat.Close++
	}
	dbs.muRead.Unlock()
}

func (dbs *DB) closeStmtWrite(t time.Time) {
	dbs.muWrite.Lock()
	if dbs.stmtWriteState == stateOpen {
		if err := dbs.stmtWrite.Close(); err != nil && dbs.cfg.Log != nil && dbs.cfg.Log.IsInfo() {
			dbs.cfg.Log.Info("config.storage.DB.stmtWrite.Close", log.Stringer("ticker", t), log.Err(err))
		}
		dbs.stmtWriteState = stateClosed
		dbs.stmtWriteStat.Close++
	}
	dbs.muWrite.Unlock()
}

func (dbs *DB) runStateCheckers() {

	dbs.tickerRead = time.NewTicker(dbs.cfg.IdleRead)
	dbs.tickerWrite = time.NewTicker(dbs.cfg.IdleWrite)

	go func() {
		for {
			select {
			case <-dbs.tickerDaemonStop:
				dbs.tickerRead.Stop()
				dbs.tickerWrite.Stop()
				t := time.Now()
				dbs.closeStmtRead(t)
				dbs.closeStmtWrite(t)
				return
			case t := <-dbs.tickerRead.C:
				dbs.closeStmtRead(t)
			case t := <-dbs.tickerWrite.C:
				dbs.closeStmtWrite(t)
			}
		}
	}()
}

// Close terminates the prepared statements and internal go routines.
func (dbs *DB) Close() error {
	dbs.tickerDaemonStop <- struct{}{} // no need to close this chan because we might restart later the goroutine.
	return nil
}

// Set puts a value with its key. Database errors get logged as Info message.
// Enabled debug level logs the insert ID or rows affected.
func (dbs *DB) Set(p *config.Path, value []byte) error {
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
			dbs.stmtWrite.WithPreparedStmt(stmt.Stmt)
		}
		dbs.stmtWriteStat.Open++
	}

	ctx, cancel := context.WithTimeout(ctx, dbs.cfg.ContextTimeoutWrite)
	defer cancel()
	scp, path := p.ScopeRoute()
	res, err := dbs.stmtWrite.Uint64(scp.ToUint64()).String(path).Bytes(value).ExecContext(ctx)
	dbs.stmtWrite.Reset()
	if dbs.cfg.Log != nil && dbs.cfg.Log.IsDebug() {
		li, err1 := res.LastInsertId()
		ra, err2 := res.RowsAffected()
		dbs.cfg.Log.Debug(
			"config.storage.DB.Set.Write.Result",
			log.Int64("lastInsertID", li),
			log.ErrWithKey("lastInsertIDErr", err1),
			log.Int64("rowsAffected", ra),
			log.ErrWithKey("rowsAffectedErr", err2),
			log.String("path", p.String()),
			log.Int("value_len", len(value)),
		)
	}

	return err
}

// Get performs a read operation from the database and returns a value from
// the table. The `ok` return argument can be true even if byte slice `v` is
// nil, which means that the path and scope are stored in the database table.
func (dbs *DB) Get(p *config.Path) (v []byte, ok bool, err error) {
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
			dbs.stmtRead.WithPreparedStmt(stmt.Stmt)
		}
		dbs.stmtReadStat.Open++
	}

	ctx, cancel := context.WithTimeout(ctx, dbs.cfg.ContextTimeoutRead)
	defer cancel()
	scp, path := p.ScopeRoute()
	s, id := scp.Unpack()
	nv, found, err := dbs.stmtRead.String(s.StrType()).Int64(id).String(path).LoadNullString(ctx)
	if err != nil {
		return nil, false, errors.Wrapf(err, "[config/storage] DB Scope %q Path %q", scp.String(), path)
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

// Statistics returns live statistics about opening and closing prepared statements.
func (dbs *DB) Statistics() (value dbStats, set dbStats) {
	dbs.muRead.Lock()
	value = dbs.stmtReadStat
	dbs.muRead.Unlock()

	dbs.muWrite.Lock()
	set = dbs.stmtWriteStat
	dbs.muWrite.Unlock()

	return
}

// WithLoadFromDB reads the table core_config_data into the Service and
// overrides existing values. Stops on errors.
func WithLoadFromDB(tbls *ddl.Tables, o DBOptions) config.LoadDataOption {
	return config.MakeLoadDataOption(func(s *config.Service) error {

		tn := o.TableName
		if tn == "" {
			tn = TableNameCoreConfigData
		}

		tbl, err := tbls.Table(tn)
		if err != nil {
			return errors.WithStack(err)
		}

		if o.ContextTimeoutRead == 0 {
			o.ContextTimeoutRead = time.Second * 10 // just a guess
		}

		ctx, cancel := context.WithTimeout(context.Background(), o.ContextTimeoutRead)
		defer cancel()

		return tbl.Select("*").WithArgs().IterateSerial(ctx, func(cm *dml.ColumnMap) error {
			var ccd CoreConfigData
			if err := ccd.MapColumns(cm); err != nil {
				return errors.Wrapf(err, "[config/storage] dbs.stmtAll.IterateSerial at row %d", cm.Count)
			}

			var v []byte
			if ccd.Value.Valid {
				v = []byte(ccd.Value.String)
			}
			scp := scope.FromString(ccd.Scope).WithID(int64(ccd.ScopeID))
			p, err := config.NewPathWithScope(scp, ccd.Path)
			if err != nil {
				return errors.Wrapf(err, "[config/storage] WithLoadFromDB.config.NewPathWithScope Path %q Scope: %q ID: %d", ccd.Path, scp, ccd.ConfigID)
			}
			if err = s.Set(p, v); err != nil {
				return errors.Wrapf(err, "[config/storage] WithLoadFromDB.Service.Write Path %q Scope: %q ID: %d", ccd.Path, scp, ccd.ConfigID)
			}

			return nil
		})
	}).WithUseStorageLevel(1)
}
