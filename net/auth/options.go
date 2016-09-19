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

package auth

import (
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// WithDefaultConfig applies the default configuration settings based for
// a specific scope. This function overwrites any previous set options.
//
// Default values are:
//		- ?
func WithDefaultConfig(h scope.Hash) Option {
	return func(s *Service) error {
		var err error
		if h == scope.DefaultHash {
			s.defaultScopeCache, err = defaultScopedConfig()
			return errors.Wrap(err, "[mwauth] Default Scope with Default Config")
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		s.scopeCache[h], err = defaultScopedConfig()
		return errors.Wrapf(err, "[mwauth] Scope %s with Default Config", h)
	}
}

// WithIsActive enables or disables the authentication for a specific scope.
func WithIsActive(h scope.Hash, active bool) Option {
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.enable = active
			return nil
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.enable = active

		if sc, ok := s.scopeCache[h]; ok {
			sc.enable = scNew.enable
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}
