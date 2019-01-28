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
	"os"
	"sort"
	"strings"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/go-sql-driver/mysql"
)

type uniqueIDFn func() string

func uniqueIDNoOp() string             { return "" }
func mapTableNameNoOp(n string) string { return n }

const (
	eventOnOpen = iota // must start with zero
	eventOnClose
)

type connCommon struct {
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
	mapTableName func(oldName string) (newName string)
	runOnClose   []ConnPoolOption
}

// ConnPool at a connection to the database with an EventReceiver to send
// events, errors, and timings to
type ConnPool struct {
	connCommon
	// DB must be set using one of the ConnPoolOption function.
	DB  *sql.DB
	dsn *mysql.Config
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
	connCommon
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
// Practical Guide to SQL Transaction Isolation:
// https://begriffs.com/posts/2017-08-01-practical-guide-sql-isolation.html
type Tx struct {
	connCommon
	DB *sql.Tx
}

// ConnPoolOption can be used at an argument in NewConnPool to configure a
// connection.
type ConnPoolOption struct {
	eventType uint8
	sortOrder uint8
	fn        func(*ConnPool) error
	// WithUniqueIDFn applies a unique ID generator function without an applied
	// logger as in WithLogger. For more details see WithLogger function.
	// Sort Order 8.
	UniqueIDFn func() string
	// TableNameMapper maps the old name in the DML query to a new name. E.g.
	// for adding a prefix and/or a suffix.
	TableNameMapper func(oldName string) (newName string)
	// OptimisticLock is enabled all queries with Exec will have a `version` field.
	// UPDATE user SET ..., version = version + 1 WHERE id = ? AND version = ?
	// TODO implement OptimisticLock
	OptimisticLock bool
	// OptimisticLockFieldName custom global column name, defaults to `version
	// uint64`.
	// TODO implement OptimisticLock
	OptimisticLockColumnName string
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

// WithVerifyConnection checks if the connection to the server is valid and can
// be established.
func WithVerifyConnection() ConnPoolOption {
	return ConnPoolOption{
		sortOrder: 253,
		fn: func(c *ConnPool) error {
			return errors.WithStack(c.DB.Ping())
		},
	}
}

// WithCreateDatabase creates the database and sets the utf8mb4 option. It does
// not drop the database.
func WithCreateDatabase(ctx context.Context, databaseName string) ConnPoolOption {
	return ConnPoolOption{
		sortOrder: 253,
		fn: func(c *ConnPool) error {

			if databaseName == "" {
				databaseName = c.dsn.DBName
			}
			if err := WithSetNamesUTF8MB4().fn(c); err != nil {
				return errors.WithStack(err)
			}
			qdb := Quoter.Name(databaseName)
			if _, err := c.DB.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+qdb); err != nil {
				return errors.WithStack(err)
			}
			if _, err := c.DB.ExecContext(ctx, "ALTER DATABASE "+qdb+" DEFAULT CHARACTER SET='utf8mb4' COLLATE='utf8mb4_unicode_ci'"); err != nil {
				return errors.WithStack(err)
			}
			return nil
		},
	}
}

// TODO func WithRequireUTF8MB4() ConnPoolOption {
// 	return ConnPoolOption{
// 		sortOrder: 253,
// 		fn: func(c *ConnPool) error {
// 			// For Schemas:
// 			//
// 			// SELECT default_character_set_name FROM information_schema.SCHEMATA
// 			// WHERE schema_name = "schemaname";
// 			// For Tables:
// 			//
// 			// SELECT CCSA.character_set_name FROM information_schema.`TABLES` T,
// 			// 	information_schema.`COLLATION_CHARACTER_SET_APPLICABILITY` CCSA
// 			// WHERE CCSA.collation_name = T.table_collation
// 			// AND T.table_schema = "schemaname"
// 			// AND T.table_name = "tablename";
// 			// For Columns:
// 			//
// 			// SELECT character_set_name FROM information_schema.`COLUMNS`
// 			// WHERE table_schema = "schemaname"
// 			// AND table_name = "tablename"
// 			// AND column_name = "columnname";
//
// 			return nil
// 		},
// 	}
// }

// WithExecSQLOnConnOpen runs the sqlQuery arguments after successful opening a
// DB connection. All queries are running in a transaction.
func WithExecSQLOnConnOpen(ctx context.Context, sqlQuery ...string) ConnPoolOption {
	return withExecSQL(ctx, eventOnOpen, sqlQuery...)
}

// WithExecSQLOnConnClose runs the sqlQuery arguments before closing a DB
// connection. All queries are running in a transaction.
func WithExecSQLOnConnClose(ctx context.Context, sqlQuery ...string) ConnPoolOption {
	return withExecSQL(ctx, eventOnClose, sqlQuery...)
}

func withExecSQL(ctx context.Context, event uint8, sqlQuery ...string) ConnPoolOption {
	return ConnPoolOption{
		eventType: event,
		sortOrder: 250,
		fn: func(c *ConnPool) error {
			if len(sqlQuery) == 0 {
				return errors.Empty.Newf("[dml] WithInitialExecSQL argument sqlQuery is empty.")
			}

			fns := make([]func(*Tx) error, len(sqlQuery))
			for i, sq := range sqlQuery {
				sq := sq // prevent bug while scoping
				fns[i] = func(tx *Tx) error {
					if _, err := tx.DB.ExecContext(ctx, sq); err != nil {
						return errors.Wrapf(err, "[dml] WithInitialExecSQL Query: %q", sq)
					}
					return nil
				}
			}
			return c.Transaction(ctx, &sql.TxOptions{}, fns...)
		},
	}
}

// WithSetNamesUTF8MB4 sets the utf8mb4 charset and collation.
func WithSetNamesUTF8MB4() ConnPoolOption {
	return ConnPoolOption{
		sortOrder: 2, // must run after WithDSN and WithDB
		fn: func(c *ConnPool) error {
			_, err := c.DB.ExecContext(context.Background(), "SET NAMES 'utf8mb4' COLLATE 'utf8mb4_unicode_ci'")
			return errors.WithStack(err)
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

// WithDSN sets the data source name for a connection.
// Second argument DriverCallBack adds a low level call back function on MySQL driver level to
// create a a new instrumented driver. No need to call `sql.Register`!
func WithDSN(dsn string, cb ...DriverCallBack) ConnPoolOption {
	if len(cb) > 1 {
		panic(errors.NotImplemented.Newf("[dml] Only one DriverCallBack function does currently work. You provided: %d", len(cb)))
	}
	return ConnPoolOption{
		sortOrder: 0,
		fn: func(c *ConnPool) (err error) {
			if !strings.Contains(dsn, "parseTime") {
				return errors.NotImplemented.Newf("[dml] The DSN for go-sql-driver/mysql must contain the parameters `?parseTime=true[&loc=YourTimeZone]`")
			}
			if c.dsn, err = mysql.ParseDSN(dsn); err != nil {
				return errors.WithStack(err)
			}
			var drv driver.Driver = mysql.MySQLDriver{}
			if len(cb) == 1 {
				drv = wrapDriver(drv, cb[0])
			}
			c.DB = sql.OpenDB(dsnConnector{dsn: dsn, driver: drv})
			return nil
		},
	}
}

// EnvDSN is the name of the environment variable
const EnvDSN string = "CS_DSN"

// WithDSNfromEnv loads the DSN string from an environment variable named by
// `dsnEnvName`. If `dsnEnvName` is empty, then it falls back to the environment
// variable name of constant `EnvDSN`.
func WithDSNfromEnv(dsnEnvName string, cb ...DriverCallBack) ConnPoolOption {
	if dsnEnvName == "" {
		dsnEnvName = EnvDSN
	}
	return WithDSN(os.Getenv(dsnEnvName), cb...)
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
//
// Quote: http://techblog.en.klab-blogs.com/archives/31093990.html
// Recommended sql.DB Settings:
//
// Definitely set SetMaxOpenConns(). You need this in order to stop opening new
// connections and sending queries when the load is high and server response
// slows. If possible, it’s good to do a load test and set the minimum number of
// connections to ensure maximum throughput, but even if you can’t do that, you
// should decide on a reasonably appropriate number based on max_connection and
// the number of cores.
//
// Configure SetMaxIdleConns() to be equal to or higher than SetMaxOpenConns().
// Let SetConnMaxLifetime handle closing idle connections.
//
// Set SetConnMaxLifetime() to be the maximum number of connections x 1 second.
// In most environments, a load of one connection per second won’t be a problem.
// When you want to set it for longer than an hour, discuss that with an
// infrastructure/network engineer.
func NewConnPool(opts ...ConnPoolOption) (*ConnPool, error) {
	c := &ConnPool{}
	if err := c.Options(opts...); err != nil {
		return nil, errors.WithStack(err)
	}

	if c.makeUniqueID == nil {
		c.makeUniqueID = uniqueIDNoOp
	}
	if c.mapTableName == nil {
		c.mapTableName = mapTableNameNoOp
	}
	// validate that DSN contains the utf8mb4 setting

	return c, nil
}

// MustConnectAndVerify at like NewConnPool but it verifies the connection
// and panics on errors.
func MustConnectAndVerify(opts ...ConnPoolOption) *ConnPool {
	c, err := NewConnPool(opts...)
	if err != nil {
		panic(err)
	}
	if err := c.Options(WithVerifyConnection()); err != nil {
		panic(err)
	}
	return c
}

// Options applies options to a connection.
func (c *ConnPool) Options(opts ...ConnPoolOption) error {
	for i, opt := range opts {
		if opt.UniqueIDFn != nil {
			opts[i].sortOrder = 8 // must be this number
			opt := opt
			opts[i].fn = func(cp *ConnPool) error {
				cp.makeUniqueID = opt.UniqueIDFn
				return nil
			}
		}
		if opt.TableNameMapper != nil {
			opts[i].sortOrder = 20 // just a number
			opt := opt
			opts[i].fn = func(cp *ConnPool) error {
				cp.mapTableName = opt.TableNameMapper
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
		switch opt.eventType {
		case eventOnOpen:
			if err := opt.fn(c); err != nil {
				return errors.WithStack(err)
			}
		case eventOnClose:
			c.runOnClose = append(c.runOnClose, opt)
		}
	}

	return nil
}

// Close closes the database, releasing any open resources.
//
// It is rare to Close a DB, as the DB handle is meant to be long-lived and
// shared between many goroutines. It logs the time taken, if a logger has been
// set with Info logging enabled. It runs the ConnPoolOption, marked for running
// before close.
func (c *ConnPool) Close() error {
	if c.Log != nil && c.Log.IsDebug() {
		defer c.Log.Debug("Close", log.Duration("duration", now().Sub(c.start)))
	}
	for _, opt := range c.runOnClose {
		if err := opt.fn(c); err != nil {
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
		connCommon: connCommon{
			start:        start,
			Log:          l,
			makeUniqueID: c.makeUniqueID,
			mapTableName: c.mapTableName,
		},
		DB: dbTx,
	}, nil
}

// Transaction is a helper method that will automatically BEGIN a transaction
// and COMMIT or ROLLBACK once the supplied functions are done executing.
//
//      if err := con.Transaction(
// 			func(tx *dml.Tx) error {
//          	// SQL
// 		        return nil
//      	}[,
// 			func(tx *dml.Tx) error {
//          	// more SQL
// 		        return nil
//      	},]
// 		); err != nil{
//           panic(err.Error()) // you could gracefully handle the error also
//      }
// It logs the time taken, if a logger has been set with Debug logging enabled.
// The provided context gets used only for starting the transaction.
func (c *ConnPool) Transaction(ctx context.Context, opts *sql.TxOptions, fns ...func(*Tx) error) error {
	tx, err := c.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	for i, f := range fns {
		if err := f(tx); err != nil {
			err = errors.Wrapf(err, "[dml] ConnPool.Transaction.error at index %d", i)
			if rErr := tx.Rollback(); rErr != nil && err == nil {
				err = errors.Wrapf(rErr, "[dml] ConnPool.Transaction.Rollback.error at index %d", i)
			}
			return err
		}
	}
	return errors.WithStack(tx.Commit())
}

// WithQueryBuilder creates a new Artisan for handling the arguments with the
// assigned connection and builds the SQL string. The returned arguments and
// errors of the QueryBuilder will be forwarded to the Artisan type.
func (c *ConnPool) WithQueryBuilder(qb QueryBuilder) *Artisan {
	sql, argsRaw, err := qb.ToSQL()
	var args [defaultArgumentsCapacity]argument
	return &Artisan{
		base: builderCommon{
			cachedSQL: []string{"", sql},
			Log:       c.Log,
			id:        c.makeUniqueID(),
			DB:        c.DB,
			ärgErr:    errors.WithStack(err),
		},
		raw:       argsRaw,
		arguments: args[:0],
	}
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
		connCommon: connCommon{
			start:        now(),
			Log:          l,
			makeUniqueID: c.makeUniqueID,
			mapTableName: c.mapTableName,
		},
		DB: dbc,
	}, errors.WithStack(err)
}

// WithRawSQL creates a new Artisan for the given SQL string. It does not
// prepare the query nor runs place holder substitution.
func (c *ConnPool) WithRawSQL(query string) *Artisan {
	id := c.makeUniqueID()
	l := c.Log
	if l != nil {
		l = l.With(log.String("conn_pool_raw_sql_id", id), log.String("query", query))
	}
	var args [defaultArgumentsCapacity]argument
	return &Artisan{
		base: builderCommon{
			cachedSQL: []string{"", query},
			Log:       l,
			id:        id,
			DB:        c.DB,
		},
		arguments: args[:0],
	}
}

// Prepare executes the statement represented by the Select to create a prepared
// statement. It returns a custom statement type or an error if there was one.
// Provided arguments or records in the Select are getting ignored. The provided
// context is used for the preparation of the statement, not for the execution
// of the statement. The returned Stmter is not safe for concurrent use, despite
// the underlying *sql.Stmt is.
func (c *ConnPool) WithPrepare(ctx context.Context, query string) *Artisan {
	id := c.makeUniqueID()
	l := c.Log
	if l != nil {
		l = l.With(log.String("conn_pool_prepare_sql_id", id), log.String("query", query))
	}

	stmt, err := c.DB.PrepareContext(ctx, query)

	var args [defaultArgumentsCapacity]argument
	a := &Artisan{
		base: builderCommon{
			id:     id,
			ärgErr: err,
			Log:    l,
			DB:     stmtWrapper{stmt: stmt},
		},
		arguments:  args[:0],
		isPrepared: true,
	}
	return a
}

// WithDisabledForeignKeyChecks runs the callBack with disabled foreign key
// checks in a dedicated connection session. Foreign key checks are getting
// automatically re-enabled. The context is used to disable and enable the FK
// check.
func (c *ConnPool) WithDisabledForeignKeyChecks(ctx context.Context, callBack func(*Conn) error) (err error) {
	dbc, err := c.Conn(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if _, err2 := dbc.DB.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1"); err == nil && err2 != nil {
			err = errors.WithStack(err2)
		}
		if err2 := dbc.Close(); err == nil && err2 != nil {
			err = errors.WithStack(err2)
		}
	}()
	if _, err = dbc.DB.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0"); err != nil {
		err = errors.WithStack(err)
		return
	}
	err = callBack(dbc)
	return
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
		connCommon: connCommon{
			start:        start,
			Log:          l,
			makeUniqueID: c.makeUniqueID,
			mapTableName: c.mapTableName,
		},
		DB: dbTx,
	}, nil
}

// Transaction is a helper method that will automatically BEGIN a transaction
// and COMMIT or ROLLBACK once the supplied functions are done executing.
//
//      if err := con.Transaction(
// 			func(tx *dml.Tx) error {
//          	// SQL
// 		        return nil
//      	}[,
// 			func(tx *dml.Tx) error {
//          	// more SQL
// 		        return nil
//      	},]
// 		); err != nil{
//           panic(err.Error()) // you could gracefully handle the error also
//      }
// It logs the time taken, if a logger has been set with Debug logging enabled.
// The provided context gets used only for starting the transaction.
func (c *Conn) Transaction(ctx context.Context, opts *sql.TxOptions, fns ...func(*Tx) error) error {
	tx, err := c.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	for i, f := range fns {
		if err := f(tx); err != nil {
			err = errors.Wrapf(err, "[dml] ConnPool.Transaction.error at index %d", i)
			if rErr := tx.Rollback(); rErr != nil {
				err = errors.Wrapf(rErr, "[dml] ConnPool.Transaction.Rollback.error at index %d", i)
			}
			return err
		}
	}
	return errors.WithStack(tx.Commit())
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

// WithQueryBuilder creates a new Artisan for handling the arguments with the
// assigned connection and builds the SQL string. The returned arguments and
// errors of the QueryBuilder will be forwarded to the Artisan type.
func (c *Conn) WithQueryBuilder(qb QueryBuilder) *Artisan {
	sql, argsRaw, err := qb.ToSQL()
	id := c.makeUniqueID()
	l := c.Log
	if l != nil {
		l = l.With(log.String("query_builder_id", id), log.String("sql", sql))
	}
	var args [defaultArgumentsCapacity]argument
	return &Artisan{
		base: builderCommon{
			cachedSQL: []string{"", sql},
			Log:       l,
			id:        id,
			DB:        c.DB,
			ärgErr:    errors.WithStack(err),
		},
		raw:       argsRaw,
		arguments: args[:0],
	}
}

// WithRawSQL creates a new Artisan for the given SQL string in the current
// connection.
func (c *Conn) WithRawSQL(sql string) *Artisan {
	id := c.makeUniqueID()
	l := c.Log
	if l != nil {
		l = l.With(log.String("conn_pool_raw_sql_id", id), log.String("sql", sql))
	}
	var args [defaultArgumentsCapacity]argument
	return &Artisan{
		base: builderCommon{
			cachedSQL: []string{"", sql},
			Log:       l,
			id:        id,
			DB:        c.DB,
		},
		arguments: args[:0],
	}
}

// WithRawSQL creates a new Artisan for the given SQL string in the current
// transaction.
func (tx *Tx) WithRawSQL(sql string) *Artisan {
	id := tx.makeUniqueID()
	l := tx.Log
	if l != nil {
		l = l.With(log.String("tx_raw_sql_id", id), log.String("sql", sql))
	}
	var args [defaultArgumentsCapacity]argument
	return &Artisan{
		base: builderCommon{
			cachedSQL: []string{"", sql},
			Log:       l,
			id:        id,
			DB:        tx.DB,
		},
		arguments: args[:0],
	}
}

// Prepare executes the statement represented by the Select to create a prepared
// statement. It returns a custom statement type or an error if there was one.
// Provided arguments or records in the Select are getting ignored. The provided
// context is used for the preparation of the statement, not for the execution
// of the statement. The returned Stmter is not safe for concurrent use, despite
// the underlying *sql.Stmt is.
func (tx *Tx) WithPrepare(ctx context.Context, query string) *Artisan {
	id := tx.makeUniqueID()
	l := tx.Log
	if l != nil {
		l = l.With(log.String("tx_prepare_sql_id", id), log.String("query", query))
	}

	stmt, err := tx.DB.PrepareContext(ctx, query)

	var args [defaultArgumentsCapacity]argument
	a := &Artisan{
		base: builderCommon{
			id:     id,
			ärgErr: err,
			Log:    l,
			DB:     stmtWrapper{stmt: stmt},
		},
		arguments:  args[:0],
		isPrepared: true,
	}
	return a
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

// WithQueryBuilder creates a new Artisan for handling the arguments with the
// assigned connection and builds the SQL string. The returned arguments and
// errors of the QueryBuilder will be forwarded to the Artisan type.
func (tx *Tx) WithQueryBuilder(qb QueryBuilder) *Artisan {
	sql, argsRaw, err := qb.ToSQL()
	var args [defaultArgumentsCapacity]argument
	return &Artisan{
		base: builderCommon{
			cachedSQL: []string{"", sql},
			Log:       tx.Log,
			id:        tx.makeUniqueID(),
			DB:        tx.DB,
			ärgErr:    errors.WithStack(err),
		},
		raw:       argsRaw,
		arguments: args[:0],
	}
}
