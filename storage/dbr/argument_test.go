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

package dbr

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ fmt.Stringer = Op(0)
var _ Value = Int(0)
var _ Value = Int64(0)
var _ Value = Float64(0)
var _ Value = Bool(true)
var _ Value = (*BytesSlice)(nil)
var _ Value = (*Ints)(nil)
var _ Value = (*Int64s)(nil)
var _ Value = (*Float64s)(nil)
var _ Value = (*Times)(nil)
var _ Value = (*Strings)(nil)
var _ Value = (*NullString)(nil)
var _ Value = (*ArgNullStrings)(nil)
var _ Value = (*NullFloat64)(nil)
var _ Value = (*ArgNullFloat64s)(nil)
var _ Value = (*NullTime)(nil)
var _ Value = (*ArgNullTimes)(nil)
var _ Value = (*NullInt64)(nil)
var _ Value = (*ArgNullInt64s)(nil)
var _ Value = (*NullBool)(nil)
var _ Value = (*DriverValues)(nil)
var _ driver.Valuer = (*Bool)(nil)
var _ driver.Valuer = (*Int)(nil)
var _ driver.Valuer = (*Int64)(nil)
var _ driver.Valuer = (*Float64)(nil)
var _ driver.Valuer = (*Time)(nil)
var _ driver.Valuer = (*Bytes)(nil)
var _ driver.Valuer = (*String)(nil)

func TestDriverValuer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have driver.Valuer
		want driver.Value
	}{
		{Bool(true), true},
		{Int(7), int64(7)},
		{Int64(8), int64(8)},
		{Float64(8.9), float64(8.9)},
		{MakeTime(now()), now()},
		{Bytes(`Go2`), []byte(`Go2`)},
		{String(`Go2`), `Go2`},
	}
	for i, test := range tests {
		v, err := test.have.Value()
		if err != nil {
			t.Fatalf("index %d with %+v", i, err)
		}
		assert.Exactly(t, test.want, v, "Index %d", i)
	}
}

func TestOpRune(t *testing.T) {
	t.Parallel()
	s := NewSelect().From("tableA").AddColumns("a", "b").
		Where(
			Column("a1").Like().String("H_ll_"),
			Column("a1").Like().NullString(NullString{}),
			Column("a1").Like().NullString(MakeNullString("NullString")),
			Column("a1").Like().Float64(2.718281),
			Column("a1").Like().NullFloat64(NullFloat64{}),
			Column("a1").Like().NullFloat64(MakeNullFloat64(-2.718281)),
			Column("a1").Like().Int64(2718281),
			Column("a1").Like().NullInt64(NullInt64{}),
			Column("a1").Like().NullInt64(MakeNullInt64(-987)),
			Column("a1").Like().Int(2718281),
			Column("a1").Like().Bool(true),
			Column("a1").Like().NullBool(NullBool{}),
			Column("a1").Like().NullBool(MakeNullBool(false)),
			Column("a1").Like().Time(now()),
			Column("a1").Like().NullTime(MakeNullTime(now().Add(time.Minute))),
			Column("a1").Like().Null(),
			Column("a1").Like().Bytes([]byte(`H3llo`)),
			Column("a1").Like().Value(MakeNullInt64(2345)),

			Column("a2").NotLike().String("H_ll_"),
			Column("a2").NotLike().NullString(NullString{}),
			Column("a2").NotLike().NullString(MakeNullString("NullString")),
			Column("a2").NotLike().Float64(2.718281),
			Column("a2").NotLike().NullFloat64(NullFloat64{}),
			Column("a2").NotLike().NullFloat64(MakeNullFloat64(-2.718281)),
			Column("a2").NotLike().Int64(2718281),
			Column("a2").NotLike().NullInt64(NullInt64{}),
			Column("a2").NotLike().NullInt64(MakeNullInt64(-987)),
			Column("a2").NotLike().Int(2718281),
			Column("a2").NotLike().Bool(true),
			Column("a2").NotLike().NullBool(NullBool{}),
			Column("a2").NotLike().NullBool(MakeNullBool(false)),
			Column("a2").NotLike().Time(now()),
			Column("a2").NotLike().NullTime(MakeNullTime(now().Add(time.Minute))),
			Column("a2").NotLike().Null(),
			Column("a2").NotLike().Bytes([]byte(`H3llo`)),
			Column("a2").NotLike().Value(MakeNullInt64(2345)),

			Column("a301").In().Strings("Go1", "Go2"),
			Column("a302").In().NullString(NullString{}, NullString{}),
			Column("a303").In().NullString(MakeNullString("NullString")),
			Column("a304").In().Float64s(2.718281, 3.14159),
			Column("a305").In().NullFloat64(NullFloat64{}),
			Column("a306").In().NullFloat64(MakeNullFloat64(-2.718281), MakeNullFloat64(-3.14159)),
			Column("a307").In().Int64s(2718281, 314159),
			Column("a308").In().NullInt64(NullInt64{}),
			Column("a309").In().NullInt64(MakeNullInt64(-987), MakeNullInt64(-654)),
			Column("a310").In().Ints(2718281, 314159),
			Column("a311").In().Bools(true, false),
			Column("a312").In().NullBool(NullBool{}),
			Column("a313").In().NullBool(MakeNullBool(true)),
			Column("a314").In().Times(now(), now()),
			Column("a315").In().NullTime(MakeNullTime(now().Add(time.Minute))),
			Column("a316").In().Null(),
			Column("a317").In().Bytes([]byte(`H3llo1`)),
			Column("a320").In().Value(MakeNullInt64(2345), MakeNullFloat64(3.14159)),

			Column("a401").SpaceShip().String("H_ll_"),
			Column("a402").SpaceShip().NullString(NullString{}),
			Column("a403").SpaceShip().NullString(MakeNullString("NullString")),
		)
	compareToSQL(t, s, nil,
		"SELECT `a`, `b` FROM `tableA` WHERE (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` IS NULL) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` IS NULL) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a301` IN (?,?)) AND (`a302` IN (?,?)) AND (`a303` IN (?)) AND (`a304` IN (?,?)) AND (`a305` IN (?)) AND (`a306` IN (?,?)) AND (`a307` IN (?,?)) AND (`a308` IN (?)) AND (`a309` IN (?,?)) AND (`a310` IN (?,?)) AND (`a311` IN (?,?)) AND (`a312` IN (?)) AND (`a313` IN (?)) AND (`a314` IN (?,?)) AND (`a315` IN (?)) AND (`a316` IS NULL) AND (`a317` IN (?)) AND (`a320` IN (?,?)) AND (`a401` <=> ?) AND (`a402` <=> ?) AND (`a403` <=> ?)",
		"SELECT `a`, `b` FROM `tableA` WHERE (`a1` LIKE 'H_ll_') AND (`a1` LIKE NULL) AND (`a1` LIKE 'NullString') AND (`a1` LIKE 2.718281) AND (`a1` LIKE NULL) AND (`a1` LIKE -2.718281) AND (`a1` LIKE 2718281) AND (`a1` LIKE NULL) AND (`a1` LIKE -987) AND (`a1` LIKE 2718281) AND (`a1` LIKE 1) AND (`a1` LIKE NULL) AND (`a1` LIKE 0) AND (`a1` LIKE '2006-01-02 15:04:05') AND (`a1` LIKE '2006-01-02 15:05:05') AND (`a1` IS NULL) AND (`a1` LIKE 'H3llo') AND (`a1` LIKE 2345) AND (`a2` NOT LIKE 'H_ll_') AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE 'NullString') AND (`a2` NOT LIKE 2.718281) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE -2.718281) AND (`a2` NOT LIKE 2718281) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE -987) AND (`a2` NOT LIKE 2718281) AND (`a2` NOT LIKE 1) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE 0) AND (`a2` NOT LIKE '2006-01-02 15:04:05') AND (`a2` NOT LIKE '2006-01-02 15:05:05') AND (`a2` IS NULL) AND (`a2` NOT LIKE 'H3llo') AND (`a2` NOT LIKE 2345) AND (`a301` IN ('Go1','Go2')) AND (`a302` IN (NULL,NULL)) AND (`a303` IN ('NullString')) AND (`a304` IN (2.718281,3.14159)) AND (`a305` IN (NULL)) AND (`a306` IN (-2.718281,-3.14159)) AND (`a307` IN (2718281,314159)) AND (`a308` IN (NULL)) AND (`a309` IN (-987,-654)) AND (`a310` IN (2718281,314159)) AND (`a311` IN (1,0)) AND (`a312` IN (NULL)) AND (`a313` IN (1)) AND (`a314` IN ('2006-01-02 15:04:05','2006-01-02 15:04:05')) AND (`a315` IN ('2006-01-02 15:05:05')) AND (`a316` IS NULL) AND (`a317` IN ('H3llo1')) AND (`a320` IN (2345,3.14159)) AND (`a401` <=> 'H_ll_') AND (`a402` <=> NULL) AND (`a403` <=> 'NullString')",
		"H_ll_", interface{}(nil), "NullString", 2.718281, interface{}(nil),
		-2.718281, int64(2718281), interface{}(nil), int64(-987), int64(2718281), true,
		interface{}(nil), false, now(), now().Add(time.Minute),
		[]uint8{0x48, 0x33, 0x6c, 0x6c, 0x6f}, int64(2345), "H_ll_",
		interface{}(nil), "NullString", 2.718281, interface{}(nil), -2.718281, int64(2718281),
		interface{}(nil), int64(-987), int64(2718281), true, interface{}(nil), false, now(), now().Add(time.Minute),
		[]uint8{0x48, 0x33, 0x6c, 0x6c, 0x6f}, int64(2345),
		"Go1", "Go2", interface{}(nil), interface{}(nil), "NullString", 2.718281, 3.14159,
		interface{}(nil), -2.718281, -3.14159, int64(2718281), int64(314159), interface{}(nil),
		int64(-987), int64(-654), int64(2718281), int64(314159), true, false, interface{}(nil), true,
		now(), now(), now().Add(time.Minute), []uint8{0x48, 0x33, 0x6c, 0x6c, 0x6f, 0x31},
		int64(2345), 3.14159, "H_ll_", interface{}(nil), "NullString",
	)
}

func TestOp_String(t *testing.T) {
	t.Parallel()
	var o Op
	assert.Exactly(t, "=", o.String())
	assert.Exactly(t, "🚀", SpaceShip.String())
}

func TestOpArgs(t *testing.T) {
	t.Parallel()
	t.Run("Null with place holder IN,Regexp", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a315").In().Null(),
				Column("a316").In().PlaceHolder(),
				Column("a317").Regexp().PlaceHolder(),
				Column("a317").NotRegexp().PlaceHolder(),
			),
			nil,
			"SELECT `a`, `b` FROM `t1` WHERE (`a315` IS NULL) AND (`a316` IN (?)) AND (`a317` REGEXP ?) AND (`a317` NOT REGEXP ?)",
			"SELECT `a`, `b` FROM `t1` WHERE (`a315` IS NULL) AND (`a316` IN (?)) AND (`a317` REGEXP ?) AND (`a317` NOT REGEXP ?)",
			[]interface{}{}...,
		)
	})

	t.Run("Args In", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a311").Xor().Int(9),
				Column("a313").In().Float64(3.3),
				Column("a314").In().Int64(33),
				Column("a312").In().Int(44),
				Column("a315").In().String(`Go1`),
				Column("a316").In().BytesSlice([]byte(`Go`), []byte(`Rust`)),
			),
			nil,
			"SELECT `a`, `b` FROM `t1` WHERE (`a311` XOR ?) AND (`a313` IN (?)) AND (`a314` IN (?)) AND (`a312` IN (?)) AND (`a315` IN (?)) AND (`a316` IN (?,?))",
			"SELECT `a`, `b` FROM `t1` WHERE (`a311` XOR 9) AND (`a313` IN (3.3)) AND (`a314` IN (33)) AND (`a312` IN (44)) AND (`a315` IN ('Go1')) AND (`a316` IN ('Go','Rust'))",
			int64(9), 3.3, int64(33), int64(44), "Go1", []uint8{0x47, 0x6f}, []uint8{0x52, 0x75, 0x73, 0x74},
		)
	})

	t.Run("BytesSlice BETWEEN strings", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a316").Between().BytesSlice([]byte(`Go`), []byte(`Rust`)),
			),
			nil,
			"SELECT `a`, `b` FROM `t1` WHERE (`a316` BETWEEN ? AND ?)",
			"SELECT `a`, `b` FROM `t1` WHERE (`a316` BETWEEN 'Go' AND 'Rust')",
			[]uint8{0x47, 0x6f}, []uint8{0x52, 0x75, 0x73, 0x74},
		)
	})
	t.Run("BytesSlice IN binary", func(t *testing.T) {

		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a316").In().BytesSlice([]byte{66, 250, 67}, []byte(`Rust`), []byte("\xFB\xBF\xBF\xBF\xBF")),
			),
			nil,
			"SELECT `a`, `b` FROM `t1` WHERE (`a316` IN (?,?,?))",
			"SELECT `a`, `b` FROM `t1` WHERE (`a316` IN (0x42fa43,'Rust',0xfbbfbfbfbf))",
			[]uint8{0x42, 0xfa, 0x43}, []uint8{0x52, 0x75, 0x73, 0x74}, []uint8{0xfb, 0xbf, 0xbf, 0xbf, 0xbf},
		)
	})
	t.Run("ArgValue IN", func(t *testing.T) {

		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a319").In().Value(
					MakeNullFloat64(3.141),
					MakeNullString("G'o"),
					Bytes{66, 250, 67},
					MakeNullTime(now()),
					Bytes([]byte("x\x00y")),
				),
			),
			nil,
			"SELECT `a`, `b` FROM `t1` WHERE (`a319` IN (?,?,?,?,?))",
			"SELECT `a`, `b` FROM `t1` WHERE (`a319` IN (3.141,'G\\'o',0x42fa43,'2006-01-02 15:04:05',0x780079))",
			3.141, `G'o`, []uint8{0x42, 0xfa, 0x43}, now(), []uint8{0x78, 0x0, 0x79},
		)
	})
	t.Run("ArgValue BETWEEN values", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("t1").Where(
				Column("a319").Between().Value(MakeNullFloat64(3.141), MakeNullString("G'o")),
			),
			nil,
			"SELECT `a`, `b` FROM `t1` WHERE (`a319` BETWEEN ? AND ?)",
			"SELECT `a`, `b` FROM `t1` WHERE (`a319` BETWEEN 3.141 AND 'G\\'o')",
			3.141, `G'o`,
		)
	})
}

func TestNullStringFrom(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "product", MakeNullString("product").String)
	assert.True(t, MakeNullString("product").Valid)
	//assert.False(t, NullStringFromPtr(nil).Valid)
	assert.True(t, MakeNullString("").Valid)
	v, err := MakeNullString("product").Value()
	assert.NoError(t, err)
	assert.Equal(t, "product", v)
}

func TestNewNullInt64(t *testing.T) {
	t.Parallel()
	assert.EqualValues(t, 1257894000, MakeNullInt64(1257894000).Int64)
	assert.True(t, MakeNullInt64(1257894000).Valid)
	assert.True(t, MakeNullInt64(0).Valid)
	v, err := MakeNullInt64(1257894000).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, 1257894000, v)
}

func TestNewNullFloat64(t *testing.T) {
	t.Parallel()
	var test = 1257894000.93445000001
	assert.Equal(t, test, MakeNullFloat64(test).Float64)
	assert.True(t, MakeNullFloat64(test).Valid)
	assert.True(t, MakeNullFloat64(0).Valid)
	v, err := MakeNullFloat64(test).Value()
	assert.NoError(t, err)
	assert.Equal(t, test, v)
}

func TestNewNullTime(t *testing.T) {
	t.Parallel()
	var test = time.Now()
	assert.Equal(t, test, MakeNullTime(test).Time)
	assert.True(t, MakeNullTime(test).Valid)
	assert.True(t, MakeNullTime(time.Time{}).Valid)

	v, err := MakeNullTime(test).Value()
	assert.NoError(t, err)
	assert.Equal(t, test, v)
}

func TestNewNullBool(t *testing.T) {
	t.Parallel()

	assert.Equal(t, true, MakeNullBool(true).Bool)
	assert.True(t, MakeNullBool(true).Valid)
	assert.True(t, MakeNullBool(false).Valid)
	v, err := MakeNullBool(true).Value()
	assert.NoError(t, err)
	assert.Equal(t, true, v)
}

func TestNullTypeScanning(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)

	type nullTypeScanningTest struct {
		record *nullTypedRecord
		valid  bool
	}

	tests := []nullTypeScanningTest{
		{
			record: &nullTypedRecord{ID: 1},
			valid:  false,
		},
		{
			record: newNullTypedRecordWithData(),
			valid:  true,
		},
	}

	for _, test := range tests {
		// Create the record in the db
		res, err := s.InsertInto("null_types").
			AddColumns("string_val", "int64_val", "float64_val", "time_val", "bool_val").
			AddRecords(test.record).Exec(context.TODO())
		require.NoError(t, err)
		id, err := res.LastInsertId()
		assert.NoError(t, err)

		// Scan it back and check that all fields are of the correct validity and are
		// equal to the reference record
		nullTypeSet := &nullTypedRecord{}
		_, err = s.Select("*").From("null_types").Where(Expression("id = ?", Int64(id))).Load(context.TODO(), nullTypeSet)
		assert.NoError(t, err)

		assert.Equal(t, test.record, nullTypeSet)
		assert.Equal(t, test.valid, nullTypeSet.StringVal.Valid)
		assert.Equal(t, test.valid, nullTypeSet.Int64Val.Valid)
		assert.Equal(t, test.valid, nullTypeSet.Float64Val.Valid)
		assert.Equal(t, test.valid, nullTypeSet.TimeVal.Valid)
		assert.Equal(t, test.valid, nullTypeSet.BoolVal.Valid)

		nullTypeSet.StringVal.String = "newStringVal"
		assert.NotEqual(t, test.record, nullTypeSet)
	}
}

func TestNullTypeJSONMarshal(t *testing.T) {
	t.Parallel()
	type nullTypeJSONTest struct {
		record       *nullTypedRecord
		expectedJSON []byte
	}

	tests := []nullTypeJSONTest{
		{
			record:       &nullTypedRecord{},
			expectedJSON: []byte(`{"ID":0,"StringVal":null,"Int64Val":null,"Float64Val":null,"TimeVal":null,"BoolVal":null}`),
		},
		{
			record:       newNullTypedRecordWithData(),
			expectedJSON: []byte(`{"ID":2,"StringVal":"wow","Int64Val":42,"Float64Val":1.618,"TimeVal":"2009-01-03T18:15:05Z","BoolVal":true}`),
		},
	}

	for _, test := range tests {
		// Marshal the record
		rawJSON, err := json.Marshal(test.record)
		assert.NoError(t, err)
		assert.Equal(t, string(test.expectedJSON), string(rawJSON))

		// Unmarshal it back
		newRecord := &nullTypedRecord{}
		err = json.Unmarshal([]byte(rawJSON), newRecord)
		assert.NoError(t, err)
		assert.Equal(t, test.record, newRecord)
	}
}

func newNullTypedRecordWithData() *nullTypedRecord {
	return &nullTypedRecord{
		ID:         2,
		StringVal:  NullString{NullString: sql.NullString{String: "wow", Valid: true}},
		Int64Val:   NullInt64{NullInt64: sql.NullInt64{Int64: 42, Valid: true}},
		Float64Val: NullFloat64{NullFloat64: sql.NullFloat64{Float64: 1.618, Valid: true}},
		TimeVal:    NullTime{Time: time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC), Valid: true},
		BoolVal:    NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}},
	}
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
