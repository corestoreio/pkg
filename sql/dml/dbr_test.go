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
	"bytes"
	"database/sql"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

func toIFaceSlice(args ...interface{}) []interface{} {
	for i, a := range args {
		if a == nil {
			args[i] = internalNULLNIL{}
		}
	}
	return args
}

func TestArguments_Interfaces(t *testing.T) {
	container := make([]interface{}, 0, 48)

	t.Run("no slices, nulls valid", func(t *testing.T) {
		args := toIFaceSlice(
			nil, -1, int64(1), uint64(2), 3.1, true, "eCom1", []byte(`eCom2`), now(),
			null.MakeString("eCom3"), null.MakeInt64(4), null.MakeFloat64(2.7),
			null.MakeBool(true), null.MakeTime(now()))

		assert.Exactly(t,
			[]interface{}{
				nil, int64(-1), int64(1), int64(2), 3.1, true, "eCom1",
				[]uint8{0x65, 0x43, 0x6f, 0x6d, 0x32},
				now(),
				"eCom3", int64(4), 2.7, true, now(),
			},
			expandInterfaces(args))
		container = container[:0]
	})
	t.Run("no slices, nulls invalid", func(t *testing.T) {
		args := toIFaceSlice(
			nil, -1, int64(1), uint64(2), 3.1, true, "eCom1", []byte(`eCom2`), now(),
			null.String{}, null.Int64{}, null.Float64{},
			null.Bool{}, null.Time{})
		assert.Exactly(t,
			[]interface{}{
				nil, int64(-1), int64(1), int64(2), 3.1, true, "eCom1",
				[]uint8{0x65, 0x43, 0x6f, 0x6d, 0x32},
				now(),
				nil, nil, nil, nil, nil,
			},
			expandInterfaces(args))
		container = container[:0]
	})
	t.Run("slices, nulls valid", func(t *testing.T) {
		args := toIFaceSlice(
			nil, []int{-1, -2}, []int64{1, 2}, []uint{568, 766}, []uint64{2}, []float64{1.2, 3.1}, []bool{false, true},
			[]string{"eCom1", "eCom11"}, [][]byte{[]byte(`eCom2`)}, []time.Time{now(), now()},
			[]null.String{null.MakeString("eCom3"), null.MakeString("eCom3")}, []null.Int64{null.MakeInt64(4), null.MakeInt64(4)},
			[]null.Float64{null.MakeFloat64(2.7), null.MakeFloat64(2.7)},
			[]null.Bool{null.MakeBool(true)}, []null.Time{null.MakeTime(now()), null.MakeTime(now())})
		assert.Exactly(t,
			[]interface{}{
				nil, int64(-1), int64(-2), int64(1), int64(2), int64(568), int64(766), int64(2), 1.2, 3.1, false, true,
				"eCom1", "eCom11",
				[]uint8{0x65, 0x43, 0x6f, 0x6d, 0x32},
				now(), now(),
				"eCom3", "eCom3", int64(4), int64(4),
				2.7, 2.7,
				true, now(), now(),
			},
			expandInterfaces(args))
	})
	t.Run("returns nil interface", func(t *testing.T) {
		assert.Nil(t, expandInterfaces([]interface{}{}), "args.expandInterfaces() must return nil")
	})
}

func TestArguments_DriverValue(t *testing.T) {
	t.Run("Driver.Values supported types", func(t *testing.T) {
		args := toIFaceSlice(
			driverValueNil(0),
			driverValueBytes(nil), null.MakeInt64(3), null.MakeFloat64(2.7), null.MakeBool(true),
			driverValueBytes(`Invoice`), null.MakeString("Creditmemo"), nowSentinel{}, null.MakeTime(now()),
		)
		assert.Exactly(t,
			[]interface{}{
				nil, []uint8(nil), int64(3), 2.7, true,
				[]uint8{0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65},
				"Creditmemo", "2006-01-02 19:04:05", now(),
			},
			expandInterfaces(args))
	})

	t.Run("Driver.Value supported types", func(t *testing.T) {
		args := toIFaceSlice(
			driverValueNil(0),
			driverValueBytes(nil),
			null.MakeInt64(3),
			null.MakeFloat64(2.7),
			null.MakeBool(true),
			driverValueBytes(`Invoice`),
			null.MakeString("Creditmemo"),
			nowSentinel{},
			null.MakeTime(now()))

		assert.Exactly(t,
			[]interface{}{
				nil, []uint8(nil), int64(3), 2.7, true,
				[]uint8{0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65},
				"Creditmemo", "2006-01-02 19:04:05", now(),
			},
			expandInterfaces(args))
	})

	t.Run("Driver.Values panics because not supported", func(t *testing.T) {
		_, err := driverValue(nil, driverValueNotSupported(4))
		assert.ErrorIsKind(t, errors.NotSupported, err)
	})

	t.Run("Driver.Values panics because Value error", func(t *testing.T) {
		_, err := driverValue(nil, driverValueError(0))
		assert.ErrorIsKind(t, errors.Fatal, err)
	})
}

func TestArguments_WriteTo(t *testing.T) {
	t.Run("no slices, nulls valid", func(t *testing.T) {
		args := toIFaceSlice(
			nil, -1, int64(1), uint64(2), 3.1, true, "eCom1", []byte(`eCom2`), now(),
			null.MakeString("eCom3"), null.MakeInt64(4), null.MakeFloat64(2.7),
			null.MakeBool(true), null.MakeTime(now()))

		var buf bytes.Buffer
		assert.NoError(t, writeInterfaces(&buf, args))
		assert.Exactly(t,
			"(NULL,-1,1,2,3.1,1,'eCom1','eCom2','2006-01-02 15:04:05','eCom3',4,2.7,1,'2006-01-02 15:04:05')",
			buf.String())
	})
	t.Run("no slices, nulls invalid", func(t *testing.T) {
		args := toIFaceSlice(
			nil, -1, int64(1), uint64(2), 3.1, true, "eCom1", []byte(`eCom2`), now(),
			null.String{}, null.Int64{}, null.Float64{},
			null.Bool{}, null.Time{})

		var buf bytes.Buffer
		assert.NoError(t, writeInterfaces(&buf, args))
		assert.Exactly(t,
			"(NULL,-1,1,2,3.1,1,'eCom1','eCom2','2006-01-02 15:04:05',NULL,NULL,NULL,NULL,NULL)",
			buf.String())
	})
	t.Run("slices, nulls valid", func(t *testing.T) {
		args := toIFaceSlice(
			nil, []int{-1, -2}, []int64{1, 2}, []uint{568, 766}, []uint64{2}, []float64{1.2, 3.1}, []bool{false, true},
			[]string{"eCom1", "eCom11"}, [][]byte{[]byte(`eCom2`)}, []time.Time{now(), now()},
			[]null.String{null.MakeString("eCom3"), null.MakeString("eCom3")}, []null.Int64{null.MakeInt64(4), null.MakeInt64(4)},
			[]null.Float64{null.MakeFloat64(2.7), null.MakeFloat64(2.7)},
			[]null.Bool{null.MakeBool(true)}, []null.Time{null.MakeTime(now()), null.MakeTime(now())})

		var buf bytes.Buffer
		assert.NoError(t, writeInterfaces(&buf, args))
		assert.Exactly(t,
			"(NULL,(-1,-2),(1,2),(568,766),(2),(1.2,3.1),(0,1),('eCom1','eCom11'),('eCom2'),('2006-01-02 15:04:05','2006-01-02 15:04:05'),('eCom3','eCom3'),(4,4),(2.7,2.7),(1),('2006-01-02 15:04:05','2006-01-02 15:04:05'))",
			buf.String(), "%q", buf.String())
	})
	t.Run("non-utf8 string", func(t *testing.T) {
		args := toIFaceSlice("\xc0\x80")
		var buf bytes.Buffer
		err := writeInterfaces(&buf, args)
		assert.Empty(t, buf.String(), "Buffer should be empty")
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
	t.Run("non-utf8 strings", func(t *testing.T) {
		args := toIFaceSlice([]string{"Go", "\xc0\x80"})
		var buf bytes.Buffer
		err := writeInterfaces(&buf, args)
		assert.Exactly(t, `('Go',)`, buf.String())
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
	t.Run("non-utf8 NullStrings", func(t *testing.T) {
		args := toIFaceSlice([]null.String{null.MakeString("Go2"), null.MakeString("Hello\xc0\x80World")})
		var buf bytes.Buffer
		err := writeInterfaces(&buf, args)
		assert.Exactly(t, "('Go2',)", buf.String())
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
	t.Run("non-utf8 NullString", func(t *testing.T) {
		args := toIFaceSlice(null.MakeString("Hello\xc0\x80World"))
		var buf bytes.Buffer
		err := writeInterfaces(&buf, args)
		assert.Empty(t, buf.String())
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
	t.Run("bytes as binary", func(t *testing.T) {
		args := toIFaceSlice([][]byte{[]byte("\xc0\x80")})
		var buf bytes.Buffer
		assert.NoError(t, writeInterfaces(&buf, args))
		assert.Exactly(t, `(0xc080)`, buf.String())
	})
	t.Run("bytesSlice as binary", func(t *testing.T) {
		args := toIFaceSlice([][]byte{[]byte(`Rusty`), []byte("Go\xc0\x80")})
		var buf bytes.Buffer
		assert.NoError(t, writeInterfaces(&buf, args))
		assert.Exactly(t, "('Rusty',0x476fc080)", buf.String())
	})
	t.Run("should panic because unknown field type", func(t *testing.T) {
		var buf bytes.Buffer
		assert.ErrorIsKind(t, errors.NotSupported, writeInterfaceValue(complex64(1), &buf, 0))
		assert.Empty(t, buf.String(), "buffer should be empty")
	})
}

func TestArguments_HasNamedArgs(t *testing.T) {
	// TODO fix test resp. hasNamedArgs
	t.Run("hasNamedArgs in expression", func(t *testing.T) {
		p := &dmlPerson{
			Name: "a'bc",
		}

		a := NewSelect().
			AddColumnsConditions(
				Expr("?").Alias("n").Int64(1),
				Expr("CAST(:name AS CHAR(20))").Alias("str"),
			).WithDBR(dbMock{}).TestWithArgs(Qualify("", p))
		_, _, err := a.ToSQL()
		assert.NoError(t, err)
		// assert.Exactly(t, uint8(2), a.hasNamedArgs)
	})
	t.Run("hasNamedArgs in condition, no args", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().PlaceHolder(),
			Column("email").Like().NamedArg("ema1l")).WithDBR(dbMock{})
		_, _, err := a.ToSQL()
		assert.NoError(t, err)
		// assert.Exactly(t, uint8(0), a.hasNamedArgs)
	})
	t.Run("hasNamedArgs in condition, with args", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().PlaceHolder(),
			Column("email").Like().NamedArg("ema1l")).WithDBR(dbMock{}).TestWithArgs("my@email.org")
		_, _, err := a.ToSQL()
		assert.NoError(t, err)
		// assert.Exactly(t, uint8(1), a.hasNamedArgs)
	})
	t.Run("hasNamedArgs none", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().Int(221),
			Column("email").Like().Str("em@1l.de")).WithDBR(dbMock{})
		_, _, err := a.ToSQL()
		assert.NoError(t, err)
		// assert.Exactly(t, uint8(0), a.hasNamedArgs)
	})
}

func TestArguments_NextUnnamedArg(t *testing.T) {
	t.Run("three occurrences", func(t *testing.T) {
		args := toIFaceSlice(sql.Named("colZ", int64(3)), uint64(6), sql.Named("colB", 2.2), "c", sql.Named("colA", []string{"a", "b"}))

		dbr := &DBR{}
		var nextUnnamedArgPos int
		a, nextUnnamedArgPos, ok := dbr.nextUnnamedArg(nextUnnamedArgPos, args)
		assert.True(t, ok, "Should find an unnamed argument")
		assert.Exactly(t, uint64(6), a)

		a, nextUnnamedArgPos, ok = dbr.nextUnnamedArg(nextUnnamedArgPos, args)
		assert.True(t, ok, "Should find an unnamed argument")
		assert.Exactly(t, "c", a)

		a, nextUnnamedArgPos, ok = dbr.nextUnnamedArg(nextUnnamedArgPos, args)
		assert.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, nil, a)

		args = args[:0]
		args = append(args, 3.14159, sql.Named("price", 2.7182), now())
		nextUnnamedArgPos = 0
		a, nextUnnamedArgPos, ok = dbr.nextUnnamedArg(nextUnnamedArgPos, args)
		assert.True(t, ok, "Should find an unnamed argument")
		assert.Exactly(t, 3.14159, a)

		a, nextUnnamedArgPos, ok = dbr.nextUnnamedArg(nextUnnamedArgPos, args)
		assert.True(t, ok, "Should find an unnamed argument")
		assert.Exactly(t, now(), a)

		a, _, ok = dbr.nextUnnamedArg(nextUnnamedArgPos, args)
		assert.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, nil, a)
	})

	t.Run("zero occurrences", func(t *testing.T) {
		args := toIFaceSlice(sql.Named("colZ", int64(3)), sql.Named("colB", 2.2), sql.Named("colA", []string{"a", "b"}))
		dbr := &DBR{}

		a, _, ok := dbr.nextUnnamedArg(0, args)
		assert.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, nil, a)

		a, _, ok = dbr.nextUnnamedArg(0, args)
		assert.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, nil, a)
	})
}

func TestDBR_OrderByLimit(t *testing.T) {
	t.Run("WithoutArgs", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().Int(221),
			Column("email").Like().Str("em@1l.de")).WithDBR(dbMock{}).Limit(44, 55)

		t.Run("ASC", func(t *testing.T) {
			a.OrderBy("email", "id")
			compareToSQL2(t, a, errors.NoKind,
				"SELECT `a`, `b` FROM `c` WHERE (`id` > 221) AND (`email` LIKE 'em@1l.de') ORDER BY `email`, `id` LIMIT 44,55",
			)
		})
		t.Run("DESC", func(t *testing.T) {
			a.cachedSQL.OrderBys = a.cachedSQL.OrderBys[:1]
			a.OrderByDesc("firstname")
			compareToSQL2(t, a, errors.NoKind,
				"SELECT `a`, `b` FROM `c` WHERE (`id` > 221) AND (`email` LIKE 'em@1l.de') ORDER BY `email`, `firstname` DESC LIMIT 44,55",
			)
		})
	})

	t.Run("WithDBR", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().PlaceHolder(),
			Column("email").Like().Str("em@1l.de")).WithDBR(dbMock{}).Limit(44, 55)

		t.Run("ASC", func(t *testing.T) {
			a.OrderBy("email", "id")
			compareToSQL2(t, a, errors.NoKind,
				"SELECT `a`, `b` FROM `c` WHERE (`id` > ?) AND (`email` LIKE 'em@1l.de') ORDER BY `email`, `id` LIMIT 44,55",
			)
		})
		t.Run("DESC", func(t *testing.T) {
			a.cachedSQL.OrderBys = a.cachedSQL.OrderBys[:1]
			a.OrderByDesc("firstname")
			compareToSQL2(t, a, errors.NoKind,
				"SELECT `a`, `b` FROM `c` WHERE (`id` > ?) AND (`email` LIKE 'em@1l.de') ORDER BY `email`, `firstname` DESC LIMIT 44,55",
			)
		})
	})
}

func TestDBR_PreGeneratedQueries(t *testing.T) {
	cp, err := NewConnPool()
	assert.NoError(t, err)

	sel := NewSelect("a", "b").From("c").Where(
		Column("id").Greater().PlaceHolder(),
		Column("email").Like().PlaceHolder(),
	)
	sel2 := sel.Clone()
	sel2.Wheres = sel2.Wheres[:1]
	sel2.Wheres[0] = Column("id").Less().PlaceHolder()

	err = cp.RegisterByQueryBuilder(map[string]QueryBuilder{
		"id_greater": sel,
		"id_less":    sel2,
	})
	assert.NoError(t, err)

	// modify SQL

	compareToSQL2(t, cp.WithCacheKey("id_greater"), errors.NoKind,
		"SELECT `a`, `b` FROM `c` WHERE (`id` > ?) AND (`email` LIKE ?)",
	)
	compareToSQL2(t, cp.WithCacheKey("id_less"), errors.NoKind,
		"SELECT `a`, `b` FROM `c` WHERE (`id` < ?)",
	)
	compareToSQL2(t, cp.WithCacheKey("id_not_found"), errors.NotFound, "")
}

func TestExecValidateOneAffectedRow(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m := mockSQLRes{int64: 1}
		err := ExecValidateOneAffectedRow(m, nil)
		assert.NoError(t, err)
	})
	t.Run("RowsAffected fails", func(t *testing.T) {
		m := mockSQLRes{int64: 1, error: errors.ConnectionFailed.Newf("ups")}
		err := ExecValidateOneAffectedRow(m, nil)
		assert.ErrorIsKind(t, errors.ConnectionFailed, err)
	})
	t.Run("mismatch", func(t *testing.T) {
		m := mockSQLRes{int64: 2}
		err := ExecValidateOneAffectedRow(m, nil)
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
	t.Run("first error", func(t *testing.T) {
		m := mockSQLRes{int64: 2}
		err := ExecValidateOneAffectedRow(m, errors.AlreadyInUse.Newf("uppps"))
		assert.ErrorIsKind(t, errors.AlreadyInUse, err)
	})
}

type mockSQLRes struct {
	int64
	error
}

func (mockSQLRes) LastInsertId() (int64, error) {
	panic("implement me")
}

func (m mockSQLRes) RowsAffected() (int64, error) {
	return m.int64, m.error
}
