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
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/gogo/protobuf/proto"
)

var (
	intJSON      = []byte(`12345`)
	timeString   = "1977-05-25T20:21:21Z"
	timeJSON     = []byte(`"` + timeString + `"`)
	nullTimeJSON = []byte(sqlStrNullLC)
	timeValue, _ = time.Parse(time.RFC3339, timeString)
	timeObject   = []byte(`{"Time":"1977-05-25T20:21:21Z","Valid":true}`)
	nullObject   = []byte(`{"Time":"0001-01-01T00:00:00Z","Valid":false}`)
	badObject    = []byte(`{"hello": "world"}`)
)

var (
	_ fmt.GoStringer             = (*Time)(nil)
	_ fmt.Stringer               = (*Time)(nil)
	_ json.Marshaler             = (*Time)(nil)
	_ json.Unmarshaler           = (*Time)(nil)
	_ encoding.BinaryMarshaler   = (*Time)(nil)
	_ encoding.BinaryUnmarshaler = (*Time)(nil)
	_ encoding.TextMarshaler     = (*Time)(nil)
	_ encoding.TextUnmarshaler   = (*Time)(nil)
	_ driver.Valuer              = (*Time)(nil)
	_ proto.Marshaler            = (*Time)(nil)
	_ proto.Unmarshaler          = (*Time)(nil)
	_ proto.Sizer                = (*Time)(nil)
	_ protoMarshalToer           = (*Time)(nil)
)

func TestNullTime_JsonUnmarshal(t *testing.T) {
	var ti Time
	err := json.Unmarshal(timeJSON, &ti)
	maybePanic(err)
	assertTime(t, ti, "UnmarshalJSON() json")

	var null Time
	err = json.Unmarshal(nullTimeJSON, &null)
	maybePanic(err)
	assertNullTime(t, null, "null time json")

	var fromObject Time
	err = json.Unmarshal(timeObject, &fromObject)
	maybePanic(err)
	assertTime(t, fromObject, "time from object json")

	var nullFromObj Time
	err = json.Unmarshal(nullObject, &nullFromObj)
	maybePanic(err)
	assertNullTime(t, nullFromObj, "null from object json")

	var invalid Time
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullTime(t, invalid, "invalid from object json")

	var bad Time
	err = json.Unmarshal(badObject, &bad)
	if err == nil {
		t.Errorf("expected error: bad object")
	}
	assertNullTime(t, bad, "bad from object json")

	var wrongType Time
	err = json.Unmarshal(intJSON, &wrongType)
	if err == nil {
		t.Errorf("expected error: wrong type JSON")
	}
	assertNullTime(t, wrongType, "wrong type object json")
}

func TestNullTime_UnmarshalText(t *testing.T) {
	ti := MakeTime(timeValue)
	txt, err := ti.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, timeString, "marshal text")

	var unmarshal Time
	err = unmarshal.UnmarshalText(txt)
	maybePanic(err)
	assertTime(t, unmarshal, "unmarshal text")

	var null Time
	err = null.UnmarshalText(nullJSON)
	maybePanic(err)
	assertNullTime(t, null, "unmarshal null text")
	txt, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, string(nullJSON), "marshal null text")

	var invalid Time
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNullTime(t, invalid, "bad string")
}

func TestNullTime_JsonMarshal(t *testing.T) {
	ti := MakeTime(timeValue)
	data, err := json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(timeJSON), "non-empty json marshal")

	ti.Valid = false
	data, err = json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(nullJSON), "null json marshal")
}

func TestNullTime_BinaryEncoding(t *testing.T) {
	runner := func(nv Time, wantMB, wantM []byte) func(*testing.T) {
		return func(t *testing.T) {
			dataMB, err := nv.MarshalBinary()
			assert.NoError(t, err)
			assert.Exactly(t, wantMB, dataMB, t.Name()+": MarshalBinary %q", dataMB)
			dataM, err := nv.Marshal()
			assert.NoError(t, err)
			assert.Exactly(t, wantM, dataM, t.Name()+": Marshal %q", dataM)

			var decoded Time
			assert.NoError(t, decoded.UnmarshalBinary(dataMB), "UnmarshalBinary")

			haveS := nv.String()
			wantS := decoded.String()
			if len(haveS) > 35 {
				// Not sure but there can be a bug in the Go stdlib ...
				// expected: "2006-01-02 15:04:05.000000002 -0400 UTC-4"
				// actual: "2006-01-02 15:04:05.000000002 -0400 -0400"
				// So cutting of after 35 chars does not compare the UTC-4 with the -0400 value.
				haveS = haveS[:35]
				wantS = wantS[:35]
			}
			assert.Exactly(t, haveS, wantS)
		}
	}
	t.Run("now fixed", runner(MakeTime(now()), []byte("\x01\x00\x00\x00\x0e\xbbK7\xe5\x00\x00\x00\x02\x00\x00"), []byte("\n\b\b\xe5\x81\xe5\x9d\x04\x10\x02\x10\x01")))
	t.Run("null", runner(Time{}, nil, []byte("\n\v\b\x80\x92\xb8Ø\xfe\xff\xff\xff\x01")))
}

func TestNullTime_Size(t *testing.T) {
	assert.Exactly(t, 13, Time{}.Size())
	assert.Exactly(t, 12, MakeTime(now()).Size())
}

func TestTimeFrom(t *testing.T) {
	ti := MakeTime(timeValue)
	assertTime(t, ti, "MakeTime() time.Time")
}

func TestTimeSetValid(t *testing.T) {
	var ti time.Time
	change := MakeTime(ti).SetNull() // stupid code
	assertNullTime(t, change, "SetValid()")
	assertTime(t, change.SetValid(timeValue), "SetValid()")
}

func TestTimePointer(t *testing.T) {
	ti := MakeTime(timeValue)
	ptr := ti.Ptr()
	if *ptr != timeValue {
		t.Errorf("bad %s time: %#v ≠ %v\n", "pointer", ptr, timeValue)
	}

	var null Time
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestTimeScanValue(t *testing.T) {
	var ti Time
	maybePanic(ti.Scan(timeValue))
	assertTime(t, ti, "scanned time")
	if v, err := ti.Value(); v != timeValue || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var null Time
	maybePanic(null.Scan(nil))
	assertNullTime(t, null, "scanned null")
	if v, err := null.Value(); v != nil || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var wrong Time
	if err := wrong.Scan(int64(42)); err == nil {
		t.Error("expected error")
	}
	assertNullTime(t, wrong, "scanned wrong")
}

func assertTime(t *testing.T, ti Time, from string) {
	if ti.Time != timeValue {
		t.Errorf("bad %v time: %v ≠ %v\n", from, ti.Time, timeValue)
	}
	if !ti.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullTime(t *testing.T, ti Time, from string) {
	if ti.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func TestNewNullTime(t *testing.T) {
	test := time.Now()
	assert.Equal(t, test, MakeTime(test).Time)
	assert.True(t, MakeTime(test).Valid)
	assert.True(t, MakeTime(time.Time{}).Valid)

	v, err := MakeTime(test).Value()
	assert.NoError(t, err)
	assert.Equal(t, test, v)
}
