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
	"testing"
)

func TestParseFloatPtr(t *testing.T) {
	b := []byte(`5.1`)
	f, n := ParseFloatPtr(&b)
	if n != 3 {
		t.Errorf("Have %d Want %d", n, 3)
	}
	if f != 5.1 {
		t.Errorf("Have %f Want %f", f, 5.1)
	}

	f, n = ParseFloatPtr(nil)
	if n != 0 {
		t.Errorf("Have %d Want %d", n, 0)
	}
	if f != 0 {
		t.Errorf("Have %f Want %f", f, 0)
	}
}

func TestParseFloat(t *testing.T) {
	floatTests := []struct {
		f        string
		expected float64
	}{
		{"5", 5},
		{"5.1", 5.1},
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
		f, n := ParseFloat([]byte(tt.f))
		if n != len(tt.f) {
			t.Fatalf("Index %d invalid length for %q", i, tt.f)
		}
		if f != tt.expected {
			t.Fatalf("Index %d\nHave %f\nWant %f", i, f, tt.expected)
		}
	}
}
