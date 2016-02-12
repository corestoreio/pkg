package dbr

import (
	"database/sql"
	"github.com/juju/errors"
)

// DefaultDriverName is MySQL
const DefaultDriverName = DriverNameMySQL

// Connection is a connection to the database with an EventReceiver
// to send events, errors, and timings to
type Connection struct {
	DB *sql.DB
	EventReceiver
	// dn internal driver name
	dn string
	// dsn Data Source Name
	dsn string
}

// Session represents a business unit of execution for some connection
type Session struct {
	cxn *Connection
	EventReceiver
}

// ConnectionOption can be used as an argument in NewConnection to configure a connection.
type ConnectionOption func(c *Connection)

// SetDB sets the DB value to a connection. If set ignores the DSN values.
func SetDB(db *sql.DB) ConnectionOption {
	if db == nil {
		panic("DB argument cannot be nil")
	}
	return func(c *Connection) {
		c.DB = db
	}
}

// SetEventReceiver sets the event receiver for a connection.
func SetEventReceiver(log EventReceiver) ConnectionOption {
	if log == nil {
		log = nullReceiver
	}
	return func(c *Connection) {
		c.EventReceiver = log
	}
}

// SetDriver sets the driver name for a connection. At the moment only MySQL
// is supported.
func SetDriver(driverName string) ConnectionOption {
	return func(c *Connection) {
		c.dn = driverName
	}
}

// SetDSN sets the data source name for a connection.
func SetDSN(dsn string) ConnectionOption {
	if dsn == "" {
		panic("DSN argument cannot be empty")
	}
	return func(c *Connection) {
		c.dsn = dsn
	}
}

// NewConnection instantiates a Connection for a given database/sql connection
// and event receiver
func NewConnection(opts ...ConnectionOption) (*Connection, error) {
	c := &Connection{
		dn:            DriverNameMySQL,
		EventReceiver: nullReceiver,
	}
	c.ApplyOpts(opts...)

	switch c.dn {
	case DriverNameMySQL:
	default:
		return nil, errors.NotImplementedf("unsupported driver: %s", c.dn)
	}

	if c.DB != nil {
		return c, nil
	}

	if c.dsn != "" {
		var err error
		if c.DB, err = sql.Open(c.dn, c.dsn); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// MustConnectAndVerify is like NewConnection but it verifies the connection
// and panics on errors.
func MustConnectAndVerify(opts ...ConnectionOption) *Connection {
	c, err := NewConnection(opts...)
	if err != nil {
		panic(err)
	}
	if err := c.Ping(); err != nil {
		panic(err)
	}
	return c
}

// ApplyOpts applies options to a connection
func (c *Connection) ApplyOpts(opts ...ConnectionOption) *Connection {
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
	return c
}

// NewSession instantiates a Session for the Connection
func (c *Connection) NewSession(opts ...SessionOption) *Session {
	s := &Session{
		cxn:           c,
		EventReceiver: c.EventReceiver, // Use parent instrumentation
	}
	s.ApplyOpts(opts...)
	return s
}

// Close closes the database, releasing any open resources.
func (c *Connection) Close() error {
	return c.EventErr("dbr.connection.close", c.DB.Close())
}

// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
func (c *Connection) Ping() error {
	return c.EventErr("dbr.connection.ping", c.DB.Ping())
}

// SessionOption can be used as an argument in NewSession to configure a session.
type SessionOption func(cxn *Connection, s *Session) SessionOption

// SetSessionEventReceiver sets an event receiver securely to a session. Falls
// back to the parent event receiver if argument is nil.
// This function adheres http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func SetSessionEventReceiver(log EventReceiver) SessionOption {
	return func(cxn *Connection, s *Session) SessionOption {
		previous := s.EventReceiver
		if log == nil {
			log = cxn.EventReceiver // Use parent instrumentation
		}
		s.EventReceiver = log
		return SetSessionEventReceiver(previous)
	}
}

// NewSession instantiates a Session for the Connection
func (s *Session) ApplyOpts(opts ...SessionOption) (previous SessionOption) {
	for _, opt := range opts {
		if opt != nil {
			previous = opt(s.cxn, s)
		}
	}
	return previous
}

// SessionRunner can do anything that a Session can except start a transaction.
type SessionRunner interface {
	Select(cols ...string) *SelectBuilder
	SelectBySql(sql string, args ...interface{}) *SelectBuilder

	InsertInto(into string) *InsertBuilder
	Update(table ...string) *UpdateBuilder
	UpdateBySql(sql string, args ...interface{}) *UpdateBuilder
	DeleteFrom(from ...string) *DeleteBuilder
}

type runner interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}
