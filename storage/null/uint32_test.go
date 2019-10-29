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
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/gogo/protobuf/proto"
)

var (
	_ fmt.GoStringer             = (*Uint32)(nil)
	_ fmt.Stringer               = (*Uint32)(nil)
	_ json.Marshaler             = (*Uint32)(nil)
	_ json.Unmarshaler           = (*Uint32)(nil)
	_ encoding.BinaryMarshaler   = (*Uint32)(nil)
	_ encoding.BinaryUnmarshaler = (*Uint32)(nil)
	_ encoding.TextMarshaler     = (*Uint32)(nil)
	_ encoding.TextUnmarshaler   = (*Uint32)(nil)
	_ gob.GobEncoder             = (*Uint32)(nil)
	_ gob.GobDecoder             = (*Uint32)(nil)
	_ driver.Valuer              = (*Uint32)(nil)
	_ proto.Marshaler            = (*Uint32)(nil)
	_ proto.Unmarshaler          = (*Uint32)(nil)
	_ proto.Sizer                = (*Uint32)(nil)
	_ protoMarshalToer           = (*Uint32)(nil)
	_ sql.Scanner                = (*Uint32)(nil)
)

var (
	uint32JSON     = []byte(fmt.Sprintf("%d", math.MaxUint32))
	nullUint32JSON = []byte(fmt.Sprintf(`{"Uint32":%d,"Valid":true}`, math.MaxUint32))
	float32JSON    = []byte(`1.23456`)
)

func TestMakeNullUint32(t *testing.T) {
	i := MakeUint32(math.MaxUint32)
	assertUint32(t, i, "MakeUint32()")

	zero := MakeUint32(0)
	if !zero.Valid {
		t.Error("MakeUint32(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "null", Uint32{}.String())
	assert.Exactly(t, 8, zero.Size())
	assert.Exactly(t, 8, MakeUint32(125).Size())
	assert.Exactly(t, 8, MakeUint32(128).Size())
	assert.Exactly(t, "0", zero.String())
	assert.Exactly(t, "4294967295", i.String(), "Want: %q", i.String())
	assert.Exactly(t, 0, Uint32{}.Size())
}

func TestUint32_GoString(t *testing.T) {
	tests := []struct {
		i32  Uint32
		want string
	}{
		{Uint32{}, "null.Uint32{}"},
		{MakeUint32(2), "null.MakeUint32(2)"},
	}
	for i, test := range tests {
		if have, want := fmt.Sprintf("%#v", test.i32), test.want; have != want {
			t.Errorf("%d: Have: %v Want: %v", i, have, want)
		}
	}
}

func TestMakeUint32FromByte(t *testing.T) {
	ui, err := MakeUint32FromByte([]byte(`987654321`))
	assert.NoError(t, err)
	assert.Exactly(t, Uint32{Uint32: 987654321, Valid: true}, ui)
}

func TestNullUint32_JsonUnmarshal(t *testing.T) {
	var err error
	var i Uint32
	err = json.Unmarshal(uint32JSON, &i)
	assert.NoError(t, err, "%+v", err)
	assertUint32(t, i, "int32 json")

	var ni Uint32
	err = json.Unmarshal(nullUint32JSON, &ni)
	assert.NoError(t, err, "%+v", err)
	assertUint32(t, ni, "null.Uint32 json")

	var null Uint32
	err = json.Unmarshal(nullJSON, &null)
	assert.NoError(t, err, "%+v", err)
	assertNullUint32(t, null, "null json")

	var badType Uint32
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint32(t, badType, "wrong type json")

	var invalid Uint32
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint32(t, invalid, "invalid json")
}

func TestNullUint32_JsonUnmarshalNonIntegerNumber(t *testing.T) {
	var i Uint32
	err := json.Unmarshal(float32JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int32")
	}
}

func TestNullUint32_UnmarshalText(t *testing.T) {
	var i Uint32
	err := i.UnmarshalText(uint32JSON)
	maybePanic(err)
	assertUint32(t, i, "UnmarshalText() int32")

	var blank Uint32
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUint32(t, blank, "UnmarshalText() empty int32")

	var null Uint32
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullUint32(t, null, `UnmarshalText() "null"`)
}

func TestNullUint32_JsonMarshal(t *testing.T) {
	i := MakeUint32(math.MaxUint32)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, string(uint32JSON), "non-empty json marshal")

	// invalid values should be encoded as null
	null := Uint32{}
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullUint32_MarshalText(t *testing.T) {
	i := MakeUint32(math.MaxUint32)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, string(uint32JSON), "non-empty text marshal")

	// invalid values should be encoded as null
	null := MakeUint32(0).SetNull()
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullUint32_BinaryEncoding(t *testing.T) {
	runner := func(b Uint32, want []byte) func(*testing.T) {
		return func(t *testing.T) {
			data, err := b.GobEncode()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": GobEncode: %q", data)
			data, err = b.MarshalBinary()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": MarshalBinary: %q", data)
			data, err = b.Marshal()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": Marshal: %q", data)

			var decoded Uint32
			assert.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("987654321", runner(MakeUint32(987654321), []byte{0xb1, 0x68, 0xde, 0x3a, 0x0, 0x0, 0x0, 0x0}))
	t.Run("maxUint32", runner(MakeUint32(math.MaxUint32), []byte("\xff\xff\xff\xff\x00\x00\x00\x00")))
	t.Run("null", runner(Uint32{}, nil))
}

func TestUint32Pointer(t *testing.T) {
	i := MakeUint32(math.MaxUint32)
	ptr := i.Ptr()
	if *ptr != math.MaxUint32 {
		t.Errorf("bad %s int32: %#v ≠ %d\n", "pointer", ptr, math.MaxUint32)
	}

	null := Uint32{}
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int32: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUint32IsZero(t *testing.T) {
	i := MakeUint32(math.MaxUint32)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := Uint32{}
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeUint32(0)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestUint32SetValid(t *testing.T) {
	change := MakeUint32(0).SetNull()
	assertNullUint32(t, change, "SetValid()")

	assertUint32(t, change.SetValid(math.MaxUint32), "SetValid()")
}

func TestUint32Scan(t *testing.T) {
	var i Uint32
	err := i.Scan(uint32JSON)
	maybePanic(err)
	assertUint32(t, i, "scanned int32")

	var null Uint32
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUint32(t, null, "scanned null")
}

func assertUint32(t *testing.T, i Uint32, from string) {
	if i.Uint32 != math.MaxUint32 {
		t.Errorf("bad %q int32: %d ≠ %d\n", from, i.Uint32, math.MaxUint32)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint32(t *testing.T, i Uint32, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullUint32(t *testing.T) {
	assert.EqualValues(t, 1257894000, MakeUint32(1257894000).Uint32)
	assert.True(t, MakeUint32(1257894000).Valid)
	assert.True(t, MakeUint32(0).Valid)
	v, err := MakeUint32(1257894000).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, 1257894000, v)
}

func TestNullUint32_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var nv Uint32
		assert.NoError(t, nv.Scan(nil))
		assert.Exactly(t, Uint32{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv Uint32
		assert.NoError(t, nv.Scan(uint32JSON))
		assert.Exactly(t, MakeUint32(math.MaxUint32), nv)
	})
	t.Run("int32", func(t *testing.T) {
		var nv Uint32
		assert.NoError(t, nv.Scan(int32(math.MaxInt32)))
		assert.Exactly(t, MakeUint32(math.MaxInt32), nv)
	})
	t.Run("int", func(t *testing.T) {
		var nv Uint32
		assert.NoError(t, nv.Scan(int(math.MaxInt32)))
		assert.Exactly(t, MakeUint32(math.MaxInt32), nv)
	})
	t.Run("string unsupported", func(t *testing.T) {
		var nv Uint32
		err := nv.Scan(`1234567`)
		assert.True(t, errors.Is(err, errors.NotSupported), "Error behaviour should be errors.NotSupported")
		assert.Exactly(t, Uint32{}, nv)
	})
	t.Run("parse error negative", func(t *testing.T) {
		var nv Uint32
		err := nv.Scan([]byte(`-1234567`))
		assert.EqualError(t, err, `strconv.ParseUint: parsing "-1234567": invalid syntax`)
		assert.Exactly(t, Uint32{}, nv)
	})
}
