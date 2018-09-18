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

package ddl

import (
	"context"
	"math"
	"sort"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
)

var (
	_ dml.ColumnMapper = (*Variables)(nil)
	_ dml.QueryBuilder = (*Variables)(nil)
	_ errors.Kinder    = (*errTableNotFound)(nil)
	_ error            = (*errTableNotFound)(nil)
)

func TestErrTableNotFound(t *testing.T) {
	t.Parallel()
	err := errTableNotFound("Errr")
	assert.Exactly(t, errors.NotFound, err.ErrorKind())
	assert.Exactly(t, "[ddl] Table \"Errr\" not found or not yet added.", err.Error())
}

func TestNewVariables_Integration(t *testing.T) {
	t.Parallel()

	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	vs := NewVariables()
	_, err := db.WithQueryBuilder(vs).Load(context.TODO(), vs)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Exactly(t, "InnoDB", vs.Data["storage_engine"])
	assert.True(t, len(vs.Data) > 400, "Should have more than 400 map entries")
}

func TestNewVariables_Mock(t *testing.T) {
	t.Parallel()

	dbc, dbMock := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, dbMock)

	t.Run("one with LIKE", func(t *testing.T) {
		var mockedRows = sqlmock.NewRows([]string{"Variable_name", "Value"}).
			FromCSVString("keyVal11,helloAustralia")

		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SHOW VARIABLES WHERE (`Variable_name` LIKE 'keyVal11')")).
			WillReturnRows(mockedRows)

		vs := NewVariables("keyVal11")
		rc, err := dbc.WithQueryBuilder(vs).Load(context.TODO(), vs)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Exactly(t, uint64(1), rc, "Should load one row")

		assert.Exactly(t, `helloAustralia`, vs.Data["keyVal11"])
		assert.Len(t, vs.Data, 1)
	})

	t.Run("many with WHERE", func(t *testing.T) {
		var mockedRows = sqlmock.NewRows([]string{"Variable_name", "Value"}).
			FromCSVString("keyVal11,helloAustralia\nkeyVal22,helloNewZealand")

		dbMock.ExpectQuery(dmltest.SQLMockQuoteMeta("SHOW VARIABLES WHERE (`Variable_name` IN ('keyVal11','keyVal22'))")).
			WillReturnRows(mockedRows)

		vs := NewVariables("keyVal11", "keyVal22")
		rc, err := dbc.WithQueryBuilder(vs).Load(context.TODO(), vs)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		assert.Exactly(t, uint64(2), rc, "Shoud load two rows")

		assert.Exactly(t, `helloAustralia`, vs.Data["keyVal11"])
		assert.Exactly(t, `helloNewZealand`, vs.Data["keyVal22"])
		assert.Len(t, vs.Data, 2)
		keys := vs.Keys()
		sort.Strings(keys)
		assert.Exactly(t, []string{"keyVal11", "keyVal22"}, keys)
	})
}

func TestVariables_Equal(t *testing.T) {
	t.Parallel()

	v := NewVariables("any_key")

	v.Data["any_key"] = "A"
	assert.True(t, v.Equal("any_key", "A"))
	assert.False(t, v.Equal("any_key", "B"))
	assert.False(t, v.Equal("any_key", "a"))

	assert.True(t, v.EqualFold("any_key", "a"))
	assert.False(t, v.EqualFold("any_key", "B"))
}

func TestVariables_Types(t *testing.T) {
	t.Parallel()

	v := NewVariables("any_key")
	v.Data["float64_ok"] = "3.14159"
	v.Data["float64_nok"] = "ø"
	v.Data["int64_ok"] = strconv.FormatInt(-math.MaxInt64, 10)
	v.Data["int64_nok"] = "ø"
	v.Data["uint64_ok"] = strconv.FormatUint(math.MaxUint64, 10)
	v.Data["uint64_nok"] = "ø"

	v.Data["bool_YES"] = "YES"
	v.Data["bool_NO"] = "NO"
	v.Data["bool_ON"] = "ON"
	v.Data["bool_OFF"] = "OFF"
	v.Data["bool_nok1"] = "1"
	v.Data["bool_nok0"] = "0"

	t.Run("string", func(t *testing.T) {
		val, ok := v.String("float64_ok")
		assert.True(t, ok)
		assert.Exactly(t, "3.14159", val)

		val, ok = v.String("not_found")
		assert.False(t, ok)
		assert.Empty(t, val)
	})

	t.Run("bool", func(t *testing.T) {
		val, ok := v.Bool("bool_YES")
		assert.True(t, ok)
		assert.True(t, val)

		val, ok = v.Bool("bool_NO")
		assert.True(t, ok)
		assert.False(t, val)

		val, ok = v.Bool("bool_ON")
		assert.True(t, ok)
		assert.True(t, val)

		val, ok = v.Bool("bool_OFF")
		assert.True(t, ok)
		assert.False(t, val)

		val, ok = v.Bool("bool_nok1")
		assert.False(t, ok)
		assert.False(t, val)

		val, ok = v.Bool("bool_nok0")
		assert.False(t, ok)
		assert.False(t, val)

		val, ok = v.Bool("not_found")
		assert.False(t, ok)
		assert.False(t, val)
	})

	t.Run("float64", func(t *testing.T) {
		val, ok := v.Float64("float64_ok")
		assert.True(t, ok)
		assert.Exactly(t, 3.14159, val)

		val, ok = v.Float64("float64_nok")
		assert.False(t, ok)
		assert.Exactly(t, float64(0), val)

		val, ok = v.Float64("not_found")
		assert.False(t, ok)
		assert.Exactly(t, float64(0), val)
	})

	t.Run("int64", func(t *testing.T) {
		val, ok := v.Int64("int64_ok")
		assert.True(t, ok)
		assert.Exactly(t, int64(-math.MaxInt64), val)

		val, ok = v.Int64("int64_nok")
		assert.False(t, ok)
		assert.Exactly(t, int64(0), val)

		val, ok = v.Int64("not_found")
		assert.False(t, ok)
		assert.Exactly(t, int64(0), val)
	})

	t.Run("uint64", func(t *testing.T) {
		val, ok := v.Uint64("uint64_ok")
		assert.True(t, ok)
		assert.Exactly(t, uint64(math.MaxUint64), val)

		val, ok = v.Uint64("uint64_nok")
		assert.False(t, ok)
		assert.Exactly(t, uint64(0), val)

		val, ok = v.Uint64("not_found")
		assert.False(t, ok)
		assert.Exactly(t, uint64(0), val)
	})

}
