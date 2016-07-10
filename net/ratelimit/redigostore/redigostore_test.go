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

package redigostore_test

import (
	"testing"

	"github.com/corestoreio/csfw/net/ratelimit"
	"github.com/corestoreio/csfw/net/ratelimit/redigostore"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestWithGCRARedis(t *testing.T) {
	s4 := scope.NewHash(scope.Store, 4)

	t.Run("CalcErrorRedis", func(t *testing.T) {
		s, err := ratelimit.New(redigostore.WithGCRA(scope.Store, 4, "redis://localhost/", 's', 100, 10))
		assert.Nil(t, s)
		assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	})

	t.Run("Ok", func(t *testing.T) {
		s := ratelimit.MustNew(
			ratelimit.WithDefaultConfig(scope.Store, 4),
			redigostore.WithGCRA(scope.Store, 4, "redis://localhost/1", 's', 100, 10),
			redigostore.WithGCRA(scope.Default, 0, "redis://localhost/2", 's', 100, 10),
		)
		assert.NotNil(t, s.scopeCache[s4].RateLimiter, "Scope Website")
		assert.NotNil(t, s.scopeCache[scope.DefaultHash].RateLimiter, "Scope Default")
	})

	t.Run("OverwrittenByWithDefaultConfig", func(t *testing.T) {
		s := ratelimit.MustNew(
			redigostore.WithGCRA(scope.Store, 4, "redis://localhost/1", 's', 100, 10),
			ratelimit.WithDefaultConfig(scope.Store, 4),
		)
		assert.Nil(t, s.scopeCache[s4].RateLimiter)
		err := s.getConfigByHash(s4, 0).isValid()
		assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
	})

}
