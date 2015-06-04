// Copyright 2015 CoreStore Authors
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

package config_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/utils"
	"github.com/stretchr/testify/assert"
)

func TestValueLabelSlice(t *testing.T) {

	tests := []struct {
		have      config.ValueLabelSlice
		wantValue string
		wantLabel string
		order     utils.SortDirection
	}{
		{
			config.ValueLabelSlice{config.ValueLabel{"kb", "l2"}, config.ValueLabel{"ka", "l1"}, config.ValueLabel{"kc", "l3"}, config.ValueLabel{"kY", "l5"}, config.ValueLabel{"k0", "l4"}},
			`[{"Value":"k0","Label":"l4"},{"Value":"kY","Label":"l5"},{"Value":"ka","Label":"l1"},{"Value":"kb","Label":"l2"},{"Value":"kc","Label":"l3"}]` + "\n",
			`[{"Value":"ka","Label":"l1"},{"Value":"kb","Label":"l2"},{"Value":"kc","Label":"l3"},{"Value":"k0","Label":"l4"},{"Value":"kY","Label":"l5"}]` + "\n",
			utils.SortAsc,
		},
		{
			config.ValueLabelSlice{config.ValueLabel{"x3", "l2"}, config.ValueLabel{"xg", "l1"}, config.ValueLabel{"xK", "l3"}, config.ValueLabel{"x0", "l5"}, config.ValueLabel{"x-", "l4"}},
			`[{"Value":"xg","Label":"l1"},{"Value":"xK","Label":"l3"},{"Value":"x3","Label":"l2"},{"Value":"x0","Label":"l5"},{"Value":"x-","Label":"l4"}]` + "\n",
			`[{"Value":"x0","Label":"l5"},{"Value":"x-","Label":"l4"},{"Value":"xK","Label":"l3"},{"Value":"x3","Label":"l2"},{"Value":"xg","Label":"l1"}]` + "\n",
			utils.SortDesc,
		},
	}

	for i, test := range tests {
		test.have.SortByValue(test.order)
		assert.Exactly(t, test.wantValue, test.have.ToJSON(), "SortByValue Index %d", i)
		test.have.SortByLabel(test.order)
		assert.Exactly(t, test.wantLabel, test.have.ToJSON(), "SortByLabel Index %d", i)
	}
}
