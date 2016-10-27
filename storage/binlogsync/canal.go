package binlogsync

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/siddontang/go-mysql/replication"
)

// Use flavor for different MySQL versions,
const (
	MySQLFlavor   = "mysql"
	MariaDBFlavor = "mariadb"
)

type DBer interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Close() error
}

// Canal can sync your MySQL data. MySQL must use the binlog format ROW.
type Canal struct {
	// mclose acts only during the call to Close().
	mclose sync.Mutex
	// dsn contains the parsed DSN
	dsn         *mysql.Config
	canalParams map[string]string

	cfgw config.Writer

	masterMu           sync.RWMutex
	masterStatus       csdb.MasterStatus
	masterLastSaveTime time.Time

	syncer *replication.BinlogSyncer

	rsMu       sync.RWMutex
	rsHandlers []RowsEventHandler

	db DBer

	// Tables contains the overall SQL table cache. If a table gets modified
	// during runtime of this program then somehow we must clear the cache to
	// reload the table structures.
	Tables *csdb.Tables

	closed *int32
	Log    log.Logger
	wg     sync.WaitGroup
}

// Option applies multiple options to the Canal type.
type Option func(*Canal) error

// WithMySQL adds the database/sql.DB driver including a ping to the database.
func WithMySQL() Option {
	return func(c *Canal) error {
		db, err := sql.Open("mysql", c.dsn.FormatDSN())
		if err != nil {
			return errors.Wrap(err, "[binlogsync] sql.Open")
		}
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, "[binlogsync] sql ping failed")
		}
		c.db = db
		return nil
	}
}

// WithDB allows to set your own DB connection. Mostly used for testing.
func WithDB(db *sql.DB) Option {
	return func(c *Canal) (err error) {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, "[binlogsync] sql ping failed")
		}
		c.db = db
		return nil
	}
}

func WithConfigurationWriter(w config.Writer) Option {
	//Write(p cfgpath.Path, value interface{}) error
	return func(c *Canal) error {
		c.cfgw = w
		return nil
	}
}

func withUpdateBinlogStartPosition(c *Canal) error {
	var ms csdb.MasterStatus
	if err := ms.Load(context.TODO(), c.db); err != nil {
		return errors.Wrap(err, "[binlogsync] ShowMasterStatus Load")
	}

	c.masterStatus = ms

	if v, ok := c.canalParams["BinlogStartPosition"]; ok && v != "" {
		hasPos := conv.ToUint(c.canalParams["BinlogStartPosition"])
		if hasPos >= 4 {
			c.masterStatus.Position = uint(hasPos)
		}
	}
	return nil
}

// withPrepareSyncer creates its own database connection.
func withPrepareSyncer(c *Canal) error {
	host, port, err := net.SplitHostPort(c.dsn.Addr)
	if err != nil {
		return errors.Wrap(err, "[binlogsync] withPrepareSyncer SplitHostPort")
	}
	var blSlaveID = 100
	if v, ok := c.canalParams["BinlogSlaveId"]; ok && v != "" {
		blSlaveID = conv.ToInt(c.canalParams["BinlogSlaveId"])
	}

	cfg := replication.BinlogSyncerConfig{
		ServerID: uint32(blSlaveID),
		Flavor:   c.flavor(),
		Host:     host,
		Port:     uint16(conv.ToInt(port)),
		User:     c.dsn.User,
		Password: c.dsn.Passwd,
	}
	c.syncer = replication.NewBinlogSyncer(&cfg)
	return nil
}

func withCheckBinlogRowFormat(c *Canal) error {
	v := csdb.Variable{}
	if err := v.LoadOne(context.TODO(), c.db, "binlog_format"); err != nil {
		return errors.Wrap(err, "[binlogsync] checkBinlogRowFormat row.Scan")
	}
	if v.Value != "ROW" {
		return errors.NewNotSupportedf("[binlogsync] binlog variable %q must have the configured ROW format, but got %q", v.Name, v.Value)
	}
	return nil
}

var customMySQLParams = []string{"BinlogStartPosition", "BinlogSlaveId", "flavor"}

// NewCanal creates a new canal object to start reading the MySQL binary log. If
// you don't provide a database connection option this function will panic.
// export CS_DSN='root:PASSWORD@tcp(localhost:3306)/DATABASE_NAME?BinlogSlaveId=100&BinlogStartPosition=0'
func NewCanal(dsn *mysql.Config, db Option, opts ...Option) (*Canal, error) {
	c := new(Canal)
	c.dsn = dsn
	c.closed = new(int32)
	atomic.StoreInt32(c.closed, 0)

	// remove custom parameters from DSN and copy them into our own map because
	// otherwise MySQL connection fails due to unknown connection parameters.
	if c.dsn.Params != nil {
		c.canalParams = make(map[string]string)
		for _, p := range customMySQLParams {
			if v, ok := c.dsn.Params[p]; ok && v != "" {
				c.canalParams[p] = v
				delete(c.dsn.Params, p)
			}
		}

	}

	c.Tables = csdb.MustNewTables()
	c.Tables.Schema = c.dsn.DBName
	c.Log = log.BlackHole{}

	opts2 := []Option{db}
	opts2 = append(opts2, opts...)
	opts2 = append(opts2, withUpdateBinlogStartPosition, withPrepareSyncer, withCheckBinlogRowFormat)

	for _, opt := range opts2 {
		if err := opt(c); err != nil {
			return nil, errors.Wrap(err, "[binlogsync] Applied options")
		}
	}

	return c, nil
}

func (m *Canal) masterSave() error {

	n := time.Now()
	if n.Sub(m.masterLastSaveTime) < time.Second {
		return nil
	}
	m.masterMu.Lock()
	defer m.masterMu.Unlock()

	//if err := m.cfgw.Write(); err != nil {
	//	return errors.Wrap(err, "[binlogsync] failed to write into config")
	//}

	m.masterLastSaveTime = n

	return nil
}

func (m *Canal) masterUpdate(fileName string, pos uint) {
	m.masterMu.Lock()
	defer m.masterMu.Unlock()
	m.masterStatus.File = fileName
	m.masterStatus.Position = pos
}

func (m *Canal) SyncedPosition() csdb.MasterStatus {
	m.masterMu.RLock()
	defer m.masterMu.RUnlock()
	return m.masterStatus
}

func (c *Canal) Start(ctx context.Context) error {
	c.wg.Add(1)
	go c.run(ctx)

	return nil
}

// run gets executed in its own goroutine
func (c *Canal) run(ctx context.Context) error {
	// refactor for better error handling
	defer c.wg.Done()

	if err := c.startSyncBinlog(ctx); err != nil {
		if !c.isClosed() {
			c.Log.Info("[binlogsync] Canal start has encountered a sync binlog error", log.Err(err))
		}
		return errors.Wrap(err, "[binlogsync] run.startSyncBinlog")
	}
	return nil
}

func (c *Canal) isClosed() bool {
	return atomic.LoadInt32(c.closed) == int32(1)
}

// Close closes all underlying connections
func (c *Canal) Close() error {
	c.mclose.Lock()
	defer c.mclose.Unlock()

	if c.isClosed() {
		return nil
	}

	atomic.StoreInt32(c.closed, 1)

	if c.syncer != nil {
		c.syncer.Close()
		c.syncer = nil
	}

	if err := c.db.Close(); err != nil {
		return errors.Wrap(err, "[binlogsync] DB close error")
	}
	c.wg.Wait()
	return nil
}

// FindTable tries to find a table by its ID. If the table cannot be found, it
// will add the table to the internal map and performs a column load from the
// information_schema.
func (c *Canal) FindTable(ctx context.Context, id int, tableName string) (*csdb.Table, error) {

	t, err := c.Tables.Table(id)
	if err == nil {
		return t, nil
	}
	if !errors.IsNotFound(err) {
		return nil, errors.Wrapf(err, "[binlogsync] FindTable.Table error")
	}

	if err := c.Tables.Options(csdb.WithTableLoadColumns(ctx, c.db, id, tableName)); err != nil {
		return nil, errors.Wrapf(err, "[binlogsync] FindTable.WithTableLoadColumns error")
	}

	t, err = c.Tables.Table(id)
	t.Schema = c.dsn.DBName
	if err != nil {
		return nil, errors.Wrapf(err, "[binlogsync] FindTable.Table2 error")
	}

	return t, nil
}

// Check MySQL binlog row image, must be in FULL, MINIMAL, NOBLOB
func (c *Canal) CheckBinlogRowImage(ctx context.Context, image string) error {
	// need to check MySQL binlog row image? full, minimal or noblob?
	// now only log
	if c.flavor() == MySQLFlavor {
		var v csdb.Variable
		if err := v.LoadOne(ctx, c.db, "binlog_row_image"); err != nil {
			return errors.Wrap(err, "[binlogsync] CheckBinlogRowImage LoadOne")
		}

		// MySQL has binlog row image from 5.6, so older will return empty
		if v.Value != "" && !strings.EqualFold(v.Value, image) {
			return errors.NewNotSupportedf("[binlogsync] MySQL uses %q binlog row image, but we want %q", v.Value, image)
		}
	}

	return nil
}

func (c *Canal) flavor() string {
	var f string
	if v, ok := c.canalParams["flavor"]; ok && v != "" {
		f = c.canalParams["flavor"]
	}
	if f == "" {
		f = MySQLFlavor
	}
	switch f {
	case MySQLFlavor:
		return MySQLFlavor
	case MariaDBFlavor:
		return MariaDBFlavor
	}
	// todo remove panic
	panic(fmt.Sprintf("[binlogsync] Unknown flavor: %q", f))
}
