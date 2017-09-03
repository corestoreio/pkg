// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"log"
	"os"
	"testing"

	"github.com/corestoreio/errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

///////////////////////////////////////////////////////////////////////////////
// TEST HELPERS
///////////////////////////////////////////////////////////////////////////////

// Returns a session that's not backed by a database
func createFakeSession() *ConnPool {
	cxn, err := NewConnPool()
	if err != nil {
		panic(err)
	}
	return cxn
}

func createRealSession(t testing.TB) *ConnPool {
	dsn := os.Getenv("CS_DSN")
	if dsn == "" {
		t.Skip("Environment variable CS_DSN not found. Skipping ...")
	}
	cxn, err := NewConnPool(
		WithDSN(dsn),
	)
	if err != nil {
		panic(err)
	}
	return cxn
}

func createRealSessionWithFixtures(t testing.TB, c *installFixturesConfig) *ConnPool {
	sess := createRealSession(t)
	installFixtures(sess.DB, c)
	return sess
}

var _ ArgumentsAppender = (*dmlPerson)(nil)
var _ Scanner = (*dmlPerson)(nil)
var _ LastInsertIDAssigner = (*dmlPerson)(nil)
var _ Scanner = (*dmlPersons)(nil)
var _ ArgumentsAppender = (*nullTypedRecord)(nil)
var _ Scanner = (*nullTypedRecord)(nil)

type dmlPerson struct {
	convert RowConvert
	ID      uint64
	Name    string
	Email   NullString
	Key     NullString
}

func (p *dmlPerson) AssignLastInsertID(id int64) {
	p.ID = uint64(id)
}

// RowScan loads a single row from a SELECT statement returning only one row
func (p *dmlPerson) RowScan(r *sql.Rows) error {
	if err := p.convert.Scan(r); err != nil {
		return errors.WithStack(err)
	}
	return p.assign(&p.convert)
}

func (p *dmlPerson) assign(rc *RowConvert) (err error) {
	for i, c := range rc.Columns {
		b := rc.Index(i)
		switch c {
		case "id":
			p.ID, err = b.Uint64()
		case "name":
			p.Name, err = b.Str()
		case "email":
			p.Email.NullString, err = b.NullString()
		case "key":
			p.Key.NullString, err = b.NullString()
			//default:
			//	return errors.NewNotFoundf("[dml_test] Column %q not found", c)
		}
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (p *dmlPerson) AppendArgs(args Arguments, columns []string) (_ Arguments, err error) {
	l := len(columns)
	if l == 1 {
		return p.appendArgs(args, columns[0])
	}
	if l == 0 {
		return args.Uint64(p.ID).Str(p.Name).NullString(p.Email), nil // except auto inc column ;-)
	}
	for _, col := range columns {
		if args, err = p.appendArgs(args, col); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return args, err
}

func (p *dmlPerson) appendArgs(args Arguments, column string) (_ Arguments, err error) {
	switch column {
	case "id":
		args = args.Uint64(p.ID)
	case "name":
		args = args.Str(p.Name)
	case "email":
		args = args.NullString(p.Email)
		// case "key": don't add key, it triggers a test failure condition
	default:
		return nil, errors.NewNotFoundf("[dml_test] dmlPerson Column %q not found", column)
	}
	return args, err
}

type dmlPersons struct {
	convert RowConvert
	Data    []*dmlPerson
}

func (ps *dmlPersons) AppendArgs(args Arguments, columns []string) (_ Arguments, err error) {
	if len(columns) != 1 {
		// INSERT STATEMENT requesting all columns or specific columns
		for _, p := range ps.Data {
			if args, err = p.AppendArgs(args, columns); err != nil {
				return nil, errors.WithStack(err)
			}
		}
		return args, err
	}

	// SELECT, DELETE or UPDATE or INSERT with one column
	column := columns[0]
	var ids []uint64
	var names []string
	var emails []NullString
	for _, p := range ps.Data {
		switch column {
		case "id":
			ids = append(ids, p.ID)
		case "name":
			names = append(names, p.Name)
		case "email":
			emails = append(emails, p.Email)
			// case "key": don't add key, it triggers a test failure condition
		default:
			return nil, errors.NewNotFoundf("[dml_test] dmlPerson Column %q not found", column)
		}
	}

	switch column {
	case "id":
		args = args.Uint64s(ids...)
	case "name":
		args = args.Strs(names...)
	case "email":
		args = args.NullString(emails...)
	}

	return args, nil
}

//func (ps *dmlPersons) AssignLastInsertID(uint64) error {
//	// todo iterate and assign to the last item in the slice and assign
//	// decremented IDs to the previous items in the slice.
//	return nil
//}
//

func (ps *dmlPersons) RowScan(r *sql.Rows) error {
	if err := ps.convert.Scan(r); err != nil {
		return errors.WithStack(err)
	}

	p := new(dmlPerson)
	if err := p.assign(&ps.convert); err != nil {
		return errors.WithStack(err)
	}
	ps.Data = append(ps.Data, p)
	return nil
}

type nullTypedRecord struct {
	ID         int64
	StringVal  NullString
	Int64Val   NullInt64
	Float64Val NullFloat64
	TimeVal    NullTime
	BoolVal    NullBool
}

func (p *nullTypedRecord) RowScan(r *sql.Rows) error {
	return r.Scan(&p.ID, &p.StringVal, &p.Int64Val, &p.Float64Val, &p.TimeVal, &p.BoolVal)
}

func (p *nullTypedRecord) AppendArgs(args Arguments, columns []string) (Arguments, error) {
	for _, column := range columns {
		switch column {
		case "id":
			args = args.Int64(p.ID)
		case "string_val":
			args = args.NullString(p.StringVal)
		case "int64_val":
			if p.Int64Val.Valid {
				args = args.Int64(p.Int64Val.Int64)
			} else {
				args = args.Null()
			}
		case "float64_val":
			if p.Float64Val.Valid {
				args = args.Float64(p.Float64Val.Float64)
			} else {
				args = args.Null()
			}
		case "time_val":
			if p.TimeVal.Valid {
				args = args.Time(p.TimeVal.Time)
			} else {
				args = args.Null()
			}
		case "bool_val":
			if p.BoolVal.Valid {
				args = args.Bool(p.BoolVal.Bool)
			} else {
				args = args.Null()
			}
		default:
			return nil, errors.NewNotFoundf("[dml_test] Column %q not found", columns)
		}
	}
	return args, nil
}

type installFixturesConfig struct {
	AddPeopleWithMaxUint64 bool
}

func installFixtures(db *sql.DB, c *installFixturesConfig) {
	createPeopleTable := fmt.Sprintf(`
		CREATE TABLE dml_people (
			id bigint(8) unsigned NOT NULL auto_increment PRIMARY KEY,
			name varchar(255) NOT NULL,
			email varchar(255),
			%s varchar(255),
			store_id smallint(5) unsigned DEFAULT 0 COMMENT 'Store Id',
			created_at timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT 'Created At',
			total_income decimal(12,4) NOT NULL DEFAULT 0.0000 COMMENT 'Total Income Amount'
		)
	`, "`key`")

	createNullTypesTable := `
		CREATE TABLE null_types (
			id int(11) NOT NULL auto_increment PRIMARY KEY,
			string_val varchar(255) NULL,
			int64_val int(11) NULL,
			float64_val float NULL,
			time_val datetime NULL,
			bool_val bool NULL
		)
	`
	// see also test case "LoadUint64 max Uint64 found"
	sqlToRun := []string{
		"DROP TABLE IF EXISTS dml_people",
		createPeopleTable,
		"INSERT INTO dml_people (name,email) VALUES ('Jonathan', 'jonathan@uservoice.com')",
		"INSERT INTO dml_people (name,email) VALUES ('Dmitri', 'zavorotni@jadius.com')",

		"DROP TABLE IF EXISTS null_types",
		createNullTypesTable,
	}
	if c != nil && c.AddPeopleWithMaxUint64 {
		sqlToRun = append(sqlToRun, "INSERT INTO dml_people (id,name,email) VALUES (18446744073700551613,'Cyrill', 'firstname@lastname.fm')")
	}

	for _, v := range sqlToRun {
		_, err := db.Exec(v)
		if err != nil {
			log.Fatalln("Failed to execute statement: ", v, " Got error: ", err)
		}
	}
}

var _ Querier = (*dbMock)(nil)
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
	t testing.TB, qb QueryBuilder, wantErr errors.BehaviourFunc,
	wantSQLPlaceholders, wantSQLInterpolated string,
	wantArgs ...interface{},
) {

	sqlStr, args, err := qb.ToSQL()
	if wantErr == nil {
		require.NoError(t, err)
	} else {
		require.True(t, wantErr(err), "%+v", err)
	}

	if wantSQLPlaceholders != "" {
		assert.Equal(t, wantSQLPlaceholders, sqlStr, "Placeholder SQL strings do not match")
		assert.Equal(t, wantArgs, args, "Placeholder Arguments do not match")
	}

	if wantSQLInterpolated == "" {
		return
	}

	// If you care regarding the duplication ... send us a PR ;-)
	// Enables Interpolate feature and resets it after the test has been
	// executed.
	switch dml := qb.(type) {
	case *Delete:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *Update:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *Insert:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *Select:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *Union:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *With:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *Show:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	default:
		t.Fatalf("func compareToSQL: the type %#v is not (yet) supported.", qb)
	}

	sqlStr, args, err = qb.ToSQL() // Call with enabled interpolation
	require.Nil(t, args, "Arguments should be nil when the SQL string gets interpolated")
	if wantErr == nil {
		require.NoError(t, err)
	} else {
		require.True(t, wantErr(err), "%+v")
	}
	require.Equal(t, wantSQLInterpolated, sqlStr, "Interpolated SQL strings do not match")
}
