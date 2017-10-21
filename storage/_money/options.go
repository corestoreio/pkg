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

// WithPrecision sets the precision. 2 decimal places => 10^2; 3 decimal places
// => 10^3; x decimal places => 10^x. If not a decimal power then falls back to
// the default value.
//func WithPrecision(p int) Option {
//	return func(c *Money) Option {
//		global.Lock()
//		defer global.Unlock()
//		previous := int(c.dp)
//		c.dp, c.dpf, c.prec = precision(p)
//		return WithPrecision(previous)
//	}
//}

// precision internal prec generator. Optimized to reduce allocs
func precision(p int) (int64, float64, int) {
	p64 := int64(p)
	l := int64(math.Log(float64(p64)))
	if p64 != 0 && (l%2) != 0 {
		p64 = int64(10000)
	}
	if p64 == 0 { // check for division by zero
		p64 = 1
	}
	return p64, float64(p64), decimals(p64)
}
