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
	int16JSON     = []byte(fmt.Sprintf("%d", math.MaxInt16))
	int16JSONNeg  = []byte(fmt.Sprintf("%d", -math.MaxInt16))
	nullInt16JSON = []byte(fmt.Sprintf(`{"Int16":%d,"Valid":true}`, math.MaxInt16))
)

var (
	_ fmt.GoStringer             = (*Int16)(nil)
	_ fmt.Stringer               = (*Int16)(nil)
	_ json.Marshaler             = (*Int16)(nil)
	_ json.Unmarshaler           = (*Int16)(nil)
	_ encoding.BinaryMarshaler   = (*Int16)(nil)
	_ encoding.BinaryUnmarshaler = (*Int16)(nil)
	_ encoding.TextMarshaler     = (*Int16)(nil)
	_ encoding.TextUnmarshaler   = (*Int16)(nil)
	_ driver.Valuer              = (*Int16)(nil)
	_ proto.Marshaler            = (*Int16)(nil)
	_ proto.Unmarshaler          = (*Int16)(nil)
	_ proto.Sizer                = (*Int16)(nil)
	_ protoMarshalToer           = (*Int16)(nil)
	_ sql.Scanner                = (*Int16)(nil)
)

func TestMakeNullInt16(t *testing.T) {
	i := MakeInt16(math.MaxInt16)
	assertInt16(t, i, "MakeInt16()")

	zero := MakeInt16(0)
	if !zero.Valid {
		t.Error("MakeInt16(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "null", Int16{}.String())
	assert.Exactly(t, 2, zero.Size())
	assert.Exactly(t, 4, MakeInt16(125).Size())
	assert.Exactly(t, 5, MakeInt16(128).Size())
	assert.Exactly(t, "0", zero.String())
	assert.Exactly(t, string(int16JSON), i.String())
	assert.Exactly(t, 0, Int16{}.Size())
}

func TestInt16_GoString(t *testing.T) {
	tests := []struct {
		i16  Int16
		want string
	}{
		{Int16{}, "null.Int16{}"},
		{MakeInt16(2), "null.MakeInt16(2)"},
	}
	for i, test := range tests {
		if have, want := fmt.Sprintf("%#v", test.i16), test.want; have != want {
			t.Errorf("%d: Have: %v Want: %v", i, have, want)
		}
	}
}

func TestNullInt16_JsonUnmarshal(t *testing.T) {
	var i Int16
	err := json.Unmarshal(int16JSON, &i)
	maybePanic(err)
	assertInt16(t, i, "int16 json")

	var ni Int16
	err = json.Unmarshal(nullInt16JSON, &ni)
	maybePanic(err)
	assertInt16(t, ni, "sql.Int16 json")

	var null Int16
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt16(t, null, "null json")

	var badType Int16
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt16(t, badType, "wrong type json")

	var invalid Int16
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullInt16(t, invalid, "invalid json")
}

func TestNullInt16_JsonUnmarshalNonIntegerNumber(t *testing.T) {
	var i Int16
	err := json.Unmarshal(float16JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int16")
	}
}

func TestNullInt16_JsonUnmarshalInt16Overflow(t *testing.T) {
	int16Overflow := uint16(math.MaxInt16)

	// Max int16 should decode successfully
	var i Int16
	err := json.Unmarshal([]byte(strconv.FormatUint(uint64(int16Overflow), 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int16Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(uint64(int16Overflow), 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int16")
	}
}

func TestNullInt16_UnmarshalText(t *testing.T) {
	var i Int16
	err := i.UnmarshalText(int16JSON)
	maybePanic(err)
	assertInt16(t, i, "UnmarshalText() int16")

	var blank Int16
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt16(t, blank, "UnmarshalText() empty int16")

	var null Int16
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullInt16(t, null, `UnmarshalText() "null"`)
}

func TestNullInt16_JsonMarshal(t *testing.T) {
	i := MakeInt16(math.MaxInt16)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, string(int16JSON), "non-empty json marshal")

	// invalid values should be encoded as null
	null := Int16{}
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullInt16_MarshalText(t *testing.T) {
	i := MakeInt16(math.MaxInt16)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, string(int16JSON), "non-empty text marshal")

	// invalid values should be encoded as null
	null := MakeInt16(0).SetNull()
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullInt16_BinaryEncoding(t *testing.T) {
	runner := func(b Int16, want []byte) func(*testing.T) {
		return func(t *testing.T) {
			data, err := b.MarshalBinary()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": MarshalBinary: %q", data)
			data, err = b.Marshal()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": Marshal: %q", data)

			var decoded Int16
			assert.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("-987654161", runner(MakeInt16(-math.MaxInt16), []byte("\b\x81\x80\xfe\xff\xff\xff\xff\xff\xff\x01\x10\x01")))
	t.Run("987654161", runner(MakeInt16(math.MaxInt16), []byte("\b\xff\xff\x01\x10\x01")))
	t.Run("-maxInt16", runner(MakeInt16(-math.MaxInt16), []byte("\b\x81\x80\xfe\xff\xff\xff\xff\xff\xff\x01\x10\x01")))
	t.Run("maxInt16", runner(MakeInt16(math.MaxInt16), []byte("\b\xff\xff\x01\x10\x01")))
	t.Run("null", runner(Int16{}, []byte("")))
}

func TestInt16Pointer(t *testing.T) {
	i := MakeInt16(math.MaxInt16)
	ptr := i.Ptr()
	if *ptr != math.MaxInt16 {
		t.Errorf("bad %s int16: %#v ≠ %d\n", "pointer", ptr, math.MaxInt16)
	}

	null := Int16{}
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int16: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt16IsZero(t *testing.T) {
	i := MakeInt16(math.MaxInt16)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := MakeInt16(0).SetNull()
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeInt16(0)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestInt16SetValid(t *testing.T) {
	change := MakeInt16(0).SetNull()
	assertNullInt16(t, change, "SetValid()")

	assertInt16(t, change.SetValid(math.MaxInt16), "SetValid()")
}

func TestInt16Scan(t *testing.T) {
	var i Int16
	err := i.Scan(math.MaxInt16)
	maybePanic(err)
	assertInt16(t, i, "scanned int16")

	var null Int16
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt16(t, null, "scanned null")
}

func assertInt16(t *testing.T, i Int16, from string) {
	if i.Int16 != math.MaxInt16 {
		t.Errorf("bad %s int16: %d ≠ %d\n", from, i.Int16, math.MaxInt16)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt16(t *testing.T, i Int16, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullInt16(t *testing.T) {
	assert.EqualValues(t, math.MaxInt16, MakeInt16(math.MaxInt16).Int16)
	assert.True(t, MakeInt16(math.MaxInt16).Valid)
	assert.True(t, MakeInt16(0).Valid)
	v, err := MakeInt16(math.MaxInt16).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, math.MaxInt16, v)
}

func TestNullInt16_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var nv Int16
		assert.NoError(t, nv.Scan(nil))
		assert.Exactly(t, Int16{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv Int16
		assert.NoError(t, nv.Scan(int16JSONNeg))
		assert.Exactly(t, MakeInt16(-math.MaxInt16), nv)
	})
	t.Run("int16", func(t *testing.T) {
		var nv Int16
		assert.NoError(t, nv.Scan(int16(-math.MaxInt16)))
		assert.Exactly(t, MakeInt16(-math.MaxInt16), nv)
	})
	t.Run("int", func(t *testing.T) {
		var nv Int16
		assert.NoError(t, nv.Scan(int(-math.MaxInt16)))
		assert.Exactly(t, MakeInt16(-math.MaxInt16), nv)
	})
	t.Run("string unsupported", func(t *testing.T) {
		var nv Int16
		err := nv.Scan(string(int16JSONNeg))
		assert.True(t, errors.Is(err, errors.NotSupported), "Error behaviour should be errors.NotSupported")
		assert.Exactly(t, Int16{}, nv)
	})
}
