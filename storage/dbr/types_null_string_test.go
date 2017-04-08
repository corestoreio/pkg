package dbr_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

var (
	boolJSON        = []byte(`true`)
	stringJSON      = []byte(`"test"`)
	blankStringJSON = []byte(`""`)
	nullStringJSON  = []byte(`{"NullString":"test","Valid":true}`)

	nullJSON    = []byte(`null`)
	invalidJSON = []byte(`:)`)
)

func init() {
	dbr.JSONMarshalFn = json.Marshal
	dbr.JSONUnMarshalFn = json.Unmarshal
}

type stringInStruct struct {
	Test dbr.NullString `json:"test,omitempty"`
}

func TestStringFrom(t *testing.T) {
	str := dbr.MakeNullString("test")
	assertStr(t, str, "MakeNullString() string")

	zero := dbr.MakeNullString("")
	if !zero.Valid {
		t.Error("MakeNullString(0)", "is invalid, but should be valid")
	}
}

func TestUnmarshalString(t *testing.T) {
	var str dbr.NullString
	maybePanic(json.Unmarshal(stringJSON, &str))
	assertStr(t, str, "string json")

	var ns dbr.NullString
	maybePanic(json.Unmarshal(nullStringJSON, &ns))
	assertStr(t, ns, "sql.NullString json")

	var blank dbr.NullString
	maybePanic(json.Unmarshal(blankStringJSON, &blank))
	if !blank.Valid {
		t.Error("blank string should be valid")
	}

	var null dbr.NullString
	maybePanic(json.Unmarshal(nullJSON, &null))
	assertNullStr(t, null, "null json")

	var badType dbr.NullString
	err := json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullStr(t, badType, "wrong type json")

	var invalid dbr.NullString
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullStr(t, invalid, "invalid json")
}

func TestTextUnmarshalString(t *testing.T) {
	var str dbr.NullString
	err := str.UnmarshalText([]byte("test"))
	maybePanic(err)
	assertStr(t, str, "UnmarshalText() string")

	var null dbr.NullString
	err = null.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullStr(t, null, "UnmarshalText() empty string")

	var iv dbr.NullString
	err = iv.UnmarshalText([]byte{0x44, 0xff, 0x01})
	assert.True(t, errors.IsNotValid(err), "%+v", err)
}

func TestMarshalString(t *testing.T) {
	str := dbr.MakeNullString("test")
	data, err := json.Marshal(str)
	maybePanic(err)
	assertJSONEquals(t, data, `"test"`, "non-empty json marshal")
	data, err = str.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "test", "non-empty text marshal")

	// empty values should be encoded as an empty string
	zero := dbr.MakeNullString("")
	data, err = json.Marshal(zero)
	maybePanic(err)
	assertJSONEquals(t, data, `""`, "empty json marshal")
	data, err = zero.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "string marshal text")

	//null := dbr.NullStringFromPtr(nil)
	//data, err = json.Marshal(null)
	//maybePanic(err)
	//assertJSONEquals(t, data, `null`, "null json marshal")
	//data, err = null.MarshalText()
	//maybePanic(err)
	//assertJSONEquals(t, data, "", "string marshal text")
}

func TestStringPointer(t *testing.T) {
	str := dbr.MakeNullString("test")
	ptr := str.Ptr()
	if *ptr != "test" {
		t.Errorf("bad %s string: %#v ≠ %s\n", "pointer", ptr, "test")
	}

	null := dbr.MakeNullString("", false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s string: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestStringIsZero(t *testing.T) {
	str := dbr.MakeNullString("test")
	if str.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	blank := dbr.MakeNullString("")
	if blank.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	empty := dbr.MakeNullString("", true)
	if empty.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	//null := StringFromPtr(nil)
	//if !null.IsZero() {
	//	t.Errorf("IsZero() should be true")
	//}
}

func TestStringSetValid(t *testing.T) {
	change := dbr.MakeNullString("", false)
	assertNullStr(t, change, "SetValid()")
	change.SetValid("test")
	assertStr(t, change, "SetValid()")
}

func TestStringScan(t *testing.T) {
	var str dbr.NullString
	err := str.Scan("test")
	maybePanic(err)
	assertStr(t, str, "scanned string")

	var null dbr.NullString
	err = null.Scan(nil)
	maybePanic(err)
	assertNullStr(t, null, "scanned null")
}

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

var _ fmt.GoStringer = (*dbr.NullString)(nil)

func TestString_GoString(t *testing.T) {
	s := dbr.MakeNullString("test", true)
	assert.Exactly(t, "dbr.MakeNullString(`test`)", s.GoString())

	s = dbr.MakeNullString("test", false)
	assert.Exactly(t, "dbr.NullString{}", s.GoString())

	s = dbr.MakeNullString("te`st", true)
	gsWant := []byte("dbr.MakeNullString(`te`+\"`\"+`st`)")
	if !bytes.Equal(gsWant, []byte(s.GoString())) {
		t.Errorf("Have: %#v Want: %v", s.GoString(), string(gsWant))
	}
}

func assertStr(t *testing.T, s dbr.NullString, from string) {
	if s.String != "test" {
		t.Errorf("bad %s string: %s ≠ %s\n", from, s.String, "test")
	}
	if !s.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullStr(t *testing.T, s dbr.NullString, from string) {
	if s.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}
