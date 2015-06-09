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
Package money uses a fixed-length guard for precision arithmetic.
Implements Un/Marshaller and Scan() method for database
columns including null, optimized for decimal(12, 4) fields.

Rounding is done on float64 to int64 by	the Rnd() function truncating
at values less than (.5 + (1 / Guardf))	or greater than -(.5 + (1 / Guardf))
in the case of negative numbers. The Guard adds four decimal places
of protection to rounding.
Decimal precision can be changed in the Precision() option
function. Precision() hold the places after the decimal place in teh active money struct field m.

http://en.wikipedia.org/wiki/Floating_point#Accuracy_problems

Options

The following options can be set while calling New():

	m := New(Swedish(Interval005), Guard(100), Precision(100))

Those values are really optional and even the order they appear ;-).
Default settings are:

	Precision 10000 which reflects decimal(12,4) database field
	Guard 	  10000 which reflects decimal(12,4) database field
	Swedish   No rounding

If you need to temporarily set a different option value you can stick to this pattern:
http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html

	prev := m.Option(Swedish(Interval005))
	defer m.Option(prev)
	// do something with the different Swedish rounding

Initial Idea: Copyright (c) 2011 Jad Dittmar
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
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"

	"database/sql"

	"bytes"

	"github.com/corestoreio/csfw/utils/log"
)

var (
	ErrOverflow = errors.New("Integer Overflow")

	guard   int64   = 10000
	guardi          = int(guard)
	guardf  float64 = float64(guard)
	dp      int64   = 10000
	dpi             = int(dp)
	dpf     float64 = float64(dp)
	swedish         = Interval000

	Round = .5
	//	Round  = .5 + (1 / Guardf)
	Roundn = Round * -1
)

// Interval* constants http://en.wikipedia.org/wiki/Swedish_rounding
const (
	// Interval000 no swedish rounding (default)
	Interval000 Interval = iota
	// Interval005 rounding with 0.05 intervals
	Interval005
	// Interval010 rounding with 0.10 intervals
	Interval010
	// Interval050 rounding with 0.50 intervals
	Interval050
	// Interval100 rounding with 1.00 intervals
	Interval100
	// interval999 max, not available
	interval999
)

type (
	// Interval defines the type for the Swedish rounding.
	Interval uint8

	// Currency represents a money aka currency type to avoid rounding errors with floats.
	// Takes also care of http://en.wikipedia.org/wiki/Swedish_rounding
	Currency struct {
		// m money in Guard/DP
		m int64
		// Locale to allow language specific output formats @todo
		Locale string
		// Valid if false the internal value is NULL
		Valid bool
		// swedish is the interval for the swedish rounding.
		swedish Interval

		guard  int64
		guardf float64
		dp     int64
		dpf    float64
	}

	OptionFunc func(*Currency) OptionFunc
)

// DefaultSwedish sets the global and New() defaults swedish rounding
// http://en.wikipedia.org/wiki/Swedish_rounding
func DefaultSwedish(i Interval) {
	if i < interval999 {
		swedish = i
	} else {
		log.Error("money=SetSwedishRounding", "err", errors.New("Interval out of scope"), "interval", i)
	}
}

// DefaultGuard sets the global default guard. A fixed-length guard for precision arithmetic.
func DefaultGuard(g int64) {
	guard = g
	guardf = float64(g)
}

// DefaultPrecision sets the global default decimal precision.
// 2 decimal places => 10^2; 3 decimal places => 10^3; x decimal places => 10^x
func DefaultPrecision(p int64) {
	p64 := int64(p)
	l := int64(math.Log(float64(p64)))
	if p64 == 0 || (p64 != 0 && (l%2) != 0) {
		p64 = dp
	}
	dp = p64
	dpf = float64(p64)
}

// Swedish sets the Swedish rounding
// http://en.wikipedia.org/wiki/Swedish_rounding
func Swedish(i Interval) OptionFunc {
	if i >= interval999 {
		log.Error("Currency=SetSwedishRounding", "err", errors.New("Interval out of scope. Resetting."), "interval", i)
		i = Interval000
	}
	return func(c *Currency) OptionFunc {
		previous := c.swedish
		c.swedish = i
		return Swedish(previous)
	}
}

// SetGuard sets the guard
func Guard(g int) OptionFunc {
	if g == 0 { // check for division by zero
		g = 1
	}
	return func(c *Currency) OptionFunc {
		previous := int(c.guard)
		c.guard = int64(g)
		c.guardf = float64(g)
		return Guard(previous)
	}
}

// Precision sets the precision.
// 2 decimal places => 10^2; 3 decimal places => 10^3; x decimal places => 10^x
// If not a decimal power then falls back to the default value.
func Precision(p int) OptionFunc {
	p64 := int64(p)
	l := int64(math.Log(float64(p64)))
	if p64 != 0 && (l%2) != 0 {
		p64 = dp
	}
	if p64 == 0 { // check for division by zero
		p64 = 1
	}
	return func(c *Currency) OptionFunc {
		previous := int(c.dp)
		c.dp = p64
		c.dpf = float64(p64)
		return Precision(previous)
	}
}

// New creates a new empty Currency struct with package default values of
// Guard and decimal precision.
func New(opts ...OptionFunc) Currency {
	c := Currency{
		guard:  guard,
		guardf: guardf,
		dp:     dp,
		dpf:    dpf,
	}
	c.Option(opts...)
	return c
}

// Options besides New() also Option() can apply options to the current
// struct. It returns the last set option. More info about the returned function:
// http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
func (c *Currency) Option(opts ...OptionFunc) (previous OptionFunc) {
	for _, o := range opts {
		if o != nil {
			previous = o(c)
		}
	}
	return previous
}

// Abs Returns the absolute value of Currency
func (c Currency) Abs() Currency {
	if c.m < 0 {
		return c.Neg()
	}
	return c
}

// Add Adds two Currency types. Returns empty Currency on integer overflow
func (c Currency) Add(d Currency) Currency {
	r := c.m + d.m
	if (r^c.m)&(r^d.m) < 0 {
		if log.IsTrace() {
			log.Trace("Currency=Add", "err", ErrOverflow, "m", c, "n", d)
		}
		log.Error("Currency=Add", "err", ErrOverflow, "m", c, "n", d)
		return New()
	}
	c.m = r
	c.Valid = true
	return c
}

// Getf gets the float64 value of money (see Raw() for int64)
func (c Currency) Getf() float64 {
	return float64(c.m) / c.dpf
}

// Geti gets value of money truncating after decimal precision (see Raw() for no truncation).
// Rounds always down
func (c Currency) Geti() int64 {
	return c.m / c.dp
}

// Raw returns in int64 the value of Currency (also see Gett(), See Get() for float64)
func (c Currency) Raw() int64 {
	return c.m
}

// Set sets the raw Currency field m
func (c Currency) Set(i int64) Currency {
	c.m = i
	c.Valid = true
	return c
}

// Setf sets a float64 into a Currency type for precision calculations
func (c Currency) Setf(f float64) Currency {
	fDPf := f * c.dpf
	r := int64(f * c.dpf)
	c.Valid = true
	return c.Set(Rnd(r, fDPf-float64(r)))
}

// Sign returns the Sign of Currency 1 if positive, -1 if negative
func (c Currency) Sign() int {
	if c.m < 0 {
		return -1
	}
	return 1
}

// String for money type representation in basic monetary unit (DOLLARS CENTS)
// @todo consider locale
func (c Currency) String() string {
	return fmt.Sprintf("%d.%02d", c.Raw()/c.dp, c.Abs().Raw()%c.dp)
}

// Sub subtracts one Currency type from another. Returns empty Currency on integer overflow
func (c Currency) Sub(d Currency) Currency {
	r := c.m - d.m
	if (r^c.m)&^(r^d.m) < 0 {
		if log.IsTrace() {
			log.Trace("Currency=Sub", "err", ErrOverflow, "m", c, "n", d)
		}
		log.Error("Currency=Sub", "err", ErrOverflow, "m", c, "n", d)
		return New()
	}
	c.m = r
	return c
}

// Mul Multiplies two Currency types
func (c Currency) Mul(d Currency) Currency {
	return c.Set(c.m * d.m / c.dp)
}

// Div Divides one Currency type from another
func (c Currency) Div(d Currency) Currency {
	f := c.guardf * c.dpf * float64(c.m) / float64(d.m) / c.guardf
	i := int64(f)
	return c.Set(Rnd(i, f-float64(i)))
}

// Mulf Multiplies a Currency with a float to return a money-stored type
func (c Currency) Mulf(f float64) Currency {
	i := c.m * int64(f*c.guardf*c.dpf)
	r := i / c.guard / c.dp
	return c.Set(Rnd(r, float64(i)/c.guardf/c.dpf-float64(r)))
}

// Neg Returns the negative value of Currency
func (c Currency) Neg() Currency {
	if c.m != 0 {
		c.m *= -1
	}
	return c
}

// Pow is the power of Currency
func (c Currency) Pow(f float64) Currency {
	return c.Setf(math.Pow(c.Getf(), f))
}

// Rnd rounds int64 remainder rounded half towards plus infinity
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
	_          json.Unmarshaler = (*Currency)(nil)
	_          json.Marshaler   = (*Currency)(nil)
	_          sql.Scanner      = (*Currency)(nil)
	nullString                  = []byte("null")

//	_ driver.ValueConverter = (Currency)(nil)
//	_ driver.Valuer         = (Currency)(nil)
)

func (c Currency) MarshalJSON() ([]byte, error) {
	// @todo should be possible to output the value without the currency sign
	// or output it as an array e.g.: [1234.56, "1.234,56€", "€"]
	// hmmmm
	if false == c.Valid {
		return nullString, nil
	}
	return []byte(`"` + c.String() + `"`), nil
}

func (c *Currency) UnmarshalJSON(b []byte) error {
	// @todo rewrite and optimize unmarshalling but for now json.Unmarshal is fine
	var s interface{}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return c.Scan(s)
}

// @todo quick write down without tests so add tests 8-)
func (c *Currency) Scan(value interface{}) error {
	if value == nil {
		c.m, c.Valid = 0, false
		return nil
	}
	if c.guard == 0 {
		c.Option(Guard(guardi))
	}
	if c.dp == 0 {
		c.Option(Precision(dpi))
	}

	if rb, ok := value.(*sql.RawBytes); ok {
		f, err := atof64([]byte(*rb))
		if err != nil {
			return log.Error("Currency=Scan", "err", err)
		}
		c.Valid = true
		c.Setf(f)
	}
	return nil
}

var colon = []byte(",")

func atof64(bVal []byte) (f float64, err error) {
	bVal = bytes.Replace(bVal, colon, nil, -1)
	//	s := string(bVal)
	//	s1 := strings.Replace(s, ",", "", -1)
	f, err = strconv.ParseFloat(string(bVal), 64)
	return f, err
}

//// ConvertValue @todo ?
//func (c Currency) ConvertValue(v interface{}) (driver.Value, error) {
//	return nil, nil
//}
//
//// Value @todo ?
//func (c Currency) Value() (driver.Value, error) {
//	return nil, nil
//}
