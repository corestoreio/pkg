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
	"encoding"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	_            json.Marshaler             = (*NullTime)(nil)
	_            json.Unmarshaler           = (*NullTime)(nil)
	_            encoding.TextMarshaler     = (*NullTime)(nil)
	_            encoding.TextUnmarshaler   = (*NullTime)(nil)
	_            encoding.BinaryMarshaler   = (*NullTime)(nil)
	_            encoding.BinaryUnmarshaler = (*NullTime)(nil)
	_            fmt.GoStringer             = (*NullTime)(nil)
	intJSON                                 = []byte(`12345`)
	timeString                              = "1977-05-25T20:21:21Z"
	timeJSON                                = []byte(`"` + timeString + `"`)
	nullTimeJSON                            = []byte(`null`)
	timeValue, _                            = time.Parse(time.RFC3339, timeString)
	timeObject                              = []byte(`{"Time":"1977-05-25T20:21:21Z","Valid":true}`)
	nullObject                              = []byte(`{"Time":"0001-01-01T00:00:00Z","Valid":false}`)
	badObject                               = []byte(`{"hello": "world"}`)
)

func TestUnmarshalTimeJSON(t *testing.T) {
	t.Parallel()
	var ti NullTime
	err := json.Unmarshal(timeJSON, &ti)
	maybePanic(err)
	assertTime(t, ti, "UnmarshalJSON() json")

	var null NullTime
	err = json.Unmarshal(nullTimeJSON, &null)
	maybePanic(err)
	assertNullTime(t, null, "null time json")

	var fromObject NullTime
	err = json.Unmarshal(timeObject, &fromObject)
	maybePanic(err)
	assertTime(t, fromObject, "time from object json")

	var nullFromObj NullTime
	err = json.Unmarshal(nullObject, &nullFromObj)
	maybePanic(err)
	assertNullTime(t, nullFromObj, "null from object json")

	var invalid NullTime
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullTime(t, invalid, "invalid from object json")

	var bad NullTime
	err = json.Unmarshal(badObject, &bad)
	if err == nil {
		t.Errorf("expected error: bad object")
	}
	assertNullTime(t, bad, "bad from object json")

	var wrongType NullTime
	err = json.Unmarshal(intJSON, &wrongType)
	if err == nil {
		t.Errorf("expected error: wrong type JSON")
	}
	assertNullTime(t, wrongType, "wrong type object json")
}

func TestUnmarshalTimeText(t *testing.T) {
	t.Parallel()
	ti := MakeNullTime(timeValue)
	txt, err := ti.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, timeString, "marshal text")

	var unmarshal NullTime
	err = unmarshal.UnmarshalText(txt)
	maybePanic(err)
	assertTime(t, unmarshal, "unmarshal text")

	var null NullTime
	err = null.UnmarshalText(nullJSON)
	maybePanic(err)
	assertNullTime(t, null, "unmarshal null text")
	txt, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, string(nullJSON), "marshal null text")

	var invalid NullTime
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNullTime(t, invalid, "bad string")
}

func TestMarshalTime(t *testing.T) {
	t.Parallel()
	ti := MakeNullTime(timeValue)
	data, err := json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(timeJSON), "non-empty json marshal")

	ti.Valid = false
	data, err = json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(nullJSON), "null json marshal")
}

func TestTimeFrom(t *testing.T) {
	t.Parallel()
	ti := MakeNullTime(timeValue)
	assertTime(t, ti, "MakeNullTime() time.Time")
}

func TestTimeSetValid(t *testing.T) {
	t.Parallel()
	var ti time.Time
	change := MakeNullTime(ti, false)
	assertNullTime(t, change, "SetValid()")
	change.SetValid(timeValue)
	assertTime(t, change, "SetValid()")
}

func TestTimePointer(t *testing.T) {
	t.Parallel()
	ti := MakeNullTime(timeValue)
	ptr := ti.Ptr()
	if *ptr != timeValue {
		t.Errorf("bad %s time: %#v ≠ %v\n", "pointer", ptr, timeValue)
	}

	var nt time.Time
	null := MakeNullTime(nt, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestTimeScanValue(t *testing.T) {
	t.Parallel()

	var ti NullTime
	err := ti.Scan(timeValue)
	maybePanic(err)
	assertTime(t, ti, "scanned time")
	if v, err := ti.Value(); v != timeValue || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var null NullTime
	err = null.Scan(nil)
	maybePanic(err)
	assertNullTime(t, null, "scanned null")
	if v, err := null.Value(); v != nil || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var wrong NullTime
	err = wrong.Scan(int64(42))
	if err == nil {
		t.Error("expected error")
	}
	assertNullTime(t, wrong, "scanned wrong")
}

func assertTime(t *testing.T, ti NullTime, from string) {
	if ti.Time != timeValue {
		t.Errorf("bad %v time: %v ≠ %v\n", from, ti.Time, timeValue)
	}
	if !ti.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullTime(t *testing.T, ti NullTime, from string) {
	if ti.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullTime(t *testing.T) {
	t.Parallel()
	var test = time.Now()
	assert.Equal(t, test, MakeNullTime(test).Time)
	assert.True(t, MakeNullTime(test).Valid)
	assert.True(t, MakeNullTime(time.Time{}).Valid)

	v, err := MakeNullTime(test).Value()
	assert.NoError(t, err)
	assert.Equal(t, test, v)
}
