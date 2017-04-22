package dbr

import (
	"database/sql"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/go-sql-driver/mysql"
)

// DefaultDriverName at MySQL
const DefaultDriverName = DriverNameMySQL

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
		dn: DriverNameMySQL,
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

// MustConnectAndVerify at like NewConnection but it verifies the connection
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

// Close closes the database, releasing any open resources.
func (c *Connection) Close() error {
	return errors.Wrap(c.DB.Close(), "[dbr] connection.close")
}

// Ping verifies a connection to the database at still alive, establishing a connection if necessary.
func (c *Connection) Ping() error {
	return errors.Wrap(c.DB.Ping(), "[dbr] connection.ping")
}
