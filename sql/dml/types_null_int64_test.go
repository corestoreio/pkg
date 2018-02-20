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

package dml

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
	nullInt64JSON = []byte(`{"NullInt64":9223372036854775806,"Valid":true}`)
)

var (
	_ fmt.GoStringer             = (*NullInt64)(nil)
	_ fmt.Stringer               = (*NullInt64)(nil)
	_ json.Marshaler             = (*NullInt64)(nil)
	_ json.Unmarshaler           = (*NullInt64)(nil)
	_ encoding.BinaryMarshaler   = (*NullInt64)(nil)
	_ encoding.BinaryUnmarshaler = (*NullInt64)(nil)
	_ encoding.TextMarshaler     = (*NullInt64)(nil)
	_ encoding.TextUnmarshaler   = (*NullInt64)(nil)
	_ gob.GobEncoder             = (*NullInt64)(nil)
	_ gob.GobDecoder             = (*NullInt64)(nil)
	_ driver.Valuer              = (*NullInt64)(nil)
	_ proto.Marshaler            = (*NullInt64)(nil)
	_ proto.Unmarshaler          = (*NullInt64)(nil)
	_ proto.Sizer                = (*NullInt64)(nil)
	_ protoMarshalToer           = (*NullInt64)(nil)
	_ sql.Scanner                = (*NullInt64)(nil)
)

func TestMakeNullInt64(t *testing.T) {
	t.Parallel()
	i := MakeNullInt64(9223372036854775806)
	assertInt64(t, i, "MakeNullInt64()")

	zero := MakeNullInt64(0)
	if !zero.Valid {
		t.Error("MakeNullInt64(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "null", NullInt64{}.String())
	assert.Exactly(t, 8, zero.Size())
	assert.Exactly(t, 8, MakeNullInt64(125).Size())
	assert.Exactly(t, 8, MakeNullInt64(128).Size())
	assert.Exactly(t, "0", zero.String())
	assert.Exactly(t, "9223372036854775806", i.String())
	assert.Exactly(t, 0, NullInt64{}.Size())
}

func TestInt64_GoString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		i64  NullInt64
		want string
	}{
		{NullInt64{}, "dml.NullInt64{}"},
		{MakeNullInt64(2), "dml.MakeNullInt64(2)"},
	}
	for i, test := range tests {
		if have, want := fmt.Sprintf("%#v", test.i64), test.want; have != want {
			t.Errorf("%d: Have: %v Want: %v", i, have, want)
		}
	}
}

func TestNullInt64_JsonUnmarshal(t *testing.T) {
	t.Parallel()
	var i NullInt64
	err := json.Unmarshal(int64JSON, &i)
	maybePanic(err)
	assertInt64(t, i, "int64 json")

	var ni NullInt64
	err = json.Unmarshal(nullInt64JSON, &ni)
	maybePanic(err)
	assertInt64(t, ni, "sql.NullInt64 json")

	var null NullInt64
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt64(t, null, "null json")

	var badType NullInt64
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt64(t, badType, "wrong type json")

	var invalid NullInt64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullInt64(t, invalid, "invalid json")
}

func TestNullInt64_JsonUnmarshalNonIntegerNumber(t *testing.T) {
	t.Parallel()
	var i NullInt64
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int64")
	}
}

func TestNullInt64_JsonUnmarshalInt64Overflow(t *testing.T) {
	t.Parallel()
	int64Overflow := uint64(math.MaxInt64)

	// Max int64 should decode successfully
	var i NullInt64
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
	t.Parallel()
	var i NullInt64
	err := i.UnmarshalText([]byte("9223372036854775806"))
	maybePanic(err)
	assertInt64(t, i, "UnmarshalText() int64")

	var blank NullInt64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt64(t, blank, "UnmarshalText() empty int64")

	var null NullInt64
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullInt64(t, null, `UnmarshalText() "null"`)
}

func TestNullInt64_JsonMarshal(t *testing.T) {
	t.Parallel()
	i := MakeNullInt64(9223372036854775806)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "9223372036854775806", "non-empty json marshal")

	// invalid values should be encoded as null
	null := MakeNullInt64(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullInt64_MarshalText(t *testing.T) {
	t.Parallel()
	i := MakeNullInt64(9223372036854775806)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "9223372036854775806", "non-empty text marshal")

	// invalid values should be encoded as null
	null := MakeNullInt64(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullInt64_BinaryEncoding(t *testing.T) {
	t.Parallel()
	runner := func(b NullInt64, want []byte) func(*testing.T) {
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

			var decoded NullInt64
			require.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("-987654321", runner(MakeNullInt64(-987654321), []byte{0x4f, 0x97, 0x21, 0xc5, 0xff, 0xff, 0xff, 0xff}))
	t.Run("987654321", runner(MakeNullInt64(987654321), []byte{0xb1, 0x68, 0xde, 0x3a, 0x0, 0x0, 0x0, 0x0}))
	t.Run("-maxInt64", runner(MakeNullInt64(-math.MaxInt64), []byte{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x80}))
	t.Run("maxInt64", runner(MakeNullInt64(math.MaxInt64), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}))
	t.Run("null", runner(NullInt64{}, nil))
}

func TestInt64Pointer(t *testing.T) {
	t.Parallel()
	i := MakeNullInt64(9223372036854775806)
	ptr := i.Ptr()
	if *ptr != 9223372036854775806 {
		t.Errorf("bad %s int64: %#v ≠ %d\n", "pointer", ptr, 9223372036854775806)
	}

	null := MakeNullInt64(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt64IsZero(t *testing.T) {
	t.Parallel()
	i := MakeNullInt64(9223372036854775806)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := MakeNullInt64(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeNullInt64(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestInt64SetValid(t *testing.T) {
	t.Parallel()
	change := MakeNullInt64(0, false)
	assertNullInt64(t, change, "SetValid()")
	change.SetValid(9223372036854775806)
	assertInt64(t, change, "SetValid()")
}

func TestInt64Scan(t *testing.T) {
	t.Parallel()
	var i NullInt64
	err := i.Scan(9223372036854775806)
	maybePanic(err)
	assertInt64(t, i, "scanned int64")

	var null NullInt64
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt64(t, null, "scanned null")
}

func assertInt64(t *testing.T, i NullInt64, from string) {
	if i.Int64 != 9223372036854775806 {
		t.Errorf("bad %s int64: %d ≠ %d\n", from, i.Int64, 9223372036854775806)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt64(t *testing.T, i NullInt64, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullInt64(t *testing.T) {
	t.Parallel()
	assert.EqualValues(t, 1257894000, MakeNullInt64(1257894000).Int64)
	assert.True(t, MakeNullInt64(1257894000).Valid)
	assert.True(t, MakeNullInt64(0).Valid)
	v, err := MakeNullInt64(1257894000).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, 1257894000, v)
}

func TestNullInt64_Scan(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		var nv NullInt64
		require.NoError(t, nv.Scan(nil))
		assert.Exactly(t, NullInt64{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv NullInt64
		require.NoError(t, nv.Scan([]byte(`-1234567`)))
		assert.Exactly(t, MakeNullInt64(-1234567), nv)
	})
	t.Run("string unsupported", func(t *testing.T) {
		var nv NullInt64
		err := nv.Scan(`-1234567`)
		assert.True(t, errors.Is(err, errors.NotSupported), "Error behaviour should be errors.NotSupported")
		assert.Exactly(t, NullInt64{}, nv)
	})
}
