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
	float64JSON     = []byte(`1.2345`)
	nullFloat64JSON = []byte(`{"Float64":1.2345,"Valid":true}`)
)

var (
	_ fmt.GoStringer             = (*Float64)(nil)
	_ fmt.Stringer               = (*Float64)(nil)
	_ json.Marshaler             = (*Float64)(nil)
	_ json.Unmarshaler           = (*Float64)(nil)
	_ encoding.BinaryMarshaler   = (*Float64)(nil)
	_ encoding.BinaryUnmarshaler = (*Float64)(nil)
	_ encoding.TextMarshaler     = (*Float64)(nil)
	_ encoding.TextUnmarshaler   = (*Float64)(nil)
	_ gob.GobEncoder             = (*Float64)(nil)
	_ gob.GobDecoder             = (*Float64)(nil)
	_ driver.Valuer              = (*Float64)(nil)
	_ proto.Marshaler            = (*Float64)(nil)
	_ proto.Unmarshaler          = (*Float64)(nil)
	_ proto.Sizer                = (*Float64)(nil)
	_ protoMarshalToer           = (*Float64)(nil)
	_ sql.Scanner                = (*Float64)(nil)
)

func TestFloat64From(t *testing.T) {
	f := MakeFloat64(1.2345)
	assertFloat64(t, f, "MakeFloat64()")

	zero := MakeFloat64(0)
	if !zero.Valid {
		t.Error("MakeFloat64(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, "null", Float64{}.String())
	assert.Exactly(t, 8, zero.Size())
	assert.Exactly(t, "0", zero.String())
	assert.Exactly(t, "1.2345", f.String())
	assert.Exactly(t, 0, Float64{}.Size())
}

func TestNullFloat64_GoString(t *testing.T) {
	var f Float64
	assert.Exactly(t, "null.Float64{}", f.GoString())
	f = MakeFloat64(3.1415926)
	assert.Exactly(t, "null.MakeFloat64(3.1415926)", f.GoString())
}

func TestNullFloat64_JsonUnmarshal(t *testing.T) {
	var f Float64
	err := json.Unmarshal(float64JSON, &f)
	maybePanic(err)
	assertFloat64(t, f, "float64 json")

	var nf Float64
	err = json.Unmarshal(nullFloat64JSON, &nf)
	maybePanic(err)
	assertFloat64(t, nf, "sq.Float64 json")

	var null Float64
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullFloat64(t, null, "null json")

	var badType Float64
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullFloat64(t, badType, "wrong type json")

	var invalid Float64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
}

func TestNullFloat64_UnmarshalText(t *testing.T) {
	var f Float64
	err := f.UnmarshalText([]byte("1.2345"))
	maybePanic(err)
	assertFloat64(t, f, "UnmarshalText() float64")

	var blank Float64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullFloat64(t, blank, "UnmarshalText() empty float64")

	var null Float64
	err = null.UnmarshalText([]byte(sqlStrNullLC))
	maybePanic(err)
	assertNullFloat64(t, null, `UnmarshalText() "null"`)
}

func TestNullFloat64_JsonMarshal(t *testing.T) {
	f := MakeFloat64(1.2345)
	data, err := json.Marshal(f)
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := Float64{}
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, sqlStrNullLC, "null json marshal")
}

func TestNullFloat64_MarshalText(t *testing.T) {
	f := MakeFloat64(1.2345)
	data, err := f.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := MakeFloat64(0).SetNull()
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestNullFloat64_BinaryEncoding(t *testing.T) {
	runner := func(b Float64, want []byte) func(*testing.T) {
		return func(t *testing.T) {
			data, err := b.GobEncode()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": GobEncode")
			data, err = b.MarshalBinary()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": MarshalBinary")
			data, err = b.Marshal()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": Marshal")

			var decoded Float64
			assert.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("9.87654321", runner(MakeFloat64(9.87654321), []byte{0x33, 0xf6, 0x88, 0x45, 0xca, 0xc0, 0x23, 0x40}))
	t.Run("maxfloat64", runner(MakeFloat64(math.MaxFloat64), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xef, 0x7f}))
	t.Run("null", runner(Float64{}, nil))
}

func TestFloat64Pointer(t *testing.T) {
	f := MakeFloat64(1.2345)
	ptr := f.Ptr()
	if *ptr != 1.2345 {
		t.Errorf("bad %s float64: %#v ≠ %v\n", "pointer", ptr, 1.2345)
	}

	null := MakeFloat64(0).SetNull()
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s float64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestFloat64IsZero(t *testing.T) {
	f := MakeFloat64(1.2345)
	if f.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := MakeFloat64(0).SetNull()
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeFloat64(0)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestFloat64SetValid(t *testing.T) {
	var change Float64
	assertNullFloat64(t, change, "SetValid()")
	assertFloat64(t, change.SetValid(1.2345), "SetValid()")
}

func TestFloat64Scan(t *testing.T) {
	var f Float64
	err := f.Scan(1.2345)
	maybePanic(err)
	assertFloat64(t, f, "scanned float64")

	var null Float64
	err = null.Scan(nil)
	maybePanic(err)
	assertNullFloat64(t, null, "scanned null")
}

func assertFloat64(t *testing.T, f Float64, from string) {
	if f.Float64 != 1.2345 {
		t.Errorf("bad %s float64: %f ≠ %f\n", from, f.Float64, 1.2345)
	}
	if !f.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullFloat64(t *testing.T, f Float64, from string) {
	if f.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullFloat64(t *testing.T) {
	test := 1257894000.93445000001
	assert.Equal(t, test, MakeFloat64(test).Float64)
	assert.True(t, MakeFloat64(test).Valid)
	assert.True(t, MakeFloat64(0).Valid)
	v, err := MakeFloat64(test).Value()
	assert.NoError(t, err)
	assert.Equal(t, test, v)
}

func TestNullFloat64_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var nv Float64
		assert.NoError(t, nv.Scan(nil))
		assert.Exactly(t, Float64{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv Float64
		assert.NoError(t, nv.Scan([]byte(`-1234.567`)))
		assert.Exactly(t, MakeFloat64(-1234.567), nv)
	})
	t.Run("float64", func(t *testing.T) {
		var nv Float64
		assert.NoError(t, nv.Scan(-1234.569))
		assert.Exactly(t, MakeFloat64(-1234.569), nv)
	})
	t.Run("string unsupported", func(t *testing.T) {
		var nv Float64
		err := nv.Scan(`-123.4567`)
		assert.True(t, errors.Is(err, errors.NotSupported), "Error behaviour should be errors.NotSupported")
		assert.Exactly(t, Float64{}, nv)
	})
}
