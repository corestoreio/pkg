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

package scope_test

import (
	"testing"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestPermAll(t *testing.T) {
	var p scope.Perm
	pa := p.All()
	assert.True(t, pa.Has(scope.DefaultID))
	assert.True(t, pa.Has(scope.WebsiteID))
	assert.True(t, pa.Has(scope.StoreID))
}

func TestNewPerm(t *testing.T) {
	p := scope.NewPerm(scope.WebsiteID)

	assert.False(t, p.Has(scope.DefaultID))
	assert.True(t, p.Has(scope.WebsiteID))
	assert.False(t, p.Has(scope.StoreID))
}

func TestPermMarshalJSONAll(t *testing.T) {
	var p scope.Perm
	pa := p.All()
	jd, err := pa.MarshalJSON()
	assert.NoError(t, err)
	assert.Exactly(t, "[\"Default\",\"Website\",\"Store\"]", string(jd))
}

func TestPermMarshalJSONNull(t *testing.T) {
	var p scope.Perm
	jd, err := p.MarshalJSON()
	assert.NoError(t, err)
	assert.Exactly(t, "null", string(jd))
}
