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

package storage_test

import (
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type flusher interface {
	Flush() error
}

func validateFoundGet(t *testing.T, s config.Storager, scp scope.TypeID, route string, want string) {
	p := config.MustNewPathWithScope(scp, route)
	data, ok, err := s.Get(p)
	require.NoError(t, err)
	assert.True(t, ok, "env value must be found")
	assert.Exactly(t, []byte(want), data)
}

func validateNotFoundGet(t *testing.T, s config.Storager, scp scope.TypeID, route string) {
	p := config.MustNewPathWithScope(scp, route)
	data, ok, err := s.Get(p)
	require.NoError(t, err)
	assert.False(t, ok, "env value must NOT be found")
	assert.Nil(t, data, "Data must be nil")
}
