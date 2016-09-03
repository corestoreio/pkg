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

package signed

import (
	"hash"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/hashpool"
)

// WithDefaultConfig applies the default signed configuration settings based
// for a specific scope. This function overwrites any previous set options.
//
// Default values are:
//		- Denied Handler: http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
//		- VaryByer: returns an empty key
// Example:
//		s := MustNewService(WithDefaultConfig(scope.Store,1), WithVaryBy(scope.Store, 1, myVB))
func WithDefaultConfig(scp scope.Scope, id int64) Option {
	return withDefaultConfig(scp, id)
}

// WithVaryBy allows to set a custom key producer. VaryByer is called for each
// request to generate a key for the limiter. If it is nil, the middleware
// panics. The default VaryByer returns an empty string so that all requests
// uses the same key. VaryByer must be thread safe.
func WithHash64(scp scope.Scope, id int64, hf func() hash.Hash64) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.hashPool = hashpool.New64(hf)
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithDisable allows to disable a rate limit or enable it if set to false.
func WithDisable(scp scope.Scope, id int64, isDisabled bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.Disabled = isDisabled
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}
