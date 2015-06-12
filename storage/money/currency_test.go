// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package money_test

import (
	"testing"

	"math"

	"github.com/corestoreio/csfw/storage/money"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		have int64
		want int64
	}{
		{13, 13},
		{-13, 13},
		{-45628734653, 45628734653},
	}

	for i, test := range tests {
		c := money.New().Set(test.have)
		have := c.Abs().Raw()
		if have != test.want {
			t.Errorf("\nWant: %d\nHave: %d\nIndex: %d\n", test.want, have, i)
		}
	}
}

func TestPrecisionAndGet(t *testing.T) {
	tests := []struct {
		prec  int
		have  int64
		wanti int64
		wantf float64
	}{
		{0, 13, 13, 13.0000},
		{10, 13, 1, 1.30000},
		{100, 13, 0, 0.130000},
		{1000, 13, 0, 0.013000},
		{100, -13, 0, -0.130000},
		{0, -45628734653, -45628734653, -45628734653.000000},
		{10, -45628734653, -4562873465, -4562873465.300000},
		{100, -45628734653, -456287346, -456287346.530000},
		{1000, -45628734653, -45628734, -45628734.653000},
		{100, 256, 2, 2.56},
		10: {1234, -45628734653, -4562873, -4562873.46530000}, // fallback to prec 10000
		{100, -45628734655, -456287346, -456287346.550000},
		{100, -45628734611, -456287346, -456287346.110000},
		{100, -45628734699, -456287346, -456287346.990000},
		14: {10000000, 45628734699, 4562, 4562.87346989999969082419},
		15: {10000000, 45628734655, 4562, 4562.87346549999983835733},
	}

	for i, test := range tests {
		c := money.New(money.Precision(test.prec)).Set(test.have)
		haveCI := c.Geti()
		haveCF := c.Getf()

		if haveCI != test.wanti {
			t.Errorf("\nWantI: %d\nHaveI: %d\nIndex: %d\n", test.wanti, haveCI, i)
		}
		if haveCF != test.wantf {
			t.Errorf("\nWantF: %f\nHaveF: %.20f\nIndex: %d\n", test.wantf, haveCF, i)
		}
	}
}

func TestSetf(t *testing.T) {
	tests := []struct {
		prec  int
		want  int64
		havef float64
	}{
		{0, 13, 13.0000},
		{10, 13, 1.30000},
		{100, 13, 0.130000},
		{1000, 13, 0.013000},
		{100, -13, -0.130000},
		{0, -45628734653, -45628734653.000000},
		{10, -45628734653, -4562873465.300000},
		{100, -45628734653, -456287346.530000},
		{1000, -45628734653, -45628734.653000},
		{100, 256, 2.56},
		10: {1234, -45628734653, -4562873.46530000}, // fallback to prec 10000
		{100, -45628734655, -456287346.550000},
		{100, -45628734611, -456287346.110000},
		{100, -45628734699, -456287346.990000},
		14: {10000000, 45628734699, 4562.87346989999969082419},
		15: {10000000, 45628734655, 4562.87346549999983835733},
	}

	for i, test := range tests {
		c := money.New(money.Precision(test.prec)).Setf(test.havef)
		haveR := c.Raw()
		if haveR != test.want {
			t.Errorf("\nWantI: %d\nHaveI: %d\nIndex: %d\n", test.want, haveR, i)
		}
	}
}

func TestSign(t *testing.T) {
	tests := []struct {
		have int64
		want int
	}{{13, 1}, {-13, -1}, {-45628734653, -1}, {45628734699, 1}}
	for i, test := range tests {
		c := money.New().Set(test.have)
		have := c.Sign()
		if have != test.want {
			t.Errorf("\nWant: %d\nHave: %d\nIndex: %d\n", test.want, have, i)
		}
	}
}

func TestSwedish(t *testing.T) {
	tests := []struct {
		prec int
		iv   money.Interval
		have int64
		want string
	}{
		{0, money.Interval005, 25689, "25689.00"},
		{10, money.Interval005, 25689, "2568.09"},
		{100, money.Interval005, 25689, "256.90"},
		{1000, money.Interval005, 25689, "25.700"},
		{100, money.Interval005, 25642, "256.40"},
		{100, money.Interval005, 25644, "256.45"},

		{0, money.Interval010, 25689, "25689.00"},
		{10, money.Interval010, 25689, "2568.09"},
		{100, money.Interval010, 25689, "256.90"},
		{1000, money.Interval010, 25689, "25.700"},
		{100, money.Interval010, 25642, "256.40"},
		{100, money.Interval010, 25644, "256.40"},
		{100, money.Interval010, 25639, "256.40"},
		{100, money.Interval010, 25635, "256.40"},
		{100, money.Interval010, 25634, "256.30"},
		{100, money.Interval010, 256345, "2563.50"},

		{0, money.Interval015, 25689, "25689.00"},
		{10, money.Interval015, 25689, "2568.09"},
		{10, money.Interval015, 25685, "2568.04"},
		{100, money.Interval015, 25689, "256.90"},
		{1000, money.Interval015, 25689, "25.700"},
		{100, money.Interval015, 25642, "256.40"},
		{100, money.Interval015, 25644, "256.40"},
		{100, money.Interval015, 25639, "256.40"},
		{100, money.Interval015, 25635, "256.30"},
		{100, money.Interval015, 25636, "256.40"},
		{100, money.Interval015, 25634, "256.30"},
		{100, money.Interval015, 256345, "2563.40"},

		{0, money.Interval025, 25689, "25689.00"},
		{10, money.Interval025, 25689, "2569.00"},
		{10, money.Interval025, 25685, "2568.05"},
		{100, money.Interval025, 25689, "257.00"},
		{1000, money.Interval025, 25689, "25.750"},
		{100, money.Interval025, 25642, "256.50"},
		{100, money.Interval025, 25644, "256.50"},
		{100, money.Interval025, 25639, "256.50"},
		{100, money.Interval025, 25624, "256.25"},
		{100, money.Interval025, 25625, "256.25"},
		{100, money.Interval025, 25634, "256.25"},
		{100, money.Interval025, 256345, "2563.50"},

		{0, money.Interval050, 25689, "25689.00"},
		{10, money.Interval050, 25689, "2569.00"},
		{10, money.Interval050, 25685, "2568.05"},
		{100, money.Interval050, 25689, "257.00"},
		{1000, money.Interval050, 25689, "25.500"},
		{100, money.Interval050, 25642, "256.50"},
		{100, money.Interval050, 25644, "256.50"},
		{100, money.Interval050, 25639, "256.50"},
		{100, money.Interval050, 25624, "256.00"},
		{100, money.Interval050, 25625, "256.50"},
		{100, money.Interval050, 25634, "256.50"},
		{100, money.Interval050, 256345, "2563.50"},

		{0, money.Interval100, 25689, "25689.00"},
		{10, money.Interval100, 25689, "2569.00"},
		{10, money.Interval100, 25685, "2569.00"},
		{10, money.Interval100, 25684, "2568.00"},
		{100, money.Interval100, 25689, "257.00"},
		{1000, money.Interval100, 25689, "26.00"},
		{100, money.Interval100, 25642, "256.00"},
		{100, money.Interval100, 25644, "256.00"},
		{100, money.Interval100, 25639, "256.00"},
		{100, money.Interval100, 25624, "256.00"},
		{100, money.Interval100, 25625, "256.00"},
		{100, money.Interval100, 25634, "256.00"},
		{100, money.Interval100, 256345, "2563.00"},
	}
	for _, test := range tests {
		c := money.New(money.Precision(test.prec)).Set(test.have)
		have := c.Swedish(money.Swedish(test.iv)).Number()
		if have != test.want {
			t.Errorf("\nWant: %s\nHave: %s\nIndex: %v\n", test.want, have, test)
		}
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		have1 int64
		have2 int64
		want  int64
	}{
		{13, 13, 26},
		{-13, -13, -26},
		{-45628734653, -45628734653, -91257469306},
		{math.MaxInt64, 2, 0},
	}

	for _, test := range tests {
		c := money.New().Set(test.have1)
		c = c.Add(money.New().Set(test.have2))
		have := c.Raw()
		if have != test.want {
			t.Errorf("\nWant: %d\nHave: %d\nIndex: %v\n", test.want, have, test)
		}
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		have1 int64
		have2 int64
		want  int64
	}{
		{13, 13, 0},
		{-13, -13, 0},
		{-13, 13, -26},
		{-45628734653, -45628734653, 0},
		{-math.MaxInt64, 2, 0},
	}

	for _, test := range tests {
		c := money.New().Set(test.have1)
		c = c.Sub(money.New().Set(test.have2))
		have := c.Raw()
		if have != test.want {
			t.Errorf("\nWant: %d\nHave: %d\nIndex: %v\n", test.want, have, test)
		}
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		prec  int
		have1 int64
		have2 int64
		want  string
	}{
		{100, 1300, 1300, "169.00"},
		{1000, 18100, 18100, "327.610"},
		{100, 1319, 1488, "196.26"},
		{1000, 1319, 1488, "1.962"},
		{100, 13, -13, "-0.01"},
		{100, 1300, -1300, "-169.00"},
		{1000, 1300, -1300, "-1.690"},
		{100, 13, 13, "0.01"},
		{100, 45628734653, -45628734653, "250065429529630.21"}, // overflow ?
		{100, 45628734653, -456287346, "-237307016244604.93"},
		{100, math.MaxInt64, 2, "0.00"},
	}

	for _, test := range tests {
		c := money.New(money.Precision(test.prec)).Set(test.have1)
		c = c.Mul(money.New(money.Precision(test.prec)).Set(test.have2))
		have := c.Number()
		if have != test.want {
			t.Errorf("\nWant: %s\nHave: %s\nSign %d\nIndex: %v\n", test.want, have, c.Sign(), test)
		}
	}
}

func TestMulf(t *testing.T) {
	tests := []struct {
		prec  int
		have1 int64
		have2 float64
		want  string
	}{
		{100, 1300, 1300.13, "16901.69"},
		{1000, 18100, 18100.18, "327613.258"},
		{100, 1319, 1488.88, "19638.33"},
		{1000, 1319, 1488.88, "1963.833"},
		{100, 13, -13.13, "-1.71"},
		{100, 1300, -1300.01, "-16900.13"},
		{1000, 1300, -1300.01, "-1690.13"},
		{100, 13, 13.0, "1.69"},
		{100, 45628734653, -45628734653.0, "-47780798383.28"},
		{100, math.MaxInt64, 2.01, "92233720368.53"},
	}

	for _, test := range tests {
		c := money.New(money.Precision(test.prec)).Set(test.have1)
		c = c.Mulf(test.have2)
		have := c.Number()
		if have != test.want {
			t.Errorf("\nWant: %s\nHave: %s\nSign %d\nIndex: %v\n", test.want, have, c.Sign(), test)
		}
	}
}

func TestDiv(t *testing.T) {
	tests := []struct {
		have1 int64
		have2 int64
		want  int64
	}{
		{1300, 1300, 10000},
		{13, -13, -10000},
		{9000, -3000, -30000},
		{13, 13, 10000},
		{471100, 81500, 57804},
		{45628734653, -45628734653, -10000},
		{math.MaxInt64, 2, -9223372036854775807},
	}

	for _, test := range tests {
		c := money.New().Set(test.have1)
		c = c.Div(money.New().Set(test.have2))
		have := c.Raw()
		if have != test.want {
			t.Errorf("\nWant: %d\nHave: %d / %s\nIndex: %v\n", test.want, have, c.Number(), test)
		}
	}
}
