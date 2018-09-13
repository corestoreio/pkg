package myreplicator

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	gomysql "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	"github.com/siddontang/go-mysql/mysql"
)

var (
	testOutputLogs = flag.Bool("syncout", false, "output binlog event")
	runIntegration = flag.Bool("integration", false, "Enables MySQL/MariaDB integration tests, env var CS_DSN must be set")
)

func TestBinLogSyncer(t *testing.T) {
	if !*runIntegration {
		t.Skip("Skipping integration tests. You can enable them with via CLI option `-integration`")
	}

	tss := &testSyncerSuite{}
	defer tss.TearDownTest(t)

	tss.TestMysqlPositionSync(t)
	tss.TestMariaDBPositionSync(t)

	tss.TestMysqlGTIDSync(t)
	tss.TestMariaDBGTIDSync(t)
	//
	tss.TestMysqlSemiPositionSync(t)
	tss.TestMysqlBinlogCodec(t)
}

type testSyncerSuite struct {
	bls *BinlogSyncer
	con *dml.ConnPool

	wg sync.WaitGroup

	flavor string
}

func (t *testSyncerSuite) TearDownTest(tst *testing.T) {
	defer os.RemoveAll("./testdata/var")

	if t.bls != nil {
		t.bls.Close()
		t.bls = nil
	}

	if t.con != nil {
		assert.NoError(tst, t.con.Close())
		t.con = nil
	}
}

func (t *testSyncerSuite) testExecute(tst *testing.T, query string) {
	_, err := t.con.DB.Exec(query)
	assert.NoError(tst, err, "Query: %q", query)
}

func (t *testSyncerSuite) testSync(tst *testing.T, s *BinlogStreamer) {
	t.wg.Add(1)

	// https://dba.stackexchange.com/questions/165522/does-mariadb-support-native-json-column-data-type/176959#176959
	tblJSONColumnType := "JSON"
	if t.flavor == mysql.MariaDBFlavor {
		tblJSONColumnType = `TEXT CHECK (JSON_VALID(c1))`
	}

	tables, err := ddl.NewTables(
		ddl.WithDropTable(context.Background(), t.con.DB, "test_json", "test_json_v2", "test_geo"),
		ddl.WithCreateTable(context.Background(), t.con.DB,
			"test_json", `CREATE TABLE IF NOT EXISTS `+"`test_json`"+` (
			id BIGINT(64) UNSIGNED  NOT NULL AUTO_INCREMENT,
			c1 `+tblJSONColumnType+`,
			c2 DECIMAL(10, 0),
			PRIMARY KEY (id)
			) ENGINE=InnoDB`,

			"test_json_v2", `CREATE TABLE `+"`test_json_v2`"+` (
			id INT,
			c1 `+tblJSONColumnType+`,
			PRIMARY KEY (id)
			) ENGINE=InnoDB`,

			"test_geo", `CREATE TABLE `+"`test_geo` (g GEOMETRY)",
		),
	)
	assert.NoError(tst, err, "%+v", err)
	_ = tables

	expectedEvents := new(int32)

	go func() {
		defer t.wg.Done()

		if s == nil {
			return
		}

		const maxListeningTimeForEvents = 2 * time.Second

		for {
			ctx, cancel := context.WithTimeout(context.Background(), maxListeningTimeForEvents)
			e, err := s.GetEvent(ctx)
			cancel()

			if err != nil && errors.Cause(err) == context.DeadlineExceeded {
				return
			}

			assert.NoError(tst, err, "In goroutine while running GetEvent: %+v", err)

			if *testOutputLogs {
				e.Dump(os.Stdout)
				os.Stdout.Sync()
			}
			atomic.AddInt32(expectedEvents, 1)
		}
	}()

	//use mixed format
	t.testExecute(tst, "SET SESSION binlog_format = 'MIXED'")

	str := `DROP TABLE IF EXISTS test_replication`
	t.testExecute(tst, str)

	str = `CREATE TABLE IF NOT EXISTS test_replication (
			id BIGINT(64) UNSIGNED  NOT NULL AUTO_INCREMENT,
			str VARCHAR(256),
			f FLOAT,
			d DOUBLE,
			de DECIMAL(10,2),
			i INT,
			bi BIGINT,
			e enum ("e1", "e2"),
			b BIT(8),
			y YEAR,
			da DATE,
			ts TIMESTAMP,
			dt DATETIME,
			tm TIME,
			t TEXT,
			bb BLOB,
			se SET('a', 'b', 'c'),
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8`

	t.testExecute(tst, str)

	//use row format
	t.testExecute(tst, "SET SESSION binlog_format = 'ROW'")

	t.testExecute(tst, `INSERT INTO test_replication (str, f, i, e, b, y, da, ts, dt, tm, de, t, bb, se)
		VALUES ("3", -3.14, 10, "e1", 0b0011, 1985,
		"2012-05-07", "2012-05-07 14:01:01", "2012-05-07 14:01:01",
		"14:01:01", -45363.64, "abc", "12345", "a,b")`)

	id := 100

	if t.flavor == mysql.MySQLFlavor || true {
		t.testExecute(tst, "SET SESSION binlog_row_image = 'MINIMAL'")

		t.testExecute(tst, fmt.Sprintf(`INSERT INTO test_replication (id, str, f, i, bb, de) VALUES (%d, "4", -3.14, 100, "abc", -45635.64)`, id))
		t.testExecute(tst, fmt.Sprintf(`UPDATE test_replication SET f = -12.14, de = 555.34 WHERE id = %d`, id))
		t.testExecute(tst, fmt.Sprintf(`DELETE FROM test_replication WHERE id = %d`, id))
	}

	t.testExecute(tst, `INSERT INTO test_json (c2) VALUES (1)`)
	t.testExecute(tst, `INSERT INTO test_json (c1, c2) VALUES ('{"key1": "value1", "key2": "value2"}', 1)`)

	tbls := []string{
		// Refer: https://github.com/shyiko/mysql-binlog-connector-java/blob/c8e81c879710dc19941d952f9031b0a98f8b7c02/src/test/java/com/github/shyiko/mysql/binlog/event/deserialization/json/JsonBinaryValueIntegrationTest.java#L84
		// License: https://github.com/shyiko/mysql-binlog-connector-java#license
		`INSERT INTO test_json_v2 VALUES (0, NULL)`,
		`INSERT INTO test_json_v2 VALUES (1, '{\"a\": 2}')`,
		`INSERT INTO test_json_v2 VALUES (2, '[1,2]')`,
		`INSERT INTO test_json_v2 VALUES (3, '{\"a\":\"b\", \"c\":\"d\",\"ab\":\"abc\", \"bc\": [\"x\", \"y\"]}')`,
		`INSERT INTO test_json_v2 VALUES (4, '[\"here\", [\"I\", \"am\"], \"!!!\"]')`,
		`INSERT INTO test_json_v2 VALUES (5, '\"scalar string\"')`,
		`INSERT INTO test_json_v2 VALUES (6, 'true')`,
		`INSERT INTO test_json_v2 VALUES (7, 'false')`,
		`INSERT INTO test_json_v2 VALUES (8, 'null')`,
		`INSERT INTO test_json_v2 VALUES (9, '-1')`,
		`INSERT INTO test_json_v2 VALUES (11, '32767')`,
		`INSERT INTO test_json_v2 VALUES (12, '32768')`,
		`INSERT INTO test_json_v2 VALUES (13, '-32768')`,
		`INSERT INTO test_json_v2 VALUES (14, '-32769')`,
		`INSERT INTO test_json_v2 VALUES (15, '2147483647')`,
		`INSERT INTO test_json_v2 VALUES (16, '2147483648')`,
		`INSERT INTO test_json_v2 VALUES (17, '-2147483648')`,
		`INSERT INTO test_json_v2 VALUES (18, '-2147483649')`,
		`INSERT INTO test_json_v2 VALUES (19, '18446744073709551615')`,
		`INSERT INTO test_json_v2 VALUES (20, '18446744073709551616')`,
		`INSERT INTO test_json_v2 VALUES (21, '3.14')`,
		`INSERT INTO test_json_v2 VALUES (22, '{}')`,
		`INSERT INTO test_json_v2 VALUES (23, '[]')`,
		`INSERT INTO test_json_v2 VALUES (29, CAST('[]' AS CHAR CHARACTER SET 'ascii'))`,
		//`INSERT INTO test_json_v2 VALUES (100, CONCAT('{\"', REPEAT('€', 21845-8), '\":123}'))`,
	}

	if t.flavor == mysql.MySQLFlavor { // 10.3.9-MariaDB does not support yet to cast to JSON.
		tbls = append(tbls,
			`INSERT INTO test_json_v2 VALUES (10, CAST(CAST(1 AS UNSIGNED) AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (24, CAST(CAST('2015-01-15 23:24:25' AS DATETIME) AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (25, CAST(CAST('23:24:25' AS TIME) AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (125, CAST(CAST('23:24:25.12' AS TIME(3)) AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (225, CAST(CAST('23:24:25.0237' AS TIME(3)) AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (26, CAST(CAST('2015-01-15' AS DATE) AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (27, CAST(TIMESTAMP'2015-01-15 23:24:25' AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (127, CAST(TIMESTAMP'2015-01-15 23:24:25.12' AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (227, CAST(TIMESTAMP'2015-01-15 23:24:25.0237' AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (327, CAST(UNIX_TIMESTAMP('2015-01-15 23:24:25') AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (28, CAST(ST_GeomFromText('POINT(1 1)') AS JSON))`,
			// TODO: 30 and 31 are BIT type from JSON_TYPE, may support later.
			`INSERT INTO test_json_v2 VALUES (30, CAST(x'cafe' AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (31, CAST(x'cafebabe' AS JSON))`,
		)
	}

	for _, query := range tbls {
		t.testExecute(tst, query)
	}

	tbls = []string{
		`INSERT INTO test_geo VALUES (POINT(1, 1))`,
		`INSERT INTO test_geo VALUES (LINESTRING(POINT(0,0), POINT(1,1), POINT(2,2)))`,
		// TODO: add more geometry tests
	}

	for _, query := range tbls {
		t.testExecute(tst, query)
	}

	t.wg.Wait()

	switch t.flavor {
	case mysql.MariaDBFlavor:
		assert.True(tst, int32(145) < atomic.LoadInt32(expectedEvents), "expectedEvents count check min 145 events")
	case mysql.MySQLFlavor:
		assert.Exactly(tst, int32(9999), atomic.LoadInt32(expectedEvents), "TODO expectedEvents count check")
	}

}

func (t *testSyncerSuite) setupTest(tst *testing.T, testedFlavor string) bool {

	var err error
	if t.con != nil {
		assert.NoError(tst, t.con.Close())
	}

	t.con, err = dml.NewConnPool(dml.WithVerifyConnection(), dml.WithDSNfromEnv(dml.EnvDSN))
	if err != nil {
		tst.Fatal(err.Error())
	}

	_, err = t.con.DB.Exec("CREATE DATABASE IF NOT EXISTS test")
	assert.NoError(tst, err, "%+v", err)

	_, err = t.con.DB.Exec("USE test")
	assert.NoError(tst, err)

	if t.bls != nil {
		t.bls.Close()
	}

	dbVersion := ddl.NewVariables("version")
	_, err = t.con.WithQueryBuilder(dbVersion).Load(context.Background(), dbVersion)
	assert.NoError(tst, err)

	t.flavor = mysql.MySQLFlavor
	if dbVersion.Contains("version", "MariaDB") {
		t.flavor = mysql.MariaDBFlavor
	}

	if t.flavor != testedFlavor {
		t.flavor = ""
		tst.Logf("Skipping %q test because we run a different database server.", testedFlavor)
		return false
	}

	dsn, err := gomysql.ParseDSN(os.Getenv(dml.EnvDSN))
	assert.NoError(tst, err)

	host, port, _ := net.SplitHostPort(dsn.Addr)
	po, err := strconv.ParseUint(port, 10, 32)
	assert.NoError(tst, err)

	cfg := BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   t.flavor,
		Host:     host,
		Port:     uint16(po),
		User:     dsn.User,
		Password: dsn.Passwd,
	}

	t.bls = NewBinlogSyncer(&cfg)
	return true
}

func (t *testSyncerSuite) testPositionSync(tst *testing.T) {
	if t.flavor == "" {
		return // skipping because wrong database server
	}
	//get current master binlog file and position

	var ms ddl.MasterStatus
	_, err := t.con.WithQueryBuilder(ms).Load(context.Background(), &ms)
	assert.NoError(tst, err)

	s, err := t.bls.StartSync(ms)
	assert.NoError(tst, err, "%+v", err)

	// Test re-sync.
	time.Sleep(100 * time.Millisecond)
	t.bls.con.SetReadDeadline(time.Now().Add(time.Millisecond))
	time.Sleep(100 * time.Millisecond)

	t.testSync(tst, s)
}

func (t *testSyncerSuite) TestMysqlPositionSync(tst *testing.T) {
	t.setupTest(tst, mysql.MySQLFlavor)
	t.testPositionSync(tst)
}

func (t *testSyncerSuite) TestMysqlGTIDSync(tst *testing.T) {
	if !t.setupTest(tst, mysql.MySQLFlavor) {
		return
	}

	gtMode, ok, err := t.con.WithRawSQL("SELECT @@gtid_mode").LoadNullString(context.Background())
	assert.NoError(tst, err, "%+v", err)

	if !ok || gtMode.String != "ON" {
		tst.Skipf("GTID mode is not ON; got: %#v", gtMode)
	}

	srvUUID, ok, err := t.con.WithRawSQL("SHOW GLOBAL VARIABLES LIKE 'SERVER_UUID'").LoadNullString(context.Background())
	assert.NoError(tst, err, "%+v", err)

	var masterUuid uuid.UUID
	if srvUUID.String != "" && srvUUID.String != "NONE" {
		masterUuid, err = uuid.FromString(srvUUID.String)
		assert.NoError(tst, err, "%+v", err)
	}

	set, _ := mysql.ParseMysqlGTIDSet(fmt.Sprintf("%s:%d-%d", masterUuid.String(), 1, 2))

	s, err := t.bls.StartSyncGTID(set)
	assert.NoError(tst, err, "%+v", err)

	t.testSync(tst, s)
}

func (t *testSyncerSuite) TestMariaDBPositionSync(tst *testing.T) {
	if !t.setupTest(tst, mysql.MariaDBFlavor) {
		return
	}
	t.testPositionSync(tst)
}

func (t *testSyncerSuite) TestMariaDBGTIDSync(tst *testing.T) {
	if !t.setupTest(tst, mysql.MariaDBFlavor) {
		return
	}

	gtid_binlog_pos := ddl.NewVariables("gtid_binlog_pos")

	_, err := t.con.WithQueryBuilder(gtid_binlog_pos).Load(context.Background(), gtid_binlog_pos)
	assert.NoError(tst, err, "%+v", err)

	set, _ := mysql.ParseMariadbGTIDSet(gtid_binlog_pos.Data["gtid_binlog_pos"])

	s, err := t.bls.StartSyncGTID(set)
	assert.NoError(tst, err, "%+v", err)

	t.testSync(tst, s)
}

func (t *testSyncerSuite) TestMysqlSemiPositionSync(tst *testing.T) {
	if !t.setupTest(tst, mysql.MySQLFlavor) {
		return
	}

	t.bls.cfg.SemiSyncEnabled = true

	t.testPositionSync(tst)
}

func (t *testSyncerSuite) TestMysqlBinlogCodec(tst *testing.T) {
	if !t.setupTest(tst, mysql.MySQLFlavor) {
		return
	}

	t.testExecute(tst, "RESET MASTER")

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	go func() {
		defer wg.Done()

		t.testSync(tst, nil)

		t.testExecute(tst, "FLUSH LOGS")

		t.testSync(tst, nil)
	}()

	assert.NoError(tst, os.RemoveAll("./testdata/var"))

	err := t.bls.StartBackup("./testdata/var", ddl.MasterStatus{Position: uint(0)}, 2*time.Second)
	assert.NoError(tst, err, "%+v", err)

	p := NewBinlogParser()

	f := func(e *BinlogEvent) error {
		if *testOutputLogs {
			e.Dump(os.Stdout)
			os.Stdout.Sync()
		}
		return nil
	}

	err = p.ParseFile("./var/mysql.000001", 0, f)
	assert.NoError(tst, err, "%+v", err)

	err = p.ParseFile("./var/mysql.000002", 0, f)
	assert.NoError(tst, err, "%+v", err)
}
