package dbr

import (
	"database/sql"

	"github.com/juju/errgo"
)

// DefaultDriverName is MySQL
const DefaultDriverName = "mysql"

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

// ConnOpts type for setting options on a connection.
type ConnOpts func(c *Connection)

// ConnDB sets the DB value to a connection. If set ignores the DSN values.
func ConnDB(db *sql.DB) ConnOpts {
	if db == nil {
		panic("DB argument cannot be nil")
	}
	return func(c *Connection) {
		c.DB = db
	}
}

// ConnEvent sets the event receiver for a connection.
func ConnEvent(log EventReceiver) ConnOpts {
	if log == nil {
		log = nullReceiver
	}
	return func(c *Connection) {
		c.EventReceiver = log
	}
}

// ConnDriver sets the driver name for a connection. At the moment only MySQL
// is supported.
func ConnDriver(driverName string) ConnOpts {
	return func(c *Connection) {
		c.dn = driverName
	}
}

// ConnDSN sets the data source name for a connection.
func ConnDSN(dsn string) ConnOpts {
	if dsn == "" {
		panic("DSN argument cannot be empty")
	}
	return func(c *Connection) {
		c.dsn = dsn
	}
}

// NewConnection instantiates a Connection for a given database/sql connection
// and event receiver
func NewConnection(opts ...ConnOpts) (*Connection, error) {
	c := &Connection{
		dn:            DefaultDriverName,
		EventReceiver: nullReceiver,
	}
	c.ApplyOpts(opts...)

	switch c.dn {
	case "mysql":
	default:
		return nil, errgo.Newf("unsupported driver: %s", c.dn)
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

// SessionOpts function type to apply options to a session
type SessionOpts func(cxn *Connection, s *Session)

// SessionEvent sets an event receiver securely to a session. Falls back to the
// parent event receiver if argument is nil.
func SessionEvent(log EventReceiver) SessionOpts {
	return func(cxn *Connection, s *Session) {
		if log == nil {
			log = cxn.EventReceiver // Use parent instrumentation
		}
		s.EventReceiver = log
	}
}

// ApplyOpts applies options to a connection
func (c *Connection) ApplyOpts(opts ...ConnOpts) *Connection {
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
	return c
}

// NewSession instantiates a Session for the Connection
func (c *Connection) NewSession(opts ...SessionOpts) *Session {
	s := &Session{
		cxn:           c,
		EventReceiver: c.EventReceiver, // Use parent instrumentation
	}
	for _, opt := range opts {
		if opt != nil {
			opt(c, s)
		}
	}
	return s
}

// MustConnectAndVerify is like NewConnection but it verifies the connection
// and panics on errors.
func MustConnectAndVerify(opts ...ConnOpts) *Connection {
	c, err := NewConnection(opts...)
	if err != nil {
		panic(err)
	}
	if err := c.Ping(); err != nil {
		panic(err)
	}
	return c
}

// Close closes the database, releasing any open resources.
func (c *Connection) Close() error {
	return c.EventErr("dbr.connection.close", c.DB.Close())
}

// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
func (c *Connection) Ping() error {
	return c.EventErr("dbr.connection.ping", c.DB.Ping())
}

// SessionRunner can do anything that a Session can except start a transaction.
type SessionRunner interface {
	Select(cols ...string) *SelectBuilder
	SelectBySql(sql string, args ...interface{}) *SelectBuilder

	InsertInto(into string) *InsertBuilder
	Update(table string) *UpdateBuilder
	UpdateBySql(sql string, args ...interface{}) *UpdateBuilder
	DeleteFrom(from string) *DeleteBuilder
}

type runner interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}
