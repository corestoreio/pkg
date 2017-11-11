// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"database/sql/driver"
	"encoding"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
)

var (
	float64JSON     = []byte(`1.2345`)
	nullFloat64JSON = []byte(`{"NullFloat64":1.2345,"Valid":true}`)
)

var (
	_ fmt.GoStringer             = (*NullFloat64)(nil)
	_ fmt.Stringer               = (*NullFloat64)(nil)
	_ json.Marshaler             = (*NullFloat64)(nil)
	_ json.Unmarshaler           = (*NullFloat64)(nil)
	_ encoding.BinaryMarshaler   = (*NullFloat64)(nil)
	_ encoding.BinaryUnmarshaler = (*NullFloat64)(nil)
	_ encoding.TextMarshaler     = (*NullFloat64)(nil)
	_ encoding.TextUnmarshaler   = (*NullFloat64)(nil)
	_ gob.GobEncoder             = (*NullFloat64)(nil)
	_ gob.GobDecoder             = (*NullFloat64)(nil)
	_ driver.Valuer              = (*NullFloat64)(nil)
	_ proto.Marshaler            = (*NullFloat64)(nil)
	_ proto.Unmarshaler          = (*NullFloat64)(nil)
	_ proto.Sizer                = (*NullFloat64)(nil)
	_ protoMarshalToer           = (*NullFloat64)(nil)
)

func TestFloat64From(t *testing.T) {
	t.Parallel()
	f := MakeNullFloat64(1.2345)
	assertFloat64(t, f, "MakeNullFloat64()")

	zero := MakeNullFloat64(0)
	if !zero.Valid {
		t.Error("MakeNullFloat64(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "null", NullFloat64{}.String())
	assert.Exactly(t, 8, zero.Size())
	assert.Exactly(t, "0", zero.String())
	assert.Exactly(t, "1.2345", f.String())
	assert.Exactly(t, 0, NullFloat64{}.Size())
}

func TestNullFloat64_GoString(t *testing.T) {
	var f NullFloat64
	assert.Exactly(t, "dml.NullFloat64{}", f.GoString())
	f = MakeNullFloat64(3.1415926)
	assert.Exactly(t, "dml.MakeNullFloat64(3.1415926)", f.GoString())
}

func TestNullFloat64_JsonUnmarshal(t *testing.T) {
	t.Parallel()
	var f NullFloat64
	err := json.Unmarshal(float64JSON, &f)
	maybePanic(err)
	assertFloat64(t, f, "float64 json")

	var nf NullFloat64
	err = json.Unmarshal(nullFloat64JSON, &nf)
	maybePanic(err)
	assertFloat64(t, nf, "sq.NullFloat64 json")

	var null NullFloat64
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullFloat64(t, null, "null json")

	var badType NullFloat64
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullFloat64(t, badType, "wrong type json")

	var invalid NullFloat64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
}

func TestNullFloat64_UnmarshalText(t *testing.T) {
	t.Parallel()
	var f NullFloat64
	err := f.UnmarshalText([]byte("1.2345"))
	maybePanic(err)
	assertFloat64(t, f, "UnmarshalText() float64")

	var blank NullFloat64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullFloat64(t, blank, "UnmarshalText() empty float64")

	var null NullFloat64
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullFloat64(t, null, `UnmarshalText() "null"`)
}

func TestNullFloat64_JsonMarshal(t *testing.T) {
	t.Parallel()
	f := MakeNullFloat64(1.2345)
	data, err := json.Marshal(f)
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := MakeNullFloat64(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullFloat64_MarshalText(t *testing.T) {
	t.Parallel()
	f := MakeNullFloat64(1.2345)
	data, err := f.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := MakeNullFloat64(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullFloat64_BinaryEncoding(t *testing.T) {
	t.Parallel()
	runner := func(b NullFloat64, want []byte) func(*testing.T) {
		return func(t *testing.T) {
			data, err := b.GobEncode()
			require.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": GobEncode")
			data, err = b.MarshalBinary()
			require.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": MarshalBinary")
			data, err = b.Marshal()
			require.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": Marshal")

			var decoded NullFloat64
			require.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("9.87654321", runner(MakeNullFloat64(9.87654321), []byte{0x33, 0xf6, 0x88, 0x45, 0xca, 0xc0, 0x23, 0x40}))
	t.Run("maxfloat64", runner(MakeNullFloat64(math.MaxFloat64), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xef, 0x7f}))
	t.Run("null", runner(NullFloat64{}, nil))
}

func TestFloat64Pointer(t *testing.T) {
	f := MakeNullFloat64(1.2345)
	ptr := f.Ptr()
	if *ptr != 1.2345 {
		t.Errorf("bad %s float64: %#v ≠ %v\n", "pointer", ptr, 1.2345)
	}

	null := MakeNullFloat64(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s float64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestFloat64IsZero(t *testing.T) {
	f := MakeNullFloat64(1.2345)
	if f.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := MakeNullFloat64(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeNullFloat64(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestFloat64SetValid(t *testing.T) {
	change := MakeNullFloat64(0, false)
	assertNullFloat64(t, change, "SetValid()")
	change.SetValid(1.2345)
	assertFloat64(t, change, "SetValid()")
}

func TestFloat64Scan(t *testing.T) {
	var f NullFloat64
	err := f.Scan(1.2345)
	maybePanic(err)
	assertFloat64(t, f, "scanned float64")

	var null NullFloat64
	err = null.Scan(nil)
	maybePanic(err)
	assertNullFloat64(t, null, "scanned null")
}

func assertFloat64(t *testing.T, f NullFloat64, from string) {
	if f.Float64 != 1.2345 {
		t.Errorf("bad %s float64: %f ≠ %f\n", from, f.Float64, 1.2345)
	}
	if !f.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullFloat64(t *testing.T, f NullFloat64, from string) {
	if f.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullFloat64(t *testing.T) {
	t.Parallel()
	var test = 1257894000.93445000001
	assert.Equal(t, test, MakeNullFloat64(test).Float64)
	assert.True(t, MakeNullFloat64(test).Valid)
	assert.True(t, MakeNullFloat64(0).Valid)
	v, err := MakeNullFloat64(test).Value()
	assert.NoError(t, err)
	assert.Equal(t, test, v)
}
