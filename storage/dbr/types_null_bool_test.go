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
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	boolJSON     = []byte(`true`)
	nullBoolJSON = []byte(`{"NullBool":true,"Valid":true}`)
)

func TestMakeNullBool(t *testing.T) {
	t.Parallel()
	b := MakeNullBool(true)
	assertBool(t, b, "MakeNullBool()")

	zero := MakeNullBool(false)
	if !zero.Valid {
		t.Error("MakeNullBool(false)", "is invalid, but should be valid")
	}
}

func TestUnmarshalBool(t *testing.T) {
	t.Parallel()
	var b NullBool
	err := json.Unmarshal(boolJSON, &b)
	maybePanic(err)
	assertBool(t, b, "bool json")

	var nb NullBool
	err = json.Unmarshal(nullBoolJSON, &nb)
	maybePanic(err)
	assertBool(t, nb, "sq.NullBool json")

	var null NullBool
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullBool(t, null, "null json")

	var badType NullBool
	err = json.Unmarshal(intJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullBool(t, badType, "wrong type json")

	var invalid NullBool
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
}

func TestTextUnmarshalBool(t *testing.T) {
	t.Parallel()

	var b NullBool
	err := b.UnmarshalText([]byte("true"))
	maybePanic(err)
	assertBool(t, b, "UnmarshalText() bool")

	var zero NullBool
	err = zero.UnmarshalText([]byte("false"))
	maybePanic(err)
	assertFalseBool(t, zero, "UnmarshalText() false")

	var blank NullBool
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullBool(t, blank, "UnmarshalText() empty bool")

	var null NullBool
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullBool(t, null, `UnmarshalText() "null"`)

	var invalid NullBool
	err = invalid.UnmarshalText([]byte(":D"))
	if err == nil {
		panic("err should not be nil")
	}
	assertNullBool(t, invalid, "invalid json")
}

func TestMarshalBool(t *testing.T) {
	t.Parallel()

	b := MakeNullBool(true)
	data, err := json.Marshal(b)
	maybePanic(err)
	assertJSONEquals(t, data, "true", "non-empty json marshal")

	zero := MakeNullBool(false, true)
	data, err = json.Marshal(zero)
	maybePanic(err)
	assertJSONEquals(t, data, "false", "zero json marshal")

	// invalid values should be encoded as null
	null := MakeNullBool(false, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalBoolText(t *testing.T) {
	t.Parallel()

	b := MakeNullBool(true)
	data, err := b.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "true", "non-empty text marshal")

	zero := MakeNullBool(false, true)
	data, err = zero.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "false", "zero text marshal")

	// invalid values should be encoded as null
	null := MakeNullBool(false, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestBoolPointer(t *testing.T) {
	t.Parallel()

	b := MakeNullBool(true)
	ptr := b.Ptr()
	if !*ptr {
		t.Errorf("bad %s bool: %#v ≠ %v\n", "pointer", ptr, true)
	}

	null := MakeNullBool(false, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s bool: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestBoolIsZero(t *testing.T) {
	t.Parallel()

	b := MakeNullBool(true)
	if b.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := MakeNullBool(false, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeNullBool(false, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestBoolSetValid(t *testing.T) {
	t.Parallel()

	change := MakeNullBool(false, false)
	assertNullBool(t, change, "SetValid()")
	change.SetValid(true)
	assertBool(t, change, "SetValid()")
}

func TestBoolScan(t *testing.T) {
	t.Parallel()

	var b NullBool
	err := b.Scan(true)
	maybePanic(err)
	assertBool(t, b, "scanned bool")

	var null NullBool
	err = null.Scan(nil)
	maybePanic(err)
	assertNullBool(t, null, "scanned null")
}

func assertBool(t *testing.T, b NullBool, from string) {
	if !b.Bool {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Bool, true)
	}
	if !b.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertFalseBool(t *testing.T, b NullBool, from string) {
	if b.Bool {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Bool, false)
	}
	if !b.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullBool(t *testing.T, b NullBool, from string) {
	if b.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNullBool_Argument(t *testing.T) {
	t.Parallel()

	nss := []NullBool{
		{
			NullBool: sql.NullBool{
				Bool: false,
			},
		},
		{
			NullBool: sql.NullBool{
				Bool:  true,
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
	assert.Exactly(t, []interface{}{interface{}(nil), true}, args)
	assert.Exactly(t, "NULL1", buf.String())
}
