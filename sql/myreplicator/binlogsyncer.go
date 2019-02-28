package myreplicator

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/satori/go.uuid"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
)

var (
	errSyncRunning = errors.AlreadyExists.Newf("[myreplicator] Sync is already running. You should close it first.")
)

// BinlogSyncerConfig is the configuration for BinlogSyncer.
type BinlogSyncerConfig struct {
	// ServerID is the unique ID in cluster.
	ServerID uint32
	// Flavor is "mysql" or "mariadb", if not set, use "mysql" default.
	Flavor string

	// Host is for MySQL server host.
	Host string
	// Port is for MySQL server port.
	Port uint16
	// User is for MySQL user.
	User string
	// Password is for MySQL password.
	Password string

	// Localhost is local hostname if register salve.
	// If not set, use os.Hostname() instead.
	Localhost string

	// Charset is for MySQL client character set
	Charset string

	// SemiSyncEnabled enables semi-sync or not.
	SemiSyncEnabled bool

	// RawModeEnabled is for not parsing binlog event.
	RawModeEnabled bool

	Log log.Logger
	// TLSConfig if not nil, use the provided tls.Config to connect to the
	// database using TLS/SSL.
	TLSConfig *tls.Config

	// Use replication.Time structure for timestamp and datetime.
	// We will use Local location for timestamp and UTC location for datatime.
	ParseTime bool

	// If ParseTime is false, convert TIMESTAMP into this specified timezone. If
	// ParseTime is true, this option will have no effect and TIMESTAMP data will
	// be parsed into the local timezone and a full time.Time struct will be
	// returned.
	//
	// Note that MySQL TIMESTAMP columns are offset from the machine local
	// timezone while DATETIME columns are offset from UTC. This is consistent
	// with documented MySQL behaviour as it return TIMESTAMP in local timezone
	// and DATETIME in UTC.
	//
	// Setting this to UTC effectively equalizes the TIMESTAMP and DATETIME time
	// strings obtained from MySQL.
	TimestampStringLocation *time.Location

	// RecvBufferSize sets the size in bytes of the operating system's receive buffer associated with the connection.
	RecvBufferSize int

	// master heartbeat period
	HeartbeatPeriod time.Duration

	// read timeout
	ReadTimeout time.Duration

	// maximum number of attempts to re-establish a broken connection
	MaxReconnectAttempts int

	// Only works when MySQL/MariaDB variable binlog_checksum=CRC32.
	// For MySQL, binlog_checksum was introduced since 5.6.2, but CRC32 was set as default value since 5.6.6 .
	// https://dev.mysql.com/doc/refman/5.6/en/replication-options-binary-log.html#option_mysqld_binlog-checksum
	// For MariaDB, binlog_checksum was introduced since MariaDB 5.3, but CRC32 was set as default value since MariaDB 10.2.1 .
	// https://mariadb.com/kb/en/library/replication-and-binary-log-server-system-variables/#binlog_checksum
	VerifyChecksum bool
}

// BinlogSyncer syncs binlog event from server.
type BinlogSyncer struct {
	m sync.RWMutex

	cfg *BinlogSyncerConfig

	con *client.Conn

	wg sync.WaitGroup

	parser *BinlogParser

	nextPos ddl.MasterStatus

	gset mysql.GTIDSet

	running bool

	ctx    context.Context
	cancel context.CancelFunc

	lastConnectionID uint32

	retryCount int
}

// NewBinlogSyncer creates the BinlogSyncer with cfg.
func NewBinlogSyncer(cfg *BinlogSyncerConfig) *BinlogSyncer {
	if cfg.Log == nil {
		cfg.Log = log.BlackHole{}
	}
	if cfg.Log.IsDebug() {
		cfg2 := new(BinlogSyncerConfig)
		*cfg2 = *cfg
		cfg2.Password = "********"
		cfg.Log.Debug("NewBinlogSyncer.BinlogSyncerConfig", log.Object("config", cfg2))
	}

	b := &BinlogSyncer{
		cfg:    cfg,
		parser: NewBinlogParser(),
	}
	b.parser.SetRawMode(b.cfg.RawModeEnabled)
	b.parser.SetParseTime(b.cfg.ParseTime)
	b.parser.SetTimestampStringLocation(b.cfg.TimestampStringLocation)
	b.parser.SetVerifyChecksum(b.cfg.VerifyChecksum)
	b.running = false

	b.ctx, b.cancel = context.WithCancel(context.Background())

	return b
}

// Close closes the BinlogSyncer.
func (b *BinlogSyncer) Close() error {
	b.m.Lock()
	defer b.m.Unlock()
	return b.close()
}

func (b *BinlogSyncer) close() (err error) {
	if b.cfg.Log.IsDebug() {
		defer log.WhenDone(b.cfg.Log).Debug("BinlogSyncer.close", log.Err(err))
	}

	if b.isClosed() {
		return nil
	}

	b.running = false
	b.cancel()

	if b.con != nil {
		if err = b.con.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); err != nil {
			err = errors.WithStack(err)
			return
		}
	}

	b.wg.Wait()

	if b.con != nil {
		if err = b.con.Close(); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	return
}

func (b *BinlogSyncer) isClosed() bool {
	select {
	case <-b.ctx.Done():
		return true
	default:
		return false
	}
}

func (b *BinlogSyncer) registerSlave() error {
	if b.con != nil {
		if err := b.con.Close(); err != nil && b.cfg.Log.IsInfo() {
			b.cfg.Log.Info("BinlogSyncer.registerSlave.con.close.error", log.Err(err))
		}
	}

	if b.cfg.Log.IsDebug() {
		b.cfg.Log.Debug("BinlogSyncer.registerSlave.new", log.String("slave_host", b.cfg.Host), log.Uint("slave_port", uint(b.cfg.Port)))
	}
	var err error
	b.con, err = client.Connect(fmt.Sprintf("%s:%d", b.cfg.Host, b.cfg.Port), b.cfg.User, b.cfg.Password, "", func(c *client.Conn) {
		c.SetTLSConfig(b.cfg.TLSConfig)
	})
	if err != nil {
		return errors.WithStack(err)
	}

	if len(b.cfg.Charset) != 0 {
		if err := b.con.SetCharset(b.cfg.Charset); err != nil {
			return errors.WithStack(err)
		}
	}

	// set read timeout
	if b.cfg.ReadTimeout > 0 {
		if err := b.con.SetReadDeadline(time.Now().Add(b.cfg.ReadTimeout)); err != nil {
			return errors.WithStack(err)
		}
	}

	if b.cfg.RecvBufferSize > 0 {
		if tcp, ok := b.con.Conn.Conn.(*net.TCPConn); ok {
			if err := tcp.SetReadBuffer(b.cfg.RecvBufferSize); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	// kill last connection id
	if b.lastConnectionID > 0 {
		cmd := fmt.Sprintf("KILL %d", b.lastConnectionID)
		if _, err := b.con.Execute(cmd); err != nil {
			if b.cfg.Log.IsInfo() {
				b.cfg.Log.Info("BinlogSyncer.registerSlave.connection.kill.error", log.Err(err), log.Uint("last_connection_id", uint(b.lastConnectionID)))
			}
			// Unknown thread id
			if code := mysql.ErrorCode(err.Error()); code != mysql.ER_NO_SUCH_THREAD {
				return errors.WithStack(err)
			}
		}
		if b.cfg.Log.IsDebug() {
			b.cfg.Log.Debug("BinlogSyncer.registerSlave.connection.killed", log.Uint("last_connection_id", uint(b.lastConnectionID)))
		}
	}

	// save last last connection id for kill
	b.lastConnectionID = b.con.GetConnectionID()

	//for mysql 5.6+, binlog has a crc32 checksum
	//before mysql 5.6, this will not work, don't matter.:-)
	if r, err := b.con.Execute("SHOW GLOBAL VARIABLES LIKE 'BINLOG_CHECKSUM'"); err != nil {
		return errors.WithStack(err)
	} else {
		s, _ := r.GetString(0, 1)
		if s != "" {
			// maybe CRC32 or NONE

			// mysqlbinlog.cc use NONE, see its below comments:
			// Make a notice to the server that this client
			// is checksum-aware. It does not need the first fake Rotate
			// necessary checksummed.
			// That preference is specified below.

			if _, err = b.con.Execute(`SET @master_binlog_checksum='NONE'`); err != nil {
				return errors.WithStack(err)
			}

			// if _, err = b.con.Execute(`SET @master_binlog_checksum=@@global.binlog_checksum`); err != nil {
			// 	return errors.WithStack(err)
			// }

		}
	}

	if b.cfg.Flavor == mysql.MariaDBFlavor {
		// Refer https://github.com/alibaba/canal/wiki/BinlogChange(MariaDB5&10)
		// Tell the server that we understand GTIDs by setting our slave capability
		// to MARIA_SLAVE_CAPABILITY_GTID = 4 (MariaDB >= 10.0.1).
		if _, err := b.con.Execute("SET @mariadb_slave_capability=4"); err != nil {
			return errors.Fatal.Newf("[myreplicator] failed to set @mariadb_slave_capability=4: %v", err)
		}
	}

	if b.cfg.HeartbeatPeriod > 0 {
		_, err = b.con.Execute(fmt.Sprintf("SET @master_heartbeat_period=%d;", b.cfg.HeartbeatPeriod))
		if err != nil {
			if b.cfg.Log.IsInfo() {
				b.cfg.Log.Info("BinlogSyncer.registerSlave.heartbeatPeriod.error", log.Err(err), log.Duration("last_connection_id", b.cfg.HeartbeatPeriod))
			}
			return errors.WithStack(err)
		}
	}

	if err = b.writeRegisterSlaveCommand(); err != nil {
		return errors.WithStack(err)
	}

	if _, err = b.con.ReadOKPacket(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (b *BinlogSyncer) enableSemiSync() error {
	if !b.cfg.SemiSyncEnabled {
		return nil
	}

	if r, err := b.con.Execute("SHOW VARIABLES LIKE 'rpl_semi_sync_master_enabled';"); err != nil {
		return errors.WithStack(err)
	} else {
		s, _ := r.GetString(0, 1)
		if s != "ON" {
			if b.cfg.Log.IsDebug() {
				b.cfg.Log.Debug("BinlogSyncer.enableSemiSync.failed", log.String("reason", "master does not support semi synchronous myreplicator, use non-semi-sync"))
			}
			b.cfg.SemiSyncEnabled = false
			return nil
		}
	}

	_, err := b.con.Execute(`SET @rpl_semi_sync_slave = 1;`)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (b *BinlogSyncer) prepare() error {
	if b.isClosed() {
		return errors.AlreadyClosed.Newf("[myreplicator] Syncer already closed")
	}

	if err := b.registerSlave(); err != nil {
		return errors.WithStack(err)
	}

	if err := b.enableSemiSync(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (b *BinlogSyncer) startDumpStream() *BinlogStreamer {
	b.running = true

	s := newBinlogStreamer(b.cfg.Log)

	b.wg.Add(1)
	go b.onStream(s)
	return s
}

// StartSync starts syncing from the `pos` position.
func (b *BinlogSyncer) StartSync(pos ddl.MasterStatus) (*BinlogStreamer, error) {
	if b.cfg.Log.IsDebug() {
		b.cfg.Log.Debug("BinlogSyncer.StartSync", log.Stringer("position", pos))
	}

	b.m.Lock()
	defer b.m.Unlock()

	if b.running {
		return nil, errors.WithStack(errSyncRunning)
	}

	if err := b.prepareSyncPos(pos); err != nil {
		return nil, errors.WithStack(err)
	}

	return b.startDumpStream(), nil
}

// StartSyncGTID starts syncing from the `gset` mysql.GTIDSet.
func (b *BinlogSyncer) StartSyncGTID(gset mysql.GTIDSet) (*BinlogStreamer, error) {
	if b.cfg.Log.IsDebug() {
		b.cfg.Log.Debug("BinlogSyncer.StartSyncGTID", log.Stringer("gtid", gset))
	}

	b.m.Lock()
	defer b.m.Unlock()

	if b.running {
		return nil, errors.WithStack(errSyncRunning)
	}

	if err := b.prepare(); err != nil {
		return nil, errors.WithStack(err)
	}

	var err error
	switch b.cfg.Flavor {
	case mysql.MariaDBFlavor:
		err = b.writeBinlogDumpMariadbGTIDCommand(gset)
	default:
		// default use MySQL
		err = b.writeBinlogDumpMysqlGTIDCommand(gset)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return b.startDumpStream(), nil
}

func (b *BinlogSyncer) writeBinglogDumpCommand(p ddl.MasterStatus) error {
	b.con.ResetSequence()

	data := make([]byte, 4+1+4+2+4+len(p.File))

	pos := 4
	data[pos] = mysql.COM_BINLOG_DUMP
	pos++

	binary.LittleEndian.PutUint32(data[pos:], uint32(p.Position))
	pos += 4

	binary.LittleEndian.PutUint16(data[pos:], BINLOG_DUMP_NEVER_STOP)
	pos += 2

	binary.LittleEndian.PutUint32(data[pos:], b.cfg.ServerID)
	pos += 4

	copy(data[pos:], p.File)

	return b.con.WritePacket(data)
}

func (b *BinlogSyncer) writeBinlogDumpMysqlGTIDCommand(gset mysql.GTIDSet) error {
	p := ddl.MasterStatus{Position: 4}
	gtidData := gset.Encode()

	b.con.ResetSequence()

	data := make([]byte, 4+1+2+4+4+len(p.File)+8+4+len(gtidData))
	pos := 4
	data[pos] = mysql.COM_BINLOG_DUMP_GTID
	pos++

	binary.LittleEndian.PutUint16(data[pos:], 0)
	pos += 2

	binary.LittleEndian.PutUint32(data[pos:], b.cfg.ServerID)
	pos += 4

	binary.LittleEndian.PutUint32(data[pos:], uint32(len(p.File)))
	pos += 4

	n := copy(data[pos:], p.File)
	pos += n

	binary.LittleEndian.PutUint64(data[pos:], uint64(p.Position))
	pos += 8

	binary.LittleEndian.PutUint32(data[pos:], uint32(len(gtidData)))
	pos += 4
	n = copy(data[pos:], gtidData)
	pos += n

	data = data[0:pos]

	return b.con.WritePacket(data)
}

func (b *BinlogSyncer) writeBinlogDumpMariadbGTIDCommand(gset mysql.GTIDSet) error {
	// Copy from vitess

	startPos := gset.String()

	// Set the slave_connect_state variable before issuing COM_BINLOG_DUMP to
	// provide the start position in GTID form.
	query := fmt.Sprintf("SET @slave_connect_state='%s'", startPos)

	if _, err := b.con.Execute(query); err != nil {
		return errors.Fatal.Newf("[myreplicator] failed to set @slave_connect_state='%s': %v", startPos, err)
	}

	// Real slaves set this upon connecting if their gtid_strict_mode option was
	// enabled. We always use gtid_strict_mode because we need it to make our
	// internal GTID comparisons safe.
	if _, err := b.con.Execute("SET @slave_gtid_strict_mode=1"); err != nil {
		return errors.Fatal.Newf("[myreplicator] failed to set @slave_gtid_strict_mode=1: %v", err)
	}

	// Since we use @slave_connect_state, the file and position here are ignored.
	return b.writeBinglogDumpCommand(ddl.MasterStatus{})
}

// localHostname returns the hostname that register slave would register as.
func (b *BinlogSyncer) localHostname() string {
	if len(b.cfg.Localhost) == 0 {
		h, _ := os.Hostname()
		return h
	}
	return b.cfg.Localhost
}

// if this command fails then the permission `REPLICATION SLAVE` is missing.
func (b *BinlogSyncer) writeRegisterSlaveCommand() error {
	b.con.ResetSequence()

	hostname := b.localHostname()

	// This should be the name of slave host not the host we are connecting to.
	data := make([]byte, 4+1+4+1+len(hostname)+1+len(b.cfg.User)+1+len(b.cfg.Password)+2+4+4)
	pos := 4

	data[pos] = mysql.COM_REGISTER_SLAVE
	pos++

	binary.LittleEndian.PutUint32(data[pos:], b.cfg.ServerID)
	pos += 4

	// This should be the name of slave hostname not the host we are connecting to.
	data[pos] = uint8(len(hostname))
	pos++
	n := copy(data[pos:], hostname)
	pos += n

	data[pos] = uint8(len(b.cfg.User))
	pos++
	n = copy(data[pos:], b.cfg.User)
	pos += n

	data[pos] = uint8(len(b.cfg.Password))
	pos++
	n = copy(data[pos:], b.cfg.Password)
	pos += n

	binary.LittleEndian.PutUint16(data[pos:], b.cfg.Port)
	pos += 2

	// myreplicator rank, not used
	binary.LittleEndian.PutUint32(data[pos:], 0)
	pos += 4

	// master ID, 0 is OK
	binary.LittleEndian.PutUint32(data[pos:], 0)

	return b.con.WritePacket(data)
}

func (b *BinlogSyncer) replySemiSyncACK(p ddl.MasterStatus) error {
	b.con.ResetSequence()

	data := make([]byte, 4+1+8+len(p.File))
	pos := 4
	// semi sync indicator
	data[pos] = SemiSyncIndicator
	pos++

	binary.LittleEndian.PutUint64(data[pos:], uint64(p.Position))
	pos += 8

	copy(data[pos:], p.File)

	err := b.con.WritePacket(data)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (b *BinlogSyncer) retrySync() error {
	b.m.Lock()
	defer b.m.Unlock()

	if b.cfg.Log.IsDebug() {
		b.cfg.Log.Debug("BinlogSyncer.retrySync", log.Stringer("position", b.nextPos))
	}

	b.parser.Reset()

	if b.gset != nil {
		if err := b.prepareSyncGTID(b.gset); err != nil {
			return errors.WithStack(err)
		}
		if b.cfg.Log.IsDebug() {
			b.cfg.Log.Debug("BinlogSyncer.retrySync.start", log.Stringer("pos", b.gset))
		}
	} else {
		if err := b.prepareSyncPos(b.nextPos); err != nil {
			return errors.WithStack(err)
		}
		if b.cfg.Log.IsDebug() {
			b.cfg.Log.Debug("BinlogSyncer.retrySync.start", log.Stringer("pos", b.nextPos))
		}
	}

	return nil
}

func (b *BinlogSyncer) prepareSyncPos(pos ddl.MasterStatus) error {
	// always start from position 4
	if pos.Position < 4 {
		pos.Position = 4
	}

	if err := b.prepare(); err != nil {
		return errors.WithStack(err)
	}

	if err := b.writeBinglogDumpCommand(pos); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (b *BinlogSyncer) prepareSyncGTID(gset mysql.GTIDSet) (err error) {
	if err = b.prepare(); err != nil {
		return errors.WithStack(err)
	}

	switch b.cfg.Flavor {
	case mysql.MariaDBFlavor:
		err = b.writeBinlogDumpMariadbGTIDCommand(gset)
	default:
		// default use MySQL
		err = b.writeBinlogDumpMysqlGTIDCommand(gset)
	}
	return
}

// onStream runs in a goroutine
func (b *BinlogSyncer) onStream(s *BinlogStreamer) {
	defer func() {
		if e := recover(); e != nil {
			s.closeWithError(errors.Fatal.Newf("[myreplicator] onStream.Recovered with error: %v", e))
		}
		b.wg.Done()
	}()

	for {
		data, err := b.con.ReadPacket()
		if err != nil {
			if b.cfg.Log.IsInfo() {
				b.cfg.Log.Info("BinlogSyncer.onStream.connection.readpacket", log.Err(err))
			}

			// we meet connection error, should re-connect again with
			// last nextPos we got.
			if b.nextPos.File == "" {
				// we can't get the correct position, close.
				s.closeWithError(err)
				return
			}

			// TODO: add a max retry count.
			for {
				select {
				case <-b.ctx.Done():
					if b.cfg.Log.IsDebug() {
						b.cfg.Log.Debug("BinlogSyncer.onStream.connection.context.done")
					}
					return
				case <-time.After(time.Second):
					b.retryCount++
					if err = b.retrySync(); err != nil {
						if b.cfg.MaxReconnectAttempts > 0 && b.retryCount >= b.cfg.MaxReconnectAttempts {
							if b.cfg.Log.IsInfo() {
								b.cfg.Log.Info("BinlogSyncer.onStream.retrySync.waitAsecond", log.Err(err), log.Int("max_reconnect_attempts", b.cfg.MaxReconnectAttempts))
							}
							s.closeWithError(err)
							return
						}
						if b.cfg.Log.IsInfo() {
							b.cfg.Log.Info("BinlogSyncer.onStream.retrySync.err", log.Err(err))
						}
						continue
					}
				}
				break
			}

			// we connect the server and begin to re-sync again.
			continue
		}

		// set read timeout
		if b.cfg.ReadTimeout > 0 {
			if err := b.con.SetReadDeadline(time.Now().Add(b.cfg.ReadTimeout)); err != nil {
				s.closeWithError(err)
				return
			}
		}

		// Reset retry count on successful packet receieve
		b.retryCount = 0

		switch data[0] {
		case mysql.OK_HEADER:
			if err = b.parseEvent(s, data); err != nil {
				s.closeWithError(err)
				return
			}
		case mysql.ERR_HEADER:
			err = b.con.HandleErrorPacket(data)
			s.closeWithError(err)
			return
		case mysql.EOF_HEADER:
			// Refer http://dev.mysql.com/doc/internals/en/packet-EOF_Packet.html
			// In the MySQL client/server protocol, EOF and OK packets serve the same purpose.
			// Some users told me that they received EOF packet here, but I don't know why.
			// So we only log a message and retry ReadPacket.
			if b.cfg.Log.IsInfo() {
				b.cfg.Log.Info("BinlogSyncer.onStream.eof_header", log.String("info", "receive EOF packet, retry ReadPacket"),
					log.Err(err))
			}
			continue
		default:
			s.closeWithError(fmt.Errorf("invalid stream header %c", data[0]))
			return
		}
	}
}

func (b *BinlogSyncer) parseEvent(s *BinlogStreamer, data []byte) error {
	// skip OK byte, 0x00
	data = data[1:]

	needACK := false
	if b.cfg.SemiSyncEnabled && (data[0] == SemiSyncIndicator) {
		needACK = data[1] == 0x01
		// skip semi sync header
		data = data[2:]
	}

	e, err := b.parser.Parse(data)
	if b.cfg.Log.IsDebug() {
		b.cfg.Log.Debug("BinlogSyncer.onStream.parseEvent",
			log.Int("data_len", len(data)),
			log.String("type", fmt.Sprintf("%T", e.Event)),
			log.Err(err),
		)
	}
	if err != nil {
		return errors.WithStack(err)
	}

	if e.Header.LogPos > 0 {
		// Some events like FormatDescriptionEvent return 0, ignore.
		b.nextPos.Position = uint(e.Header.LogPos)
	}
	switch event := e.Event.(type) {
	case *RotateEvent:
		b.nextPos.File = string(event.NextLogName)
		b.nextPos.Position = uint(event.Position)
		if b.cfg.Log.IsInfo() {
			b.cfg.Log.Info("BinlogSyncer.parseEvent.rotate_to", log.Stringer("position", b.nextPos))
		}
	case *GTIDEvent:
		if b.gset == nil {
			break
		}
		u, _ := uuid.FromBytes(event.SID)
		err := b.gset.Update(fmt.Sprintf("%s:%d", u.String(), event.GNO))
		if err != nil {
			return errors.WithStack(err)
		}
	case *MariadbGTIDEvent:
		if b.gset == nil {
			break
		}
		GTID := event.GTID
		err := b.gset.Update(fmt.Sprintf("%d-%d-%d", GTID.DomainID, GTID.ServerID, GTID.SequenceNumber))
		if err != nil {
			return errors.WithStack(err)
		}
	case *XIDEvent:
		event.GSet = b.getGtidSet()
	case *QueryEvent:
		event.GSet = b.getGtidSet()
	}

	needStop := false
	select {
	case s.bleChan <- e:
	case <-b.ctx.Done():
		needStop = true
	}

	if needACK {
		err := b.replySemiSyncACK(b.nextPos)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if needStop {
		return errors.InProgress.Newf("[myreplicator] sync is been closing...")
	}

	return nil
}

func (b *BinlogSyncer) getGtidSet() mysql.GTIDSet {
	if b.gset == nil {
		return nil
	}
	return b.gset.Clone()
}

// LastConnectionID returns last connectionID.
func (b *BinlogSyncer) LastConnectionID() uint32 {
	return b.lastConnectionID
}
