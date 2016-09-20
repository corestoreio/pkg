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

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/ratelimit"
	"github.com/corestoreio/csfw/net/ratelimit/backendratelimit"
	"github.com/corestoreio/csfw/net/ratelimit/memstore"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestWithGCRAMemStore(t *testing.T) {
	s4 := scope.NewHash(scope.Store, 4)

	t.Run("CalcErrorRate", func(t *testing.T) {
		s, err := ratelimit.New(memstore.WithGCRA(scope.Store.ToHash(4), 3333, 's', 100, -1))
		assert.Nil(t, s)
		assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	})

	t.Run("Ok", func(t *testing.T) {
		s := ratelimit.MustNew(
			ratelimit.WithDefaultConfig(scope.Store.ToHash(4)),
			memstore.WithGCRA(scope.Store.ToHash(4), 3333, 's', 100, 10),
			memstore.WithGCRA(scope.DefaultHash, 2222, 's', 100, 20),
		)
		assert.NotNil(t, s.ConfigByScopeHash(s4, 0).RateLimiter, "Scope Website")
		assert.NotNil(t, s.ConfigByScopeHash(scope.DefaultHash, 0).RateLimiter, "Scope Default")
	})

	t.Run("OverwrittenByWithDefaultConfig", func(t *testing.T) {
		s := ratelimit.MustNew(
			memstore.WithGCRA(scope.Store.ToHash(4), 1111, 's', 100, 10),
			ratelimit.WithDefaultConfig(scope.Store.ToHash(4)),
		)
		assert.Nil(t, s.ConfigByScopeHash(s4, 0).RateLimiter)
		err := s.ConfigByScopeHash(s4, 0).IsValid()
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
		toPath func(s scope.Scope, scopeID int64) string
		val    interface{}
		errBhf errors.BehaviourFunc
	}{
		{backend.Burst.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.Requests.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.Duration.MustFQ, "[a-z+", errors.IsFatal},
		{backend.Duration.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.StorageGCRAMaxMemoryKeys.MustFQ, struct{}{}, errors.IsNotValid},
		{backend.StorageGCRAMaxMemoryKeys.MustFQ, 0, errors.IsEmpty},
	}
	for i, test := range tests {

		name, scpFnc := memstore.NewOptionFactory(backend.Burst, backend.Requests, backend.Duration, backend.StorageGCRAMaxMemoryKeys)
		if have, want := name, memstore.OptionName; have != want {
			t.Errorf("Have: %v Want: %v", have, want)
		}

		cfgSrv := cfgmock.NewService(cfgmock.PathValue{
			test.toPath(scope.Website, 2): test.val,
		})
		cfgScp := cfgSrv.NewScoped(2, 0)

		_, err := ratelimit.New(scpFnc(cfgScp)...)
		assert.True(t, test.errBhf(err), "Index %d Error: %+v", i, err)
	}
}
