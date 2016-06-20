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

import "math"

var (
	RoundTo = .5
	//	RoundTo  = .5 + (1 / Guardf)
	RoundToN = RoundTo * -1
)

// Interval* constants http://en.wikipedia.org/wiki/Swedish_rounding
const (
	// Interval000 no swedish rounding (default)
	Interval000 Interval = iota
	// Interval005 rounding with 0.05 intervals
	Interval005
	// Interval010 rounding with 0.10 intervals
	Interval010
	// Interval015 same as Interval010 except that 5 will be rounded down.
	// 0.45 => 0.40 or 0.46 => 0.50
	// Special case for New Zealand (a must visit!), it is up to the business to
	// decide if they will round 5¢ intervals up or down. The majority of
	// retailers follow government advice and round it down. Use then
	// Interval015. otherwise use Interval010.
	Interval015
	// Interval025 rounding with 0.25 intervals
	Interval025
	// Interval050 rounding with 0.50 intervals
	Interval050
	// Interval100 rounding with 1.00 intervals
	Interval100
	interval999
)

// Interval defines the type for the Swedish rounding.
type Interval uint8

// Option used to apply options to the Money struct
type Option func(*Money) Option

// WithSwedish sets the Swedish rounding
// http://en.wikipedia.org/wiki/Swedish_rounding Invalid interval falls back to
// Interval000.
func WithSwedish(i Interval) Option {
	if i >= interval999 {
		i = Interval000
	}
	return func(c *Money) Option {
		previous := c.Interval
		c.Interval = i
		return WithSwedish(previous)
	}
}

// WithCashRounding same as Swedish() option function, but: Rounding increment,
// in units of 10-digits. The default is 0, which means no rounding is to be
// done. Therefore, rounding=0 and rounding=1 have identical behavior. Thus with
// fraction digits of 2 and rounding increment of 5, numeric values are rounded
// to the nearest 0.05 units in formatting. With fraction digits of 0 and
// rounding increment of 50, numeric values are rounded to the nearest 50.
// Possible values: 5, 10, 15, 25, 50, 100.
// todo: refactor to use text/currency package
func WithCashRounding(rounding int) Option {
	// somehow that feels like ... not very nice
	i := Interval000
	switch rounding {
	case 5:
		i = Interval005
	case 10:
		i = Interval010
	case 15:
		i = Interval015
	case 25:
		i = Interval025
	case 50:
		i = Interval050
	case 100:
		i = Interval100
	}

	return func(c *Money) Option {
		var p int
		switch c.Interval {
		case Interval005:
			p = 5
		case Interval010:
			p = 10
		case Interval015:
			p = 15
		case Interval025:
			p = 25
		case Interval050:
			p = 50
		case Interval100:
			p = 100
		}
		c.Interval = i
		return WithCashRounding(p)
	}
}

// WithGuard sets the guard
func WithGuard(g int) Option {
	return func(c *Money) Option {
		previous := int(c.guard)
		c.guard, c.guardf = guard(g)
		return WithGuard(previous)
	}
}

// guard generates the guard value. Optimized to reduce allocs
func guard(g int) (int64, float64) {
	if g == 0 {
		g = 1
	}
	return int64(g), float64(g)
}

// WithPrecision sets the precision. 2 decimal places => 10^2; 3 decimal places
// => 10^3; x decimal places => 10^x. If not a decimal power then falls back to
// the default value.
func WithPrecision(p int) Option {
	return func(c *Money) Option {
		global.Lock()
		defer global.Unlock()
		previous := int(c.dp)
		c.dp, c.dpf, c.prec = precision(p)
		return WithPrecision(previous)
	}
}

// precision internal prec generator. Optimized to reduce allocs
func precision(p int) (int64, float64, int) {
	p64 := int64(p)
	l := int64(math.Log(float64(p64)))
	if p64 != 0 && (l%2) != 0 {
		p64 = int64(global.dpi)
	}
	if p64 == 0 { // check for division by zero
		p64 = 1
	}
	return p64, float64(p64), decimals(p64)
}
