// The MIT License (MIT)
//
// Copyright (c) 2014 Alex Saskevich
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package validation

import "testing"

func TestAbs(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    float64
		expected float64
	}{
		{0, 0},
		{-1, 1},
		{10, 10},
		{3.14, 3.14},
		{-96, 96},
		{-10e-12, 10e-12},
	}
	for _, test := range tests {
		actual := Abs(test.param)
		if actual != test.expected {
			t.Errorf("Expected Abs(%v) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestSign(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    float64
		expected float64
	}{
		{0, 0},
		{-1, -1},
		{10, 1},
		{3.14, 1},
		{-96, -1},
		{-10e-12, -1},
	}
	for _, test := range tests {
		actual := Sign(test.param)
		if actual != test.expected {
			t.Errorf("Expected Sign(%v) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsNegative(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    float64
		expected bool
	}{
		{0, false},
		{-1, true},
		{10, false},
		{3.14, false},
		{-96, true},
		{-10e-12, true},
	}
	for _, test := range tests {
		actual := IsNegative(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsNegative(%v) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsNonNegative(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    float64
		expected bool
	}{
		{0, true},
		{-1, false},
		{10, true},
		{3.14, true},
		{-96, false},
		{-10e-12, false},
	}
	for _, test := range tests {
		actual := IsNonNegative(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsNonNegative(%v) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsPositive(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    float64
		expected bool
	}{
		{0, false},
		{-1, false},
		{10, true},
		{3.14, true},
		{-96, false},
		{-10e-12, false},
	}
	for _, test := range tests {
		actual := IsPositive(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsPositive(%v) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsNonPositive(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    float64
		expected bool
	}{
		{0, true},
		{-1, true},
		{10, false},
		{3.14, false},
		{-96, true},
		{-10e-12, true},
	}
	for _, test := range tests {
		actual := IsNonPositive(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsNonPositive(%v) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsWhole(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    float64
		expected bool
	}{
		{0, true},
		{-1, true},
		{10, true},
		{3.14, false},
		{-96, true},
		{-10e-12, false},
	}
	for _, test := range tests {
		actual := IsWhole(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsWhole(%v) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsNatural(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    float64
		expected bool
	}{
		{0, false},
		{-1, false},
		{10, true},
		{3.14, false},
		{96, true},
		{-10e-12, false},
	}
	for _, test := range tests {
		actual := IsNatural(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsNatural(%v) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestInRangeInt64(t *testing.T) {
	t.Parallel()

	var testAsInt64s = []struct {
		param    int64
		left     int64
		right    int64
		expected bool
	}{
		{0, 0, 0, true},
		{1, 0, 0, false},
		{-1, 0, 0, false},
		{0, -1, 1, true},
		{0, 0, 1, true},
		{0, -1, 0, true},
		{0, 0, -1, true},
		{0, 10, 5, false},
	}
	for _, test := range testAsInt64s {
		actual := InRangeInt64(test.param, test.left, test.right)
		if actual != test.expected {
			t.Errorf("Expected InRangeInt64(%v, %v, %v) to be %v, got %v using type int64", test.param, test.left, test.right, test.expected, actual)
		}
	}

	var testAsUInt64s = []struct {
		param    uint64
		left     uint64
		right    uint64
		expected bool
	}{
		{0, 0, 0, true},
		{1, 0, 0, false},
		{0, 0, 1, true},
		{0, 10, 5, false},
	}
	for _, test := range testAsUInt64s {
		actual := InRangeUint64(test.param, test.left, test.right)
		if actual != test.expected {
			t.Errorf("Expected InRangeInt64(%v, %v, %v) to be %v, got %v using type uint", test.param, test.left, test.right, test.expected, actual)
		}
	}

}

func TestInRangeFloat64(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    float64
		left     float64
		right    float64
		expected bool
	}{
		{0, 0, 0, true},
		{1, 0, 0, false},
		{-1, 0, 0, false},
		{0, -1, 1, true},
		{0, 0, 1, true},
		{0, -1, 0, true},
		{0, 0, -1, true},
		{0, 10, 5, false},
	}
	for _, test := range tests {
		actual := InRangeFloat64(test.param, test.left, test.right)
		if actual != test.expected {
			t.Errorf("Expected InRangeFloat64(%v, %v, %v) to be %v, got %v", test.param, test.left, test.right, test.expected, actual)
		}
	}
}
