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

package dml_test

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var now = func() time.Time {
	return time.Date(2006, 1, 2, 15, 4, 5, 02, time.FixedZone("hardcoded", -7))
}

func init() {
	// Freeze time in package log
	log.Now = now
}

var _ dml.ArgumentsAppender = (*dmlPerson)(nil)
var _ dml.Scanner = (*dmlPerson)(nil)

type dmlPerson struct {
	ID          int64
	Name        string
	Email       dml.NullString
	Key         dml.NullString
	StoreID     int64
	CreatedAt   time.Time
	TotalIncome float64
}

func (p *dmlPerson) RowScan(r *sql.Rows) error {
	// noop
	return nil
}

func (p *dmlPerson) AssignLastInsertID(id int64) {
	p.ID = id
}

func (p *dmlPerson) AppendArgs(args dml.Arguments, columns []string) (_ dml.Arguments, err error) {
	l := len(columns)
	if l == 1 {
		return p.appendArgs(args, columns[0])
	}
	if l == 0 {
		return args.Int64(p.ID).String(p.Name).NullString(p.Email).NullString(p.Key).Int64(p.StoreID).Time(p.CreatedAt).Float64(p.TotalIncome), nil
	}
	for _, col := range columns {
		if args, err = p.appendArgs(args, col); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return args, err
}

func (p *dmlPerson) appendArgs(args dml.Arguments, column string) (_ dml.Arguments, err error) {
	switch column {
	case "id":
		args = args.Int64(p.ID)
	case "name":
		args = args.String(p.Name)
	case "email":
		args = args.NullString(p.Email)
	case "key":
		args = args.NullString(p.Key)
	case "store_id":
		args = args.Int64(p.StoreID)
	case "created_at":
		args = args.Time(p.CreatedAt)
	case "total_income":
		args = args.Float64(p.TotalIncome)
	default:
		return nil, errors.NewNotFoundf("[dml_test] dmlPerson Column %q not found", column)
	}
	return args, err
}

func createRealSession(t testing.TB) *dml.ConnPool {
	dsn := os.Getenv("CS_DSN")
	if dsn == "" {
		t.Skip("Environment variable CS_DSN not found. Skipping ...")
	}
	cxn, err := dml.NewConnPool(
		dml.WithDSN(dsn),
	)
	if err != nil {
		panic(err)
	}
	return cxn
}

// compareToSQL compares a SQL object with a placeholder string and an optional
// interpolated string. This function also exists in file dml_public_test.go to
// avoid import cycles when using a single package dedicated for testing.
func compareToSQL(
	t testing.TB, qb dml.QueryBuilder, wantErr errors.BehaviourFunc,
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
	case *dml.Delete:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *dml.Update:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *dml.Insert:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *dml.Select:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *dml.Union:
		dml.Interpolate()
		defer func() { dml.IsInterpolate = false }()
	case *dml.With:
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
