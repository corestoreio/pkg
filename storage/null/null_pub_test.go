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

package null_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecimal_Select_Integration(t *testing.T) {
	dbc := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, dbc)

	rec := newNullTypedRecordWithData()
	in := dbc.InsertInto("dml_null_types").
		AddColumns("id", "string_val", "int64_val", "float64_val", "time_val", "bool_val", "decimal_val")

	res, err := in.WithArgs().Record("", rec).ExecContext(context.TODO())
	require.NoError(t, err)
	id, err := res.LastInsertId()
	require.NoError(t, err)
	assert.Exactly(t, int64(2), id)

	nullTypeSet := &nullTypedRecord{}
	dec := null.Decimal{Precision: 12345, Scale: 3, Valid: true}

	sel := dbc.SelectFrom("dml_null_types").Star().Where(
		dml.Column("decimal_val").Decimal(dec),
	)

	rc, err := sel.WithArgs().Load(context.TODO(), nullTypeSet)
	require.NoError(t, err)
	assert.Exactly(t, uint64(1), rc)

	assert.Exactly(t, rec, nullTypeSet)
}

func TestNullTypeScanning(t *testing.T) {
	dbc := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, dbc)

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
		res, err := dbc.InsertInto("dml_null_types").
			AddColumns("string_val", "int64_val", "float64_val", "time_val", "bool_val", "decimal_val").
			WithArgs().Record("", test.record).ExecContext(context.TODO())
		require.NoError(t, err)
		id, err := res.LastInsertId()
		require.NoError(t, err)

		// Scan it back and check that all fields are of the correct validity and are
		// equal to the reference record
		nullTypeSet := &nullTypedRecord{}
		_, err = dbc.SelectFrom("dml_null_types").Star().Where(
			dml.Expr("id = ?").Int64(id),
		).WithArgs().Load(context.TODO(), nullTypeSet)
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

var _ dml.ColumnMapper = (*nullTypedRecord)(nil)

type nullTypedRecord struct {
	ID         int64
	StringVal  null.String
	Int64Val   null.Int64
	Float64Val null.Float64
	TimeVal    null.Time
	BoolVal    null.Bool
	DecimalVal null.Decimal
}

func (p *nullTypedRecord) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm.Int64(&p.ID).NullString(&p.StringVal).NullInt64(&p.Int64Val).NullFloat64(&p.Float64Val).NullTime(&p.TimeVal).NullBool(&p.BoolVal).Decimal(&p.DecimalVal).Err()
	}
	for cm.Next() {
		c := cm.Column()
		switch c {
		case "id":
			cm.Int64(&p.ID)
		case "string_val":
			cm.NullString(&p.StringVal)
		case "int64_val":
			cm.NullInt64(&p.Int64Val)
		case "float64_val":
			cm.NullFloat64(&p.Float64Val)
		case "time_val":
			cm.NullTime(&p.TimeVal)
		case "bool_val":
			cm.NullBool(&p.BoolVal)
		case "decimal_val":
			cm.Decimal(&p.DecimalVal)
		default:
			return errors.NotFound.Newf("[dml_test] Column %q not found", c)
		}
	}
	return cm.Err()
}
func newNullTypedRecordWithData() *nullTypedRecord {
	return &nullTypedRecord{
		ID:         2,
		StringVal:  null.String{String: "wow", Valid: true},
		Int64Val:   null.Int64{Int64: 42, Valid: true},
		Float64Val: null.Float64{Float64: 1.618, Valid: true},
		TimeVal:    null.Time{Time: time.Date(2009, 1, 3, 18, 15, 5, 0, time.UTC), Valid: true},
		BoolVal:    null.Bool{Bool: true, Valid: true},
		DecimalVal: null.Decimal{Precision: 12345, Scale: 3, Valid: true},
	}
}
