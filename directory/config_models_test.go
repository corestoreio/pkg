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

package directory_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/directory"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigCurrencyGet(t *testing.T) {
	cc := directory.NewConfigCurrency(directory.Path.CurrencyOptionsBase.String())

	cr := config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			directory.Path.CurrencyOptionsBase.FQPathInt64(scope.StrStores, 1): "EUR",
			directory.Path.CurrencyOptionsBase.FQPathInt64(scope.StrStores, 2): "WIR", // Special Swiss currency
		}),
	)

	cur, err := cc.Get(cr.NewScoped(1, 1, 1))
	assert.NoError(t, err)
	assert.Exactly(t, directory.MustNewCurrencyISO("EUR"), cur)

	cur, err = cc.Get(cr.NewScoped(1, 1, 2))
	assert.EqualError(t, err, "currency: tag is not a recognized currency")
	assert.Exactly(t, directory.Currency{}, cur)
}

func TestNewConfigCurrencyWrite(t *testing.T) {
	cc := directory.NewConfigCurrency(directory.Path.CurrencyOptionsBase.String())
	c := directory.MustNewCurrencyISO("EUR")

	w := new(config.MockWrite)

	assert.NoError(t, cc.Write(w, c, scope.WebsiteID, 33))

	assert.Exactly(t, directory.Path.CurrencyOptionsBase.FQPathInt64(scope.StrWebsites, 33), w.ArgPath)
	assert.Exactly(t, "EUR", w.ArgValue)
}
