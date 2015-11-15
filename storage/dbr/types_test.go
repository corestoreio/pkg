package dbr

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestNewNullString(t *testing.T) {
	assert.Equal(t, "product", NewNullString("product").String)
	assert.True(t, NewNullString("product", false).Valid)
	assert.False(t, NewNullString("", false).Valid)
	v, err := NewNullString("product", false).Value()
	assert.NoError(t, err)
	assert.Equal(t, "product", v)
}

func TestNewNullInt64(t *testing.T) {
	assert.EqualValues(t, 1257894000, NewNullInt64(1257894000).Int64)
	assert.True(t, NewNullInt64(1257894000, false).Valid)
	assert.False(t, NewNullInt64(0, false).Valid)
	v, err := NewNullInt64(1257894000, false).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, 1257894000, v)
}

func TestNewNullFloat64(t *testing.T) {
	var test float64 = 1257894000.93445000001
	assert.Equal(t, test, NewNullFloat64(test).Float64)
	assert.True(t, NewNullFloat64(test, false).Valid)
	assert.False(t, NewNullFloat64(0, false).Valid)
	v, err := NewNullFloat64(test, false).Value()
	assert.NoError(t, err)
	assert.Equal(t, test, v)
}

func TestNewNullTime(t *testing.T) {
	var test = time.Now()
	assert.Equal(t, test, NewNullTime(test).Time)
	assert.True(t, NewNullTime(test, false).Valid)
	assert.False(t, NewNullTime(time.Time{}, false).Valid)
	v, err := NewNullTime(test, false).Value()
	assert.NoError(t, err)
	assert.Equal(t, test, v)
}

func TestNewNullBool(t *testing.T) {

	assert.Equal(t, true, NewNullBool(true, true).Bool)
	assert.True(t, NewNullBool(true, false).Valid)
	assert.False(t, NewNullBool(false, false).Valid)
	v, err := NewNullBool(true, false).Value()
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
			record: &nullTypedRecord{},
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
		err = s.Select("*").From("null_types").Where("id = ?", id).LoadStruct(nullTypeSet)
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

func newNullTypedRecordWithData() *nullTypedRecord {
	return &nullTypedRecord{
		StringVal:  NullString{sql.NullString{String: "wow", Valid: true}},
		Int64Val:   NullInt64{sql.NullInt64{Int64: 42, Valid: true}},
		Float64Val: NullFloat64{sql.NullFloat64{Float64: 1.618, Valid: true}},
		TimeVal:    NullTime{mysql.NullTime{Time: time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC), Valid: true}},
		BoolVal:    NullBool{sql.NullBool{Bool: true, Valid: true}},
	}
}
