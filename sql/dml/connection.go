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
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/corestoreio/pkg/util/bufferpool"

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

type queryCache struct {
	// makeUniqueID generates for each call a new unique ID. Those IDs will be
	// assigned to a new connection or a new statement. The function signature is
	// equal to fmt.Stringer so one can use for example:
	//		uuid.NewV4().String
	// The returned unique ID gets used in logging and inserted as a comment
	// into the SQL string. The returned string must not contain the
	// comment-end-termination pattern: `*/`.
	makeUniqueID uniqueIDFn
	mapTableName func(oldName string) (newName string)

	mu sync.RWMutex
	// cachedSQL contains the final SQL string which gets send to the server.
	// Using the CacheKey allows a dml type (insert,update,select ... ) to build
	// multiple different versions from object.
	queries map[string]*cachedSQL
}

type connCommon struct {
	start      time.Time
	Log        log.Logger
	runOnClose []ConnPoolOption
}

// ConnPool at a connection to the database with an EventReceiver to send
// events, errors, and timings to
type ConnPool struct {
	queryCache *queryCache // must be a pointer because we forward that pointer to Conn and Tx structs.
	connCommon
	driverCallBack DriverCallBack
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
	queryCache *queryCache
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
	queryCache *queryCache
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
			c.queryCache.makeUniqueID = uniqueIDFn
			if l != nil {
				c.Log = l.With(log.String("conn_pool_id", c.queryCache.makeUniqueID()))
			}
			return nil
		},
	}
}

// WithVerifyConnection checks if the connection to the server is valid and can
// be established.
func WithVerifyConnection(ctx context.Context, pingRetry time.Duration) ConnPoolOption {
	return ConnPoolOption{
		sortOrder: 149,
		fn: func(c *ConnPool) error {
			if err := c.DB.PingContext(ctx); err == nil {
				return nil
			}

			tkr := time.NewTicker(pingRetry)
			defer tkr.Stop()
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-tkr.C:
					if err := c.DB.PingContext(ctx); err == nil {
						return nil
					}
				}
			}
		},
	}
}

// WithCreateDatabase creates the database and sets the utf8mb4 option. It does
// not drop the database. If databaseName is empty, the DB name gets derived
// from the DSN.
func WithCreateDatabase(ctx context.Context, databaseName string) ConnPoolOption {
	return ConnPoolOption{
		sortOrder: 150,
		fn: func(c *ConnPool) error {
			if databaseName == "" && c.dsn != nil {
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

			if c.dsn != nil && c.dsn.User == "sqlmock" && c.dsn.Passwd == "sqlmock" {
				return nil
			}

			// explanation: close db and reconnect because MySQL driver binds to a
			// database. the SQL statement: `USE dbname` cannot be used because of new
			// connections where you have to call USE DB every time.
			if err := c.DB.Close(); err != nil {
				return errors.WithStack(err)
			}

			var drv driver.Driver = mysql.MySQLDriver{}
			if c.driverCallBack != nil {
				drv = wrapDriver(drv, c.driverCallBack, c.queryCache.makeUniqueID != nil)
			}
			if c.dsn != nil {
				dsn := c.dsn.FormatDSN()
				c.DB = sql.OpenDB(dsnConnector{
					dsn:    dsn,
					driver: drv,
				})
			} else {
				c.DB = nil
			}
			return nil
		},
	}
}

// WithExecSQLOnConnOpen runs the sqlQuery arguments after successful opening a
// DB connection. More than one queries are running in a transaction, a single
// query not.
func WithExecSQLOnConnOpen(ctx context.Context, sqlQuery ...string) ConnPoolOption {
	return withExecSQL(ctx, eventOnOpen, sqlQuery...)
}

// WithExecSQLOnConnClose runs the sqlQuery arguments before closing a DB
// connection. More than one queries are running in a transaction, a single
// query not.
func WithExecSQLOnConnClose(ctx context.Context, sqlQuery ...string) ConnPoolOption {
	return withExecSQL(ctx, eventOnClose, sqlQuery...)
}

func withExecSQL(ctx context.Context, event uint8, sqlQuery ...string) ConnPoolOption {
	return ConnPoolOption{
		eventType: event,
		sortOrder: 154,
		fn: func(c *ConnPool) error {
			switch len(sqlQuery) {

			case 0:
				return errors.Empty.Newf("[dml] WithInitialExecSQL argument sqlQuery is empty.")

			case 1:
				_, err := c.DB.ExecContext(ctx, sqlQuery[0])
				return err

			default:
				return c.Transaction(ctx, nil, func(tx *Tx) error {
					for i, sq := range sqlQuery {
						if _, err := tx.DB.ExecContext(ctx, sq); err != nil {
							return errors.Wrapf(err, "[dml] WithInitialExecSQL Query at index %d", i)
						}
					}
					return nil
				})
			}
		},
	}
}

// WithSetNamesUTF8MB4 sets the utf8mb4 charset and collation.
func WithSetNamesUTF8MB4() ConnPoolOption {
	return ConnPoolOption{
		sortOrder: 3, // must run after WithDSN and WithDB
		fn: func(c *ConnPool) error {
			_, err := c.DB.ExecContext(context.Background(), "SET NAMES 'utf8mb4' COLLATE 'utf8mb4_unicode_ci'")
			return errors.WithStack(err)
		},
	}
}

const randomTestDBPrefix = "test_"

// WithDB sets the DB value to an existing connection. Mainly used for testing.
// Does not support DriverCallBack.
func WithDB(db *sql.DB) ConnPoolOption {
	// this function can be called multiple times within different contexts.
	return ConnPoolOption{
		sortOrder: 2,
		fn: func(c *ConnPool) error {
			if c.DB == nil {
				c.DB = db
			}
			if c.DB == nil && c.dsn != nil {
				var drv driver.Driver = mysql.MySQLDriver{}
				if c.driverCallBack != nil {
					drv = wrapDriver(drv, c.driverCallBack, c.queryCache.makeUniqueID != nil)
				}

				dsn := c.dsn.FormatDSN()
				if strings.HasPrefix(c.dsn.DBName, randomTestDBPrefix) {
					// sql.OpenDB must connect with an empty DB name to not use the DB or
					// CREATE database statement will fail.
					db := c.dsn.DBName
					c.dsn.DBName = ""
					dsn = c.dsn.FormatDSN()
					c.dsn.DBName = db
				}

				c.DB = sql.OpenDB(dsnConnector{
					dsn:    dsn,
					driver: drv,
				})
			}
			return nil
		},
	}
}

// WithDSN sets the data source name for a connection. Second argument
// DriverCallBack adds a low level call back function on MySQL driver level to
// create a a new instrumented driver. No need to call `sql.Register`! If the
// DSN contains as database name the word "random", then the name will be
// "test_[unixtimestamp_nano]", especially useful in tests.
// The environment variable SKIP_CLEANUP=1 skips dropping the test database.
//		$ SKIP_CLEANUP=1 go test -v -run=TestX
func WithDSN(dsn string) ConnPoolOption {
	return ConnPoolOption{
		sortOrder: 1,
		fn: func(c *ConnPool) (err error) {
			if !strings.Contains(dsn, "parseTime") {
				return errors.NotImplemented.Newf("[dml] The DSN for go-sql-driver/mysql must contain the parameters `?parseTime=true[&loc=YourTimeZone]`")
			}
			if c.dsn, err = mysql.ParseDSN(dsn); err != nil {
				return errors.WithStack(err)
			}
			if c.dsn.DBName == "random" {
				db := fmt.Sprintf(randomTestDBPrefix+"%d", time.Now().UnixNano())
				c.dsn.DBName = db
				if os.Getenv("SKIP_CLEANUP") != "1" {
					opt := WithExecSQLOnConnClose(context.Background(), "DROP DATABASE IF EXISTS "+Quoter.Name(db))
					opt.sortOrder = 200 // must run at the very end. or other queries will fail.
					c.runOnClose = append(c.runOnClose, opt)
				}
			}

			return nil
		},
	}
}

// EnvDSN is the name of the environment variable
const EnvDSN string = "CS_DSN"

// WithDSNFromEnv loads the DSN string from an environment variable named by
// `dsnEnvName`. If `dsnEnvName` is empty, then it falls back to the environment
// variable name of constant `EnvDSN`.
func WithDSNFromEnv(dsnEnvName string) ConnPoolOption {
	if dsnEnvName == "" {
		dsnEnvName = EnvDSN
	}
	env := os.Getenv(dsnEnvName)
	if env == "" {
		return ConnPoolOption{
			sortOrder: 0,
			fn: func(c *ConnPool) (err error) {
				return errors.NotExists.Newf("[dml] The environment variable %q does not exists.", dsnEnvName)
			},
		}
	}
	return WithDSN(env)
}

// WithDriverCallBack allows low level query logging and argument inspection.
func WithDriverCallBack(cb DriverCallBack) ConnPoolOption {
	return ConnPoolOption{
		sortOrder: 0,
		fn: func(c *ConnPool) (err error) {
			c.driverCallBack = cb
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
	c := ConnPool{
		queryCache: &queryCache{
			queries: make(map[string]*cachedSQL),
		},
	}
	opts = append(opts, WithDB(nil))
	if err := c.options(opts...); err != nil {
		return nil, errors.WithStack(err)
	}

	if c.queryCache.makeUniqueID == nil {
		c.queryCache.makeUniqueID = uniqueIDNoOp
	}
	if c.queryCache.mapTableName == nil {
		c.queryCache.mapTableName = mapTableNameNoOp
	}
	// validate that DSN contains the utf8mb4 setting, if DSN is set

	return &c, nil
}

// MustConnectAndVerify at like NewConnPool but it verifies the connection
// and panics on errors.
func MustConnectAndVerify(opts ...ConnPoolOption) *ConnPool {
	c, err := NewConnPool(opts...)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second) // not an elegant solution
	defer cancel()
	if err := c.options(WithVerifyConnection(ctx, 10*time.Second)); err != nil {
		panic(err)
	}
	return c
}

// Options applies options to a connection.
func (c *ConnPool) options(opts ...ConnPoolOption) error {
	for i, opt := range opts {
		if opt.UniqueIDFn != nil {
			opts[i].sortOrder = 8 // must be this number
			opt := opt
			opts[i].fn = func(cp *ConnPool) error {
				cp.queryCache.makeUniqueID = opt.UniqueIDFn
				return nil
			}
		}
		if opt.TableNameMapper != nil {
			opts[i].sortOrder = 20 // just a number
			opt := opt
			opts[i].fn = func(cp *ConnPool) error {
				cp.queryCache.mapTableName = opt.TableNameMapper
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

const (
	sqlIDPrefix    = "/*$ID$"
	sqlIDPrefixLen = 6
	sqlIDSuffix    = "*/"
	sqlIDSuffixLen = 2
)

func extractSQLIDPrefix(rawSQL string) (prefix string, lastPos int) {
	if !strings.HasPrefix(rawSQL, sqlIDPrefix) {
		return "", 0
	}
	// remove unique query ID /*$ID$....*/
	lastPos = strings.Index(rawSQL, sqlIDSuffix)
	if lastPos < 0 {
		lastPos = 0
	} else {
		lastPos += sqlIDSuffixLen
	}
	if lastPos < sqlIDSuffixLen {
		return "", 0
	}
	prefix = rawSQL[sqlIDPrefixLen : lastPos-sqlIDSuffixLen]

	return prefix, lastPos
}

func (qc *queryCache) prependUniqueID(rawSQL string) (id, _rawSQL string) {
	if qc.makeUniqueID != nil {
		id = qc.makeUniqueID()
		if id != "" {
			var buf strings.Builder
			buf.WriteString(sqlIDPrefix)
			buf.WriteString(id)
			buf.WriteString(sqlIDSuffix)
			buf.WriteString(rawSQL)
			rawSQL = buf.String()
		}
	}
	return id, rawSQL
}

func (qc *queryCache) cacheKeyExists(cacheKey string) bool {
	qc.mu.RLock()
	_, ok := qc.queries[cacheKey]
	qc.mu.RUnlock()
	return ok
}

func (qc *queryCache) initDBRCacheKey(
	ctx context.Context,
	l log.Logger,
	connSource, cacheKey string,
	isPrepared bool,
	db QueryExecPreparer,
	opts []DBRFunc,
) *DBR {
	qc.mu.RLock()
	defer qc.mu.RUnlock()

	if l != nil && l.IsDebug() {
		defer log.WhenDone(l).Debug("Prepare")
	}

	dbr := &DBR{
		customCacheKey: cacheKey,
		DB:             db,
		ResultCheckFn:  strictAffectedRowsResultCheck,
		isPrepared:     isPrepared,
	}
	for _, opt := range opts {
		opt(dbr)
	}
	sqlCache, ok := qc.queries[dbr.customCacheKey]
	if !ok {
		return &DBR{
			previousErr: errors.NotFound.Newf("CacheKey %q not found", dbr.customCacheKey),
		}
	}
	if l != nil {
		l = l.With(log.String("conn_source", connSource), log.Bool("is_prepared", isPrepared),
			log.String("query_id", sqlCache.id), log.String("query", sqlCache.rawSQL))
	}
	dbr.cachedSQL = *sqlCache
	dbr.log = l

	if isPrepared {
		stmt, err := db.PrepareContext(ctx, dbr.cachedSQL.rawSQL)
		if err != nil {
			return &DBR{
				previousErr: err,
			}
		}
		dbr.DB = stmtWrapper{stmt: stmt}
	}

	return dbr
}

func (qc *queryCache) initDBRQB(
	ctx context.Context,
	l log.Logger,
	connSource string,
	isPrepared bool,
	qb QueryBuilder,
	db QueryExecPreparer,
	opts []DBRFunc,
) *DBR {
	prepareQueryBuilder(qc.mapTableName, qb)
	rawSQL, _, err := qb.ToSQL()
	if err != nil {
		return &DBR{
			previousErr: errors.WithStack(err),
		}
	}

	if isPrepared {
		stmt, err := db.PrepareContext(ctx, rawSQL)
		if err != nil {
			return &DBR{
				previousErr: errors.WithStack(err),
			}
		}
		db = stmtWrapper{stmt: stmt}
	}

	dbr := &DBR{
		customCacheKey: hashSQL(rawSQL),
		DB:             db,
		ResultCheckFn:  strictAffectedRowsResultCheck,
		isPrepared:     isPrepared,
	}

	for _, opt := range opts {
		opt(dbr)
	}

	qc.mu.Lock()
	defer qc.mu.Unlock()
	sqlCache, ok := qc.queries[dbr.customCacheKey]
	if !ok {
		id, rawSQL := qc.prependUniqueID(rawSQL)
		sqlCache = makeCachedSQL(qb, rawSQL, id)
		qc.queries[dbr.customCacheKey] = sqlCache
	}
	if l != nil {
		l = l.With(log.String("conn_source", connSource), log.Bool("is_prepared", isPrepared),
			log.String("query_id", sqlCache.id), log.String("query", sqlCache.rawSQL))
	}
	dbr.cachedSQL = *sqlCache
	// https://github.com/go101/go101/wiki
	dbr.cachedSQL.qualifiedColumns = append(sqlCache.qualifiedColumns[:0:0], sqlCache.qualifiedColumns...)
	dbr.log = l

	return dbr
}

// hashSQL removes spaces and transforms the SQL to all lowercase and hashes it,
// to avoid minor duplicates ... maybe that is not worth. refactor later.
func hashSQL(rawSQL string) string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if strings.HasPrefix(rawSQL, sqlIDPrefix) { // remove unique query ID /*$ID$....*/
		idx := strings.Index(rawSQL, sqlIDSuffix)
		if idx < 0 {
			idx = 0
		} else {
			idx += len(sqlIDSuffix)
		}
		rawSQL = rawSQL[idx:]
	}
	for _, r := range rawSQL {
		switch {
		case unicode.IsSpace(r):
		// do nothing, removes all spaces
		default:
			buf.WriteRune(unicode.ToUpper(r))
		}
	}
	h := fnv.New64a()
	buf.WriteTo(h)
	idx := strings.IndexFunc(rawSQL, unicode.IsSpace)
	if idx < 0 {
		idx = 0
	}
	return rawSQL[:idx] + hex.EncodeToString(h.Sum(nil)) // fastest code
}

// Schema returns the database name as provided in the DSN. Returns an empty
// string if no DSN has been set.
func (c *ConnPool) Schema() string {
	if c.dsn != nil {
		return c.dsn.DBName
	}
	return ""
}

// DSN returns the formatted DSN. Will leak the password.
func (c *ConnPool) DSN() string {
	if c.dsn != nil {
		return c.dsn.FormatDSN()
	}
	return ""
}

// Close closes the database, releasing any open resources.
//
// It is rare to Close a DB, as the DB handle is meant to be long-lived and
// shared between many goroutines. It logs the time taken, if a logger has been
// set with Info logging enabled. It runs the ConnPoolOption, marked for running
// before close.
func (c *ConnPool) Close() (err error) {
	if c.Log != nil && c.Log.IsDebug() {
		defer c.Log.Debug("Close", log.Err(err), log.Duration("duration", now().Sub(c.start)))
	}

	sort.SliceStable(c.runOnClose, func(i, j int) bool {
		return c.runOnClose[i].sortOrder < c.runOnClose[j].sortOrder
	})

	for _, opt := range c.runOnClose {
		if err = opt.fn(c); err != nil {
			return errors.WithStack(err)
		}
	}
	if c.DB != nil {
		err = c.DB.Close() // no stack wrap otherwise error is hard to compare
	}
	return
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
		l = l.With(log.String("tx_id", c.queryCache.makeUniqueID()))
		if l.IsDebug() {
			l.Debug("BeginTx")
		}
	}
	return &Tx{
		queryCache: c.queryCache,
		connCommon: connCommon{
			start: start,
			Log:   l,
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
//      	}
// 		); err != nil{
//           panic(err.Error()) // you could gracefully handle the error also
//      }
// It logs the time taken, if a logger has been set with Debug logging enabled.
// The provided context gets used only for starting the transaction.
func (c *ConnPool) Transaction(ctx context.Context, opts *sql.TxOptions, fn func(*Tx) error) (err error) {
	tx, err := c.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = errors.Wrapf(err, "Rollback failed too: %+v", rbErr)
			}
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (c *ConnPool) CachedQueries() map[string]string {
	c.queryCache.mu.RLock()
	defer c.queryCache.mu.RUnlock()

	queries := make(map[string]string, len(c.queryCache.queries))
	for key, cq := range c.queryCache.queries {
		queries[key] = cq.rawSQL
	}

	return queries
}

// RegisterByQueryBuilder adds the SQL queries to the local internal cache. The
// cacheKeyQB map gets iterated in alphabetical order.
func (c *ConnPool) RegisterByQueryBuilder(cacheKeyQB map[string]QueryBuilder) error {
	c.queryCache.mu.Lock()
	defer c.queryCache.mu.Unlock()

	keys := make([]string, 0, len(cacheKeyQB))
	for cacheKey := range cacheKeyQB {
		keys = append(keys, cacheKey)
	}
	sort.Strings(keys)
	for _, cacheKey := range keys {
		qb := cacheKeyQB[cacheKey]
		if _, ok := c.queryCache.queries[cacheKey]; ok {
			return errors.AlreadyExists.Newf("[dml] CacheKey %q already exists", cacheKey)
		}

		prepareQueryBuilder(c.queryCache.mapTableName, qb)
		rawSQL, _, err := qb.ToSQL()
		if err != nil {
			return errors.Fatal.New(err, "Failed to build SQL for cache key %q", cacheKey)
		}
		id, rawSQL := c.queryCache.prependUniqueID(rawSQL)
		c.queryCache.queries[cacheKey] = makeCachedSQL(qb, rawSQL, id)
	}
	return nil
}

func (c *ConnPool) DeregisterByCacheKey(cacheKey string) error {
	c.queryCache.mu.Lock()
	delete(c.queryCache.queries, cacheKey)
	c.queryCache.mu.Unlock()
	return nil
}

// WithCacheKey creates a DBR object from a cached query.
func (c *ConnPool) WithCacheKey(cacheKey string, opts ...DBRFunc) *DBR {
	return c.queryCache.initDBRCacheKey(context.Background(), c.Log, "ConnPool", cacheKey, false, c.DB, opts)
}

// CacheKeyExists returns true if a given key already exists.
func (c *ConnPool) CacheKeyExists(cacheKey string) bool {
	return c.queryCache.cacheKeyExists(cacheKey)
}

// WithPrepareCacheKey creates a DBR object from a prepared cached query.
func (c *ConnPool) WithPrepareCacheKey(ctx context.Context, cacheKey string, opts ...DBRFunc) *DBR {
	return c.queryCache.initDBRCacheKey(ctx, c.Log, "ConnPool", cacheKey, true, c.DB, opts)
}

// WithQueryBuilder creates a new DBR for handling the arguments with the
// assigned connection and builds the SQL string. The returned arguments and
// errors of the QueryBuilder will be forwarded to the DBR type. It generates a
// unique cache key based on the SQL string. The cache key can be retrieved via
// DBR object.
func (c *ConnPool) WithQueryBuilder(qb QueryBuilder, opts ...DBRFunc) *DBR {
	return c.queryCache.initDBRQB(context.Background(), c.Log, "ConnPool", false, qb, c.DB, opts)
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
		l = c.Log.With(log.String("conn_id", c.queryCache.makeUniqueID()))
	}
	return &Conn{
		queryCache: c.queryCache,
		connCommon: connCommon{
			start: now(),
			Log:   l,
		},
		DB: dbc,
	}, errors.WithStack(err)
}

// Prepare executes the statement represented by the Select to create a prepared
// statement. It returns a custom statement type or an error if there was one.
// Provided arguments or records in the Select are getting ignored. The provided
// context is used for the preparation of the statement, not for the execution
// of the statement. The returned DBR is not safe for concurrent use, despite
// the underlying *sql.Stmt is.
// It generates a unique cache key based on the SQL string. The cache key can be
// retrieved via DBR object.
func (c *ConnPool) WithPrepare(ctx context.Context, qb QueryBuilder, opts ...DBRFunc) *DBR {
	return c.queryCache.initDBRQB(ctx, c.Log, "ConnPool", true, qb, c.DB, opts)
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
		l = l.With(log.String("tx_id", c.queryCache.makeUniqueID()))
		if l.IsDebug() {
			l.Debug("BeginTx")
		}
	}
	return &Tx{
		queryCache: c.queryCache,
		connCommon: connCommon{
			start: start,
			Log:   l,
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
//      	},
// 		); err != nil{
//           panic(err.Error()) // you could gracefully handle the error also
//      }
// It logs the time taken, if a logger has been set with Debug logging enabled.
// The provided context gets used only for starting the transaction.
func (c *Conn) Transaction(ctx context.Context, opts *sql.TxOptions, f func(*Tx) error) error {
	tx, err := c.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	if err := f(tx); err != nil {
		err = errors.Wrapf(err, "[dml] ConnPool.Transaction.error")
		if rErr := tx.Rollback(); rErr != nil {
			err = errors.Wrapf(rErr, "[dml] ConnPool.Transaction.Rollback.error")
		}
		return err
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

// WithQueryBuilder creates a new DBR for handling the arguments with the
// assigned connection and builds the SQL string. The returned arguments and
// errors of the QueryBuilder will be forwarded to the DBR type. It generates a
// unique cache key based on the SQL string. The cache key can be retrieved via
// DBR object.
func (c *Conn) WithQueryBuilder(qb QueryBuilder, opts ...DBRFunc) *DBR {
	return c.queryCache.initDBRQB(context.Background(), c.Log, "Conn", false, qb, c.DB, opts)
}

// WithPrepare adds the query to the cache and returns a prepared statement
// which must be closed after its use.
func (c *Conn) WithPrepare(ctx context.Context, qb QueryBuilder, opts ...DBRFunc) *DBR {
	return c.queryCache.initDBRQB(ctx, c.Log, "Conn", true, qb, c.DB, opts)
}

// WithCacheKey creates a DBR object from a cached query.
func (c *Conn) WithCacheKey(cacheKey string, opts ...DBRFunc) *DBR {
	return c.queryCache.initDBRCacheKey(context.Background(), c.Log, "Conn", cacheKey, false, c.DB, opts)
}

// CacheKeyExists returns true if a given key already exists.
func (c *Conn) CacheKeyExists(cacheKey string) bool {
	return c.queryCache.cacheKeyExists(cacheKey)
}

// WithPrepareCacheKey creates a DBR object from a prepared cached query. The
// statement must be closed after its use.
func (c *Conn) WithPrepareCacheKey(ctx context.Context, cacheKey string, opts ...DBRFunc) *DBR {
	return c.queryCache.initDBRCacheKey(ctx, c.Log, "Conn", cacheKey, true, c.DB, opts)
}

// WithCacheKey creates a DBR object from a cached query.
func (tx *Tx) WithCacheKey(cacheKey string, opts ...DBRFunc) *DBR {
	return tx.queryCache.initDBRCacheKey(context.Background(), tx.Log, "Tx", cacheKey, false, tx.DB, opts)
}

// CacheKeyExists returns true if a given key already exists.
func (tx *Tx) CacheKeyExists(cacheKey string) bool {
	return tx.queryCache.cacheKeyExists(cacheKey)
}

// WithPrepareCacheKey creates a DBR object from a prepared cached query. After
// use the query statement must be closed.
func (tx *Tx) WithPrepareCacheKey(ctx context.Context, cacheKey string, opts ...DBRFunc) *DBR {
	return tx.queryCache.initDBRCacheKey(ctx, tx.Log, "Tx", cacheKey, true, tx.DB, opts)
}

// WithPrepare executes the statement represented by the Select to create a
// prepared statement. It returns a custom statement type or an error if there
// was one. Provided arguments or records in the Select are getting ignored. The
// provided context is used for the preparation of the statement, not for the
// execution of the statement. The returned Stmter is not safe for concurrent
// use, despite the underlying *sql.Stmt is. You must close DBR after its use.
// It generates a unique cache key based on the SQL string. The cache key can be
// retrieved via DBR object.
func (tx *Tx) WithPrepare(ctx context.Context, qb QueryBuilder, opts ...DBRFunc) *DBR {
	return tx.queryCache.initDBRQB(ctx, tx.Log, "Tx", true, qb, tx.DB, opts)
}

// WithQueryBuilder creates a new DBR for handling the arguments with the
// assigned connection and builds the SQL string. The returned arguments and
// errors of the QueryBuilder will be forwarded to the DBR type. It generates a
// unique cache key based on the SQL string. The cache key can be retrieved via
// DBR object.
func (tx *Tx) WithQueryBuilder(qb QueryBuilder, opts ...DBRFunc) *DBR {
	return tx.queryCache.initDBRQB(context.Background(), tx.Log, "Tx", false, qb, tx.DB, opts)
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

// TODO func WithRequireUTF8MB4() ConnPoolOption {
// 	return ConnPoolOption{
// 		sortOrder: 152,
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
