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
	"fmt"
	"testing"

	"github.com/corestoreio/pkg/storage/money"
)

var _ fmt.Stringer = (*money.Money)(nil)

//var _ fmt.GoStringer = (*money.Money)(nil)

func TestMoney_String(t *testing.T) {

	tests := []struct {
		prec int
		have int64
		want string
	}{
		{0, 13, "13.00"},
		{1, 13, "1.3"},
		{1, 26, "2.6"},
		{2, 13, "0.13"},
		{2, 1300, "13.00"},
		{3, 13, "0.013"},
		{2, -13, "-0.13"},
		{0, -45628734653, "-45628734653.00"},
		{1, -45628734653, "-4562873465.3"},
		{2, -456287346530, "-4562873465.30"},
		{2, 256, "2.56"},
		{7, 45628734655, "4562.8734655"},
	}

	for i, test := range tests {
		c := money.New(test.have, test.prec)

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
		c := money.New(test.have, 0)
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
		{1, 13, 1, 1.30000},
		{2, 13, 0, 0.130000},
		{3, 13, 0, 0.013000},
		{2, -13, 0, -0.130000},
		{0, -45628734653, -45628734653, -45628734653.000000},
		{1, -45628734653, -4562873465, -4562873465.300000},
		{2, -45628734653, -456287346, -456287346.530000},
		{3, -45628734653, -45628734, -45628734.653000},
		{2, 256, 2, 2.56},
		{3, -45628734653, -45628734, -45628734.6530000},
		{2, -45628734655, -456287346, -456287346.550000},
		{2, -45628734611, -456287346, -456287346.110000},
		{2, -45628734699, -456287346, -456287346.990000},
		{7, 45628734699, 4562, 4562.87346989999969082419},
		{7, 45628734655, 4562, 4562.87346549999983835733},
	}

	for i, test := range tests {
		c := money.New(test.have, test.prec)
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
		{1, 13, 1.30000},
		{2, 13, 0.130000},
		{3, 13, 0.013000},
		{2, -13, -0.130000},
		{0, -45628734653, -45628734653.000000},
		{1, -45628734653, -4562873465.300000},
		{2, -45628734653, -456287346.530000},
		{3, -45628734653, -45628734.653000},
		{2, 256, 2.56},
		{3, -45628734653, -45628734.6530000},
		{2, -45628734655, -456287346.550000},
		{2, -45628734611, -456287346.110000},
		{2, -45628734699, -456287346.990000},
		{7, 45628734699, 4562.87346989999969082419},
		{7, 45628734655, 4562.87346549999983835733},
	}

	for i, test := range tests {
		c := money.New(0, test.prec).Setf(test.havef)
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
		c := money.New(test.have, 0)
		have := c.Sign()
		if have != test.want {
			t.Errorf("\nWant: %d\nHave: %d\nIndex: %d\n", test.want, have, i)
		}
	}
}

func TestSwedishNumber(t *testing.T) {

	tests := []struct {
		prec int
		iv   uint8
		have int64
		want string
	}{
		{0, 5, 25689, "25689.00"},
		{2, 5, 25600, "256.00"},
		{1, 5, 25689, "2568.9"},
		{2, 5, 25689, "256.90"},
		{3, 5, 25689, "25.700"},
		{2, 5, 25642, "256.40"},
		{2, 5, 25644, "256.45"},

		{0, 10, 25689, "25689.00"},
		{1, 10, 25689, "2568.9"},
		{2, 10, 25689, "256.90"},
		{3, 10, 25689, "25.700"},
		{2, 10, 25642, "256.40"},
		{2, 10, 25644, "256.40"},
		{2, 10, 25639, "256.40"},
		{2, 10, 25635, "256.40"},
		{2, 10, 25634, "256.30"},
		{2, 10, 256345, "2563.50"},

		{0, 15, 25689, "25689.00"},
		{1, 15, 25689, "2568.9"},
		{1, 15, 25685, "2568.4"},
		{2, 15, 25689, "256.90"},
		{3, 15, 25689, "25.700"},
		{2, 15, 25642, "256.40"},
		{2, 15, 25644, "256.40"},
		{2, 15, 25639, "256.40"},
		{2, 15, 25635, "256.30"},
		{2, 15, 25636, "256.40"},
		{2, 15, 25634, "256.30"},
		{2, 15, 256345, "2563.40"},

		{0, 25, 25689, "25689.00"},
		{1, 25, 25689, "2569.0"},
		{1, 25, 25685, "2568.5"},
		{2, 25, 25689, "257.00"},
		{3, 25, 25689, "25.750"},
		{2, 25, 25642, "256.50"},
		{2, 25, 25644, "256.50"},
		{2, 25, 25639, "256.50"},
		{2, 25, 25624, "256.25"},
		{2, 25, 25625, "256.25"},
		{2, 25, 25634, "256.25"},
		{2, 25, 256345, "2563.50"},

		{0, 50, 25689, "25689.00"},
		{1, 50, 25689, "2569.0"},
		{1, 50, 25685, "2568.5"},
		{2, 50, 25689, "257.00"},
		{3, 50, 25689, "25.500"},
		{2, 50, 25642, "256.50"},
		{2, 50, 25644, "256.50"},
		{2, 50, 25639, "256.50"},
		{2, 50, 25624, "256.00"},
		{2, 50, 25625, "256.50"},
		{2, 50, 25634, "256.50"},
		{2, 50, 256345, "2563.50"},

		{0, 100, 25689, "25689.00"},
		{1, 100, 25689, "2569.0"},
		{1, 100, 25685, "2569.0"},
		{1, 100, 25684, "2568.0"},
		{2, 100, 25689, "257.00"},
		{3, 100, 25689, "26.000"},
		{2, 100, 25642, "256.00"},
		{2, 100, 25644, "256.00"},
		{2, 100, 25639, "256.00"},
		{2, 100, 25624, "256.00"},
		{2, 100, 25625, "256.00"},
		{2, 100, 25634, "256.00"},
		{2, 100, 256345, "2563.00"},
	}
	for _, test := range tests {
		c := money.New(test.have, test.prec)

		haveB := c.RoundCash(test.iv).String()

		if haveB != test.want {
			t.Errorf("\nWant: %s\nHave: %s\nIndex: %v\n", test.want, haveB, test)
		}
	}
}

func TestMoney_Rescale(t *testing.T) {

	tests := []struct {
		have1   int64
		prec1   int
		newPrec int
		want    string
	}{
	//{13, 0, 1, "13.0"},
	//{256, 1, 1, "25.6"},
	//{256, 2, 1, "2.6"},
	//{2545, 3, 1, "2.6"}, // 2.545 => 2.6
	//{2535, 3, 1, "2.5"},
	}

	for i, test := range tests {
		m := money.New(test.have1, test.prec1)
		haveM := m.Rescale(test.newPrec)

		if have := haveM.String(); have != test.want {
			t.Errorf("\nWant: %d\nHave: %d\nIndex: %d\n", test.want, have, i)
		}
	}
}

func TestMoney_Add(t *testing.T) {

	tests := []struct {
		have1 int64
		prec1 int
		have2 int64
		prec2 int
		want  string
	}{
		{13, 0, 13, 0, "26.00"},
		{13, 1, 13, 2, "1.4"},
		{14, 1, 13, 2, "1.5"},
		{2545, 3, 5, 4, "2.546"},  // 2.545+0.0005 = 2.546
		{2545, 3, -5, 4, "2.545"}, // 2.545-0.0005 = 2.546
		{-13, 0, -13, 0, "-26.00"},
		{-45628734653, 7, 45628734653, 7, "0.0000000"},
		{-45628734653, 7, 3, 7, "-4562.8734650"},
	}

	for i, test := range tests {
		haveM := money.New(test.have1, test.prec1)
		haveM = haveM.Add(money.New(test.have2, test.prec2))
		if have := haveM.String(); have != test.want {
			t.Errorf("\nWant: %q\nHave: %q\nIndex: %d\n", test.want, have, i)
		}
	}
}

//func TestMoney_Add_Overflow(t *testing.T) {
//	defer func() {
//		if r := recover(); r != nil {
//			if err, ok := r.(error); ok {
//				assert.True(t, errors.IsNotValid(err), "Error %+v", err)
//			} else {
//				t.Fatal("Expecting an error in the panic")
//			}
//		} else {
//			t.Fatal("Expecting a panic")
//		}
//	}()
//	c := money.New().Set(math.MaxInt64)
//	c.Add(money.New().Set(2))
//}
//
//func TestMoney_Sub(t *testing.T) {
//
//	tests := []struct {
//		have1 int64
//		have2 int64
//		want  int64
//	}{
//		{13, 13, 0},
//		{-13, -13, 0},
//		{-13, 13, -26},
//		{-45628734653, -45628734653, 0},
//	}
//
//	for i, test := range tests {
//		c := money.New().Set(test.have1)
//		c = c.Sub(money.New().Set(test.have2))
//		have := c.Raw()
//		if have != test.want {
//			t.Errorf("\nWant: %d\nHave: %d\nIndex: %d\n", test.want, have, i)
//		}
//	}
//}
//
//func TestMoney_Sub_Overflow(t *testing.T) {
//	defer func() {
//		if r := recover(); r != nil {
//			if err, ok := r.(error); ok {
//				assert.True(t, errors.IsNotValid(err), "Error %+v", err)
//			} else {
//				t.Fatal("Expecting an error in the panic")
//			}
//		} else {
//			t.Fatal("Expecting a panic")
//		}
//	}()
//	c := money.New().Set(-math.MaxInt64)
//	c.Sub(money.New().Set(2))
//}
//
//func TestMulNumber(t *testing.T) {
//
//	tests := []struct {
//		prec  int
//		have1 int64
//		have2 int64
//		want  string
//	}{
//		{0, 1300, 1300, "1690000.000"},
//		{100, 1300, 1300, "169.000"},
//		{1000, 18100, 18100, "327.610"},
//		{100, 1319, 1488, "196.270"},
//		{1000, 1319, 1488, "1.963"},
//		{100, 13, -13, "-0.020"},
//		{100, 1300, -1300, "-169.000"},
//		{1000, 1300, -1300, "-1.690"},
//		{100, 13, 13, "0.020"},
//		{100, 45628734653, -45628734653, "250065429529630.200"}, // overflow of int64 ?
//		{100, 45628734653, -456287346, "-237307016244604.920"},
//		{100, math.MaxInt64, 2, "0.000"},
//	}
//
//	for _, test := range tests {
//		c := money.New(
//			money.WithPrecision(test.prec),
//		).Set(test.have1)
//
//		c = c.Mul(money.New(money.WithPrecision(test.prec)).Set(test.have2))
//
//		haveB, err := c.Number()
//		assert.NoError(t, err)
//
//		if haveB.String() != test.want {
//			t.Errorf("\nWant: %s\nHave: %s\nSign %d\nIndex: %v\n", test.want, haveB.String(), c.Sign(), test)
//		}
//	}
//}
//
//func TestMulf(t *testing.T) {
//
//	tests := []struct {
//		prec  int
//		have1 int64
//		have2 float64
//		want  string
//	}{
//		{100, 1300, 1300.13, "16,901.690"},
//		{1000, 18100, 18100.18, "327,613.258"},
//		{0, 18100, 18100.18, "327,613,258.000"},
//		{0, 18103, 18.307, "331,412.000"}, // rounds up
//		{100, 1319, 1488.88, "19,638.330"},
//		{1000, 1319, 1488.88, "1,963.833"},
//		{100, 13, -13.13, "-1.710"},
//		{100, 1300, -1300.01, "-16,900.130"},
//		{1000, 1300, -1300.01, "-1,690.013"},
//		{100, 13, 13.0, "1.690"},
//		{100, 45628734653, -45628734653.0, "-47,780,798,383.280"},
//		{100, math.MaxInt64, 2.01, "92,233,720,368.530"},
//	}
//
//	for _, test := range tests {
//		c := money.New(money.WithPrecision(test.prec)).Set(test.have1)
//		c = c.Mulf(test.have2)
//		haveB, err := c.Number()
//		assert.NoError(t, err)
//
//		if haveB.String() != test.want {
//			t.Errorf("\nWant: %s\nHave: %s\nSign %d\nIndex: %v\n", test.want, haveB.String(), c.Sign(), test)
//		}
//	}
//}
//
//func TestDiv(t *testing.T) {
//
//	tests := []struct {
//		have1 int64
//		have2 int64
//		want  int64
//	}{
//		{1300, 1300, 10000},
//		{13, -13, -10000},
//		{9000, -3000, -30000},
//		{13, 13, 10000},
//		{471100, 81500, 57804},
//		{45628734653, -45628734653, -10000},
//		{math.MaxInt64, 2, -9223372036854775807},
//	}
//
//	for _, test := range tests {
//		c := money.New().Set(test.have1)
//		c = c.Div(money.New().Set(test.have2))
//		have := c.Raw()
//		nob, err := c.Number()
//		assert.NoError(t, err)
//
//		if have != test.want {
//			t.Errorf("\nWant: %d\nHave: %d / %s\nIndex: %v\n", test.want, have, nob.String(), test)
//		}
//	}
//}
