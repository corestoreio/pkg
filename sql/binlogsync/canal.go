package binlogsync

import (
	"bytes"
	"context"
	"database/sql"
	"net"
	"regexp"
	"sync"
	"sync/atomic"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/myreplicator"
	"github.com/corestoreio/pkg/sync/singleflight"
	"github.com/corestoreio/pkg/util/conv"
	"github.com/go-sql-driver/mysql"
)

// Use flavor for different MySQL versions,
const (
	MySQLFlavor   = "mysql"
	MariaDBFlavor = "mariadb"
)

var ConfigPathBackendPosition = config.MustNewPath(`sql/binlogsync/master_position`)

// Canal can sync your MySQL data. MySQL must use the binlog format ROW.
type Canal struct {
	// mclose acts only during the call to Close().
	mclose sync.Mutex
	// DSN contains the parsed DSN
	DSN         *mysql.Config
	canalParams map[string]string

	cfgBackend config.Setter

	masterMu           sync.RWMutex
	masterStatus       ddl.MasterStatus
	masterLastSaveTime time.Time

	// expAlterTable defines the regex to be used to detect ALTER TABLE
	// statements to reinitialize the internal table structure cache.
	expAlterTable *regexp.Regexp
	syncer        *myreplicator.BinlogSyncer

	rsMu       sync.RWMutex
	rsHandlers []RowsEventHandler

	// dbcp is a database connection pool
	dbcp *dml.ConnPool

	// Tables contains the overall SQL table cache. If a table gets modified
	// during runtime of this program then somehow we must clear the cache to
	// reload the table structures.
	tables *ddl.Tables
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
		dbc, err := dml.NewConnPool(dml.WithDSN(c.DSN.FormatDSN()), dml.WithVerifyConnection())
		if err != nil {
			return errors.Wrap(err, "[binlogsync] sql.Open")
		}
		c.dbcp = dbc
		return nil
	}
}

// WithDB allows to set your own DB connection.
func WithDB(db *sql.DB) Option {
	return func(c *Canal) (err error) {
		if err = db.Ping(); err != nil {
			return errors.Wrap(err, "[binlogsync] sql ping failed")
		}
		c.dbcp, err = dml.NewConnPool(dml.WithDB(db))
		return err
	}
}

// WithConfigurationSetter used to persists the current binlog position.
// If discarded, the last position won't be saved.
func WithConfigurationSetter(s config.Setter) Option {
	return func(c *Canal) error {
		c.cfgBackend = s
		return nil
	}
}

// TODO(CyS) add a WithContext() option function or just only a parameter for a time out.

func withUpdateBinlogStart(c *Canal) error {
	ctx := context.TODO()
	var ms ddl.MasterStatus

	if _, err := c.dbcp.WithQueryBuilder(&ms).Load(ctx, &ms); err != nil {
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
	const varName = "binlog_format"
	ctx := context.Background()

	v := ddl.NewVariables(varName)
	if _, err := c.dbcp.WithQueryBuilder(v).Load(ctx, v); err != nil {
		return errors.Wrap(err, "[binlogsync] checkBinlogRowFormat row.Scan")
	}
	if !v.EqualFold(varName, "ROW") {
		return errors.NotSupported.Newf("[binlogsync] binlog variable %q must have the configured ROW format, but got %q", varName, v.Data[varName])
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
	c.expAlterTable = regexp.MustCompile("(?i)^ALTER\\sTABLE\\s.*?`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}\\s.*")

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

	c.tables = ddl.MustNewTables()
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

// TODO continue sync from last stored master position

func (c *Canal) masterSave(fileName string, pos uint) error {
	c.masterMu.Lock()
	defer c.masterMu.Unlock()

	c.masterStatus.File = fileName
	c.masterStatus.Position = pos

	now := time.Now()
	if now.Sub(c.masterLastSaveTime) < time.Second {
		return nil
	}

	if c.cfgBackend == nil {
		if c.Log.IsDebug() {
			c.Log.Debug("[binlogsync] Warning: Master Status cannot be saved because config.Setter is nil",
				log.String("database", c.DSN.DBName), log.Stringer("master_status", c.masterStatus))
		}
		return nil
	}

	var buf bytes.Buffer
	c.masterStatus.WriteTo(&buf)
	if err := c.cfgBackend.Set(ConfigPathBackendPosition, buf.Bytes()); err != nil {
		if c.Log.IsInfo() {
			c.Log.Info("[binlogsync] Failed to store Master Status",
				log.Time("master_last_save_time", c.masterLastSaveTime),
				log.Err(err), log.String("database", c.DSN.DBName), log.Stringer("master_status", c.masterStatus))
		}
		return errors.WithStack(err)
	}

	c.masterLastSaveTime = now

	return nil
}

// SyncedPosition returns the current synced position as retrieved from the SQl
// server.
func (c *Canal) SyncedPosition() ddl.MasterStatus {
	c.masterMu.RLock()
	defer c.masterMu.RUnlock()
	return c.masterStatus
}

// Start starts the sync process in the background as a goroutine. You can stop
// the goroutine via the context.
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

	if err := c.dbcp.Close(); err != nil {
		return errors.Wrap(err, "[binlogsync] DB close error")
	}
	c.wg.Wait()
	return nil
}

// FindTable tries to find a table by its ID. If the table cannot be found by
// the first search, it will add the table to the internal map and performs a
// column load from the information_schema and then returns the fully defined
// table.
func (c *Canal) FindTable(ctx context.Context, tableName string) (ddl.Table, error) {
	// deference the table pointer to avoid race conditions and devs modifying the
	// table ;-)
	t, err := c.tables.Table(tableName)
	if err == nil {
		return *t, nil
	}
	if !errors.Is(err, errors.NotFound) {
		return ddl.Table{}, errors.Wrapf(err, "[binlogsync] FindTable.Table error")
	}

	val, err, _ := c.tableSFG.Do(tableName, func() (interface{}, error) {
		if err := c.tables.Options(ddl.WithTableLoadColumns(ctx, c.dbcp.DB, tableName)); err != nil {
			return ddl.Table{}, errors.Wrapf(err, "[binlogsync] FindTable.WithTableLoadColumns error")
		}

		t, err = c.tables.Table(tableName)
		if err != nil {
			return ddl.Table{}, errors.Wrapf(err, "[binlogsync] FindTable.Table2 error")
		}
		return *t, nil
	})

	if err != nil {
		return ddl.Table{}, errors.Wrapf(err, "[binlogsync] FindTable.SingleFlight error")
	}

	return val.(ddl.Table), nil
}

// ClearTableCache clear table cache
func (c *Canal) ClearTableCache(db string, table string) {
	// TODO implement
	// c.tables.DeleteAllFromCache()
	//key := fmt.Sprintf("%s.%s", dbcp, table)
	//c.tableLock.Lock()
	//delete(c.tables, key)
	//c.tableLock.Unlock()
}

// CheckBinlogRowImage checks MySQL binlog row image, must be in FULL, MINIMAL, NOBLOB
func (c *Canal) CheckBinlogRowImage(ctx context.Context, image string) error {
	// need to check MySQL binlog row image? full, minimal or noblob?
	// now only log.
	//  TODO what about MariaDB?
	const varName = "binlog_row_image"
	if c.flavor() == MySQLFlavor {
		v := ddl.NewVariables(varName)
		if _, err := c.dbcp.WithQueryBuilder(v).Load(ctx, v); err != nil {
			return errors.Wrap(err, "[binlogsync] CheckBinlogRowImage LoadOne")
		}

		// MySQL has binlog row image from 5.6, so older will return empty
		if v.EqualFold(varName, image) {
			return errors.NotSupported.Newf("[binlogsync] MySQL uses %q binlog row image, but we want %q", v.Data[varName], image)
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
