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

package config_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ config.Storager = config.NewInMemoryStore()

func TestSimpleStorage(t *testing.T) {

	sp := config.NewInMemoryStore()

	p1 := cfgpath.MustNewByParts("aa/bb/cc")

	assert.NoError(t, sp.Set(p1, 19.99))
	f, err := sp.Get(p1)
	assert.NoError(t, err)
	assert.Exactly(t, 19.99, f.(float64))

	p2 := cfgpath.MustNewByParts("xx/yy/zz").BindStore(2)

	assert.NoError(t, sp.Set(p2, 4711))
	i, err := sp.Get(p2)
	assert.NoError(t, err)
	assert.Exactly(t, 4711, i.(int))

	ni, err := sp.Get(cfgpath.Path{})
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
	assert.Nil(t, ni)

	keys, err := sp.AllKeys()
	assert.NoError(t, err)
	keys.Sort()

	wantKeys := cfgpath.PathSlice{
		cfgpath.Path{Route: cfgpath.NewRoute(`aa/bb/cc`), ScopeID: scope.DefaultTypeID},
		cfgpath.Path{Route: cfgpath.NewRoute(`xx/yy/zz`), ScopeID: scope.MakeTypeID(scope.Store, 2)},
	}
	assert.Exactly(t, wantKeys, keys)

	p3 := cfgpath.MustNewByParts("rr/ss/tt").BindStore(1)
	ni, err = sp.Get(p3)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
	assert.Nil(t, ni)
}
