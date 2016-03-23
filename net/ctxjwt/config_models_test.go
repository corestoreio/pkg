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

package ctxjwt_test

import (
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigSigningMethodGetDefaultPathError(t *testing.T) {
	t.Parallel()

	ccModel := ctxjwt.NewConfigSigningMethod("a/x/c")

	cr := cfgmock.NewService()

	sm, err := ccModel.Get(cr.NewScoped(1, 1))
	assert.EqualError(t, err, "Route a/x/c: Incorrect Path. Either to short or missing path separator.")
	assert.Nil(t, sm)
}

func TestNewConfigSigningMethodGetPathError(t *testing.T) {
	t.Parallel()

	ccModel := ctxjwt.NewConfigSigningMethod("a//c")

	cr := cfgmock.NewService()

	sm, err := ccModel.Get(cr.NewScoped(0, 0))
	assert.EqualError(t, err, "Route a/\uf8ff/c: This character \"\\uf8ff\" is not allowed in Route a/\uf8ff/c")
	assert.Nil(t, sm)
}

//func TestNewConfigSigningMethodGetEmpty(t *testing.T) {
//	t.Parallel()
//
//	cobPath, err := backend.CurrencyOptionsBase.ToPath(0, 0)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// this test shows a discouraged use of the NewConfigSigningMethod() model.
//	ccModel := ctxjwt.NewConfigSigningMethod(backend.CurrencyOptionsBase.String())
//
//	cr := cfgmock.NewService(
//		cfgmock.WithPV(cfgmock.PathValue{
//			// default scope is enforced because NewConfigSigningMethod() has been created
//			// with the ConfigStructure slice and so we're missing the *element.Field
//			// with the special configuration
//			cobPath.Bind(scope.WebsiteID, 1).String(): "CHF",
//			cobPath.Bind(scope.StoreID, 1).String():   "EUR",
//		}),
//	)
//
//	cur, err := ccModel.Get(cr.NewScoped(1, 1))
//	assert.EqualError(t, err, `Empty currency for path: "currency/options/base", scope: "Store", scopeID: 1`)
//	assert.Exactly(t, "XXX", cur.String())
//}
//
//func TestNewConfigSigningMethodGet(t *testing.T) {
//	t.Parallel()
//
//	cobPath, err := backend.CurrencyOptionsBase.ToPath(0, 0)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	cr := cfgmock.NewService(
//		cfgmock.WithPV(cfgmock.PathValue{
//			cobPath.Bind(scope.WebsiteID, 1).String(): "EUR",
//			cobPath.Bind(scope.WebsiteID, 2).String(): "WIR", // Special Swiss currency
//		}),
//	)
//
//	// scope of CurrencyOptionsBase set to website, so no store config values are possible
//	cur, err := backend.CurrencyOptionsBase.Get(cr.NewScoped(1, 1))
//	if err != nil {
//		t.Fatal(err, cur)
//	}
//
//	assert.Exactly(t, ctxjwt.MustNewCurrencyISO("EUR"), cur)
//
//	cur, err = backend.CurrencyOptionsBase.Get(cr.NewScoped(2, 1))
//	assert.EqualError(t, err, "currency: tag is not a recognized currency")
//	assert.Exactly(t, ctxjwt.Currency{}, cur)
//}
//
