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

package catconfig

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/juju/errgo"
)

type ConfigPriceScope struct {
	model.Int
}

func NewConfigPriceScope(path string) ConfigPriceScope {
	return ConfigPriceScope{
		Int: model.NewInt(path),
	}

	//<source_model>Magento\Catalog\Model\Config\Source\Price\Scope</source_model>
}

func (p ConfigPriceScope) Write(w config.Writer, v int, s scope.Scope, id int64, idx interface {
	Invalidate()
}) error {

	if err := p.Int.Write(w, v, s, id); err != nil {
		return errgo.Mask(err)
	}

	idx.Invalidate()

	// invalid price indexer and fully reindex
	//<backend_model>Magento\Catalog\Model\Indexer\Product\Price\System\Config\PriceScope</backend_model>

	return nil
}
