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

func createRealSession() *Connection {
	dsn := os.Getenv("CS_DSN")
	cxn, err := NewConnection(
		WithDSN(dsn),
	)
	if err != nil {
		panic(err)
	}
	return cxn
}

func createRealSessionWithFixtures() *Connection {
	sess := createRealSession()
	installFixtures(sess.DB)
	return sess
}

var _ ArgumentAssembler = (*dbrPerson)(nil)
var _ Scanner = (*dbrPerson)(nil)
var _ Scanner = (*dbrPersons)(nil)
var _ ArgumentAssembler = (*nullTypedRecord)(nil)
var _ Scanner = (*nullTypedRecord)(nil)

type dbrPerson struct {
	ID    int64
	Name  string
	Email NullString
	Key   NullString
}

// ScanRow loads a single row from a SELECT statement returning only one row
func (p *dbrPerson) ScanRow(idx int, columns []string, scan func(dest ...interface{}) error) error {
	if idx > 0 {
		return errors.NewExceededf("[dbr_test] Can only load one row. Got a next row.")
	}
	vp := make([]interface{}, 0, 4) // vp == valuePointers
	for _, c := range columns {
		switch c {
		case "id":
			vp = append(vp, &p.ID)
		case "name":
			vp = append(vp, &p.Name)
		case "email":
			vp = append(vp, &p.Email)
		case "key":
			vp = append(vp, &p.Key)
		default:
			return errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}
	return scan(vp...)
}

func (p *dbrPerson) AssembleArguments(stmtType int, args Arguments, columns []string) (_ Arguments, err error) {
	for _, c := range columns {
		switch c {
		case "id", "dp.id":
			//if t == 'i' {
			args = append(args, ArgInt64(p.ID))
			//}
		case "name":
			args = append(args, ArgString(p.Name))
		case "email":
			args = append(args, ArgNullString(p.Email))
			// case "key": don't add key, it triggers a test failure condition
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}
	return args, err
}

type dbrPersons struct {
	Data     []*dbrPerson
	scanArgs []interface{}
	dto      dbrPerson
}

func (ps *dbrPersons) ScanRow(idx int, columns []string, scan func(dest ...interface{}) error) error {
	if idx == 0 {
		ps.Data = make([]*dbrPerson, 0, 5)
		ps.scanArgs = make([]interface{}, 0, 4) // four fields in the struct

		for _, c := range columns {
			switch c {
			case "id":
				ps.scanArgs = append(ps.scanArgs, &ps.dto.ID)
			case "name":
				ps.scanArgs = append(ps.scanArgs, &ps.dto.Name)
			case "email":
				ps.scanArgs = append(ps.scanArgs, &ps.dto.Email)
			case "key":
				ps.scanArgs = append(ps.scanArgs, &ps.dto.Key)
			default:
				return errors.NewNotFoundf("[dbr_test] Column %q not found", c)
			}
		}
	}

	if err := scan(ps.scanArgs...); err != nil {
		return errors.Wrap(err, "[dbr_test] dbrPersons.ScanRow")
	}
	ps.Data = append(ps.Data, &dbrPerson{
		ID:    ps.dto.ID,
		Name:  ps.dto.Name,
		Email: ps.dto.Email,
		Key:   ps.dto.Key,
	})
	return nil
}

var _ ArgumentAssembler = (*nullTypedRecord)(nil)

type nullTypedRecord struct {
	ID         int64
	StringVal  NullString
	Int64Val   NullInt64
	Float64Val NullFloat64
	TimeVal    NullTime
	BoolVal    NullBool
}

func (p *nullTypedRecord) ScanRow(idx int, columns []string, scan func(dest ...interface{}) error) error {
	if idx > 0 {
		return errors.NewExceededf("[dbr_test] Can only load one row. Got a next row.")
	}
	return scan(&p.ID, &p.StringVal, &p.Int64Val, &p.Float64Val, &p.TimeVal, &p.BoolVal)
}

func (p *nullTypedRecord) AssembleArguments(stmtType int, args Arguments, columns []string) (Arguments, error) {
	for _, c := range columns {
		switch c {
		case "id":
			args = append(args, ArgInt64(p.ID))
		case "string_val":
			args = append(args, ArgNullString(p.StringVal))
		case "int64_val":
			if p.Int64Val.Valid {
				args = append(args, ArgInt64(p.Int64Val.Int64))
			} else {
				args = append(args, ArgNull())
			}
		case "float64_val":
			if p.Float64Val.Valid {
				args = append(args, ArgFloat64(p.Float64Val.Float64))
			} else {
				args = append(args, ArgNull())
			}
		case "time_val":
			if p.TimeVal.Valid {
				args = append(args, ArgTime(p.TimeVal.Time))
			} else {
				args = append(args, ArgNull())
			}
		case "bool_val":
			if p.BoolVal.Valid {
				args = append(args, ArgBool(p.BoolVal.Bool))
			} else {
				args = append(args, ArgNull())
			}
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}

	return args, nil
}

func installFixtures(db *sql.DB) {
	createPeopleTable := fmt.Sprintf(`
		CREATE TABLE dbr_people (
			id int(11) NOT NULL auto_increment PRIMARY KEY,
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

	sqlToRun := []string{
		"DROP TABLE IF EXISTS dbr_people",
		createPeopleTable,
		"INSERT INTO dbr_people (name,email) VALUES ('Jonathan', 'jonathan@uservoice.com')",
		"INSERT INTO dbr_people (name,email) VALUES ('Dmitri', 'zavorotni@jadius.com')",

		"DROP TABLE IF EXISTS null_types",
		createNullTypesTable,
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

	if wantSQLPlaceholders != "" {
		sqlStr, args, err := qb.ToSQL()
		if wantErr == nil {
			require.NoError(t, err, "%+v", err)
		} else {
			require.True(t, wantErr(err), "%+v")
		}
		assert.Equal(t, wantSQLPlaceholders, sqlStr, "Placeholder SQL strings do not match")
		assert.Equal(t, wantArgs, args.Interfaces(), "Placeholder Arguments do not match")
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
	default:
		t.Fatalf("Type %#v not (yet) supported.", qb)
	}

	sqlStr, args, err := qb.ToSQL()
	require.Nil(t, args, "Arguments should be nil when the SQL string gets interpolated")
	if wantErr == nil {
		require.NoError(t, err, "%+v", err)
	} else {
		require.True(t, wantErr(err), "%+v")
	}
	require.Equal(t, wantSQLInterpolated, sqlStr, "Interpolated SQL strings do not match")
}
