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

package dbr_test

import (
	"os"
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ dbr.ArgumentsAppender = (*dbrPerson)(nil)

type dbrPerson struct {
	ID    int64 `db:"id"`
	Name  string
	Email dbr.NullString
	Key   dbr.NullString
}

func (p *dbrPerson) AppendArguments(_ dbr.SQLStmt, args dbr.Arguments, columns []string) (dbr.Arguments, error) {
	for _, c := range columns {
		switch c {
		case "id":
			args = args.Int64(p.ID)
		case "name":
			args = args.Str(p.Name)
		case "email":
			args = args.NullString(p.Email)
		case "key":
			args = args.NullString(p.Key)
		default:
			return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
		}
	}
	return args, nil
}

func createRealSession(t testing.TB) *dbr.Connection {
	dsn := os.Getenv("CS_DSN")
	if dsn == "" {
		t.Skip("Environment variable CS_DSN not found. Skipping ...")
	}
	cxn, err := dbr.NewConnection(
		dbr.WithDSN(dsn),
	)
	if err != nil {
		panic(err)
	}
	return cxn
}

// compareToSQL compares a SQL object with a placeholder string and an optional
// interpolated string. This function also exists in file dbr_public_test.go to
// avoid import cycles when using a single package dedicated for testing.
func compareToSQL(
	t testing.TB, qb dbr.QueryBuilder, wantErr errors.BehaviourFunc,
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
	case *dbr.Delete:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *dbr.Update:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *dbr.Insert:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *dbr.Select:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *dbr.Union:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *dbr.With:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	default:
		t.Fatalf("Type %#v not (yet) supported.", qb)
	}

	sqlStr, args, err = qb.ToSQL()
	require.Nil(t, args, "Arguments should be nil when the SQL string gets interpolated")
	if wantErr == nil {
		require.NoError(t, err)
	} else {
		require.True(t, wantErr(err), "%+v")
	}
	require.Equal(t, wantSQLInterpolated, sqlStr, "Interpolated SQL strings do not match")
}
