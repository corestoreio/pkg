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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestColumnMap_Nil_Pointers(t *testing.T) {
	t.Parallel()

	args := MakeArgs(17)
	cm := newColumnMap(args)

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
		cm.Args.GoString())
}

func TestColumnMap_Scan_Empty_Bytes(t *testing.T) {
	t.Parallel()

	cm := newColumnMap(nil)
	cm.index = 0
	cm.current = []byte(nil)
	cm.columns = []string{"SomeColumn"}

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
		assert.EqualError(t, cm.Time(&v).Err(), "[dml] Column \"SomeColumn\": parsing time \"\" as \"2006-01-02T15:04:05.999999999Z07:00\": cannot parse \"\" as \"2006\"")
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("Uint", func(t *testing.T) {
		var v uint
		assert.NoError(t, cm.Uint(&v).Err())
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("Uint8", func(t *testing.T) {
		var v uint8
		assert.NoError(t, cm.Uint8(&v).Err())
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("Uint16", func(t *testing.T) {
		var v uint16
		assert.NoError(t, cm.Uint16(&v).Err())
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("Uint32", func(t *testing.T) {
		var v uint32
		assert.NoError(t, cm.Uint32(&v).Err())
		assert.Empty(t, v)
		cm.scanErr = nil
	})
	t.Run("Uint64", func(t *testing.T) {
		var v uint64
		assert.NoError(t, cm.Uint64(&v).Err())
		assert.Empty(t, v)
		cm.scanErr = nil
	})
}
