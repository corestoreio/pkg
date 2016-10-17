package binlogsync

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/siddontang/go-mysql/replication"
	"github.com/siddontang/go-mysql/schema"
)

// Use flavor for different MySQL versions,
const (
	MySQLFlavor   = "mysql"
	MariaDBFlavor = "mariadb"
)

// db matches type database/sql.DB
type DBQueryier interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}

// Canal can sync your MySQL data. MySQL must use the binlog format ROW.
type Canal struct {
	// mclose acts only during the call to Close().
	mclose sync.Mutex
	// dsn contains the parsed DSN
	dsn *url.URL
	// database sets the name to which the current Canal object has been bound to.
	database string

	cfgw config.Writer

	masterMu           sync.RWMutex
	masterPos          Position
	masterLastSaveTime time.Time

	syncer *replication.BinlogSyncer

	rsMu       sync.RWMutex
	rsHandlers []RowsEventHandler

	db DBQueryier

	tableLock sync.RWMutex
	tables    map[string]*schema.Table

	closed *int32
	Log    log.Logger
	wg     sync.WaitGroup
}

// Option applies multiple options to the Canal type.
type Option func(*Canal) error

// WithMySQL adds the database/sql.DB driver including a ping to the database.
func WithMySQL() Option {
	return func(c *Canal) error {
		db, err := sql.Open("mysql", c.dsn.String())
		if err != nil {
			return errors.Wrap(err, "[binlogsync] sql.Open")
		}
		if err := db.Ping(); err != nil {
			return errors.Wrap(err, "[binlogsync] sql.Ping")
		}
		c.db = db
		return nil
	}
}

// WithDB allows to set your own DB connection. Mostly used for testing.
func WithDB(db DBQueryier) Option {
	return func(c *Canal) error {
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
	ms, err := c.ShowMasterStatus()
	if err != nil {
		return errors.Wrap(err, "[binlogsync] ShowMasterStatus")
	}

	c.masterPos = ms

	hasPos := conv.ToInt(c.dsn.Query().Get("BinlogStartPosition"))
	if hasPos >= 4 {
		c.masterPos.Pos = uint32(hasPos)
	}
	return nil
}

func withPrepareSyncer(c *Canal) error {
	pw, _ := c.dsn.User.Password()
	cfg := replication.BinlogSyncerConfig{
		ServerID: uint32(conv.ToInt(c.dsn.Query().Get("BinlogSlaveId"))),
		Flavor:   c.flavor(),
		Host:     c.dsn.Hostname(),
		Port:     uint16(conv.ToInt(c.dsn.Port())),
		User:     c.dsn.User.Username(),
		Password: pw,
	}
	c.syncer = replication.NewBinlogSyncer(&cfg)
	return nil
}

func withCheckBinlogRowFormat(c *Canal) error {
	row := c.db.QueryRow(`SHOW GLOBAL VARIABLES LIKE "binlog_format";`)
	res := &struct {
		VariableName, Value string
	}{}
	if err := row.Scan(res.VariableName, res.Value); err != nil {
		return errors.Wrap(err, "[binlogsync] checkBinlogRowFormat row.Scan")
	}
	if res.Value != "ROW" {
		return errors.NewNotSupportedf("[binlogsync] binlog must have the configured ROW format, but got %q", res.Value)
	}
	return nil
}

// NewCanal creates a new canal object to start reading the MySQL binary log. If
// you don't provide a database connection option this function will panic.
// export CS_DSN='mysql://root:PASSWORD@localhost:3306/DATABASE_NAME?BinlogSlaveId=100&BinlogStartPosition=0'
func NewCanal(dsn *url.URL, db Option, opts ...Option) (*Canal, error) {
	c := new(Canal)
	c.dsn = dsn
	c.closed = new(int32)
	atomic.StoreInt32(c.closed, 0)

	c.database = dsn.Path[1:] // strip first slash, let it panic if the DSN has been provided incorrectly.
	c.tables = make(map[string]*schema.Table)
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

func (m *Canal) masterUpdate(name string, pos uint32) {
	m.masterMu.Lock()
	defer m.masterMu.Unlock()
	m.masterPos.Name = name
	m.masterPos.Pos = pos
}

func (m *Canal) SyncedPosition() Position {
	m.masterMu.RLock()
	defer m.masterMu.RUnlock()

	return m.masterPos
}

func (c *Canal) Start() error {
	c.wg.Add(1)
	go c.run()

	return nil
}

// run gets executed in its own goroutine
func (c *Canal) run() error {
	// refactor for better error handling
	defer c.wg.Done()

	if err := c.startSyncBinlog(); err != nil {
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

func (c *Canal) FlushTables() error {
	c.tableLock.Lock()
	defer c.tableLock.Unlock()
	c.tables = make(map[string]*schema.Table)

	return nil
}
func (c *Canal) GetTable(table string) (*schema.Table, error) {
	key := c.database + "." + table
	c.tableLock.RLock()
	t, ok := c.tables[key]
	c.tableLock.RUnlock()

	if ok {
		return t, nil
	}

	t, err := schema.NewTable(c, c.database, table)
	if err != nil {
		return nil, errors.Wrapf(err, "[binlogsync] GetTable schema.NewTable")
	}

	c.tableLock.Lock()
	c.tables[key] = t
	c.tableLock.Unlock()

	return t, nil
}

// Check MySQL binlog row image, must be in FULL, MINIMAL, NOBLOB
func (c *Canal) CheckBinlogRowImage(image string) error {
	// need to check MySQL binlog row image? full, minimal or noblob?
	// now only log
	if c.flavor() == MySQLFlavor {
		row := c.db.QueryRow(`SHOW GLOBAL VARIABLES LIKE "binlog_row_image"`)
		res := &struct{ VariableName, Value string }{}
		if err := row.Scan(res.VariableName, res.Value); err != nil {
			return errors.Wrap(err, "[binlogsync] CheckBinlogRowImage Execute")
		}
		// MySQL has binlog row image from 5.6, so older will return empty
		if res.Value != "" && !strings.EqualFold(res.Value, image) {
			return errors.NewNotSupportedf("[binlogsync] MySQL uses %q binlog row image, but we want %q", res.Value, image)
		}
	}

	return nil
}

func (c *Canal) flavor() string {
	f := c.dsn.Query().Get("flavor")
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

func (c *Canal) ShowMasterStatus() (p Position, _ error) {
	row := c.db.QueryRow("SHOW MASTER STATUS")
	res := &struct {
		VariableName, Value string
	}{}
	if err := row.Scan(res.VariableName, res.Value); err != nil {
		return p, errors.Wrap(err, "[binlogsync] ShowMasterStatus")
	}

	pos, err := strconv.ParseUint(res.Value, 10, 32)
	if err != nil {
		return p, errors.Wrap(err, "[binlogsync] ShowMasterStatus ParseUint")
	}
	return Position{Name: res.VariableName, Pos: uint32(pos)}, nil
}
