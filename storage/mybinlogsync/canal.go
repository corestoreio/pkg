package mybinlogsync

import (
	"fmt"
	"strings"
	"sync"

	"net/url"

	"github.com/corestoreio/csfw/util/conv"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"github.com/siddontang/go-mysql/schema"
	"github.com/siddontang/go/sync2"
)

var errCanalClosed = errors.New("canal was closed")

// Canal can sync your MySQL data. MySQL must use the binlog format ROW.
type Canal struct {
	m sync.Mutex

	cfg *url.URL

	master *masterInfo

	syncer *replication.BinlogSyncer

	rsLock     sync.Mutex
	rsHandlers []RowsEventHandler

	connLock sync.Mutex
	conn     *client.Conn

	wg sync.WaitGroup

	tableLock sync.RWMutex
	tables    map[string]*schema.Table

	closed sync2.AtomicBool
}

// export CS_DSN='mysql://root:PASSWORD@localhost:3306/DATABASE_NAME?BinlogSlaveId=100&BinlogStartPosition=0'
func NewCanal(cfg *url.URL) (*Canal, error) {
	c := new(Canal)
	c.cfg = cfg
	c.closed.Set(false)

	c.rsHandlers = make([]RowsEventHandler, 0, 4)
	c.tables = make(map[string]*schema.Table)

	var err error
	if c.master, err = loadMasterInfo(c); err != nil {
		return nil, errors.Trace(err)
	}
	c.updateBinlogStartPosition()

	if err = c.prepareSyncer(); err != nil {
		return nil, errors.Trace(err)
	}

	if err := c.checkBinlogRowFormat(); err != nil {
		return nil, errors.Trace(err)
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
			log.Errorf("canal start sync binlog err: %v", err)
		}
		return errors.Trace(err)
	}

	return nil
}

func (c *Canal) isClosed() bool {
	return c.closed.Get()
}

func (c *Canal) Close() {
	c.m.Lock()
	defer c.m.Unlock()

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

func (c *Canal) GetTable(db string, table string) (*schema.Table, error) {
	key := db + "." + table
	c.tableLock.RLock()
	t, ok := c.tables[key]
	c.tableLock.RUnlock()

	if ok {
		return t, nil
	}

	t, err := schema.NewTable(c, db, table)
	if err != nil {
		return nil, errors.Trace(err)
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
			return errors.Trace(err)
		} else {
			// MySQL has binlog row image from 5.6, so older will return empty
			rowImage, _ := res.GetString(0, 1)
			if rowImage != "" && !strings.EqualFold(rowImage, image) {
				return errors.Errorf("MySQL uses %s binlog row image, but we want %s", rowImage, image)
			}
		}
	}

	return nil
}

func (c *Canal) checkBinlogRowFormat() error {
	res, err := c.Execute(`SHOW GLOBAL VARIABLES LIKE "binlog_format";`)
	if err != nil {
		return errors.Trace(err)
	} else if f, _ := res.GetString(0, 1); f != "ROW" {
		return errors.Errorf("binlog must ROW format, but %s now", f)
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
				return nil, errors.Trace(err)
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
	panic(fmt.Sprintf("[mybinlogsync] Unknown flavor: %q", f))
}
