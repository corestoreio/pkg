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
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

func TestArguments_Interfaces(t *testing.T) {
	t.Parallel()

	container := make([]interface{}, 0, 48)

	t.Run("no slices, nulls valid", func(t *testing.T) {
		args := makeArguments10().
			add(nil, -1, int64(1), uint64(2), 3.1, true, "eCom1", []byte(`eCom2`), now()).
			add(null.MakeString("eCom3"), null.MakeInt64(4), null.MakeFloat64(2.7)).
			add(null.MakeBool(true), null.MakeTime(now()))

		assert.Exactly(t,
			[]interface{}{
				nil, int64(-1), int64(1), int64(2), 3.1, true, "eCom1", []uint8{0x65, 0x43, 0x6f, 0x6d, 0x32}, now(),
				"eCom3", int64(4), 2.7, true, now(),
			},
			args.toInterfaces(container))
		container = container[:0]
	})
	t.Run("no slices, nulls invalid", func(t *testing.T) {
		args := makeArguments10().
			add(nil, -1, int64(1), uint64(2), 3.1, true, "eCom1", []byte(`eCom2`), now()).
			add(null.String{}, null.Int64{}, null.Float64{}).
			add(null.Bool{}, null.Time{})
		assert.Exactly(t,
			[]interface{}{
				nil, int64(-1), int64(1), int64(2), 3.1, true, "eCom1", []uint8{0x65, 0x43, 0x6f, 0x6d, 0x32}, now(),
				nil, nil, nil, nil, nil,
			},
			args.toInterfaces(container))
		container = container[:0]
	})
	t.Run("slices, nulls valid", func(t *testing.T) {
		args := makeArguments10().
			add(nil, []int{-1, -2}, []int64{1, 2}, []uint{568, 766}, []uint64{2}, []float64{1.2, 3.1}, []bool{false, true}).
			add([]string{"eCom1", "eCom11"}, [][]byte{[]byte(`eCom2`)}, []time.Time{now(), now()}).
			add([]null.String{null.MakeString("eCom3"), null.MakeString("eCom3")}, []null.Int64{null.MakeInt64(4), null.MakeInt64(4)}).
			add([]null.Float64{null.MakeFloat64(2.7), null.MakeFloat64(2.7)}).
			add([]null.Bool{null.MakeBool(true)}, []null.Time{null.MakeTime(now()), null.MakeTime(now())})
		assert.Exactly(t,
			[]interface{}{
				nil, int64(-1), int64(-2), int64(1), int64(2), int64(568), int64(766), int64(2), 1.2, 3.1, false, true,
				"eCom1", "eCom11", []uint8{0x65, 0x43, 0x6f, 0x6d, 0x32}, now(), now(),
				"eCom3", "eCom3", int64(4), int64(4),
				2.7, 2.7,
				true, now(), now(),
			},
			args.toInterfaces(nil))
	})
	t.Run("returns nil interface", func(t *testing.T) {
		args := makeArguments10()
		assert.Nil(t, args.toInterfaces(nil), "args.toInterfaces() must return nil")
	})
}

func TestArguments_DriverValue(t *testing.T) {
	t.Parallel()

	t.Run("Driver.Values supported types", func(t *testing.T) {
		args := makeArguments10().
			add(
				driverValueNil(0),
				driverValueBytes(nil), null.MakeInt64(3), null.MakeFloat64(2.7), null.MakeBool(true),
				driverValueBytes(`Invoice`), null.MakeString("Creditmemo"), nowSentinel{}, null.MakeTime(now()),
			)
		assert.Exactly(t,
			[]interface{}{
				nil, []uint8(nil), int64(3), 2.7, true,
				[]uint8{0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65}, "Creditmemo", "2006-01-02 19:04:05", now(),
			},
			args.toInterfaces(nil))
	})

	t.Run("Driver.Value supported types", func(t *testing.T) {
		args := makeArguments10().
			add(driverValueNil(0)).
			add(driverValueBytes(nil)).
			add(null.MakeInt64(3)).
			add(null.MakeFloat64(2.7)).
			add(null.MakeBool(true)).
			add(driverValueBytes(`Invoice`)).
			add(null.MakeString("Creditmemo")).
			add(nowSentinel{}).
			add(null.MakeTime(now()))

		assert.Exactly(t,
			[]interface{}{
				nil, []uint8(nil), int64(3), 2.7, true,
				[]uint8{0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65}, "Creditmemo", "2006-01-02 19:04:05", now(),
			},
			args.toInterfaces(nil))
	})

	t.Run("Driver.Values panics because not supported", func(t *testing.T) {
		_, err := driverValue(
			arguments{},
			driverValueNotSupported(4),
		)
		assert.ErrorIsKind(t, errors.NotSupported, err)
	})

	t.Run("Driver.Values panics because Value error", func(t *testing.T) {
		_, err := driverValue(
			arguments{},
			driverValueError(0),
		)
		assert.ErrorIsKind(t, errors.Fatal, err)
	})
}

func TestArguments_WriteTo(t *testing.T) {
	t.Parallel()

	t.Run("no slices, nulls valid", func(t *testing.T) {
		args := makeArguments10().
			add(nil, -1, int64(1), uint64(2), 3.1, true, "eCom1", []byte(`eCom2`), now()).
			add(null.MakeString("eCom3"), null.MakeInt64(4), null.MakeFloat64(2.7)).
			add(null.MakeBool(true), null.MakeTime(now()))

		buf := new(bytes.Buffer)
		err := args.write(buf)
		assert.NoError(t, err)
		assert.Exactly(t,
			"(NULL,-1,1,2,3.1,1,'eCom1','eCom2','2006-01-02 15:04:05','eCom3',4,2.7,1,'2006-01-02 15:04:05')",
			buf.String())
	})
	t.Run("no slices, nulls invalid", func(t *testing.T) {
		args := makeArguments10().
			add(nil, -1, int64(1), uint64(2), 3.1, true, "eCom1", []byte(`eCom2`), now()).
			add(null.String{}, null.Int64{}, null.Float64{}).
			add(null.Bool{}, null.Time{})

		buf := new(bytes.Buffer)
		err := args.write(buf)
		assert.NoError(t, err)
		assert.Exactly(t,
			"(NULL,-1,1,2,3.1,1,'eCom1','eCom2','2006-01-02 15:04:05',NULL,NULL,NULL,NULL,NULL)",
			buf.String())
	})
	t.Run("slices, nulls valid", func(t *testing.T) {
		args := makeArguments10().
			add(nil, []int{-1, -2}, []int64{1, 2}, []uint{568, 766}, []uint64{2}, []float64{1.2, 3.1}, []bool{false, true}).
			add([]string{"eCom1", "eCom11"}, [][]byte{[]byte(`eCom2`)}, []time.Time{now(), now()}).
			add([]null.String{null.MakeString("eCom3"), null.MakeString("eCom3")}, []null.Int64{null.MakeInt64(4), null.MakeInt64(4)}).
			add([]null.Float64{null.MakeFloat64(2.7), null.MakeFloat64(2.7)}).
			add([]null.Bool{null.MakeBool(true)}, []null.Time{null.MakeTime(now()), null.MakeTime(now())})

		buf := new(bytes.Buffer)
		err := args.write(buf)
		assert.NoError(t, err)
		assert.Exactly(t,
			"(NULL,(-1,-2),(1,2),(568,766),(2),(1.2,3.1),(0,1),('eCom1','eCom11'),('eCom2'),('2006-01-02 15:04:05','2006-01-02 15:04:05'),('eCom3','eCom3'),(4,4),(2.7,2.7),(1),('2006-01-02 15:04:05','2006-01-02 15:04:05'))",
			buf.String(), "%q", buf.String())
	})
	t.Run("non-utf8 string", func(t *testing.T) {
		args := makeArguments10().add("\xc0\x80")
		buf := new(bytes.Buffer)
		err := args.write(buf)
		assert.Empty(t, buf.String(), "Buffer should be empty")
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
	t.Run("non-utf8 strings", func(t *testing.T) {
		args := makeArguments10().add([]string{"Go", "\xc0\x80"})
		buf := new(bytes.Buffer)
		err := args.write(buf)
		assert.Exactly(t, `('Go',)`, buf.String())
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
	t.Run("non-utf8 NullStrings", func(t *testing.T) {
		args := makeArguments10().add([]null.String{null.MakeString("Go2"), null.MakeString("Hello\xc0\x80World")})
		buf := new(bytes.Buffer)
		err := args.write(buf)
		assert.Exactly(t, "('Go2',)", buf.String())
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
	t.Run("non-utf8 NullString", func(t *testing.T) {
		args := makeArguments10().add(null.MakeString("Hello\xc0\x80World"))
		buf := new(bytes.Buffer)
		err := args.write(buf)
		assert.Empty(t, buf.String())
		assert.ErrorIsKind(t, errors.NotValid, err)
	})
	t.Run("bytes as binary", func(t *testing.T) {
		args := makeArguments10().add([][]byte{[]byte("\xc0\x80")})
		buf := new(bytes.Buffer)
		assert.NoError(t, args.write(buf))
		assert.Exactly(t, `(0xc080)`, buf.String())
	})
	t.Run("bytesSlice as binary", func(t *testing.T) {
		args := makeArguments10().add([][]byte{[]byte(`Rusty`), []byte("Go\xc0\x80")})
		buf := new(bytes.Buffer)
		assert.NoError(t, args.write(buf))
		assert.Exactly(t, "('Rusty',0x476fc080)", buf.String())
	})
	t.Run("should panic because unknown field type", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.ErrorIsKind(t, errors.NotSupported, err)
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()

		au := argument{value: complex64(1), isSet: true}
		buf := new(bytes.Buffer)
		assert.NoError(t, au.writeTo(buf, 0))
		assert.Empty(t, buf.String(), "buffer should be empty")
	})
}

func TestArguments_HasNamedArgs(t *testing.T) {
	t.Parallel()
	// TODO fix test resp. hasNamedArgs
	t.Run("hasNamedArgs in expression", func(t *testing.T) {
		p := &dmlPerson{
			Name: "a'bc",
		}

		a := NewSelect().
			AddColumnsConditions(
				Expr("?").Alias("n").Int64(1),
				Expr("CAST(:name AS CHAR(20))").Alias("str"),
			).WithDBR().TestWithArgs(Qualify("", p))
		_, _, err := a.ToSQL()
		assert.NoError(t, err)
		// assert.Exactly(t, uint8(2), a.hasNamedArgs)
	})
	t.Run("hasNamedArgs in condition, no args", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().PlaceHolder(),
			Column("email").Like().NamedArg("ema1l")).WithDBR()
		_, _, err := a.ToSQL()
		assert.NoError(t, err)
		// assert.Exactly(t, uint8(0), a.hasNamedArgs)
	})
	t.Run("hasNamedArgs in condition, with args", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().PlaceHolder(),
			Column("email").Like().NamedArg("ema1l")).WithDBR().TestWithArgs("my@email.org")
		_, _, err := a.ToSQL()
		assert.NoError(t, err)
		// assert.Exactly(t, uint8(1), a.hasNamedArgs)
	})
	t.Run("hasNamedArgs none", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().Int(221),
			Column("email").Like().Str("em@1l.de")).WithDBR()
		_, _, err := a.ToSQL()
		assert.NoError(t, err)
		// assert.Exactly(t, uint8(0), a.hasNamedArgs)
	})
}

func TestArguments_NextUnnamedArg(t *testing.T) {
	t.Parallel()

	t.Run("three occurrences", func(t *testing.T) {
		args := makeArguments10().add(sql.Named("colZ", int64(3)), uint64(6), sql.Named("colB", 2.2), "c", sql.Named("colA", []string{"a", "b"}))

		dbr := &DBR{}
		var nextUnnamedArgPos int
		a, nextUnnamedArgPos, ok := dbr.nextUnnamedArg(nextUnnamedArgPos, args)
		assert.True(t, ok, "Should find an unnamed argument")
		assert.Exactly(t, uint64(6), a.value)

		a, nextUnnamedArgPos, ok = dbr.nextUnnamedArg(nextUnnamedArgPos, args)
		assert.True(t, ok, "Should find an unnamed argument")
		assert.Exactly(t, "c", a.value)

		a, nextUnnamedArgPos, ok = dbr.nextUnnamedArg(nextUnnamedArgPos, args)
		assert.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, argument{}, a)

		args = args[:0]
		args = args.add(3.14159, sql.Named("price", 2.7182), now())
		nextUnnamedArgPos = 0
		a, nextUnnamedArgPos, ok = dbr.nextUnnamedArg(nextUnnamedArgPos, args)
		assert.True(t, ok, "Should find an unnamed argument")
		assert.Exactly(t, 3.14159, a.value)

		a, nextUnnamedArgPos, ok = dbr.nextUnnamedArg(nextUnnamedArgPos, args)
		assert.True(t, ok, "Should find an unnamed argument")
		assert.Exactly(t, now(), a.value)

		a, _, ok = dbr.nextUnnamedArg(nextUnnamedArgPos, args)
		assert.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, argument{}, a)
	})

	t.Run("zero occurrences", func(t *testing.T) {
		args := makeArguments10().add(sql.Named("colZ", int64(3)), sql.Named("colB", 2.2), sql.Named("colA", []string{"a", "b"}))
		dbr := &DBR{}

		a, _, ok := dbr.nextUnnamedArg(0, args)
		assert.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, argument{}, a)

		a, _, ok = dbr.nextUnnamedArg(0, args)
		assert.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, argument{}, a)
	})
}

func TestDBR_Clone(t *testing.T) {
	t.Parallel()
	sel := NewSelect("a", "b").From("c").WithDB(dbMock{})
	sel.qualifiedColumns = []string{"x", "y"}
	selA := sel.WithDBR()
	selA.base.Log = log.BlackHole{}

	selB := selA.Clone()
	assert.NotNil(t, selB.base.DB)
	assert.Exactly(t, selA.base.Log, selB.base.Log)
	assert.Exactly(t, selA.base.cachedSQL, selB.base.cachedSQL)

	assert.Exactly(t, selA.QualifiedColumnsAliases, selB.QualifiedColumnsAliases)
}

func TestDBR_OrderByLimit(t *testing.T) {
	t.Parallel()

	t.Run("WithoutArgs", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().Int(221),
			Column("email").Like().Str("em@1l.de")).WithDBR().Limit(44, 55)

		t.Run("ASC", func(t *testing.T) {
			a.OrderBy("email", "id")
			compareToSQL2(t, a, errors.NoKind,
				"SELECT `a`, `b` FROM `c` WHERE (`id` > 221) AND (`email` LIKE 'em@1l.de') ORDER BY `email`, `id` LIMIT 44,55",
			)
		})
		t.Run("DESC", func(t *testing.T) {
			a.OrderBys = a.OrderBys[:1]
			a.OrderByDesc("firstname")
			compareToSQL2(t, a, errors.NoKind,
				"SELECT `a`, `b` FROM `c` WHERE (`id` > 221) AND (`email` LIKE 'em@1l.de') ORDER BY `email`, `firstname` DESC LIMIT 44,55",
			)
		})
	})

	t.Run("WithDBR", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().PlaceHolder(),
			Column("email").Like().Str("em@1l.de")).WithDBR().Limit(44, 55)

		t.Run("ASC", func(t *testing.T) {
			a.OrderBy("email", "id")
			compareToSQL2(t, a, errors.NoKind,
				"SELECT `a`, `b` FROM `c` WHERE (`id` > ?) AND (`email` LIKE 'em@1l.de') ORDER BY `email`, `id` LIMIT 44,55",
			)
		})
		t.Run("DESC", func(t *testing.T) {
			a.OrderBys = a.OrderBys[:1]
			a.OrderByDesc("firstname")
			compareToSQL2(t, a, errors.NoKind,
				"SELECT `a`, `b` FROM `c` WHERE (`id` > ?) AND (`email` LIKE 'em@1l.de') ORDER BY `email`, `firstname` DESC LIMIT 44,55",
			)
		})
	})
}

func TestDBR_PreGeneratedQueries(t *testing.T) {
	t.Parallel()

	t.Run("SELECT", func(t *testing.T) {
		sel := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().PlaceHolder(),
			Column("email").Like().PlaceHolder())

		_, _, err := sel.WithCacheKey("id_greater").ToSQL() // build cache 1
		assert.NoError(t, err)

		// modify SQL
		sel.Wheres[0] = Column("id").Less().PlaceHolder()
		_, _, err = sel.WithCacheKey("id_less").ToSQL() // build cache 2
		assert.NoError(t, err)

		selA := sel.WithDBR()

		compareToSQL2(t, selA.WithCacheKey("id_greater"), errors.NoKind,
			"SELECT `a`, `b` FROM `c` WHERE (`id` > ?) AND (`email` LIKE ?)",
		)
		compareToSQL2(t, selA.WithCacheKey("id_less"), errors.NoKind,
			"SELECT `a`, `b` FROM `c` WHERE (`id` < ?) AND (`email` LIKE ?)",
		)
		compareToSQL2(t, selA.WithCacheKey("id_not_found"), errors.Empty, "")
	})
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
