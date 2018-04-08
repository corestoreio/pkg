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

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ config.Storager = config.NewInMemoryStore()

func TestSimpleStorage_OneKey(t *testing.T) {
	t.Parallel()

	const path = "aa/bb/cc"
	var testTypeID = scope.Store.Pack(55)
	sp := config.NewInMemoryStore()

	assert.NoError(t, sp.Set(testTypeID, "aa/bb/cc", []byte(`19.99`)))
	vb, ok, err := sp.Value(testTypeID, path)
	require.NoError(t, err)
	require.True(t, ok)
	assert.True(t, ok)
	assert.Exactly(t, []byte(`19.99`), vb)

	ni, ok, err := sp.Value(testTypeID, "")
	require.NoError(t, err, "Error: %s", err)
	require.False(t, ok)
	assert.Nil(t, ni)

	scps, paths, err := sp.AllKeys()
	require.NoError(t, err)
	sort.Strings(paths)

	assert.Exactly(t, "Type(Store) ID(55)", scps.String())

	assert.Exactly(t, []string{"aa/bb/cc"}, paths)

	ni, ok, err = sp.Value(0, "rr/ss/tt")
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, ni)
}

func TestSimpleStorage_MultiKey(t *testing.T) {
	t.Parallel()

	sp := config.NewInMemoryStore()
	sp.Set(scope.Website.Pack(22), "aa/bb/cc", []byte(`path22`))
	sp.Set(scope.Store.Pack(33), "dd/ee/ff", []byte(`path33`))
	sp.Set(scope.DefaultTypeID, "gg/hh/ii", []byte(`path44`))

	scps, paths, err := sp.AllKeys()
	require.NoError(t, err)
	sort.Strings(paths)
	sort.Sort(scps)

	assert.Exactly(t, "Type(Default) ID(0); Type(Website) ID(22); Type(Store) ID(33)", scps.String())
	assert.Exactly(t, []string{"aa/bb/cc", "dd/ee/ff", "gg/hh/ii"}, paths)
}
