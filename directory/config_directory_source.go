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

	"github.com/corestoreio/csfw/config/source"
	"github.com/corestoreio/csfw/storage/dbr"
	"golang.org/x/text/currency"
)

// PkgSource just exported for the sake of documentation. See fields
// for more information.
type PkgSource struct {
	sync.Mutex
	// Currency contains all possible currencies on this planet.
	Currency source.Slice
	// Country should contain all countries
	Country source.Slice
}

// Sources global slice containing available sources to be used in
// the Path models. Changing a source.Slice here affects all
// other Paths models.
var Sources *PkgSource

// InitSources initializes the global variable Sources.
func InitSources(dbrSess dbr.SessionRunner) error {
	o := new(PkgSource)
	o.Lock()
	defer o.Unlock()

	initSourceCurrency(dbrSess, o)
	initSourceCountry(dbrSess, o)

	Sources = o
	return nil
}

// initSourceCurrency sets the Options() on all PathCurrency* configuration
// global variables.
func initSourceCurrency(dbrSess dbr.SessionRunner, o *PkgSource) error {
	Backend.Lock()
	defer Backend.Unlock()

	o.Currency = source.NewByStringValue(currency.All()...)

	Backend.SystemCurrencyInstalled.Source = o.Currency
	Backend.CurrencyOptionsBase.Source = o.Currency
	Backend.CurrencyOptionsAllow.Source = o.Currency
	Backend.CurrencyOptionsDefault.Source = o.Currency
	return nil
}

// initSourceCountry should run every time your service initializes or
// a value in the database changes.
func initSourceCountry(dbrsess dbr.SessionRunner, o *PkgSource) error {
	Backend.Lock()
	defer Backend.Unlock()

	// o.Country
	// TODO load from database the iso code and as value the names
	o.Country = source.NewByString("AU", "Australia", "FI", "Finland", "DE", "Germany")

	// apply the list of country codes to:
	Backend.GeneralCountryDefault.Source = o.Country
	Backend.GeneralCountryAllow.Source = o.Country
	Backend.GeneralCountryOptionalZipCountries.Source = o.Country
	Backend.GeneralCountryEuCountries.Source = o.Country
	Backend.GeneralCountryDestinations.Source = o.Country
	return nil
}
