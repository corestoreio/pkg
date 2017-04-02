package dbr

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/corestoreio/csfw/util/null"
	"github.com/stretchr/testify/assert"
)

func TestNullStringFrom(t *testing.T) {
	assert.Equal(t, "product", null.StringFrom("product").String)
	assert.True(t, null.StringFrom("product").Valid)
	assert.False(t, null.StringFromPtr(nil).Valid)
	assert.True(t, null.StringFrom("").Valid)
	v, err := null.StringFrom("product").Value()
	assert.NoError(t, err)
	assert.Equal(t, "product", v)
}

func TestNewNullInt64(t *testing.T) {
	assert.EqualValues(t, 1257894000, null.Int64From(1257894000).Int64)
	assert.True(t, null.Int64From(1257894000).Valid)
	assert.True(t, null.Int64From(0).Valid)
	assert.False(t, null.Int64FromPtr(nil).Valid)
	v, err := null.Int64From(1257894000).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, 1257894000, v)
}

func TestNewNullFloat64(t *testing.T) {
	var test = 1257894000.93445000001
	assert.Equal(t, test, null.Float64From(test).Float64)
	assert.True(t, null.Float64From(test).Valid)
	assert.True(t, null.Float64From(0).Valid)
	assert.False(t, null.Float64FromPtr(nil).Valid)
	v, err := null.Float64From(test).Value()
	assert.NoError(t, err)
	assert.Equal(t, test, v)
}

func TestNewNullTime(t *testing.T) {
	var test = time.Now()
	assert.Equal(t, test, null.TimeFrom(test).Time)
	assert.True(t, null.TimeFrom(test).Valid)
	assert.True(t, null.TimeFrom(time.Time{}).Valid)
	assert.False(t, null.TimeFromPtr(nil).Valid)
	v, err := null.TimeFrom(test).Value()
	assert.NoError(t, err)
	assert.Equal(t, test, v)
}

func TestNewNullBool(t *testing.T) {

	assert.Equal(t, true, null.BoolFrom(true).Bool)
	assert.True(t, null.BoolFrom(true).Valid)
	assert.True(t, null.BoolFrom(false).Valid)
	assert.False(t, null.BoolFromPtr(nil).Valid)
	v, err := null.BoolFrom(true).Value()
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
		res, err := s.InsertInto("null_types").Columns("string_val", "int64_val", "float64_val", "time_val", "bool_val").Record(test.record).Exec()
		assert.NoError(t, err)
		id, err := res.LastInsertId()
		assert.NoError(t, err)

		// Scan it back and check that all fields are of the correct validity and are
		// equal to the reference record
		nullTypeSet := &nullTypedRecord{}
		err = s.Select("*").From("null_types").Where(ConditionRaw("id = ?", ArgInt64(id))).LoadStruct(nullTypeSet)
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
		StringVal:  null.String{sql.NullString{String: "wow", Valid: true}},
		Int64Val:   null.Int64{sql.NullInt64{Int64: 42, Valid: true}},
		Float64Val: null.Float64{sql.NullFloat64{Float64: 1.618, Valid: true}},
		TimeVal:    null.Time{Time: time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC), Valid: true},
		BoolVal:    null.Bool{sql.NullBool{Bool: true, Valid: true}},
	}
}
