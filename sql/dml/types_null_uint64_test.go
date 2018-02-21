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
	"testing"

	"github.com/corestoreio/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ fmt.GoStringer             = (*NullUint64)(nil)
	_ fmt.Stringer               = (*NullUint64)(nil)
	_ json.Marshaler             = (*NullUint64)(nil)
	_ json.Unmarshaler           = (*NullUint64)(nil)
	_ encoding.BinaryMarshaler   = (*NullUint64)(nil)
	_ encoding.BinaryUnmarshaler = (*NullUint64)(nil)
	_ encoding.TextMarshaler     = (*NullUint64)(nil)
	_ encoding.TextUnmarshaler   = (*NullUint64)(nil)
	_ gob.GobEncoder             = (*NullUint64)(nil)
	_ gob.GobDecoder             = (*NullUint64)(nil)
	_ driver.Valuer              = (*NullUint64)(nil)
	_ proto.Marshaler            = (*NullUint64)(nil)
	_ proto.Unmarshaler          = (*NullUint64)(nil)
	_ proto.Sizer                = (*NullUint64)(nil)
	_ protoMarshalToer           = (*NullUint64)(nil)
	_ sql.Scanner                = (*NullUint64)(nil)
)
var (
	nullUint64JSON = []byte(`{"NullUint64":9223372036854775806,"Valid":true}`)
)

func TestMakeNullUint64(t *testing.T) {
	t.Parallel()
	i := MakeNullUint64(9223372036854775806)
	assertUint64(t, i, "MakeNullUint64()")

	zero := MakeNullUint64(0)
	if !zero.Valid {
		t.Error("MakeNullUint64(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "null", NullUint64{}.String())
	assert.Exactly(t, 8, zero.Size())
	assert.Exactly(t, 8, MakeNullUint64(125).Size())
	assert.Exactly(t, 8, MakeNullUint64(128).Size())
	assert.Exactly(t, "0", zero.String())
	assert.Exactly(t, "9223372036854775806", i.String())
	assert.Exactly(t, 0, NullUint64{}.Size())
}

func TestUint64_GoString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		i64  NullUint64
		want string
	}{
		{NullUint64{}, "dml.NullUint64{}"},
		{MakeNullUint64(2), "dml.MakeNullUint64(2)"},
	}
	for i, test := range tests {
		if have, want := fmt.Sprintf("%#v", test.i64), test.want; have != want {
			t.Errorf("%d: Have: %v Want: %v", i, have, want)
		}
	}
}

func TestNullUint64_JsonUnmarshal(t *testing.T) {
	t.Parallel()
	var i NullUint64
	err := json.Unmarshal(int64JSON, &i)
	maybePanic(err)
	assertUint64(t, i, "int64 json")

	var ni NullUint64
	err = json.Unmarshal(nullUint64JSON, &ni)
	maybePanic(err)
	assertUint64(t, ni, "sql.NullUint64 json")

	var null NullUint64
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullUint64(t, null, "null json")

	var badType NullUint64
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint64(t, badType, "wrong type json")

	var invalid NullUint64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint64(t, invalid, "invalid json")
}

func TestNullUint64_JsonUnmarshalNonIntegerNumber(t *testing.T) {
	t.Parallel()
	var i NullUint64
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int64")
	}
}

func TestNullUint64_UnmarshalText(t *testing.T) {
	t.Parallel()
	var i NullUint64
	err := i.UnmarshalText([]byte("9223372036854775806"))
	maybePanic(err)
	assertUint64(t, i, "UnmarshalText() int64")

	var blank NullUint64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUint64(t, blank, "UnmarshalText() empty int64")

	var null NullUint64
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullUint64(t, null, `UnmarshalText() "null"`)
}

func TestNullUint64_JsonMarshal(t *testing.T) {
	t.Parallel()
	i := MakeNullUint64(9223372036854775806)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "9223372036854775806", "non-empty json marshal")

	// invalid values should be encoded as null
	null := MakeNullUint64(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullUint64_MarshalText(t *testing.T) {
	t.Parallel()
	i := MakeNullUint64(9223372036854775806)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "9223372036854775806", "non-empty text marshal")

	// invalid values should be encoded as null
	null := MakeNullUint64(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullUint64_BinaryEncoding(t *testing.T) {
	t.Parallel()
	runner := func(b NullUint64, want []byte) func(*testing.T) {
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

			var decoded NullUint64
			require.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("987654321", runner(MakeNullUint64(987654321), []byte{0xb1, 0x68, 0xde, 0x3a, 0x0, 0x0, 0x0, 0x0}))
	t.Run("maxUint64", runner(MakeNullUint64(math.MaxUint64), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}))
	t.Run("null", runner(NullUint64{}, nil))
}

func TestUint64Pointer(t *testing.T) {
	t.Parallel()
	i := MakeNullUint64(9223372036854775806)
	ptr := i.Ptr()
	if *ptr != 9223372036854775806 {
		t.Errorf("bad %s int64: %#v ≠ %d\n", "pointer", ptr, 9223372036854775806)
	}

	null := MakeNullUint64(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUint64IsZero(t *testing.T) {
	t.Parallel()
	i := MakeNullUint64(9223372036854775806)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := MakeNullUint64(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeNullUint64(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestUint64SetValid(t *testing.T) {
	t.Parallel()
	change := MakeNullUint64(0, false)
	assertNullUint64(t, change, "SetValid()")
	change.SetValid(9223372036854775806)
	assertUint64(t, change, "SetValid()")
}

func TestUint64Scan(t *testing.T) {
	t.Parallel()
	var i NullUint64
	err := i.Scan([]byte(`9223372036854775806`))
	maybePanic(err)
	assertUint64(t, i, "scanned int64")

	var null NullUint64
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUint64(t, null, "scanned null")
}

func assertUint64(t *testing.T, i NullUint64, from string) {
	if i.Uint64 != 9223372036854775806 {
		t.Errorf("bad %q int64: %d ≠ %d\n", from, i.Uint64, 9223372036854775806)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint64(t *testing.T, i NullUint64, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullUint64(t *testing.T) {
	t.Parallel()
	assert.EqualValues(t, 1257894000, MakeNullUint64(1257894000).Uint64)
	assert.True(t, MakeNullUint64(1257894000).Valid)
	assert.True(t, MakeNullUint64(0).Valid)
	v, err := MakeNullUint64(1257894000).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, 1257894000, v)
}

func TestNullUint64_Scan(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		var nv NullUint64
		require.NoError(t, nv.Scan(nil))
		assert.Exactly(t, NullUint64{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv NullUint64
		require.NoError(t, nv.Scan([]byte(`12345678910`)))
		assert.Exactly(t, MakeNullUint64(12345678910), nv)
	})
	t.Run("int64", func(t *testing.T) {
		var nv NullUint64
		require.NoError(t, nv.Scan(int64(12345678911)))
		assert.Exactly(t, MakeNullUint64(12345678911), nv)
	})
	t.Run("int", func(t *testing.T) {
		var nv NullUint64
		require.NoError(t, nv.Scan(int(12345678912)))
		assert.Exactly(t, MakeNullUint64(12345678912), nv)
	})
	t.Run("string unsupported", func(t *testing.T) {
		var nv NullUint64
		err := nv.Scan(`1234567`)
		assert.True(t, errors.Is(err, errors.NotSupported), "Error behaviour should be errors.NotSupported")
		assert.Exactly(t, NullUint64{}, nv)
	})
	t.Run("parse error negative", func(t *testing.T) {
		var nv NullUint64
		err := nv.Scan([]byte(`-1234567`))
		assert.EqualError(t, err, `strconv.ParseUint: parsing "-1234567": invalid syntax`)
		assert.Exactly(t, NullUint64{}, nv)
	})
}
