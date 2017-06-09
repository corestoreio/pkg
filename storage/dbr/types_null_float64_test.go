package dbr

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_               fmt.GoStringer = (*NullFloat64)(nil)
	float64JSON                    = []byte(`1.2345`)
	nullFloat64JSON                = []byte(`{"NullFloat64":1.2345,"Valid":true}`)
)

func TestFloat64From(t *testing.T) {
	f := MakeNullFloat64(1.2345)
	assertFloat64(t, f, "MakeNullFloat64()")

	zero := MakeNullFloat64(0)
	if !zero.Valid {
		t.Error("MakeNullFloat64(0)", "is invalid, but should be valid")
	}
}

func TestNullFloat64_GoString(t *testing.T) {
	var f NullFloat64
	assert.Exactly(t, "dbr.NullFloat64{}", f.GoString())
	f = MakeNullFloat64(3.1415926)
	assert.Exactly(t, "dbr.MakeNullFloat64(3.1415926)", f.GoString())
}

func TestUnmarshalFloat64(t *testing.T) {
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

func TestTextUnmarshalFloat64(t *testing.T) {
	var f NullFloat64
	err := f.UnmarshalText([]byte("1.2345"))
	maybePanic(err)
	assertFloat64(t, f, "UnmarshalText() float64")

	var blank NullFloat64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullFloat64(t, blank, "UnmarshalText() empty float64")

	var null NullFloat64
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullFloat64(t, null, `UnmarshalText() "null"`)
}

func TestMarshalFloat64(t *testing.T) {
	f := MakeNullFloat64(1.2345)
	data, err := json.Marshal(f)
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := MakeNullFloat64(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalFloat64Text(t *testing.T) {
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

func TestNullFloat64_Argument(t *testing.T) {
	t.Parallel()

	nss := []NullFloat64{
		{
			NullFloat64: sql.NullFloat64{
				Float64: math.Pi,
			},
		},
		{
			NullFloat64: sql.NullFloat64{
				Float64: math.Ln10,
				Valid:   true,
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
	assert.Exactly(t, []interface{}{interface{}(nil), math.Ln10}, args)
	assert.Exactly(t, "NULL2.302585092994046", buf.String())
}

func TestArgNullFloat64(t *testing.T) {
	t.Parallel()

	args := ArgNullFloat64(MakeNullFloat64(math.Phi), MakeNullFloat64(math.E, false), MakeNullFloat64(math.SqrtE))
	assert.Exactly(t, 3, args.len())
	args = args.applyOperator(NotIn)
	assert.Exactly(t, 3, args.len())

	t.Run("writeTo", func(t *testing.T) {
		args = args.applyOperator(NotEqual)
		var buf bytes.Buffer
		argIF := make([]interface{}, 0, 2)
		for i := 0; i < args.len(); i++ {
			if err := args.writeTo(&buf, i); err != nil {
				t.Fatalf("%+v", err)
			}
		}
		argIF = args.toIFace(argIF)
		assert.Exactly(t, []interface{}{math.Phi, interface{}(nil), math.SqrtE}, argIF)
		assert.Exactly(t, "1.618033988749895NULL1.6487212707001282", buf.String())
	})

	t.Run("single arg", func(t *testing.T) {
		args = ArgNullFloat64(MakeNullFloat64(math.Sqrt2))
		args = args.applyOperator(NotEqual)
		var buf bytes.Buffer
		argIF := make([]interface{}, 0, 2)
		for i := 0; i < args.len(); i++ {
			if err := args.writeTo(&buf, i); err != nil {
				t.Fatalf("%+v", err)
			}
		}
		argIF = args.toIFace(argIF)
		assert.Exactly(t, []interface{}{1.4142135623730951}, argIF)
		assert.Exactly(t, "1.4142135623730951", buf.String())
	})
}
