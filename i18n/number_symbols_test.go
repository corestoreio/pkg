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

package i18n_test

import (
	"testing"

	"github.com/corestoreio/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

func TestSymbolsString(t *testing.T) {
	assert.Equal(t,
		"Decimal\t\t\t\t\t.\nGroup\t\t\t\t\t,\nList\t\t\t\t\t;\nPercentSign\t\t\t\t%\nCurrencySign\t\t\t¤\nPlusSign\t\t\t\t+\nMinusSign\t\t\t\t—\nExponential\t\t\t\tE\nSuperscriptingExponent\t×\nPerMille\t\t\t\t‰\nInfinity\t\t\t\t∞\nNaN\t\t\t\t\t\tNaN\n",
		i18n.NewSymbols().String())
}

func TestSymbolsMerge(t *testing.T) {
	to := i18n.NewSymbols()
	from := i18n.Symbols{
		Decimal:                '',
		Group:                  '',
		List:                   '',
		PercentSign:            '',
		CurrencySign:           '',
		PlusSign:               '',
		MinusSign:              '',
		Exponential:            '',
		SuperscriptingExponent: '',
		PerMille:               '',
		Infinity:               '',
		Nan:                    []byte(`Pear`),
	}
	to.Merge(from)
	assert.EqualValues(t, from, to)
	to.Merge(i18n.Symbols{})
	assert.EqualValues(t, from, to)
}
