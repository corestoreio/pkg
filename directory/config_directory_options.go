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

import (
	"sync"

	"github.com/corestoreio/csfw/config/valuelabel"
	"github.com/corestoreio/csfw/storage/dbr"
	"golang.org/x/text/currency"
)

type PkgOptions struct {
	sync.Mutex
	Currency valuelabel.Slice
	Country  valuelabel.Slice
}

// Options global valuelabel slice containg available options for a
// configuration type. Changing a valuelabel.Slice here affects all
// other Paths.
var Options *PkgOptions

func InitOptions(dbrSess dbr.SessionRunner) error {
	o := new(PkgOptions)
	o.Lock()
	defer o.Unlock()

	initOptionCurrency(dbrSess, o)
	initOptionCountry(dbrSess, o)

	Options = o
	return nil
}

// initOptionCurrency sets the Options() on all PathCurrency* configuration
// global variables.
func initOptionCurrency(dbrSess dbr.SessionRunner, o *PkgOptions) error {
	Path.Lock()
	defer Path.Unlock()

	o.Currency = valuelabel.NewByStringValue(currency.All()...)

	Path.SystemCurrencyInstalled.ValueLabel = o.Currency
	Path.CurrencyOptionsBase.ValueLabel = o.Currency
	Path.CurrencyOptionsAllow.ValueLabel = o.Currency
	Path.CurrencyOptionsDefault.ValueLabel = o.Currency
	return nil
}

// initOptionCountry should run every time your service initializes or
// a value in the database changes.
func initOptionCountry(dbrsess dbr.SessionRunner, o *PkgOptions) error {
	Path.Lock()
	defer Path.Unlock()

	// o.Country
	// TODO load from database the iso code and as value the names
	o.Country = valuelabel.NewByString("AU", "Australia", "FI", "Finland", "DE", "Germany")

	// apply the list of country codes to:
	Path.GeneralCountryDefault.ValueLabel = o.Country
	Path.GeneralCountryAllow.ValueLabel = o.Country
	Path.GeneralCountryOptionalZipCountries.ValueLabel = o.Country
	Path.GeneralCountryEuCountries.ValueLabel = o.Country
	Path.GeneralCountryDestinations.ValueLabel = o.Country
	return nil
}
