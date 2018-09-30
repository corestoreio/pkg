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

package null

import (
	"database/sql/driver"
	"encoding"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/gogo/protobuf/proto"
)

var (
	boolJSON     = []byte(`true`)
	nullBoolJSON = []byte(`{"Bool":true,"Valid":true}`)
)

var (
	_ fmt.GoStringer             = (*Bool)(nil)
	_ fmt.Stringer               = (*Bool)(nil)
	_ json.Marshaler             = (*Bool)(nil)
	_ json.Unmarshaler           = (*Bool)(nil)
	_ encoding.BinaryMarshaler   = (*Bool)(nil)
	_ encoding.BinaryUnmarshaler = (*Bool)(nil)
	_ encoding.TextMarshaler     = (*Bool)(nil)
	_ encoding.TextUnmarshaler   = (*Bool)(nil)
	_ gob.GobEncoder             = (*Bool)(nil)
	_ gob.GobDecoder             = (*Bool)(nil)
	_ driver.Valuer              = (*Bool)(nil)
	_ proto.Marshaler            = (*Bool)(nil)
	_ proto.Unmarshaler          = (*Bool)(nil)
	_ proto.Sizer                = (*Bool)(nil)
	_ protoMarshalToer           = (*Bool)(nil)
)

func TestMakeNullBool(t *testing.T) {

	b := MakeBool(true)
	assertBool(t, b, "MakeBool()")
	assert.Exactly(t, "true", b.String())

	zero := MakeBool(false)
	if !zero.Valid {
		t.Error("MakeBool(false)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "false", zero.String())
	assert.Exactly(t, 1, zero.Size())
	assert.Exactly(t, "null", Bool{}.String())
	assert.Exactly(t, 0, Bool{}.Size())
}

func TestNullBool_UnmarshalJSON(t *testing.T) {

	var b Bool
	err := json.Unmarshal(boolJSON, &b)
	maybePanic(err)
	assertBool(t, b, "bool json")

	var nb Bool
	err = json.Unmarshal(nullBoolJSON, &nb)
	maybePanic(err)
	assertBool(t, nb, "sq.Bool json")

	var null Bool
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullBool(t, null, "null json")

	var badType Bool
	err = json.Unmarshal(intJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullBool(t, badType, "wrong type json")

	var invalid Bool
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
}

func TestNullBool_UnmarshalText(t *testing.T) {

	var b Bool
	err := b.UnmarshalText([]byte("true"))
	maybePanic(err)
	assertBool(t, b, "UnmarshalText() bool")

	var zero Bool
	err = zero.UnmarshalText([]byte("false"))
	maybePanic(err)
	assertFalseBool(t, zero, "UnmarshalText() false")

	var blank Bool
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullBool(t, blank, "UnmarshalText() empty bool")

	var null Bool
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullBool(t, null, `UnmarshalText() "null"`)

	var invalid Bool
	err = invalid.UnmarshalText([]byte(":D"))
	if err == nil {
		panic("err should not be nil")
	}
	assertNullBool(t, invalid, "invalid json")
}

func TestNullBool_JsonMarshal(t *testing.T) {

	b := MakeBool(true)
	data, err := json.Marshal(b)
	maybePanic(err)
	assertJSONEquals(t, data, "true", "non-empty json marshal")

	zero := MakeBool(false)
	data, err = json.Marshal(zero)
	maybePanic(err)
	assertJSONEquals(t, data, "false", "zero json marshal")

	// invalid values should be encoded as null
	null := Bool{}
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullBool_MarshalText(t *testing.T) {

	b := MakeBool(true)
	data, err := b.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "true", "non-empty text marshal")

	zero := MakeBool(false)
	data, err = zero.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "false", "zero text marshal")

	// invalid values should be encoded as null
	null := MakeBool(false).SetNull()
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullBool_BinaryEncoding(t *testing.T) {

	runner := func(b Bool, want []byte) func(*testing.T) {
		return func(t *testing.T) {
			data, err := b.GobEncode()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, "GobEncode")
			data, err = b.MarshalBinary()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, "MarshalBinary")
			data, err = b.Marshal()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, "Marshal")

			var decoded Bool
			assert.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("true", runner(MakeBool(true), []byte{1}))
	t.Run("false", runner(MakeBool(false), []byte{0}))
	t.Run("null", runner(Bool{}, nil))
}

func TestNullBool_BinaryDecoding(t *testing.T) {

	runner := func(data []byte, want Bool) func(*testing.T) {
		return func(t *testing.T) {
			var have Bool
			assert.NoError(t, have.GobDecode(data), "GobDecode")
			assert.Exactly(t, want, have, "GobDecode")
			assert.NoError(t, have.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, want, have, "UnmarshalBinary")
			assert.NoError(t, have.Unmarshal(data), "Unmarshal")
			assert.Exactly(t, want, have, "Unmarshal")
		}
	}
	t.Run("true", runner([]byte{1}, MakeBool(true)))
	t.Run("false", runner([]byte{0}, MakeBool(false)))
	t.Run("null", runner(nil, Bool{}))
	t.Run("junk", runner([]byte{2, 1, 3}, Bool{}))
}

func TestBoolPointer(t *testing.T) {

	b := MakeBool(true)
	ptr := b.Ptr()
	if !*ptr {
		t.Errorf("bad %s bool: %#v ≠ %v\n", "pointer", ptr, true)
	}

	null := MakeBool(false).SetNull()
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s bool: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestBoolIsZero(t *testing.T) {

	b := MakeBool(true)
	if b.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := MakeBool(false).SetNull()
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeBool(false)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestBoolSetValid(t *testing.T) {

	change := MakeBool(false).SetNull()
	assertNullBool(t, change, "SetValid()")
	assertBool(t, change.SetValid(true), "SetValid()")
}

func TestBoolScan(t *testing.T) {

	var b Bool
	err := b.Scan(true)
	maybePanic(err)
	assertBool(t, b, "scanned bool")

	var null Bool
	err = null.Scan(nil)
	maybePanic(err)
	assertNullBool(t, null, "scanned null")
}

func assertBool(t *testing.T, b Bool, from string) {
	if !b.Bool {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Bool, true)
	}
	if !b.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertFalseBool(t *testing.T, b Bool, from string) {
	if b.Bool {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Bool, false)
	}
	if !b.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullBool(t *testing.T, b Bool, from string) {
	if b.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullBool(t *testing.T) {

	assert.Equal(t, true, MakeBool(true).Bool)
	assert.True(t, MakeBool(true).Valid)
	assert.True(t, MakeBool(false).Valid)
	v, err := MakeBool(true).Value()
	assert.NoError(t, err)
	assert.Equal(t, true, v)
}
