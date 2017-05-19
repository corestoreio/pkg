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

package dbr

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	int64JSON     = []byte(`9223372036854775806`)
	nullInt64JSON = []byte(`{"NullInt64":9223372036854775806,"Valid":true}`)
)

func TestMakeNullInt64(t *testing.T) {
	t.Parallel()
	i := MakeNullInt64(9223372036854775806)
	assertInt64(t, i, "MakeNullInt64()")

	zero := MakeNullInt64(0)
	if !zero.Valid {
		t.Error("MakeNullInt64(0)", "is invalid, but should be valid")
	}
}

func TestInt64_GoString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		i64  NullInt64
		want string
	}{
		{NullInt64{}, "dbr.NullInt64{}"},
		{MakeNullInt64(2), "dbr.MakeNullInt64(2)"},
	}
	for i, test := range tests {
		if have, want := fmt.Sprintf("%#v", test.i64), test.want; have != want {
			t.Errorf("%d: Have: %v Want: %v", i, have, want)
		}
	}
}

func TestUnmarshalInt64(t *testing.T) {
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

func TestUnmarshalNonIntegerNumber64(t *testing.T) {
	t.Parallel()
	var i NullInt64
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int64")
	}
}

func TestUnmarshalInt64Overflow(t *testing.T) {
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

func TestTextUnmarshalInt64(t *testing.T) {
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
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullInt64(t, null, `UnmarshalText() "null"`)
}

func TestMarshalInt64(t *testing.T) {
	t.Parallel()
	i := MakeNullInt64(9223372036854775806)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "9223372036854775806", "non-empty json marshal")

	// invalid values should be encoded as null
	null := MakeNullInt64(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalInt64Text(t *testing.T) {
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

func TestNullInt64_Argument(t *testing.T) {
	t.Parallel()

	nss := []NullInt64{
		{
			NullInt64: sql.NullInt64{
				Int64: 987654,
			},
		},
		{
			NullInt64: sql.NullInt64{
				Int64: 987653,
				Valid: true,
			},
		},
	}
	var buf bytes.Buffer
	args := make([]interface{}, 0, 2)
	for i, ns := range nss {
		args = ns.toIFace(args)
		ns.writeTo(&buf, i)

		arg := ns.applyOperator(NotBetween)
		assert.Exactly(t, NotBetween, arg.operator(), "Index %d", i)
		assert.Exactly(t, 1, arg.len(), "Length must be always one")
	}
	assert.Exactly(t, []interface{}{interface{}(nil), int64(987653)}, args)
	assert.Exactly(t, "NULL987653", buf.String())
}

func TestArgNullInt64(t *testing.T) {
	t.Parallel()

	args := ArgNullInt64(MakeNullInt64(987651), MakeNullInt64(987652, false), MakeNullInt64(987653))
	assert.Exactly(t, 3, args.len())
	args = args.applyOperator(NotIn)
	assert.Exactly(t, 1, args.len())

	t.Run("IN operator", func(t *testing.T) {
		args = args.applyOperator(In)
		var buf bytes.Buffer
		argIF := make([]interface{}, 0, 2)
		if err := args.writeTo(&buf, 0); err != nil {
			t.Fatalf("%+v", err)
		}
		argIF = args.toIFace(argIF)
		assert.Exactly(t, []interface{}{int64(987651), interface{}(nil), int64(987653)}, argIF)
		assert.Exactly(t, "(987651,NULL,987653)", buf.String())
	})

	t.Run("Not Equal operator", func(t *testing.T) {
		args = args.applyOperator(NotEqual)
		var buf bytes.Buffer
		argIF := make([]interface{}, 0, 2)
		for i := 0; i < args.len(); i++ {
			if err := args.writeTo(&buf, i); err != nil {
				t.Fatalf("%+v", err)
			}
		}
		argIF = args.toIFace(argIF)
		assert.Exactly(t, []interface{}{int64(987651), interface{}(nil), int64(987653)}, argIF)
		assert.Exactly(t, "987651NULL987653", buf.String())
	})

	t.Run("single arg", func(t *testing.T) {
		args = ArgNullInt64(MakeNullInt64(1234567))
		args = args.applyOperator(NotEqual)
		var buf bytes.Buffer
		argIF := make([]interface{}, 0, 2)
		for i := 0; i < args.len(); i++ {
			if err := args.writeTo(&buf, i); err != nil {
				t.Fatalf("%+v", err)
			}
		}
		argIF = args.toIFace(argIF)
		assert.Exactly(t, []interface{}{int64(1234567)}, argIF)
		assert.Exactly(t, "1234567", buf.String())
	})
}
