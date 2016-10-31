package binlogsync

import (
	"context"
	"database/sql"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/myreplicator"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sync/singleflight"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/go-sql-driver/mysql"
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
	// BackendPosition initial idea. writing supported but not loading
	BackendPosition cfgmodel.Str

	// mclose acts only during the call to Close().
	mclose sync.Mutex
	// DSN contains the parsed DSN
	DSN         *mysql.Config
	canalParams map[string]string

	cfgw config.Writer

	masterMu           sync.RWMutex
	masterStatus       csdb.MasterStatus
	masterLastSaveTime time.Time

	syncer *myreplicator.BinlogSyncer

	rsMu       sync.RWMutex
	rsHandlers []RowsEventHandler

	db DBer

	// Tables contains the overall SQL table cache. If a table gets modified
	// during runtime of this program then somehow we must clear the cache to
	// reload the table structures.
	tables *csdb.Tables
	// tableSFG takes to only execute one SQL query per table in parallel
	// situations. No need for a pointer because Canal is already a pointer. So
	// simple embedding.
	tableSFG *singleflight.Group

	closed *int32
	Log    log.Logger
	wg     sync.WaitGroup
}

// Option applies multiple options to the Canal type.
type Option func(*Canal) error

// WithMySQL adds the database/sql.DB driver including a ping to the database.
func WithMySQL() Option {
	return func(c *Canal) error {
		db, err := sql.Open("mysql", c.DSN.FormatDSN())
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

// WithDB allows to set your own DB connection.
func WithDB(db *sql.DB) Option {
	return func(c *Canal) (err error) {
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, "[binlogsync] sql ping failed")
		}
		c.db = db
		return nil
	}
}

// WithConfigurationWriter used to persists the current binlog position.
func WithConfigurationWriter(w config.Writer) Option {
	//Write(p cfgpath.Path, value interface{}) error
	return func(c *Canal) error {
		c.cfgw = w
		return nil
	}
}

func withUpdateBinlogStart(c *Canal) error {
	var ms csdb.MasterStatus
	if err := ms.Load(context.TODO(), c.db); err != nil {
		return errors.Wrap(err, "[binlogsync] ShowMasterStatus Load")
	}

	c.masterStatus = ms

	if v, ok := c.canalParams["BinlogStartFile"]; ok && v != "" {
		c.masterStatus.File = v
	}
	if v, ok := c.canalParams["BinlogStartPosition"]; ok && v != "" {
		if hasPos := conv.ToUint(v); hasPos >= 4 {
			c.masterStatus.Position = hasPos
		}
	}
	return nil
}

// withPrepareSyncer creates its own database connection.
func withPrepareSyncer(c *Canal) error {
	host, port, err := net.SplitHostPort(c.DSN.Addr)
	if err != nil {
		return errors.Wrap(err, "[binlogsync] withPrepareSyncer SplitHostPort")
	}
	var blSlaveID = 100
	if v, ok := c.canalParams["BinlogSlaveId"]; ok && v != "" {
		blSlaveID = conv.ToInt(v)
	}

	cfg := myreplicator.BinlogSyncerConfig{
		ServerID: uint32(blSlaveID),
		Flavor:   c.flavor(),
		Host:     host,
		Port:     uint16(conv.ToInt(port)),
		User:     c.DSN.User,
		Password: c.DSN.Passwd,
	}
	c.syncer = myreplicator.NewBinlogSyncer(&cfg)
	return nil
}

func withCheckBinlogRowFormat(c *Canal) error {
	v := csdb.Variable{}
	if err := v.LoadOne(context.TODO(), c.db, "binlog_format"); err != nil {
		return errors.Wrap(err, "[binlogsync] checkBinlogRowFormat row.Scan")
	}
	if !strings.EqualFold(v.Value, "ROW") {
		return errors.NewNotSupportedf("[binlogsync] binlog variable %q must have the configured ROW format, but got %q", v.Name, v.Value)
	}
	return nil
}

var customMySQLParams = []string{"BinlogStartFile", "BinlogStartPosition", "BinlogSlaveId", "flavor"}

// NewCanal creates a new canal object to start reading the MySQL binary log. If
// you don't provide a database connection option this function will panic.
// export CS_DSN='root:PASSWORD@tcp(localhost:3306)/DATABASE_NAME?BinlogSlaveId=100&BinlogStartFile=mysql-bin.000002&BinlogStartPosition=4'
func NewCanal(dsn *mysql.Config, db Option, opts ...Option) (*Canal, error) {
	c := new(Canal)
	c.DSN = dsn
	c.closed = new(int32)
	atomic.StoreInt32(c.closed, 0)

	c.BackendPosition = cfgmodel.NewStr("storage/binlogsync/position")

	// remove custom parameters from DSN and copy them into our own map because
	// otherwise MySQL connection fails due to unknown connection parameters.
	if c.DSN.Params != nil {
		c.canalParams = make(map[string]string)
		for _, p := range customMySQLParams {
			if v, ok := c.DSN.Params[p]; ok && v != "" {
				c.canalParams[p] = v
				delete(c.DSN.Params, p)
			}
		}

	}

	c.tables = csdb.MustNewTables()
	c.tables.Schema = c.DSN.DBName
	c.tableSFG = new(singleflight.Group)
	c.Log = log.BlackHole{}

	opts2 := []Option{db}
	opts2 = append(opts2, opts...)
	opts2 = append(opts2, withUpdateBinlogStart, withPrepareSyncer, withCheckBinlogRowFormat)

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

	if m.cfgw == nil {
		if m.Log.IsDebug() {
			m.Log.Debug("[binlogsync] Master Status cannot be saved because config.Writer is nil",
				log.String("database", m.DSN.DBName), log.Stringer("master_status", m.masterStatus))
		}
		return nil
	}

	// todo refactor to find a different way by not importing package config and scope
	if err := m.BackendPosition.Write(m.cfgw, m.masterStatus.String(), scope.DefaultTypeID); err != nil {
		return errors.Wrap(err, "[binlogsync] failed to write into config")
	}

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

// FindTable tries to find a table by its ID. If the table cannot be found by
// the first search, it will add the table to the internal map and performs a
// column load from the information_schema and then returns the fully defined
// table.
func (c *Canal) FindTable(ctx context.Context, id int, tableName string) (csdb.Table, error) {
	// deference the table pointer to avoid race conditions and devs modifying the
	// table ;-)
	t, err := c.tables.Table(id)
	if err == nil {
		return *t, nil
	}
	if !errors.IsNotFound(err) {
		return csdb.Table{}, errors.Wrapf(err, "[binlogsync] FindTable.Table error")
	}

	val, err, _ := c.tableSFG.Do(tableName, func() (interface{}, error) {
		if err := c.tables.Options(csdb.WithTableLoadColumns(ctx, c.db, id, tableName)); err != nil {
			return csdb.Table{}, errors.Wrapf(err, "[binlogsync] FindTable.WithTableLoadColumns error")
		}

		t, err = c.tables.Table(id)
		if err != nil {
			return csdb.Table{}, errors.Wrapf(err, "[binlogsync] FindTable.Table2 error")
		}
		return *t, nil
	})

	if err != nil {
		return csdb.Table{}, errors.Wrapf(err, "[binlogsync] FindTable.SingleFlight error")
	}

	return val.(csdb.Table), nil
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
		f = v
	}
	if f == "" {
		f = MySQLFlavor
	}
	switch f {
	case MariaDBFlavor:
		return MariaDBFlavor
	}
	return MySQLFlavor
}
