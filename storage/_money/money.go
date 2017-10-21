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

package money

/*
	@todo
	- http://unicode.org/reports/tr35/tr35-numbers.html#Supplemental_Currency_Data to automatically
	- set the Swedish rounding
	- Currency
	- https://github.com/golang/go/issues/12127 decimal type coming to math/big package
	- https://github.com/golang/go/issues/12332
		This means that for a large company, a decimal data type must cope with
		figures like: 100,000,000,000.00000000, that is, 20 significant digits.
	- https://github.com/shopspring/decimal -> get inspiration
	- https://github.com/EricLagergren/decimal -> get inspiration
	- https://github.com/moneyphp/money
	- http://verraes.net/2016/02/type-safety-and-money/
	- https://github.com/Rhymond/go-money
*/

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/corestoreio/errors"
)

const guardf = 100000

// Currency represents a money aka currency type to avoid rounding errors with
// floats. Includes options for printing, Swedish rounding, database scanning
// and JSON en/decoding.
type Money struct {
	// m money in Guard/DP
	m int64
	// Valid if false the internal value is NULL
	Valid bool

	// guardf needed for divide and multiply operation
	guardf float64
	prec   int32 // precision of decimal places, eg: 10,100,1000
	pos    int8  // position exponent
}

// New creates a new empty Money struct with package default values. Formatter
// can be overridden after you have created the new type. Implements the
// interfaces: database.Scanner, driver.Valuer, json.Marshaller,
// json.Unmarshaller
func New(value int64, exp int) Money {
	return Money{
		m:      value,
		guardf: guardf,
		prec:   int32(math.Pow10(exp)),
		pos:    int8(exp),
	}
}

// Abs returns the absolute value of Currency
func (m Money) Abs() Money {
	if m.m < 0 {
		return m.Neg()
	}
	return m
}

// Getf gets the float64 value of money (see Raw() for int64)
func (m Money) Getf() float64 {
	return float64(m.m) / float64(m.prec)
}

// Geti gets value of money truncating after decimal precision (see Raw() for no
// truncation). Rounds always down.
func (m Money) Geti() int64 {
	return m.m / int64(m.prec)
}

// Dec returns the decimals
func (m Money) Dec() int64 {
	return m.Abs().Raw() % int64(m.prec)
}

// Raw returns in int64 the value of Currency (also see Geti(), See Getf() for
// float64)
func (m Money) Raw() int64 {
	return m.m
}

// Set sets the raw Currency field m
func (m Money) Set(i int64) Money {
	m.m = i
	m.Valid = true
	return m
}

// Set sets the raw Currency field m
func (m Money) setRnd(f float64) Money {
	return m.Set(int64(math.Round(f)))
}

// Setf sets a float64 into a Currency type for precision calculations
func (m Money) Setf(f float64) Money {
	return m.setRnd(f * float64(m.prec))
}

// ParseFloat transforms a string float value into a real float64 value and sets
// it. Current value will be overridden. Returns a logged error.
func (m *Money) ParseFloat(s string) error {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	m.Valid = true
	*m = m.Setf(f)
	return nil
}

// Sign returns:
//
//	-1 if x <  0
//	+1 if x >=  0
//
func (m Money) Sign() int {
	if m.m < 0 {
		return -1
	}
	return 1
}

// Precision returns the amount of decimal digits
func (m Money) Precision() int {
	return int(m.pos)
}

// String for money type representation in a specific locale. Errors will be
// written to the buffer.
func (m Money) String() string {
	prec64 := int64(m.prec)
	neg := ""
	v := m.m / prec64
	if m.Sign() < 0 && v == 0 {
		neg = "-"
	}
	pPos := "2"
	lZero := "0"
	if m.pos > 0 {
		pPos = strconv.FormatInt(int64(m.pos), 10)
		lZero = strings.Repeat("0", int(m.pos-1))
	}
	// TODO optimize; remove decimals when there are none
	return fmt.Sprintf("%s%d.%"+lZero+pPos+"d", neg, v, m.Abs().m%prec64)

}

// GoString implements fmt.GoStringer
func (m Money) xGoString() string {
	return fmt.Sprintf("money.New(money.WithPrecision(%d)).Set(%d)", m.prec, m.Raw())
}

// FtoaAppend converts the internal floating-point number to a byte slice
// without any applied formatting and appends it to dst and returns the extended
// buffer.
func (m Money) FtoaAppend(dst []byte) []byte {
	if false == m.Valid {
		return append(dst, gNaN...)
	}
	return strconv.AppendFloat(dst, m.Getf(), 'f', int(m.prec), 64)
}

func (m Money) Rescale(exp int) Money {
	//exp8 := int8(exp)
	//if m.pos == exp8 {
	//	return m
	//}
	//
	//// up scale; add more decimal places
	//if diff := exp8 - m.pos; exp8 > m.pos && diff > 0 {
	//	m.pos += diff
	//	m.prec = int32(math.Pow10(int(m.pos)))
	//	m.m *= int64(math.Pow10(int(diff)))
	//	return m
	//}
	//// down scale
	//if diff := exp8 - m.pos; exp8 < m.pos && diff < 0 {
	//	// down scale: 2.545 to 2.6
	//	dec := float64(m.Dec()) / float64(m.prec) // 0.545
	//
	//	diff *= -1
	//	m.pos -= diff
	//
	//	deci64 := int64(math.Round(dec)) // / math.Pow10(int(m.pos)) // 0.1
	//	m2 := m.m / int64(math.Pow10(int(diff)))
	//	m.m = m2 + deci64
	//	m.prec = int32(math.Pow10(int(m.pos)))
	//
	//	//println("m.pos -= diff", m.pos, diff, "old prec", m.prec)
	//	//m.prec = int32(math.Pow10(int(m.pos)))
	//	//println("new prc", m.prec, "new pos", m.pos)
	//	//
	//	//toRound := float64(m.m) / math.Pow10(int(diff))
	//	//println("toRound", toRound)
	//	//newM := math.Round(toRound)
	//	//m.m = int64(newM)
	//}

	return m
}

// Add adds two Currency types. Returns empty Currency on integer overflow.
// Errors gets appended to the Multi Error type. Panics on integer overflow.
func (m Money) Add(d Money) Money {
	if r := m.m + d.m; (r^m.m)&(r^d.m) < 0 {
		panic(errors.NewOverflowedf("[money] Integer overflow"))
	}
	mRaw := float64(m.m)*m.guardf/float64(m.prec) + (float64(d.m) * d.guardf / float64(d.prec))
	mRaw *= float64(m.prec) / m.guardf
	mRaw = math.Round(mRaw)
	m.m = int64(mRaw)
	m.Valid = true
	return m
}

// Sub subtracts one Currency type from another. Returns empty Currency on
// integer overflow. Errors gets appended to the Multi Error type. Panics on
// integer overflow.
func (m Money) Sub(d Money) Money {
	r := m.m - d.m
	if (r^m.m)&^(r^d.m) < 0 {
		panic(errors.NewOverflowedf("[money] Integer overflow"))
	}
	m.m = r
	return m
}

// Mul multiplies two Currency types. Both types must have the same precision.
// Panics on integer overflow.
func (m Money) Mul(d Money) Money {
	mP := m.m * d.m
	if m.m > 0 && d.m > 0 && mP < 0 {
		panic(fmt.Sprintf("[money] %d overflows Int (negative)", mP))
	}
	if m.m < 0 && d.m < 0 && mP > 0 {
		panic(fmt.Sprintf("[money] %d overflows Int (positive)", mP))
	}
	r := math.Round(float64(m.m*d.m) / float64(m.prec))
	return m.Set(int64(r))
}

// Div divides one Currency type from another
func (m Money) Div(d Money) Money {
	f := (m.guardf * float64(m.prec) * float64(m.m)) / float64(d.m) / m.guardf
	return m.Set(int64(math.Round(f)))
}

// Mulf multiplies a Currency with a float to return a money-stored type
func (m Money) Mulf(f float64) Money {
	// TODO: Check if implementation is flawed after my refactoring. See git history.
	i := float64(m.m) * f * m.guardf * float64(m.prec)
	r := i / m.guardf / float64(m.prec)
	return m.Set(int64(math.Round(r)))
}

// Neg returns the negative value of Currency
func (m Money) Neg() Money {
	if m.m != 0 {
		m.m *= -1
	}
	return m
}

// Pow is the power of Currency
func (m Money) Pow(f float64) Money {
	return m.Setf(math.Pow(m.Getf(), f))
}

// Swedish applies the Swedish rounding. You may set the usual options.
// Interval* constants http://en.wikipedia.org/wiki/Swedish_rounding
// Interval000 no swedish rounding (default)
// Interval005 rounding with 0.05 intervals
// Interval010 rounding with 0.10 intervals
// Interval015 same as Interval010 except that 5 will be rounded down.
// 0.45 => 0.40 or 0.46 => 0.50
// Special case for New Zealand (a must visit!), it is up to the business to
// decide if they will round 5¢ intervals up or down. The majority of
// retailers follow government advice and round it down. Use then
// Interval015. otherwise use Interval010.
// Interval025 rounding with 0.25 intervals
// Interval050 rounding with 0.50 intervals
// Interval100 rounding with 1.00 intervals
func (m Money) RoundCash(interval uint8) Money {

	switch interval {
	case 5:
		// NL, SG, SA, CH, TR, CL, IE
		// 5 cent rounding
		return m.Setf(math.Round(m.Getf()*20) / 20) // base 5
	case 10:
		// New Zealand & Hong Kong
		// 10 cent rounding
		// In Sweden between 1985 and 1992, prices were rounded up for sales
		// ending in 5 öre.
		return m.Setf(math.Round(m.Getf()*10) / 10)
	case 15:
		// 10 cent rounding, special case
		// Special case: In NZ, it is up to the business to decide if they
		// will round 5¢ intervals up or down. The majority of retailers follow
		// government advice and round it down.
		if m.m%5 == 0 {
			m.m = m.m - 1
		}
		return m.Setf(math.Round(m.Getf()*10) / 10)
	case 25:
		// round to quarter
		return m.Setf(math.Round(m.Getf()*4) / 4)
	case 50:
		// 50 cent rounding
		// The system used in Sweden from 1992 to 2010, in Norway from 1993 to 2012,
		// and in Denmark since 1 October 2008 is the following:
		// Sales ending in 1–24 öre round down to 0 öre.
		// Sales ending in 25–49 öre round up to 50 öre.
		// Sales ending in 51–74 öre round down to 50 öre.
		// Sales ending in 75–99 öre round up to the next whole Krone/krona.
		return m.Setf(math.Round(m.Getf()*2) / 2)
	case 100:
		// The system used in Sweden since 30 September 2010 and used in Norway since 1 May 2012.
		// Sales ending in 1–49 öre/øre round down to 0 öre/øre.
		// Sales ending in 50–99 öre/øre round up to the next whole krona/krone.
		return m.Setf(math.Round(m.Getf()*1) / 1) // ;-)
	}
	return m
}

// CompareTo depends on the Valuta field (TODO)
func (m Money) CompareTo(d Money) bool {
	return false
}

// decimals returns the length of the 10^n calculation, ignoring the leading 1.
// For other int64 you must add + 1.
func decimals(dec int64) int {
	if dec < 1 {
		return 0
	}
	return int(math.Floor(math.Log10(float64(dec))))
}

// max returns the max number from two numbers
func maxInt8(a, b int8) int8 {
	if a > b {
		return a
	}
	return b
}
