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

package dbr

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

//
// Test helpers
//

// Returns a session that's not backed by a database
func createFakeSession() *Connection {
	cxn, err := NewConnection()
	if err != nil {
		panic(err)
	}
	return cxn
}

func createRealSession(t testing.TB) *Connection {
	dsn := os.Getenv("CS_DSN")
	if dsn == "" {
		t.Skip("Environment variable CS_DSN not found. Skipping ...")
	}
	cxn, err := NewConnection(
		WithDSN(dsn),
	)
	if err != nil {
		panic(err)
	}
	return cxn
}

func createRealSessionWithFixtures(t testing.TB, c *installFixturesConfig) *Connection {
	sess := createRealSession(t)
	installFixtures(sess.DB, c)
	return sess
}

var _ ArgumentsAppender = (*dbrPerson)(nil)
var _ Scanner = (*dbrPerson)(nil)
var _ Scanner = (*dbrPersons)(nil)
var _ ArgumentsAppender = (*nullTypedRecord)(nil)
var _ Scanner = (*nullTypedRecord)(nil)

type dbrPerson struct {
	convert RowConvert
	ID      uint64
	Name    string
	Email   NullString
	Key     NullString
}

func assignDbrPerson(p *dbrPerson, rc *RowConvert) error {
	for i, c := range rc.Columns {
		b := rc.Index(i)
		var err error
		switch c {
		case "id":
			p.ID, err = b.Uint64()
		case "name":
			p.Name, err = b.Str()
		case "email":
			p.Email.NullString, err = b.NullString()
		case "key":
			p.Key.NullString, err = b.NullString()
		default:
			return errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// RowScan loads a single row from a SELECT statement returning only one row
func (p *dbrPerson) RowScan(r *sql.Rows) error {
	if err := p.convert.Scan(r); err != nil {
		return errors.WithStack(err)
	}
	return assignDbrPerson(p, &p.convert)
}

func personAppendArguments(p *dbrPerson, args Arguments, columns []string) (_ Arguments, err error) {
	for _, c := range columns {
		switch c {
		case "id", "dp.id":
			args = append(args, Int64(p.ID))
		case "name":
			args = append(args, String(p.Name))
		case "email":
			args = append(args, NullString(p.Email))
			// case "key": don't add key, it triggers a test failure condition
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}
	return args, err
}

func (p *dbrPerson) AppendArguments(stmtType int, args Arguments, columns []string) (_ Arguments, err error) {
	return personAppendArguments(p, args, columns)
}

type dbrPersons struct {
	convert RowConvert
	Data    []*dbrPerson
}

func (ps *dbrPersons) AppendArguments(stmtType int, args Arguments, columns []string) (_ Arguments, err error) {
	for _, p := range ps.Data {
		args, err = personAppendArguments(p, args, columns)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return args, nil
}

func (ps *dbrPersons) RowScan(r *sql.Rows) error {
	if err := ps.convert.Scan(r); err != nil {
		return errors.WithStack(err)
	}

	p := new(dbrPerson)
	if err := assignDbrPerson(p, &ps.convert); err != nil {
		return errors.WithStack(err)
	}
	ps.Data = append(ps.Data, p)
	return nil
}

var _ ArgumentsAppender = (*nullTypedRecord)(nil)

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

func (p *nullTypedRecord) AppendArguments(stmtType int, args Arguments, columns []string) (Arguments, error) {
	for _, c := range columns {
		switch c {
		case "id":
			args = append(args, Int64(p.ID))
		case "string_val":
			args = append(args, NullString(p.StringVal))
		case "int64_val":
			if p.Int64Val.Valid {
				args = append(args, Int64(p.Int64Val.Int64))
			} else {
				args = append(args, NullValue())
			}
		case "float64_val":
			if p.Float64Val.Valid {
				args = append(args, Float64(p.Float64Val.Float64))
			} else {
				args = append(args, NullValue())
			}
		case "time_val":
			if p.TimeVal.Valid {
				args = append(args, MakeTime(p.TimeVal.Time))
			} else {
				args = append(args, NullValue())
			}
		case "bool_val":
			if p.BoolVal.Valid {
				args = append(args, Bool(p.BoolVal.Bool))
			} else {
				args = append(args, NullValue())
			}
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}

	return args, nil
}

type installFixturesConfig struct {
	AddPeopleWithMaxUint64 bool
}

func installFixtures(db *sql.DB, c *installFixturesConfig) {
	createPeopleTable := fmt.Sprintf(`
		CREATE TABLE dbr_people (
			id bigint(8) unsigned NOT NULL auto_increment PRIMARY KEY,
			name varchar(255) NOT NULL,
			email varchar(255),
			%s varchar(255)
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
		"DROP TABLE IF EXISTS dbr_people",
		createPeopleTable,
		"INSERT INTO dbr_people (name,email) VALUES ('Jonathan', 'jonathan@uservoice.com')",
		"INSERT INTO dbr_people (name,email) VALUES ('Dmitri', 'zavorotni@jadius.com')",

		"DROP TABLE IF EXISTS null_types",
		createNullTypesTable,
	}
	if c != nil && c.AddPeopleWithMaxUint64 {
		sqlToRun = append(sqlToRun, "INSERT INTO dbr_people (id,name,email) VALUES (18446744073700551613,'Cyrill', 'firstname@lastname.fm')")
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
// interpolated string. This function also exists in file dbr_public_test.go to
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
