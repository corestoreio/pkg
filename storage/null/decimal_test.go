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
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/pseudo"
	"github.com/gogo/protobuf/proto"
)

// Holy guacamole. Those are many interface implementations. Maybe too much but who knows.
var (
	_ fmt.GoStringer             = (*Decimal)(nil)
	_ fmt.Stringer               = (*Decimal)(nil)
	_ json.Marshaler             = (*Decimal)(nil)
	_ json.Unmarshaler           = (*Decimal)(nil)
	_ encoding.BinaryMarshaler   = (*Decimal)(nil)
	_ encoding.BinaryUnmarshaler = (*Decimal)(nil)
	_ encoding.TextMarshaler     = (*Decimal)(nil)
	_ encoding.TextUnmarshaler   = (*Decimal)(nil)
	_ driver.Valuer              = (*Decimal)(nil)
	_ proto.Marshaler            = (*Decimal)(nil)
	_ proto.Unmarshaler          = (*Decimal)(nil)
	_ proto.Sizer                = (*Decimal)(nil)
	_ protoMarshalToer           = (*Decimal)(nil)
	_ sql.Scanner                = (*Decimal)(nil)
	_ pseudo.Faker               = (*Decimal)(nil)
)

func TestMakeDecimalBytes(t *testing.T) {
	tests := []struct {
		data          string
		wantPrecision uint64
		wantScale     int32
		wantNegative  bool
		wantErr       error
		wantStr       string
		wantValid     bool
	}{
		{"2681.7000", 26817000, 4, false, nil, "2681.7", true},
		{"-10.550000000000000000001", 0, 21, true, nil, "-10.550000000000000000001", true},
		{"-10.55000000000000000000", 0, 20, true, nil, "-10.55", true},
		{"-10.5500000000000000000000000", 0, 25, true, nil, "-10.55", true},
		{"0010.5500000000000000000000000", 0, 25, false, nil, "10.55", true},
		{"-0010.651234560000000000000000", 0, 24, true, nil, "-10.65123456", true},
		{"0010.55", 1055, 2, false, nil, "10.55", true},
		{"0010.00", 1000, 2, false, nil, "10", true},
		{"10000", 10000, 0, false, nil, "10000", true},
		{"47.11", 4711, 2, false, nil, "47.11", true},
		{"0010", 10, 0, false, nil, "10", true},
		{"0.000", 0, 3, false, nil, "0", true},
		{"0.010", 10, 3, false, nil, "0.01", true},
		{"00000.0000000", 0, 7, false, nil, "0", true},
		{"0", 0, 0, false, nil, "0", true},
		{".0", 0, 1, false, nil, "0", true},
		{"", 0, 0, false, nil, "NULL", false},
		{"0.1234567890123456789", 1234567890123456789, 19, false, nil, "0.1234567890123456789", true},
		{"0.01234567890123456789", 1234567890123456789, 20, false, nil, "0.01234567890123456789", true},
		{"-0.012345678901234567891", 12345678901234567891, 21, true, nil, "-0.012345678901234567891", true},
		{"-0.18446744073709551615", math.MaxUint64, 20, true, nil, "-0.18446744073709551615", true},
		{"-0.184467440737095516151", 0, 21, true, nil, "-0.184467440737095516151", true},
		{"0.0123456789012345678912345", 0, 25, false, nil, "0.0123456789012345678912345", true},
		{"123456789012345678912345678901234", 0, 0, false, nil, "123456789012345678912345678901234", true},
		{"123456789012345678912345678901234.123456789012345678912345678901234", 0, 33, false, nil, "123456789012345678912345678901234.123456789012345678912345678901234", true},
	}
	for i, test := range tests {
		haveD, haveErr := MakeDecimalBytes([]byte(test.data))
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "[%d] Err: %+v", i, haveErr)
		assert.Exactly(t, test.wantValid, haveD.Valid, "Index %d not valid", i)
		assert.Exactly(t, test.wantPrecision, haveD.Precision, "Index %d does not match precision", i)
		assert.Exactly(t, test.wantScale, haveD.Scale, "Index %d does not match scale", i)
		assert.Exactly(t, test.wantNegative, haveD.Negative, "Index %d does not match negative", i)
		assert.Exactly(t, test.wantStr, haveD.String(), "Index %d does not match String", i)
	}
}

func TestMakeDecimalInt64(t *testing.T) {
	d := MakeDecimalInt64(-math.MaxInt64, 13)
	assert.True(t, d.Negative)
	assert.Exactly(t, uint64(math.MaxInt64), d.Precision)
	assert.Exactly(t, int32(13), d.Scale)
}

func TestMakeDecimalFloat64(t *testing.T) {
	tests := []struct {
		have    float64
		want    string
		wantErr error
	}{
		{math.NaN(), "0", errors.New(`strconv.ParseUint: parsing "NaN": invalid syntax`)},
		{math.Inf(1), "0", errors.New(`strconv.ParseUint: parsing "Inf": invalid syntax`)},
		{math.Inf(-1), "-0", errors.New(`strconv.ParseUint: parsing "Inf": invalid syntax`)},
		{.00000000000000001, "0.00000000000000001", nil},
		{123.45678901234567, "123.45678901234567", nil},
		{123.456789012345678, "123.45678901234568", nil},
		{123.456789012345671, "123.45678901234567", nil},
		{987, "987", nil},
		{math.Phi * 4.01 * 5 / 9.099999, "3.565009344993927", nil},
	}
	for i, test := range tests {
		d, err := MakeDecimalFloat64(test.have)
		if test.wantErr != nil {
			assert.EqualError(t, err, test.wantErr.Error(), "Index %d", i)
			d.Negative = false
			assert.Exactly(t, Decimal{}, d, "Index %d", i)
		} else {
			assert.NoError(t, err, "Index %d", i)
			assert.Exactly(t, test.want, d.String(), "Index %d", i)
		}
	}
}

func TestDecimal_GoString(t *testing.T) {
	tests := []struct {
		have Decimal
		want string
	}{
		{Decimal{}, "null.Decimal{}"},
		{Decimal{
			Valid:     true,
			Precision: math.MaxUint64,
		}, "null.Decimal{Precision:18446744073709551615,Valid:true,}"},
		{Decimal{
			Valid:     true,
			Precision: math.MaxUint32,
			Scale:     16,
		}, "null.Decimal{Precision:4294967295,Scale:16,Valid:true,}"},
		{Decimal{
			Valid:     true,
			Precision: math.MaxUint16,
			Scale:     8,
			Negative:  true,
		}, "null.Decimal{Precision:65535,Scale:8,Negative:true,Valid:true,}"},
		{Decimal{
			Valid:     true,
			Precision: math.MaxUint16,
			Scale:     8,
			Negative:  true,
			Quote:     true,
		}, "null.Decimal{Precision:65535,Scale:8,Negative:true,Valid:true,Quote:true,}"},
		{
			MustMakeDecimalBytes([]byte("12345678912345.12345678")),
			`null.Decimal{PrecisionStr:"1234567891234512345678",Scale:8,Valid:true,}`,
		},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.have.GoString(), "Index %d", i)
	}
}

func TestMustMakeDecimalBytes(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.EqualError(t, err, "strconv.ParseUint: parsing \"helloWorld\": invalid syntax")
			} else {
				t.Errorf("Panic should contain an error but got:\n%+v", r)
			}
		} else {
			t.Error("Expecting a panic but got nothing")
		}
	}()
	MustMakeDecimalBytes([]byte(`helloWorld`))
}

func TestDecimal_String(t *testing.T) {
	tests := []struct {
		have Decimal
		want string
	}{
		{Decimal{}, "NULL"},
		{Decimal{
			Valid:     true,
			Precision: math.MaxUint64,
		}, "18446744073709551615"},
		{Decimal{
			Valid:     true,
			Precision: 1234,
		}, "1234"},
		{Decimal{
			Valid:     true,
			Precision: 1234,
			Negative:  true,
		}, "-1234"},
		{Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     1,
			Negative:  true,
		}, "-123.4"},
		{Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     1,
		}, "123.4"},
		{Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     2,
		}, "12.34"},
		{Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     4,
			Negative:  false,
		}, "0.1234"},
		{Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     11,
			Negative:  false,
		}, "0.00000001234"}, // 1234*10^-11
		{Decimal{
			Valid:     true,
			Precision: math.MaxUint64,
			Scale:     150,
			Negative:  true,
			// 18446744073709551615*10^-150
		}, "-0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000018446744073709551615"},
	}
	for i, test := range tests {
		val, err := test.have.Value()
		assert.NoError(t, err, "Index %d", i)
		assert.Exactly(t, test.want, val, "Index %d", i)

	}
}

func TestDecimal_MarshalJSON(t *testing.T) {
	runner := func(d Decimal, want string) func(*testing.T) {
		return func(t *testing.T) {
			raw, err := d.MarshalJSON()
			assert.NoError(t, err, t.Name())
			assert.Exactly(t, want, string(raw), t.Name())

			var d2 Decimal
			assert.NoError(t, d2.UnmarshalJSON(raw), t.Name())
			assert.Exactly(t, d, d2, t.Name())
		}
	}

	// TODO: Fuzzy testing

	t.Run("not valid", runner(Decimal{}, "null"))

	t.Run("quoted minus", runner(Decimal{
		Valid:     true,
		Precision: math.MaxUint64,
		Scale:     7, // large Scales not yet supported
		Negative:  true,
		Quote:     true,
	}, "\"-1844674407370.9551615\""))

	t.Run("quoted plus", runner(Decimal{
		Valid:     true,
		Precision: math.MaxUint32,
		Scale:     8, // large Scales not yet supported
		Quote:     true,
	}, "\"42.94967295\""))

	t.Run("unquoted", runner(Decimal{
		Valid:     true,
		Precision: 1234,
		Scale:     1,
		Negative:  true,
	}, "-123.4"))

	t.Run("-0.073", runner(Decimal{
		Valid:     true,
		Precision: 73,
		Scale:     3,
		Negative:  true,
	}, "-0.073"))

	t.Run("+9", runner(Decimal{
		Valid:     true,
		Precision: 9,
		Scale:     0,
	}, "9"))

	t.Run("Unmarshal null", func(t *testing.T) {
		dNull := Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     2,
		}
		assert.NoError(t, dNull.UnmarshalJSON([]byte("null")))
		assert.Exactly(t, Decimal{}, dNull)
	})
}

func TestDecimal_MarshalText(t *testing.T) {
	runner := func(d Decimal, want string) func(*testing.T) {
		return func(t *testing.T) {
			raw, err := d.MarshalText()
			assert.NoError(t, err, t.Name())
			assert.Exactly(t, want, string(raw), t.Name())
			d.Quote = false

			var d2 Decimal
			assert.NoError(t, d2.UnmarshalText(raw), t.Name())
			assert.Exactly(t, d, d2, t.Name())
		}
	}

	// TODO: Fuzzy testing

	t.Run("not valid", runner(Decimal{}, "NULL"))

	t.Run("quoted", runner(Decimal{
		Valid:     true,
		Precision: math.MaxUint64,
		Scale:     7, // large Scales not yet supported
		Negative:  true,
		Quote:     true,
	}, "-1844674407370.9551615")) // does not quote

	t.Run("unquoted", runner(Decimal{
		Valid:     true,
		Precision: 1234,
		Scale:     1,
		Negative:  true,
	}, "-123.4"))

	t.Run("Unmarshal empty", func(t *testing.T) {
		dNull := Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     2,
		}
		assert.NoError(t, dNull.UnmarshalText([]byte("")))
		assert.Exactly(t, Decimal{}, dNull)
	})
}

func TestDecimal_Int64(t *testing.T) {
	t.Run("1234", func(t *testing.T) {
		d := Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     2,
		}
		i, s := d.Int64()
		assert.Exactly(t, int64(1234), i)
		assert.Exactly(t, int32(2), s)
	})
	t.Run("-987654321", func(t *testing.T) {
		d := Decimal{
			Valid:     true,
			Precision: 987654321,
			Scale:     5,
			Negative:  true,
		}
		i, s := d.Int64()
		assert.Exactly(t, int64(-987654321), i)
		assert.Exactly(t, int32(5), s)
	})
	t.Run("overflow", func(t *testing.T) {
		d := Decimal{
			Valid:     true,
			Precision: math.MaxInt64 + 9876,
			Scale:     5,
			Negative:  true,
		}
		i, s := d.Int64()
		assert.Exactly(t, int64(0), i)
		assert.Exactly(t, int32(0), s)
	})
}

func TestDecimal_Float64(t *testing.T) {
	t.Run("0.0", func(t *testing.T) {
		d := Decimal{
			Valid: true,
		}
		f := d.Float64()
		assert.Exactly(t, 0.0, f)
	})
	t.Run("12.34", func(t *testing.T) {
		d := Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     2,
		}
		f := d.Float64()
		assert.Exactly(t, 12.34, f)
	})
	t.Run("-9876.54321", func(t *testing.T) {
		d := Decimal{
			Valid:     true,
			Precision: 987654321,
			Scale:     5,
			Negative:  true,
		}
		f := d.Float64()
		assert.Exactly(t, -9876.543210000002, f)
	})
	t.Run("overflow", func(t *testing.T) {
		d := Decimal{
			Valid:     true,
			Precision: math.MaxInt64 + 9876,
			Scale:     5,
			Negative:  true,
		}
		f := d.Float64()
		assert.Exactly(t, -9.223372036854788e+13, f)
	})
}

func TestDecimal_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var nv Decimal
		assert.NoError(t, nv.Scan(nil))
		assert.Exactly(t, Decimal{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv Decimal
		assert.NoError(t, nv.Scan([]byte(`-1234.567`)))
		assert.Exactly(t, MakeDecimalInt64(-1234567, 3), nv)
	})
	t.Run("string", func(t *testing.T) {
		var nv Decimal
		assert.NoError(t, nv.Scan(`-1234.567`))
		assert.Exactly(t, MakeDecimalInt64(-1234567, 3), nv)
	})
	t.Run("float64", func(t *testing.T) {
		var nv Decimal
		assert.NoError(t, nv.Scan(-1234.569))
		assert.Exactly(t, MakeDecimalInt64(-1234569, 3), nv)
	})
}

func TestDecimal_Equal(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		a := Decimal{Precision: 11, Scale: 1, Negative: true, Valid: true}
		b := Decimal{Precision: 11, Scale: 1, Negative: true, Valid: true}
		assert.True(t, a.Equal(b))
	})
	t.Run("unequal", func(t *testing.T) {
		a := Decimal{Precision: 13, Scale: 1, Negative: true, Valid: true}
		b := Decimal{Precision: 11, Scale: 1, Negative: true, Valid: true}
		assert.False(t, a.Equal(b))
	})
	t.Run("not valid1", func(t *testing.T) {
		a := Decimal{Precision: 13, Scale: 1, Negative: true}
		b := Decimal{Precision: 11, Scale: 1, Negative: true, Valid: true}
		assert.False(t, a.Equal(b))
	})
	t.Run("not valid2", func(t *testing.T) {
		a := Decimal{Precision: 13, Scale: 1, Negative: true, Valid: true}
		b := Decimal{Precision: 11, Scale: 1, Negative: true}
		assert.False(t, a.Equal(b))
	})
}

func TestDecimal_Fake(t *testing.T) {
	t.Run("PrecisionStr", func(t *testing.T) {
		d := &Decimal{}
		hasFakeDataApplied, err := d.Fake("PrecisionStr")
		assert.NoError(t, err)
		assert.True(t, hasFakeDataApplied)
	})
	t.Run("Precision", func(t *testing.T) {
		d := &Decimal{}
		hasFakeDataApplied, err := d.Fake("Precision")
		assert.NoError(t, err)
		assert.False(t, hasFakeDataApplied)
	})
}
