// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

	"github.com/corestoreio/csfw/i18n"
	"github.com/corestoreio/csfw/storage/money"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

var testDefaultSymbols = i18n.Symbols{
	Decimal:                '.',
	Group:                  0,
	List:                   ';',
	PercentSign:            '%',
	CurrencySign:           '¤',
	PlusSign:               '+',
	MinusSign:              '-',
	Exponential:            'E',
	SuperscriptingExponent: '×',
	PerMille:               '‰',
	Infinity:               '∞',
	Nan:                    []byte(`NaN`),
}

var testFmtCur = i18n.NewCurrency(
	i18n.SetCurrencyFormat("¤ #,##0.00"),
	i18n.SetCurrencySymbols(testDefaultSymbols),
)
var testFmtNum = i18n.NewNumber(
	i18n.SetNumberFormat("###0.###"),
	i18n.SetNumberSymbols(testDefaultSymbols),
)

func TestString(t *testing.T) {

	tests := []struct {
		prec int
		have int64
		want string
	}{
		{0, 13, "$ 13.00"},
		{10, 13, "$ 1.30"},
		{100, 13, "$ 0.13"},
		{1000, 13, "$ 0.01"},
		{100, -13, "$ -0.13"},
		{0, -45628734653, "$ -45,628,734,653.00"},
		{10, -45628734653, "$ -4,562,873,465.30"},
		{100, -45628734653, "$ -456,287,346.53"},
		{1000, -45628734653, "$ -45,628,734.65"},
		{100, 256, "$ 2.56"},
		10: {1234, -45628734653, "$ -4,562,873.47"},
		{100, -45628734655, "$ -456,287,346.55"},
		{100, -45628734611, "$ -456,287,346.11"},
		{100, -45628734699, "$ -456,287,346.99"},
		14: {10000000, 45628734699, "$ 4,562.87"},
		15: {10000000, 45628734655, "$ 4,562.87"},
	}

	for i, test := range tests {
		c := money.New(
			money.WithPrecision(test.prec),
		).Set(test.have)
		c.FmtCur = testFmtCur
		c.FmtNum = testFmtNum

		have := c.String()
		if have != test.want {
			t.Errorf("\nWant: %s\nHave: %s\nIndex: %d\n", test.want, have, i)
		}
	}

}

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
		c := money.New(money.WithPrecision(test.prec)).Set(test.have)
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
		c := money.New(money.WithPrecision(test.prec)).Setf(test.havef)
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

func TestSwedishNumber(t *testing.T) {

	tests := []struct {
		prec int
		iv   money.Interval
		have int64
		want string
	}{
		{0, money.Interval005, 25689, "25689.000"},
		{100, money.Interval005, 25600, "256.000"},
		{10, money.Interval005, 25689, "2568.900"},
		{100, money.Interval005, 25689, "256.900"},
		{1000, money.Interval005, 25689, "25.700"},
		{100, money.Interval005, 25642, "256.400"},
		{100, money.Interval005, 25644, "256.450"},

		{0, money.Interval010, 25689, "25689.000"},
		{10, money.Interval010, 25689, "2568.900"},
		{100, money.Interval010, 25689, "256.900"},
		{1000, money.Interval010, 25689, "25.700"},
		{100, money.Interval010, 25642, "256.400"},
		{100, money.Interval010, 25644, "256.400"},
		{100, money.Interval010, 25639, "256.400"},
		{100, money.Interval010, 25635, "256.400"},
		{100, money.Interval010, 25634, "256.300"},
		{100, money.Interval010, 256345, "2563.500"},

		{0, money.Interval015, 25689, "25689.000"},
		{10, money.Interval015, 25689, "2568.900"},
		{10, money.Interval015, 25685, "2568.400"},
		{100, money.Interval015, 25689, "256.900"},
		{1000, money.Interval015, 25689, "25.700"},
		{100, money.Interval015, 25642, "256.400"},
		{100, money.Interval015, 25644, "256.400"},
		{100, money.Interval015, 25639, "256.400"},
		{100, money.Interval015, 25635, "256.300"},
		{100, money.Interval015, 25636, "256.400"},
		{100, money.Interval015, 25634, "256.300"},
		{100, money.Interval015, 256345, "2563.400"},

		{0, money.Interval025, 25689, "25689.000"},
		{10, money.Interval025, 25689, "2569.000"},
		{10, money.Interval025, 25685, "2568.500"},
		{100, money.Interval025, 25689, "257.000"},
		{1000, money.Interval025, 25689, "25.750"},
		{100, money.Interval025, 25642, "256.500"},
		{100, money.Interval025, 25644, "256.500"},
		{100, money.Interval025, 25639, "256.500"},
		{100, money.Interval025, 25624, "256.250"},
		{100, money.Interval025, 25625, "256.250"},
		{100, money.Interval025, 25634, "256.250"},
		{100, money.Interval025, 256345, "2563.500"},

		{0, money.Interval050, 25689, "25689.000"},
		{10, money.Interval050, 25689, "2569.000"},
		{10, money.Interval050, 25685, "2568.500"},
		{100, money.Interval050, 25689, "257.000"},
		{1000, money.Interval050, 25689, "25.500"},
		{100, money.Interval050, 25642, "256.500"},
		{100, money.Interval050, 25644, "256.500"},
		{100, money.Interval050, 25639, "256.500"},
		{100, money.Interval050, 25624, "256.000"},
		{100, money.Interval050, 25625, "256.500"},
		{100, money.Interval050, 25634, "256.500"},
		{100, money.Interval050, 256345, "2563.500"},

		{0, money.Interval100, 25689, "25689.000"},
		{10, money.Interval100, 25689, "2569.000"},
		{10, money.Interval100, 25685, "2569.000"},
		{10, money.Interval100, 25684, "2568.000"},
		{100, money.Interval100, 25689, "257.000"},
		{1000, money.Interval100, 25689, "26.000"},
		{100, money.Interval100, 25642, "256.000"},
		{100, money.Interval100, 25644, "256.000"},
		{100, money.Interval100, 25639, "256.000"},
		{100, money.Interval100, 25624, "256.000"},
		{100, money.Interval100, 25625, "256.000"},
		{100, money.Interval100, 25634, "256.000"},
		{100, money.Interval100, 256345, "2563.000"},
	}
	for _, test := range tests {
		c := money.New(
			money.WithPrecision(test.prec),
		).Set(test.have)
		c.FmtCur = testFmtCur
		c.FmtNum = testFmtNum

		haveB, err := c.Swedish(money.WithSwedish(test.iv)).Number()
		assert.NoError(t, err, "%v", test)

		if haveB.String() != test.want {
			t.Errorf("\nWant: %s\nHave: %s\nIndex: %v\n", test.want, haveB.String(), test)
		}
	}
}

func TestMoney_Add(t *testing.T) {

	tests := []struct {
		have1 int64
		have2 int64
		want  int64
	}{
		{13, 13, 26},
		{-13, -13, -26},
		{-45628734653, -45628734653, -91257469306},
	}

	for i, test := range tests {
		c := money.New().Set(test.have1)
		c = c.Add(money.New().Set(test.have2))
		have := c.Raw()
		if have != test.want {
			t.Errorf("\nWant: %d\nHave: %d\nIndex: %d\n", test.want, have, i)
		}
	}
}

func TestMoney_Add_Overflow(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.True(t, errors.IsNotValid(err), "Error %+v", err)
			} else {
				t.Fatal("Expecting an error in the panic")
			}
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	c := money.New().Set(math.MaxInt64)
	c.Add(money.New().Set(2))
}

func TestMoney_Sub(t *testing.T) {

	tests := []struct {
		have1 int64
		have2 int64
		want  int64
	}{
		{13, 13, 0},
		{-13, -13, 0},
		{-13, 13, -26},
		{-45628734653, -45628734653, 0},
	}

	for i, test := range tests {
		c := money.New().Set(test.have1)
		c = c.Sub(money.New().Set(test.have2))
		have := c.Raw()
		if have != test.want {
			t.Errorf("\nWant: %d\nHave: %d\nIndex: %d\n", test.want, have, i)
		}
	}
}

func TestMoney_Sub_Overflow(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.True(t, errors.IsNotValid(err), "Error %+v", err)
			} else {
				t.Fatal("Expecting an error in the panic")
			}
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	c := money.New().Set(-math.MaxInt64)
	c.Sub(money.New().Set(2))
}

func TestMulNumber(t *testing.T) {

	tests := []struct {
		prec  int
		have1 int64
		have2 int64
		want  string
	}{
		{0, 1300, 1300, "1690000.000"},
		{100, 1300, 1300, "169.000"},
		{1000, 18100, 18100, "327.610"},
		{100, 1319, 1488, "196.270"},
		{1000, 1319, 1488, "1.963"},
		{100, 13, -13, "-0.020"},
		{100, 1300, -1300, "-169.000"},
		{1000, 1300, -1300, "-1.690"},
		{100, 13, 13, "0.020"},
		{100, 45628734653, -45628734653, "250065429529630.200"}, // overflow of int64 ?
		{100, 45628734653, -456287346, "-237307016244604.920"},
		{100, math.MaxInt64, 2, "0.000"},
	}

	for _, test := range tests {
		c := money.New(
			money.WithPrecision(test.prec),
		).Set(test.have1)
		c.FmtCur = testFmtCur
		c.FmtNum = testFmtNum

		c = c.Mul(money.New(money.WithPrecision(test.prec)).Set(test.have2))

		haveB, err := c.Number()
		assert.NoError(t, err)

		if haveB.String() != test.want {
			t.Errorf("\nWant: %s\nHave: %s\nSign %d\nIndex: %v\n", test.want, haveB.String(), c.Sign(), test)
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
		{100, 1300, 1300.13, "16,901.690"},
		{1000, 18100, 18100.18, "327,613.258"},
		{0, 18100, 18100.18, "327,613,258.000"},
		{0, 18103, 18.307, "331,412.000"}, // rounds up
		{100, 1319, 1488.88, "19,638.330"},
		{1000, 1319, 1488.88, "1,963.833"},
		{100, 13, -13.13, "-1.710"},
		{100, 1300, -1300.01, "-16,900.130"},
		{1000, 1300, -1300.01, "-1,690.013"},
		{100, 13, 13.0, "1.690"},
		{100, 45628734653, -45628734653.0, "-47,780,798,383.280"},
		{100, math.MaxInt64, 2.01, "92,233,720,368.530"},
	}

	for _, test := range tests {
		c := money.New(money.WithPrecision(test.prec)).Set(test.have1)
		c.FmtCur = i18n.DefaultCurrency
		c.FmtNum = i18n.DefaultNumber
		c = c.Mulf(test.have2)
		haveB, err := c.Number()
		assert.NoError(t, err)

		if haveB.String() != test.want {
			t.Errorf("\nWant: %s\nHave: %s\nSign %d\nIndex: %v\n", test.want, haveB.String(), c.Sign(), test)
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
		nob, err := c.Number()
		assert.NoError(t, err)

		if have != test.want {
			t.Errorf("\nWant: %d\nHave: %d / %s\nIndex: %v\n", test.want, have, nob.String(), test)
		}
	}
}
