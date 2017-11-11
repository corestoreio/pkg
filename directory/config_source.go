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

	"github.com/corestoreio/cspkg/config/source"
	"github.com/corestoreio/cspkg/storage/dbr"
	"golang.org/x/text/currency"
)

// PkgSource a container for all available source.Slice within this package.
// These sources will be applied to the models
// See fields for more information.
type PkgSource struct {
	sync.Mutex
	// Currency contains all possible currencies on this planet.
	Currency source.Slice
	// Country should contain all countries
	Country source.Slice
}

// InitSources initializes the global variable Sources in all models.
// Changing a source.Slice here affects all
// other models.
func (be *PkgBackend) InitSources(dbrSess dbr.SessionRunner) (*PkgSource, error) {
	o := new(PkgSource)

	if err := be.InitCurrency(dbrSess, o); err != nil {
		return nil, err
	}
	if err := be.InitCountry(dbrSess, o); err != nil {
		return nil, err
	}

	return o, nil
}

// InitCurrency sets the Options() on all PathCurrency* configuration
// global variables.
func (be *PkgBackend) InitCurrency(dbrSess dbr.SessionRunner, o *PkgSource) error {
	be.Lock()
	defer be.Unlock()
	o.Lock()
	defer o.Unlock()

	// apply all world wide available currencies but extract them from the DB
	o.Currency = source.NewByStringValue(currency.All()...) // DB query

	be.SystemCurrencyInstalled.Source = o.Currency
	be.CurrencyOptionsBase.Source = o.Currency
	be.CurrencyOptionsAllow.Source = o.Currency
	be.CurrencyOptionsDefault.Source = o.Currency
	return nil
}

// InitCountry should run every time your service initializes or
// a value in the database changes.
func (be *PkgBackend) InitCountry(dbrsess dbr.SessionRunner, o *PkgSource) error {
	be.Lock()
	defer be.Unlock()
	o.Lock()
	defer o.Unlock()

	// o.Country
	// TODO load from database the iso code and as value the names
	o.Country = source.NewByString("AU", "Australia", "FI", "Finland", "DE", "Germany")

	// apply the list of country codes to:
	be.GeneralCountryDefault.Source = o.Country
	be.GeneralCountryAllow.Source = o.Country
	be.GeneralCountryOptionalZipCountries.Source = o.Country
	be.GeneralCountryEuCountries.Source = o.Country
	be.GeneralCountryDestinations.Source = o.Country
	return nil
}
