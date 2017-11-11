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

package catconfig

import (
	"github.com/corestoreio/cspkg/config"
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/config/source"
	"github.com/corestoreio/cspkg/store/scope"
	"github.com/juju/errors"
)

const (
	// PriceScopeGlobal prices are for all stores and websites the same.
	PriceScopeGlobal int = 0
	// PriceScopeWebsite prices are in each website different.
	PriceScopeWebsite int = 1
)

// PriceScope takes of the correct values for the price scope.
// A price scope can only be global or website based.
type PriceScope struct {
	cfgmodel.Int
}

// NewPriceScope defines the base currency scope
// ("Currency Setup" > "Currency Options" > "Base Currency").
// can be 0 = Global or 1 = Website
// See constants PriceScopeGlobal and PriceScopeWebsite.
func NewPriceScope(path string, opts ...cfgmodel.Option) PriceScope {
	return PriceScope{
		Int: cfgmodel.NewInt(path, append(opts, cfgmodel.WithSourceByInt(source.Ints{
			0: {PriceScopeGlobal, "Global Scope"},
			1: {PriceScopeWebsite, "Website Scope"},
		}))...),
	}

	//<source_model>Magento\Catalog\Model\Config\Source\Price\Scope</source_model>
}

// IsGlobal true if global scope for base currency
func (p PriceScope) IsGlobal(sg config.Scoped) (bool, error) {
	g, err := p.Get(sg)
	if err != nil {
		return false, errors.Mask(err)
	}
	return g == PriceScopeGlobal, nil
}

func (p PriceScope) Write(w config.Writer, v int, s scope.Type, id int64, idx interface {
	Invalidate()
}) error {

	if err := p.ValidateInt(v); err != nil {
		return err
	}

	if err := p.Int.Write(w, v, s, id); err != nil {
		return errors.Mask(err)
	}

	if idx != nil {
		idx.Invalidate()
	}

	// invalid price indexer and fully reindex
	//<backend_model>Magento\Catalog\Model\Indexer\Product\Price\System\Config\PriceScope</backend_model>

	return nil
}
