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
	"strconv"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	int64JSON     = []byte(`9223372036854775806`)
	nullInt64JSON = []byte(`{"Int64":9223372036854775806,"Valid":true}`)
)

var (
	_ fmt.GoStringer             = (*Int64)(nil)
	_ fmt.Stringer               = (*Int64)(nil)
	_ json.Marshaler             = (*Int64)(nil)
	_ json.Unmarshaler           = (*Int64)(nil)
	_ encoding.BinaryMarshaler   = (*Int64)(nil)
	_ encoding.BinaryUnmarshaler = (*Int64)(nil)
	_ encoding.TextMarshaler     = (*Int64)(nil)
	_ encoding.TextUnmarshaler   = (*Int64)(nil)
	_ gob.GobEncoder             = (*Int64)(nil)
	_ gob.GobDecoder             = (*Int64)(nil)
	_ driver.Valuer              = (*Int64)(nil)
	_ proto.Marshaler            = (*Int64)(nil)
	_ proto.Unmarshaler          = (*Int64)(nil)
	_ proto.Sizer                = (*Int64)(nil)
	_ protoMarshalToer           = (*Int64)(nil)
	_ sql.Scanner                = (*Int64)(nil)
)

func TestMakeNullInt64(t *testing.T) {

	i := MakeInt64(9223372036854775806)
	assertInt64(t, i, "MakeInt64()")

	zero := MakeInt64(0)
	if !zero.Valid {
		t.Error("MakeInt64(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "null", Int64{}.String())
	assert.Exactly(t, 8, zero.Size())
	assert.Exactly(t, 8, MakeInt64(125).Size())
	assert.Exactly(t, 8, MakeInt64(128).Size())
	assert.Exactly(t, "0", zero.String())
	assert.Exactly(t, "9223372036854775806", i.String())
	assert.Exactly(t, 0, Int64{}.Size())
}

func TestInt64_GoString(t *testing.T) {

	tests := []struct {
		i64  Int64
		want string
	}{
		{Int64{}, "null.Int64{}"},
		{MakeInt64(2), "null.MakeInt64(2)"},
	}
	for i, test := range tests {
		if have, want := fmt.Sprintf("%#v", test.i64), test.want; have != want {
			t.Errorf("%d: Have: %v Want: %v", i, have, want)
		}
	}
}

func TestNullInt64_JsonUnmarshal(t *testing.T) {

	var i Int64
	err := json.Unmarshal(int64JSON, &i)
	maybePanic(err)
	assertInt64(t, i, "int64 json")

	var ni Int64
	err = json.Unmarshal(nullInt64JSON, &ni)
	maybePanic(err)
	assertInt64(t, ni, "sql.Int64 json")

	var null Int64
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt64(t, null, "null json")

	var badType Int64
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt64(t, badType, "wrong type json")

	var invalid Int64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullInt64(t, invalid, "invalid json")
}

func TestNullInt64_JsonUnmarshalNonIntegerNumber(t *testing.T) {

	var i Int64
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int64")
	}
}

func TestNullInt64_JsonUnmarshalInt64Overflow(t *testing.T) {

	int64Overflow := uint64(math.MaxInt64)

	// Max int64 should decode successfully
	var i Int64
	err := json.Unmarshal([]byte(strconv.FormatUint(uint64(int64Overflow), 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(uint64(int64Overflow), 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int64")
	}
}

func TestNullInt64_UnmarshalText(t *testing.T) {

	var i Int64
	err := i.UnmarshalText([]byte("9223372036854775806"))
	maybePanic(err)
	assertInt64(t, i, "UnmarshalText() int64")

	var blank Int64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt64(t, blank, "UnmarshalText() empty int64")

	var null Int64
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullInt64(t, null, `UnmarshalText() "null"`)
}

func TestNullInt64_JsonMarshal(t *testing.T) {

	i := MakeInt64(9223372036854775806)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "9223372036854775806", "non-empty json marshal")

	// invalid values should be encoded as null
	null := Int64{}
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullInt64_MarshalText(t *testing.T) {

	i := MakeInt64(9223372036854775806)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "9223372036854775806", "non-empty text marshal")

	// invalid values should be encoded as null
	null := MakeInt64(0).SetNull()
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullInt64_BinaryEncoding(t *testing.T) {

	runner := func(b Int64, want []byte) func(*testing.T) {
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

			var decoded Int64
			require.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("-987654321", runner(MakeInt64(-987654321), []byte{0x4f, 0x97, 0x21, 0xc5, 0xff, 0xff, 0xff, 0xff}))
	t.Run("987654321", runner(MakeInt64(987654321), []byte{0xb1, 0x68, 0xde, 0x3a, 0x0, 0x0, 0x0, 0x0}))
	t.Run("-maxInt64", runner(MakeInt64(-math.MaxInt64), []byte{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x80}))
	t.Run("maxInt64", runner(MakeInt64(math.MaxInt64), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}))
	t.Run("null", runner(Int64{}, nil))
}

func TestInt64Pointer(t *testing.T) {

	i := MakeInt64(9223372036854775806)
	ptr := i.Ptr()
	if *ptr != 9223372036854775806 {
		t.Errorf("bad %s int64: %#v ≠ %d\n", "pointer", ptr, 9223372036854775806)
	}

	null := Int64{}
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt64IsZero(t *testing.T) {

	i := MakeInt64(9223372036854775806)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := MakeInt64(0).SetNull()
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeInt64(0)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestInt64SetValid(t *testing.T) {

	change := MakeInt64(0).SetNull()
	assertNullInt64(t, change, "SetValid()")

	assertInt64(t, change.SetValid(9223372036854775806), "SetValid()")
}

func TestInt64Scan(t *testing.T) {

	var i Int64
	err := i.Scan(9223372036854775806)
	maybePanic(err)
	assertInt64(t, i, "scanned int64")

	var null Int64
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt64(t, null, "scanned null")
}

func assertInt64(t *testing.T, i Int64, from string) {
	if i.Int64 != 9223372036854775806 {
		t.Errorf("bad %s int64: %d ≠ %d\n", from, i.Int64, 9223372036854775806)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt64(t *testing.T, i Int64, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullInt64(t *testing.T) {

	assert.EqualValues(t, 1257894000, MakeInt64(1257894000).Int64)
	assert.True(t, MakeInt64(1257894000).Valid)
	assert.True(t, MakeInt64(0).Valid)
	v, err := MakeInt64(1257894000).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, 1257894000, v)
}

func TestNullInt64_Scan(t *testing.T) {

	t.Run("nil", func(t *testing.T) {
		var nv Int64
		require.NoError(t, nv.Scan(nil))
		assert.Exactly(t, Int64{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv Int64
		require.NoError(t, nv.Scan([]byte(`-1234567`)))
		assert.Exactly(t, MakeInt64(-1234567), nv)
	})
	t.Run("int64", func(t *testing.T) {
		var nv Int64
		require.NoError(t, nv.Scan(int64(-1234568)))
		assert.Exactly(t, MakeInt64(-1234568), nv)
	})
	t.Run("int", func(t *testing.T) {
		var nv Int64
		require.NoError(t, nv.Scan(int(-1234569)))
		assert.Exactly(t, MakeInt64(-1234569), nv)
	})
	t.Run("string unsupported", func(t *testing.T) {
		var nv Int64
		err := nv.Scan(`-1234567`)
		assert.True(t, errors.Is(err, errors.NotSupported), "Error behaviour should be errors.NotSupported")
		assert.Exactly(t, Int64{}, nv)
	})
}
