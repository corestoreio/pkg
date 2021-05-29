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
	_ fmt.GoStringer             = (*Uint8)(nil)
	_ fmt.Stringer               = (*Uint8)(nil)
	_ json.Marshaler             = (*Uint8)(nil)
	_ json.Unmarshaler           = (*Uint8)(nil)
	_ encoding.BinaryMarshaler   = (*Uint8)(nil)
	_ encoding.BinaryUnmarshaler = (*Uint8)(nil)
	_ encoding.TextMarshaler     = (*Uint8)(nil)
	_ encoding.TextUnmarshaler   = (*Uint8)(nil)
	_ driver.Valuer              = (*Uint8)(nil)
	_ sql.Scanner                = (*Uint8)(nil)
)

var (
	uint8JSON     = []byte(fmt.Sprintf("%d", math.MaxUint8))
	nullUint8JSON = []byte(fmt.Sprintf(`{"Uint8":%d,"Valid":true}`, math.MaxUint8))
	float8JSON    = []byte(`1.23456`)
)

func TestMakeNullUint8(t *testing.T) {
	i := MakeUint8(math.MaxUint8)
	assertUint8(t, i, "MakeUint8()")

	zero := MakeUint8(0)
	if !zero.Valid {
		t.Error("MakeUint8(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "null", Uint8{}.String())
	assert.Exactly(t, 2, zero.Size())
	assert.Exactly(t, 4, MakeUint8(125).Size())
	assert.Exactly(t, 5, MakeUint8(128).Size())
	assert.Exactly(t, "0", zero.String())
	assert.Exactly(t, "255", i.String(), "Want: %q", i.String())
	assert.Exactly(t, 0, Uint8{}.Size())
}

func TestUint8_GoString(t *testing.T) {
	tests := []struct {
		i8   Uint8
		want string
	}{
		{Uint8{}, "null.Uint8{}"},
		{MakeUint8(2), "null.MakeUint8(2)"},
	}
	for i, test := range tests {
		if have, want := fmt.Sprintf("%#v", test.i8), test.want; have != want {
			t.Errorf("%d: Have: %v Want: %v", i, have, want)
		}
	}
}

func TestMakeUint8FromByte(t *testing.T) {
	ui, err := MakeUint8FromByte([]byte(`255`))
	assert.NoError(t, err)
	assert.Exactly(t, Uint8{Uint8: math.MaxUint8, Valid: true}, ui)
}

func TestNullUint8_JsonUnmarshal(t *testing.T) {
	var err error
	var i Uint8
	err = json.Unmarshal(uint8JSON, &i)
	assert.NoError(t, err, "%+v", err)
	assertUint8(t, i, "int8 json")

	var ni Uint8
	err = json.Unmarshal(nullUint8JSON, &ni)
	assert.NoError(t, err, "%+v", err)
	assertUint8(t, ni, "null.Uint8 json")

	var null Uint8
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err, "%+v", err)
	assertNullUint8(t, null, "null json")

	var badType Uint8
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint8(t, badType, "wrong type json")

	var invalid Uint8
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint8(t, invalid, "invalid json")
}

func TestNullUint8_JsonUnmarshalNonIntegerNumber(t *testing.T) {
	var i Uint8
	err := json.Unmarshal(float8JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int8")
	}
}

func TestNullUint8_UnmarshalText(t *testing.T) {
	var i Uint8
	err := i.UnmarshalText(uint8JSON)
	maybePanic(err)
	assertUint8(t, i, "UnmarshalText() int8")

	var blank Uint8
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUint8(t, blank, "UnmarshalText() empty int8")

	var null Uint8
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullUint8(t, null, `UnmarshalText() "null"`)
}

func TestNullUint8_JsonMarshal(t *testing.T) {
	i := MakeUint8(math.MaxUint8)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, string(uint8JSON), "non-empty json marshal")

	// invalid values should be encoded as null
	null := Uint8{}
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullUint8_MarshalText(t *testing.T) {
	i := MakeUint8(math.MaxUint8)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, string(uint8JSON), "non-empty text marshal")

	// invalid values should be encoded as null
	var null Uint8
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullUint8_BinaryEncoding(t *testing.T) {
	runner := func(b Uint8, want []byte) func(*testing.T) {
		return func(t *testing.T) {
			data, err := b.MarshalBinary()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": MarshalBinary: %q", data)
			data, err = b.Marshal()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": Marshal: %q", data)

			var decoded Uint8
			assert.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("98765481", runner(MakeUint8(math.MaxUint8), []byte("\b\xff\x01\x10\x01")))
	t.Run("maxUint8", runner(MakeUint8(math.MaxUint8), []byte("\b\xff\x01\x10\x01")))
	t.Run("null", runner(Uint8{}, []byte("")))
}

func TestUint8Pointer(t *testing.T) {
	i := MakeUint8(math.MaxUint8)
	ptr := i.Ptr()
	if *ptr != math.MaxUint8 {
		t.Errorf("bad %s int8: %#v ≠ %d\n", "pointer", ptr, math.MaxUint8)
	}

	null := Uint8{}
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int8: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUint8IsZero(t *testing.T) {
	i := MakeUint8(math.MaxUint8)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := Uint8{}
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeUint8(0)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestUint8SetValid(t *testing.T) {
	var change Uint8
	assertNullUint8(t, change, "SetValid()")
	change.SetValid(math.MaxUint8)
	assertUint8(t, change, "SetValid()")
}

func TestUint8SetPtr(t *testing.T) {
	var change Uint8
	v := uint8(math.MaxUint8)
	change.SetPtr(&v)
	assertUint8(t, change, "SetPtr()")
}

func TestUint8Scan(t *testing.T) {
	var i Uint8
	err := i.Scan(uint8JSON)
	maybePanic(err)
	assertUint8(t, i, "scanned int8")

	var null Uint8
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUint8(t, null, "scanned null")
}

func assertUint8(t *testing.T, i Uint8, from string) {
	if i.Uint8 != math.MaxUint8 {
		t.Errorf("bad %q int8: %d ≠ %d\n", from, i.Uint8, math.MaxUint8)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint8(t *testing.T, i Uint8, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullUint8(t *testing.T) {
	assert.EqualValues(t, math.MaxUint8, MakeUint8(math.MaxUint8).Uint8)
	assert.True(t, MakeUint8(math.MaxUint8).Valid)
	assert.True(t, MakeUint8(0).Valid)
	v, err := MakeUint8(math.MaxUint8).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, uint8JSON, v)
}

func TestNullUint8_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var nv Uint8
		assert.NoError(t, nv.Scan(nil))
		assert.Exactly(t, Uint8{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv Uint8
		assert.NoError(t, nv.Scan(uint8JSON))
		assert.Exactly(t, MakeUint8(math.MaxUint8), nv)
	})
	t.Run("int8", func(t *testing.T) {
		var nv Uint8
		assert.NoError(t, nv.Scan(int8(math.MaxInt8)))
		assert.Exactly(t, MakeUint8(math.MaxInt8), nv)
	})
	t.Run("int", func(t *testing.T) {
		var nv Uint8
		assert.NoError(t, nv.Scan(int(math.MaxInt8)))
		assert.Exactly(t, MakeUint8(math.MaxInt8), nv)
	})
	t.Run("string unsupported", func(t *testing.T) {
		var nv Uint8
		err := nv.Scan(`1234567`)
		assert.True(t, errors.MatchKind(err, errors.NotSupported), "Error behaviour should be errors.NotSupported")
		assert.Exactly(t, Uint8{}, nv)
	})
	t.Run("parse error negative", func(t *testing.T) {
		var nv Uint8
		err := nv.Scan([]byte(`-1234567`))
		assert.EqualError(t, err, `strconv.ParseUint: parsing "-1234567": invalid syntax`)
		assert.Exactly(t, Uint8{}, nv)
	})
}
