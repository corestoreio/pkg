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
	"fmt"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

var _ ColumnMapper = (*Artisan)(nil)
var _ fmt.GoStringer = (*Artisan)(nil)
var _ fmt.GoStringer = (*argument)(nil)

func TestArguments_Interfaces(t *testing.T) {
	t.Parallel()

	container := make([]interface{}, 0, 48)

	t.Run("no slices, nulls valid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Int(-1).Int64(1).Uint64(2).Float64(3.1).Bool(true).String("eCom1").Bytes([]byte(`eCom2`)).Time(now()).
			NullString(null.MakeString("eCom3")).NullInt64(null.MakeInt64(4)).NullFloat64(null.MakeFloat64(2.7)).
			NullBool(null.MakeBool(true)).NullTime(null.MakeTime(now()))

		assert.Exactly(t,
			[]interface{}{
				nil, int64(-1), int64(1), int64(2), 3.1, true, "eCom1", []uint8{0x65, 0x43, 0x6f, 0x6d, 0x32}, now(),
				"eCom3", int64(4), 2.7, true, now(),
			},
			args.Interfaces(container...))
		container = container[:0]
	})
	t.Run("no slices, nulls invalid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Int(-1).Int64(1).Uint64(2).Float64(3.1).Bool(true).String("eCom1").Bytes([]byte(`eCom2`)).Time(now()).
			NullString(null.String{}).NullInt64(null.Int64{}).NullFloat64(null.Float64{}).
			NullBool(null.Bool{}).NullTime(null.Time{})
		assert.Exactly(t,
			[]interface{}{nil, int64(-1), int64(1), int64(2), 3.1, true, "eCom1", []uint8{0x65, 0x43, 0x6f, 0x6d, 0x32}, now(),
				nil, nil, nil, nil, nil},
			args.Interfaces(container...))
		container = container[:0]
	})
	t.Run("slices, nulls valid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Ints(-1, -2).Int64s(1, 2).Uints(568, 766).Uint64s(2).Float64s(1.2, 3.1).Bools(false, true).
			Strings("eCom1", "eCom11").BytesSlice([]byte(`eCom2`)).Times(now(), now()).
			NullStrings(null.MakeString("eCom3"), null.MakeString("eCom3")).NullInt64s(null.MakeInt64(4), null.MakeInt64(4)).
			NullFloat64s(null.MakeFloat64(2.7), null.MakeFloat64(2.7)).
			NullBools(null.MakeBool(true)).NullTimes(null.MakeTime(now()), null.MakeTime(now()))
		assert.Exactly(t,
			[]interface{}{nil, int64(-1), int64(-2), int64(1), int64(2), int64(568), int64(766), int64(2), 1.2, 3.1, false, true,
				"eCom1", "eCom11", []uint8{0x65, 0x43, 0x6f, 0x6d, 0x32}, now(), now(),
				"eCom3", "eCom3", int64(4), int64(4),
				2.7, 2.7,
				true, now(), now()},
			args.Interfaces())
	})
	t.Run("returns nil interface", func(t *testing.T) {
		args := MakeArgs(10)
		assert.Nil(t, args.Interfaces(), "args.Interfaces() must return nil")
	})
}

func TestArguments_DriverValue(t *testing.T) {
	t.Parallel()

	t.Run("Driver.Values supported types", func(t *testing.T) {
		args := MakeArgs(10).
			DriverValues(
				driverValueNil(0),
				driverValueBytes(nil), null.MakeInt64(3), null.MakeFloat64(2.7), null.MakeBool(true),
				driverValueBytes(`Invoice`), null.MakeString("Creditmemo"), nowSentinel{}, null.MakeTime(now()),
			)
		assert.Exactly(t,
			[]interface{}{nil, []uint8(nil), int64(3), 2.7, true,
				[]uint8{0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65}, "Creditmemo", "2006-01-02 19:04:05", now()},
			args.Interfaces())
	})

	t.Run("Driver.Value supported types", func(t *testing.T) {
		args := MakeArgs(10).
			DriverValue(driverValueNil(0)).
			DriverValue(driverValueBytes(nil)).
			DriverValue(null.MakeInt64(3)).
			DriverValue(null.MakeFloat64(2.7)).
			DriverValue(null.MakeBool(true)).
			DriverValue(driverValueBytes(`Invoice`)).
			DriverValue(null.MakeString("Creditmemo")).
			DriverValue(nowSentinel{}).
			DriverValue(null.MakeTime(now()))

		assert.Exactly(t,
			[]interface{}{nil, []uint8(nil), int64(3), 2.7, true,
				[]uint8{0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65}, "Creditmemo", "2006-01-02 19:04:05", now()},
			args.Interfaces())
	})

	t.Run("Driver.Values panics because not supported", func(t *testing.T) {
		_, _, err := MakeArgs(10).
			DriverValue(
				driverValueNotSupported(4),
			).ToSQL()
		assert.True(t, errors.Is(err, errors.NotSupported), "Should have behaviour errors.NotSupported")
	})

	t.Run("Driver.Values panics because Value error", func(t *testing.T) {
		_, _, err := MakeArgs(10).
			DriverValue(
				driverValueError(0),
			).ToSQL()
		assert.True(t, errors.Is(err, errors.Fatal), "Should have behaviour errors.Fatal")
	})
}

func TestArguments_WriteTo(t *testing.T) {
	t.Parallel()

	t.Run("no slices, nulls valid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Int(-1).Int64(1).Uint64(2).Float64(3.1).Bool(true).String("eCom1").Bytes([]byte(`eCom2`)).Time(now()).
			NullString(null.MakeString("eCom3")).NullInt64(null.MakeInt64(4)).NullFloat64(null.MakeFloat64(2.7)).
			NullBool(null.MakeBool(true)).NullTime(null.MakeTime(now()))

		buf := new(bytes.Buffer)
		err := args.Write(buf)
		assert.NoError(t, err)
		assert.Exactly(t,
			"(NULL,-1,1,2,3.1,1,'eCom1','eCom2','2006-01-02 15:04:05','eCom3',4,2.7,1,'2006-01-02 15:04:05')",
			buf.String())
	})
	t.Run("no slices, nulls invalid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Int(-1).Int64(1).Uint64(2).Float64(3.1).Bool(true).String("eCom1").Bytes([]byte(`eCom2`)).Time(now()).
			NullString(null.String{}).NullInt64(null.Int64{}).NullFloat64(null.Float64{}).
			NullBool(null.Bool{}).NullTime(null.Time{})

		buf := new(bytes.Buffer)
		err := args.Write(buf)
		assert.NoError(t, err)
		assert.Exactly(t,
			"(NULL,-1,1,2,3.1,1,'eCom1','eCom2','2006-01-02 15:04:05',NULL,NULL,NULL,NULL,NULL)",
			buf.String())
	})
	t.Run("slices, nulls valid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Ints(-1, -2).Int64s(1, 2).Uint64s(2).Float64s(1.2, 3.1).Bools(false, true).Strings("eCom1", "eCom11").BytesSlice([]byte(`eCom2`)).Times(now(), now()).
			NullStrings(null.MakeString("eCom3"), null.MakeString("eCom3")).NullInt64s(null.MakeInt64(4), null.MakeInt64(5)).NullFloat64s(null.MakeFloat64(2.71), null.MakeFloat64(2.72)).
			NullBools(null.MakeBool(true)).NullTimes(null.MakeTime(now()), null.MakeTime(now()))

		buf := new(bytes.Buffer)
		err := args.Write(buf)
		assert.NoError(t, err)
		assert.Exactly(t,
			"(NULL,(-1,-2),(1,2),(2),(1.2,3.1),(0,1),('eCom1','eCom11'),('eCom2'),('2006-01-02 15:04:05','2006-01-02 15:04:05'),('eCom3','eCom3'),(4,5),(2.71,2.72),(1),('2006-01-02 15:04:05','2006-01-02 15:04:05'))",
			buf.String())
	})
	t.Run("non-utf8 string", func(t *testing.T) {
		args := MakeArgs(2).String("\xc0\x80")
		buf := new(bytes.Buffer)
		err := args.Write(buf)
		assert.Empty(t, buf.String(), "Buffer should be empty")
		assert.True(t, errors.NotValid.Match(err), "Should have a not valid error behaviour %+v", err)
	})
	t.Run("non-utf8 strings", func(t *testing.T) {
		args := MakeArgs(2).Strings("Go", "\xc0\x80")
		buf := new(bytes.Buffer)
		err := args.Write(buf)
		assert.Exactly(t, `('Go',)`, buf.String())
		assert.True(t, errors.NotValid.Match(err), "Should have a not valid error behaviour %+v", err)
	})
	t.Run("non-utf8 NullStrings", func(t *testing.T) {
		args := MakeArgs(2).NullStrings(null.MakeString("Go2"), null.MakeString("Hello\xc0\x80World"))
		buf := new(bytes.Buffer)
		err := args.Write(buf)
		assert.Exactly(t, "('Go2',)", buf.String())
		assert.True(t, errors.NotValid.Match(err), "Should have a not valid error behaviour %+v", err)
	})
	t.Run("non-utf8 NullString", func(t *testing.T) {
		args := MakeArgs(2).NullString(null.MakeString("Hello\xc0\x80World"))
		buf := new(bytes.Buffer)
		err := args.Write(buf)
		assert.Empty(t, buf.String())
		assert.True(t, errors.NotValid.Match(err), "Should have a not valid error behaviour %+v", err)
	})
	t.Run("bytes as binary", func(t *testing.T) {
		args := MakeArgs(2).Bytes([]byte("\xc0\x80"))
		buf := new(bytes.Buffer)
		assert.NoError(t, args.Write(buf))
		assert.Exactly(t, "0xc080", buf.String())
	})
	t.Run("bytesSlice as binary", func(t *testing.T) {
		args := MakeArgs(2).BytesSlice([]byte(`Rusty`), []byte("Go\xc0\x80"))
		buf := new(bytes.Buffer)
		assert.NoError(t, args.Write(buf))
		assert.Exactly(t, "('Rusty',0x476fc080)", buf.String())
	})
	t.Run("should panic because unknown field type", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.True(t, errors.NotSupported.Match(err), "Should be a not supported error; got %+v", err)
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

func TestArguments_Named(t *testing.T) {
	t.Parallel()

	assert.Exactly(t,
		"dml.MakeArgs(4).String(\"Rusty\").Name(\"entity_id\").Null().Name(\"entity_sku\").Int64(4).Float64(3.141000)",
		MakeArgs(2).
			String("Rusty").
			Name("entity_id").Name("entity_sku").Int64(4).
			Float64(3.141).
			GoString())

	assert.Exactly(t,
		"dml.MakeArgs(3).Name(\"entity_id\").Null().Name(\"entity_sku\").Int64(4).Float64(3.141000)",
		MakeArgs(2).
			Name("entity_id").Name("entity_sku").Int64(4).
			Float64(3.141).
			GoString())

	assert.Exactly(t,
		"dml.MakeArgs(2).Name(\"entity_id\").Int64(4).Float64(3.141000)",
		MakeArgs(2).
			Name("entity_id").Int64(4).
			Float64(3.141).
			GoString())

	assert.Exactly(t,
		"dml.MakeArgs(2).Float64(3.141000).Name(\"entity_id\").Int64(4)",
		MakeArgs(2).
			Float64(3.141).
			Name("entity_id").Int64(4).
			GoString())

	assert.Exactly(t,
		"dml.MakeArgs(4).Float64s([]float64{2.76, 3.141}...).Name(\"entity_id\").Int64(4).Name(\"store_id\").Uint64(5678).Time(time.Unix(1136228645,2))",
		MakeArgs(2).
			Float64s(2.76, 3.141).
			Name("entity_id").Int64(4).
			Name("store_id").Uint64(5678).
			Time(now()).
			GoString())
}

func TestArguments_HasNamedArgs(t *testing.T) {
	t.Parallel()

	t.Run("hasNamedArgs in expression", func(t *testing.T) {
		a := NewSelect().
			AddColumnsConditions(
				Expr("?").Alias("n").Int64(1),
				Expr("CAST(:abc AS CHAR(20))").Alias("str"),
			).WithArgs().Record("", MakeArgs(1).Name("abc").String("a'bc"))
		_, _, err := a.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, uint8(2), a.hasNamedArgs)
	})
	t.Run("hasNamedArgs in condition, no args", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().PlaceHolder(),
			Column("email").Like().NamedArg("ema1l")).WithArgs()
		_, _, err := a.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, uint8(1), a.hasNamedArgs)
	})
	t.Run("hasNamedArgs in condition, with args", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().PlaceHolder(),
			Column("email").Like().NamedArg("ema1l")).WithArgs().String("my@email.org")
		_, _, err := a.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, uint8(1), a.hasNamedArgs)
	})
	t.Run("hasNamedArgs none", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().Int(221),
			Column("email").Like().Str("em@1l.de")).WithArgs()
		_, _, err := a.ToSQL()
		assert.NoError(t, err)
		assert.Exactly(t, uint8(1), a.hasNamedArgs)
	})
}

func TestArguments_MapColumns(t *testing.T) {
	t.Parallel()

	from := MakeArgs(3)

	t.Run("len=1", func(t *testing.T) {
		from = from.Reset().Int64(3).Float64(2.2).Name("colA").Strings("a", "b")
		cm := NewColumnMap(1, "colA")
		if err := from.MapColumns(cm); err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t,
			"dml.MakeArgs(1).Name(\"colA\").Strings(\"a\",\"b\")",
			cm.arguments.GoString())
	})

	t.Run("len=0", func(t *testing.T) {
		from = from.Reset().Name("colZ").Int64(3).Float64(2.2).Name("colA").Strings("a", "b")
		cm := NewColumnMap(1)
		if err := from.MapColumns(cm); err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t,
			"dml.MakeArgs(3).Name(\"colZ\").Int64(3).Float64(2.200000).Name(\"colA\").Strings(\"a\",\"b\")",
			cm.arguments.GoString())
	})

	t.Run("len>1", func(t *testing.T) {
		from = from.Reset().Name("colZ").Int64(3).Uint64(6).Name("colB").Float64(2.2).String("c").Name("colA").Strings("a", "b")
		cm := NewColumnMap(1, "colA", "colB")
		if err := from.MapColumns(cm); err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t,
			"dml.MakeArgs(2).Name(\"colA\").Strings(\"a\",\"b\").Name(\"colB\").Float64(2.200000)",
			cm.arguments.GoString())
	})
}

func TestArguments_NextUnnamedArg(t *testing.T) {
	t.Parallel()

	t.Run("three occurrences", func(t *testing.T) {
		args := MakeArgs(5).Name("colZ").Int64(3).Uint64(6).Name("colB").Float64(2.2).String("c").Name("colA").Strings("a", "b")

		a, ok := args.nextUnnamedArg()
		assert.True(t, ok, "Should find an unnamed argument")
		assert.Empty(t, a.name)
		assert.Exactly(t, uint64(6), a.value)

		a, ok = args.nextUnnamedArg()
		assert.True(t, ok, "Should find an unnamed argument")
		assert.Empty(t, a.name)
		assert.Exactly(t, "c", a.value)

		a, ok = args.nextUnnamedArg()
		assert.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, argument{}, a)

		args.Reset().Float64(3.14159).Name("price").Float64(2.7182).Time(now())

		a, ok = args.nextUnnamedArg()
		assert.True(t, ok, "Should find an unnamed argument")
		assert.Empty(t, a.name)
		assert.Exactly(t, 3.14159, a.value)

		a, ok = args.nextUnnamedArg()
		assert.True(t, ok, "Should find an unnamed argument")
		assert.Empty(t, a.name)
		assert.Exactly(t, now(), a.value)

		a, ok = args.nextUnnamedArg()
		assert.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, argument{}, a)
	})

	t.Run("zero occurrences", func(t *testing.T) {
		args := MakeArgs(5).Name("colZ").Int64(3).Name("colB").Float64(2.2).Name("colA").Strings("a", "b")

		a, ok := args.nextUnnamedArg()
		assert.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, argument{}, a)

		a, ok = args.nextUnnamedArg()
		assert.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, argument{}, a)
	})
}

func TestArtisan_Clone(t *testing.T) {
	t.Parallel()
	sel := NewSelect("a", "b").From("c").WithDB(dbMock{})
	sel.qualifiedColumns = []string{"x", "y"}
	selA := sel.WithArgs()
	selA.base.Log = log.BlackHole{}

	selB := selA.Clone()
	assert.Nil(t, selB.base.DB)
	assert.Exactly(t, selA.base.Log, selB.base.Log)
	assert.Exactly(t, selA.base.cachedSQL, selB.base.cachedSQL)

	assert.Exactly(t, selA.QualifiedColumnsAliases, selB.QualifiedColumnsAliases)
}

func TestArtisan_OrderByLimit(t *testing.T) {
	t.Parallel()

	t.Run("WithoutArgs", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().Int(221),
			Column("email").Like().Str("em@1l.de")).WithArgs().Limit(44, 55)

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

	t.Run("WithArgs", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().PlaceHolder(),
			Column("email").Like().Str("em@1l.de")).WithArgs().Int(87653).Limit(44, 55)

		t.Run("ASC", func(t *testing.T) {
			a.OrderBy("email", "id")
			compareToSQL2(t, a, errors.NoKind,
				"SELECT `a`, `b` FROM `c` WHERE (`id` > ?) AND (`email` LIKE 'em@1l.de') ORDER BY `email`, `id` LIMIT 44,55",
				int64(87653),
			)
		})
		t.Run("DESC", func(t *testing.T) {
			a.OrderBys = a.OrderBys[:1]
			a.OrderByDesc("firstname")
			compareToSQL2(t, a, errors.NoKind,
				"SELECT `a`, `b` FROM `c` WHERE (`id` > ?) AND (`email` LIKE 'em@1l.de') ORDER BY `email`, `firstname` DESC LIMIT 44,55",
				int64(87653),
			)
		})
	})
}
