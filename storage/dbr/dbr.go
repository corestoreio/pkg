package dbr

import (
	"database/sql"

	"github.com/corestoreio/csfw/util/errors"
	"github.com/go-sql-driver/mysql"
)

// DefaultDriverName is MySQL
const DefaultDriverName = DriverNameMySQL

// Connection is a connection to the database with an EventReceiver to send
// events, errors, and timings to
type Connection struct {
	DB *sql.DB
	EventReceiver
	// dn internal driver name
	dn string
	// dsn Data Source Name
	dsn *mysql.Config
	// DatabaseName contains the database name to which this connection has been
	// bound to. It will only be set when a DSN has been parsed.
	DatabaseName string
}

// Session represents a business unit of execution for some connection
type Session struct {
	cxn *Connection
	EventReceiver
}

// ConnectionOption can be used as an argument in NewConnection to configure a
// connection.
type ConnectionOption func(*Connection) error

// WithDB sets the DB value to a connection. If set ignores the DSN values.
func WithDB(db *sql.DB) ConnectionOption {
	return func(c *Connection) error {
		c.DB = db
		return nil
	}
}

// WithEventReceiver sets the event receiver for a connection.
func WithEventReceiver(log EventReceiver) ConnectionOption {
	return func(c *Connection) error {
		c.EventReceiver = log
		return nil
	}
}

// WithDSN sets the data source name for a connection.
func WithDSN(dsn string) ConnectionOption {
	return func(c *Connection) error {
		myc, err := mysql.ParseDSN(dsn)
		if err != nil {
			return errors.Wrap(err, "[dbr] mysql.ParseDSN")
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
		dn:            DriverNameMySQL,
		EventReceiver: nullReceiver,
	}
	if err := c.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[dbr] NewConnection.ApplyOpts")
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
		return nil, errors.Wrap(err, "[dbr] sql.Open")
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

// Options applies options to a connection
func (c *Connection) Options(opts ...ConnectionOption) error {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return errors.Wrap(err, "[dbr] Connection ApplyOpts")
		}
	}
	return nil
}

// NewSession instantiates a Session for the Connection
func (c *Connection) NewSession(opts ...SessionOption) *Session {
	s := &Session{
		cxn:           c,
		EventReceiver: c.EventReceiver, // Use parent instrumentation
	}
	s.Options(opts...)
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
type SessionOption func(cxn *Connection, s *Session) error

// SetSessionEventReceiver sets an event receiver securely to a session. Falls
// back to the parent event receiver if argument is nil.
func SetSessionEventReceiver(log EventReceiver) SessionOption {
	return func(cxn *Connection, s *Session) error {
		if log == nil {
			log = cxn.EventReceiver // Use parent instrumentation
		}
		s.EventReceiver = log
		return nil
	}
}

// Options applies options to a session
func (s *Session) Options(opts ...SessionOption) error {
	for _, opt := range opts {
		if err := opt(s.cxn, s); err != nil {
			return errors.Wrap(err, "[dbr] Session.Options")
		}
	}
	return nil
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
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	//QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	//QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	//ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
