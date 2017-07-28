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
	"database/sql"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/go-sql-driver/mysql"
)

// Connection at a connection to the database with an EventReceiver to send
// events, errors, and timings to
type Connection struct {
	DB  *sql.DB
	Log log.Logger
	// dn internal driver name
	dn string
	// dsn Data Source Name
	dsn *mysql.Config
	// DatabaseName contains the database name to which this connection has been
	// bound to. It will only be set when a DSN has been parsed.
	DatabaseName string
}

// ConnectionOption can be used at an argument in NewConnection to configure a
// connection.
type ConnectionOption func(*Connection) error

// WithDB sets the DB value to a connection. If set ignores the DSN values.
func WithDB(db *sql.DB) ConnectionOption {
	return func(c *Connection) error {
		c.DB = db
		return nil
	}
}

// WithDSN sets the data source name for a connection.
func WithDSN(dsn string) ConnectionOption {
	return func(c *Connection) error {
		myc, err := mysql.ParseDSN(dsn)
		if err != nil {
			return errors.WithStack(err)
		}
		c.dsn = myc
		return nil
	}
}

// NewConnection instantiates a Connection for a given database/sql connection
// and event receiver. An invalid drivername causes a NotImplemented error to be
// returned. You can either apply a DSN or a pre configured *sql.DB type.
func NewConnection(opts ...ConnectionOption) (*Connection, error) {
	c := &Connection{
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

	if c.dsn != nil {
		c.DatabaseName = c.dsn.DBName
	}

	if c.DB != nil || c.dsn == nil {
		return c, nil
	}

	var err error
	if c.DB, err = sql.Open(c.dn, c.dsn.FormatDSN()); err != nil {
		return nil, errors.WithStack(err)
	}

	return c, nil
}

// MustConnectAndVerify at like NewConnection but it verifies the connection
// and panics on errors.
func MustConnectAndVerify(opts ...ConnectionOption) *Connection {
	c, err := NewConnection(opts...)
	if err != nil {
		panic(err)
	}
	if err := c.DB.Ping(); err != nil {
		panic(err)
	}
	return c
}

// Options applies options to a connection
func (c *Connection) Options(opts ...ConnectionOption) error {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Close closes the database, releasing any open resources.
func (c *Connection) Close() error {
	return errors.WithStack(c.DB.Close())
}

// BuilderBase contains fields which all SQL query builder have in common, the
// same base. Exported for documentation reasons.
type BuilderBase struct {
	// ID of a statement. Used in logging. If empty the generated SQL string
	// gets used which can might contain sensitive information which should not
	// get logged. TODO implement
	ID           string
	Log          log.Logger // Log optional logger
	RawFullSQL   string
	RawArguments Arguments // Arguments used by RawFullSQL

	Table identifier

	// PropagationStopped set to true if you would like to interrupt the
	// listener chain. Once set to true all sub sequent calls of the next
	// listeners will be suppressed.
	PropagationStopped bool
	IsInterpolate      bool // See Interpolate()
	IsBuildCache       bool // see BuildCache()
	cacheSQL           []byte
	cacheArgs          Arguments // like a buffer, gets reused
	// propagationStoppedAt position in the slice where the stopped propagation
	// has been requested. for every new iteration the propagation must stop at
	// this position.
	propagationStoppedAt int
}

// BuilderConditional defines base fields used in statements which can have
// conditional constraints like WHERE, JOIN, ORDER, etc.
type BuilderConditional struct {
	// Record if set retrieves the necessary arguments from the interface.
	Record      ArgumentsAppender
	Joins       Joins
	Wheres      Conditions
	OrderBys    identifiers
	LimitCount  uint64
	OffsetCount uint64
	LimitValid  bool
	OffsetValid bool
}
