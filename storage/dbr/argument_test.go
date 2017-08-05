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
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ fmt.Stringer = Op(0)

//var _ Argument = Int(0)
//var _ Argument = Int64(0)
//var _ Argument = Float64(0)
//var _ Argument = Bool(true)
//var _ Argument = (*BytesSlice)(nil)
//var _ Argument = (*Ints)(nil)
//var _ Argument = (*Int64s)(nil)
//var _ Argument = (*Float64s)(nil)
//var _ Argument = (*Times)(nil)
//var _ Argument = (*Strings)(nil)
//var _ Argument = (*NullString)(nil)
//var _ Argument = (*NullStrings)(nil)
//var _ Argument = (*NullFloat64)(nil)
//var _ Argument = (*NullFloat64s)(nil)
//var _ Argument = (*NullTime)(nil)
//var _ Argument = (*NullTimes)(nil)
//var _ Argument = (*NullInt64)(nil)
//var _ Argument = (*NullInt64s)(nil)
//var _ Argument = (*NullBool)(nil)
//var _ Argument = (*DriverValues)(nil)
//var _ driver.Valuer = (*Bool)(nil)
//var _ driver.Valuer = (*Int)(nil)
//var _ driver.Valuer = (*Int64)(nil)
//var _ driver.Valuer = (*Float64)(nil)
//var _ driver.Valuer = (*Time)(nil)
//var _ driver.Valuer = (*Bytes)(nil)
//var _ driver.Valuer = (*String)(nil)

//func TestDriverValuer(t *testing.T) {
//	t.Parallel()
//	tests := []struct {
//		have driver.Valuer
//		want driver.Value
//	}{
//		{Bool(true), true},
//		{Int(7), int64(7)},
//		{Int64(8), int64(8)},
//		{Float64(8.9), float64(8.9)},
//		{MakeTime(now()), now()},
//		{Bytes(`Go2`), []byte(`Go2`)},
//		{String(`Go2`), `Go2`},
//	}
//	for i, test := range tests {
//		v, err := test.have.Value()
//		if err != nil {
//			t.Fatalf("index %d with %+v", i, err)
//		}
//		assert.Exactly(t, test.want, v, "Index %d", i)
//	}
//}

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
		_, err = s.Select("*").From("null_types").Where(
			Expression("id = ?").Int64(id),
		).Load(context.TODO(), nullTypeSet)
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
