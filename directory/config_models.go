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
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/juju/errgo"
	"golang.org/x/text/currency"
)

// ConfigCurrency currency type for the configuration based on text/currency pkg.
type ConfigCurrency struct {
	model.Str
}

// NewConfigCurrency creates a new currency configuration type.
func NewConfigCurrency(path string, opts ...model.Option) ConfigCurrency {
	return ConfigCurrency{
		Str: model.NewStr(path, opts...),
	}
}

// Get tries to retrieve a currency
func (p ConfigCurrency) Get(sg config.ScopedGetter) (Currency, error) {
	cur := p.Str.Get(sg)
	u, err := currency.ParseISO(cur)
	if err != nil {
		return Currency{}, errgo.Mask(err)
	}
	return Currency{Unit: u}, nil
}

// Writes a currency to the configuration storage.
func (p ConfigCurrency) Write(w config.Writer, v Currency, s scope.Scope, id int64) error {
	return p.Str.Write(w, v.String(), s, id)
}
