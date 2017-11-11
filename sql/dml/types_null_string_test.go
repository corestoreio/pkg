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
	"bytes"
	"database/sql/driver"
	"encoding"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ fmt.GoStringer             = (*NullString)(nil)
	_ json.Marshaler             = (*NullString)(nil)
	_ json.Unmarshaler           = (*NullString)(nil)
	_ encoding.BinaryMarshaler   = (*NullString)(nil)
	_ encoding.BinaryUnmarshaler = (*NullString)(nil)
	_ encoding.TextMarshaler     = (*NullString)(nil)
	_ encoding.TextUnmarshaler   = (*NullString)(nil)
	_ gob.GobEncoder             = (*NullString)(nil)
	_ gob.GobDecoder             = (*NullString)(nil)
	_ driver.Valuer              = (*NullString)(nil)
	_ proto.Marshaler            = (*NullString)(nil)
	_ proto.Unmarshaler          = (*NullString)(nil)
	_ proto.Sizer                = (*NullString)(nil)
	_ protoMarshalToer           = (*NullString)(nil)
)
var (
	stringJSON      = []byte(`"test"`)
	blankStringJSON = []byte(`""`)
	nullStringJSON  = []byte(`{"NullString":"test","Valid":true}`)

	nullJSON    = []byte(sqlStrNullLC)
	invalidJSON = []byte(`:)`)
)

func TestStringFrom(t *testing.T) {
	t.Parallel()
	str := MakeNullString("test")
	assertStr(t, str, "MakeNullString() string")
	assert.Exactly(t, 4, str.Size())

	zero := MakeNullString("")
	if !zero.Valid {
		t.Error("MakeNullString(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, 0, zero.Size())
}

func TestNullString_JsonUnmarshal(t *testing.T) {
	t.Parallel()
	var str NullString
	maybePanic(json.Unmarshal(stringJSON, &str))
	assertStr(t, str, "string json")

	var ns NullString
	maybePanic(json.Unmarshal(nullStringJSON, &ns))
	assertStr(t, ns, "sql.NullString json")

	var blank NullString
	maybePanic(json.Unmarshal(blankStringJSON, &blank))
	if !blank.Valid {
		t.Error("blank string should be valid")
	}

	var null NullString
	maybePanic(json.Unmarshal(nullJSON, &null))
	assertNullStr(t, null, "null json")

	var badType NullString
	err := json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullStr(t, badType, "wrong type json")

	var invalid NullString
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullStr(t, invalid, "invalid json")
}

func TestNullString_TextUnmarshal(t *testing.T) {
	t.Parallel()
	var str NullString
	err := str.UnmarshalText([]byte("test"))
	maybePanic(err)
	assertStr(t, str, "UnmarshalText() string")

	var null NullString
	err = null.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullStr(t, null, "UnmarshalText() empty string")

	var iv NullString
	err = iv.UnmarshalText([]byte{0x44, 0xff, 0x01})
	assert.True(t, errors.IsNotValid(err), "%+v", err)
}

func TestNullString_MarshalText(t *testing.T) {
	t.Parallel()
	str := MakeNullString("test")
	data, err := json.Marshal(str)
	maybePanic(err)
	assertJSONEquals(t, data, `"test"`, "non-empty json marshal")
	data, err = str.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "test", "non-empty text marshal")

	// empty values should be encoded as an empty string
	zero := MakeNullString("")
	data, err = json.Marshal(zero)
	maybePanic(err)
	assertJSONEquals(t, data, `""`, "empty json marshal")
	data, err = zero.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "string marshal text")

	zero.Valid = false
	data, err = zero.MarshalText()
	maybePanic(err)
	assert.Nil(t, data)
}

func TestNullString_BinaryEncoding(t *testing.T) {
	t.Parallel()
	runner := func(b NullString, want []byte) func(*testing.T) {
		return func(t *testing.T) {
			data, err := b.GobEncode()
			require.NoError(t, err)
			require.Exactly(t, want, data, t.Name()+": GobEncode")
			data, err = b.MarshalBinary()
			require.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": MarshalBinary")
			data, err = b.Marshal()
			require.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": Marshal")

			var decoded NullString
			require.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("HelloWorld", runner(MakeNullString("HelloWorld"), []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0xef, 0xa3, 0xbf, 0x57, 0x6f, 0x72, 0x6c, 0x64}))
	t.Run("null", runner(NullString{}, nil))
}

func TestNullString_MarshalTo(t *testing.T) {
	t.Parallel()
	str := MakeNullString("HelloWorld")
	var buf4 [4]byte
	n, err := str.MarshalTo(buf4[:])
	maybePanic(err)
	assert.Exactly(t, 4, n)
	assert.Exactly(t, []byte(`Hell`), buf4[:])

	bufFit := make([]byte, str.Size())
	n, err = str.MarshalTo(bufFit)
	maybePanic(err)
	assert.Exactly(t, 13, n)
	assert.Exactly(t, []byte(`HelloWorld`), bufFit)
}

func TestStringPointer(t *testing.T) {
	t.Parallel()
	str := MakeNullString("test")
	ptr := str.Ptr()
	if *ptr != "test" {
		t.Errorf("bad %s string: %#v ≠ %s\n", "pointer", ptr, "test")
	}

	null := MakeNullString("", false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s string: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestStringIsZero(t *testing.T) {
	t.Parallel()
	str := MakeNullString("test")
	if str.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	blank := MakeNullString("")
	if blank.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	empty := MakeNullString("", true)
	if empty.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestStringSetValid(t *testing.T) {
	t.Parallel()
	change := MakeNullString("", false)
	assertNullStr(t, change, "SetValid()")
	change.SetValid("test")
	assertStr(t, change, "SetValid()")
}

func TestStringScan(t *testing.T) {
	t.Parallel()
	var str NullString
	err := str.Scan("test")
	maybePanic(err)
	assertStr(t, str, "scanned string")

	var null NullString
	err = null.Scan(nil)
	maybePanic(err)
	assertNullStr(t, null, "scanned null")
}

func TestString_GoString(t *testing.T) {
	t.Parallel()
	s := MakeNullString("test", true)
	assert.Exactly(t, "dml.MakeNullString(`test`)", s.GoString())

	s = MakeNullString("test", false)
	assert.Exactly(t, "dml.NullString{}", s.GoString())

	s = MakeNullString("te`st", true)
	gsWant := []byte("dml.MakeNullString(`te`+\"`\"+`st`)")
	if !bytes.Equal(gsWant, []byte(s.GoString())) {
		t.Errorf("Have: %#v Want: %v", s.GoString(), string(gsWant))
	}
}

func assertStr(t *testing.T, s NullString, from string) {
	if s.String != "test" {
		t.Errorf("bad %s string: %s ≠ %s\n", from, s.String, "test")
	}
	if !s.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullStr(t *testing.T, s NullString, from string) {
	if s.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}

func TestNullStringFrom(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "product", MakeNullString("product").String)
	assert.True(t, MakeNullString("product").Valid)
	//assert.False(t, NullStringFromPtr(nil).Valid)
	assert.True(t, MakeNullString("").Valid)
	v, err := MakeNullString("product").Value()
	assert.NoError(t, err)
	assert.Equal(t, "product", v)
}
