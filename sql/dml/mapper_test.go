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
	"bytes"
	"encoding"
	"fmt"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

var _ fmt.Stringer = (*scannedColumn)(nil)

var (
	_ encoding.TextMarshaler     = (*textBinaryEncoder)(nil)
	_ encoding.TextUnmarshaler   = (*textBinaryEncoder)(nil)
	_ encoding.BinaryMarshaler   = (*textBinaryEncoder)(nil)
	_ encoding.BinaryUnmarshaler = (*textBinaryEncoder)(nil)
)

type textBinaryEncoder struct {
	data []byte
}

func (be textBinaryEncoder) MarshalBinary() (data []byte, err error) {
	if bytes.Equal(be.data, []byte(`error`)) {
		return nil, errors.DecryptionFailed.Newf("decryption failed test error")
	}
	return be.data, nil
}

func (be *textBinaryEncoder) UnmarshalBinary(data []byte) error {
	if bytes.Equal(data, []byte(`error`)) {
		return errors.Empty.Newf("test error empty")
	}
	be.data = append(be.data, data...)
	return nil
}

func (be textBinaryEncoder) MarshalText() (text []byte, err error) {
	if bytes.Equal(be.data, []byte(`error`)) {
		return nil, errors.DecryptionFailed.Newf("internal validation failed test error")
	}
	return be.data, nil
}

func (be *textBinaryEncoder) UnmarshalText(text []byte) error {
	if bytes.Equal(text, []byte(`error`)) {
		return errors.Empty.Newf("test error empty")
	}
	be.data = append(be.data, text...)
	return nil
}

func TestColumnMap_BinaryText(t *testing.T) {
	cm := NewColumnMap(1)

	assert.NoError(t, cm.Binary(&textBinaryEncoder{data: []byte(`BinaryTest`)}).Err())
	assert.Exactly(t,
		[]any{[]byte{0x42, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x54, 0x65, 0x73, 0x74}},
		expandInterfaces(cm.args))
	assert.NoError(t, cm.Text(&textBinaryEncoder{data: []byte(`TextTest`)}).Err())
	assert.Exactly(t,
		[]any{[]byte{0x42, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x54, 0x65, 0x73, 0x74}, []byte{0x54, 0x65, 0x78, 0x74, 0x54, 0x65, 0x73, 0x74}},
		expandInterfaces(cm.args))

	cm.CheckValidUTF8 = true
	err := cm.Text(&textBinaryEncoder{data: []byte("\xc0\x80")}).Err()
	assert.Error(t, err)
}

func TestColumnMap_Nil_Pointers(t *testing.T) {
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

	assert.Exactly(t, []any{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		expandInterfaces(cm.args))
}

func TestScannedColumn_String(t *testing.T) {
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
	sc := scannedColumn{}

	assert.NoError(t, sc.Scan(int64(4711)))
	assert.Exactly(t, "4711", sc.String())

	assert.NoError(t, sc.Scan(int(4711)))
	assert.Exactly(t, "4711", sc.String())

	assert.NoError(t, sc.Scan(float32(47.11)))
	assert.Exactly(t, "47.11000061035156", sc.String())

	assert.NoError(t, sc.Scan(float64(47.11)))
	assert.Exactly(t, "47.11", sc.String())

	assert.NoError(t, sc.Scan(true))
	assert.Exactly(t, "true", sc.String())
	assert.NoError(t, sc.Scan(false))
	assert.Exactly(t, "false", sc.String())

	assert.NoError(t, sc.Scan([]byte(`@`)))
	assert.Exactly(t, "@", sc.String())

	assert.NoError(t, sc.Scan(`@`))
	assert.Exactly(t, "@", sc.String())

	assert.NoError(t, sc.Scan(now()))
	assert.Exactly(t, "2006-01-02 15:04:05.000000002 -0400 UTC-4", sc.String())

	assert.NoError(t, sc.Scan(nil))
	assert.Exactly(t, "<nil>", sc.String())

	err := sc.Scan(uint8(1))
	assert.Error(t, err)
}

func TestColumnMap_Scan_Empty_Bytes(t *testing.T) {
	cm := NewColumnMap(0, "SomeColumn")
	cm.index = 0
	cm.scanCol = make([]scannedColumn, 1)
	cm.scanCol[0].field = 'y'

	t.Run("Bool", func(t *testing.T) {
		var v bool
		assert.NoError(t, cm.Bool(&v).Err())
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
		var v null.Bool
		assert.NoError(t, cm.NullBool(&v).Err())
		assert.False(t, v.Valid)
		cm.scanErr = nil
	})
	t.Run("NullFloat64", func(t *testing.T) {
		var v null.Float64
		assert.NoError(t, cm.NullFloat64(&v).Err())
		assert.False(t, v.Valid)
		cm.scanErr = nil
	})
	t.Run("NullInt64", func(t *testing.T) {
		var v null.Int64
		assert.NoError(t, cm.NullInt64(&v).Err())
		assert.False(t, v.Valid)
		cm.scanErr = nil
	})
	t.Run("NullString", func(t *testing.T) {
		var v null.String
		assert.NoError(t, cm.NullString(&v).Err())
		assert.False(t, v.Valid)
		cm.scanErr = nil
	})
	t.Run("NullTime", func(t *testing.T) {
		var v null.Time
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
		assert.Error(t, cm.Time(&v).Err())
		assert.True(t, v.IsZero(), "Got: %s", v.String())
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
