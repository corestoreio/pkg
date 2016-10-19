package null

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

var (
	stringJSON      = []byte(`"test"`)
	blankStringJSON = []byte(`""`)
	nullStringJSON  = []byte(`{"String":"test","Valid":true}`)

	nullJSON    = []byte(`null`)
	invalidJSON = []byte(`:)`)
)

type stringInStruct struct {
	Test String `json:"test,omitempty"`
}

func TestStringFrom(t *testing.T) {
	str := StringFrom("test")
	assertStr(t, str, "StringFrom() string")

	zero := StringFrom("")
	if !zero.Valid {
		t.Error("StringFrom(0)", "is invalid, but should be valid")
	}
}

func TestStringFromPtr(t *testing.T) {
	s := "test"
	sptr := &s
	str := StringFromPtr(sptr)
	assertStr(t, str, "StringFromPtr() string")

	null := StringFromPtr(nil)
	assertNullStr(t, null, "StringFromPtr(nil)")
}

func TestUnmarshalString(t *testing.T) {
	var str String
	err := json.Unmarshal(stringJSON, &str)
	maybePanic(err)
	assertStr(t, str, "string json")

	var ns String
	err = json.Unmarshal(nullStringJSON, &ns)
	maybePanic(err)
	assertStr(t, ns, "sql.NullString json")

	var blank String
	err = json.Unmarshal(blankStringJSON, &blank)
	maybePanic(err)
	if !blank.Valid {
		t.Error("blank string should be valid")
	}

	var null String
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullStr(t, null, "null json")

	var badType String
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullStr(t, badType, "wrong type json")

	var invalid String
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullStr(t, invalid, "invalid json")
}

func TestTextUnmarshalString(t *testing.T) {
	var str String
	err := str.UnmarshalText([]byte("test"))
	maybePanic(err)
	assertStr(t, str, "UnmarshalText() string")

	var null String
	err = null.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullStr(t, null, "UnmarshalText() empty string")
}

func TestMarshalString(t *testing.T) {
	str := StringFrom("test")
	data, err := json.Marshal(str)
	maybePanic(err)
	assertJSONEquals(t, data, `"test"`, "non-empty json marshal")
	data, err = str.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "test", "non-empty text marshal")

	// empty values should be encoded as an empty string
	zero := StringFrom("")
	data, err = json.Marshal(zero)
	maybePanic(err)
	assertJSONEquals(t, data, `""`, "empty json marshal")
	data, err = zero.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "string marshal text")

	null := StringFromPtr(nil)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, `null`, "null json marshal")
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "string marshal text")
}

// Tests omitempty... broken until Go 1.4
// func TestMarshalStringInStruct(t *testing.T) {
// 	obj := stringInStruct{Test: StringFrom("")}
// 	data, err := json.Marshal(obj)
// 	maybePanic(err)
// 	assertJSONEquals(t, data, `{}`, "null string in struct")
// }

func TestStringPointer(t *testing.T) {
	str := StringFrom("test")
	ptr := str.Ptr()
	if *ptr != "test" {
		t.Errorf("bad %s string: %#v ≠ %s\n", "pointer", ptr, "test")
	}

	null := MakeString("", false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s string: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestStringIsZero(t *testing.T) {
	str := StringFrom("test")
	if str.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	blank := StringFrom("")
	if blank.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	empty := MakeString("", true)
	if empty.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := StringFromPtr(nil)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestStringSetValid(t *testing.T) {
	change := MakeString("", false)
	assertNullStr(t, change, "SetValid()")
	change.SetValid("test")
	assertStr(t, change, "SetValid()")
}

func TestStringScan(t *testing.T) {
	var str String
	err := str.Scan("test")
	maybePanic(err)
	assertStr(t, str, "scanned string")

	var null String
	err = null.Scan(nil)
	maybePanic(err)
	assertNullStr(t, null, "scanned null")
}

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

var _ fmt.GoStringer = (*String)(nil)

func TestString_GoString(t *testing.T) {
	s := MakeString("test", true)
	if have, want := s.GoString(), "null.StringFrom(`test`)"; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
	s = MakeString("test", false)
	if have, want := s.GoString(), "null.String{}"; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
	s = MakeString("te`st", true)
	// null.MakeString(`te`+"`"+`st`)
	gsWant := []byte{0x6e, 0x75, 0x6c, 0x6c, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x72, 0x6f, 0x6d, 0x28, 0x60, 0x74, 0x65, 0x60, 0x2b, 0x22, 0x60, 0x22, 0x2b, 0x60, 0x73, 0x74, 0x60, 0x29}
	if !bytes.Equal(gsWant, []byte(s.GoString())) {
		t.Errorf("Have: %#v Want: %v", s.GoString(), string(gsWant))
	}
}

func assertStr(t *testing.T, s String, from string) {
	if s.String != "test" {
		t.Errorf("bad %s string: %s ≠ %s\n", from, s.String, "test")
	}
	if !s.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullStr(t *testing.T, s String, from string) {
	if s.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}
