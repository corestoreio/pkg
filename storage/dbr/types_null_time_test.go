package dbr

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	intJSON      = []byte(`12345`)
	timeString   = "1977-05-25T20:21:21Z"
	timeJSON     = []byte(`"` + timeString + `"`)
	nullTimeJSON = []byte(`null`)
	timeValue, _ = time.Parse(time.RFC3339, timeString)
	timeObject   = []byte(`{"Time":"1977-05-25T20:21:21Z","Valid":true}`)
	nullObject   = []byte(`{"Time":"0001-01-01T00:00:00Z","Valid":false}`)
	badObject    = []byte(`{"hello": "world"}`)
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

func TestNullTime_Argument(t *testing.T) {
	t.Parallel()

	nss := []NullTime{
		{
			Time: timeValue,
		},
		{
			Time:  timeValue,
			Valid: true,
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
	assert.Exactly(t, []interface{}{interface{}(nil), timeValue}, args)
	assert.Exactly(t, "NULL'1977-05-25 20:21:21'", buf.String())
}

func TestArgNullTime(t *testing.T) {
	t.Parallel()

	args := ArgNullTime(MakeNullTime(timeValue), MakeNullTime(timeValue, false), MakeNullTime(timeValue))
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
		assert.Exactly(t, []interface{}{timeValue, interface{}(nil), timeValue}, argIF)
		assert.Exactly(t, "('1977-05-25 20:21:21',NULL,'1977-05-25 20:21:21')", buf.String())
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
		assert.Exactly(t, []interface{}{timeValue, interface{}(nil), timeValue}, argIF)
		assert.Exactly(t, "'1977-05-25 20:21:21'NULL'1977-05-25 20:21:21'", buf.String())
	})

	t.Run("single arg", func(t *testing.T) {
		args = ArgNullTime(MakeNullTime(timeValue))
		args = args.applyOperator(NotEqual)
		var buf bytes.Buffer
		argIF := make([]interface{}, 0, 2)
		for i := 0; i < args.len(); i++ {
			if err := args.writeTo(&buf, i); err != nil {
				t.Fatalf("%+v", err)
			}
		}
		argIF = args.toIFace(argIF)
		assert.Exactly(t, []interface{}{timeValue}, argIF)
		assert.Exactly(t, "'1977-05-25 20:21:21'", buf.String())
	})
}
