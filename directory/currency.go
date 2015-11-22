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

package directory

import (
	"strings"

	"github.com/corestoreio/csfw/config"
	"golang.org/x/text/currency"
)

type Currency struct {
	// https://godoc.org/golang.org/x/text/language
	c currency.Currency
}

// BaseCurrencyCode retrieves application base currency code
func BaseCurrencyCode(cr config.Reader) (currency.Currency, error) {
	base, err := cr.GetString(config.Path(PathCurrencyBase))
	if config.NotKeyNotFoundError(err) {
		return currency.Currency{}, err
	}
	return currency.ParseISO(base)
}

// AllowedCurrencies returns all installed currencies from global scope.
func AllowedCurrencies(cr config.Reader) ([]string, error) {
	installedCur, err := cr.GetString(config.Path(PathSystemCurrencyInstalled))
	if config.NotKeyNotFoundError(err) {
		return nil, err
	}
	// TODO use internal model of PathSystemCurrencyInstalled defined in package directory
	return strings.Split(installedCur, ","), nil
}
