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
	"os"
	"testing"
	"time"

	"fmt"
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

var _ dml.ColumnMapper = (*dmlPerson)(nil)

type dmlPerson struct {
	ID          int64
	Name        string
	Email       dml.NullString
	Key         dml.NullString
	StoreID     int64
	CreatedAt   time.Time
	TotalIncome float64
}

func (p *dmlPerson) AssignLastInsertID(id int64) {
	p.ID = id
}

func (p *dmlPerson) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm.Int64(&p.ID).String(&p.Name).NullString(&p.Email).NullString(&p.Key).Int64(&p.StoreID).Time(&p.CreatedAt).Float64(&p.TotalIncome).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "id":
			cm.Int64(&p.ID)
		case "name":
			cm.String(&p.Name)
		case "email":
			cm.NullString(&p.Email)
		case "key":
			cm.NullString(&p.Key)
		case "store_id":
			cm.Int64(&p.StoreID)
		case "created_at":
			cm.Time(&p.CreatedAt)
		case "total_income":
			cm.Float64(&p.TotalIncome)
		default:
			return errors.NewNotFoundf("[dml_test] dmlPerson Column %q not found", c)
		}
	}
	return cm.Err()
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

	qb, reset := enableInterpolate(qb)
	defer reset()

	sqlStr, args, err = qb.ToSQL()
	require.Nil(t, args, "Arguments should be nil when the SQL string gets interpolated")
	if wantErr == nil {
		require.NoError(t, err)
	} else {
		require.True(t, wantErr(err), "%+v")
	}
	require.Equal(t, wantSQLInterpolated, sqlStr, "Interpolated SQL strings do not match")
}

func enableInterpolate(qb dml.QueryBuilder) (_ dml.QueryBuilder, reset func()) {
	switch cqb := qb.(type) {
	case *dml.Delete:
		cqb.Interpolate()
		reset = func() { cqb.IsInterpolate = false }
	case *dml.Update:
		cqb.Interpolate()
		reset = func() { cqb.IsInterpolate = false }
	case *dml.Insert:
		cqb.Interpolate()
		reset = func() { cqb.IsInterpolate = false }
	case *dml.Select:
		cqb.Interpolate()
		reset = func() { cqb.IsInterpolate = false }
	case *dml.Union:
		cqb.Interpolate()
		reset = func() { cqb.IsInterpolate = false }
	case *dml.With:
		cqb.Interpolate()
		reset = func() { cqb.IsInterpolate = false }
	default:
		panic(fmt.Sprintf("Type %#v not (yet) supported.", qb))
	}
	return qb, reset
}
