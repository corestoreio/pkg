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

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/assert"
)

type flusher interface {
	Flush() error
}

func validateFoundGet(t *testing.T, s config.Storager, scp scope.TypeID, route, want string) {
	p := config.MustNewPathWithScope(scp, route)
	data, ok, err := s.Get(p)
	assert.NoError(t, err)
	assert.True(t, ok, "env value must be found")
	assert.Exactly(t, []byte(want), data)
}

func validateNotFoundGet(t *testing.T, s config.Storager, scp scope.TypeID, route string) {
	p := config.MustNewPathWithScope(scp, route)
	data, ok, err := s.Get(p)
	assert.NoError(t, err)
	assert.False(t, ok, "env value must NOT be found")
	assert.Nil(t, data, "Data must be nil")
}

func TestWithLoadStrings(t *testing.T) {
	t.Parallel()

	t.Run("unbalanced", func(t *testing.T) {
		cfgSrv, err := config.NewService(
			nil, config.Options{
				Level1: storage.NewMap(),
			},
			storage.WithLoadStrings("Baaam"),
		)
		assert.Nil(t, cfgSrv)
		assert.True(t, errors.NotAcceptable.Match(err), "%+v", err)
	})

	t.Run("successful level 1", func(t *testing.T) {
		pUserName := config.MustNewPath("payment/stripe/user_name").BindWebsite(2)

		cfgSrv, err := config.NewService(
			nil, config.Options{
				Level1: storage.NewMap(),
			},
			storage.WithLoadStrings(pUserName.String(), "alphZ"),
		)
		assert.NoError(t, err)
		assert.Exactly(t, "\"alph\\uf8ffZ\"", cfgSrv.Get(pUserName).String())
	})
	t.Run("successful level 1+2", func(t *testing.T) {
		pUserName := config.MustNewPath("payment/stripe/user_name").BindWebsite(2)

		cfgSrv, err := config.NewService(
			storage.NewMap(), config.Options{
				Level1: storage.NewMap(),
			},
			storage.WithLoadStrings(pUserName.String(), "alphX").WithUseStorageLevel(2),
			storage.WithLoadStrings(pUserName.String(), "alphZ"),
		)
		assert.NoError(t, err)
		assert.Exactly(t, "\"alph\\uf8ffZ\"", cfgSrv.Get(pUserName).String())
	})
	t.Run("successful sort order", func(t *testing.T) {
		pUserName := config.MustNewPath("payment/stripe/user_name").BindWebsite(2)

		cfgSrv, err := config.NewService(
			storage.NewMap(), config.Options{},
			storage.WithLoadStrings(pUserName.String(), "alphX").WithUseStorageLevel(2).WithSortOrder(2),
			storage.WithLoadStrings(pUserName.String(), "alphZ").WithUseStorageLevel(2).WithSortOrder(-1),
		)
		assert.NoError(t, err)
		assert.Exactly(t, "\"alph\\uf8ffX\"", cfgSrv.Get(pUserName).String())
	})
}
