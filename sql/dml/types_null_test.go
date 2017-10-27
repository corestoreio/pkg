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
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecimal_Select_Integration(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)

	rec := newNullTypedRecordWithData()
	in := s.InsertInto("dml_null_types").
		AddColumns("id", "string_val", "int64_val", "float64_val", "time_val", "bool_val", "decimal_val")

	res, err := in.BindRecord(rec).Exec(context.TODO())
	require.NoError(t, err)
	id, err := res.LastInsertId()
	require.NoError(t, err)
	assert.Exactly(t, int64(2), id)

	nullTypeSet := &nullTypedRecord{}
	dec := Decimal{Precision: 12345, Scale: 3, Valid: true}

	sel := s.SelectFrom("dml_null_types").Star().Where(
		Column("decimal_val").Decimal(dec),
	)

	rc, err := sel.Load(context.TODO(), nullTypeSet)
	require.NoError(t, err)
	assert.Exactly(t, uint64(1), rc)

	assert.Exactly(t, rec, nullTypeSet)
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
		res, err := s.InsertInto("dml_null_types").
			AddColumns("string_val", "int64_val", "float64_val", "time_val", "bool_val", "decimal_val").
			BindRecord(test.record).Exec(context.TODO())
		require.NoError(t, err)
		id, err := res.LastInsertId()
		require.NoError(t, err)

		// Scan it back and check that all fields are of the correct validity and are
		// equal to the reference record
		nullTypeSet := &nullTypedRecord{}
		_, err = s.SelectFrom("dml_null_types").Star().Where(
			Expr("id = ?").Int64(id),
		).Load(context.TODO(), nullTypeSet)
		require.NoError(t, err)

		assert.Equal(t, test.record, nullTypeSet)
		assert.Equal(t, test.valid, nullTypeSet.StringVal.Valid)
		assert.Equal(t, test.valid, nullTypeSet.Int64Val.Valid)
		assert.Equal(t, test.valid, nullTypeSet.Float64Val.Valid)
		assert.Equal(t, test.valid, nullTypeSet.TimeVal.Valid)
		assert.Equal(t, test.valid, nullTypeSet.BoolVal.Valid)
		assert.Equal(t, test.valid, nullTypeSet.DecimalVal.Valid)

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
			expectedJSON: []byte("{\"ID\":0,\"StringVal\":null,\"Int64Val\":null,\"Float64Val\":null,\"TimeVal\":null,\"BoolVal\":null,\"DecimalVal\":null}"),
		},
		{
			record:       newNullTypedRecordWithData(),
			expectedJSON: []byte("{\"ID\":2,\"StringVal\":\"wow\",\"Int64Val\":42,\"Float64Val\":1.618,\"TimeVal\":\"2009-01-03T18:15:05Z\",\"BoolVal\":true,\"DecimalVal\":12.345}"),
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
		DecimalVal: Decimal{Precision: 12345, Scale: 3, Valid: true},
	}
}
