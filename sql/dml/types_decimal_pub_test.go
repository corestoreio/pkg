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

package dml_test

import (
	"database/sql/driver"
	"encoding"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/corestoreio/cspkg/sql/dml"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Holy guacamole. Those are many interface implementations. Maybe too much but who knows.
var (
	_ fmt.GoStringer             = (*dml.Decimal)(nil)
	_ fmt.Stringer               = (*dml.Decimal)(nil)
	_ json.Marshaler             = (*dml.Decimal)(nil)
	_ json.Unmarshaler           = (*dml.Decimal)(nil)
	_ encoding.BinaryMarshaler   = (*dml.Decimal)(nil)
	_ encoding.BinaryUnmarshaler = (*dml.Decimal)(nil)
	_ encoding.TextMarshaler     = (*dml.Decimal)(nil)
	_ encoding.TextUnmarshaler   = (*dml.Decimal)(nil)
	_ gob.GobEncoder             = (*dml.Decimal)(nil)
	_ gob.GobDecoder             = (*dml.Decimal)(nil)
	_ driver.Valuer              = (*dml.Decimal)(nil)
	_ proto.Marshaler            = (*dml.Decimal)(nil)
	_ proto.Unmarshaler          = (*dml.Decimal)(nil)
	_ proto.Sizer                = (*dml.Decimal)(nil)
	_ protoMarshalToer           = (*dml.Decimal)(nil)
)

func TestMakeDecimalInt64(t *testing.T) {
	t.Parallel()
	d := dml.MakeDecimalInt64(-math.MaxInt64, 13)
	assert.True(t, d.Negative)
	assert.Exactly(t, uint64(math.MaxInt64), d.Precision)
	assert.Exactly(t, int32(13), d.Scale)
}

func TestMakeDecimalFloat64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have    float64
		want    string
		wantErr error
	}{
		{math.NaN(), "0", nil},
		{math.Inf(1), "0", nil},
		{math.Inf(-1), "-0", nil},
		{.00000000000000001, "0.00000000000000001", nil},
		{123.45678901234567, "123.45678901234567", nil},
		{123.456789012345678, "123.45678901234568", nil},
		{123.456789012345671, "123.45678901234567", nil},
		{987, "987", nil},
		{math.MaxFloat64, strconv.FormatUint(math.MaxUint64, 10), nil},
		{math.Phi * 4.01 * 5 / 9.099999, "3.565009344993927", nil},
	}
	for i, test := range tests {
		d, err := dml.MakeDecimalFloat64(test.have)
		if test.wantErr != nil {
			assert.EqualError(t, err, test.wantErr.Error())
			assert.Exactly(t, dml.Decimal{}, d)
			continue
		}
		assert.Exactly(t, test.want, d.String(), "Index %d", i)
	}
}

func TestDecimal_GoString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		have dml.Decimal
		want string
	}{
		{dml.Decimal{}, "dml.Decimal{}"},
		{dml.Decimal{
			Valid:     true,
			Precision: math.MaxUint64,
		}, "dml.Decimal{Precision:18446744073709551615,Valid:true,}"},
		{dml.Decimal{
			Valid:     true,
			Precision: math.MaxUint32,
			Scale:     16,
		}, "dml.Decimal{Precision:4294967295,Scale:16,Valid:true,}"},
		{dml.Decimal{
			Valid:     true,
			Precision: math.MaxUint16,
			Scale:     8,
			Negative:  true,
		}, "dml.Decimal{Precision:65535,Scale:8,Negative:true,Valid:true,}"},
		{dml.Decimal{
			Valid:     true,
			Precision: math.MaxUint16,
			Scale:     8,
			Negative:  true,
			Quote:     true,
		}, "dml.Decimal{Precision:65535,Scale:8,Negative:true,Valid:true,Quote:true,}"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.have.GoString(), "Index %d", i)
	}
}

func TestDecimal_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		have dml.Decimal
		want string
	}{
		{dml.Decimal{}, "NULL"},
		{dml.Decimal{
			Valid:     true,
			Precision: math.MaxUint64,
		}, "18446744073709551615"},
		{dml.Decimal{
			Valid:     true,
			Precision: 1234,
		}, "1234"},
		{dml.Decimal{
			Valid:     true,
			Precision: 1234,
			Negative:  true,
		}, "-1234"},
		{dml.Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     1,
			Negative:  true,
		}, "-123.4"},
		{dml.Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     1,
		}, "123.4"},
		{dml.Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     2,
		}, "12.34"},
		{dml.Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     4,
			Negative:  false,
		}, "0.1234"},
		{dml.Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     11,
			Negative:  false,
		}, "0.00000001234"}, // 1234*10^-11
		{dml.Decimal{
			Valid:     true,
			Precision: math.MaxUint64,
			Scale:     150,
			Negative:  true,
			// 18446744073709551615*10^-150
		}, "-0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000018446744073709551615"},
	}
	for i, test := range tests {
		val, err := test.have.Value()
		require.NoError(t, err, "Index %d", i)
		assert.Exactly(t, test.want, val, "Index %d", i)

	}
}

func TestDecimal_MarshalJSON(t *testing.T) {
	t.Parallel()

	runner := func(d dml.Decimal, want string) func(*testing.T) {
		return func(t *testing.T) {
			raw, err := d.MarshalJSON()
			require.NoError(t, err, t.Name())
			assert.Exactly(t, want, string(raw), t.Name())

			var d2 dml.Decimal
			require.NoError(t, d2.UnmarshalJSON(raw), t.Name())
			assert.Exactly(t, d, d2, t.Name())
		}
	}

	// TODO: Fuzzy testing

	t.Run("not valid", runner(dml.Decimal{}, "null"))

	t.Run("quoted minus", runner(dml.Decimal{
		Valid:     true,
		Precision: math.MaxUint64,
		Scale:     7, // large Scales not yet supported
		Negative:  true,
		Quote:     true,
	}, "\"-1844674407370.9551615\""))

	t.Run("quoted plus", runner(dml.Decimal{
		Valid:     true,
		Precision: math.MaxUint32,
		Scale:     8, // large Scales not yet supported
		Quote:     true,
	}, "\"42.94967295\""))

	t.Run("unquoted", runner(dml.Decimal{
		Valid:     true,
		Precision: 1234,
		Scale:     1,
		Negative:  true,
	}, "-123.4"))

	t.Run("-0.073", runner(dml.Decimal{
		Valid:     true,
		Precision: 73,
		Scale:     3,
		Negative:  true,
	}, "-0.073"))

	t.Run("+9", runner(dml.Decimal{
		Valid:     true,
		Precision: 9,
		Scale:     0,
	}, "9"))

	t.Run("Unmarshal null", func(t *testing.T) {
		dNull := dml.Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     2,
		}
		require.NoError(t, dNull.UnmarshalJSON([]byte("null")))
		assert.Exactly(t, dml.Decimal{}, dNull)
	})
}

func TestDecimal_MarshalText(t *testing.T) {
	t.Parallel()

	runner := func(d dml.Decimal, want string) func(*testing.T) {
		return func(t *testing.T) {
			raw, err := d.MarshalText()
			require.NoError(t, err, t.Name())
			assert.Exactly(t, want, string(raw), t.Name())
			d.Quote = false

			var d2 dml.Decimal
			require.NoError(t, d2.UnmarshalText(raw), t.Name())
			assert.Exactly(t, d, d2, t.Name())
		}
	}

	// TODO: Fuzzy testing

	t.Run("not valid", runner(dml.Decimal{}, "NULL"))

	t.Run("quoted", runner(dml.Decimal{
		Valid:     true,
		Precision: math.MaxUint64,
		Scale:     7, // large Scales not yet supported
		Negative:  true,
		Quote:     true,
	}, "-1844674407370.9551615")) // does not quote

	t.Run("unquoted", runner(dml.Decimal{
		Valid:     true,
		Precision: 1234,
		Scale:     1,
		Negative:  true,
	}, "-123.4"))

	t.Run("Unmarshal empty", func(t *testing.T) {
		dNull := dml.Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     2,
		}
		require.NoError(t, dNull.UnmarshalText([]byte("")))
		assert.Exactly(t, dml.Decimal{}, dNull)
	})
}

func TestDecimal_GobEncode(t *testing.T) {
	t.Parallel()

	runner := func(d dml.Decimal, want []byte) func(*testing.T) {
		return func(t *testing.T) {
			raw, err := d.GobEncode()
			require.NoError(t, err, t.Name())
			assert.Exactly(t, want, raw, t.Name())

			var d2 dml.Decimal
			require.NoError(t, d2.GobDecode(raw), t.Name())
			assert.Exactly(t, d, d2, t.Name())
		}
	}

	// TODO: Fuzzy testing

	t.Run("not valid", runner(dml.Decimal{}, nil))

	t.Run("quoted", runner(dml.Decimal{
		Valid:     true,
		Precision: math.MaxUint64 - 987654,
		Scale:     7, // large Scales not yet supported
		Negative:  true,
		Quote:     true,
	}, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xf0, 0xed, 0xf9, 0x0, 0x0, 0x0, 0x7, 0x0, 0xf})) // does not quote

	t.Run("unquoted", runner(dml.Decimal{
		Valid:     true,
		Precision: 1234,
		Scale:     2,
		Negative:  true,
	}, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0xd2, 0x0, 0x0, 0x0, 0x2, 0x0, 0xb}))

	t.Run("GobDecode nil", func(t *testing.T) {
		dNull := dml.Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     2,
		}
		require.NoError(t, dNull.GobDecode([]byte("")))
		assert.Exactly(t, dml.Decimal{}, dNull)
	})
}

func TestDecimal_Int64(t *testing.T) {
	t.Parallel()

	t.Run("1234", func(t *testing.T) {
		d := dml.Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     2,
		}
		i, s := d.Int64()
		assert.Exactly(t, int64(1234), i)
		assert.Exactly(t, int32(2), s)
	})
	t.Run("-987654321", func(t *testing.T) {
		d := dml.Decimal{
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
		d := dml.Decimal{
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
	t.Parallel()

	t.Run("0.0", func(t *testing.T) {
		d := dml.Decimal{
			Valid: true,
		}
		f := d.Float64()
		assert.Exactly(t, 0.0, f)
	})
	t.Run("12.34", func(t *testing.T) {
		d := dml.Decimal{
			Valid:     true,
			Precision: 1234,
			Scale:     2,
		}
		f := d.Float64()
		assert.Exactly(t, 12.34, f)
	})
	t.Run("-9876.54321", func(t *testing.T) {
		d := dml.Decimal{
			Valid:     true,
			Precision: 987654321,
			Scale:     5,
			Negative:  true,
		}
		f := d.Float64()
		assert.Exactly(t, -9876.543210000002, f)
	})
	t.Run("overflow", func(t *testing.T) {
		d := dml.Decimal{
			Valid:     true,
			Precision: math.MaxInt64 + 9876,
			Scale:     5,
			Negative:  true,
		}
		f := d.Float64()
		assert.Exactly(t, -9.223372036854788e+13, f)
	})
}
