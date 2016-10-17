package binlogsync

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"github.com/siddontang/go-mysql/schema"
	"github.com/siddontang/go/sync2"
)

// Canal can sync your MySQL data. MySQL must use the binlog format ROW.
type Canal struct {
	// mclose acts only during the call to Close().
	mclose sync.Mutex
	// cfg contains the parsed DSN
	cfg *url.URL
	// database sets the name to which
	database string
	master   *masterInfo

	syncer *replication.BinlogSyncer

	rsLock     sync.RWMutex
	rsHandlers []RowsEventHandler

	connLock sync.Mutex
	conn     *client.Conn
	//SQLRow   interface {
	//	QueryRow(query string, args ...interface{}) *sql.Row
	//}

	wg sync.WaitGroup

	tableLock sync.RWMutex
	tables    map[string]*schema.Table

	closed sync2.AtomicBool
	Log    log.Logger
}

// export CS_DSN='mysql://root:PASSWORD@localhost:3306/DATABASE_NAME?BinlogSlaveId=100&BinlogStartPosition=0'
func NewCanal(cfg *url.URL) (*Canal, error) {
	c := new(Canal)
	c.cfg = cfg
	c.closed.Set(false)

	c.database = cfg.Path[1:] // strip first slash, let it panic if the DSN has been provided incorrectly.

	c.rsHandlers = make([]RowsEventHandler, 0, 4)
	c.tables = make(map[string]*schema.Table)
	c.Log = log.BlackHole{}

	var err error
	if c.master, err = loadMasterInfo(c); err != nil {
		return nil, errors.Wrap(err, "[binlogsync] loadMasterInfo")
	}
	c.updateBinlogStartPosition()

	if err = c.prepareSyncer(); err != nil {
		return nil, errors.Wrap(err, "[binlogsync] prepareSyncer")
	}

	if err := c.checkBinlogRowFormat(); err != nil {
		return nil, errors.Wrap(err, "[binlogsync] checkBinlogRowFormat")
	}

	return c, nil
}

func (c *Canal) updateBinlogStartPosition() {
	hasPos := conv.ToInt(c.cfg.Query().Get("BinlogStartPosition"))
	if hasPos >= 4 {
		c.master.Position = uint32(hasPos)
	}
}

func (c *Canal) Start() error {
	c.wg.Add(1)
	go c.run()

	return nil
}

// run gets executed in its own goroutine
func (c *Canal) run() error {
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
	return c.closed.Get()
}

func (c *Canal) Close() {
	c.mclose.Lock()
	defer c.mclose.Unlock()

	if c.isClosed() {
		return
	}

	c.closed.Set(true)

	c.connLock.Lock()
	c.conn.Close()
	c.conn = nil
	c.connLock.Unlock()

	if c.syncer != nil {
		c.syncer.Close()
		c.syncer = nil
	}

	c.master.Close()

	c.wg.Wait()
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
	if c.flavor() == mysql.MySQLFlavor {
		if res, err := c.Execute(`SHOW GLOBAL VARIABLES LIKE "binlog_row_image"`); err != nil {
			return errors.Wrap(err, "[binlogsync] CheckBinlogRowImage Execute")
		} else {
			// MySQL has binlog row image from 5.6, so older will return empty
			rowImage, _ := res.GetString(0, 1)
			if rowImage != "" && !strings.EqualFold(rowImage, image) {
				return errors.NewNotSupportedf("[binlogsync] MySQL uses %q binlog row image, but we want %q", rowImage, image)
			}
		}
	}

	return nil
}

func (c *Canal) checkBinlogRowFormat() error {
	res, err := c.Execute(`SHOW GLOBAL VARIABLES LIKE "binlog_format";`)
	if err != nil {
		return errors.Wrap(err, "[binlogsync] checkBinlogRowFormat Execute")
	} else if f, _ := res.GetString(0, 1); f != "ROW" {
		return errors.NewNotSupportedf("[binlogsync] binlog must have the ROW format, but got %q", f)
	}

	return nil
}

func (c *Canal) prepareSyncer() error {

	pw, _ := c.cfg.User.Password()
	cfg := replication.BinlogSyncerConfig{
		ServerID: uint32(conv.ToInt(c.cfg.Query().Get("BinlogSlaveId"))),
		Flavor:   c.flavor(),
		Host:     c.cfg.Hostname(),
		Port:     uint16(conv.ToInt(c.cfg.Port())),
		User:     c.cfg.User.Username(),
		Password: pw,
	}

	c.syncer = replication.NewBinlogSyncer(&cfg)

	return nil
}

// Execute a SQL
func (c *Canal) Execute(cmd string, args ...interface{}) (rr *mysql.Result, err error) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	retryNum := 3
	for i := 0; i < retryNum; i++ {
		if c.conn == nil {
			pw, _ := c.cfg.User.Password()
			c.conn, err = client.Connect(c.cfg.Host, c.cfg.User.Username(), pw, "")
			if err != nil {
				return nil, errors.Wrap(err, "[binlogsync] checkBinlogRowFormat Execute")

			}
		}

		rr, err = c.conn.Execute(cmd, args...)
		if err != nil && !mysql.ErrorEqual(err, mysql.ErrBadConn) {
			return
		} else if mysql.ErrorEqual(err, mysql.ErrBadConn) {
			c.conn.Close()
			c.conn = nil
			continue
		} else {
			return
		}
	}
	return
}

func (c *Canal) SyncedPosition() mysql.Position {
	return c.master.Pos()
}

func (c *Canal) flavor() string {
	f := c.cfg.Query().Get("flavor")
	if f == "" {
		f = mysql.MySQLFlavor
	}
	switch f {
	case mysql.MySQLFlavor:
		return mysql.MySQLFlavor
	case mysql.MariaDBFlavor:
		return mysql.MariaDBFlavor
	}
	// todo remove panic
	panic(fmt.Sprintf("[binlogsync] Unknown flavor: %q", f))
}
