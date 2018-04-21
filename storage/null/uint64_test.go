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
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ fmt.GoStringer             = (*Uint64)(nil)
	_ fmt.Stringer               = (*Uint64)(nil)
	_ json.Marshaler             = (*Uint64)(nil)
	_ json.Unmarshaler           = (*Uint64)(nil)
	_ encoding.BinaryMarshaler   = (*Uint64)(nil)
	_ encoding.BinaryUnmarshaler = (*Uint64)(nil)
	_ encoding.TextMarshaler     = (*Uint64)(nil)
	_ encoding.TextUnmarshaler   = (*Uint64)(nil)
	_ gob.GobEncoder             = (*Uint64)(nil)
	_ gob.GobDecoder             = (*Uint64)(nil)
	_ driver.Valuer              = (*Uint64)(nil)
	_ proto.Marshaler            = (*Uint64)(nil)
	_ proto.Unmarshaler          = (*Uint64)(nil)
	_ proto.Sizer                = (*Uint64)(nil)
	_ protoMarshalToer           = (*Uint64)(nil)
	_ sql.Scanner                = (*Uint64)(nil)
)
var (
	nullUint64JSON = []byte(`{"Uint64":9223372036854775806,"Valid":true}`)
)

func TestMakeNullUint64(t *testing.T) {
	t.Parallel()
	i := MakeUint64(9223372036854775806)
	assertUint64(t, i, "MakeUint64()")

	zero := MakeUint64(0)
	if !zero.Valid {
		t.Error("MakeUint64(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "null", Uint64{}.String())
	assert.Exactly(t, 8, zero.Size())
	assert.Exactly(t, 8, MakeUint64(125).Size())
	assert.Exactly(t, 8, MakeUint64(128).Size())
	assert.Exactly(t, "0", zero.String())
	assert.Exactly(t, "9223372036854775806", i.String())
	assert.Exactly(t, 0, Uint64{}.Size())
}

func TestUint64_GoString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		i64  Uint64
		want string
	}{
		{Uint64{}, "null.Uint64{}"},
		{MakeUint64(2), "null.MakeUint64(2)"},
	}
	for i, test := range tests {
		if have, want := fmt.Sprintf("%#v", test.i64), test.want; have != want {
			t.Errorf("%d: Have: %v Want: %v", i, have, want)
		}
	}
}

func TestNullUint64_JsonUnmarshal(t *testing.T) {
	t.Parallel()
	var i Uint64
	err := json.Unmarshal(int64JSON, &i)
	maybePanic(err)
	assertUint64(t, i, "int64 json")

	var ni Uint64
	err = json.Unmarshal(nullUint64JSON, &ni)
	maybePanic(err)
	assertUint64(t, ni, "sql.Uint64 json")

	var null Uint64
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullUint64(t, null, "null json")

	var badType Uint64
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint64(t, badType, "wrong type json")

	var invalid Uint64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint64(t, invalid, "invalid json")
}

func TestNullUint64_JsonUnmarshalNonIntegerNumber(t *testing.T) {
	t.Parallel()
	var i Uint64
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int64")
	}
}

func TestNullUint64_UnmarshalText(t *testing.T) {
	t.Parallel()
	var i Uint64
	err := i.UnmarshalText([]byte("9223372036854775806"))
	maybePanic(err)
	assertUint64(t, i, "UnmarshalText() int64")

	var blank Uint64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUint64(t, blank, "UnmarshalText() empty int64")

	var null Uint64
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullUint64(t, null, `UnmarshalText() "null"`)
}

func TestNullUint64_JsonMarshal(t *testing.T) {
	t.Parallel()
	i := MakeUint64(9223372036854775806)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "9223372036854775806", "non-empty json marshal")

	// invalid values should be encoded as null
	null := MakeUint64(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullUint64_MarshalText(t *testing.T) {
	t.Parallel()
	i := MakeUint64(9223372036854775806)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "9223372036854775806", "non-empty text marshal")

	// invalid values should be encoded as null
	null := MakeUint64(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullUint64_BinaryEncoding(t *testing.T) {
	t.Parallel()
	runner := func(b Uint64, want []byte) func(*testing.T) {
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

			var decoded Uint64
			require.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("987654321", runner(MakeUint64(987654321), []byte{0xb1, 0x68, 0xde, 0x3a, 0x0, 0x0, 0x0, 0x0}))
	t.Run("maxUint64", runner(MakeUint64(math.MaxUint64), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}))
	t.Run("null", runner(Uint64{}, nil))
}

func TestUint64Pointer(t *testing.T) {
	t.Parallel()
	i := MakeUint64(9223372036854775806)
	ptr := i.Ptr()
	if *ptr != 9223372036854775806 {
		t.Errorf("bad %s int64: %#v ≠ %d\n", "pointer", ptr, 9223372036854775806)
	}

	null := MakeUint64(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUint64IsZero(t *testing.T) {
	t.Parallel()
	i := MakeUint64(9223372036854775806)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := MakeUint64(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeUint64(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestUint64SetValid(t *testing.T) {
	t.Parallel()
	change := MakeUint64(0, false)
	assertNullUint64(t, change, "SetValid()")
	change.SetValid(9223372036854775806)
	assertUint64(t, change, "SetValid()")
}

func TestUint64Scan(t *testing.T) {
	t.Parallel()
	var i Uint64
	err := i.Scan([]byte(`9223372036854775806`))
	maybePanic(err)
	assertUint64(t, i, "scanned int64")

	var null Uint64
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUint64(t, null, "scanned null")
}

func assertUint64(t *testing.T, i Uint64, from string) {
	if i.Uint64 != 9223372036854775806 {
		t.Errorf("bad %q int64: %d ≠ %d\n", from, i.Uint64, 9223372036854775806)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint64(t *testing.T, i Uint64, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullUint64(t *testing.T) {
	t.Parallel()
	assert.EqualValues(t, 1257894000, MakeUint64(1257894000).Uint64)
	assert.True(t, MakeUint64(1257894000).Valid)
	assert.True(t, MakeUint64(0).Valid)
	v, err := MakeUint64(1257894000).Value()
	assert.NoError(t, err)
	assert.EqualValues(t, 1257894000, v)
}

func TestNullUint64_Scan(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		var nv Uint64
		require.NoError(t, nv.Scan(nil))
		assert.Exactly(t, Uint64{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv Uint64
		require.NoError(t, nv.Scan([]byte(`12345678910`)))
		assert.Exactly(t, MakeUint64(12345678910), nv)
	})
	t.Run("int64", func(t *testing.T) {
		var nv Uint64
		require.NoError(t, nv.Scan(int64(12345678911)))
		assert.Exactly(t, MakeUint64(12345678911), nv)
	})
	t.Run("int", func(t *testing.T) {
		var nv Uint64
		require.NoError(t, nv.Scan(int(12345678912)))
		assert.Exactly(t, MakeUint64(12345678912), nv)
	})
	t.Run("string unsupported", func(t *testing.T) {
		var nv Uint64
		err := nv.Scan(`1234567`)
		assert.True(t, errors.Is(err, errors.NotSupported), "Error behaviour should be errors.NotSupported")
		assert.Exactly(t, Uint64{}, nv)
	})
	t.Run("parse error negative", func(t *testing.T) {
		var nv Uint64
		err := nv.Scan([]byte(`-1234567`))
		assert.EqualError(t, err, `strconv.ParseUint: parsing "-1234567": invalid syntax`)
		assert.Exactly(t, Uint64{}, nv)
	})
}
