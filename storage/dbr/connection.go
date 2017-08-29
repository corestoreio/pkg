// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr

import (
	"context"
	"database/sql"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/go-sql-driver/mysql"
)

type uniqueIDFn func() string

func uniqueIDNoop() string { return "" }

type logWithID struct {
	start time.Time
	Log   log.Logger
	// makeUniqueID generates for each call a new unique ID. Those IDs will be
	// assigned to a new connection or a new statement. The function signature is
	// equal to fmt.Stringer so one can use for example:
	//		uuid.NewV4().String
	makeUniqueID uniqueIDFn
}

// ConnPool at a connection to the database with an EventReceiver to send
// events, errors, and timings to
type ConnPool struct {
	logWithID
	DB *sql.DB
	// dn internal driver name
	dn string
	// dsn Data Source Name
	dsn *mysql.Config
}

// ConnPoolOption can be used at an argument in NewConnPool to configure a
// connection.
type ConnPoolOption func(*ConnPool) error

// WithLogger sets the customer logger to be used across the package. The logger
// gets inherited to type Conn and Tx and also to all statement types. Each
// heredity creates new fields as a prefix. Argument `uniqueID` generates for
// each heredity a new unique ID for tracing in Info logging. Those IDs will be
// assigned to a new connection or a new statement. The function signature is
// equal to fmt.Stringer so one can use for example:
//		uuid.NewV4().String
func WithLogger(l log.Logger, uniqueID func() string) ConnPoolOption {
	return func(c *ConnPool) error {
		c.makeUniqueID = uniqueID
		c.Log = l.With(log.String("id", c.makeUniqueID()))
		return nil
	}
}

// WithDB sets the DB value to a connection. If set ignores the DSN values.
// Mainly used for testing.
func WithDB(db *sql.DB) ConnPoolOption {
	return func(c *ConnPool) error {
		c.DB = db
		return nil
	}
}

// WithDSN sets the data source name for a connection.
func WithDSN(dsn string) ConnPoolOption {
	return func(c *ConnPool) error {
		myc, err := mysql.ParseDSN(dsn)
		if err != nil {
			return errors.WithStack(err)
		}
		c.dsn = myc
		return nil
	}
}

// NewConnPool instantiates a ConnPool for a given database/sql connection
// and event receiver. An invalid driver name causes a NotImplemented error to be
// returned. You can either apply a DSN or a pre configured *sql.DB type. For
// full UTF-8 support you must set the charset in the SQL driver to utf8mb4.
func NewConnPool(opts ...ConnPoolOption) (*ConnPool, error) {
	c := &ConnPool{
		dn: DriverNameMySQL,
	}
	if err := c.Options(opts...); err != nil {
		return nil, errors.WithStack(err)
	}

	switch c.dn {
	case DriverNameMySQL:
	default:
		return nil, errors.NewNotImplementedf("[dbr] unsupported driver: %q", c.dn)
	}
	if c.makeUniqueID == nil {
		c.makeUniqueID = uniqueIDNoop
	}
	if c.DB != nil || c.dsn == nil {
		return c, nil
	}

	// validate that DSN contains the utf8mb4 setting

	var err error
	if c.DB, err = sql.Open(c.dn, c.dsn.FormatDSN()); err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO: Validate that we run with utf8mb4 the normal utf8 is only 3 bytes
	// where utf8mb4 is full 4byte support.
	// SHOW VARIABLES WHERE Variable_name LIKE 'character\_set\_%' OR Variable_name LIKE 'collation%';

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
	for _, opt := range opts {
		if err := opt(c); err != nil {
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
	if c.Log != nil && c.Log.IsInfo() {
		defer c.Log.Info("ConnPool", log.String("type", "close"), log.Duration("duration", time.Since(c.start)))
	}
	return c.DB.Close() // no stack wrap otherwise error is hard to compare
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
		l = c.Log.With(log.String("ConnPool", "Conn"), log.String("id", c.makeUniqueID()))
	}
	return &Conn{
		logWithID: logWithID{
			start:        time.Now(),
			Log:          l,
			makeUniqueID: c.makeUniqueID,
		},
		Conn: dbc,
	}, errors.WithStack(err)
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
	*sql.Conn
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
	start := time.Now()

	dbTx, err := c.Conn.BeginTx(ctx, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	l := c.Log
	if l != nil {
		l = l.With(log.String("Conn", "Transaction"), log.String("id", c.makeUniqueID()))
		if l.IsInfo() {
			l.Info("Transaction", log.String("type", "begin"))
		}
	}
	return &Tx{
		logWithID: logWithID{
			start:        start,
			Log:          l,
			makeUniqueID: c.makeUniqueID,
		},
		Tx: dbTx,
	}, nil
}

// Close returns the connection to the connection pool. All operations after a
// Close will return with ErrConnDone. Close is safe to call concurrently with
// other operations and will block until all other operations finish. It may be
// useful to first cancel any used context and then call close directly after.
// It logs the time taken, if a logger has been set with Info logging enabled.
func (c *Conn) Close() error {
	if c.Log != nil && c.Log.IsInfo() {
		defer c.Log.Info("Conn", log.String("type", "close"), log.Duration("duration", time.Since(c.start)))
	}
	return c.Conn.Close() // no stack wrap otherwise error is hard to compare
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
	*sql.Tx
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
	start := time.Now()

	dbTx, err := c.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	l := c.Log
	if l != nil {
		l = l.With(log.String("ConnPool", "Transaction"), log.String("id", c.makeUniqueID()))
		if l.IsInfo() {
			l.Info("Transaction", log.String("type", "begin"))
		}
	}
	return &Tx{
		logWithID: logWithID{
			start:        start,
			Log:          l,
			makeUniqueID: c.makeUniqueID,
		},
		Tx: dbTx,
	}, nil
}

// Commit finishes the transaction. It logs the time taken, if a logger has been
// set with Info logging enabled.
func (tx *Tx) Commit() error {
	if tx.Log != nil && tx.Log.IsInfo() {
		defer tx.Log.Info("Transaction", log.String("type", "commit"), log.Duration("duration", time.Since(tx.start)))
	}
	return tx.Tx.Commit()
}

// Rollback cancels the transaction. It logs the time taken, if a logger has
// been set with Info logging enabled.
func (tx *Tx) Rollback() error {
	if tx.Log != nil && tx.Log.IsInfo() {
		defer tx.Log.Info("Transaction", log.String("type", "rollback"), log.Duration("duration", time.Since(tx.start)))
	}
	return tx.Tx.Rollback()
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
				return errors.Wrapf(err, "[dbr] transaction.wrap.Rollback.error at index %d", i)
			}
			return errors.Wrapf(err, "[dbr] transaction.wrap.error at index %d", i)
		}
	}
	return errors.WithStack(tx.Commit())
}
