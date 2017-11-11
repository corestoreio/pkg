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

	"github.com/corestoreio/cspkg/config/cfgmock"
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/config/cfgpath"
	"github.com/corestoreio/cspkg/directory"
	"github.com/corestoreio/cspkg/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigCurrencyGetDefault(t *testing.T) {
	t.Parallel()

	cobPath, err := backend.CurrencyOptionsBase.ToPath(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	cr := cfgmock.NewService(cfgmock.PathValue{
		cobPath.Bind(scope.Default, 0).String(): "CHF",
	})

	cur, err := backend.CurrencyOptionsBase.GetDefault(cr)
	assert.NoError(t, err)
	assert.Exactly(t, "CHF", cur.String())
}

func TestNewConfigCurrencyGetDefaultPathError(t *testing.T) {
	t.Parallel()

	ccModel := directory.NewConfigCurrency("a/b/c")

	cr := cfgmock.NewService()

	cur, err := ccModel.GetDefault(cr)
	assert.EqualError(t, err, cfgpath.errIncorrectPath.Error())
	assert.Exactly(t, "XXX", cur.String())
}

func TestNewConfigCurrencyGetPathError(t *testing.T) {
	t.Parallel()

	ccModel := directory.NewConfigCurrency("a//c")

	cr := cfgmock.NewService()

	cur, err := ccModel.Get(cr.NewScoped(0, 0))
	assert.EqualError(t, err, "Route a/\uf8ff/c: This character \"\\uf8ff\" is not allowed in Route a/\uf8ff/c")
	assert.Exactly(t, "XXX", cur.String())
}

func TestNewConfigCurrencyGetEmpty(t *testing.T) {
	t.Parallel()

	cobPath, err := backend.CurrencyOptionsBase.ToPath(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	// this test shows a discouraged use of the NewConfigCurrency() model.
	ccModel := directory.NewConfigCurrency(backend.CurrencyOptionsBase.String())

	cr := cfgmock.NewService(cfgmock.PathValue{
		// default scope is enforced because NewConfigCurrency() has been created
		// with the ConfigStructure slice and so we're missing the *element.Field
		// with the special configuration
		cobPath.Bind(scope.Website, 1).String(): "CHF",
		cobPath.Bind(scope.Store, 1).String():   "EUR",
	})

	cur, err := ccModel.Get(cr.NewScoped(1, 1))
	assert.EqualError(t, err, `Empty currency for path: "currency/options/base", scope: "Store", scopeID: 1`)
	assert.Exactly(t, "XXX", cur.String())
}

func TestNewConfigCurrencyGet(t *testing.T) {
	t.Parallel()

	cobPath, err := backend.CurrencyOptionsBase.ToPath(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	cr := cfgmock.NewService(cfgmock.PathValue{
		cobPath.Bind(scope.Website, 1).String(): "EUR",
		cobPath.Bind(scope.Website, 2).String(): "WIR", // Special Swiss currency
	})

	// scope of CurrencyOptionsBase set to website, so no store config values are possible
	cur, err := backend.CurrencyOptionsBase.Get(cr.NewScoped(1, 1))
	if err != nil {
		t.Fatal(err, cur)
	}

	assert.Exactly(t, directory.MustNewCurrencyISO("EUR"), cur)

	cur, err = backend.CurrencyOptionsBase.Get(cr.NewScoped(2, 1))
	assert.EqualError(t, err, "currency: tag is not a recognized currency")
	assert.Exactly(t, directory.Currency{}, cur)
}

func TestNewConfigCurrencyWrite(t *testing.T) {
	t.Parallel()

	cfgStruct, err := directory.NewConfigStructure()
	if err != nil {
		t.Fatal(err)
	}

	// special setup for testing
	cc := directory.NewConfigCurrency(
		backend.CurrencyOptionsBase.String(),
		cfgmodel.WithFieldFromSectionSlice(cfgStruct),
		cfgmodel.WithSourceByString("EUR", "Euro", "CHF", "Swiss Franc", "AUD", "Australian Dinar ;-)"),
	)

	c := directory.MustNewCurrencyISO("EUR")

	cobPath, err := backend.CurrencyOptionsBase.ToPath(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	w := new(cfgmock.Write)
	assert.NoError(t, cc.Write(w, c, scope.Website, 33))

	assert.Exactly(t, cobPath.Bind(scope.Website, 33).String(), w.ArgPath)
	assert.Exactly(t, "EUR", w.ArgValue)

	assert.EqualError(t,
		cc.Write(w, directory.Currency{}, scope.Website, 33),
		"The value 'XXX' cannot be found within the allowed Options():\n[{\"Value\":\"EUR\",\"Label\":\"Euro\"},{\"Value\":\"CHF\",\"Label\":\"Swiss Franc\"},{\"Value\":\"AUD\",\"Label\":\"Australian Dinar ;-)\"}]\n",
	)
}
