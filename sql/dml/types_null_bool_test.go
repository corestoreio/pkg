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
	"database/sql/driver"
	"encoding"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	boolJSON     = []byte(`true`)
	nullBoolJSON = []byte(`{"NullBool":true,"Valid":true}`)
)

var (
	_ fmt.GoStringer             = (*NullBool)(nil)
	_ fmt.Stringer               = (*NullBool)(nil)
	_ json.Marshaler             = (*NullBool)(nil)
	_ json.Unmarshaler           = (*NullBool)(nil)
	_ encoding.BinaryMarshaler   = (*NullBool)(nil)
	_ encoding.BinaryUnmarshaler = (*NullBool)(nil)
	_ encoding.TextMarshaler     = (*NullBool)(nil)
	_ encoding.TextUnmarshaler   = (*NullBool)(nil)
	_ gob.GobEncoder             = (*NullBool)(nil)
	_ gob.GobDecoder             = (*NullBool)(nil)
	_ driver.Valuer              = (*NullBool)(nil)
	_ proto.Marshaler            = (*NullBool)(nil)
	_ proto.Unmarshaler          = (*NullBool)(nil)
	_ proto.Sizer                = (*NullBool)(nil)
	_ protoMarshalToer           = (*NullBool)(nil)
)

func TestMakeNullBool(t *testing.T) {
	t.Parallel()
	b := MakeNullBool(true)
	assertBool(t, b, "MakeNullBool()")
	assert.Exactly(t, "true", b.String())

	zero := MakeNullBool(false)
	if !zero.Valid {
		t.Error("MakeNullBool(false)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "false", zero.String())
	assert.Exactly(t, 1, zero.Size())
	assert.Exactly(t, "null", NullBool{}.String())
	assert.Exactly(t, 0, NullBool{}.Size())
}

func TestNullBool_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	var b NullBool
	err := json.Unmarshal(boolJSON, &b)
	maybePanic(err)
	assertBool(t, b, "bool json")

	var nb NullBool
	err = json.Unmarshal(nullBoolJSON, &nb)
	maybePanic(err)
	assertBool(t, nb, "sq.NullBool json")

	var null NullBool
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullBool(t, null, "null json")

	var badType NullBool
	err = json.Unmarshal(intJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullBool(t, badType, "wrong type json")

	var invalid NullBool
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
}

func TestNullBool_UnmarshalText(t *testing.T) {
	t.Parallel()

	var b NullBool
	err := b.UnmarshalText([]byte("true"))
	maybePanic(err)
	assertBool(t, b, "UnmarshalText() bool")

	var zero NullBool
	err = zero.UnmarshalText([]byte("false"))
	maybePanic(err)
	assertFalseBool(t, zero, "UnmarshalText() false")

	var blank NullBool
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullBool(t, blank, "UnmarshalText() empty bool")

	var null NullBool
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullBool(t, null, `UnmarshalText() "null"`)

	var invalid NullBool
	err = invalid.UnmarshalText([]byte(":D"))
	if err == nil {
		panic("err should not be nil")
	}
	assertNullBool(t, invalid, "invalid json")
}

func TestNullBool_JsonMarshal(t *testing.T) {
	t.Parallel()

	b := MakeNullBool(true)
	data, err := json.Marshal(b)
	maybePanic(err)
	assertJSONEquals(t, data, "true", "non-empty json marshal")

	zero := MakeNullBool(false, true)
	data, err = json.Marshal(zero)
	maybePanic(err)
	assertJSONEquals(t, data, "false", "zero json marshal")

	// invalid values should be encoded as null
	null := MakeNullBool(false, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullBool_MarshalText(t *testing.T) {
	t.Parallel()

	b := MakeNullBool(true)
	data, err := b.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "true", "non-empty text marshal")

	zero := MakeNullBool(false, true)
	data, err = zero.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "false", "zero text marshal")

	// invalid values should be encoded as null
	null := MakeNullBool(false, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullBool_BinaryEncoding(t *testing.T) {
	t.Parallel()
	runner := func(b NullBool, want []byte) func(*testing.T) {
		return func(t *testing.T) {
			data, err := b.GobEncode()
			require.NoError(t, err)
			assert.Exactly(t, want, data, "GobEncode")
			data, err = b.MarshalBinary()
			require.NoError(t, err)
			assert.Exactly(t, want, data, "MarshalBinary")
			data, err = b.Marshal()
			require.NoError(t, err)
			assert.Exactly(t, want, data, "Marshal")

			var decoded NullBool
			require.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("true", runner(MakeNullBool(true), []byte{1}))
	t.Run("false", runner(MakeNullBool(false), []byte{0}))
	t.Run("null", runner(NullBool{}, nil))
}

func TestNullBool_BinaryDecoding(t *testing.T) {
	t.Parallel()
	runner := func(data []byte, want NullBool) func(*testing.T) {
		return func(t *testing.T) {
			var have NullBool
			require.NoError(t, have.GobDecode(data), "GobDecode")
			assert.Exactly(t, want, have, "GobDecode")
			require.NoError(t, have.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, want, have, "UnmarshalBinary")
			require.NoError(t, have.Unmarshal(data), "Unmarshal")
			assert.Exactly(t, want, have, "Unmarshal")
		}
	}
	t.Run("true", runner([]byte{1}, MakeNullBool(true)))
	t.Run("false", runner([]byte{0}, MakeNullBool(false)))
	t.Run("null", runner(nil, NullBool{}))
	t.Run("junk", runner([]byte{2, 1, 3}, NullBool{}))
}

func TestBoolPointer(t *testing.T) {
	t.Parallel()

	b := MakeNullBool(true)
	ptr := b.Ptr()
	if !*ptr {
		t.Errorf("bad %s bool: %#v ≠ %v\n", "pointer", ptr, true)
	}

	null := MakeNullBool(false, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s bool: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestBoolIsZero(t *testing.T) {
	t.Parallel()

	b := MakeNullBool(true)
	if b.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := MakeNullBool(false, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeNullBool(false, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestBoolSetValid(t *testing.T) {
	t.Parallel()

	change := MakeNullBool(false, false)
	assertNullBool(t, change, "SetValid()")
	change.SetValid(true)
	assertBool(t, change, "SetValid()")
}

func TestBoolScan(t *testing.T) {
	t.Parallel()

	var b NullBool
	err := b.Scan(true)
	maybePanic(err)
	assertBool(t, b, "scanned bool")

	var null NullBool
	err = null.Scan(nil)
	maybePanic(err)
	assertNullBool(t, null, "scanned null")
}

func assertBool(t *testing.T, b NullBool, from string) {
	if !b.Bool {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Bool, true)
	}
	if !b.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertFalseBool(t *testing.T, b NullBool, from string) {
	if b.Bool {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Bool, false)
	}
	if !b.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullBool(t *testing.T, b NullBool, from string) {
	if b.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
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
