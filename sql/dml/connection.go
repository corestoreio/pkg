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
	"database/sql/driver"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/modl"
)

type uniqueIDFn func() string

func uniqueIDNoOp() string { return "" }

type logWithID struct {
	start time.Time
	Log   log.Logger
	// makeUniqueID generates for each call a new unique ID. Those IDs will be
	// assigned to a new connection or a new statement. The function signature is
	// equal to fmt.Stringer so one can use for example:
	//		uuid.NewV4().String
	// The returned unique ID gets used in logging and inserted as a comment
	// into the SQL string. The returned string must not contain the
	// comment-end-termination pattern: `*/`.
	makeUniqueID uniqueIDFn
}

// ConnPool at a connection to the database with an EventReceiver to send
// events, errors, and timings to
type ConnPool struct {
	logWithID
	DB  *sql.DB
	dsn string

	rwmu          sync.RWMutex
	preparedStmts map[string]*StmtRedux
}

// Conn represents a single database session rather a pool of database sessions.
// Prefer running queries from DB unless there is a specific need for a
// continuous single database session.
//
// A Conn must call Close to return the connection to the database pool and may
// do so concurrently with a running query.
//
// After a call to Close, all operations on the connection fail with
// ErrConnDone.
type Conn struct {
	logWithID
	DB *sql.Conn
}

// Tx is an in-progress database transaction.
//
// A transaction must end with a call to Commit or Rollback.
//
// After a call to Commit or Rollback, all operations on the transaction fail
// with ErrTxDone.
//
// The statements prepared for a transaction by calling the transaction's
// Prepare or Stmt methods are closed by the call to Commit or Rollback.
// Practical Guide to SQL Transaction Isolation: https://begriffs.com/posts/2017-08-01-practical-guide-sql-isolation.html
type Tx struct {
	logWithID
	DB *sql.Tx
}

// ConnPoolOption can be used at an argument in NewConnPool to configure a
// connection.
type ConnPoolOption struct {
	sortOrder uint8
	fn        func(*ConnPool) error
	// WithUniqueIDFn applies a unique ID generator function without an applied
	// logger as in WithLogger. For more details see WithLogger function.
	// Sort Order 8.
	UniqueIDFn func() string
	// TableNameMapper maps the old name in the DML query to a new name. E.g.
	// for adding a prefix and/or a suffix.
	// TODO implement TableNameMapper
	TableNameMapper func(oldName string) (newName string)
	// OptimisticLock is enabled all queries with Exec will have a `version` field.
	// UPDATE user SET ..., version = version + 1 WHERE id = ? AND version = ?
	// TODO implement OptimisticLock
	OptimisticLock bool
	// OptimisticLockFieldName custom global column name, defaults to `version uint64`
	OptimisticLockColumnName string
	A                        gorp.OptimisticLockError
	B                        modl.OptimisticLockError
}

// WithLogger sets the customer logger to be used across the package. The logger
// gets inherited to type Conn and Tx and also to all statement types. Each
// heredity creates new fields as a prefix. Argument `uniqueID` generates for
// each heredity a new unique ID for tracing in Info logging. Those IDs will be
// assigned to a new connection or a new statement. The function signature is
// equal to fmt.Stringer so one can use for example:
//		uuid.NewV4().String
// The returned unique ID from `uniqueIDFn` gets used in logging and inserted as
// a comment into the SQL string for tracing in server log files and PROCESS
// LIST. The returned string must not contain the comment-end-termination
// pattern: `*/`. The `uniqueIDFn` must be thread safe.
func WithLogger(l log.Logger, uniqueIDFn func() string) ConnPoolOption {
	return ConnPoolOption{
		sortOrder: 10,
		fn: func(c *ConnPool) error {
			c.makeUniqueID = uniqueIDFn
			c.Log = l.With(log.String("conn_pool_id", c.makeUniqueID()))
			return nil
		},
	}
}

// WithDB sets the DB value to an existing connection. Mainly used for testing.
// Does not support DriverCallBack.
func WithDB(db *sql.DB) ConnPoolOption {
	return ConnPoolOption{
		sortOrder: 1,
		fn: func(c *ConnPool) error {
			c.DB = db
			return nil
		},
	}
}

// WithPreparedStatement uses the name as unique identifier to create a
// redux/resurrectable prepared statement. This prepared statement is bound to a
// single connection. After an idle time the statement and the connection gets
// closed and resources on the DB server freed. Once the statement will be used
// again, it reduxes, gets re-prepared, via its dedicated connection. If there
// are no more connections available, it waits until it can connect. Due to the
// connection handling implementation of database/sql/DB object we must grab a
// dedicated connection. To call this statement use function ConnPool.Stmt.
func WithPreparedStatement(name string, qb QueryBuilder, idleTime time.Duration) ConnPoolOption {
	// last parameter of above API signature might turn into a config struct
	// because we will have more options for newReduxStmt.
	return ConnPoolOption{
		sortOrder: 200,
		fn: func(c *ConnPool) error {
			c.rwmu.Lock()
			defer c.rwmu.Unlock()
			if c.preparedStmts == nil {
				c.preparedStmts = make(map[string]*StmtRedux, 20)
			}

			if rs, ok := c.preparedStmts[name]; ok {
				if err := rs.Close(); err != nil {
					return errors.Wrapf(err, "[dml] WithPrepareStatements failed for %q", name)
				}
			}

			rs, err := newReduxStmt(c.DB, name, qb, idleTime, c.Log)
			if err != nil {
				return errors.WithStack(err)
			}
			c.preparedStmts[name] = rs
			// check max_prepared_stmt_count and con count
			// TODO(CyS) consider: http://dev.mysql.com/doc/refman/5.7/en/server-system-variables.html#sysvar_max_prepared_stmt_count

			return nil
		},
	}
}

// WithDSN sets the data source name for a connection.
// Second argument DriverCallBack adds a low level call back function on MySQL driver level to
// create a a new instrumented driver. No need to call `sql.Register`!
func WithDSN(dsn string, cb ...DriverCallBack) ConnPoolOption {
	if len(cb) > 1 {
		panic(errors.NotImplemented.Newf("[dml] Only one DriverCallBack function does currently work. You provided: %d", len(cb)))
	}
	return ConnPoolOption{
		sortOrder: 0,
		fn: func(c *ConnPool) error {
			if !strings.Contains(dsn, "parseTime") {
				return errors.NotImplemented.Newf("[dml] The DSN for go-sql-driver/mysql must contain the parameters `?parseTime=true[&loc=YourTimeZone]`")
			}
			c.dsn = dsn
			var drv driver.Driver = mysql.MySQLDriver{}
			if len(cb) == 1 {
				drv = wrapDriver(drv, cb[0])
			}
			c.DB = sql.OpenDB(dsnConnector{dsn: dsn, driver: drv})
			return nil
		},
	}
}

// dsnConnector implements a type to open a connection to the DB. It makes the
// call to sql.Register superfluous.
type dsnConnector struct {
	dsn    string
	driver driver.Driver
}

func (t dsnConnector) Connect(_ context.Context) (driver.Conn, error) {
	return t.driver.Open(t.dsn)
}

func (t dsnConnector) Driver() driver.Driver {
	return t.driver
}

// NewConnPool instantiates a ConnPool for a given database/sql connection
// and event receiver. An invalid driver name causes a NotImplemented error to be
// returned. You can either apply a DSN or a pre configured *sql.DB type. For
// full UTF-8 support you must set the charset in the SQL driver to utf8mb4.
func NewConnPool(opts ...ConnPoolOption) (*ConnPool, error) {
	c := &ConnPool{}
	if err := c.Options(opts...); err != nil {
		return nil, errors.WithStack(err)
	}

	if c.makeUniqueID == nil {
		c.makeUniqueID = uniqueIDNoOp
	}
	// validate that DSN contains the utf8mb4 setting

	// TODO: Validate that we run with utf8mb4 the normal utf8 is only 3 bytes
	// where utf8mb4 is full 4byte support.
	// SHOW VARIABLES WHERE Variable_name LIKE 'character\_set\_%' OR Variable_name LIKE 'collation%';
	// TODO: Set SQL mode to strict https://dev.mysql.com/doc/refman/5.7/en/sql-mode.html#sql-mode-strict

	return c, nil
}

// MustConnectAndVerify at like NewConnPool but it verifies the connection
// and panics on errors.
func MustConnectAndVerify(opts ...ConnPoolOption) *ConnPool {
	c, err := NewConnPool(opts...)
	if err != nil {
		panic(err)
	}
	if err := c.DB.Ping(); err != nil {
		panic(err)
	}
	return c
}

// Options applies options to a connection
func (c *ConnPool) Options(opts ...ConnPoolOption) error {

	for i, opt := range opts {
		if opt.UniqueIDFn != nil {
			opts[i].sortOrder = 8
			opt := opt
			opts[i].fn = func(cp *ConnPool) error {
				cp.makeUniqueID = opt.UniqueIDFn
				return nil
			}
		}
	}

	// SliceStable must be stable to maintain the order of all options where
	// sortOrder is zero.
	sort.SliceStable(opts, func(i, j int) bool {
		return opts[i].sortOrder < opts[j].sortOrder
	})

	for _, opt := range opts {
		if err := opt.fn(c); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Close closes the database, releasing any open resources.
//
// It is rare to Close a DB, as the DB handle is meant to be long-lived and
// shared between many goroutines. It logs the time taken, if a logger has been
// set with Info logging enabled.
func (c *ConnPool) Close() error {
	if c.Log != nil && c.Log.IsDebug() {
		defer c.Log.Debug("Close", log.Duration("duration", now().Sub(c.start)))
	}
	c.rwmu.Lock()
	defer c.rwmu.Unlock()
	for _, rs := range c.preparedStmts {
		if err := rs.Close(); err != nil {
			return errors.WithStack(err)
		}
	}
	return c.DB.Close() // no stack wrap otherwise error is hard to compare
}

// BeginTx starts a transaction.
//
// The provided context is used until the transaction is committed or rolled
// back. If the context is canceled, the sql package will roll back the
// transaction. Tx.Commit will return an error if the context provided to
// BeginTx is canceled.
//
// The provided TxOptions is optional and may be nil if defaults should be used.
// If a non-default isolation level is used that the driver doesn't support, an
// error will be returned.
//
// Practical Guide to SQL Transaction Isolation: https://begriffs.com/posts/2017-08-01-practical-guide-sql-isolation.html
func (c *ConnPool) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	start := now()

	dbTx, err := c.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	l := c.Log
	if l != nil {
		l = l.With(log.String("tx_id", c.makeUniqueID()))
		if l.IsDebug() {
			l.Debug("BeginTx")
		}
	}
	return &Tx{
		logWithID: logWithID{
			start:        start,
			Log:          l,
			makeUniqueID: c.makeUniqueID,
		},
		DB: dbTx,
	}, nil
}

// TODO add this to all connection types, like Conn and Tx
func (c *ConnPool) WithQueryBuilder(qb QueryBuilder) *Arguments {
	sqlStr, argsRaw, err := qb.ToSQL()
	var args [defaultArgumentsCapacity]argument
	return &Arguments{
		base: builderCommon{
			cachedSQL: []byte(sqlStr),
			Log:       c.Log,
			id:        c.makeUniqueID(),
			DB:        c.DB,
			ärgErr:    errors.WithStack(err),
		},
		raw:  argsRaw,
		args: args[:0],
	}
}

// StmtRedux returns a redux prepared statement. The object represents a
// prepared statement which can resurrect/redux and is bound to a single
// connection. After an idle time the statement and the connection gets closed
// and resources on the DB server freed. Once the statement will be used again,
// it reduxes, gets re-prepared, via its dedicated connection. If there are no
// more connections available, it waits until it can connect. Due to the
// connection handling implementation of database/sql/DB object we must grab a
// dedicated connection. A later aim would that multiple prepared statements can
// share a single connection.
func (c *ConnPool) StmtRedux(name string) (*StmtRedux, error) {
	c.rwmu.RLock()
	defer c.rwmu.RUnlock()

	rs, ok := c.preparedStmts[name]
	if !ok {
		return nil, errors.NotFound.Newf("[dml] Stmt %q not found", name)
	}

	return rs, nil
}

// StmtReduxPrepare same as functional option WithPreparedStatement but returns
// the lazy prepared Stmter.
func (c *ConnPool) StmtReduxPrepare(name string, qb QueryBuilder, idleTime time.Duration) (*StmtRedux, error) {
	c.rwmu.RLock()
	_, ok := c.preparedStmts[name]
	c.rwmu.RUnlock()
	if !ok {
		if err := WithPreparedStatement(name, qb, idleTime).fn(c); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return c.StmtRedux(name)
}

// Conn returns a single connection by either opening a new connection
// or returning an existing connection from the connection pool. Conn will
// block until either a connection is returned or ctx is canceled.
// Queries run on the same Conn will be run in the same database session.
//
// Every Conn must be returned to the database pool after use by
// calling Conn.Close.
func (c *ConnPool) Conn(ctx context.Context) (*Conn, error) {
	dbc, err := c.DB.Conn(ctx)
	l := c.Log
	if l != nil {
		l = c.Log.With(log.String("conn_id", c.makeUniqueID()))
	}
	return &Conn{
		logWithID: logWithID{
			start:        now(),
			Log:          l,
			makeUniqueID: c.makeUniqueID,
		},
		DB: dbc,
	}, errors.WithStack(err)
}

// BeginTx starts a transaction.
//
// The provided context is used until the transaction is committed or rolled back.
// If the context is canceled, the sql package will roll back
// the transaction. Tx.Commit will return an error if the context provided to
// BeginTx is canceled.
//
// The provided TxOptions is optional and may be nil if defaults should be used.
// If a non-default isolation level is used that the driver doesn't support,
// an error will be returned.
func (c *Conn) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	start := now()

	dbTx, err := c.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	l := c.Log
	if l != nil {
		l = l.With(log.String("tx_id", c.makeUniqueID()))
		if l.IsDebug() {
			l.Debug("BeginTx")
		}
	}
	return &Tx{
		logWithID: logWithID{
			start:        start,
			Log:          l,
			makeUniqueID: c.makeUniqueID,
		},
		DB: dbTx,
	}, nil
}

// Close returns the connection to the connection pool. All operations after a
// Close will return with ErrConnDone. Close is safe to call concurrently with
// other operations and will block until all other operations finish. It may be
// useful to first cancel any used context and then call close directly after.
// It logs the time taken, if a logger has been set with Info logging enabled.
func (c *Conn) Close() error {
	if c.Log != nil && c.Log.IsDebug() {
		defer c.Log.Debug("Close", log.Duration("duration", now().Sub(c.start)))
	}
	return c.DB.Close() // no stack wrap otherwise error is hard to compare
}

// Commit finishes the transaction. It logs the time taken, if a logger has been
// set with Info logging enabled.
func (tx *Tx) Commit() error {
	if tx.Log != nil && tx.Log.IsDebug() {
		defer tx.Log.Debug("Commit", log.Duration("duration", now().Sub(tx.start)))
	}
	return tx.DB.Commit()
}

// Rollback cancels the transaction. It logs the time taken, if a logger has
// been set with Info logging enabled.
func (tx *Tx) Rollback() error {
	if tx.Log != nil && tx.Log.IsDebug() {
		defer tx.Log.Debug("Rollback", log.Duration("duration", now().Sub(tx.start)))
	}
	return tx.DB.Rollback()
}

// Wrap is a helper method that will automatically COMMIT or ROLLBACK once the
// supplied functions are done executing.
//
//      tx, err := db.Begin()
//      if err != nil{
//           panic(err.Error()) // you could gracefully handle the error also
//      }
//      if err := tx.Wrap(func() error {
//          // SQL
//          return nil
//      }); err != nil{
//           panic(err.Error()) // you could gracefully handle the error also
//      }
// It logs the time taken, if a logger has been set with Info logging enabled.
func (tx *Tx) Wrap(fns ...func() error) error {
	for i, f := range fns {
		if err := f(); err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				return errors.Wrapf(rErr, "[dml] transaction.wrap.Rollback.error at index %d", i)
			}
			return errors.Wrapf(err, "[dml] transaction.wrap.error at index %d", i)
		}
	}
	return errors.WithStack(tx.Commit())
}

// Architecture bug in this function
//func (tx *Tx) WrapBuilder(bldrs ...QueryBuilder) (err error) {
//	defer func() {
//		if err != nil {
//			if rErr := tx.Rollback(); rErr != nil {
//				err = errors.Wrapf(rErr, "[dml] transaction.wrap.Rollback.error")
//			}
//		}
//	}()
//	for i := 0; i < len(bldrs) && err == nil; i++ {
//		bldr := bldrs[i]
//		switch b := bldr.(type) {
//		case *Arguments:
//			b.WithDB(tx.DB)
//		case *Delete:
//			b.WithDB(tx.DB)
//		case *Insert:
//			b.WithDB(tx.DB)
//		case *Update:
//			b.WithDB(tx.DB)
//		case *Select:
//			b.WithDB(tx.DB)
//		case *Union:
//			b.WithDB(tx.DB)
//		case *With:
//			b.WithDB(tx.DB)
//		//case *Stmt:
//		//	b.Stmt = tx.DB.Stmt(b.Stmt)
//		//	b.base.DB = stmtWrapper{stmt: b.Stmt}
//		default:
//			err = errors.NotSupported.Newf("[dml] WrapBuilder does not support this type: %T", bldr)
//		}
//
//	}
//	if err == nil {
//		err = errors.WithStack(tx.Commit())
//	}
//	return
//}
