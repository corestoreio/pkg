// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"sort"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/stretchr/testify/assert"
)

var _ config.Storager = config.NewInMemoryStore()

func TestSimpleStorage(t *testing.T) {
	t.Parallel()

	const path = "aa/bb/cc"
	sp := config.NewInMemoryStore()

	assert.NoError(t, sp.Set(0, "aa/bb/cc", []byte(`19.99`)))
	vb, ok, err := sp.Value(0, path)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Exactly(t, []byte(`19.99`), vb)

	ni, ok, err := sp.Value(0, "")
	assert.True(t, errors.NotValid.Match(err), "Error: %s", err)
	assert.Nil(t, ni)

	scps, paths, err := sp.AllKeys()
	assert.NoError(t, err)
	sort.Strings(paths)

	assert.Exactly(t, "x TODO xx", scps.String())

	wantKeys := config.PathSlice{
		config.MustMakePath(`aa/bb/cc`),
		config.MustMakePath(`xx/yy/zz`).BindStore(2),
	}
	assert.Exactly(t, wantKeys, paths)

	ni, ok, err = sp.Value(0, "rr/ss/tt")
	assert.True(t, errors.NotFound.Match(err), "Error: %s", err)
	assert.False(t, ok)
	assert.Nil(t, ni)
}
