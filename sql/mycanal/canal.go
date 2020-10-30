package mycanal

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
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
	"github.com/corestoreio/pkg/store/scope"
	"golang.org/x/sync/singleflight"
	"github.com/corestoreio/pkg/util/conv"
	"github.com/go-sql-driver/mysql"
	simysql "github.com/siddontang/go-mysql/mysql"
)

// Use flavor for different MySQL versions,
const (
	MySQLFlavor   = "mysql"
	MariaDBFlavor = "mariadb"
)

// Configuration paths for config.Service
const (
	ConfigPathBackendPosition     = `sql/mycanal/master_position`
	ConfigPathIncludeTableRegex   = `sql/mycanal/include_table_regex`
	ConfigPathExcludeTableRegex   = `sql/mycanal/exclude_table_regex`
	ConfigPathBinlogStartFile     = `sql/mycanal/binlog_start_file`
	ConfigPathBinlogStartPosition = `sql/mycanal/binlog_start_position`
	ConfigPathBinlogSlaveID       = `sql/mycanal/binlog_slave_id`
	ConfigPathServerFlavor        = `sql/mycanal/server_flavor`
)

// Canal can sync your MySQL data. MySQL must use the binlog format ROW.
type Canal struct {
	opts                      Options
	configPathBackendPosition *config.Path
	// mclose acts only during the call to Close().
	mclose sync.Mutex
	// DSN contains the parsed DSN
	dsn *mysql.Config

	cfgScope config.Scoped // required
	syncer   *myreplicator.BinlogSyncer

	masterMu                          sync.RWMutex
	masterStatus                      ddl.MasterStatus
	masterGTID                        simysql.GTIDSet
	masterLastSaveTime                time.Time
	masterWarningConfigSetIsNilLogged bool

	rsMu sync.RWMutex
	// the empty map key declares event handler for all tables, filtered  by the regexes.
	// Otherwise an event handler is only registered for a specific table.
	rsHandlers map[string][]RowsEventHandler

	// dbcp is a database connection pool
	dbcp *dml.ConnPool

	// Tables contains the overall SQL table cache. If a table gets modified
	// during runtime of this program then somehow we must clear the cache to
	// reload the table structures.
	tables *ddl.Tables
	// tableSFG takes to only execute one SQL query per table in parallel
	// situations. No need for a pointer because Canal is already a pointer. So
	// simple embedding.
	tableSFG singleflight.Group

	tableLock         sync.RWMutex
	tableAllowedCache map[string]bool
	includeTableRegex []*regexp.Regexp
	excludeTableRegex []*regexp.Regexp

	closed *int32
	wg     sync.WaitGroup
}

// DBConFactory creates a new database connection.
type DBConFactory func(dsn string) (*dml.ConnPool, error)

// WithMySQL adds the database/sql.DB driver including a ping to the database
// from the provided DSN.
func WithMySQL() DBConFactory {
	return func(dsn string) (*dml.ConnPool, error) {
		dbc, err := dml.NewConnPool(dml.WithDSN(dsn), dml.WithVerifyConnection())
		return dbc, errors.WithStack(err)
	}
}

// WithDB allows to set your own DB connection.
func WithDB(db *sql.DB) DBConFactory {
	return func(_ string) (*dml.ConnPool, error) {
		if err := db.Ping(); err != nil {
			return nil, errors.WithStack(err)
		}
		dbc, err := dml.NewConnPool(dml.WithDB(db))
		return dbc, errors.WithStack(err)
	}
}

func withIncludeTables(regexes []string) func(c *Canal) error {
	return func(c *Canal) error {
		if len(regexes) == 0 {
			return nil
		}
		c.tableLock.Lock()
		defer c.tableLock.Unlock()
		c.includeTableRegex = make([]*regexp.Regexp, len(regexes))
		for i, val := range regexes {
			reg, err := regexp.Compile(val)
			if err != nil {
				return errors.WithStack(err)
			}
			c.includeTableRegex[i] = reg
		}
		return nil
	}
}

func withExcludeTables(regexes []string) func(c *Canal) error {
	return func(c *Canal) error {
		if len(regexes) == 0 {
			return nil
		}
		c.tableLock.Lock()
		defer c.tableLock.Unlock()
		c.excludeTableRegex = make([]*regexp.Regexp, len(regexes))
		for i, val := range regexes {
			reg, err := regexp.Compile(val)
			if err != nil {
				return errors.WithStack(err)
			}
			c.excludeTableRegex[i] = reg
		}
		return nil
	}
}

// withUpdateBinlogStart enables to start from a specific position or just start
// from the current master position. See startSyncBinlog
func withUpdateBinlogStart(c *Canal) error {

	if c.opts.BinlogStartFile != "" && c.opts.BinlogStartPosition > 0 {
		c.masterStatus.File = c.opts.BinlogStartFile
		c.masterStatus.Position = uint(c.opts.BinlogStartPosition)
		return nil
	}

	if c.opts.MasterStatusQueryTimeout == 0 {
		c.opts.MasterStatusQueryTimeout = time.Second * 20
	}

	var ms ddl.MasterStatus
	ctx, cancel := context.WithTimeout(context.Background(), c.opts.MasterStatusQueryTimeout)
	defer cancel()
	if _, err := c.dbcp.WithQueryBuilder(&ms).Load(ctx, &ms); err != nil {
		return errors.WithStack(err)
	}

	c.masterStatus = ms

	return nil
}

// withPrepareSyncer creates its own database connection.
func withPrepareSyncer(c *Canal) error {

	host, port, err := net.SplitHostPort(c.dsn.Addr)
	if err != nil {
		return errors.Wrapf(err, "[mycanal] withPrepareSyncer SplitHostPort %q", c.dsn.Addr)
	}

	if c.opts.BinlogSlaveId == 0 {
		c.opts.BinlogSlaveId = 100
	}

	cfg := myreplicator.BinlogSyncerConfig{
		ServerID:             uint32(c.opts.BinlogSlaveId),
		Flavor:               c.opts.Flavor,
		Host:                 host,
		Port:                 uint16(conv.ToUint(port)),
		User:                 c.dsn.User,
		Password:             c.dsn.Passwd,
		Log:                  c.opts.Log,
		TLSConfig:            c.opts.TLSConfig,
		MaxReconnectAttempts: c.opts.MaxReconnectAttempts,
	}

	c.syncer = myreplicator.NewBinlogSyncer(&cfg)

	return nil
}

func withCheckBinlogRowFormat(c *Canal) error {
	const varName = "binlog_format"
	ctx := context.Background()

	v := ddl.NewVariables(varName)
	if _, err := c.dbcp.WithQueryBuilder(v).Load(ctx, v); err != nil {
		return errors.WithStack(err)
	}
	if !v.EqualFold(varName, "ROW") {
		return errors.NotSupported.Newf("[mycanal] binlog variable %q must have the configured ROW format, but got %q. ROW means: Records events affecting individual table rows.", varName, v.Data[varName])
	}
	return nil
}

// Options provides multiple options for NewCanal. Part of those options can get
// loaded via config.Scoped.
type Options struct {
	// ConfigScoped defines the configuration to load the following fields from.
	// If not set the data won't be loaded.
	ConfigScoped config.Scoped
	// ConfigSet used to persists the master position of the binlog stream.
	ConfigSet config.Setter
	Log       log.Logger
	TLSConfig *tls.Config // Needs some rework
	// IncludeTableRegex defines the regex which matches the allowed table
	// names. Default state of WithIncludeTables is empty, this will include all
	// tables.
	IncludeTableRegex []string
	// ExcludeTableRegex defines the regex which matches the excluded table
	// names. Default state of WithExcludeTables is empty, ignores exluding and
	// includes all tables.
	ExcludeTableRegex []string

	// Set to change the maximum number of attempts to re-establish a broken
	// connection
	MaxReconnectAttempts int

	BinlogStartFile     string
	BinlogStartPosition uint64
	BinlogSlaveId       uint64
	// Flavor defines if `mariadb` or `mysql` should be used. Defaults to
	// `mariadb`.
	Flavor                   string
	MasterStatusQueryTimeout time.Duration
	// OnClose runs before the database connection gets closed and after the
	// syncer has been closed. The syncer does not "see" the changes comming
	// from the queries executed in the call back.
	OnClose func(*dml.ConnPool) error
}

func (o *Options) loadFromConfigService() (err error) {
	defer func() {
		switch o.Flavor {
		case MySQLFlavor:
			o.Flavor = MySQLFlavor
		default:
			o.Flavor = MariaDBFlavor
		}
	}()

	if o.Log == nil {
		o.Log = log.BlackHole{}
	}

	if !o.ConfigScoped.IsValid() {
		return nil
	}

	if o.IncludeTableRegex == nil {
		v := o.ConfigScoped.Get(scope.Default, ConfigPathIncludeTableRegex)
		if o.IncludeTableRegex, err = v.Strs(o.IncludeTableRegex...); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	if o.ExcludeTableRegex == nil {
		v := o.ConfigScoped.Get(scope.Default, ConfigPathExcludeTableRegex)
		if o.ExcludeTableRegex, err = v.Strs(o.ExcludeTableRegex...); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	if o.BinlogStartFile == "" {
		v := o.ConfigScoped.Get(scope.Default, ConfigPathBinlogStartFile)
		o.BinlogStartFile = v.UnsafeStr()
	}
	if o.BinlogStartPosition == 0 {
		v := o.ConfigScoped.Get(scope.Default, ConfigPathBinlogStartPosition)
		o.BinlogStartPosition = v.UnsafeUint64()
	}
	if o.BinlogSlaveId == 0 {
		v := o.ConfigScoped.Get(scope.Default, ConfigPathBinlogSlaveID)
		o.BinlogSlaveId = v.UnsafeUint64()
	}
	if o.Flavor == "" {
		v := o.ConfigScoped.Get(scope.Default, ConfigPathServerFlavor)
		o.Flavor = v.UnsafeStr()
	}

	return nil
}

// NewCanal creates a new canal object to start reading the MySQL binary log.
// The DSN is need to setup two different connections. One connection for
// reading the binary stream and the 2nd connection to execute queries. The 2nd
// argument `db` gets used to executed the queries, like setting variables or
// getting table information. Default database flavor is `mariadb`.
// export CS_DSN='root:PASSWORD@tcp(localhost:3306)/DATABASE_NAME
func NewCanal(dsn string, db DBConFactory, opt *Options) (*Canal, error) {
	pDSN, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := opt.loadFromConfigService(); err != nil {
		return nil, errors.WithStack(err)
	}

	c := &Canal{
		opts:                      *opt,
		configPathBackendPosition: config.MustMakePath(ConfigPathBackendPosition),
		dsn:                       pDSN,
		closed:                    new(int32),
		tables:                    ddl.MustNewTables(),
	}

	atomic.StoreInt32(c.closed, 0)

	c.tables.Schema = c.dsn.DBName

	c.dbcp, err = db(dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := c.tables.Options(ddl.WithConnPool(c.dbcp)); err != nil {
		return nil, errors.WithStack(err)
	}

	initOptFn := [...]func(c *Canal) error{
		withUpdateBinlogStart, withPrepareSyncer, withCheckBinlogRowFormat,
		withIncludeTables(opt.IncludeTableRegex), withExcludeTables(opt.ExcludeTableRegex),
	}
	for _, optFn := range initOptFn {
		if err := optFn(c); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	c.tableLock.Lock()
	if c.includeTableRegex != nil || c.excludeTableRegex != nil {
		c.tableAllowedCache = make(map[string]bool)
	}
	c.tableLock.Unlock()

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

	if c.opts.ConfigSet == nil {
		if !c.masterWarningConfigSetIsNilLogged && c.opts.Log.IsInfo() {
			c.masterWarningConfigSetIsNilLogged = true
			c.opts.Log.Info("mycanal.masterSave.config.setter",
				log.Bool("config_setter_is_nil", true), log.String("database", c.dsn.DBName), log.Stringer("master_status", c.masterStatus))
		}
		return nil
	}

	var buf bytes.Buffer
	_, _ = c.masterStatus.WriteTo(&buf)
	if err := c.opts.ConfigSet.Set(c.configPathBackendPosition, buf.Bytes()); err != nil {
		if c.opts.Log.IsInfo() {
			c.opts.Log.Info("mycanal.masterSave.error",
				log.Time("master_last_save_time", c.masterLastSaveTime),
				log.Err(err), log.String("database", c.dsn.DBName), log.Stringer("master_status", c.masterStatus))
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
	go c.run(ctx)
	return nil
}

// run gets executed in its own goroutine
func (c *Canal) run(ctx context.Context) {
	// refactor for better error handling
	defer c.wg.Done()
	c.wg.Add(1)
	if err := c.startSyncBinlog(ctx); err != nil {
		if !c.isClosed() && c.opts.Log.IsInfo() {
			c.opts.Log.Info("[mycanal] Canal start has encountered a sync binlog error", log.Err(err))
		}
	}
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
		if err := c.syncer.Close(); err != nil {
			return errors.WithStack(err)
		}
		c.syncer = nil
	}

	if err := c.opts.OnClose(c.dbcp); err != nil {
		return errors.WithStack(err)
	}

	if err := c.dbcp.Close(); err != nil {
		return errors.WithStack(err)
	}
	c.wg.Wait()
	return nil
}

func (c *Canal) isTableAllowed(tblName string) bool {
	// no filter, return true
	if c.tableAllowedCache == nil {
		return true
	}

	c.tableLock.RLock()
	isAllowed, ok := c.tableAllowedCache[tblName]
	c.tableLock.RUnlock()
	if ok {
		// cache hit
		return isAllowed
	}
	matchFlag := false
	// check include
	if c.includeTableRegex != nil {
		for _, reg := range c.includeTableRegex {
			if reg.MatchString(tblName) {
				matchFlag = true
				break
			}
		}
	}
	// check exclude
	if matchFlag && c.excludeTableRegex != nil {
		for _, reg := range c.excludeTableRegex {
			if reg.MatchString(tblName) {
				matchFlag = false
				break
			}
		}
	}
	c.tableLock.Lock()
	c.tableAllowedCache[tblName] = matchFlag
	c.tableLock.Unlock()
	return matchFlag
}

type errTableNotAllowed string

func (t errTableNotAllowed) ErrorKind() errors.Kind { return errors.NotAllowed }
func (t errTableNotAllowed) Error() string {
	return fmt.Sprintf("[mycanal] Table %q is not allowed", string(t))
}

// FindTable tries to find a table by its ID. If the table cannot be found by
// the first search, it will add the table to the internal map and performs a
// column load from the information_schema and then returns the fully defined
// table. Only tables which are found in the database name of the DSN get
// loaded.
func (c *Canal) FindTable(ctx context.Context, tableName string) (*ddl.Table, error) {

	if !c.isTableAllowed(tableName) {
		return nil, errTableNotAllowed(tableName)
	}

	t, err := c.tables.Table(tableName)
	if err == nil {
		return t, nil
	}
	if !errors.NotFound.Match(err) {
		return nil, errors.WithStack(err)
	}

	val, err, _ := c.tableSFG.Do(tableName, func() (interface{}, error) {
		if err := c.tables.Options(ddl.WithCreateTable(ctx, tableName, "")); err != nil {
			return nil, errors.WithStack(err)
		}

		t, err = c.tables.Table(tableName)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return t, nil
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return val.(*ddl.Table), nil
}

// ClearTableCache clear table cache
func (c *Canal) ClearTableCache(db string, table string) {
	c.tables.DeleteAllFromCache()
}

// CheckBinlogRowImage checks MySQL binlog row image, must be in FULL, MINIMAL, NOBLOB
func (c *Canal) CheckBinlogRowImage(ctx context.Context, image string) error {
	// need to check MySQL binlog row image? full, minimal or noblob?
	// now only log.
	//  TODO what about MariaDB?
	const varName = "binlog_row_image"
	if c.opts.Flavor == MySQLFlavor {
		v := ddl.NewVariables(varName)
		if _, err := c.dbcp.WithQueryBuilder(v).Load(ctx, v); err != nil {
			return errors.WithStack(err)
		}

		// MySQL has binlog row image from 5.6, so older will return empty
		if v.EqualFold(varName, image) {
			return errors.NotSupported.Newf("[mycanal] MySQL uses %q binlog row image, but we want %q", v.Data[varName], image)
		}
	}
	return nil
}
