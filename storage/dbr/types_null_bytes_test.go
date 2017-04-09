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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	bytesJSON = []byte(`"hello"`)
)

func TestBytesFrom(t *testing.T) {
	t.Parallel()
	i := MakeNullBytes([]byte(`"hello"`))
	assertBytes(t, i, "MakeNullBytes()")

	zero := MakeNullBytes(nil)
	if zero.Valid {
		t.Error("MakeNullBytes(nil)", "is valid, but should be invalid")
	}

	zero = MakeNullBytes([]byte{})
	if !zero.Valid {
		t.Error("MakeNullBytes([]byte{})", "is invalid, but should be valid")
	}
}

func TestUnmarshalBytes(t *testing.T) {
	t.Parallel()
	var i NullBytes
	err := json.Unmarshal(bytesJSON, &i)
	maybePanic(err)
	assertBytes(t, i, "[]byte json")

	var ni NullBytes
	err = ni.UnmarshalJSON([]byte{})
	if ni.Valid == false {
		t.Errorf("expected Valid to be true, got false")
	}
	if !bytes.Equal(ni.Bytes, []byte("null")) {
		t.Errorf("Expected NullBytes to be nil slice, but was not: %#v %#v", ni.Bytes, []byte(`null`))
	}

	var null NullBytes
	err = null.UnmarshalJSON(nil)
	if null.Valid == false {
		t.Errorf("expected Valid to be true, got false")
	}
	if !bytes.Equal(null.Bytes, []byte(`null`)) {
		t.Errorf("Expected NullBytes to be []byte nil, but was not: %#v %#v", null.Bytes, []byte(`null`))
	}
}

func TestTextUnmarshalBytes(t *testing.T) {
	t.Parallel()
	var i NullBytes
	err := i.UnmarshalText([]byte(`"hello"`))
	maybePanic(err)
	assertBytes(t, i, "UnmarshalText() []byte")

	var blank NullBytes
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullBytes(t, blank, "UnmarshalText() empty []byte")
}

func TestMarshalBytes(t *testing.T) {
	t.Parallel()
	i := MakeNullBytes([]byte(`"hello"`))
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, `"hello"`, "non-empty json marshal")

	// invalid values should be encoded as null
	null := MakeNullBytes(nil, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalBytesText(t *testing.T) {
	t.Parallel()
	i := MakeNullBytes([]byte(`"hello"`))
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, `"hello"`, "non-empty text marshal")

	// invalid values should be encoded as null
	null := MakeNullBytes(nil, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestBytesPointer(t *testing.T) {
	t.Parallel()
	i := MakeNullBytes([]byte(`"hello"`))
	ptr := i.Ptr()
	if !bytes.Equal(*ptr, []byte(`"hello"`)) {
		t.Errorf("bad %s []byte: %#v ≠ %s\n", "pointer", ptr, `"hello"`)
	}

	null := MakeNullBytes(nil, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s []byte: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestBytesIsZero(t *testing.T) {
	t.Parallel()
	i := MakeNullBytes([]byte(`"hello"`))
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := MakeNullBytes(nil, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := MakeNullBytes(nil, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestBytesSetValid(t *testing.T) {
	t.Parallel()
	change := MakeNullBytes(nil, false)
	assertNullBytes(t, change, "SetValid()")
	change.SetValid([]byte(`"hello"`))
	assertBytes(t, change, "SetValid()")
}

func TestBytesScan(t *testing.T) {
	t.Parallel()
	var i NullBytes
	err := i.Scan(`"hello"`)
	maybePanic(err)
	assertBytes(t, i, "scanned []byte")

	var null NullBytes
	err = null.Scan(nil)
	maybePanic(err)
	assertNullBytes(t, null, "scanned null")
}

func assertBytes(t *testing.T, i NullBytes, from string) {
	if !bytes.Equal(i.Bytes, []byte(`"hello"`)) {
		t.Errorf("bad %s []byte: %#v ≠ %#v\n", from, string(i.Bytes), string([]byte(`"hello"`)))
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullBytes(t *testing.T, i NullBytes, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

var _ fmt.GoStringer = (*NullBytes)(nil)

func TestBytes_GoString(t *testing.T) {
	t.Parallel()
	b := MakeNullBytes([]byte("test"), true)
	if have, want := b.GoString(), "dbr.MakeNullBytes([]byte{0x74, 0x65, 0x73, 0x74})"; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
	b = MakeNullBytes([]byte("test"), false)
	if have, want := b.GoString(), "dbr.NullBytes{}"; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
	b = MakeNullBytes([]byte("te`st"), true)
	// null.MakeString(`te`+"`"+`st`)
	gsWant := []byte(`dbr.MakeNullBytes([]byte{0x74, 0x65, 0x60, 0x73, 0x74})`)
	if !bytes.Equal(gsWant, []byte(b.GoString())) {
		t.Errorf("Have: %#v Want: %v", b.GoString(), string(gsWant))
	}
}

func TestNullBytes_Argument(t *testing.T) {
	t.Parallel()

	nss := []NullBytes{
		{
			Bytes: []byte(`Not valid`),
		},
		{
			Bytes: []byte(`I'm valid'`),
			Valid: true,
		},
	}
	var buf bytes.Buffer
	args := make([]interface{}, 0, 2)
	for i, ns := range nss {
		ns.toIFace(&args)
		ns.writeTo(&buf, i)

		arg := ns.Operator(OperatorNotBetween)
		assert.Exactly(t, OperatorNotBetween, arg.operator(), "Index %d", i)
		assert.Exactly(t, 1, arg.len(), "Length must be always one")
	}
	assert.Exactly(t, []interface{}{interface{}(nil), []uint8{0x49, 0x27, 0x6d, 0x20, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x27}}, args)
	assert.Exactly(t, "NULL0x49276d2076616c696427", buf.String())
}

func TestArgNullBytes(t *testing.T) {
	t.Parallel()

	// enable these tests once we implement ArgNullBytes

	//args := ArgNullBytes(MakeNullBytes(math.Phi), MakeNullBytes(math.E, false), MakeNullBytes(math.SqrtE))
	//assert.Exactly(t, 3, args.len())
	//args = args.Operator(OperatorNotIn)
	//assert.Exactly(t, 1, args.len())
	//
	//t.Run("IN operator", func(t *testing.T) {
	//	args = args.Operator(OperatorIn)
	//	var buf bytes.Buffer
	//	argIF := make([]interface{}, 0, 2)
	//	if err := args.writeTo(&buf, 0); err != nil {
	//		t.Fatalf("%+v", err)
	//	}
	//	args.toIFace(&argIF)
	//	assert.Exactly(t, []interface{}{math.Phi, interface{}(nil), math.SqrtE}, argIF)
	//	assert.Exactly(t, "(1.618033988749895,NULL,1.6487212707001282)", buf.String())
	//})
	//
	//t.Run("Not Equal operator", func(t *testing.T) {
	//	args = args.Operator(OperatorNotEqual)
	//	var buf bytes.Buffer
	//	argIF := make([]interface{}, 0, 2)
	//	for i := 0; i < args.len(); i++ {
	//		if err := args.writeTo(&buf, i); err != nil {
	//			t.Fatalf("%+v", err)
	//		}
	//	}
	//	args.toIFace(&argIF)
	//	assert.Exactly(t, []interface{}{math.Phi, interface{}(nil), math.SqrtE}, argIF)
	//	assert.Exactly(t, "1.618033988749895NULL1.6487212707001282", buf.String())
	//})

	t.Run("single arg", func(t *testing.T) {
		args := MakeNullBytes([]byte("The quic\b\b\b\b\b\bk brown fo\u0007\u0007\u0007\u0007\u0007\u0007\u0007\u0007\u0007\u0007\u0007x... [Beeeep]")).
			Operator(OperatorNotEqual)

		var buf bytes.Buffer
		argIF := make([]interface{}, 0, 2)
		for i := 0; i < args.len(); i++ {
			if err := args.writeTo(&buf, i); err != nil {
				t.Fatalf("%+v", err)
			}
		}
		args.toIFace(&argIF)
		assert.Exactly(t, []interface{}{[]uint8{0x54, 0x68, 0x65, 0x20, 0x71, 0x75, 0x69, 0x63, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x6b, 0x20, 0x62, 0x72, 0x6f, 0x77, 0x6e, 0x20, 0x66, 0x6f, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x78, 0x2e, 0x2e, 0x2e, 0x20, 0x5b, 0x42, 0x65, 0x65, 0x65, 0x65, 0x70, 0x5d}}, argIF)
		assert.Exactly(t, "0x54686520717569630808080808086b2062726f776e20666f0707070707070707070707782e2e2e205b4265656565705d", buf.String())
	})
}
