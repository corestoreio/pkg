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

package directory

import "golang.org/x/text/currency"

// Currency represents a corestore currency type which may add more features.
type Currency struct {
	currency.Unit
}

// NewCurrencyISO creates a new Currency by parsing a 3-letter ISO 4217 currency
// code. It returns an error if s is not well-formed or not a recognized
// currency code.
func NewCurrencyISO(iso string) (c Currency, err error) {
	var u currency.Unit
	u, err = currency.ParseISO(iso)
	c.Unit = u
	return
}

// MustNewCurrencyISO same as NewCurrencyISO() but panics on error.
func MustNewCurrencyISO(iso string) Currency {
	c, err := NewCurrencyISO(iso)
	if err != nil {
		panic(err)
	}
	return c
}
