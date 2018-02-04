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
	"fmt"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ fmt.Stringer = (*scannedColumn)(nil)

func TestColumnMap_Nil_Pointers(t *testing.T) {
	t.Parallel()

	cm := NewColumnMap(20)
	cm.
		Bool(nil).
		Byte(nil).
		Float64(nil).
		Int(nil).
		Int64(nil).
		NullBool(nil).
		NullFloat64(nil).
		NullInt64(nil).
		NullString(nil).
		NullTime(nil).
		String(nil).
		Time(nil).
		Uint(nil).
		Uint16(nil).
		Uint32(nil).
		Uint64(nil).
		Uint8(nil)

	assert.Exactly(t,
		"dml.MakeArgs(17).Null().Null().Null().Null().Null().Null().Null().Null().Null().Null().Null().Null().Null().Null().Null().Null().Null()",
		cm.arguments.GoString())
}

func TestScannedColumn_String(t *testing.T) {
	t.Parallel()
	sc := scannedColumn{
		field:   'b',
		bool:    true,
		int64:   -17,
		float64: 3.14159,
		string:  "",
		time:    now(),
		byte:    []byte(`Håi`),
	}
	assert.Exactly(t, "true", sc.String())
	sc.field = 'i'
	assert.Exactly(t, "-17", sc.String())
	sc.field = 'f'
	assert.Exactly(t, "3.14159", sc.String())
	sc.field = 's'
	assert.Exactly(t, "", sc.String())
	sc.field = 't'
	assert.Exactly(t, "2006-01-02 15:04:05.000000002 -0400 UTC-4", sc.String())
	sc.field = 'y'
	assert.Exactly(t, "Håi", sc.String())
	sc.field = 0
	assert.Exactly(t, "Field Type '\\x00' not supported", sc.String())
}

func TestScannedColumn_Scan(t *testing.T) {
	t.Parallel()
	sc := scannedColumn{}

	require.NoError(t, sc.Scan(int64(4711)))
	assert.Exactly(t, "4711", sc.String())

	require.NoError(t, sc.Scan(int(4711)))
	assert.Exactly(t, "4711", sc.String())

	require.NoError(t, sc.Scan(float32(47.11)))
	assert.Exactly(t, "47.11000061035156", sc.String())

	require.NoError(t, sc.Scan(float64(47.11)))
	assert.Exactly(t, "47.11", sc.String())

	require.NoError(t, sc.Scan(true))
	assert.Exactly(t, "true", sc.String())
	require.NoError(t, sc.Scan(false))
	assert.Exactly(t, "false", sc.String())

	require.NoError(t, sc.Scan([]byte(`@`)))
	assert.Exactly(t, "@", sc.String())

	require.NoError(t, sc.Scan(`@`))
	assert.Exactly(t, "@", sc.String())

	require.NoError(t, sc.Scan(now()))
	assert.Exactly(t, "2006-01-02 15:04:05.000000002 -0400 UTC-4", sc.String())

	require.NoError(t, sc.Scan(nil))
	assert.Exactly(t, "<nil>", sc.String())

	err := sc.Scan(uint8(1))
	require.True(t, errors.Is(err, errors.NotSupported), "Should be error kind NotSupported")

}

func TestColumnMap_Scan_Empty_Bytes(t *testing.T) {
	t.Parallel()

	cm := NewColumnMap(0, "SomeColumn")
	cm.index = 0
	cm.scanCol = make([]scannedColumn, 1)
	cm.scanCol[0].field = 'y'

	t.Run("Bool", func(t *testing.T) {
		var v bool
		assert.EqualError(t, cm.Bool(&v).Err(), "[dml] Column \"SomeColumn\": strconv.ParseBool: parsing \"\": invalid syntax")
		cm.scanErr = nil
	})
	t.Run("Byte", func(t *testing.T) {
		var v []byte
		assert.NoError(t, cm.Byte(&v).Err())
		assert.Nil(t, v)
		cm.scanErr = nil
	})
	t.Run("Float64", func(t *testing.T) {
		var v float64
		assert.NoError(t, cm.Float64(&v).Err())
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("Int", func(t *testing.T) {
		var v int
		assert.NoError(t, cm.Int(&v).Err())
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("Int64", func(t *testing.T) {
		var v int64
		assert.NoError(t, cm.Int64(&v).Err())
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("NullBool", func(t *testing.T) {
		var v NullBool
		assert.NoError(t, cm.NullBool(&v).Err())
		assert.False(t, v.Valid)
		cm.scanErr = nil
	})
	t.Run("NullFloat64", func(t *testing.T) {
		var v NullFloat64
		assert.NoError(t, cm.NullFloat64(&v).Err())
		assert.False(t, v.Valid)
		cm.scanErr = nil
	})
	t.Run("NullInt64", func(t *testing.T) {
		var v NullInt64
		assert.NoError(t, cm.NullInt64(&v).Err())
		assert.False(t, v.Valid)
		cm.scanErr = nil
	})
	t.Run("NullString", func(t *testing.T) {
		var v NullString
		assert.NoError(t, cm.NullString(&v).Err())
		assert.False(t, v.Valid)
		cm.scanErr = nil
	})
	t.Run("NullTime", func(t *testing.T) {
		var v NullTime
		assert.NoError(t, cm.NullTime(&v).Err())
		assert.False(t, v.Valid)
		cm.scanErr = nil
	})
	t.Run("String", func(t *testing.T) {
		var v string
		assert.NoError(t, cm.String(&v).Err())
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("Time", func(t *testing.T) {
		var v time.Time
		assert.EqualError(t, cm.Time(&v).Err(),
			"[dml] Column \"SomeColumn\": invalid time string: \"\"")
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("Uint", func(t *testing.T) {
		var v uint
		assert.EqualError(t, cm.Uint(&v).Err(), "[dml] Column \"SomeColumn\": strconv.ParseUint: parsing \"\": invalid syntax")
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("Uint8", func(t *testing.T) {
		var v uint8
		assert.EqualError(t, cm.Uint8(&v).Err(), "[dml] Column \"SomeColumn\": strconv.ParseUint: parsing \"\": invalid syntax")
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("Uint16", func(t *testing.T) {
		var v uint16
		assert.EqualError(t, cm.Uint16(&v).Err(), "[dml] Column \"SomeColumn\": strconv.ParseUint: parsing \"\": invalid syntax")
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("Uint32", func(t *testing.T) {
		var v uint32
		assert.EqualError(t, cm.Uint32(&v).Err(), "[dml] Column \"SomeColumn\": strconv.ParseUint: parsing \"\": invalid syntax")
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("Uint64", func(t *testing.T) {
		var v uint64
		assert.EqualError(t, cm.Uint64(&v).Err(), "[dml] Column \"SomeColumn\": strconv.ParseUint: parsing \"\": invalid syntax")
		assert.Empty(t, v)
		cm.scanErr = nil
	})
}
