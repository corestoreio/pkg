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

package config

import (
	"testing"

	"github.com/corestoreio/csfw/config/path"
	"github.com/stretchr/testify/assert"
)

func TestSimpleStorage(t *testing.T) {
	sp := newSimpleStorage()

	p1 := path.MustNewByParts("aa/bb/cc")

	assert.NoError(t, sp.Set(p1, 19.99))
	f, err := sp.Get(p1)
	assert.NoError(t, err)
	assert.Exactly(t, 19.99, f.(float64))

	assert.NoError(t, sp.Set(p1, 4711))
	i, err := sp.Get(p1)
	assert.NoError(t, err)
	assert.Exactly(t, 4711, i.(int))

	ni, err := sp.Get(path.Path{})
	assert.NoError(t, err)
	assert.Nil(t, ni)

	keys, err := sp.AllKeys()
	assert.NoError(t, err)
	assert.Exactly(t, []string{"k1", "k2"}, keys)
}
