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
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/directory"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigCurrencyGet(t *testing.T) {
	t.Parallel()
	cc := directory.NewConfigCurrency(directory.Backend.CurrencyOptionsBase.String())

	cobPath, err := directory.Backend.CurrencyOptionsBase.ToPath(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	cr := config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			cobPath.Bind(scope.StoreID, 1).String(): "EUR",
			cobPath.Bind(scope.StoreID, 2).String(): "WIR", // Special Swiss currency
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
	t.Parallel()
	// special setup for testing
	cc := directory.NewConfigCurrency(
		directory.Backend.CurrencyOptionsBase.String(),
		model.WithConfigStructure(directory.ConfigStructure),
		model.WithSourceByString("EUR", "Euro", "CHF", "Swiss Franc", "AUD", "Australian Dinar ;-)"),
	)

	c := directory.MustNewCurrencyISO("EUR")

	cobPath, err := directory.Backend.CurrencyOptionsBase.ToPath(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	w := new(config.MockWrite)
	assert.NoError(t, cc.Write(w, c, scope.WebsiteID, 33))

	assert.Exactly(t, cobPath.Bind(scope.WebsiteID, 33).String(), w.ArgPath)
	assert.Exactly(t, "EUR", w.ArgValue)

	assert.EqualError(t,
		cc.Write(w, directory.Currency{}, scope.WebsiteID, 33),
		"The value 'XXX' cannot be found within the allowed Options():\n[{\"Value\":\"EUR\",\"Label\":\"Euro\"},{\"Value\":\"CHF\",\"Label\":\"Swiss Franc\"},{\"Value\":\"AUD\",\"Label\":\"Australian Dinar ;-)\"}]\n",
	)
}
