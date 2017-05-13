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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//var _ Argument = argInt(0)
var _ fmt.Stringer = Op(0)
var _ Argument = argInt64(0)
var _ Argument = argFloat64(0)
var _ Argument = argBool(true)
var _ Argument = (*argBytes)(nil)
var _ Argument = (*argInts)(nil)
var _ Argument = (*argInt64s)(nil)
var _ Argument = (*argFloat64s)(nil)
var _ Argument = (*argTimes)(nil)
var _ Argument = (*argBools)(nil)
var _ Argument = (*argStrings)(nil)
var _ Argument = (*NullString)(nil)
var _ Argument = (*argNullStrings)(nil)
var _ Argument = (*NullFloat64)(nil)
var _ Argument = (*argNullFloat64s)(nil)
var _ Argument = (*NullBytes)(nil)
var _ Argument = (*NullTime)(nil)
var _ Argument = (*argNullTimes)(nil)
var _ Argument = (*NullInt64)(nil)
var _ Argument = (*argNullInt64s)(nil)
var _ Argument = (*NullBool)(nil)
var _ Argument = (*argValue)(nil)
var _ driver.Value = (*Arguments)(nil)

func TestOpRune(t *testing.T) {
	t.Parallel()
	s := NewSelect().From("tableA").AddColumns("a", "b").
		Where(
			Column("a1", Like.Str("H_ll_")),
			Column("a1", Like.NullString(NullString{})),
			Column("a1", Like.NullString(MakeNullString("NullString"))),
			Column("a1", Like.Float64(2.718281)),
			Column("a1", Like.NullFloat64(NullFloat64{})),
			Column("a1", Like.NullFloat64(MakeNullFloat64(-2.718281))),
			Column("a1", Like.Int64(2718281)),
			Column("a1", Like.NullInt64(NullInt64{})),
			Column("a1", Like.NullInt64(MakeNullInt64(-987))),
			Column("a1", Like.Int(2718281)),
			Column("a1", Like.Bool(true)),
			Column("a1", Like.NullBool(NullBool{})),
			Column("a1", Like.NullBool(MakeNullBool(false))),
			Column("a1", Like.Time(now())),
			Column("a1", Like.NullTime(MakeNullTime(now().Add(time.Minute)))),
			Column("a1", Like.Null()),
			Column("a1", Like.Bytes([]byte(`H3llo`))),
			Column("a1", Like.Value(MakeNullInt64(2345))),

			Column("a2", NotLike.Str("H_ll_")),
			Column("a2", NotLike.NullString(NullString{})),
			Column("a2", NotLike.NullString(MakeNullString("NullString"))),
			Column("a2", NotLike.Float64(2.718281)),
			Column("a2", NotLike.NullFloat64(NullFloat64{})),
			Column("a2", NotLike.NullFloat64(MakeNullFloat64(-2.718281))),
			Column("a2", NotLike.Int64(2718281)),
			Column("a2", NotLike.NullInt64(NullInt64{})),
			Column("a2", NotLike.NullInt64(MakeNullInt64(-987))),
			Column("a2", NotLike.Int(2718281)),
			Column("a2", NotLike.Bool(true)),
			Column("a2", NotLike.NullBool(NullBool{})),
			Column("a2", NotLike.NullBool(MakeNullBool(false))),
			Column("a2", NotLike.Time(now())),
			Column("a2", NotLike.NullTime(MakeNullTime(now().Add(time.Minute)))),
			Column("a2", NotLike.Null()),
			Column("a2", NotLike.Bytes([]byte(`H3llo`))),
			Column("a2", NotLike.Value(MakeNullInt64(2345))),

			Column("a301", In.Str("Go1", "Go2")),
			Column("a302", In.NullString(NullString{}, NullString{})),
			Column("a303", In.NullString(MakeNullString("NullString"))),
			Column("a304", In.Float64(2.718281, 3.14159)),
			Column("a305", In.NullFloat64(NullFloat64{})),
			Column("a306", In.NullFloat64(MakeNullFloat64(-2.718281), MakeNullFloat64(-3.14159))),
			Column("a307", In.Int64(2718281, 314159)),
			Column("a308", In.NullInt64(NullInt64{})),
			Column("a309", In.NullInt64(MakeNullInt64(-987), MakeNullInt64(-654))),
			Column("a310", In.Int(2718281, 314159)),
			Column("a311", In.Bool(true, false)),
			Column("a312", In.NullBool(NullBool{})),
			Column("a313", In.NullBool(MakeNullBool(true))),
			Column("a314", In.Time(now(), now())),
			Column("a315", In.NullTime(MakeNullTime(now().Add(time.Minute)))),
			Column("a316", In.Null()),
			Column("a317", In.Bytes([]byte(`H3llo1`))),
			Column("a320", In.Value(MakeNullInt64(2345), MakeNullFloat64(3.14159))),

			Column("a401", SpaceShip.Str("H_ll_")),
			Column("a402", SpaceShip.NullString(NullString{})),
			Column("a403", SpaceShip.NullString(MakeNullString("NullString"))),
		)
	sStr, args, err := s.ToSQL()
	require.NoError(t, err)
	assert.Equal(t, "SELECT `a`, `b` FROM `tableA` WHERE (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a1` LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a2` NOT LIKE ?) AND (`a301` IN ?) AND (`a302` IN ?) AND (`a303` IN ?) AND (`a304` IN ?) AND (`a305` IN ?) AND (`a306` IN ?) AND (`a307` IN ?) AND (`a308` IN ?) AND (`a309` IN ?) AND (`a310` IN ?) AND (`a311` IN ?) AND (`a312` IN ?) AND (`a313` IN ?) AND (`a314` IN ?) AND (`a315` IN ?) AND (`a316` IN ?) AND (`a317` IN ?) AND (`a320` IN ?) AND (`a401` <=> ?) AND (`a402` <=> ?) AND (`a403` <=> ?)",
		sStr)
	assert.Equal(t, 57, len(args), "Argument count")

	sStr, args, err = s.Interpolate().ToSQL()
	require.NoError(t, err)
	assert.Nil(t, args)
	assert.Equal(t,
		"SELECT `a`, `b` FROM `tableA` WHERE (`a1` LIKE 'H_ll_') AND (`a1` LIKE NULL) AND (`a1` LIKE 'NullString') AND (`a1` LIKE 2.718281) AND (`a1` LIKE NULL) AND (`a1` LIKE -2.718281) AND (`a1` LIKE 2718281) AND (`a1` LIKE NULL) AND (`a1` LIKE -987) AND (`a1` LIKE 2718281) AND (`a1` LIKE 1) AND (`a1` LIKE NULL) AND (`a1` LIKE 0) AND (`a1` LIKE '2006-01-02 15:04:05') AND (`a1` LIKE '2006-01-02 15:05:05') AND (`a1` LIKE NULL) AND (`a1` LIKE 'H3llo') AND (`a1` LIKE 2345) AND (`a2` NOT LIKE 'H_ll_') AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE 'NullString') AND (`a2` NOT LIKE 2.718281) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE -2.718281) AND (`a2` NOT LIKE 2718281) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE -987) AND (`a2` NOT LIKE 2718281) AND (`a2` NOT LIKE 1) AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE 0) AND (`a2` NOT LIKE '2006-01-02 15:04:05') AND (`a2` NOT LIKE '2006-01-02 15:05:05') AND (`a2` NOT LIKE NULL) AND (`a2` NOT LIKE 'H3llo') AND (`a2` NOT LIKE 2345) AND (`a301` IN ('Go1','Go2')) AND (`a302` IN (NULL,NULL)) AND (`a303` IN 'NullString') AND (`a304` IN (2.718281,3.14159)) AND (`a305` IN NULL) AND (`a306` IN (-2.718281,-3.14159)) AND (`a307` IN (2718281,314159)) AND (`a308` IN NULL) AND (`a309` IN (-987,-654)) AND (`a310` IN (2718281,314159)) AND (`a311` IN (1,0)) AND (`a312` IN NULL) AND (`a313` IN 1) AND (`a314` IN ('2006-01-02 15:04:05','2006-01-02 15:04:05')) AND (`a315` IN '2006-01-02 15:05:05') AND (`a316` IN (NULL)) AND (`a317` IN ('H3llo1')) AND (`a320` IN (2345,3.14159)) AND (`a401` <=> 'H_ll_') AND (`a402` <=> NULL) AND (`a403` <=> 'NullString')",
		sStr)
}

func TestOpArgs(t *testing.T) {
	t.Parallel()
	t.Run("ArgNull IN", func(t *testing.T) {
		s := NewSelect("a", "b").From("t1").Where(
			Column("a316", In.Null()),
			Column("a317", Regexp.Null()),
		).Interpolate()
		sStr, args, err := s.ToSQL()
		require.NoError(t, err)
		assert.Equal(t, "SELECT `a`, `b` FROM `t1` WHERE (`a316` IN (NULL)) AND (`a317` REGEXP NULL)", sStr)
		assert.Nil(t, args)
	})
	t.Run("Args IN", func(t *testing.T) {
		s := NewSelect("a", "b").From("t1").Where(
			Column("a313", In.Float64(3.3)),
			Column("a314", In.Int64(33)),
			Column("a312", In.Int(44)),
			Column("a315", In.Str(`Go1`)),
			Column("a316", In.Bytes([]byte(`Go`), []byte(`Rust`))),
		).Interpolate()
		sStr, args, err := s.ToSQL()
		require.NoError(t, err)
		assert.Equal(t, "SELECT `a`, `b` FROM `t1` WHERE (`a313` IN (3.3)) AND (`a314` IN (33)) AND (`a312` IN (44)) AND (`a315` IN ('Go1')) AND (`a316` IN ('Go','Rust'))", sStr)
		assert.Nil(t, args)
	})
	t.Run("ArgBytes BETWEEN strings", func(t *testing.T) {
		s := NewSelect("a", "b").From("t1").Where(
			Column("a316", Between.Bytes([]byte(`Go`), []byte(`Rust`))),
		).Interpolate()
		sStr, args, err := s.ToSQL()
		require.NoError(t, err)
		assert.Equal(t, "SELECT `a`, `b` FROM `t1` WHERE (`a316` BETWEEN 'Go' AND 'Rust')", sStr)
		assert.Nil(t, args)
	})
	t.Run("ArgBytes IN binary", func(t *testing.T) {
		s := NewSelect("a", "b").From("t1").Where(
			Column("a316", In.Bytes([]byte{66, 250, 67}, []byte(`Rust`), []byte("\xFB\xBF\xBF\xBF\xBF"))),
		).Interpolate()
		sStr, args, err := s.ToSQL()
		require.NoError(t, err)
		assert.Equal(t, "SELECT `a`, `b` FROM `t1` WHERE (`a316` IN (0x42fa43,'Rust',0xfbbfbfbfbf))", sStr)
		assert.Nil(t, args)
	})
	t.Run("ArgValue IN", func(t *testing.T) {
		s := NewSelect("a", "b").From("t1").Where(
			Column("a319", In.Value(MakeNullFloat64(3.141), MakeNullString("G'o"), MakeNullBytes([]byte{66, 250, 67}),
				MakeNullTime(now()), MakeNullBytes([]byte("x\x00y")))),
		).Interpolate()
		sStr, args, err := s.ToSQL()
		require.NoError(t, err)
		assert.Equal(t, "SELECT `a`, `b` FROM `t1` WHERE (`a319` IN (3.141,'G\\'o',0x42fa43,'2006-01-02 15:04:05',0x780079))", sStr)
		assert.Nil(t, args)
	})
	t.Run("ArgValue BETWEEN values", func(t *testing.T) {
		s := NewSelect("a", "b").From("t1").Where(
			Column("a319", Between.Value(MakeNullFloat64(3.141), MakeNullString("G'o"))),
		).Interpolate()
		sStr, args, err := s.ToSQL()
		require.NoError(t, err)
		assert.Equal(t, "SELECT `a`, `b` FROM `t1` WHERE (`a319` BETWEEN 3.141 AND 'G\\'o')", sStr)
		assert.Nil(t, args)
	})
}

func TestArguments_DriverValues(t *testing.T) {
	t.Parallel()
	args := Arguments{
		ArgInt64(1),
		ArgInt(2),
		ArgString("BlackForest"),
		ArgFloat64(3.14159),
		ArgBool(true),
	}
	assert.Exactly(t,
		[]driver.Value{int64(1), int64(2), "BlackForest", 3.14159, true},
		args.DriverValues())
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
	s := createRealSessionWithFixtures()

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
		assert.NoError(t, err)
		id, err := res.LastInsertId()
		assert.NoError(t, err)

		// Scan it back and check that all fields are of the correct validity and are
		// equal to the reference record
		nullTypeSet := &nullTypedRecord{}
		_, err = s.Select("*").From("null_types").Where(Expression("id = ?", ArgInt64(id))).Load(context.TODO(), nullTypeSet)
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
