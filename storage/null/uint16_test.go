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
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

var (
	_ fmt.GoStringer             = (*Uint16)(nil)
	_ fmt.Stringer               = (*Uint16)(nil)
	_ json.Marshaler             = (*Uint16)(nil)
	_ json.Unmarshaler           = (*Uint16)(nil)
	_ encoding.BinaryMarshaler   = (*Uint16)(nil)
	_ encoding.BinaryUnmarshaler = (*Uint16)(nil)
	_ encoding.TextMarshaler     = (*Uint16)(nil)
	_ encoding.TextUnmarshaler   = (*Uint16)(nil)
	_ driver.Valuer              = (*Uint16)(nil)
	_ sql.Scanner                = (*Uint16)(nil)
)

var (
	uint16JSON     = []byte(fmt.Sprintf("%d", math.MaxUint16))
	nullUint16JSON = []byte(fmt.Sprintf(`{"Uint16":%d,"Valid":true}`, math.MaxUint16))
	float16JSON    = []byte(`1.23456`)
)

func TestMakeNullUint16(t *testing.T) {
	i := MakeUint16(math.MaxUint16)
	assertUint16(t, i, "MakeUint16()")

	zero := MakeUint16(0)
	if !zero.Valid {
		t.Error("MakeUint16(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "null", Uint16{}.String())
	assert.Exactly(t, 2, zero.Size())
	assert.Exactly(t, 4, MakeUint16(125).Size())
	assert.Exactly(t, 5, MakeUint16(128).Size())
	assert.Exactly(t, "0", zero.String())
	assert.Exactly(t, "65535", i.String(), "Want: %q", i.String())
	assert.Exactly(t, 0, Uint16{}.Size())
}

func TestUint16_GoString(t *testing.T) {
	tests := []struct {
		i16  Uint16
		want string
	}{
		{Uint16{}, "null.Uint16{}"},
		{MakeUint16(2), "null.MakeUint16(2)"},
	}
	for i, test := range tests {
		if have, want := fmt.Sprintf("%#v", test.i16), test.want; have != want {
			t.Errorf("%d: Have: %v Want: %v", i, have, want)
		}
	}
}

func TestMakeUint16FromByte(t *testing.T) {
	ui, err := MakeUint16FromByte([]byte(`65535`))
	assert.NoError(t, err)
	assert.Exactly(t, Uint16{Uint16: math.MaxUint16, Valid: true}, ui)
}

func TestNullUint16_JsonUnmarshal(t *testing.T) {
	var err error
	var i Uint16
	err = json.Unmarshal(uint16JSON, &i)
	assert.NoError(t, err, "%+v", err)
	assertUint16(t, i, "int16 json")

	var ni Uint16
	err = json.Unmarshal(nullUint16JSON, &ni)
	assert.NoError(t, err, "%+v", err)
	assertUint16(t, ni, "null.Uint16 json")

	var null Uint16
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err, "%+v", err)
	assertNullUint16(t, null, "null json")

	var badType Uint16
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint16(t, badType, "wrong type json")

	var invalid Uint16
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint16(t, invalid, "invalid json")
}

func TestNullUint16_JsonUnmarshalNonIntegerNumber(t *testing.T) {
	var i Uint16
	err := json.Unmarshal(float16JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int16")
	}
}

func TestNullUint16_UnmarshalText(t *testing.T) {
	var i Uint16
	err := i.UnmarshalText(uint16JSON)
	maybePanic(err)
	assertUint16(t, i, "UnmarshalText() int16")

	var blank Uint16
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUint16(t, blank, "UnmarshalText() empty int16")

	var null Uint16
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullUint16(t, null, `UnmarshalText() "null"`)
}

func TestNullUint16_JsonMarshal(t *testing.T) {
	i := MakeUint16(math.MaxUint16)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, string(uint16JSON), "non-empty json marshal")

	// invalid values should be encoded as null
	null := Uint16{}
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullUint16_MarshalText(t *testing.T) {
	i := MakeUint16(math.MaxUint16)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, string(uint16JSON), "non-empty text marshal")

	// invalid values should be encoded as null
	var null Uint16
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullUint16_BinaryEncoding(t *testing.T) {
	runner := func(b Uint16, want []byte) func(*testing.T) {
		return func(t *testing.T) {
			data, err := b.MarshalBinary()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": MarshalBinary: %q", data)
			data, err = b.Marshal()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": Marshal: %q", data)

			var decoded Uint16
			assert.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("987654161", runner(MakeUint16(math.MaxUint16), []byte("\b\xff\xff\x03\x10\x01")))
	t.Run("maxUint16", runner(MakeUint16(math.MaxUint16), []byte("\b\xff\xff\x03\x10\x01")))
	t.Run("null", runner(Uint16{}, []byte("")))
}

func TestUint16Pointer(t *testing.T) {
	i := MakeUint16(math.MaxUint16)
	ptr := i.Ptr()
	if *ptr != math.MaxUint16 {
		t.Errorf("bad %s int16: %#v ≠ %d\n", "pointer", ptr, math.MaxUint16)
	}

	null := Uint16{}
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int16: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUint16IsZero(t *testing.T) {
	i := MakeUint16(math.MaxUint16)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := Uint16{}
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeUint16(0)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestUint16SetValid(t *testing.T) {
	var change Uint16
	assertNullUint16(t, change, "SetValid()")
	change.SetValid(math.MaxUint16)
	assertUint16(t, change, "SetValid()")
}

func TestUint16SetPtr(t *testing.T) {
	var change Uint16
	v := uint16(math.MaxUint16)
	change.SetPtr(&v)
	assertUint16(t, change, "SetPtr()")
}

func TestUint16Scan(t *testing.T) {
	var i Uint16
	err := i.Scan(uint16JSON)
	maybePanic(err)
	assertUint16(t, i, "scanned int16")

	var null Uint16
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUint16(t, null, "scanned null")
}

func assertUint16(t *testing.T, i Uint16, from string) {
	if i.Uint16 != math.MaxUint16 {
		t.Errorf("bad %q int16: %d ≠ %d\n", from, i.Uint16, math.MaxUint16)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint16(t *testing.T, i Uint16, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullUint16(t *testing.T) {
	assert.EqualValues(t, 65531, MakeUint16(65531).Uint16)
	assert.True(t, MakeUint16(65531).Valid)
	assert.True(t, MakeUint16(0).Valid)
	v, err := MakeUint16(65531).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, 65531, v)
}

func TestNullUint16_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var nv Uint16
		assert.NoError(t, nv.Scan(nil))
		assert.Exactly(t, Uint16{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv Uint16
		assert.NoError(t, nv.Scan(uint16JSON))
		assert.Exactly(t, MakeUint16(math.MaxUint16), nv)
	})
	t.Run("int16", func(t *testing.T) {
		var nv Uint16
		assert.NoError(t, nv.Scan(int16(math.MaxInt16)))
		assert.Exactly(t, MakeUint16(math.MaxInt16), nv)
	})
	t.Run("int", func(t *testing.T) {
		var nv Uint16
		assert.NoError(t, nv.Scan(int(math.MaxInt16)))
		assert.Exactly(t, MakeUint16(math.MaxInt16), nv)
	})
	t.Run("string unsupported", func(t *testing.T) {
		var nv Uint16
		err := nv.Scan(`1234567`)
		assert.True(t, errors.MatchKind(err, errors.NotSupported), "Error behaviour should be errors.NotSupported")
		assert.Exactly(t, Uint16{}, nv)
	})
	t.Run("parse error negative", func(t *testing.T) {
		var nv Uint16
		err := nv.Scan([]byte(`-1234567`))
		assert.EqualError(t, err, `strconv.ParseUint: parsing "-1234567": invalid syntax`)
		assert.Exactly(t, Uint16{}, nv)
	})
}
