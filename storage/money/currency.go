// Copyright 2015 CoreStore Authors
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

/*
Package money usese a fixed-length guard for precision arithmetic: the
int64 variable Guard (and its float64 and int-related variables Guardf
and Guardi.

Rounding is done on float64 to int64 by	the Rnd() function truncating
at values less than (.5 + (1 / Guardf))	or greater than -(.5 + (1 / Guardf))
in the case of negative numbers. The Guard adds four decimal places
of protection to rounding.
DP is the decimal precision, which can be changed in the DecimalPrecision()
function.  DP hold the places after the decimalplace in teh active money struct field M.

Initial copyright: Copyright (c) 2011 Jad Dittmar
https://github.com/Confunctionist/finance

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/
package money

import (
	"errors"
	"fmt"
	"math"

	"database/sql/driver"

	"encoding/json"

	"github.com/corestoreio/csfw/utils/log"
)

var (
	ErrOverflow         = errors.New("Integer Overflow")
	Guardi      int     = 100
	Guard       int64   = int64(Guardi)
	Guardf      float64 = float64(Guardi)
	DP          int64   = 100         // for default of 2 decimal places => 10^2 (can be reset)
	DPf         float64 = float64(DP) // for default of 2 decimal places => 10^2 (can be reset)
	Round               = .5
	//	Round  = .5 + (1 / Guardf)
	Roundn = Round * -1
)

// Currency represents a money aka currency type to avoid rounding errors with floats
type Currency struct {
	// M money
	M int64
	// L locale to allow language specific output formats
	L string
}

// Abs Returns the absolute value of Currency
func (m *Currency) Abs() *Currency {
	if m.M < 0 {
		m.Neg()
	}
	return m
}

// Add Adds two Currency types. Returns nil on integer overflow
func (m *Currency) Add(n *Currency) *Currency {
	r := m.M + n.M
	if (r^m.M)&(r^n.M) < 0 {
		if log.IsTrace() {
			log.Trace("Currency=Add", "err", ErrOverflow, "m", m, "n", n)
		}
		log.Error("Currency=Add", "err", ErrOverflow, "m", m, "n", n)
		return nil
	}
	m.M = r
	return m
}

// Get gets the float64 value of money (see Raw() for int64)
func (m *Currency) Get() float64 {
	return float64(m.M) / DPf
}

// Gett gets value of money truncating after DP (see Raw() for no truncation)
func (m *Currency) Gett() int64 {
	return m.M / DP
}

// Raw returns in int64 the value of Currency (also see Gett(), See Get() for float64)
func (m *Currency) Raw() int64 {
	return m.M
}

// Set sets the Currency field M
func (m *Currency) Set(x int64) *Currency {
	m.M = x
	return m
}

// Setf sets a float64 into a Currency type for precision calculations
func (m *Currency) Setf(f float64) *Currency {
	fDPf := f * DPf
	r := int64(f * DPf)
	return m.Set(Rnd(r, fDPf-float64(r)))
}

// Sign returns the Sign of Currency 1 if positive, -1 if negative
func (m *Currency) Sign() int {
	if m.M < 0 {
		return -1
	}
	return 1
}

// String for money type representation in basic monetary unit (DOLLARS CENTS)
// @todo consider locale
func (m *Currency) String() string {
	return fmt.Sprintf("%d.%02d", m.Value()/DP, m.Abs().Value()%DP)
}

// Sub subtracts one Currency type from another. Returns nil on integer overflow
func (m *Currency) Sub(n *Currency) *Currency {
	r := m.M - n.M
	if (r^m.M)&^(r^n.M) < 0 {
		if log.IsTrace() {
			log.Trace("Currency=Sub", "err", ErrOverflow, "m", m, "n", n)
		}
		log.Error("Currency=Sub", "err", ErrOverflow, "m", m, "n", n)
		return nil
	}
	m.M = r
	return m
}

// Mul Multiplies two Currency types
func (m *Currency) Mul(n *Currency) *Currency {
	return m.Set(m.M * n.M / DP)
}

// Mulf Multiplies a Currency with a float to return a money-stored type
func (m *Currency) Mulf(f float64) *Currency {
	i := m.M * int64(f*Guardf*DPf)
	r := i / Guard / DP
	return m.Set(Rnd(r, float64(i)/Guardf/DPf-float64(r)))
}

// Neg Returns the negative value of Currency
func (m *Currency) Neg() *Currency {
	if m.M != 0 {
		m.M *= -1
	}
	return m
}

// Pow is the power of Currency
func (m *Currency) Pow(r float64) *Currency {
	return m.Setf(math.Pow(m.Get(), r))
}

// RND rounds int64 remainder rounded half towards plus infinity
// trunc = the remainder of the float64 calc
// r     = the result of the int64 cal
func Rnd(r int64, trunc float64) int64 {

	//fmt.Printf("RND 1 r = % v, trunc = %v Round = %v\n", r, trunc, Round)
	if trunc > 0 {
		if trunc >= Round {
			r++
		}
	} else {
		if trunc < Roundn {
			r--
		}
	}
	//fmt.Printf("RND 2 r = % v, trunc = %v Round = %v\n", r, trunc, Round)
	return r
}

var (
	_ json.Unmarshaler      = (*Currency)(nil)
	_ json.Marshaler        = (*Currency)(nil)
	_ driver.ValueConverter = (*Currency)(nil)
	_ driver.Valuer         = (*Currency)(nil)
)

func (m *Currency) MarshalJSON() ([]byte, error) {

}

func (m *Currency) UnmarshalJSON(b []byte) error {

}

func (m *Currency) ConvertValue(v interface{}) (driver.Value, error) {

}
func (m *Currency) Value() (driver.Value, error) {

}
