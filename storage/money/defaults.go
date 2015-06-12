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

package money

import (
	"errors"
	"math"

	"github.com/corestoreio/csfw/i18n"
	"github.com/corestoreio/csfw/utils/log"
)

var (
	guard   int64   = 10000
	guardi          = int(guard)
	guardf  float64 = float64(guard)
	dp      int64   = 10000
	dpi             = int(dp)
	dpf     float64 = float64(dp)
	swedish         = Interval000
)

// DefaultFormat sets the package wide default locale specific currency formatter
var DefaultFormat i18n.CurrencyFormatter = i18n.DefaultCurrency

// DefaultJSONEncode is struct JSONLocale
var DefaultJSONEncode JSONMarshaller = JSONLocale{}

// DefaultJSONDecode is struct JSONLocale
var DefaultJSONDecode JSONUnmarshaller = JSONLocale{}

// DefaultSwedish sets the global and New() defaults swedish rounding. Errors will be logged.
// http://en.wikipedia.org/wiki/Swedish_rounding
func DefaultSwedish(i Interval) {
	if i < interval999 {
		swedish = i
	} else {
		log.Error("money=SetSwedishRounding", "err", errors.New("Interval out of scope"), "interval", i)
	}
}

// DefaultGuard sets the global default guard. A fixed-length guard for precision
// arithmetic. Returns the successful applied value.
func DefaultGuard(g int64) int64 {
	if g == 0 {
		g = 1
	}
	guard = g
	guardf = float64(g)
	return guard
}

// DefaultPrecision sets the global default decimal precision.
// 2 decimal places => 10^2; 3 decimal places => 10^3; x decimal places => 10^x
// Returns the successful applied value.
func DefaultPrecision(p int64) int64 {
	l := int64(math.Log(float64(p)))
	if p == 0 || (p != 0 && (l%2) != 0) {
		p = dp
	}
	dp = p
	dpf = float64(p)
	return dp
}
