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

package dml_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

var now = func() time.Time {
	return time.Date(2006, 1, 2, 15, 4, 5, 02, time.FixedZone("hardcoded", -7))
}

func init() {
	// Freeze time in package log
	log.Now = now
	null.JSONMarshalFn = json.Marshal
	null.JSONUnMarshalFn = json.Unmarshal
}

var _ dml.ColumnMapper = (*dmlPerson)(nil)

type dmlPerson struct {
	ID          int64
	Name        string
	Email       null.String
	Key         null.String
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
		case "name", "name2": // name2 used in TestWithLogger_WithCTE
			cm.String(&p.Name)
		case "email", "email2": // email2 used in TestWithLogger_WithCTE
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
			return errors.NotFound.Newf("[dml_test] dmlPerson Column %q not found", c)
		}
	}
	return cm.Err()
}

func createRealSession(t testing.TB, opts ...dml.ConnPoolOption) *dml.ConnPool {
	dsn := dmltest.MustGetDSN(t)
	cxn, err := dml.NewConnPool(
		append([]dml.ConnPoolOption{dml.WithDSN(dsn)}, opts...)...,
	)
	if err != nil {
		t.Fatal(err)
	}
	return cxn
}

// compareToSQL compares a SQL object with a placeholder string and an optional
// interpolated string. This function also exists in file dml_public_test.go to
// avoid import cycles when using a single package dedicated for testing.
func compareToSQL(
	t testing.TB, qb dml.QueryBuilder, wantErrKind errors.Kind,
	wantSQLPlaceholders, wantSQLInterpolated string,
	wantArgs ...interface{},
) {
	sqlStr, args, err := qb.ToSQL()
	if wantErrKind.Empty() {
		assert.NoError(t, err)
	} else {
		assert.ErrorIsKind(t, wantErrKind, err)
	}

	if wantSQLPlaceholders != "" {
		assert.Exactly(t, wantSQLPlaceholders, sqlStr, "Placeholder SQL strings do not match")
		assert.Exactly(t, wantArgs, args, "Placeholder Arguments do not match")
	}

	if wantSQLInterpolated == "" {
		return
	}

	if dmlArgs, ok := qb.(*dml.Artisan); ok {
		prev := dmlArgs.Options
		qb = dmlArgs.Interpolate()
		defer func() { dmlArgs.Options = prev; qb = dmlArgs }()
	}

	sqlStr, args, err = qb.ToSQL() // Call with enabled interpolation
	assert.Nil(t, args, "Artisan should be nil when the SQL string gets interpolated")
	if wantErrKind.Empty() {
		assert.NoError(t, err)
	} else {
		assert.ErrorIsKind(t, wantErrKind, err)
	}
	assert.Exactly(t, wantSQLInterpolated, sqlStr, "Interpolated SQL strings do not match")
}

func ifNotEqualPanic(have, want interface{}, msg ...string) {
	// The reason for this function is that I have no idea why testing.T is
	// blocking inside the bgwork.Wait function.
	if !reflect.DeepEqual(have, want) {
		panic(fmt.Sprintf("%q\nHave: %#v\nWant: %#v\n\n", strings.Join(msg, ""), have, want))
	}
}

func ifErrPanic(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}

func notEqualPointers(t *testing.T, o1, o2 interface{}, msgAndArgs ...interface{}) {
	p1 := reflect.ValueOf(o1)
	p2 := reflect.ValueOf(o2)
	if len(msgAndArgs) == 0 {
		msgAndArgs = []interface{}{"Pointers for type o1:%T o2:%T should not be equal", o1, o2}
	}
	assert.NotEqual(t, p1.Pointer(), p2.Pointer(), msgAndArgs...)
}
