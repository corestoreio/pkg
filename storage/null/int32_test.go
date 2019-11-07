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
	int32JSON     = []byte(fmt.Sprintf("%d", math.MaxInt32))
	nullInt32JSON = []byte(fmt.Sprintf(`{"Int32":%d,"Valid":true}`, math.MaxInt32))
)

var (
	_ fmt.GoStringer             = (*Int32)(nil)
	_ fmt.Stringer               = (*Int32)(nil)
	_ json.Marshaler             = (*Int32)(nil)
	_ json.Unmarshaler           = (*Int32)(nil)
	_ encoding.BinaryMarshaler   = (*Int32)(nil)
	_ encoding.BinaryUnmarshaler = (*Int32)(nil)
	_ encoding.TextMarshaler     = (*Int32)(nil)
	_ encoding.TextUnmarshaler   = (*Int32)(nil)
	_ driver.Valuer              = (*Int32)(nil)
	_ proto.Marshaler            = (*Int32)(nil)
	_ proto.Unmarshaler          = (*Int32)(nil)
	_ proto.Sizer                = (*Int32)(nil)
	_ protoMarshalToer           = (*Int32)(nil)
	_ sql.Scanner                = (*Int32)(nil)
)

func TestMakeNullInt32(t *testing.T) {
	i := MakeInt32(math.MaxInt32)
	assertInt32(t, i, "MakeInt32()")

	zero := MakeInt32(0)
	if !zero.Valid {
		t.Error("MakeInt32(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "null", Int32{}.String())
	assert.Exactly(t, 2, zero.Size())
	assert.Exactly(t, 4, MakeInt32(125).Size())
	assert.Exactly(t, 5, MakeInt32(128).Size())
	assert.Exactly(t, "0", zero.String())
	assert.Exactly(t, string(int32JSON), i.String())
	assert.Exactly(t, 0, Int32{}.Size())
}

func TestInt32_GoString(t *testing.T) {
	tests := []struct {
		i32  Int32
		want string
	}{
		{Int32{}, "null.Int32{}"},
		{MakeInt32(2), "null.MakeInt32(2)"},
	}
	for i, test := range tests {
		if have, want := fmt.Sprintf("%#v", test.i32), test.want; have != want {
			t.Errorf("%d: Have: %v Want: %v", i, have, want)
		}
	}
}

func TestNullInt32_JsonUnmarshal(t *testing.T) {
	var i Int32
	err := json.Unmarshal(int32JSON, &i)
	maybePanic(err)
	assertInt32(t, i, "int32 json")

	var ni Int32
	err = json.Unmarshal(nullInt32JSON, &ni)
	maybePanic(err)
	assertInt32(t, ni, "sql.Int32 json")

	var null Int32
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt32(t, null, "null json")

	var badType Int32
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt32(t, badType, "wrong type json")

	var invalid Int32
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullInt32(t, invalid, "invalid json")
}

func TestNullInt32_JsonUnmarshalNonIntegerNumber(t *testing.T) {
	var i Int32
	err := json.Unmarshal(float32JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int32")
	}
}

func TestNullInt32_JsonUnmarshalInt32Overflow(t *testing.T) {
	int32Overflow := uint32(math.MaxInt32)

	// Max int32 should decode successfully
	var i Int32
	err := json.Unmarshal([]byte(strconv.FormatUint(uint64(int32Overflow), 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int32Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(uint64(int32Overflow), 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int32")
	}
}

func TestNullInt32_UnmarshalText(t *testing.T) {
	var i Int32
	err := i.UnmarshalText(int32JSON)
	maybePanic(err)
	assertInt32(t, i, "UnmarshalText() int32")

	var blank Int32
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt32(t, blank, "UnmarshalText() empty int32")

	var null Int32
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullInt32(t, null, `UnmarshalText() "null"`)
}

func TestNullInt32_JsonMarshal(t *testing.T) {
	i := MakeInt32(math.MaxInt32)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, string(int32JSON), "non-empty json marshal")

	// invalid values should be encoded as null
	null := Int32{}
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullInt32_MarshalText(t *testing.T) {
	i := MakeInt32(math.MaxInt32)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, string(int32JSON), "non-empty text marshal")

	// invalid values should be encoded as null
	null := MakeInt32(0).SetNull()
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullInt32_BinaryEncoding(t *testing.T) {
	runner := func(b Int32, want []byte) func(*testing.T) {
		return func(t *testing.T) {
			data, err := b.MarshalBinary()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": MarshalBinary: %q", data)
			data, err = b.Marshal()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": Marshal: %q", data)

			var decoded Int32
			assert.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("-987654321", runner(MakeInt32(-987654321), []byte("\bϮ\x86\xa9\xfc\xff\xff\xff\xff\x01\x10\x01")))
	t.Run("987654321", runner(MakeInt32(987654321), []byte("\b\xb1\xd1\xf9\xd6\x03\x10\x01")))
	t.Run("-maxInt32", runner(MakeInt32(-math.MaxInt32), []byte("\b\x81\x80\x80\x80\xf8\xff\xff\xff\xff\x01\x10\x01")))
	t.Run("maxInt32", runner(MakeInt32(math.MaxInt32), []byte("\b\xff\xff\xff\xff\a\x10\x01")))
	t.Run("null", runner(Int32{}, []byte("")))
}

func TestInt32Pointer(t *testing.T) {
	i := MakeInt32(math.MaxInt32)
	ptr := i.Ptr()
	if *ptr != math.MaxInt32 {
		t.Errorf("bad %s int32: %#v ≠ %d\n", "pointer", ptr, math.MaxInt32)
	}

	null := Int32{}
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int32: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt32IsZero(t *testing.T) {
	i := MakeInt32(math.MaxInt32)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := MakeInt32(0).SetNull()
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeInt32(0)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestInt32SetValid(t *testing.T) {
	change := MakeInt32(0).SetNull()
	assertNullInt32(t, change, "SetValid()")

	assertInt32(t, change.SetValid(math.MaxInt32), "SetValid()")
}

func TestInt32Scan(t *testing.T) {
	var i Int32
	err := i.Scan(math.MaxInt32)
	maybePanic(err)
	assertInt32(t, i, "scanned int32")

	var null Int32
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt32(t, null, "scanned null")
}

func assertInt32(t *testing.T, i Int32, from string) {
	if i.Int32 != math.MaxInt32 {
		t.Errorf("bad %s int32: %d ≠ %d\n", from, i.Int32, math.MaxInt32)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt32(t *testing.T, i Int32, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullInt32(t *testing.T) {
	assert.EqualValues(t, 1257894000, MakeInt32(1257894000).Int32)
	assert.True(t, MakeInt32(1257894000).Valid)
	assert.True(t, MakeInt32(0).Valid)
	v, err := MakeInt32(1257894000).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, 1257894000, v)
}

func TestNullInt32_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var nv Int32
		assert.NoError(t, nv.Scan(nil))
		assert.Exactly(t, Int32{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv Int32
		assert.NoError(t, nv.Scan([]byte(`-1234567`)))
		assert.Exactly(t, MakeInt32(-1234567), nv)
	})
	t.Run("int32", func(t *testing.T) {
		var nv Int32
		assert.NoError(t, nv.Scan(int32(-1234568)))
		assert.Exactly(t, MakeInt32(-1234568), nv)
	})
	t.Run("int", func(t *testing.T) {
		var nv Int32
		assert.NoError(t, nv.Scan(int(-1234569)))
		assert.Exactly(t, MakeInt32(-1234569), nv)
	})
	t.Run("string unsupported", func(t *testing.T) {
		var nv Int32
		err := nv.Scan(`-1234567`)
		assert.True(t, errors.Is(err, errors.NotSupported), "Error behaviour should be errors.NotSupported")
		assert.Exactly(t, Int32{}, nv)
	})
}
