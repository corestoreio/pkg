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

package cfgmodel_test

import (
	"fmt"
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

var _ cfgmodel.MapIntResolver = (*mockMapIS)(nil)

type mockMapIS struct {
	string
	error
}

func (m mockMapIS) IntToStr(sg config.Scoped, i int) (string, error) {
	if m.string == "" {
		return fmt.Sprintf("Parent: %s => Current: %s => Value: %d", sg.ParentID(), sg.ScopeID(), i), m.error
	}
	return m.string, m.error
}

func TestNewMapIntStr(t *testing.T) {
	m := cfgmodel.NewMapIntStr(`general/store_information/region_id`, cfgmodel.WithScopeStore())
	m.MapIntResolver = mockMapIS{}

	s := cfgmock.NewService(cfgmock.PathValue{
		m.MustFQStore(5): 33,
	})

	val, err := m.Get(s.NewScoped(2, 5))
	assert.NoError(t, err, "%+v", err)
	assert.Exactly(t, scope.TypeIDs{scope.Store.Pack(5)}, s.IntInvokes().TypeIDs())
	assert.Exactly(t, `Parent: Type(Website) ID(2) => Current: Type(Store) ID(5) => Value: 33`, val)
}
