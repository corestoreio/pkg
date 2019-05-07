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

package memstore_test

import (
	"testing"

	"github.com/corestoreio/pkg/config/cfgmock"
	"github.com/corestoreio/pkg/net/ratelimit"
	"github.com/corestoreio/pkg/net/ratelimit/backendratelimit"
	"github.com/corestoreio/pkg/net/ratelimit/memstore"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

func TestWithGCRAMemStore(t *testing.T) {
	s4 := scope.MakeTypeID(scope.Store, 4)

	t.Run("CalcErrorRate", func(t *testing.T) {
		s, err := ratelimit.New(memstore.WithGCRA(3333, 's', 100, -1, scope.Store.WithID(4)))
		assert.Nil(t, s)
		assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	})

	t.Run("Ok", func(t *testing.T) {
		s := ratelimit.MustNew(
			ratelimit.WithDefaultConfig(scope.Store.WithID(4)),
			memstore.WithGCRA(3333, 's', 100, 10, scope.Store.WithID(4)),
			memstore.WithGCRA(2222, 's', 100, 20, scope.DefaultTypeID),
		)
		cfg, err := s.ConfigByScopeID(s4, 0)
		assert.NoError(t, err, "%+v", err)
		assert.NotNil(t, cfg.RateLimiter, "Scope Website")
		cfg, err = s.ConfigByScopeID(scope.DefaultTypeID, 0)
		assert.NoError(t, err, "%+v", err)
		assert.NotNil(t, cfg.RateLimiter, "Scope Default")
	})

	t.Run("OverwrittenByWithDefaultConfig", func(t *testing.T) {
		s := ratelimit.MustNew(
			memstore.WithGCRA(1111, 's', 100, 10, scope.Store.WithID(4)),
			ratelimit.WithDefaultConfig(scope.Store.WithID(4)),
		)
		cfg, err := s.ConfigByScopeID(s4, 0)
		assert.True(t, errors.IsNotValid(err), "%+v", err)
		assert.Nil(t, cfg.RateLimiter)
		_, err = s.ConfigByScopeID(s4, 0)
		assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	})
}

func TestBackend_Path_Errors(t *testing.T) {

	cfgStruct, err := backendratelimit.NewConfigStructure()
	if err != nil {
		panic(err)
	}
	backend := backendratelimit.New(cfgStruct)

	tests := []struct {
		cfgPath string
		val     interface{}
		errBhf  errors.BehaviourFunc
	}{
		{backend.Burst.MustFQWebsite(2), struct{}{}, errors.IsNotValid},
		{backend.Requests.MustFQWebsite(2), struct{}{}, errors.IsNotValid},
		{backend.Duration.MustFQWebsite(2), "[a-z+", errors.IsFatal},
		{backend.Duration.MustFQWebsite(2), struct{}{}, errors.IsNotValid},
		{backend.StorageGCRAMaxMemoryKeys.MustFQWebsite(2), struct{}{}, errors.IsNotValid},
		{backend.StorageGCRAMaxMemoryKeys.MustFQWebsite(2), 0, errors.IsEmpty},
	}
	for i, test := range tests {

		name, scpFnc := memstore.NewOptionFactory(backend.Burst, backend.Requests, backend.Duration, backend.StorageGCRAMaxMemoryKeys)
		if have, want := name, memstore.OptionName; have != want {
			t.Errorf("Have: %v Want: %v", have, want)
		}

		cfgSrv := cfgmock.NewService(cfgmock.PathValue{
			test.cfgPath: test.val,
		})
		cfgScp := cfgSrv.NewScoped(2, 0)

		_, err := ratelimit.New(scpFnc(cfgScp)...)
		assert.True(t, test.errBhf(err), "Index %d Error: %+v", i, err)
	}
}
