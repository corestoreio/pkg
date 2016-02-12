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

import (
	"math"
	"sync"

	"github.com/corestoreio/csfw/i18n"
)

type globalSettings struct {
	sync.Mutex
	guardi  int
	dpi     int
	swedish Interval
}

var global = &globalSettings{
	guardi:  10000,
	dpi:     10000,
	swedish: Interval000,
}

// gNaN will be returned if Valid is false in the Currency struct
var gNaN = []byte(`NaN`)

// DefaultFormatterCurrency sets the package wide default locale specific currency formatter.
// This variable can be overridden.
var DefaultFormatterCurrency i18n.CurrencyFormatter = i18n.DefaultCurrency

// DefaultFormatterNumber sets the package wide default locale specific number formatter
// This variable can be overridden.
var DefaultFormatterNumber i18n.NumberFormatter = i18n.DefaultNumber

// DefaultJSONEncode is JSONLocale
var DefaultJSONEncode Encoder = NewJSONEncoder()

// DefaultJSONDecode is JSONLocale
var DefaultJSONDecode Decoder = NewJSONDecoder()

// DefaultSwedish sets the global and New() defaults swedish rounding. Errors will be logged.
// Invalid intervals gets ignored. Returns the successful applied value.
// http://en.wikipedia.org/wiki/Swedish_rounding
func DefaultSwedish(i Interval) Interval {
	global.Lock()
	defer global.Unlock()
	if i < interval999 {
		global.swedish = i
	}
	return global.swedish
}

// DefaultGuard sets the global default guard. A fixed-length guard for precision
// arithmetic. Returns the successful applied value.
func DefaultGuard(g int) int {
	global.Lock()
	defer global.Unlock()
	if g == 0 {
		g = 1
	}
	global.guardi = g
	return global.guardi
}

// DefaultPrecision sets the global default decimal precision.
// 2 decimal places => 10^2; 3 decimal places => 10^3; x decimal places => 10^x
// Returns the successful applied value.
func DefaultPrecision(p int) int {
	global.Lock()
	defer global.Unlock()
	l := int64(math.Log(float64(p)))
	if p == 0 || (p != 0 && (l%2) != 0) {
		p = global.dpi
	}
	global.dpi = p
	return global.dpi
}
