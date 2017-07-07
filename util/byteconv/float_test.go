// Copyright (c) 2015 Taco de Wolff
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package byteconv

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseNullFloatSQL_ParseFloatSQL(t *testing.T) {

	runner := func(have string, want sql.NullFloat64, wantErr bool) func(*testing.T) {
		return func(t *testing.T) {
			b := sql.RawBytes(have)
			if have == "NULL" {
				b = nil
			}
			nf, err := ParseNullFloat64(b)
			if wantErr {
				assert.Error(t, err, "err: For number %q", have)
				return
			}
			require.NoError(t, err, t.Name())
			assert.Exactly(t, want, nf, t.Name())
		}
	}
	t.Run("NULL is 0 and invalid", runner("NULL", sql.NullFloat64{}, false))
	t.Run("empty is 0 and invalid", runner("", sql.NullFloat64{}, false))
	t.Run(" is 0 and invalid", runner("", sql.NullFloat64{}, true))
	t.Run("0 is valid", runner("0", sql.NullFloat64{Valid: true}, false))
	t.Run("1 valid", runner("1", sql.NullFloat64{Valid: true, Float64: 1}, false))
	t.Run("35.5456 valid", runner("35.5456", sql.NullFloat64{Valid: true, Float64: 35.5456}, false))
	t.Run("-35.5456 valid", runner("-35.5456", sql.NullFloat64{Valid: true, Float64: -35.5456}, false))
	t.Run("999 valid", runner("999", sql.NullFloat64{Valid: true, Float64: 999}, false))
	t.Run("10 is valid", runner("10", sql.NullFloat64{Valid: true, Float64: 10}, false))
	t.Run("01 is valid", runner("01", sql.NullFloat64{Valid: true, Float64: 1}, false))
}

func TestParseFloat(t *testing.T) {
	floatTests := []struct {
		f        string
		expected float64
	}{
		{"5", 5},
		{"5.1", 5.1},
		{"+5.1", 5.1},
		{"-5.1", -5.1},
		{"5.1e-2", 5.1e-2},
		{"5.1e+2", 5.1e+2},
		{"0.0e1", 0.0e1},
		{"18446744073709551620", 18446744073709551620.0},
		{"1e23", 1e23},
		// TODO: hard to test due to float imprecision
		// {"1.7976931348623e+308", 1.7976931348623e+308)
		// {"4.9406564584124e-308", 4.9406564584124e-308)
	}
	for i, tt := range floatTests {
		f, err := ParseFloat([]byte(tt.f))
		if err != nil {
			t.Fatalf("Index %d invalid length for %q with error %s", i, tt.f, err)
		}
		if f != tt.expected {
			t.Fatalf("Index %d\nHave %f\nWant %f", i, f, tt.expected)
		}
	}
}
