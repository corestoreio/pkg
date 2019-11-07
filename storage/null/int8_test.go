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
	"strconv"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/gogo/protobuf/proto"
)

var (
	int8JSON     = []byte(fmt.Sprintf("%d", math.MaxInt8))
	int8JSONNeg  = []byte(fmt.Sprintf("%d", -math.MaxInt8))
	nullInt8JSON = []byte(fmt.Sprintf(`{"Int8":%d,"Valid":true}`, math.MaxInt8))
)

var (
	_ fmt.GoStringer             = (*Int8)(nil)
	_ fmt.Stringer               = (*Int8)(nil)
	_ json.Marshaler             = (*Int8)(nil)
	_ json.Unmarshaler           = (*Int8)(nil)
	_ encoding.BinaryMarshaler   = (*Int8)(nil)
	_ encoding.BinaryUnmarshaler = (*Int8)(nil)
	_ encoding.TextMarshaler     = (*Int8)(nil)
	_ encoding.TextUnmarshaler   = (*Int8)(nil)
	_ driver.Valuer              = (*Int8)(nil)
	_ proto.Marshaler            = (*Int8)(nil)
	_ proto.Unmarshaler          = (*Int8)(nil)
	_ proto.Sizer                = (*Int8)(nil)
	_ protoMarshalToer           = (*Int8)(nil)
	_ sql.Scanner                = (*Int8)(nil)
)

func TestMakeNullInt8(t *testing.T) {
	i := MakeInt8(math.MaxInt8)
	assertInt8(t, i, "MakeInt8()")

	zero := MakeInt8(0)
	if !zero.Valid {
		t.Error("MakeInt8(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "null", Int8{}.String())
	assert.Exactly(t, 2, zero.Size())
	assert.Exactly(t, 4, MakeInt8(125).Size())
	assert.Exactly(t, 4, MakeInt8(math.MaxInt8).Size())
	assert.Exactly(t, "0", zero.String())
	assert.Exactly(t, string(int8JSON), i.String())
	assert.Exactly(t, 0, Int8{}.Size())
}

func TestInt8_GoString(t *testing.T) {
	tests := []struct {
		i8   Int8
		want string
	}{
		{Int8{}, "null.Int8{}"},
		{MakeInt8(2), "null.MakeInt8(2)"},
	}
	for i, test := range tests {
		if have, want := fmt.Sprintf("%#v", test.i8), test.want; have != want {
			t.Errorf("%d: Have: %v Want: %v", i, have, want)
		}
	}
}

func TestNullInt8_JsonUnmarshal(t *testing.T) {
	var i Int8
	err := json.Unmarshal(int8JSON, &i)
	maybePanic(err)
	assertInt8(t, i, "int8 json")

	var ni Int8
	err = json.Unmarshal(nullInt8JSON, &ni)
	maybePanic(err)
	assertInt8(t, ni, "sql.Int8 json")

	var null Int8
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt8(t, null, "null json")

	var badType Int8
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt8(t, badType, "wrong type json")

	var invalid Int8
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullInt8(t, invalid, "invalid json")
}

func TestNullInt8_JsonUnmarshalNonIntegerNumber(t *testing.T) {
	var i Int8
	err := json.Unmarshal(float8JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int8")
	}
}

func TestNullInt8_JsonUnmarshalInt8Overflow(t *testing.T) {
	int8Overflow := uint8(math.MaxInt8)

	// Max int8 should decode successfully
	var i Int8
	err := json.Unmarshal([]byte(strconv.FormatUint(uint64(int8Overflow), 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int8Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(uint64(int8Overflow), 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int8")
	}
}

func TestNullInt8_UnmarshalText(t *testing.T) {
	var i Int8
	err := i.UnmarshalText(int8JSON)
	maybePanic(err)
	assertInt8(t, i, "UnmarshalText() int8")

	var blank Int8
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt8(t, blank, "UnmarshalText() empty int8")

	var null Int8
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullInt8(t, null, `UnmarshalText() "null"`)
}

func TestNullInt8_JsonMarshal(t *testing.T) {
	i := MakeInt8(math.MaxInt8)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, string(int8JSON), "non-empty json marshal")

	// invalid values should be encoded as null
	null := Int8{}
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullInt8_MarshalText(t *testing.T) {
	i := MakeInt8(math.MaxInt8)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, string(int8JSON), "non-empty text marshal")

	// invalid values should be encoded as null
	null := MakeInt8(0).SetNull()
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullInt8_BinaryEncoding(t *testing.T) {
	runner := func(b Int8, want []byte) func(*testing.T) {
		return func(t *testing.T) {
			data, err := b.MarshalBinary()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": MarshalBinary: %q", data)
			data, err = b.Marshal()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": Marshal: %q", data)

			var decoded Int8
			assert.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("-98765481", runner(MakeInt8(-math.MaxInt8), []byte("\b\x81\xff\xff\xff\xff\xff\xff\xff\xff\x01\x10\x01")))
	t.Run("98765481", runner(MakeInt8(math.MaxInt8), []byte("\b\u007f\x10\x01")))
	t.Run("-maxInt8", runner(MakeInt8(-math.MaxInt8), []byte("\b\x81\xff\xff\xff\xff\xff\xff\xff\xff\x01\x10\x01")))
	t.Run("maxInt8", runner(MakeInt8(math.MaxInt8), []byte("\b\u007f\x10\x01")))
	t.Run("null", runner(Int8{}, []byte("")))
}

func TestInt8Pointer(t *testing.T) {
	i := MakeInt8(math.MaxInt8)
	ptr := i.Ptr()
	if *ptr != math.MaxInt8 {
		t.Errorf("bad %s int8: %#v ≠ %d\n", "pointer", ptr, math.MaxInt8)
	}

	null := Int8{}
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int8: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt8IsZero(t *testing.T) {
	i := MakeInt8(math.MaxInt8)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := MakeInt8(0).SetNull()
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeInt8(0)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestInt8SetValid(t *testing.T) {
	change := MakeInt8(0).SetNull()
	assertNullInt8(t, change, "SetValid()")

	assertInt8(t, change.SetValid(math.MaxInt8), "SetValid()")
}

func TestInt8Scan(t *testing.T) {
	var i Int8
	err := i.Scan(math.MaxInt8)
	maybePanic(err)
	assertInt8(t, i, "scanned int8")

	var null Int8
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt8(t, null, "scanned null")
}

func assertInt8(t *testing.T, i Int8, from string) {
	if i.Int8 != math.MaxInt8 {
		t.Errorf("bad %s int8: %d ≠ %d\n", from, i.Int8, math.MaxInt8)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt8(t *testing.T, i Int8, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullInt8(t *testing.T) {
	assert.EqualValues(t, math.MaxInt8, MakeInt8(math.MaxInt8).Int8)
	assert.True(t, MakeInt8(math.MaxInt8).Valid)
	assert.True(t, MakeInt8(0).Valid)
	v, err := MakeInt8(math.MaxInt8).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, math.MaxInt8, v)
}

func TestNullInt8_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var nv Int8
		assert.NoError(t, nv.Scan(nil))
		assert.Exactly(t, Int8{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv Int8
		assert.NoError(t, nv.Scan(int8JSONNeg))
		assert.Exactly(t, MakeInt8(-math.MaxInt8), nv)
	})
	t.Run("int8", func(t *testing.T) {
		var nv Int8
		assert.NoError(t, nv.Scan(int8(-math.MaxInt8)))
		assert.Exactly(t, MakeInt8(-math.MaxInt8), nv)
	})
	t.Run("int", func(t *testing.T) {
		var nv Int8
		assert.NoError(t, nv.Scan(int(-math.MaxInt8)))
		assert.Exactly(t, MakeInt8(-math.MaxInt8), nv)
	})
	t.Run("string unsupported", func(t *testing.T) {
		var nv Int8
		err := nv.Scan(string(int8JSONNeg))
		assert.True(t, errors.Is(err, errors.NotSupported), "Error behaviour should be errors.NotSupported")
		assert.Exactly(t, Int8{}, nv)
	})
}
