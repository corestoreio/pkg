// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dml

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
	_ "github.com/go-sql-driver/mysql"
)

///////////////////////////////////////////////////////////////////////////////
// TEST HELPERS
///////////////////////////////////////////////////////////////////////////////

func createRealSession(t testing.TB) *ConnPool {
	dsn := os.Getenv("CS_DSN")
	if dsn == "" {
		t.Skip("Environment variable CS_DSN not found. Skipping ...")
	}
	cxn, err := NewConnPool(
		WithDSN(dsn),
	)
	if err != nil {
		t.Fatal(err)
	}
	return cxn
}

func createRealSessionWithFixtures(t testing.TB, c *installFixturesConfig) *ConnPool {
	sess := createRealSession(t)
	installFixtures(t, sess.DB, c)
	return sess
}

// testCloser for usage in conjunction with defer.
// 		defer testCloser(t,db)
// Cannot use dmltest.Close because of circular dependency.
func testCloser(t testing.TB, c ioCloser) {
	t.Helper()
	if err := c.Close(); err != nil {
		t.Errorf("%+v", err)
	}
}

var _ ColumnMapper = (*dmlPerson)(nil)
var _ LastInsertIDAssigner = (*dmlPerson)(nil)
var _ ColumnMapper = (*dmlPersons)(nil)

type dmlPerson struct {
	ID    uint64
	Name  string
	Email null.String
	Key   null.String
}

func (p *dmlPerson) AssignLastInsertID(id int64) {
	p.ID = uint64(id)
}

// RowScan loads a single row from a SELECT statement returning only one row
func (p *dmlPerson) MapColumns(cm *ColumnMap) error {
	if cm.Mode() == ColumnMapEntityReadAll {
		return cm.Uint64(&p.ID).String(&p.Name).NullString(&p.Email).NullString(&p.Key).Err()
	}
	for cm.Next() {
		c := cm.Column()
		switch c {
		case "id":
			cm.Uint64(&p.ID)
		case "name":
			cm.String(&p.Name)
		case "email":
			cm.NullString(&p.Email)
		case "key":
			cm.NullString(&p.Key)
		case "store_id", "created_at", "total_income", "avg_income":
			// noop don't trigger the default case
		default:
			return errors.NotFound.Newf("[dml_test] dmlPerson Column %q not found", c)
		}
	}
	return errors.WithStack(cm.Err())
}

type dmlPersons struct {
	Data []*dmlPerson
}

// MapColumns gets called in the `for rows.Next()` loop each time in case of IsNew
func (ps *dmlPersons) MapColumns(cm *ColumnMap) error {
	switch m := cm.Mode(); m {
	case ColumnMapEntityReadAll, ColumnMapEntityReadSet:
		for _, p := range ps.Data {
			if err := p.MapColumns(cm); err != nil {
				return errors.WithStack(err)
			}
		}
	case ColumnMapScan:
		// case for scanning when loading certain rows, hence we write data from
		// the DB into the struct in each for-loop.
		if cm.Count == 0 {
			ps.Data = ps.Data[:0]
		}
		p := new(dmlPerson)
		if err := p.MapColumns(cm); err != nil {
			return errors.WithStack(err)
		}
		ps.Data = append(ps.Data, p)
	case ColumnMapCollectionReadSet: // See Test in select_test.go:TestSelect_SetRecord
		// SELECT, DELETE or UPDATE or INSERT with n columns
		// TODO in some INSERT statements this slice building code might not be needed.
		for cm.Next() {
			switch c := cm.Column(); c {
			case "id":
				cm.Uint64s(ps.IDs()...)
			case "name":
				cm.Strings(ps.Names()...)
			case "email":
				cm.NullStrings(ps.Emails()...)
			default:
				return errors.NotFound.Newf("[dml_test] dmlPerson Column %q not found", c)
			}
		}
	default:
		return errors.NotSupported.Newf("[dml] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

func (ps *dmlPersons) IDs(ret ...uint64) []uint64 {
	if ret == nil {
		ret = make([]uint64, 0, len(ps.Data))
	}
	for _, p := range ps.Data {
		ret = append(ret, p.ID)
	}
	return ret
}

func (ps *dmlPersons) Names(ret ...string) []string {
	if ret == nil {
		ret = make([]string, 0, len(ps.Data))
	}
	for _, p := range ps.Data {
		ret = append(ret, p.Name)
	}
	return ret
}

func (ps *dmlPersons) Emails(ret ...null.String) []null.String {
	if ret == nil {
		ret = make([]null.String, 0, len(ps.Data))
	}
	for _, p := range ps.Data {
		ret = append(ret, p.Email)
	}
	return ret
}

//func (ps *dmlPersons) AssignLastInsertID(uint64) error {
//	// todo iterate and assign to the last item in the slice and assign
//	// decremented IDs to the previous items in the slice.
//	return nil
//}
//

var _ ColumnMapper = (*nullTypedRecord)(nil)

type nullTypedRecord struct {
	ID         int64
	StringVal  null.String
	Int64Val   null.Int64
	Float64Val null.Float64
	TimeVal    null.Time
	BoolVal    null.Bool
	DecimalVal null.Decimal
}

func (p *nullTypedRecord) MapColumns(cm *ColumnMap) error {
	if cm.Mode() == ColumnMapEntityReadAll {
		return cm.Int64(&p.ID).NullString(&p.StringVal).NullInt64(&p.Int64Val).NullFloat64(&p.Float64Val).NullTime(&p.TimeVal).NullBool(&p.BoolVal).Decimal(&p.DecimalVal).Err()
	}
	for cm.Next() {
		c := cm.Column()
		switch c {
		case "id":
			cm.Int64(&p.ID)
		case "string_val":
			cm.NullString(&p.StringVal)
		case "int64_val":
			cm.NullInt64(&p.Int64Val)
		case "float64_val":
			cm.NullFloat64(&p.Float64Val)
		case "time_val":
			cm.NullTime(&p.TimeVal)
		case "bool_val":
			cm.NullBool(&p.BoolVal)
		case "decimal_val":
			cm.Decimal(&p.DecimalVal)
		default:
			return errors.NotFound.Newf("[dml_test] Column %q not found", c)
		}
	}
	return cm.Err()
}
func newNullTypedRecordWithData() *nullTypedRecord {
	return &nullTypedRecord{
		ID:         2,
		StringVal:  null.String{String: "wow", Valid: true},
		Int64Val:   null.Int64{Int64: 42, Valid: true},
		Float64Val: null.Float64{Float64: 1.618, Valid: true},
		TimeVal:    null.Time{Time: time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC), Valid: true},
		BoolVal:    null.Bool{Bool: true, Valid: true},
		DecimalVal: null.Decimal{Precision: 12345, Scale: 3, Valid: true},
	}
}

type installFixturesConfig struct {
	AddPeopleWithMaxUint64 bool
}

func installFixtures(t testing.TB, db *sql.DB, c *installFixturesConfig) {
	createPeopleTable := fmt.Sprintf(`
		CREATE TABLE dml_people (
			id bigint(8) unsigned NOT NULL auto_increment PRIMARY KEY,
			name varchar(255) NOT NULL,
			email varchar(255),
			%s varchar(255),
			store_id smallint(5) unsigned DEFAULT 0 COMMENT 'Store Id',
			created_at timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT 'Created At',
			total_income decimal(12,4) NOT NULL DEFAULT 0.0000 COMMENT 'Used as float64',
			avg_income decimal(12,5) COMMENT 'Used as Decimal'
		)
	`, "`key`")

	createNullTypesTable := `
		CREATE TABLE dml_null_types (
			id int(11) NOT NULL auto_increment PRIMARY KEY,
			string_val varchar(255) NULL,
			int64_val int(11) NULL,
			float64_val float NULL,
			time_val datetime NULL,
			bool_val bool NULL,
			decimal_val decimal(5,3) NULL
		)
	`
	// see also test case "LoadUint64 max Uint64 found"
	sqlToRun := []string{
		"DROP TABLE IF EXISTS `dml_people`",
		createPeopleTable,
		"INSERT INTO dml_people (name,email,avg_income) VALUES ('Sir George', 'SirGeorge@GoIsland.com',333.66677)",
		"INSERT INTO dml_people (name,email) VALUES ('Dmitri', 'zavorotni@jadius.com')",

		"DROP TABLE IF EXISTS `dml_null_types`",
		createNullTypesTable,
	}
	if c != nil && c.AddPeopleWithMaxUint64 {
		sqlToRun = append(sqlToRun, "INSERT INTO `dml_people` (id,name,email) VALUES (18446744073700551613,'Cyrill', 'firstname@lastname.fm')")
	}

	for _, sqlStr := range sqlToRun {
		_, err := db.Exec(sqlStr)
		assert.NoError(t, err, "With SQL statement: %q", sqlStr)
	}
}

var _ QueryExecPreparer = (*dbMock)(nil)
var _ Execer = (*dbMock)(nil)

type dbMock struct {
	error
	prepareFn func(query string) (*sql.Stmt, error)
}

func (pm dbMock) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return pm.prepareFn(query)
}

func (pm dbMock) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return nil, nil
}

func (pm dbMock) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return new(sql.Row)
}

func (pm dbMock) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return nil, nil
}

// compareToSQL compares a SQL object with a placeholder string and an optional
// interpolated string. This function also exists in file dml_public_test.go to
// avoid import cycles when using a single package dedicated for testing.
func compareToSQL(
	t testing.TB, qb QueryBuilder, wantErrKind errors.Kind,
	wantSQLPlaceholders, wantSQLInterpolated string,
	wantArgs ...interface{},
) {
	sqlStr, args, err := qb.ToSQL()
	if wantErrKind.Empty() {
		assert.NoError(t, err)
	} else {
		assert.True(t, wantErrKind.Match(err), "%+v", err)
	}

	if wantSQLPlaceholders != "" {
		assert.Equal(t, wantSQLPlaceholders, sqlStr, "Placeholder SQL strings do not match")
		assert.Equal(t, wantArgs, args, "Placeholder Arguments do not match")
	}

	if wantSQLInterpolated == "" {
		return
	}

	if dml, ok := qb.(*Artisan); ok {
		prev := dml.Options
		qb = dml.Interpolate()
		defer func() { dml.Options = prev; qb = dml }()
	}

	sqlStr, args, err = qb.ToSQL() // Call with enabled interpolation
	if wantErrKind.Empty() {
		assert.NoError(t, err)
	} else {
		assert.True(t, wantErrKind.Match(err), "%+v")
	}
	assert.Equal(t, wantSQLInterpolated, sqlStr, "Interpolated SQL strings do not match")
	assert.Nil(t, args, "Artisan should be nil when the SQL string gets interpolated")
}

// compareToSQL2 This function also exists in file dml_public_test.go to
// avoid import cycles when using a single package dedicated for testing.
func compareToSQL2(
	t testing.TB, qb QueryBuilder, wantErrKind errors.Kind,
	wantSQL string, wantArgs ...interface{},
) {
	t.Helper()
	sqlStr, args, err := qb.ToSQL()
	if wantErrKind.Empty() {
		assert.NoError(t, err, "With SQL %q", wantSQL)
	} else {
		assert.True(t, wantErrKind.Match(err), "%+v", err)
	}
	assert.Exactly(t, wantSQL, sqlStr, "SQL strings do not match")
	assert.Exactly(t, wantArgs, args, "Arguments do not match")
}

func compareExecContext(t testing.TB, ex StmtExecer, lastInsertID, rowsAffected int64) (retLastInsertID, retRowsAffected int64) {

	res, err := ex.ExecContext(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, res, "Returned result from ExecContext should not be nil")

	if lastInsertID > 0 {
		retLastInsertID, err = res.LastInsertId()
		assert.NoError(t, err)
		assert.Exactly(t, lastInsertID, retLastInsertID, "Last insert ID do not match")
	}
	if rowsAffected > 0 {
		retRowsAffected, err = res.RowsAffected()
		assert.NoError(t, err)
		assert.Exactly(t, rowsAffected, retRowsAffected, "Affected rows do not match")
	}
	return
}

func notEqualPointers(t *testing.T, o1, o2 interface{}, msgAndArgs ...interface{}) {
	p1 := reflect.ValueOf(o1)
	p2 := reflect.ValueOf(o2)
	if len(msgAndArgs) == 0 {
		msgAndArgs = []interface{}{"Pointers for type o1:%T o2:%T should not be equal", o1, o2}
	}
	assert.NotEqual(t, p1.Pointer(), p2.Pointer(), msgAndArgs...)
}
