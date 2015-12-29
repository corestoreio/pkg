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

package directory_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/valuelabel"
	"github.com/corestoreio/csfw/directory"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

//func TestCountry(t *testing.T) {
//	ch := language.MustParseRegion("CH")
//	tld, err := ch.Canonicalize().TLD()
//	assert.NoError(t, err)
//	t.Logf("\n%#v\n", tld.String())
//}

func TestDefaultCountry(t *testing.T) {
	t.Log("@todo")
}

func TestPathCountryAllowed(t *testing.T) {

	directory.PathCountryAllowed.ValueLabel = valuelabel.NewByString("DE", "Germany", "AU", "'Straya", "CH", "Switzerland")

	cr := config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			directory.PathCountryAllowed.FQPathInt64(scope.StrStores, 1): "DE,AU,CH,AT",
		}),
	)

	haveCountries := directory.PathCountryAllowed.Get(directory.PackageConfiguration, cr.NewScoped(1, 1, 1))

	assert.Exactly(t, []string{"DE", "AU", "CH", "AT"}, haveCountries)

	// todo validation

}

//func TestAllowedCountriesDefault(t *testing.T) {
//	cr := config.NewMockGetter(
//		config.WithMockValues(config.MockPV{}),
//	)
//
//	haveCountries, err := directory.AllowedCountries(cr.NewScoped(1, 1, 1))
//	assert.NoError(t, err)
//	assert.True(t, len(haveCountries) > 100)
//}
