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
	"database/sql/driver"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ ColumnMapper = (*Arguments)(nil)
var _ fmt.GoStringer = (*Arguments)(nil)
var _ fmt.GoStringer = (*argument)(nil)

type driverValueBytes []byte

// Value implements the driver.Valuer interface.
func (a driverValueBytes) Value() (driver.Value, error) {
	return []byte(a), nil
}

type driverValueNotSupported uint8

// Value implements the driver.Valuer interface.
func (a driverValueNotSupported) Value() (driver.Value, error) {
	return uint8(a), nil
}

type driverValueNil uint8

// Value implements the driver.Valuer interface.
func (a driverValueNil) Value() (driver.Value, error) {
	return nil, nil
}

type driverValueError uint8

// Value implements the driver.Valuer interface.
func (a driverValueError) Value() (driver.Value, error) {
	return nil, errors.Aborted.Newf("WE've aborted something")
}

func TestArguments_Length_and_Stringer(t *testing.T) {
	t.Parallel()

	t.Run("no slices, nulls valid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Int(-1).Int64(1).Uint(9898).Uint64(2).Float64(3.1).Bool(true).String("eCom1").Bytes([]byte(`eCom2`)).Time(now()).
			NullString(MakeNullString("eCom3")).NullInt64(MakeNullInt64(4)).NullFloat64(MakeNullFloat64(2.7)).
			NullBool(MakeNullBool(true)).NullTime(MakeNullTime(now()))
		assert.Exactly(t, 15, args.Len(), "Length mismatch")

		// like fmt.GoStringer
		assert.Exactly(t,
			"dml.MakeArgs(15).Null().Int(-1).Int64(1).Uint64(9898).Uint64(2).Float64(3.100000).Bool(true).String(\"eCom1\").Bytes([]byte{0x65, 0x43, 0x6f, 0x6d, 0x32}).Time(time.Unix(1136228645,2)).NullString(dml.MakeNullString(`eCom3`)).NullInt64(dml.MakeNullInt64(4)).NullFloat64(dml.MakeNullFloat64(2.7)).NullBool(dml.MakeNullBool(true)).NullTime(dml.MakeNullTime(time.Unix(1136228645,2))",
			fmt.Sprintf("%#v", args))
	})

	t.Run("no slices, nulls invalid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Int(-1).Int64(1).Uint64(2).Float64(3.1).Bool(true).String("eCom1").Bytes([]byte(`eCom2`)).Time(now()).
			NullString(MakeNullString("eCom3", false)).NullInt64(MakeNullInt64(4, false)).NullFloat64(MakeNullFloat64(2.7, false)).
			NullBool(MakeNullBool(true, false)).NullTime(MakeNullTime(now(), false))
		assert.Exactly(t, 14, args.Len(), "Length mismatch")
		assert.Exactly(t,
			"dml.MakeArgs(14).Null().Int(-1).Int64(1).Uint64(2).Float64(3.100000).Bool(true).String(\"eCom1\").Bytes([]byte{0x65, 0x43, 0x6f, 0x6d, 0x32}).Time(time.Unix(1136228645,2)).NullString(dml.NullString{}).NullInt64(dml.NullInt64{}).NullFloat64(dml.NullFloat64{}).NullBool(dml.NullBool{}).NullTime(dml.NullTime{})",
			fmt.Sprintf("%#v", args))
	})

	t.Run("slices, nulls valid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Int(-1).Int64s(1, 2).Uints(567, 765).Uint64s(2).Float64s(1.2, 3.1).Bools(false, true).Strings("eCom1", "eCom11").BytesSlice(nil, []byte(`eCom2`)).Times(now(), now()).
			NullStrings(MakeNullString("eCom3"), MakeNullString("eCom3")).NullInt64s(MakeNullInt64(4), MakeNullInt64(4)).NullFloat64s(MakeNullFloat64(2.7), MakeNullFloat64(2.7)).
			NullBools(MakeNullBool(true)).NullTimes(MakeNullTime(now()), MakeNullTime(now()))
		assert.Exactly(t, 26, args.Len(), "Length mismatch")
		assert.Exactly(t,
			"dml.MakeArgs(15).Null().Int(-1).Int64s([]int64{1, 2}...).Uints([]uint{0x237, 0x2fd}...).Uint64s([]uint64{0x2}...).Float64s([]float64{1.2, 3.1}...).Bools([]bool{false, true}...).Strings(\"eCom1\",\"eCom11\").BytesSlice([]byte(nil),[]byte{0x65, 0x43, 0x6f, 0x6d, 0x32}).Times(time.Unix(1136228645,2),time.Unix(1136228645,2)).NullStrings(dml.MakeNullString(`eCom3`),dml.MakeNullString(`eCom3`)).NullInt64s(dml.MakeNullInt64(4),dml.MakeNullInt64(4)).NullFloat64s(dml.MakeNullFloat64(2.7),dml.MakeNullFloat64(2.7)).NullBools(dml.MakeNullBool(true)).NullTimes(dml.MakeNullTime(time.Unix(1136228645,2),dml.MakeNullTime(time.Unix(1136228645,2))",
			fmt.Sprintf("%#v", args))
	})
}

func TestArguments_Interfaces(t *testing.T) {
	t.Parallel()

	container := make([]interface{}, 0, 48)

	t.Run("no slices, nulls valid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Int(-1).Int64(1).Uint64(2).Float64(3.1).Bool(true).String("eCom1").Bytes([]byte(`eCom2`)).Time(now()).
			NullString(MakeNullString("eCom3")).NullInt64(MakeNullInt64(4)).NullFloat64(MakeNullFloat64(2.7)).
			NullBool(MakeNullBool(true)).NullTime(MakeNullTime(now()))

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
			NullString(MakeNullString("eCom3", false)).NullInt64(MakeNullInt64(4, false)).NullFloat64(MakeNullFloat64(2.7, false)).
			NullBool(MakeNullBool(true, false)).NullTime(MakeNullTime(now(), false))
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
			NullStrings(MakeNullString("eCom3"), MakeNullString("eCom3")).NullInt64s(MakeNullInt64(4), MakeNullInt64(4)).
			NullFloat64s(MakeNullFloat64(2.7), MakeNullFloat64(2.7)).
			NullBools(MakeNullBool(true)).NullTimes(MakeNullTime(now()), MakeNullTime(now()))
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
				driverValueBytes(nil), MakeNullInt64(3), MakeNullFloat64(2.7), MakeNullBool(true),
				driverValueBytes(`Invoice`), MakeNullString("Creditmemo"), nowSentinel{}, MakeNullTime(now()),
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
			DriverValue(MakeNullInt64(3)).
			DriverValue(MakeNullFloat64(2.7)).
			DriverValue(MakeNullBool(true)).
			DriverValue(driverValueBytes(`Invoice`)).
			DriverValue(MakeNullString("Creditmemo")).
			DriverValue(nowSentinel{}).
			DriverValue(MakeNullTime(now()))

		assert.Exactly(t,
			[]interface{}{nil, []uint8(nil), int64(3), 2.7, true,
				[]uint8{0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65}, "Creditmemo", "2006-01-02 19:04:05", now()},
			args.Interfaces())
	})

	t.Run("Driver.Values panics because not supported", func(t *testing.T) {
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

		args := MakeArgs(10).
			DriverValue(
				driverValueNotSupported(4),
			)
		assert.Nil(t, args)
	})

	t.Run("Driver.Values panics because Value error", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.True(t, errors.Fatal.Match(err), "Should be a fatal error; got %+v", err)
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()

		args := MakeArgs(10).
			DriverValue(
				driverValueError(0),
			)
		assert.Nil(t, args)
	})

}

func TestArguments_WriteTo(t *testing.T) {
	t.Parallel()

	t.Run("no slices, nulls valid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Int(-1).Int64(1).Uint64(2).Float64(3.1).Bool(true).String("eCom1").Bytes([]byte(`eCom2`)).Time(now()).
			NullString(MakeNullString("eCom3")).NullInt64(MakeNullInt64(4)).NullFloat64(MakeNullFloat64(2.7)).
			NullBool(MakeNullBool(true)).NullTime(MakeNullTime(now()))

		buf := new(bytes.Buffer)
		err := args.Write(buf)
		require.NoError(t, err)
		assert.Exactly(t,
			"(NULL,-1,1,2,3.1,1,'eCom1','eCom2','2006-01-02 15:04:05','eCom3',4,2.7,1,'2006-01-02 15:04:05')",
			buf.String())
	})
	t.Run("no slices, nulls invalid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Int(-1).Int64(1).Uint64(2).Float64(3.1).Bool(true).String("eCom1").Bytes([]byte(`eCom2`)).Time(now()).
			NullString(MakeNullString("eCom3", false)).NullInt64(MakeNullInt64(4, false)).NullFloat64(MakeNullFloat64(2.7, false)).
			NullBool(MakeNullBool(true, false)).NullTime(MakeNullTime(now(), false))

		buf := new(bytes.Buffer)
		err := args.Write(buf)
		require.NoError(t, err)
		assert.Exactly(t,
			"(NULL,-1,1,2,3.1,1,'eCom1','eCom2','2006-01-02 15:04:05',NULL,NULL,NULL,NULL,NULL)",
			buf.String())
	})
	t.Run("slices, nulls valid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Ints(-1, -2).Int64s(1, 2).Uint64s(2).Float64s(1.2, 3.1).Bools(false, true).Strings("eCom1", "eCom11").BytesSlice([]byte(`eCom2`)).Times(now(), now()).
			NullStrings(MakeNullString("eCom3"), MakeNullString("eCom3")).NullInt64s(MakeNullInt64(4), MakeNullInt64(5)).NullFloat64s(MakeNullFloat64(2.71), MakeNullFloat64(2.72)).
			NullBools(MakeNullBool(true)).NullTimes(MakeNullTime(now()), MakeNullTime(now()))

		buf := new(bytes.Buffer)
		err := args.Write(buf)
		require.NoError(t, err)
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
		args := MakeArgs(2).NullStrings(MakeNullString("Go2"), MakeNullString("Hello\xc0\x80World"))
		buf := new(bytes.Buffer)
		err := args.Write(buf)
		assert.Exactly(t, "('Go2',)", buf.String())
		assert.True(t, errors.NotValid.Match(err), "Should have a not valid error behaviour %+v", err)
	})
	t.Run("non-utf8 NullString", func(t *testing.T) {
		args := MakeArgs(2).NullString(MakeNullString("Hello\xc0\x80World"))
		buf := new(bytes.Buffer)
		err := args.Write(buf)
		assert.Empty(t, buf.String())
		assert.True(t, errors.NotValid.Match(err), "Should have a not valid error behaviour %+v", err)
	})
	t.Run("bytes as binary", func(t *testing.T) {
		args := MakeArgs(2).Bytes([]byte("\xc0\x80"))
		buf := new(bytes.Buffer)
		require.NoError(t, args.Write(buf))
		assert.Exactly(t, "0xc080", buf.String())
	})
	t.Run("bytesSlice as binary", func(t *testing.T) {
		args := MakeArgs(2).BytesSlice([]byte(`Rusty`), []byte("Go\xc0\x80"))
		buf := new(bytes.Buffer)
		require.NoError(t, args.Write(buf))
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
		require.NoError(t, au.writeTo(buf, 0))
		assert.Empty(t, buf.String(), "buffer should be empty")
	})
}

func TestIFaceToArgs(t *testing.T) {
	t.Parallel()
	t.Run("not supported", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.True(t, errors.NotSupported.Match(err), "%+v", err)
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()
		_ = iFaceToArgs(time.Minute)
	})
	t.Run("all types", func(t *testing.T) {
		nt := now()
		args := iFaceToArgs(
			float32(2.3), float64(2.2),
			int64(5), int(6), int32(7), int16(8), int8(9),
			uint32(math.MaxUint32), uint16(math.MaxUint16), uint8(math.MaxUint8),
			true, "Gopher", []byte(`Hello`),
			now(), &nt, nil,
		)

		assert.Exactly(t, []interface{}{
			float64(2.299999952316284), float64(2.2),
			int64(5), int64(6), int64(7), int64(8), int64(9),
			int64(math.MaxUint32), int64(math.MaxUint16), int64(math.MaxUint8),
			true, "Gopher", []uint8{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			now(), now(), nil,
		}, args.Interfaces())
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
		require.NoError(t, err)
		assert.Exactly(t, uint8(2), a.hasNamedArgs)
	})
	t.Run("hasNamedArgs in condition, no args", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().PlaceHolder(),
			Column("email").Like().NamedArg("ema1l")).WithArgs()
		_, _, err := a.ToSQL()
		require.NoError(t, err)
		assert.Exactly(t, uint8(1), a.hasNamedArgs)
	})
	t.Run("hasNamedArgs in condition, with args", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().PlaceHolder(),
			Column("email").Like().NamedArg("ema1l")).WithArgs().String("my@email.org")
		_, _, err := a.ToSQL()
		require.NoError(t, err)
		assert.Exactly(t, uint8(1), a.hasNamedArgs)
	})
	t.Run("hasNamedArgs none", func(t *testing.T) {
		a := NewSelect("a", "b").From("c").Where(
			Column("id").Greater().Int(221),
			Column("email").Like().Str("em@1l.de")).WithArgs()
		_, _, err := a.ToSQL()
		require.NoError(t, err)
		assert.Exactly(t, uint8(1), a.hasNamedArgs)
	})
}

func TestArguments_MapColumns(t *testing.T) {
	t.Parallel()

	to := MakeArgs(4)
	from := MakeArgs(3)

	t.Run("len=1", func(t *testing.T) {

		from = from.Reset().Int64(3).Float64(2.2).Name("colA").Strings("a", "b")
		rm := newColumnMap(to.Reset(), "colA")
		if err := from.MapColumns(rm); err != nil {
			t.Fatal(err)
		}
		to = rm.Args
		assert.Exactly(t,
			"dml.MakeArgs(1).Name(\"colA\").Strings(\"a\",\"b\")",
			to.GoString())
	})

	t.Run("len=0", func(t *testing.T) {

		from = from.Reset().Name("colZ").Int64(3).Float64(2.2).Name("colA").Strings("a", "b")
		rm := newColumnMap(to.Reset())
		if err := from.MapColumns(rm); err != nil {
			t.Fatal(err)
		}
		to = rm.Args
		assert.Exactly(t,
			"dml.MakeArgs(3).Name(\"colZ\").Int64(3).Float64(2.200000).Name(\"colA\").Strings(\"a\",\"b\")",
			to.GoString())
	})

	t.Run("len>1", func(t *testing.T) {

		from = from.Reset().Name("colZ").Int64(3).Uint64(6).Name("colB").Float64(2.2).String("c").Name("colA").Strings("a", "b")
		rm := newColumnMap(to.Reset(), "colA", "colB")
		if err := from.MapColumns(rm); err != nil {
			t.Fatal(err)
		}
		to = rm.Args
		assert.Exactly(t,
			"dml.MakeArgs(2).Name(\"colA\").Strings(\"a\",\"b\").Name(\"colB\").Float64(2.200000)",
			to.GoString())
	})
}

func TestArguments_NextUnnamedArg(t *testing.T) {
	t.Parallel()

	t.Run("three occurrences", func(t *testing.T) {
		args := MakeArgs(5).Name("colZ").Int64(3).Uint64(6).Name("colB").Float64(2.2).String("c").Name("colA").Strings("a", "b")

		a, ok := args.nextUnnamedArg()
		require.True(t, ok, "Should find an unnamed argument")
		assert.Empty(t, a.name)
		assert.Exactly(t, uint64(6), a.value)

		a, ok = args.nextUnnamedArg()
		require.True(t, ok, "Should find an unnamed argument")
		assert.Empty(t, a.name)
		assert.Exactly(t, "c", a.value)

		a, ok = args.nextUnnamedArg()
		require.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, argument{}, a)

		args.Reset().Float64(3.14159).Name("price").Float64(2.7182).Time(now())

		a, ok = args.nextUnnamedArg()
		require.True(t, ok, "Should find an unnamed argument")
		assert.Empty(t, a.name)
		assert.Exactly(t, 3.14159, a.value)

		a, ok = args.nextUnnamedArg()
		require.True(t, ok, "Should find an unnamed argument")
		assert.Empty(t, a.name)
		assert.Exactly(t, now(), a.value)

		a, ok = args.nextUnnamedArg()
		require.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, argument{}, a)
	})

	t.Run("zero occurrences", func(t *testing.T) {
		args := MakeArgs(5).Name("colZ").Int64(3).Name("colB").Float64(2.2).Name("colA").Strings("a", "b")

		a, ok := args.nextUnnamedArg()
		require.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, argument{}, a)

		a, ok = args.nextUnnamedArg()
		require.False(t, ok, "Should NOT find an unnamed argument")
		assert.Exactly(t, argument{}, a)
	})

}

type myToSQL struct {
	sql  string
	args []interface{}
	error
}

func (m myToSQL) ToSQL() (string, []interface{}, error) {
	return m.sql, m.args, m.error
}

func TestArguments_ExecContext(t *testing.T) {
	t.Parallel()
	t.Skip("TODO IMPLEMENT")
	//haveErr := errors.AlreadyClosed.Newf("Who closed myself?")
	//
	//t.Run("ToSQL error", func(t *testing.T) {
	//	stmt, err := dml.Exec(context.TODO(), nil, myToSQL{error: haveErr})
	//	assert.Nil(t, stmt)
	//	assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	//})
}

func TestArguments_QueryContext(t *testing.T) {
	t.Parallel()
	t.Skip("TODO IMPLEMENT")
	//haveErr := errors.AlreadyClosed.Newf("Who closed myself?")
	//
	//t.Run("ToSQL error", func(t *testing.T) {
	//	stmt, err := dml.Prepare(context.TODO(), nil, myToSQL{error: haveErr})
	//	assert.Nil(t, stmt)
	//	assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	//})
	//t.Run("ToSQL prepare error", func(t *testing.T) {
	//	dbc, dbMock := dmltest.MockDB(t)
	//	defer dmltest.MockClose(t, dbc, dbMock)
	//
	//	dbMock.ExpectPrepare("SELECT `a` FROM `b`").WillReturnError(haveErr)
	//
	//	stmt, err := dml.Prepare(context.TODO(), dbc.DB, myToSQL{sql: "SELECT `a` FROM `b`"})
	//	assert.Nil(t, stmt)
	//	assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	//})
}
