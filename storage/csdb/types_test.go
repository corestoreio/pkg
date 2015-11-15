// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package csdb_test

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestNewNullString(t *testing.T) {
	assert.Equal(t, "product", csdb.NewNullString("product").String)
	assert.True(t, csdb.NewNullString("product").Valid)
	assert.False(t, csdb.NewNullString(nil).Valid)
	assert.True(t, csdb.NewNullString("").Valid)
	v, err := csdb.NewNullString("product").Value()
	assert.NoError(t, err)
	assert.Equal(t, "product", v)
}

func TestNewNullInt64(t *testing.T) {
	assert.EqualValues(t, 1257894000, csdb.NewNullInt64(1257894000).Int64)
	assert.True(t, csdb.NewNullInt64(1257894000).Valid)
	assert.True(t, csdb.NewNullInt64(0).Valid)
	assert.False(t, csdb.NewNullInt64(nil).Valid)
	v, err := csdb.NewNullInt64(1257894000).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, 1257894000, v)
}

func TestNewNullFloat64(t *testing.T) {
	var test float64 = 1257894000.93445000001
	assert.Equal(t, test, csdb.NewNullFloat64(test).Float64)
	assert.True(t, csdb.NewNullFloat64(test).Valid)
	assert.True(t, csdb.NewNullFloat64(0).Valid)
	assert.False(t, csdb.NewNullFloat64(nil).Valid)
	v, err := csdb.NewNullFloat64(test).Value()
	assert.NoError(t, err)
	assert.Equal(t, test, v)
}

func TestNewNullTime(t *testing.T) {
	var test = time.Now()
	assert.Equal(t, test, csdb.NewNullTime(test).Time)
	assert.True(t, csdb.NewNullTime(test).Valid)
	assert.True(t, csdb.NewNullTime(time.Time{}).Valid)
	assert.False(t, csdb.NewNullTime(nil).Valid)
	v, err := csdb.NewNullTime(test).Value()
	assert.NoError(t, err)
	assert.Equal(t, test, v)
}

func TestNewNullBool(t *testing.T) {

	assert.Equal(t, true, csdb.NewNullBool(true).Bool)
	assert.True(t, csdb.NewNullBool(true).Valid)
	assert.True(t, csdb.NewNullBool(false).Valid)
	assert.False(t, csdb.NewNullBool(nil).Valid)
	v, err := csdb.NewNullBool(true).Value()
	assert.NoError(t, err)
	assert.Equal(t, true, v)
}

//func TestNullTypeScanning(t *testing.T) {
//	s := createRealSessionWithFixtures()
//
//	type nullTypeScanningTest struct {
//		record *nullTypedRecord
//		valid  bool
//	}
//
//	tests := []nullTypeScanningTest{
//		{
//			record: &nullTypedRecord{},
//			valid:  false,
//		},
//		{
//			record: newNullTypedRecordWithData(),
//			valid:  true,
//		},
//	}
//
//	for _, test := range tests {
//		// Create the record in the db
//		res, err := s.InsertInto("null_types").Columns("string_val", "int64_val", "float64_val", "time_val", "bool_val").Record(test.record).Exec()
//		assert.NoError(t, err)
//		id, err := res.LastInsertId()
//		assert.NoError(t, err)
//
//		// Scan it back and check that all fields are of the correct validity and are
//		// equal to the reference record
//		nullTypeSet := &nullTypedRecord{}
//		err = s.Select("*").From("null_types").Where("id = ?", id).LoadStruct(nullTypeSet)
//		assert.NoError(t, err)
//
//		assert.Equal(t, test.record, nullTypeSet)
//		assert.Equal(t, test.valid, nullTypeSet.StringVal.Valid)
//		assert.Equal(t, test.valid, nullTypeSet.Int64Val.Valid)
//		assert.Equal(t, test.valid, nullTypeSet.Float64Val.Valid)
//		assert.Equal(t, test.valid, nullTypeSet.TimeVal.Valid)
//		assert.Equal(t, test.valid, nullTypeSet.BoolVal.Valid)
//
//		nullTypeSet.StringVal.String = "newStringVal"
//		assert.NotEqual(t, test.record, nullTypeSet)
//	}
//}

func TestNullTypeJSONMarshal(t *testing.T) {
	type nullTypeJSONTest struct {
		record       *nullTypedRecord
		expectedJSON []byte
	}

	tests := []nullTypeJSONTest{
		{
			record:       &nullTypedRecord{},
			expectedJSON: []byte(`{"Id":0,"StringVal":null,"Int64Val":null,"Float64Val":null,"TimeVal":null,"BoolVal":null}`),
		},
		{
			record:       newNullTypedRecordWithData(),
			expectedJSON: []byte(`{"Id":0,"StringVal":"wow","Int64Val":42,"Float64Val":1.618,"TimeVal":"2009-01-03T18:15:05Z","BoolVal":true}`),
		},
	}

	for _, test := range tests {
		// Marshal the record
		rawJSON, err := json.Marshal(test.record)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedJSON, rawJSON)

		// Unmarshal it back
		newRecord := &nullTypedRecord{}
		err = json.Unmarshal([]byte(rawJSON), newRecord)
		assert.NoError(t, err)
		assert.Equal(t, test.record, newRecord)
	}
}

type nullTypedRecord struct {
	Id         int64
	StringVal  csdb.NullString
	Int64Val   csdb.NullInt64
	Float64Val csdb.NullFloat64
	TimeVal    csdb.NullTime
	BoolVal    csdb.NullBool
}

func newNullTypedRecordWithData() *nullTypedRecord {
	return &nullTypedRecord{
		StringVal:  csdb.NullString{sql.NullString{String: "wow", Valid: true}},
		Int64Val:   csdb.NullInt64{sql.NullInt64{Int64: 42, Valid: true}},
		Float64Val: csdb.NullFloat64{sql.NullFloat64{Float64: 1.618, Valid: true}},
		TimeVal:    csdb.NullTime{mysql.NullTime{Time: time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC), Valid: true}},
		BoolVal:    csdb.NullBool{sql.NullBool{Bool: true, Valid: true}},
	}
}
