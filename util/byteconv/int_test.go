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

func TestParseNullInt64SQL_ParseIntSQL(t *testing.T) {

	runner := func(have string, want sql.NullInt64, wantErr bool) func(*testing.T) {
		return func(t *testing.T) {
			b := sql.RawBytes(have)
			if have == "NULL" {
				b = nil
			}
			ni, ok, err := ParseInt(b)

			assert.Exactly(t, want.Valid, ok, t.Name())
			assert.Exactly(t, want.Int64, ni, t.Name())

			if wantErr {
				assert.Error(t, err, "For %q", have)
				return
			}
			require.NoError(t, err, t.Name())
		}
	}
	t.Run("NULL is 0 and invalid", runner("NULL", sql.NullInt64{}, false))
	t.Run("empty is 0 and invalid", runner("", sql.NullInt64{}, false))
	t.Run(" is 0 and invalid", runner("", sql.NullInt64{}, true))
	t.Run("0 is valid", runner("0", sql.NullInt64{Valid: true}, false))
	t.Run("1 valid", runner("1", sql.NullInt64{Valid: true, Int64: 1}, false))
	t.Run("35 valid", runner("35", sql.NullInt64{Valid: true, Int64: 35}, false))
	t.Run("-35 valid", runner("-35", sql.NullInt64{Valid: true, Int64: -35}, false))
	t.Run("35.5456 valid", runner("35.5456", sql.NullInt64{}, true))
	t.Run("10 is valid", runner("10", sql.NullInt64{Valid: true, Int64: 10}, false))
	t.Run("01 is valid", runner("01", sql.NullInt64{Valid: true, Int64: 1}, false))
}

func TestParseInt(t *testing.T) {
	intTests := []struct {
		i        string
		expected int64
	}{
		{"5", 5},
		{"99", 99},
		{"999", 999},
		{"-5", -5},
		{"+5", 5},
		{"9223372036854775807", 9223372036854775807},
		{"9223372036854775808", 0},
		{"-9223372036854775807", -9223372036854775807},
		{"-9223372036854775808", -9223372036854775808},
		{"-9223372036854775809", 0},
		{"18446744073709551620", 0},
		{"a", 0},
	}
	for _, tt := range intTests {
		i, _, _ := ParseInt([]byte(tt.i))
		if i != tt.expected {
			t.Fatalf("Have %d Want %d", i, tt.expected)
		}
	}
}

func TestLenInt(t *testing.T) {
	lenIntTests := []struct {
		number   int64
		expected int
	}{
		{0, 1},
		{1, 1},
		{10, 2},
		{99, 2},

		// coverage
		{100, 3},
		{1000, 4},
		{10000, 5},
		{100000, 6},
		{1000000, 7},
		{10000000, 8},
		{100000000, 9},
		{1000000000, 10},
		{10000000000, 11},
		{100000000000, 12},
		{1000000000000, 13},
		{10000000000000, 14},
		{100000000000000, 15},
		{1000000000000000, 16},
		{10000000000000000, 17},
		{100000000000000000, 18},
		{1000000000000000000, 19},
	}
	for _, tt := range lenIntTests {

		if li := LenInt(tt.number); li != tt.expected {
			t.Fatalf("Have %d Want %d", li, tt.expected)
		}

	}
}

////////////////////////////////////////////////////////////////

//var num []int64
//
//func TestMain(t *testing.T) {
//	for j := 0; j < 1000; j++ {
//		num = append(num, rand.Int63n(1000))
//	}
//}
//
//func BenchmarkLenIntLog(b *testing.B) {
//	n := 0
//	for i := 0; i < b.N; i++ {
//		for j := 0; j < 1000; j++ {
//			n += int(math.Log10(math.Abs(float64(num[j])))) + 1
//		}
//	}
//}
//
//func BenchmarkLenIntSwitch(b *testing.B) {
//	n := 0
//	for i := 0; i < b.N; i++ {
//		for j := 0; j < 1000; j++ {
//			n += LenInt(num[j])
//		}
//	}
//}
