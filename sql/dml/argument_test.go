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
	return nil, errors.NewAbortedf("WE've aborted something")
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
			NullString(MakeNullString("eCom3"), MakeNullString("eCom3")).NullInt64(MakeNullInt64(4), MakeNullInt64(4)).NullFloat64(MakeNullFloat64(2.7), MakeNullFloat64(2.7)).
			NullBool(MakeNullBool(true)).NullTime(MakeNullTime(now()), MakeNullTime(now()))
		assert.Exactly(t, 26, args.Len(), "Length mismatch")
		assert.Exactly(t,
			"dml.MakeArgs(15).Null().Int(-1).Int64s([]int64{1, 2}...).Uints([]uint{0x237, 0x2fd}...).Uint64s([]uint64{0x2}...).Float64s([]float64{1.2, 3.1}...).Bools([]bool{false, true}...).Strings(\"eCom1\",\"eCom11\").BytesSlice([]byte(nil),[]byte{0x65, 0x43, 0x6f, 0x6d, 0x32}).Times(time.Unix(1136228645,2),time.Unix(1136228645,2)).NullString(dml.MakeNullString(`eCom3`),dml.MakeNullString(`eCom3`)).NullInt64(dml.MakeNullInt64(4),dml.MakeNullInt64(4)).NullFloat64(dml.MakeNullFloat64(2.7),dml.MakeNullFloat64(2.7)).NullBool(dml.MakeNullBool(true)).NullTime(dml.MakeNullTime(time.Unix(1136228645,2),dml.MakeNullTime(time.Unix(1136228645,2))",
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
			NullString(MakeNullString("eCom3"), MakeNullString("eCom3")).NullInt64(MakeNullInt64(4), MakeNullInt64(4)).
			NullFloat64(MakeNullFloat64(2.7), MakeNullFloat64(2.7)).
			NullBool(MakeNullBool(true)).NullTime(MakeNullTime(now()), MakeNullTime(now()))
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
					assert.True(t, errors.IsNotSupported(err), "Should be a not supported error; got %+v", err)
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
					assert.True(t, errors.IsFatal(err), "Should be a fatal error; got %+v", err)
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
			NullString(MakeNullString("eCom3"), MakeNullString("eCom3")).NullInt64(MakeNullInt64(4), MakeNullInt64(5)).NullFloat64(MakeNullFloat64(2.71), MakeNullFloat64(2.72)).
			NullBool(MakeNullBool(true)).NullTime(MakeNullTime(now()), MakeNullTime(now()))

		buf := new(bytes.Buffer)
		err := args.Write(buf)
		require.NoError(t, err)
		assert.Exactly(t,
			"(NULL,-1,-2,1,2,2,1.2,3.1,0,1,'eCom1','eCom11','eCom2','2006-01-02 15:04:05','2006-01-02 15:04:05','eCom3','eCom3',4,5,2.71,2.72,1,'2006-01-02 15:04:05','2006-01-02 15:04:05')",
			buf.String())
	})
	t.Run("non-utf8 string", func(t *testing.T) {
		args := MakeArgs(2).String("\xc0\x80")
		buf := new(bytes.Buffer)
		err := args.Write(buf)
		assert.Exactly(t, `(`, buf.String())
		assert.True(t, errors.IsNotValid(err), "Should have a not valid error behaviour %+v", err)
	})
	t.Run("non-utf8 strings", func(t *testing.T) {
		args := MakeArgs(2).Strings("Go", "\xc0\x80")
		buf := new(bytes.Buffer)
		err := args.Write(buf)
		assert.Exactly(t, `('Go',`, buf.String())
		assert.True(t, errors.IsNotValid(err), "Should have a not valid error behaviour %+v", err)
	})
	t.Run("non-utf8 NullString", func(t *testing.T) {
		args := MakeArgs(2).NullString(MakeNullString("Go2"), MakeNullString("Hello\xc0\x80World"))
		buf := new(bytes.Buffer)
		err := args.Write(buf)
		assert.Exactly(t, "('Go2',", buf.String())
		assert.True(t, errors.IsNotValid(err), "Should have a not valid error behaviour %+v", err)
	})
	t.Run("bytes as binary", func(t *testing.T) {
		args := MakeArgs(2).Bytes([]byte("\xc0\x80"))
		buf := new(bytes.Buffer)
		require.NoError(t, args.Write(buf))
		assert.Exactly(t, "(0xc080)", buf.String())
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
					assert.True(t, errors.IsNotSupported(err), "Should be a not supported error; got %+v", err)
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()

		au := argument{value: complex64(1)}
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
					assert.True(t, errors.IsNotSupported(err), "%+v", err)
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
